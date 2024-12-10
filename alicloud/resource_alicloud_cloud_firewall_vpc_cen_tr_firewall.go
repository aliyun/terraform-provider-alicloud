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

func resourceAliCloudCloudFirewallVpcCenTrFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudFirewallVpcCenTrFirewallCreate,
		Read:   resourceAliCloudCloudFirewallVpcCenTrFirewallRead,
		Update: resourceAliCloudCloudFirewallVpcCenTrFirewallUpdate,
		Delete: resourceAliCloudCloudFirewallVpcCenTrFirewallDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cen_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"firewall_description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"firewall_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"firewall_subnet_cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"firewall_vpc_cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_no": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"managed"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tr_attachment_master_cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tr_attachment_master_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tr_attachment_slave_cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tr_attachment_slave_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCloudFirewallVpcCenTrFirewallCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTrFirewallV2"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	var endpoint string
	request = make(map[string]interface{})

	request["FirewallName"] = d.Get("firewall_name")
	request["RouteMode"] = d.Get("route_mode")
	request["TransitRouterId"] = d.Get("transit_router_id")
	request["RegionNo"] = d.Get("region_no")
	request["FirewallVpcCidr"] = d.Get("firewall_vpc_cidr")
	request["FirewallSubnetCidr"] = d.Get("firewall_subnet_cidr")
	request["TrAttachmentSlaveCidr"] = d.Get("tr_attachment_slave_cidr")
	request["TrAttachmentMasterCidr"] = d.Get("tr_attachment_master_cidr")
	request["CenId"] = d.Get("cen_id")
	if v, ok := d.GetOk("firewall_description"); ok {
		request["FirewallDescription"] = v
	}
	if v, ok := d.GetOk("tr_attachment_slave_zone"); ok {
		request["TrAttachmentSlaveZone"] = v
	}
	if v, ok := d.GetOk("tr_attachment_master_zone"); ok {
		request["TrAttachmentMasterZone"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, query, request, false, endpoint)
		if err != nil {
			if IsExpectedErrors(err, []string{"ErrorTrResourceNotReady"}) || NeedRetry(err) {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_vpc_cen_tr_firewall", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["FirewallId"]))

	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Ready"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, cloudFirewallServiceV2.CloudFirewallVpcCenTrFirewallStateRefreshFunc(d.Id(), "FirewallStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCloudFirewallVpcCenTrFirewallRead(d, meta)
}

func resourceAliCloudCloudFirewallVpcCenTrFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallVpcCenTrFirewall(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_vpc_cen_tr_firewall DescribeCloudFirewallVpcCenTrFirewall Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CenId"] != nil {
		d.Set("cen_id", objectRaw["CenId"])
	}
	if objectRaw["FirewallDescription"] != nil {
		d.Set("firewall_description", objectRaw["FirewallDescription"])
	}
	if objectRaw["FirewallName"] != nil {
		d.Set("firewall_name", objectRaw["FirewallName"])
	}
	if objectRaw["FirewallSubnetCidr"] != nil {
		d.Set("firewall_subnet_cidr", objectRaw["FirewallSubnetCidr"])
	}
	if objectRaw["FirewallVpcCidr"] != nil {
		d.Set("firewall_vpc_cidr", objectRaw["FirewallVpcCidr"])
	}
	if objectRaw["RegionNo"] != nil {
		d.Set("region_no", objectRaw["RegionNo"])
	}
	if objectRaw["RouteMode"] != nil {
		d.Set("route_mode", objectRaw["RouteMode"])
	}
	if objectRaw["FirewallStatus"] != nil {
		d.Set("status", objectRaw["FirewallStatus"])
	}
	if objectRaw["TrAttachmentMasterCidr"] != nil {
		d.Set("tr_attachment_master_cidr", objectRaw["TrAttachmentMasterCidr"])
	}
	if objectRaw["TrAttachmentSlaveCidr"] != nil {
		d.Set("tr_attachment_slave_cidr", objectRaw["TrAttachmentSlaveCidr"])
	}
	if objectRaw["TransitRouterId"] != nil {
		d.Set("transit_router_id", objectRaw["TransitRouterId"])
	}

	return nil
}

func resourceAliCloudCloudFirewallVpcCenTrFirewallUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "ModifyTrFirewallV2Configuration"
	var err error
	var endpoint string
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["FirewallId"] = d.Id()

	if d.HasChange("firewall_name") {
		update = true
	}
	request["FirewallName"] = d.Get("firewall_name")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
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
	}

	return resourceAliCloudCloudFirewallVpcCenTrFirewallRead(d, meta)
}

func resourceAliCloudCloudFirewallVpcCenTrFirewallDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTrFirewallV2"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	var endpoint string
	request = make(map[string]interface{})
	query["FirewallId"] = d.Id()

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
		if IsExpectedErrors(err, []string{"ErrorTrFirewallNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	cloudFirewallServiceV2 := CloudFirewallServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, cloudFirewallServiceV2.CloudFirewallVpcCenTrFirewallStateRefreshFunc(d.Id(), "FirewallStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
