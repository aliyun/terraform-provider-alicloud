package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

func resourceAlicloudCloudFirewallVpcFirewallCen() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudFirewallVpcFirewallCenCreate,
		Read:   resourceAlicloudCloudFirewallVpcFirewallCenRead,
		Update: resourceAlicloudCloudFirewallVpcFirewallCenUpdate,
		Delete: resourceAlicloudCloudFirewallVpcFirewallCenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(31 * time.Minute),
			Update: schema.DefaultTimeout(31 * time.Minute),
			Delete: schema.DefaultTimeout(31 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"connect_type": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"lang": {
				Optional: true,
				Type:     schema.TypeString,
				Computed: true,
			},
			"local_vpc": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeList,
				MaxItems: 1,
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
							Required: true,
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
										Type:     schema.TypeSet,
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
			"member_uid": {
				Optional: true,
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
			"vpc_region": {
				Required: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudCloudFirewallVpcFirewallCenCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}
	request := make(map[string]interface{})
	conn, err := client.NewCloudfwClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("cen_id"); ok {
		request["CenId"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["FirewallSwitch"] = v
	}
	if v, ok := d.GetOk("vpc_firewall_name"); ok {
		request["VpcFirewallName"] = v
	}
	if v, ok := d.GetOk("vpc_region"); ok {
		request["VpcRegion"] = v
	}
	if v, ok := d.GetOk("local_vpc"); ok {
		networkInstanceId, err := jsonpath.Get("$[0].network_instance_id", v)
		if err != nil {
			return WrapError(err)
		}
		request["NetworkInstanceId"] = networkInstanceId
	}

	var response map[string]interface{}
	action := "CreateVpcFirewallCenConfigure"
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_vpc_firewall_cen", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.VpcFirewallId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_cloud_firewall_vpc_firewall_cen")
	} else {
		d.SetId(fmt.Sprint(v))
	}
	stateConf := BuildStateConf([]string{}, []string{"closed", "opened"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cloudfwService.CloudFirewallVpcFirewallCenStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudCloudFirewallVpcFirewallCenRead(d, meta)
}

func resourceAlicloudCloudFirewallVpcFirewallCenRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}

	object, err := cloudfwService.DescribeCloudFirewallVpcFirewallCen(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_vpc_firewall_cen cloudfwService.DescribeCloudFirewallVpcFirewallCen Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	objectExtra, err := cloudfwService.DescribeVpcFirewallCenList(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_vpc_firewall_cen cloudfwService.DescribeVpcFirewallCenList Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("member_uid", objectExtra["MemberUid"])
	d.Set("cen_id", objectExtra["CenId"])
	vpcRegion, err := jsonpath.Get("$.LocalVpc.RegionNo", objectExtra)
	if err != nil {
		return WrapError(err)
	}
	d.Set("vpc_region", vpcRegion)
	d.Set("connect_type", object["ConnectType"])
	d.Set("id", object["VpcFirewallId"])
	localVpcMaps := make([]map[string]interface{}, 0)
	localVpcMap := make(map[string]interface{})
	localVpcRaw := object["LocalVpc"].(map[string]interface{})
	localVpcMap["attachment_id"] = localVpcRaw["AttachmentId"]
	localVpcMap["attachment_name"] = localVpcRaw["AttachmentName"]
	localVpcMap["defend_cidr_list"] = localVpcRaw["DefendCidrList"].([]interface{})
	eniListMaps := make([]map[string]interface{}, 0)
	eniListRaw := localVpcRaw["EniList"]
	for _, value1 := range eniListRaw.([]interface{}) {
		eniList := value1.(map[string]interface{})
		eniListMap := make(map[string]interface{})
		eniListMap["eni_id"] = eniList["EniId"]
		eniListMap["eni_private_ip_address"] = eniList["EniPrivateIpAddress"]
		eniListMaps = append(eniListMaps, eniListMap)
	}
	localVpcMap["eni_list"] = eniListMaps
	localVpcMap["manual_vswitch_id"] = localVpcRaw["ManualVSwitchId"]
	localVpcMap["network_instance_id"] = localVpcRaw["NetworkInstanceId"]
	localVpcMap["network_instance_name"] = localVpcRaw["NetworkInstanceName"]
	localVpcMap["network_instance_type"] = localVpcRaw["NetworkInstanceType"]
	localVpcMap["owner_id"] = localVpcRaw["OwnerId"]
	localVpcMap["region_no"] = localVpcRaw["RegionNo"]
	localVpcMap["route_mode"] = localVpcRaw["RouteMode"]
	localVpcMap["support_manual_mode"] = localVpcRaw["SupportManualMode"]
	localVpcMap["transit_router_id"] = localVpcRaw["TransitRouterId"]
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
	d.Set("local_vpc", localVpcMaps)
	d.Set("status", convertCloudFirewallVpcFirewallCenStatusRequest(object["FirewallSwitchStatus"]))
	d.Set("vpc_firewall_name", object["VpcFirewallName"])

	return nil
}

func resourceAlicloudCloudFirewallVpcFirewallCenUpdate(d *schema.ResourceData, meta interface{}) error {
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
	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}
	if !d.IsNewResource() && d.HasChange("vpc_firewall_name") {
		update = true
	}
	request["VpcFirewallName"] = d.Get("vpc_firewall_name")

	if update {
		action := "ModifyVpcFirewallCenConfigure"
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
		d.SetPartial("member_uid")
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
		action := "ModifyVpcFirewallCenSwitchStatus"
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
		stateConf := BuildStateConf([]string{}, []string{"opened", "closed"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cloudfwService.CloudFirewallVpcFirewallCenStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("lang")
		d.SetPartial("member_uid")
		d.SetPartial("status")
	}

	d.Partial(false)
	return resourceAlicloudCloudFirewallVpcFirewallCenRead(d, meta)
}

func resourceAlicloudCloudFirewallVpcFirewallCenDelete(d *schema.ResourceData, meta interface{}) error {
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

	action := "DeleteVpcFirewallCenConfigure"
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cloudfwService.CloudFirewallVpcFirewallCenStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
func convertCloudFirewallVpcFirewallCenStatusRequest(source interface{}) interface{} {
	switch source {
	case "closed":
		return "close"
	case "opened":
		return "open"
	}
	return source
}
