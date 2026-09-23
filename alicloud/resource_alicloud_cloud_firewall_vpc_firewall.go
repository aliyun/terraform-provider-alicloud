package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudFirewallVpcFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallVpcFirewallCreate,
		Read:   resourceAliCloudCloudFirewallVpcFirewallRead,
		Update: resourceAliCloudCloudFirewallVpcFirewallUpdate,
		Delete: resourceAliCloudCloudFirewallVpcFirewallDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(31 * time.Minute),
			Update: schema.DefaultTimeout(31 * time.Minute),
			Delete: schema.DefaultTimeout(31 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"vpc_firewall_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"open", "close"}, false),
			},
			"member_uid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"zh", "en"}, false),
			},
			"local_vpc": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"region_no": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"local_vpc_cidr_table_list": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"local_route_table_id": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"local_route_entry_list": {
										Type:     schema.TypeList,
										Required: true,
										ForceNew: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"local_next_hop_instance_id": {
													Type:     schema.TypeString,
													Required: true,
													ForceNew: true,
												},
												"local_destination_cidr": {
													Type:     schema.TypeString,
													Required: true,
													ForceNew: true,
												},
											},
										},
									},
								},
							},
						},
						"vpc_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eni_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eni_private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"router_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"peer_vpc": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"region_no": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"peer_vpc_cidr_table_list": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"peer_route_table_id": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"peer_route_entry_list": {
										Type:     schema.TypeList,
										Required: true,
										ForceNew: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"peer_destination_cidr": {
													Type:     schema.TypeString,
													Required: true,
													ForceNew: true,
												},
												"peer_next_hop_instance_id": {
													Type:     schema.TypeString,
													Required: true,
													ForceNew: true,
												},
											},
										},
									},
								},
							},
						},
						"vpc_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eni_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eni_private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"router_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"vpc_firewall_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connect_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"region_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCloudFirewallVpcFirewallCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}
	var response map[string]interface{}
	var err error
	var endpoint string
	action := "CreateVpcFirewallConfigure"
	request := make(map[string]interface{})

	request["VpcFirewallName"] = d.Get("vpc_firewall_name")
	request["FirewallSwitch"] = d.Get("status")

	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	if v, ok := d.GetOk("local_vpc"); ok {
		localVpcCidrTableList, err := jsonpath.Get("$[0].local_vpc_cidr_table_list", v)
		if err != nil {
			return WrapError(err)
		}
		localVpcCidrTableListMaps := make([]map[string]interface{}, 0)
		localVpcCidrTableListRaw := localVpcCidrTableList
		for _, value1 := range localVpcCidrTableListRaw.([]interface{}) {
			localVpcCidrTableListTmp := value1.(map[string]interface{})
			localVpcCidrTableListMap := make(map[string]interface{})
			localRouteEntryListMaps := make([]map[string]interface{}, 0)
			peerRouteEntryListRaw := localVpcCidrTableListTmp["local_route_entry_list"]
			for _, value2 := range peerRouteEntryListRaw.([]interface{}) {
				localRouteEntryList := value2.(map[string]interface{})
				localRouteEntryListMap := make(map[string]interface{})
				localRouteEntryListMap["DestinationCidr"] = localRouteEntryList["local_destination_cidr"]
				localRouteEntryListMap["NextHopInstanceId"] = localRouteEntryList["local_next_hop_instance_id"]
				localRouteEntryListMaps = append(localRouteEntryListMaps, localRouteEntryListMap)
			}
			localVpcCidrTableListMap["RouteEntryList"] = localRouteEntryListMaps
			localVpcCidrTableListMap["RouteTableId"] = localVpcCidrTableListTmp["local_route_table_id"]
			localVpcCidrTableListMaps = append(localVpcCidrTableListMaps, localVpcCidrTableListMap)
		}
		localVpcCidrTableListJson, err := json.Marshal(localVpcCidrTableListMaps)
		if err != nil {
			return WrapError(err)
		}
		request["LocalVpcCidrTableList"] = string(localVpcCidrTableListJson)
	}

	if v, ok := d.GetOk("peer_vpc"); ok {
		peerVpcCidrTableList, err := jsonpath.Get("$[0].peer_vpc_cidr_table_list", v)
		if err != nil {
			return WrapError(err)
		}
		peerVpcCidrTableListMaps := make([]map[string]interface{}, 0)
		peerVpcCidrTableListRaw := peerVpcCidrTableList
		for _, value1 := range peerVpcCidrTableListRaw.([]interface{}) {
			peerVpcCidrTableListTmp := value1.(map[string]interface{})
			peerVpcCidrTableListMap := make(map[string]interface{})
			peerRouteEntryListMaps := make([]map[string]interface{}, 0)
			peerRouteEntryListRaw := peerVpcCidrTableListTmp["peer_route_entry_list"]
			for _, value2 := range peerRouteEntryListRaw.([]interface{}) {
				peerRouteEntryList := value2.(map[string]interface{})
				peerRouteEntryListMap := make(map[string]interface{})
				peerRouteEntryListMap["DestinationCidr"] = peerRouteEntryList["peer_destination_cidr"]
				peerRouteEntryListMap["NextHopInstanceId"] = peerRouteEntryList["peer_next_hop_instance_id"]
				peerRouteEntryListMaps = append(peerRouteEntryListMaps, peerRouteEntryListMap)
			}
			peerVpcCidrTableListMap["RouteEntryList"] = peerRouteEntryListMaps
			peerVpcCidrTableListMap["RouteTableId"] = peerVpcCidrTableListTmp["peer_route_table_id"]
			peerVpcCidrTableListMaps = append(peerVpcCidrTableListMaps, peerVpcCidrTableListMap)
		}
		peerVpcCidrTableListJson, err := json.Marshal(peerVpcCidrTableListMaps)
		if err != nil {
			return WrapError(err)
		}
		request["PeerVpcCidrTableList"] = string(peerVpcCidrTableListJson)
	}

	if v, ok := d.GetOk("local_vpc"); ok {
		localVpcId, err := jsonpath.Get("$[0].vpc_id", v)
		if err != nil {
			return WrapError(err)
		}
		request["LocalVpcId"] = localVpcId
		localVpcRegion, err := jsonpath.Get("$[0].region_no", v)
		if err != nil {
			return WrapError(err)
		}
		request["LocalVpcRegion"] = localVpcRegion
	}

	if v, ok := d.GetOk("peer_vpc"); ok {
		peerVpcId, err := jsonpath.Get("$[0].vpc_id", v)
		if err != nil {
			return WrapError(err)
		}
		request["PeerVpcId"] = peerVpcId
		peerVpcRegion, err := jsonpath.Get("$[0].region_no", v)
		if err != nil {
			return WrapError(err)
		}
		request["PeerVpcRegion"] = peerVpcRegion
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			} else if IsExpectedErrors(err, []string{"ErrorVpcFirewallNotFound", "vpc firewall not found"}) {
				if err := cloudfwService.CreateVpcFirewallTask(); err != nil {
					log.Println("[ERROR] syncing vpc configure failed by api CreateVpcFirewallTask. Error:", err)
				}
				wait()
				return resource.RetryableError(err)
			} else if IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_vpc_firewall", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.VpcFirewallId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_cloud_firewall_vpc_firewall")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	stateConf := BuildStateConf([]string{}, []string{"closed", "opened"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cloudfwService.CloudFirewallVpcFirewallStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCloudFirewallVpcFirewallRead(d, meta)
}

func resourceAliCloudCloudFirewallVpcFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}

	object, err := cloudfwService.DescribeCloudFirewallVpcFirewall(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_vpc_firewall cloudfwService.DescribeCloudFirewallVpcFirewall Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	objectExtra, err := cloudfwService.DescribeVpcFirewallList(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_vpc_firewall cloudfwService.DescribeCloudFirewallVpcFirewall Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("id", object["VpcFirewallId"])
	d.Set("vpc_firewall_id", object["VpcFirewallId"])
	d.Set("bandwidth", object["Bandwidth"])
	d.Set("connect_type", object["ConnectType"])
	d.Set("member_uid", objectExtra["MemberUid"])
	localVpcMaps := make([]map[string]interface{}, 0)
	localVpcMap := make(map[string]interface{})
	localVpcRaw := object["LocalVpc"].(map[string]interface{})
	localVpcMap["eni_id"] = localVpcRaw["EniId"]
	localVpcMap["eni_private_ip_address"] = localVpcRaw["EniPrivateIpAddress"]
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
	localVpcMap["router_interface_id"] = localVpcRaw["RouterInterfaceId"]
	localVpcMaps = append(localVpcMaps, localVpcMap)
	d.Set("local_vpc", localVpcMaps)
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
	peerVpcMap["eni_id"] = peerVpcRaw["EniId"]
	peerVpcMap["eni_private_ip_address"] = peerVpcRaw["EniPrivateIpAddress"]
	peerVpcMap["peer_vpc_cidr_table_list"] = peerVpcCidrTableListMaps
	peerVpcMap["region_no"] = peerVpcRaw["RegionNo"]
	peerVpcMap["vpc_id"] = peerVpcRaw["VpcId"]
	peerVpcMap["vpc_name"] = peerVpcRaw["VpcName"]
	peerVpcMap["router_interface_id"] = peerVpcRaw["RouterInterfaceId"]
	peerVpcMaps = append(peerVpcMaps, peerVpcMap)
	d.Set("peer_vpc", peerVpcMaps)
	d.Set("status", convertCloudFirewallVpcFirewallStatusRequest(object["FirewallSwitchStatus"]))
	d.Set("vpc_firewall_name", object["VpcFirewallName"])

	return nil
}

func resourceAliCloudCloudFirewallVpcFirewallUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"VpcFirewallId": d.Id(),
	}

	if d.HasChange("vpc_firewall_name") {
		update = true
	}
	request["VpcFirewallName"] = d.Get("vpc_firewall_name")

	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	if v, ok := d.GetOk("local_vpc"); ok {
		localVpcCidrTableList, err := jsonpath.Get("$[0].local_vpc_cidr_table_list", v)
		if err != nil {
			return WrapError(err)
		}
		localVpcCidrTableListMaps := make([]map[string]interface{}, 0)
		localVpcCidrTableListRaw := localVpcCidrTableList
		for _, value1 := range localVpcCidrTableListRaw.([]interface{}) {
			localVpcCidrTableListTmp := value1.(map[string]interface{})
			localVpcCidrTableListMap := make(map[string]interface{})
			localRouteEntryListMaps := make([]map[string]interface{}, 0)
			peerRouteEntryListRaw := localVpcCidrTableListTmp["local_route_entry_list"]
			for _, value2 := range peerRouteEntryListRaw.([]interface{}) {
				localRouteEntryList := value2.(map[string]interface{})
				localRouteEntryListMap := make(map[string]interface{})
				localRouteEntryListMap["DestinationCidr"] = localRouteEntryList["local_destination_cidr"]
				localRouteEntryListMap["NextHopInstanceId"] = localRouteEntryList["local_next_hop_instance_id"]
				localRouteEntryListMaps = append(localRouteEntryListMaps, localRouteEntryListMap)
			}
			localVpcCidrTableListMap["RouteEntryList"] = localRouteEntryListMaps
			localVpcCidrTableListMap["RouteTableId"] = localVpcCidrTableListTmp["local_route_table_id"]
			localVpcCidrTableListMaps = append(localVpcCidrTableListMaps, localVpcCidrTableListMap)
		}
		localVpcCidrTableListJson, err := json.Marshal(localVpcCidrTableListMaps)
		if err != nil {
			return WrapError(err)
		}
		request["LocalVpcCidrTableList"] = string(localVpcCidrTableListJson)
	}

	if v, ok := d.GetOk("peer_vpc"); ok {
		peerVpcCidrTableList, err := jsonpath.Get("$[0].peer_vpc_cidr_table_list", v)
		if err != nil {
			return WrapError(err)
		}
		peerVpcCidrTableListMaps := make([]map[string]interface{}, 0)
		peerVpcCidrTableListRaw := peerVpcCidrTableList
		for _, value1 := range peerVpcCidrTableListRaw.([]interface{}) {
			peerVpcCidrTableListTmp := value1.(map[string]interface{})
			peerVpcCidrTableListMap := make(map[string]interface{})
			peerRouteEntryListMaps := make([]map[string]interface{}, 0)
			peerRouteEntryListRaw := peerVpcCidrTableListTmp["peer_route_entry_list"]
			for _, value2 := range peerRouteEntryListRaw.([]interface{}) {
				peerRouteEntryList := value2.(map[string]interface{})
				peerRouteEntryListMap := make(map[string]interface{})
				peerRouteEntryListMap["DestinationCidr"] = peerRouteEntryList["peer_destination_cidr"]
				peerRouteEntryListMap["NextHopInstanceId"] = peerRouteEntryList["peer_next_hop_instance_id"]
				peerRouteEntryListMaps = append(peerRouteEntryListMaps, peerRouteEntryListMap)
			}
			peerVpcCidrTableListMap["RouteEntryList"] = peerRouteEntryListMaps
			peerVpcCidrTableListMap["RouteTableId"] = peerVpcCidrTableListTmp["peer_route_table_id"]
			peerVpcCidrTableListMaps = append(peerVpcCidrTableListMaps, peerVpcCidrTableListMap)
		}
		peerVpcCidrTableListJson, err := json.Marshal(peerVpcCidrTableListMaps)
		if err != nil {
			return WrapError(err)
		}
		request["PeerVpcCidrTableList"] = string(peerVpcCidrTableListJson)
	}

	if update {
		action := "ModifyVpcFirewallConfigure"
		var err error
		var endpoint string
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
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

			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("vpc_firewall_name")
		d.SetPartial("member_uid")
		d.SetPartial("lang")
		d.SetPartial("local_vpc")
		d.SetPartial("peer_vpc")
	}

	update = false
	request = map[string]interface{}{
		"VpcFirewallId": d.Id(),
	}

	if d.HasChange("status") {
		update = true
	}
	request["FirewallSwitch"] = d.Get("status")

	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	if update {
		action := "ModifyVpcFirewallSwitchStatus"
		var err error
		var endpoint string
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
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

			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"closed", "opened"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cloudfwService.CloudFirewallVpcFirewallStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("status")
		d.SetPartial("member_uid")
		d.SetPartial("lang")
	}

	d.Partial(false)

	return resourceAliCloudCloudFirewallVpcFirewallRead(d, meta)
}

func resourceAliCloudCloudFirewallVpcFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}
	action := "DeleteVpcFirewallConfigure"
	var response map[string]interface{}
	var err error
	var endpoint string

	request := map[string]interface{}{
		"VpcFirewallIdList.1": d.Id(),
	}

	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
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

		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cloudfwService.CloudFirewallVpcFirewallStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertCloudFirewallVpcFirewallStatusRequest(source interface{}) interface{} {
	switch source {
	case "closed":
		return "close"
	case "opened":
		return "open"
	}
	return source
}
