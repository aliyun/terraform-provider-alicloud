package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCenTransitRouterMulticastDomainSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenTransitRouterMulticastDomainSourceCreate,
		Read:   resourceAlicloudCenTransitRouterMulticastDomainSourceRead,
		Delete: resourceAlicloudCenTransitRouterMulticastDomainSourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"transit_router_multicast_domain_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"group_ip_address": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"network_interface_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"vpc_id": {
				Optional: true,
				Computed: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudCenTransitRouterMulticastDomainSourceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	request := make(map[string]interface{})
	var err error

	if v, ok := d.GetOk("group_ip_address"); ok {
		request["GroupIpAddress"] = v
	}
	if v, ok := d.GetOk("network_interface_id"); ok {
		request["NetworkInterfaceIds.1"] = v
	}
	if v, ok := d.GetOk("transit_router_multicast_domain_id"); ok {
		request["TransitRouterMulticastDomainId"] = v
	}

	request["ClientToken"] = buildClientToken("RegisterTransitRouterMulticastGroupSources")
	var response map[string]interface{}
	action := "RegisterTransitRouterMulticastGroupSources"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_transit_router_multicast_domain_source", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["TransitRouterMulticastDomainId"], ":", request["GroupIpAddress"], ":", request["NetworkInterfaceIds.1"]))
	stateConf := BuildStateConf([]string{}, []string{"Registered"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenTransitRouterMulticastDomainSourceStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudCenTransitRouterMulticastDomainSourceRead(d, meta)
}

func resourceAlicloudCenTransitRouterMulticastDomainSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	object, err := cbnService.DescribeCenTransitRouterMulticastDomainSource(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_multicast_domain_source cbnService.DescribeCenTransitRouterMulticastDomainSource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("transit_router_multicast_domain_id", parts[0])
	d.Set("group_ip_address", parts[1])
	d.Set("network_interface_id", parts[2])
	d.Set("status", object["Status"])
	d.Set("vpc_id", object["ResourceId"])

	return nil
}

func resourceAlicloudCenTransitRouterMulticastDomainSourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var err error
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"TransitRouterMulticastDomainId": parts[0],
		"GroupIpAddress":                 parts[1],
		"NetworkInterfaceIds.1":          parts[2],
	}

	request["ClientToken"] = buildClientToken("DeregisterTransitRouterMulticastGroupSources")
	action := "DeregisterTransitRouterMulticastGroupSources"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("Cbn", "2017-09-12", action, nil, request, true)
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cbnService.CenTransitRouterMulticastDomainSourceStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
