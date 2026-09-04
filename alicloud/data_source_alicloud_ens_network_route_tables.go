// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAliCloudEnsNetworkRouteTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEnsNetworkRouteTablesRead,
		Schema: map[string]*schema.Schema{
			"route_table_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"associate_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"route_table_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_table_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"associate_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_table_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_default_gateway_route_table": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudEnsNetworkRouteTablesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeEnsRouteTables"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("route_table_name"); ok {
		request["RouteTableName"] = v
	}
	if v, ok := d.GetOk("associate_type"); ok {
		request["AssociateType"] = v
	}
	if v, ok := d.GetOk("route_table_type"); ok {
		request["Type"] = v
	}
	if v, ok := d.GetOk("route_table_id"); ok {
		request["RouteTableId"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ens_network_route_tables", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.RouteTables", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.RouteTables", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	networkIdFilter := ""
	if v, ok := d.GetOk("network_id"); ok {
		networkIdFilter = v.(string)
	}
	for _, object := range objects {
		if networkIdFilter != "" && fmt.Sprint(object["NetworkId"]) != networkIdFilter {
			continue
		}
		mapping := map[string]interface{}{
			"id":                             fmt.Sprint(object["RouteTableId"]),
			"route_table_id":                 object["RouteTableId"],
			"route_table_name":               object["RouteTableName"],
			"description":                    object["Description"],
			"associate_type":                 object["AssociateType"],
			"network_id":                     object["NetworkId"],
			"route_table_type":               object["Type"],
			"status":                         object["Status"],
			"is_default_gateway_route_table": object["IsDefaultGatewayRouteTable"],
			"create_time":                    object["CreationTime"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["RouteTableName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("tables", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
