package alicloud

import (
	"fmt"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"time"
)

func resourceAliCloudCenTransitRouterMulticastDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCenTransitRouterMulticastDomainCreate,
		Read:   resourceAliCloudCenTransitRouterMulticastDomainRead,
		Update: resourceAliCloudCenTransitRouterMulticastDomainUpdate,
		Delete: resourceAliCloudCenTransitRouterMulticastDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"transit_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_router_multicast_domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_router_multicast_domain_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCenTransitRouterMulticastDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var response map[string]interface{}
	action := "CreateTransitRouterMulticastDomain"
	request := make(map[string]interface{})
	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateTransitRouterMulticastDomain")
	request["TransitRouterId"] = d.Get("transit_router_id")

	if v, ok := d.GetOk("transit_router_multicast_domain_name"); ok {
		request["TransitRouterMulticastDomainName"] = v
	}

	if v, ok := d.GetOk("transit_router_multicast_domain_description"); ok {
		request["TransitRouterMulticastDomainDescription"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_transit_router_multicast_domain", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TransitRouterMulticastDomainId"]))

	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenTransitRouterMulticastDomainStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudCenTransitRouterMulticastDomainUpdate(d, meta)
}

func resourceAliCloudCenTransitRouterMulticastDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	object, err := cbnService.DescribeCenTransitRouterMulticastDomain(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("transit_router_id", object["TransitRouterId"])
	d.Set("transit_router_multicast_domain_name", object["TransitRouterMulticastDomainName"])
	d.Set("transit_router_multicast_domain_description", object["TransitRouterMulticastDomainDescription"])
	d.Set("status", object["Status"])

	listTagResourcesObject, err := cbnService.ListTagResources(d.Id(), "TransitRouterMulticastDomain")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	return nil
}

func resourceAliCloudCenTransitRouterMulticastDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	var response map[string]interface{}
	update := false
	d.Partial(true)

	if d.HasChange("tags") {
		if err := cbnService.SetResourceTags(d, "TransitRouterMulticastDomain"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	request := map[string]interface{}{
		"ClientToken":                    buildClientToken("ModifyTransitRouterMulticastDomain"),
		"TransitRouterMulticastDomainId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("transit_router_multicast_domain_name") {
		update = true
	}
	if v, ok := d.GetOk("transit_router_multicast_domain_name"); ok {
		request["TransitRouterMulticastDomainName"] = v
	}

	if !d.IsNewResource() && d.HasChange("transit_router_multicast_domain_description") {
		update = true
	}
	if v, ok := d.GetOk("transit_router_multicast_domain_description"); ok {
		request["TransitRouterMulticastDomainDescription"] = v
	}

	if update {
		action := "ModifyTransitRouterMulticastDomain"
		conn, err := client.NewCbnClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"Operation.Blocking"}) || NeedRetry(err) {
					wait()
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

		stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, cbnService.CenTransitRouterMulticastDomainStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("transit_router_multicast_domain_name")
		d.SetPartial("transit_router_multicast_domain_description")
	}

	d.Partial(false)

	return resourceAliCloudCenTransitRouterMulticastDomainRead(d, meta)
}

func resourceAliCloudCenTransitRouterMulticastDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	action := "DeleteTransitRouterMulticastDomain"
	var response map[string]interface{}

	conn, err := client.NewCbnClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"ClientToken":                    buildClientToken("DeleteTransitRouterMulticastDomain"),
		"TransitRouterMulticastDomainId": d.Id(),
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.MulticastDomain", "InvalidOperation.MulticastGroupExist"}) || NeedRetry(err) {
				wait()
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

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cbnService.CenTransitRouterMulticastDomainStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
