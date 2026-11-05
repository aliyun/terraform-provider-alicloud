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

func resourceAliCloudCenTransitRouterMulticastDomainPeerMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCenTransitRouterMulticastDomainPeerMemberCreate,
		Read:   resourceAliCloudCenTransitRouterMulticastDomainPeerMemberRead,
		Update: resourceAliCloudCenTransitRouterMulticastDomainPeerMemberUpdate,
		Delete: resourceAliCloudCenTransitRouterMulticastDomainPeerMemberDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"group_ip_address": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"peer_transit_router_multicast_domain_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"transit_router_multicast_domain_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAliCloudCenTransitRouterMulticastDomainPeerMemberCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := NewCbnService(client)
	request := make(map[string]interface{})
	var err error

	if v, ok := d.GetOk("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("group_ip_address"); ok {
		request["GroupIpAddress"] = v
	}
	if v, ok := d.GetOk("peer_transit_router_multicast_domain_id"); ok {
		request["PeerTransitRouterMulticastDomains.1"] = v
	}
	if v, ok := d.GetOk("transit_router_multicast_domain_id"); ok {
		request["TransitRouterMulticastDomainId"] = v
	}

	request["ClientToken"] = helper.BuildClientToken("RegisterTransitRouterMulticastGroupMembers")
	var response map[string]interface{}
	action := "RegisterTransitRouterMulticastGroupMembers"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		resp, err := client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cen_transit_router_multicast_domain_peer_member", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["TransitRouterMulticastDomainId"], ":", request["GroupIpAddress"], ":", request["PeerTransitRouterMulticastDomains.1"]))
	stateConf := helper.BuildStateConf([]string{}, []string{"Registered"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenTransitRouterMulticastDomainPeerMemberStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceAliCloudCenTransitRouterMulticastDomainPeerMemberRead(d, meta)
}

func resourceAliCloudCenTransitRouterMulticastDomainPeerMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAliCloudCenTransitRouterMulticastDomainPeerMemberRead(d, meta)
}

func resourceAliCloudCenTransitRouterMulticastDomainPeerMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := NewCbnService(client)

	object, err := cbnService.DescribeCenTransitRouterMulticastDomainPeerMember(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_multicast_domain_peer_member cbnService.DescribeCenTransitRouterMulticastDomainPeerMember Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("transit_router_multicast_domain_id", parts[0])
	d.Set("group_ip_address", parts[1])
	d.Set("peer_transit_router_multicast_domain_id", parts[2])
	d.Set("status", object["Status"])

	return nil
}

func resourceAliCloudCenTransitRouterMulticastDomainPeerMemberDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := NewCbnService(client)
	var err error
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"TransitRouterMulticastDomainId": parts[0], "GroupIpAddress": parts[1], "PeerTransitRouterMulticastDomains.1": parts[2],
	}

	if v, ok := d.GetOk("dry_run"); ok {
		request["DryRun"] = v
	}

	request["ClientToken"] = helper.BuildClientToken("DeregisterTransitRouterMulticastGroupMembers")
	action := "DeregisterTransitRouterMulticastGroupMembers"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		resp, err := client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cbnService.CenTransitRouterMulticastDomainPeerMemberStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
