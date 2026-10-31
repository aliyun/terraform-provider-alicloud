package cloudfirewall

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	validators "github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	retry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceCloudFirewallVpcCenTrFirewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcCenTrFirewallCreate,
		Read:   resourceVpcCenTrFirewallRead,
		Update: resourceVpcCenTrFirewallUpdate,
		Delete: resourceVpcCenTrFirewallDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(41 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(46 * time.Minute),
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
			"firewall_eni_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"firewall_eni_vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"firewall_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"firewall_vpc_attachment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"firewall_vpc_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"firewall_vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"firewall_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
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
				ValidateFunc: validators.StringInSlice([]string{"managed", "manual", "semi-managed"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tr_attachment_master_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tr_attachment_slave_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"transit_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"firewall_subnet_cidr": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "Field `firewall_subnet_cidr` has been deprecated from provider version 2.2.0.",
			},
			"tr_attachment_master_cidr": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field `tr_attachment_master_cidr` has been deprecated from provider version 2.2.0.",
			},
			"tr_attachment_slave_cidr": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field `tr_attachment_slave_cidr` has been deprecated from provider version 2.2.0.",
			},
		},
	}
}

func resourceVpcCenTrFirewallCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTrFirewallV2"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	var endpoint string
	request = make(map[string]interface{})

	request["CenId"] = d.Get("cen_id")
	if v, ok := d.GetOk("tr_attachment_master_cidr"); ok {
		request["TrAttachmentMasterCidr"] = v
	}
	if v, ok := d.GetOk("firewall_description"); ok {
		request["FirewallDescription"] = v
	}
	if v, ok := d.GetOk("tr_attachment_slave_zone"); ok {
		request["TrAttachmentSlaveZone"] = v
	}
	if v, ok := d.GetOk("firewall_subnet_cidr"); ok {
		request["FirewallSubnetCidr"] = v
	}
	request["RouteMode"] = d.Get("route_mode")
	request["RegionNo"] = d.Get("region_no")
	request["TransitRouterId"] = d.Get("transit_router_id")
	request["FirewallName"] = d.Get("firewall_name")
	if v, ok := d.GetOk("tr_attachment_slave_cidr"); ok {
		request["TrAttachmentSlaveCidr"] = v
	}
	if v, ok := d.GetOk("tr_attachment_master_zone"); ok {
		request["TrAttachmentMasterZone"] = v
	}
	if v, ok := d.GetOk("firewall_vpc_cidr"); ok {
		request["FirewallVpcCidr"] = v
	}
	if v, ok := d.GetOk("firewall_vpc_id"); ok {
		request["FirewallVpcId"] = v
	}
	if v, ok := d.GetOk("firewall_vswitch_id"); ok {
		request["FirewallVswitchId"] = v
	}
	wait := helper.IncrementalWait(30*time.Second, 30*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, query, request, false, endpoint)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ErrorTrResourceNotReady"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			} else if errs.IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_firewall_vpc_cen_tr_firewall", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["FirewallId"]))

	cloudFirewallServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Ready"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, cloudFirewallServiceV2.CloudFirewallVpcCenTrFirewallStateRefreshFunc(d.Id(), "FirewallStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceVpcCenTrFirewallRead(d, meta)
}

func resourceVpcCenTrFirewallRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudFirewallServiceV2 := NewService(client)

	objectRaw, err := cloudFirewallServiceV2.DescribeCloudFirewallVpcCenTrFirewall(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_vpc_cen_tr_firewall DescribeCloudFirewallVpcCenTrFirewall Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("cen_id", objectRaw["CenId"])
	d.Set("firewall_description", objectRaw["FirewallDescription"])
	d.Set("firewall_eni_id", objectRaw["FirewallEniId"])
	d.Set("firewall_eni_vpc_id", objectRaw["FirewallEniVpcId"])
	d.Set("firewall_vpc_id", objectRaw["FirewallEniVpcId"])
	d.Set("firewall_vswitch_id", objectRaw["FirewallEniVswitchId"])
	d.Set("firewall_name", objectRaw["FirewallName"])
	d.Set("firewall_subnet_cidr", objectRaw["FirewallSubnetCidr"])
	d.Set("firewall_vpc_attachment_id", objectRaw["TrAttachmentId"])
	d.Set("firewall_vpc_cidr", objectRaw["FirewallVpcCidr"])
	d.Set("region_no", objectRaw["RegionNo"])
	d.Set("route_mode", objectRaw["RouteMode"])
	d.Set("status", objectRaw["FirewallStatus"])
	d.Set("tr_attachment_master_cidr", objectRaw["TrAttachmentMasterCidr"])
	d.Set("tr_attachment_slave_cidr", objectRaw["TrAttachmentSlaveCidr"])
	d.Set("transit_router_id", objectRaw["TransitRouterId"])
	d.Set("tr_attachment_master_zone", objectRaw["TrAttachmentMasterZone"])
	d.Set("tr_attachment_slave_zone", objectRaw["TrAttachmentSlaveZone"])

	return nil
}

func resourceVpcCenTrFirewallUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyTrFirewallV2Configuration"
	var endpoint string
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FirewallId"] = d.Id()

	if d.HasChange("firewall_name") {
		update = true
	}
	request["FirewallName"] = d.Get("firewall_name")
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, query, request, false, endpoint)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				} else if errs.IsExpectedErrors(err, []string{"not buy user"}) {
					endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceVpcCenTrFirewallRead(d, meta)
}

func resourceVpcCenTrFirewallDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTrFirewallV2"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	var endpoint string
	request = make(map[string]interface{})
	request["FirewallId"] = d.Id()

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, query, request, false, endpoint)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			} else if errs.IsExpectedErrors(err, []string{"not buy user"}) {
				endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ErrorTrFirewallNotExist"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	cloudFirewallServiceV2 := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 30*time.Second, cloudFirewallServiceV2.CloudFirewallVpcCenTrFirewallStateRefreshFunc(d.Id(), "FirewallStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
