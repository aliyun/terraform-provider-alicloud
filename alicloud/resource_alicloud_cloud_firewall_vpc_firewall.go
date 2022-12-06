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

func resourceAlicloudCloudFirewallVpcFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudFirewallVpcFirewallCreate,
		Read:   resourceAlicloudCloudFirewallVpcFirewallRead,
		Update: resourceAlicloudCloudFirewallVpcFirewallUpdate,
		Delete: resourceAlicloudCloudFirewallVpcFirewallDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(31 * time.Minute),
			Update: schema.DefaultTimeout(31 * time.Minute),
			Delete: schema.DefaultTimeout(31 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Computed: true,
				Type:     schema.TypeInt,
			},
			"connect_type": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"lang": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"local_vpc": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeList,
				MaxItems: 1,
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
						"region_no": {
							Required: true,
							Type:     schema.TypeString,
						},
						"router_interface_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"vpc_id": {
							Required: true,
							Type:     schema.TypeString,
						},
						"vpc_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"local_vpc_cidr_table_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"local_route_table_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"local_route_entry_list": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"local_next_hop_instance_id": {
													Type:     schema.TypeString,
													Required: true,
												},
												"local_destination_cidr": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"member_uid": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"peer_vpc": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeList,
				MaxItems: 1,
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
						"region_no": {
							Required: true,
							Type:     schema.TypeString,
						},
						"router_interface_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"vpc_id": {
							Required: true,
							Type:     schema.TypeString,
						},
						"vpc_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"peer_vpc_cidr_table_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"peer_route_table_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"peer_route_entry_list": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"peer_destination_cidr": {
													Type:     schema.TypeString,
													Required: true,
												},
												"peer_next_hop_instance_id": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"region_status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"status": {
				Required: true,
				Type:     schema.TypeString,
			},
			"vpc_firewall_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"vpc_firewall_name": {
				Required: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudCloudFirewallVpcFirewallCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}
	request := make(map[string]interface{})
	conn, err := client.NewCloudfwClient()
	if err != nil {
		return WrapError(err)
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
	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
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
	if v, ok := d.GetOk("status"); ok {
		request["FirewallSwitch"] = v
	}
	if v, ok := d.GetOk("vpc_firewall_name"); ok {
		request["VpcFirewallName"] = v
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

	var response map[string]interface{}
	action := "CreateVpcFirewallConfigure"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudCloudFirewallVpcFirewallRead(d, meta)
}

func resourceAlicloudCloudFirewallVpcFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}

	object, err := cloudfwService.DescribeCloudFirewallVpcFirewall(d.Id())
	if err != nil {
		if NotFoundError(err) {
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

func resourceAlicloudCloudFirewallVpcFirewallUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	conn, err := client.NewCloudfwClient()
	if err != nil {
		return WrapError(err)
	}
	cloudfwService := CloudfwService{client}
	d.Partial(true)
	update := false
	request := map[string]interface{}{
		"VpcFirewallId": d.Id(),
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
	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
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
	if !d.IsNewResource() && d.HasChange("vpc_firewall_name") {
		update = true
	}
	request["VpcFirewallName"] = d.Get("vpc_firewall_name")

	if update {
		action := "ModifyVpcFirewallConfigure"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("lang")
		d.SetPartial("local_vpc")
		d.SetPartial("member_uid")
		d.SetPartial("peer_vpc")
		d.SetPartial("vpc_firewall_name")
	}

	update = false
	request = map[string]interface{}{
		"VpcFirewallId": d.Id(),
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}
	if !d.IsNewResource() && d.HasChange("status") {
		update = true
	}
	request["FirewallSwitch"] = d.Get("status")

	if update {
		action := "ModifyVpcFirewallSwitchStatus"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"closed", "opened"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cloudfwService.CloudFirewallVpcFirewallStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("lang")
		d.SetPartial("member_uid")
		d.SetPartial("status")
	}

	d.Partial(false)
	return resourceAlicloudCloudFirewallVpcFirewallRead(d, meta)
}

func resourceAlicloudCloudFirewallVpcFirewallDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}
	conn, err := client.NewCloudfwClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"VpcFirewallIdList.1": d.Id(),
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}

	action := "DeleteVpcFirewallConfigure"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
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
