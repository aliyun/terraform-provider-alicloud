package cen

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudCenTransitRouterMulticastDomainMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCenTransitRouterMulticastDomainMemberCreate,
		Read:   resourceAliCloudCenTransitRouterMulticastDomainMemberRead,
		Update: resourceAliCloudCenTransitRouterMulticastDomainMemberUpdate,
		Delete: resourceAliCloudCenTransitRouterMulticastDomainMemberDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"transit_router_multicast_domain_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_ip_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCenTransitRouterMulticastDomainMemberCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := NewCbnService(client)
	var response map[string]interface{}
	action := "RegisterTransitRouterMulticastGroupMembers"
	request := make(map[string]interface{})
	var err error

	request["ClientToken"] = helper.BuildClientToken("RegisterTransitRouterMulticastGroupMembers")
	request["TransitRouterMulticastDomainId"] = d.Get("transit_router_multicast_domain_id")
	request["GroupIpAddress"] = d.Get("group_ip_address")

	networkInterfaceId := d.Get("network_interface_id").(string)
	request["NetworkInterfaceIds"] = []string{networkInterfaceId}

	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cen_transit_router_multicast_domain_member", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["TransitRouterMulticastDomainId"], request["GroupIpAddress"], networkInterfaceId))

	stateConf := helper.BuildStateConf([]string{}, []string{"Registered"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenTransitRouterMulticastDomainMemberStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudCenTransitRouterMulticastDomainMemberRead(d, meta)
}

func resourceAliCloudCenTransitRouterMulticastDomainMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := NewCbnService(client)

	object, err := cbnService.DescribeCenTransitRouterMulticastDomainMember(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_multicast_domain_member cbnService.DescribeCenTransitRouterMulticastDomainMember Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("transit_router_multicast_domain_id", object["TransitRouterMulticastDomainId"])
	d.Set("group_ip_address", object["GroupIpAddress"])
	d.Set("network_interface_id", object["NetworkInterfaceId"])
	d.Set("vpc_id", object["ResourceId"])
	d.Set("status", object["Status"])

	return nil
}

func resourceAliCloudCenTransitRouterMulticastDomainMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println(fmt.Sprintf("[WARNING] The resouce has not update operation."))
	return resourceAliCloudCenTransitRouterMulticastDomainMemberRead(d, meta)
}

func resourceAliCloudCenTransitRouterMulticastDomainMemberDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := NewCbnService(client)
	action := "DeregisterTransitRouterMulticastGroupMembers"
	var response map[string]interface{}

	var err error

	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"ClientToken":                    helper.BuildClientToken("DeregisterTransitRouterMulticastGroupMembers"),
		"TransitRouterMulticastDomainId": parts[0],
		"GroupIpAddress":                 parts[1],
		"NetworkInterfaceIds":            []string{parts[2]},
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cbnService.CenTransitRouterMulticastDomainMemberStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
