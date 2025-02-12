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

func dataSourceAliCloudCloudFirewallNatFirewalls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCloudFirewallNatFirewallRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"member_uid": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"proxy_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"proxy_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"region_no": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"firewalls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ali_uid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"member_uid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"nat_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat_gateway_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat_route_entry_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nexthop_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"route_table_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"nexthop_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"destination_cidr": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"proxy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"proxy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"strict_mode": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vpc_id": {
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
				ForceNew: true,
			},
		},
	}
}

func dataSourceAliCloudCloudFirewallNatFirewallRead(d *schema.ResourceData, meta interface{}) error {
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
	action := "DescribeNatFirewallList"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	if v, ok := d.GetOk("proxy_id"); ok {
		request["ProxyId"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOkExists("member_uid"); ok {
		request["MemberUid"] = v
	}
	request["NatGatewayId"] = d.Get("nat_gateway_id")
	if v, ok := d.GetOk("proxy_id"); ok {
		request["ProxyId"] = v
	}
	request["ProxyName"] = d.Get("proxy_name")
	request["RegionNo"] = d.Get("region_no")
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	request["VpcId"] = d.Get("vpc_id")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)

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

		resp, _ := jsonpath.Get("$.NatFirewallList[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ProxyId"])]; !ok {
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["ProxyId"]

		mapping["ali_uid"] = objectRaw["AliUid"]
		mapping["member_uid"] = objectRaw["MemberUid"]
		mapping["nat_gateway_id"] = objectRaw["NatGatewayId"]
		mapping["nat_gateway_name"] = objectRaw["NatGatewayName"]
		mapping["proxy_name"] = objectRaw["ProxyName"]
		mapping["strict_mode"] = objectRaw["StrictMode"]
		mapping["vpc_id"] = objectRaw["VpcId"]
		mapping["proxy_id"] = objectRaw["ProxyId"]

		natRouteEntryListRaw := objectRaw["NatRouteEntryList"]
		natRouteEntryListMaps := make([]map[string]interface{}, 0)
		if natRouteEntryListRaw != nil {
			for _, natRouteEntryListChildRaw := range natRouteEntryListRaw.([]interface{}) {
				natRouteEntryListMap := make(map[string]interface{})
				natRouteEntryListChildRaw := natRouteEntryListChildRaw.(map[string]interface{})
				natRouteEntryListMap["destination_cidr"] = natRouteEntryListChildRaw["DestinationCidr"]
				natRouteEntryListMap["nexthop_id"] = natRouteEntryListChildRaw["NextHopId"]
				natRouteEntryListMap["nexthop_type"] = natRouteEntryListChildRaw["NextHopType"]
				natRouteEntryListMap["route_table_id"] = natRouteEntryListChildRaw["RouteTableId"]

				natRouteEntryListMaps = append(natRouteEntryListMaps, natRouteEntryListMap)
			}
		}
		mapping["nat_route_entry_list"] = natRouteEntryListMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw[""])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("firewalls", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
