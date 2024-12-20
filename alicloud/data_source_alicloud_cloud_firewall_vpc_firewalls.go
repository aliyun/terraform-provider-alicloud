package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCloudFirewallVpcFirewalls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCloudFirewallVpcFirewallsRead,
		Schema: map[string]*schema.Schema{
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
			"region_no": {
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
			"vpc_id": {
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
			"name_regex": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.ValidateRegexp,
			},
			"names": {
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
			"firewalls": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"bandwidth": {
							Computed: true,
							Type:     schema.TypeInt,
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
									"eni_id": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"eni_private_ip_address": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"local_vpc_cidr_table_list": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"local_route_entry_list": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"local_destination_cidr": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"local_next_hop_instance_id": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"local_route_table_id": {
													Computed: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"region_no": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"router_interface_id": {
										Computed: true,
										Type:     schema.TypeString,
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
						"peer_vpc": {
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
									"peer_vpc_cidr_table_list": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"peer_route_entry_list": {
													Computed: true,
													Type:     schema.TypeList,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"peer_destination_cidr": {
																Computed: true,
																Type:     schema.TypeString,
															},
															"peer_next_hop_instance_id": {
																Computed: true,
																Type:     schema.TypeString,
															},
														},
													},
												},
												"peer_route_table_id": {
													Computed: true,
													Type:     schema.TypeString,
												},
											},
										},
									},
									"region_no": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"router_interface_id": {
										Computed: true,
										Type:     schema.TypeString,
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
						"region_status": {
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

func dataSourceAlicloudCloudFirewallVpcFirewallsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}
	if v, ok := d.GetOk("region_no"); ok {
		request["RegionNo"] = v
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
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
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

	var vpcFirewallNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		vpcFirewallNameRegex = r
	}

	var err error
	var endpoint string
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "DescribeVpcFirewallList"
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_firewall_vpc_firewalls", action, AlibabaCloudSdkGoERROR)
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

			if vpcFirewallNameRegex != nil && !vpcFirewallNameRegex.MatchString(fmt.Sprint(item["VpcFirewallName"])) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                fmt.Sprint(object["VpcFirewallId"]),
			"bandwidth":         object["Bandwidth"],
			"connect_type":      object["ConnectType"],
			"region_status":     object["RegionStatus"],
			"status":            object["FirewallSwitchStatus"],
			"vpc_firewall_id":   object["VpcFirewallId"],
			"vpc_firewall_name": object["VpcFirewallName"],
		}

		localVpcMaps := make([]map[string]interface{}, 0)
		localVpcMap := make(map[string]interface{})
		localVpcRaw := object["LocalVpc"].(map[string]interface{})
		localVpcCidrTableListMaps := make([]map[string]interface{}, 0)
		localVpcCidrTableListRaw := localVpcRaw["VpcCidrTableList"]
		for _, value1 := range localVpcCidrTableListRaw.([]interface{}) {
			localVpcCidrTableList := value1.(map[string]interface{})
			localVpcCidrTableListMap := make(map[string]interface{})
			localRouteEntryListMaps := make([]map[string]interface{}, 0)
			localRouteEntryListRaw := localVpcCidrTableList["RouteEntryList"]
			for _, value2 := range localRouteEntryListRaw.([]interface{}) {
				localRouteEntryList := value2.(map[string]interface{})
				localRouteEntryListMap := make(map[string]interface{})
				localRouteEntryListMap["local_destination_cidr"] = localRouteEntryList["DestinationCidr"]
				localRouteEntryListMap["local_next_hop_instance_id"] = localRouteEntryList["NextHopInstanceId"]
				localRouteEntryListMaps = append(localRouteEntryListMaps, localRouteEntryListMap)
			}
			localVpcCidrTableListMap["local_route_entry_list"] = localRouteEntryListMaps
			localVpcCidrTableListMap["local_route_table_id"] = localVpcCidrTableList["RouteTableId"]
			localVpcCidrTableListMaps = append(localVpcCidrTableListMaps, localVpcCidrTableListMap)
		}
		localVpcMap["local_vpc_cidr_table_list"] = localVpcCidrTableListMaps
		localVpcMap["region_no"] = localVpcRaw["RegionNo"]
		localVpcMap["vpc_id"] = localVpcRaw["VpcId"]
		localVpcMap["vpc_name"] = localVpcRaw["VpcName"]
		localVpcMaps = append(localVpcMaps, localVpcMap)
		mapping["local_vpc"] = localVpcMaps
		peerVpcMaps := make([]map[string]interface{}, 0)
		peerVpcMap := make(map[string]interface{})
		peerVpcRaw := object["PeerVpc"].(map[string]interface{})
		peerVpcCidrTableListMaps := make([]map[string]interface{}, 0)
		peerVpcCidrTableListRaw := peerVpcRaw["VpcCidrTableList"]
		for _, value1 := range peerVpcCidrTableListRaw.([]interface{}) {
			peerVpcCidrTableList := value1.(map[string]interface{})
			peerVpcCidrTableListMap := make(map[string]interface{})
			peerRouteEntryListMaps := make([]map[string]interface{}, 0)
			peerRouteEntryListRaw := peerVpcCidrTableList["RouteEntryList"]
			for _, value2 := range peerRouteEntryListRaw.([]interface{}) {
				peerRouteEntryList := value2.(map[string]interface{})
				peerRouteEntryListMap := make(map[string]interface{})
				peerRouteEntryListMap["peer_destination_cidr"] = peerRouteEntryList["DestinationCidr"]
				peerRouteEntryListMap["peer_next_hop_instance_id"] = peerRouteEntryList["NextHopInstanceId"]
				peerRouteEntryListMaps = append(peerRouteEntryListMaps, peerRouteEntryListMap)
			}
			peerVpcCidrTableListMap["peer_route_entry_list"] = peerRouteEntryListMaps
			peerVpcCidrTableListMap["peer_route_table_id"] = peerVpcCidrTableList["RouteTableId"]
			peerVpcCidrTableListMaps = append(peerVpcCidrTableListMaps, peerVpcCidrTableListMap)
		}
		peerVpcMap["peer_vpc_cidr_table_list"] = peerVpcCidrTableListMaps
		peerVpcMap["region_no"] = peerVpcRaw["RegionNo"]
		peerVpcMap["vpc_id"] = peerVpcRaw["VpcId"]
		peerVpcMap["vpc_name"] = peerVpcRaw["VpcName"]
		peerVpcMaps = append(peerVpcMaps, peerVpcMap)
		mapping["peer_vpc"] = peerVpcMaps

		ids = append(ids, fmt.Sprint(object["VpcFirewallId"]))
		names = append(names, object["VpcFirewallName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
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
