package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudFirewallVpcFirewallCen() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallVpcFirewallCenCreate,
		Read:   resourceAliCloudCloudFirewallVpcFirewallCenRead,
		Update: resourceAliCloudCloudFirewallVpcFirewallCenUpdate,
		Delete: resourceAliCloudCloudFirewallVpcFirewallCenDelete,
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
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"member_uid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"local_vpc": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_instance_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"network_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attachment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attachment_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"manual_vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_no": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"support_manual_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"defend_cidr_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"eni_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"eni_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"eni_private_ip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"vpc_cidr_table_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"route_table_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"route_entry_list": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"next_hop_instance_id": {
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
								},
							},
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
		},
	}
}

func resourceAliCloudCloudFirewallVpcFirewallCenCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}
	var response map[string]interface{}
	var err error
	var endpoint string
	action := "CreateVpcFirewallCenConfigure"
	request := make(map[string]interface{})

	request["VpcFirewallName"] = d.Get("vpc_firewall_name")
	request["CenId"] = d.Get("cen_id")
	request["VpcRegion"] = d.Get("vpc_region")
	request["FirewallSwitch"] = d.Get("status")

	if v, ok := d.GetOk("member_uid"); ok {
		request["MemberUid"] = v
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	if v, ok := d.GetOk("local_vpc"); ok {
		networkInstanceId, err := jsonpath.Get("$[0].network_instance_id", v)
		if err != nil {
			return WrapError(err)
		}
		request["NetworkInstanceId"] = networkInstanceId
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, false, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			} else if IsExpectedErrors(err, []string{"ErrorVpcFirewallNotFound"}) {
				if err := cloudfwService.CreateVpcFirewallTask(); err != nil {
					log.Println("[ERROR] syncing cen configure failed by api CreateVpcFirewallTask. Error:", err)
				}
				wait()
				return resource.RetryableError(err)
			} else if IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
				return resource.RetryableError(fmt.Errorf("%s", response))
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_vpc_firewall_cen", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.VpcFirewallId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_cloud_firewall_vpc_firewall_cen")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	stateConf := BuildStateConf([]string{}, []string{"opened", "closed"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cloudfwService.CloudFirewallVpcFirewallCenStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCloudFirewallVpcFirewallCenRead(d, meta)
}

func resourceAliCloudCloudFirewallVpcFirewallCenRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}

	object, err := cloudfwService.DescribeCloudFirewallVpcFirewallCen(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_vpc_firewall_cen cloudfwService.DescribeCloudFirewallVpcFirewallCen Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	objectExtra, err := cloudfwService.DescribeVpcFirewallCenList(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
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

func resourceAliCloudCloudFirewallVpcFirewallCenUpdate(d *schema.ResourceData, meta interface{}) error {
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

	if update {
		action := "ModifyVpcFirewallCenConfigure"
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
		action := "ModifyVpcFirewallCenSwitchStatus"
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

		stateConf := BuildStateConf([]string{}, []string{"opened", "closed"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cloudfwService.CloudFirewallVpcFirewallCenStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("status")
		d.SetPartial("member_uid")
		d.SetPartial("lang")
	}

	d.Partial(false)

	return resourceAliCloudCloudFirewallVpcFirewallCenRead(d, meta)
}

func resourceAliCloudCloudFirewallVpcFirewallCenDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}
	action := "DeleteVpcFirewallCenConfigure"
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
