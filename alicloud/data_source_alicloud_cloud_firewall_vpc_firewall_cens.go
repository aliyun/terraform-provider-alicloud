package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCloudFirewallVpcFirewallCens() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCloudFirewallVpcFirewallCensRead,
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"lang": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"member_uid": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"network_instance_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"status": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"vpc_firewall_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"vpc_firewall_name": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
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
			"cens": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"cen_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"connect_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"local_vpc": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attachment_id": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"attachment_name": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"defend_cidr_list": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"eni_list": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"eni_id": {
													Computed: true,
													Type:     schema.TypeString,
												},
												"eni_private_ip_address": {
													Computed: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"manual_vswitch_id": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"network_instance_id": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"network_instance_name": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"network_instance_type": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"owner_id": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"region_no": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"route_mode": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"support_manual_mode": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"transit_router_id": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"transit_router_type": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"vpc_cidr_table_list": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"route_entry_list": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"destination_cidr": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"next_hop_instance_id": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"route_table_id": {
													Computed: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"vpc_id": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"vpc_name": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"network_instance_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"vpc_firewall_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"vpc_firewall_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCloudFirewallVpcFirewallCensRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("cen_id"); ok {
		request["CenId"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}
	if v, ok := d.GetOk("network_instance_id"); ok {
		request["NetworkInstanceId"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["FirewallSwitchStatus"] = v
	}
	if v, ok := d.GetOk("vpc_firewall_id"); ok {
		request["VpcFirewallId"] = v
	}
	if v, ok := d.GetOk("vpc_firewall_name"); ok {
		request["VpcFirewallName"] = v
	}
	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["CurrentPage"] = v.(int)
	} else {
		request["CurrentPage"] = 1
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

	var objects []interface{}
	var response map[string]interface{}
	var err error
	var endpoint string

	for {
		action := "DescribeVpcFirewallCenList"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				} else if IsExpectedErrors(err, []string{"not buy user"}) {
					endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_firewall_vpc_firewall_cens", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.VpcFirewalls", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.VpcFirewalls", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["VpcFirewallId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                fmt.Sprint(object["VpcFirewallId"]),
			"cen_id":            object["CenId"],
			"connect_type":      object["ConnectType"],
			"status":            object["FirewallSwitchStatus"],
			"vpc_firewall_id":   object["VpcFirewallId"],
			"vpc_firewall_name": object["VpcFirewallName"],
		}

		localVpcMaps := make([]map[string]interface{}, 0)
		localVpcMap := make(map[string]interface{})
		localVpcRaw := object["LocalVpc"].(map[string]interface{})
		localVpcMap["defend_cidr_list"] = localVpcRaw["DefendCidrList"].([]interface{})
		localVpcMap["manual_vswitch_id"] = localVpcRaw["ManualVSwitchId"]
		localVpcMap["network_instance_id"] = localVpcRaw["NetworkInstanceId"]
		localVpcMap["network_instance_name"] = localVpcRaw["NetworkInstanceName"]
		localVpcMap["network_instance_type"] = localVpcRaw["NetworkInstanceType"]
		localVpcMap["owner_id"] = localVpcRaw["OwnerId"]
		localVpcMap["region_no"] = localVpcRaw["RegionNo"]
		localVpcMap["route_mode"] = localVpcRaw["RouteMode"]
		localVpcMap["support_manual_mode"] = localVpcRaw["SupportManualMode"]
		localVpcMap["transit_router_type"] = localVpcRaw["TransitRouterType"]
		vpcCidrTableListMaps := make([]map[string]interface{}, 0)
		vpcCidrTableListRaw := localVpcRaw["VpcCidrTableList"]
		for _, value1 := range vpcCidrTableListRaw.([]interface{}) {
			vpcCidrTableList := value1.(map[string]interface{})
			vpcCidrTableListMap := make(map[string]interface{})
			routeEntryListMaps := make([]map[string]interface{}, 0)
			routeEntryListRaw := vpcCidrTableList["RouteEntryList"]
			for _, value2 := range routeEntryListRaw.([]interface{}) {
				routeEntryList := value2.(map[string]interface{})
				routeEntryListMap := make(map[string]interface{})
				routeEntryListMap["destination_cidr"] = routeEntryList["DestinationCidr"]
				routeEntryListMap["next_hop_instance_id"] = routeEntryList["NextHopInstanceId"]
				routeEntryListMaps = append(routeEntryListMaps, routeEntryListMap)
			}
			vpcCidrTableListMap["route_entry_list"] = routeEntryListMaps
			vpcCidrTableListMap["route_table_id"] = vpcCidrTableList["RouteTableId"]
			vpcCidrTableListMaps = append(vpcCidrTableListMaps, vpcCidrTableListMap)
		}
		localVpcMap["vpc_cidr_table_list"] = vpcCidrTableListMaps
		localVpcMap["vpc_id"] = localVpcRaw["VpcId"]
		localVpcMap["vpc_name"] = localVpcRaw["VpcName"]
		localVpcMaps = append(localVpcMaps, localVpcMap)
		mapping["local_vpc"] = localVpcMaps

		ids = append(ids, fmt.Sprint(object["VpcFirewallId"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("cens", s); err != nil {
		return WrapError(err)
	}
	return nil
}
