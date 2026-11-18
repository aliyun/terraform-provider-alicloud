package ots

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	otsTunnel "github.com/aliyun/aliyun-tablestore-go-sdk/tunnel"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tunnel/protocol"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/go-cty/cty"
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
		CustomizeDiff: resourceAliyunOtsTunnelCustomizeDiff,
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
			// start_offset/end_offset/start_offset_flag are only supported when tunnel_type
			// is Stream; configuring them for BaseData or BaseAndStream tunnels is rejected
			// at create time. start_offset/end_offset are millisecond timestamps (uint64,
			// ~1.73e12) that exceed 2^31; this provider ships 32-bit builds (386/arm) where
			// Go's int is 32-bit, so TypeInt would overflow and crash the plugin. TypeString
			// parses to uint64 in buildOtsStreamTunnelConfig and is re-formatted on Read.
			"start_offset": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validateOtsTunnelOffset,
				ConflictsWith: []string{"start_offset_flag"},
			},
			"end_offset": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateOtsTunnelOffset,
			},
			"start_offset_flag": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{string(StartOffsetFlagLatest), string(StartOffsetFlagEarliest)}, false),
				ConflictsWith: []string{"start_offset"},
				// LATEST is the server default. Read reports the effective flag, so an
				// omitted configuration and an explicit "LATEST" are equivalent.
				DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
					// SDKv2 drops suppressed entries from the apply diff, so suppressing on
					// create would hide an explicitly configured value from Create (and from
					// the Stream-only rejection). Never suppress while there is no state.
					// For an existing non-Stream tunnel the flag is invalid too, so let the
					// diff through and route the misconfiguration into the same rejection
					// instead of silently ignoring it.
					if d.Id() == "" || d.Get("tunnel_type").(string) != string(StreamTunnel) {
						return false
					}
					// Treat omitted and explicit LATEST as the same value in both directions
					// so that adding or removing the explicit default does not replace the
					// tunnel; EARLIEST still diffs and replaces as expected.
					return (old == "" && new == string(StartOffsetFlagLatest)) ||
						(old == string(StartOffsetFlagLatest) && new == "")
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
		streamCfg, err := buildOtsStreamTunnelConfig(d)
		if err != nil {
			return errs.WrapError(err)
		}
		if streamCfg != nil {
			request.StreamTunnelConfig = streamCfg
		}
	} else if streamConfigFieldSet(d.GetRawConfig()) {
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

	// Backfill the Stream consumption settings from DescribeTunnel. The three
	// arguments only apply to Stream tunnels. The server stores offsets in
	// nanoseconds (Go SDK doc 176301) while the schema exposes milliseconds
	// (Java SDK doc 108059), so convert ns -> ms on the way out.
	if object.Tunnel.TunnelType == string(otsTunnel.TunnelTypeStream) {
		if streamCfg := object.Tunnel.StreamTunnelConfig; streamCfg != nil {
			if streamCfg.StartOffset != 0 {
				d.Set("start_offset", strconv.FormatUint(streamCfg.StartOffset/1_000_000, 10))
			}
			if streamCfg.EndOffset != 0 {
				d.Set("end_offset", strconv.FormatUint(streamCfg.EndOffset/1_000_000, 10))
			}
			if streamCfg.StartOffset == 0 && streamCfg.Flag == protocol.StartOffsetFlag_EARLIEST {
				d.Set("start_offset_flag", string(StartOffsetFlagEarliest))
			} else {
				d.Set("start_offset_flag", string(StartOffsetFlagLatest))
			}
		} else {
			d.Set("start_offset_flag", string(StartOffsetFlagLatest))
		}
	} else {
		d.Set("start_offset_flag", "")
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

func resourceAliyunOtsTunnelCustomizeDiff(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	if !diff.NewValueKnown("tunnel_type") {
		return nil
	}
	tunnelType, ok := diff.GetOk("tunnel_type")
	if !ok {
		return nil
	}
	if tunnelType.(string) != string(StreamTunnel) {
		if streamConfigFieldSet(diff.GetRawConfig()) {
			return errs.WrapError(errs.Error("start_offset, end_offset and start_offset_flag are only supported when tunnel_type is Stream"))
		}
		return nil
	}
	rawConfig := diff.GetRawConfig()
	if diff.Id() != "" && !rawConfig.IsNull() && rawConfig.Type().IsObjectType() && rawConfig.GetAttr("start_offset_flag").IsNull() {
		oldFlag, _ := diff.GetChange("start_offset_flag")
		if oldFlag == string(StartOffsetFlagEarliest) {
			if err := diff.SetNewComputed("start_offset_flag"); err != nil {
				return err
			}
		}
	}
	if !diff.NewValueKnown("start_offset") || !diff.NewValueKnown("end_offset") {
		return nil
	}
	start, startSet := diff.GetOk("start_offset")
	end, endSet := diff.GetOk("end_offset")
	if !startSet || !endSet {
		return nil
	}
	startOffset, err := parseOtsTunnelOffset("start_offset", start.(string))
	if err != nil {
		return err
	}
	endOffset, err := parseOtsTunnelOffset("end_offset", end.(string))
	if err != nil {
		return err
	}
	if endOffset <= startOffset {
		return errs.WrapError(errs.Error("end_offset must be greater than start_offset when both are configured"))
	}
	return nil
}

func streamConfigFieldSet(raw cty.Value) bool {
	if raw.IsNull() || !raw.Type().IsObjectType() {
		return false
	}
	for _, key := range []string{"start_offset", "end_offset", "start_offset_flag"} {
		if v := raw.GetAttr(key); !v.IsNull() {
			return true
		}
	}
	return false
}

// buildOtsStreamTunnelConfig builds the SDK StreamTunnelConfig from the schema fields.
// It returns nil when none of the stream config fields are configured, so that an
// existing Stream tunnel without offset configuration keeps its default (LATEST) behavior.
//
// The tablestore SDK serializes StartOffset and Flag as mutually exclusive: when
// start_offset is non-zero only StartOffset is sent, otherwise Flag is sent. end_offset
// is sent independently when non-zero. The schema also declares start_offset and
// start_offset_flag mutually exclusive, so both are never configured together.
func buildOtsStreamTunnelConfig(d *schema.ResourceData) (*otsTunnel.StreamTunnelConfig, error) {
	if !streamConfigFieldSet(d.GetRawConfig()) {
		return nil, nil
	}
	cfg := &otsTunnel.StreamTunnelConfig{}
	if v, ok := d.GetOk("start_offset"); ok {
		offset, err := parseOtsTunnelOffset("start_offset", v.(string))
		if err != nil {
			return nil, err
		}
		cfg.StartOffset = offset
	}
	if v, ok := d.GetOk("end_offset"); ok {
		offset, err := parseOtsTunnelOffset("end_offset", v.(string))
		if err != nil {
			return nil, err
		}
		cfg.EndOffset = offset
	}
	if cfg.StartOffset != 0 && cfg.EndOffset != 0 && cfg.EndOffset <= cfg.StartOffset {
		return nil, errs.WrapError(errs.Error("end_offset must be greater than start_offset when both are configured"))
	}
	if flag, ok := d.GetOk("start_offset_flag"); ok {
		switch StartOffsetFlagString(flag.(string)) {
		case StartOffsetFlagEarliest:
			cfg.Flag = protocol.StartOffsetFlag_EARLIEST
		case StartOffsetFlagLatest:
			cfg.Flag = protocol.StartOffsetFlag_LATEST
		}
	}
	return cfg, nil
}

// validateOtsTunnelOffset validates that a start_offset/end_offset string is a
// positive, canonical decimal millisecond timestamp that can be converted to
// nanoseconds without overflowing. Non-canonical input (e.g. leading zeros)
// would be reformatted by Read and, because the attribute is ForceNew, would
// otherwise replace the tunnel on every plan.
func validateOtsTunnelOffset(val interface{}, key string) (warns []string, es []error) {
	s, ok := val.(string)
	if !ok {
		return
	}
	if _, err := parseOtsTunnelOffset(key, s); err != nil {
		es = append(es, err)
	}
	return
}

func parseOtsTunnelOffset(key, value string) (uint64, error) {
	if value == "" {
		return 0, fmt.Errorf("%q must not be empty", key)
	}
	offset, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%q must be a non-negative integer millisecond timestamp, got %q", key, value)
	}
	if offset == 0 {
		return 0, fmt.Errorf("%q must be a positive millisecond timestamp", key)
	}
	if strconv.FormatUint(offset, 10) != value {
		return 0, fmt.Errorf("%q must be a canonical decimal integer without leading zeros, got %q", key, value)
	}
	if offset > math.MaxUint64/1_000_000 {
		return 0, fmt.Errorf("%q is too large to convert to nanoseconds: %q", key, value)
	}
	return offset * 1_000_000, nil
}
