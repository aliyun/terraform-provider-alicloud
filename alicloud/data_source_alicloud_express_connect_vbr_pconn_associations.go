package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudExpressConnectVbrPconnAssociations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudExpressConnectVbrPconnAssociationsRead,
		Schema: map[string]*schema.Schema{
			"vbr_id": {
				Optional: true,
				Type:     schema.TypeString,
				ForceNew: true,
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"page_number": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"page_size": {
				Optional: true,
				Type:     schema.TypeInt,
				Default:  10,
			},
			"associations": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"circuit_code": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"enable_ipv6": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"local_gateway_ip": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"local_ipv6_gateway_ip": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"peer_gateway_ip": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"peer_ipv6_gateway_ip": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"peering_ipv6_subnet_mask": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"peering_subnet_mask": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"physical_connection_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"vbr_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"vlan_id": {
							Computed: true,
							Type:     schema.TypeInt,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudExpressConnectVbrPconnAssociationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	if v, ok := d.GetOk("vbr_id"); ok {
		request["Filter.1.Value.1"] = v
		request["Filter.1.Key"] = "VbrId"
	}

	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["PageNumber"] = v.(int)
	} else {
		request["PageNumber"] = 1
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeLarge
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeVirtualBorderRouters"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_express_connect_vbr_pconn_associations", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.VirtualBorderRouterSet.VirtualBorderRouterType", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.VirtualBorderRouterSet.VirtualBorderRouterType", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}

		for _, v := range result {
			associatedPhysicalConnectionItem := v.(map[string]interface{})

			if _, ok := associatedPhysicalConnectionItem["AssociatedPhysicalConnections"]; !ok {
				continue
			}

			resp, err := jsonpath.Get("$.AssociatedPhysicalConnections.AssociatedPhysicalConnection", associatedPhysicalConnectionItem)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AssociatedPhysicalConnections.AssociatedPhysicalConnection", response)
			}

			associatedPhysicalConnectionResult, _ := resp.([]interface{})

			for _, v := range associatedPhysicalConnectionResult {
				item := v.(map[string]interface{})
				item["VbrId"] = associatedPhysicalConnectionItem["VbrId"]
				if len(idsMap) > 0 {
					if _, ok := idsMap[fmt.Sprint(associatedPhysicalConnectionItem["VbrId"], ":", item["PhysicalConnectionId"])]; !ok {
						continue
					}
				}
				objects = append(objects, item)
			}
		}

		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                       fmt.Sprint(object["VbrId"], ":", object["PhysicalConnectionId"]),
			"circuit_code":             object["CircuitCode"],
			"enable_ipv6":              object["EnableIpv6"],
			"local_gateway_ip":         object["LocalGatewayIp"],
			"local_ipv6_gateway_ip":    object["LocalIpv6GatewayIp"],
			"peer_gateway_ip":          object["PeerGatewayIp"],
			"peer_ipv6_gateway_ip":     object["PeerIpv6GatewayIp"],
			"peering_ipv6_subnet_mask": object["PeeringIpv6SubnetMask"],
			"peering_subnet_mask":      object["PeeringSubnetMask"],
			"physical_connection_id":   object["PhysicalConnectionId"],
			"status":                   object["Status"],
			"vbr_id":                   object["VbrId"],
			"vlan_id":                  formatInt(object["VlanId"]),
		}

		ids = append(ids, fmt.Sprint(object["VbrId"], ":", object["PhysicalConnectionId"]))

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("associations", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
