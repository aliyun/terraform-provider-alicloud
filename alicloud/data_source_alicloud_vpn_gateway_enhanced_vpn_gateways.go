// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudVpnGatewayEnhancedVpnGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudVpnGatewayEnhancedVpnGatewayRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpn_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_propagate": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disaster_recovery_vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_gateway_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAliCloudVpnGatewayEnhancedVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeVpnGateways"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("vpn_instance_id"); ok {
		request["VpnGatewayId"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("vpn_instance_id"); ok {
		request["VpnGatewayId"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, query, request, true)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.VpnGateways.VpnGateway[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["VpnGatewayId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["VpnGatewayId"]

		mapping["auto_propagate"] = objectRaw["AutoPropagate"]
		mapping["create_time"] = objectRaw["CreateTime"]
		mapping["description"] = objectRaw["Description"]
		mapping["disaster_recovery_vswitch_id"] = objectRaw["DisasterRecoveryVSwitchId"]
		mapping["gateway_type"] = objectRaw["GatewayType"]
		mapping["network_type"] = objectRaw["NetworkType"]
		mapping["status"] = objectRaw["Status"]
		mapping["vswitch_id"] = objectRaw["VSwitchId"]
		mapping["vpc_id"] = objectRaw["VpcId"]
		mapping["vpn_gateway_name"] = objectRaw["Name"]
		mapping["vpn_type"] = objectRaw["VpnType"]
		mapping["vpn_instance_id"] = objectRaw["VpnGatewayId"]

		tagsMaps, _ := jsonpath.Get("$.Tags.Tag", objectRaw)
		mapping["tags"] = tagsToMap(tagsMaps)

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["VpnGatewayId"])
		mapping, err = dataSourceAliCloudVpnGatewayEnhancedVpnGatewayReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("gateways", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudVpnGatewayEnhancedVpnGatewayReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	vPNGatewayServiceV2 := VPNGatewayServiceV2{client}
	getResp, err := vPNGatewayServiceV2.DescribeVpnGatewayEnhancedVpnGateway(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["auto_propagate"] = objectRaw["AutoPropagate"]
	mapping["create_time"] = objectRaw["CreateTime"]
	mapping["description"] = objectRaw["Description"]
	mapping["disaster_recovery_vswitch_id"] = objectRaw["DisasterRecoveryVSwitchId"]
	mapping["gateway_type"] = objectRaw["GatewayType"]
	mapping["network_type"] = objectRaw["NetworkType"]
	mapping["status"] = objectRaw["Status"]
	mapping["vswitch_id"] = objectRaw["VSwitchId"]
	mapping["vpc_id"] = objectRaw["VpcId"]
	mapping["vpn_gateway_name"] = objectRaw["Name"]
	mapping["vpn_type"] = objectRaw["VpnType"]
	mapping["vpn_instance_id"] = objectRaw["VpnGatewayId"]

	tagResp, err := vPNGatewayServiceV2.DescribeEnhancedVpnGatewayListTagResources(id)
	if err != nil {
		return nil, WrapError(err)
	}
	tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", tagResp)
	mapping["tags"] = tagsToMap(tagsMaps)

	return mapping, nil
}
