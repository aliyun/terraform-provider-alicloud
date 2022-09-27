package alicloud

import (
	"fmt"
	"time"

	otsTunnel "github.com/aliyun/aliyun-tablestore-go-sdk/tunnel"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// resourceAlicloudOtsTunnel Tablestore tunnel not support update
func resourceAlicloudOtsTunnel() *schema.Resource {
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
	otsService := OtsService{client}
	// check table exists
	if err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, e := otsService.DescribeOtsTable(fmt.Sprintf("%s%s%s", instanceName, COLON_SEPARATED, tableName))
		if e != nil {
			if NotFoundError(e) {
				return resource.RetryableError(e)
			}
			return resource.NonRetryableError(e)
		}
		return nil
	}); err != nil {
		return WrapError(err)
	}

	request := new(otsTunnel.CreateTunnelRequest)
	request.TableName = tableName
	request.TunnelName = tunnelName
	request.Type = tunnelType

	var requestInfo otsTunnel.TunnelClient
	if err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		raw, err := client.WithTableStoreTunnelClient(instanceName, func(tunnelClient otsTunnel.TunnelClient) (interface{}, error) {
			requestInfo = tunnelClient
			return tunnelClient.CreateTunnel(request)
		})
		if err != nil {
			if IsExpectedErrors(err, OtsTunnelIsTemporarilyUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateTunnel", raw, requestInfo, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ots_tunnel", "CreateTunnel", AliyunTablestoreGoSdk)
	}

	d.SetId(fmt.Sprintf("%s%s%s%s%s", instanceName, COLON_SEPARATED, tableName, COLON_SEPARATED, tunnelName))
	return resourceAliyunOtsTunnelRead(d, meta)
}

func resourceAliyunOtsTunnelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	object, err := otsService.DescribeOtsTunnel(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
	otsService := OtsService{client}

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	req := &otsTunnel.DeleteTunnelRequest{
		TableName:  parts[1],
		TunnelName: parts[2],
	}
	var requestInfo otsTunnel.TunnelClient
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithTableStoreTunnelClient(parts[0], func(tunnelClient otsTunnel.TunnelClient) (interface{}, error) {
			requestInfo = tunnelClient
			return tunnelClient.DeleteTunnel(req)
		})
		if err != nil {
			if IsExpectedErrors(err, OtsTunnelIsTemporarilyUnavailable) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteTunnel", raw, requestInfo, req)
		return nil
	})
	if err != nil {
		if isOtsTunnelNotFound(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteTunnel", AliyunTablestoreGoSdk)
	}
	return WrapError(otsService.WaitForOtsTunnel(d.Id(), Deleted, DefaultTimeout))
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
		return "", WrapError(Error("unknown ots tunnel type: " + orig))
	}
}
