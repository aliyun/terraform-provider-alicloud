// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudFirewallNatFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallNatFirewallCreate,
		Read:   resourceAliCloudCloudFirewallNatFirewallRead,
		Update: resourceAliCloudCloudFirewallNatFirewallUpdate,
		Delete: resourceAliCloudCloudFirewallNatFirewallDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"firewall_switch": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nat_route_entry_list": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nexthop_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"route_table_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"nexthop_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"destination_cidr": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"proxy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"normal", "opening", "closed", "open", "close", "configuring", "closing"}, false),
			},
			"strict_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_auto": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudCloudFirewallNatFirewallCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateSecurityProxy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	var endpoint string
	request = make(map[string]interface{})

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	request["ProxyName"] = d.Get("proxy_name")
	request["RegionNo"] = d.Get("region_no")
	request["NatGatewayId"] = d.Get("nat_gateway_id")
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("vswitch_auto"); ok {
		request["VswitchAuto"] = v
	}
	if v, ok := d.GetOk("nat_route_entry_list"); ok {
		natRouteEntryListMaps := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["NextHopId"] = dataLoopTmp["nexthop_id"]
			dataLoopMap["DestinationCidr"] = dataLoopTmp["destination_cidr"]
			dataLoopMap["NextHopType"] = dataLoopTmp["nexthop_type"]
			dataLoopMap["RouteTableId"] = dataLoopTmp["route_table_id"]
			natRouteEntryListMaps = append(natRouteEntryListMaps, dataLoopMap)
		}
		request["NatRouteEntryList"] = natRouteEntryListMaps
	}

	if v, ok := d.GetOk("firewall_switch"); ok {
		request["FirewallSwitch"] = v
	}
	if v, ok := d.GetOk("strict_mode"); ok {
		request["StrictMode"] = v
	}
	if v, ok := d.GetOk("vswitch_cidr"); ok {
		request["VswitchCidr"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VswitchId"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, query, request, false, endpoint)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"-360838", "-360809", "-360157", "-360839"}) {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_nat_firewall", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ProxyId"]))

	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"closed", "normal"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, cloudFirewallServiceV2.CloudFirewallNatFirewallStateRefreshFunc(d.Id(), "ProxyStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCloudFirewallNatFirewallUpdate(d, meta)
}

func resourceAliCloudCloudFirewallNatFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallNatFirewall(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_nat_firewall DescribeCloudFirewallNatFirewall Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("nat_gateway_id", objectRaw["NatGatewayId"])
	d.Set("proxy_name", objectRaw["ProxyName"])
	d.Set("region_no", objectRaw["RegionId"])
	d.Set("status", objectRaw["ProxyStatus"])
	d.Set("strict_mode", objectRaw["StrictMode"])
	d.Set("vpc_id", objectRaw["VpcId"])

	natRouteEntryList1Raw := objectRaw["NatRouteEntryList"]
	natRouteEntryListMaps := make([]map[string]interface{}, 0)
	if natRouteEntryList1Raw != nil {
		for _, natRouteEntryListChild1Raw := range natRouteEntryList1Raw.([]interface{}) {
			natRouteEntryListMap := make(map[string]interface{})
			natRouteEntryListChild1Raw := natRouteEntryListChild1Raw.(map[string]interface{})
			natRouteEntryListMap["destination_cidr"] = natRouteEntryListChild1Raw["DestinationCidr"]
			natRouteEntryListMap["nexthop_id"] = natRouteEntryListChild1Raw["NextHopId"]
			natRouteEntryListMap["nexthop_type"] = natRouteEntryListChild1Raw["NextHopType"]
			natRouteEntryListMap["route_table_id"] = natRouteEntryListChild1Raw["RouteTableId"]

			natRouteEntryListMaps = append(natRouteEntryListMaps, natRouteEntryListMap)
		}
	}
	d.Set("nat_route_entry_list", natRouteEntryListMaps)

	return nil
}

func resourceAliCloudCloudFirewallNatFirewallUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "SwitchSecurityProxy"
	var err error
	var endpoint string
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ProxyId"] = d.Id()
	if d.HasChange("status") {
		update = true
	}
	if v, ok := d.GetOk("status"); ok {
		request["Switch"] = convertCloudFirewallNatFirewallSwitchRequest(v.(string))
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, query, request, false, endpoint)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"normal", "closed"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cloudFirewallServiceV2.CloudFirewallNatFirewallStateRefreshFunc(d.Id(), "ProxyStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudCloudFirewallNatFirewallRead(d, meta)
}

func resourceAliCloudCloudFirewallNatFirewallDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSecurityProxy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	var endpoint string
	request = make(map[string]interface{})
	query["ProxyId"] = d.Id()

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, query, request, false, endpoint)

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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cloudFirewallServiceV2.CloudFirewallNatFirewallStateRefreshFunc(d.Id(), "ProxyStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}

func convertCloudFirewallNatFirewallSwitchRequest(source interface{}) interface{} {
	switch source {
	case "normal":
		return "open"
	case "closed":
		return "close"
	}
	return source
}
