package ots

import (
	"fmt"
	"strconv"
	"time"

	otsTunnel "github.com/aliyun/aliyun-tablestore-go-sdk/tunnel"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tunnel/protocol"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// resourceAlicloudOtsTunnel Tablestore tunnel not support update
func ResourceOtsTunnel() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunOtsTunnelCreate,
		Read:   resourceAliyunOtsTunnelRead,
		Delete: resourceAliyunOtsTunnelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateOTSInstanceName,
			},
			"table_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateOTSTableName,
			},
			"tunnel_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateOTSTunnelName,
			},
			"tunnel_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(BaseAndStreamTunnel), string(BaseDataTunnel), string(StreamTunnel)}, false),
			},
			// start_offset/end_offset/start_offset_flag only take effect when tunnel_type is Stream.
			// They are ignored (and rejected at create time) for BaseData and BaseAndStream tunnels.
			// start_offset/end_offset are millisecond timestamps (uint64, ~1.73e12) that exceed
			// 2^31; this provider ships 32-bit builds (386/arm) where Go's int is 32-bit, so
			// TypeInt would overflow and crash the plugin. TypeString parses to uint64 in
			// buildOtsStreamTunnelConfig and is re-formatted on Read.
			"start_offset": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateOtsTunnelOffset,
			},
			"end_offset": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateOtsTunnelOffset,
			},
			"start_offset_flag": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{string(StartOffsetFlagLatest), string(StartOffsetFlagEarliest)}, false),
				// LATEST is the server default: DescribeTunnel reports Flag=0 (LATEST) when the
				// flag is unset, so Read only backfills EARLIEST. Suppress the empty-state vs
				// "LATEST"-config diff so an explicit start_offset_flag = "LATEST" does not
				// trigger a spurious ForceNew on every plan.
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					return old == "" && new == string(StartOffsetFlagLatest)
				},
			},
			"tunnel_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tunnel_rpo": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tunnel_stage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expired": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"channels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"channel_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"channel_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"channel_rpo": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliyunOtsTunnelCreate(d *schema.ResourceData, meta interface{}) error {
	instanceName := d.Get("instance_name").(string)
	tableName := d.Get("table_name").(string)
	tunnelName := d.Get("tunnel_name").(string)
	tunnelTypeS := d.Get("tunnel_type").(string)
	tunnelType, err := parseTunnelType(tunnelTypeS)
	if err != nil {
		return err
	}

	client := meta.(*connectivity.AliyunClient)
	otsService := NewService(client)
	// check table exists
	if err := retry.Retry(1*time.Minute, func() *retry.RetryError {
		_, e := otsService.DescribeOtsTable(fmt.Sprintf("%s%s%s", instanceName, COLON_SEPARATED, tableName))
		if e != nil {
			if errs.NotFoundError(e) {
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		return nil
	}); err != nil {
		return errs.WrapError(err)
	}

	request := new(otsTunnel.CreateTunnelRequest)
	request.TableName = tableName
	request.TunnelName = tunnelName
	request.Type = tunnelType

	// StreamTunnelConfig (start_offset/end_offset/start_offset_flag) only applies to
	// Stream tunnels. For BaseData and BaseAndStream it is not applicable and rejected
	// so that misconfiguration is not silently dropped.
	if tunnelType == otsTunnel.TunnelTypeStream {
		if streamCfg := buildOtsStreamTunnelConfig(d); streamCfg != nil {
			request.StreamTunnelConfig = streamCfg
		}
	} else if streamConfigFieldSet(d) {
		return errs.WrapError(errs.Error("start_offset, end_offset and start_offset_flag are only supported when tunnel_type is Stream"))
	}

	var requestInfo otsTunnel.TunnelClient
	if err := retry.Retry(1*time.Minute, func() *retry.RetryError {
		raw, err := client.WithTableStoreTunnelClient(instanceName, func(tunnelClient otsTunnel.TunnelClient) (interface{}, error) {
			requestInfo = tunnelClient
			return tunnelClient.CreateTunnel(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, errs.OtsTunnelIsTemporarilyUnavailable) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug("CreateTunnel", raw, requestInfo, request)
		return nil
	}); err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ots_tunnel", "CreateTunnel", errs.AliyunTablestoreGoSdk)
	}

	d.SetId(fmt.Sprintf("%s%s%s%s%s", instanceName, COLON_SEPARATED, tableName, COLON_SEPARATED, tunnelName))
	return resourceAliyunOtsTunnelRead(d, meta)
}

func resourceAliyunOtsTunnelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := NewService(client)
	object, err := otsService.DescribeOtsTunnel(d.Id())

	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("instance_name", object.Tunnel.InstanceName)
	d.Set("table_name", object.Tunnel.TableName)
	d.Set("tunnel_name", object.Tunnel.TunnelName)
	d.Set("tunnel_id", object.Tunnel.TunnelId)
	d.Set("tunnel_type", object.Tunnel.TunnelType)
	d.Set("tunnel_rpo", object.TunnelRPO)
	d.Set("tunnel_stage", object.Tunnel.Stage)
	d.Set("expired", object.Tunnel.Expired)
	d.Set("create_time", object.Tunnel.CreateTime.UnixNano())

	// Backfill the stream offsets from DescribeTunnel. start_offset_flag is backfilled
	// only when the server reports EARLIEST: the SDK returns the LATEST (0) default when
	// the flag is unset, so backfilling it unconditionally would introduce a spurious
	// ForceNew diff against configurations that omit the flag. EARLIEST is unambiguous
	// (non-default) and safe to backfill. The LATEST case is handled via the schema
	// DiffSuppressFunc on start_offset_flag.
	if streamCfg := object.Tunnel.StreamTunnelConfig; streamCfg != nil {
		// Only backfill offsets when non-zero: when start_offset_flag is set the server
		// reports StartOffset=0, and writing "0" into a TypeString field would diff
		// against an omitted configuration. The server stores offsets in nanoseconds
		// (Go SDK doc 176301); the schema exposes milliseconds (Java SDK doc 108059),
		// so convert ns -> ms on the way out.
		if streamCfg.StartOffset != 0 {
			d.Set("start_offset", strconv.FormatUint(streamCfg.StartOffset/1_000_000, 10))
		}
		if streamCfg.EndOffset != 0 {
			d.Set("end_offset", strconv.FormatUint(streamCfg.EndOffset/1_000_000, 10))
		}
		if streamCfg.Flag == protocol.StartOffsetFlag_EARLIEST {
			d.Set("start_offset_flag", string(StartOffsetFlagEarliest))
		}
	}

	channels := make([]map[string]interface{}, len(object.Channels))
	for i, channel := range object.Channels {
		item := make(map[string]interface{})
		item["channel_id"] = channel.ChannelId
		item["channel_type"] = channel.ChannelType
		item["channel_status"] = channel.ChannelStatus
		item["client_id"] = channel.ClientId
		item["channel_rpo"] = channel.ChannelRPO

		channels[i] = item
	}
	d.Set("channels", channels)

	return nil
}

func resourceAliyunOtsTunnelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := NewService(client)

	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	req := &otsTunnel.DeleteTunnelRequest{
		TableName:  parts[1],
		TunnelName: parts[2],
	}
	var requestInfo otsTunnel.TunnelClient
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		raw, err := client.WithTableStoreTunnelClient(parts[0], func(tunnelClient otsTunnel.TunnelClient) (interface{}, error) {
			requestInfo = tunnelClient
			return tunnelClient.DeleteTunnel(req)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, errs.OtsTunnelIsTemporarilyUnavailable) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug("DeleteTunnel", raw, requestInfo, req)
		return nil
	})
	if err != nil {
		if isOtsTunnelNotFound(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), "DeleteTunnel", errs.AliyunTablestoreGoSdk)
	}
	return errs.WrapError(otsService.WaitForOtsTunnel(d.Id(), Deleted, DefaultTimeout))
}

func parseTunnelType(orig string) (otsTunnel.TunnelType, error) {
	switch orig {
	case string(BaseAndStreamTunnel):
		return otsTunnel.TunnelTypeBaseStream, nil
	case string(BaseDataTunnel):
		return otsTunnel.TunnelTypeBaseData, nil
	case string(StreamTunnel):
		return otsTunnel.TunnelTypeStream, nil
	default:
		return "", errs.WrapError(errs.Error("unknown ots tunnel type: %s", orig))
	}
}

// streamConfigFieldSet reports whether any of the StreamTunnelConfig-related fields
// (start_offset, end_offset, start_offset_flag) is configured.
func streamConfigFieldSet(d *schema.ResourceData) bool {
	_, hasStart := d.GetOk("start_offset")
	_, hasEnd := d.GetOk("end_offset")
	_, hasFlag := d.GetOk("start_offset_flag")
	return hasStart || hasEnd || hasFlag
}

// buildOtsStreamTunnelConfig builds the SDK StreamTunnelConfig from the schema fields.
// It returns nil when none of the stream config fields are configured, so that an
// existing Stream tunnel without offset configuration keeps its default (LATEST) behavior.
//
// The tablestore SDK serializes StartOffset and Flag as mutually exclusive: when
// start_offset is non-zero only StartOffset is sent, otherwise Flag is sent. end_offset
// is sent independently when non-zero. Therefore no mutual-exclusion validation is
// required here; the SDK handles a simultaneously configured flag and offset by
// preferring the offset.
func buildOtsStreamTunnelConfig(d *schema.ResourceData) *otsTunnel.StreamTunnelConfig {
	if !streamConfigFieldSet(d) {
		return nil
	}
	cfg := &otsTunnel.StreamTunnelConfig{}
	// The schema exposes start_offset/end_offset as millisecond timestamps (Java
	// SDK doc 108059), but the tablestore Go SDK and the OTS server expect
	// nanoseconds (Go SDK doc 176301; confirmed empirically: millisecond values
	// are rejected by CreateTunnel as pointing to 1970). Convert ms -> ns here;
	// Read converts back to ms.
	if v, ok := d.GetOk("start_offset"); ok {
		if offset, err := strconv.ParseUint(v.(string), 10, 64); err == nil {
			cfg.StartOffset = offset * 1_000_000
		}
	}
	if v, ok := d.GetOk("end_offset"); ok {
		if offset, err := strconv.ParseUint(v.(string), 10, 64); err == nil {
			cfg.EndOffset = offset * 1_000_000
		}
	}
	if flag, ok := d.GetOk("start_offset_flag"); ok {
		switch StartOffsetFlagString(flag.(string)) {
		case StartOffsetFlagEarliest:
			cfg.Flag = protocol.StartOffsetFlag_EARLIEST
		case StartOffsetFlagLatest:
			cfg.Flag = protocol.StartOffsetFlag_LATEST
		}
	}
	return cfg
}

// validateOtsTunnelOffset validates that a start_offset/end_offset string is a
// non-negative integer millisecond timestamp. Non-numeric input would otherwise be
// silently dropped (parsed to 0) when building the StreamTunnelConfig.
func validateOtsTunnelOffset(val interface{}, key string) (warns []string, errs []error) {
	s, ok := val.(string)
	if !ok || s == "" {
		return
	}
	if _, err := strconv.ParseUint(s, 10, 64); err != nil {
		errs = append(errs, fmt.Errorf("%q must be a non-negative integer millisecond timestamp, got %q", key, s))
	}
	return
}
