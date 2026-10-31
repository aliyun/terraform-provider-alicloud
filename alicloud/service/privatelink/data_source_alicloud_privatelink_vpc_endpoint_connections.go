package privatelink

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceVpcEndpointConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPrivatelinkVpcEndpointConnectionsRead,
		Schema: map[string]*schema.Schema{
			"endpoint_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoint_owner_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Connected", "Connecting", "Deleted", "Deleting", "Disconnected", "Disconnecting", "Pending", "ServiceDeleted"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"traffic_control_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_owner_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudPrivatelinkVpcEndpointConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListVpcEndpointConnections"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("endpoint_id"); ok {
		request["EndpointId"] = v
	}
	if v, ok := d.GetOk("endpoint_owner_id"); ok {
		request["EndpointOwnerId"] = v
	}
	request["RegionId"] = client.RegionId
	request["ServiceId"] = d.Get("service_id")
	if v, ok := d.GetOk("status"); ok {
		request["ConnectionStatus"] = v
	}
	request["MaxResults"] = helper.PageSizeLarge
	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	for {

		response, err = client.RpcPost("Privatelink", "2020-04-15", action, nil, request, true)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_privatelink_vpc_endpoint_connections", action, errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(action, response, request)

		resp, err := jsonpath.Get("$.Connections", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Connections", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                   fmt.Sprint(object["ServiceId"], ":", object["EndpointId"]),
			"bandwidth":            object["Bandwidth"],
			"traffic_control_mode": object["TrafficControlMode"],
			"endpoint_id":          object["EndpointId"],
			"endpoint_owner_id":    object["EndpointOwnerId"],
			"status":               object["ConnectionStatus"],
		}
		ids = append(ids, fmt.Sprint(object["ServiceId"], ":", object["EndpointId"]))
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("connections", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
