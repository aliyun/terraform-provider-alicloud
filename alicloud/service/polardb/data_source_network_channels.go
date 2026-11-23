package polardb

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DataSourceNetworkChannels queries existing channels without changing their lifecycle.
func DataSourceNetworkChannels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkChannelsRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type: schema.TypeString, Required: true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"channel_name": {
				Type: schema.TypeString, Optional: true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"resource_group_id": {
				Type: schema.TypeString, Optional: true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"vpc_id": {
				Type: schema.TypeString, Optional: true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"names": {
				Type: schema.TypeList, Computed: true,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"channels": {
				Type: schema.TypeList, Computed: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"channel_name":         {Type: schema.TypeString, Computed: true},
					"db_cluster_id":        {Type: schema.TypeString, Computed: true},
					"eni_id":               {Type: schema.TypeString, Computed: true},
					"notes":                {Type: schema.TypeString, Computed: true},
					"region_id":            {Type: schema.TypeString, Computed: true},
					"target_db_cluster_id": {Type: schema.TypeString, Computed: true},
					"target_ip":            {Type: schema.TypeString, Computed: true},
					"target_port":          {Type: schema.TypeString, Computed: true},
					"target_type":          {Type: schema.TypeString, Computed: true},
					"vpc_id":               {Type: schema.TypeString, Computed: true},
				}},
			},
		},
	}
}

func dataSourceNetworkChannelsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutRead))
	defer cancel()
	client := meta.(*connectivity.AliyunClient)
	clusterID := strings.TrimSpace(d.Get("db_cluster_id").(string))
	channelName := strings.TrimSpace(d.Get("channel_name").(string))
	request := map[string]interface{}{
		"RegionId": client.RegionId, "DBClusterId": clusterID,
	}
	for field, key := range map[string]string{
		"channel_name": "ChannelName", "resource_group_id": "ResourceGroupId", "vpc_id": "VpcId",
	} {
		if value, ok := d.GetOk(field); ok {
			request[key] = strings.TrimSpace(value.(string))
		}
	}
	var response map[string]interface{}
	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		var err error
		response, err = client.RpcPostWithContext(ctx, "polardb", "2017-08-01", "DescribeNetworkChannel", nil, request)
		if err != nil {
			if helper.NeedRetry(err) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_polardb_network_channels", "DescribeNetworkChannel", errs.AlibabaCloudSdkGoERROR))
	}
	channels, names, err := flattenNetworkChannels(response["ChannelInfos"], clusterID, client.RegionId, channelName)
	if err != nil {
		if requestID, ok := response["RequestId"].(string); ok && requestID != "" {
			err = fmt.Errorf("DescribeNetworkChannel RequestId %s: %w", requestID, err)
		}
		return diag.FromErr(errs.WrapError(err))
	}
	if err := d.Set("channels", channels); err != nil {
		return diag.FromErr(errs.WrapError(err))
	}
	if err := d.Set("names", names); err != nil {
		return diag.FromErr(errs.WrapError(err))
	}
	identity, err := json.Marshal(request)
	if err != nil {
		return diag.FromErr(errs.WrapError(err))
	}
	d.SetId(helper.DataResourceIdHash([]string{string(identity)}))
	return nil
}

// flattenNetworkChannels accepts the documented array and legacy SDK encodings.
// Unknown shapes are errors, not evidence that existing channels disappeared.
func flattenNetworkChannels(value interface{}, clusterID, regionID, channelName string) ([]map[string]interface{}, []string, error) {
	if encoded, ok := value.(string); ok {
		if err := json.Unmarshal([]byte(encoded), &value); err != nil {
			return nil, nil, fmt.Errorf("DescribeNetworkChannel returned invalid JSON in ChannelInfos: %w", err)
		}
	}
	if wrapper, ok := value.(map[string]interface{}); ok {
		var exists bool
		value, exists = wrapper["ChannelInfo"]
		if !exists {
			return nil, nil, fmt.Errorf("DescribeNetworkChannel returned an unrecognized ChannelInfos object")
		}
	}
	var objects []interface{}
	switch typed := value.(type) {
	case []interface{}:
		objects = typed
	case []map[string]interface{}:
		for _, item := range typed {
			objects = append(objects, item)
		}
	default:
		return nil, nil, fmt.Errorf("DescribeNetworkChannel returned missing or invalid ChannelInfos (expected an array, got %T)", value)
	}
	channels := make([]map[string]interface{}, 0, len(objects))
	for i, value := range objects {
		object, ok := value.(map[string]interface{})
		if !ok || object == nil {
			return nil, nil, fmt.Errorf("DescribeNetworkChannel ChannelInfos[%d] is not an object", i)
		}
		item := make(map[string]interface{}, 10)
		for field, key := range map[string]string{
			"channel_name": "ChannelName", "db_cluster_id": "DBClusterId", "notes": "Notes", "region_id": "RegionId",
			"eni_id":               "EniId",
			"target_db_cluster_id": "TargetDBClusterId", "target_ip": "TargetIp", "target_port": "TargetPort", "target_type": "TargetType", "vpc_id": "VpcId",
		} {
			var text string
			if raw := object[key]; raw != nil {
				var ok bool
				text, ok = raw.(string)
				if !ok {
					return nil, nil, fmt.Errorf("DescribeNetworkChannel ChannelInfos[%d].%s is not a string", i, key)
				}
			}
			item[field] = text
		}
		for _, field := range []string{"channel_name", "db_cluster_id", "region_id"} {
			item[field] = strings.TrimSpace(item[field].(string))
		}
		if item["channel_name"] == "" {
			return nil, nil, fmt.Errorf("DescribeNetworkChannel ChannelInfos[%d].ChannelName is missing", i)
		}
		if item["db_cluster_id"] == "" {
			item["db_cluster_id"] = clusterID
		}
		if item["region_id"] == "" {
			item["region_id"] = regionID
		}
		if item["db_cluster_id"] != clusterID || item["region_id"] != regionID {
			return nil, nil, fmt.Errorf("DescribeNetworkChannel ChannelInfos[%d] does not belong to the queried cluster and region", i)
		}
		if channelName != "" && item["channel_name"] != channelName {
			continue
		}
		channels = append(channels, item)
	}
	sort.Slice(channels, func(i, j int) bool {
		if channels[i]["channel_name"] != channels[j]["channel_name"] {
			return channels[i]["channel_name"].(string) < channels[j]["channel_name"].(string)
		}
		left, _ := json.Marshal(channels[i])
		right, _ := json.Marshal(channels[j])
		return string(left) < string(right)
	})
	names := make([]string, 0, len(channels))
	for _, channel := range channels {
		names = append(names, channel["channel_name"].(string))
	}
	return channels, names, nil
}
