// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"igmpv2_support": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"transit_router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"transit_router_multicast_domain_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_router_multicast_domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudCenTransitRouterMulticastDomainCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTransitRouterMulticastDomain"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewCenClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	request["TransitRouterId"] = d.Get("transit_router_id")

	if v, ok := d.GetOk("transit_router_multicast_domain_name"); ok {
		request["TransitRouterMulticastDomainName"] = v
	}
	if v, ok := d.GetOk("transit_router_multicast_domain_description"); ok {
		request["TransitRouterMulticastDomainDescription"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := ConvertTags(v.(map[string]interface{}))
		request = expandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("options"); ok {
		jsonPathResult3, err := jsonpath.Get("$[0].igmpv2_support", v)
		if err == nil && jsonPathResult3 != "" {
			request["Options.Igmpv2Support"] = jsonPathResult3
		}
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), query, request, &runtime)
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

	return resourceAliCloudCenTransitRouterMulticastDomainRead(d, meta)
}

func resourceAliCloudCenTransitRouterMulticastDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenServiceV2 := CenServiceV2{client}

	objectRaw, err := cenServiceV2.DescribeCenTransitRouterMulticastDomain(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_multicast_domain DescribeCenTransitRouterMulticastDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["RegionId"] != nil {
		d.Set("region_id", objectRaw["RegionId"])
	}
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["TransitRouterId"] != nil {
		d.Set("transit_router_id", objectRaw["TransitRouterId"])
	}
	if objectRaw["TransitRouterMulticastDomainDescription"] != nil {
		d.Set("transit_router_multicast_domain_description", objectRaw["TransitRouterMulticastDomainDescription"])
	}
	if objectRaw["TransitRouterMulticastDomainName"] != nil {
		d.Set("transit_router_multicast_domain_name", objectRaw["TransitRouterMulticastDomainName"])
	}

	optionsMaps := make([]map[string]interface{}, 0)
	optionsMap := make(map[string]interface{})
	options1Raw := make(map[string]interface{})
	if objectRaw["Options"] != nil {
		options1Raw = objectRaw["Options"].(map[string]interface{})
	}
	if len(options1Raw) > 0 {
		optionsMap["igmpv2_support"] = options1Raw["Igmpv2Support"]

		optionsMaps = append(optionsMaps, optionsMap)
	}
	if objectRaw["Options"] != nil {
		if err := d.Set("options", optionsMaps); err != nil {
			return err
		}
	}
	tagsMaps := objectRaw["Tags"]
	d.Set("tags", tagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudCenTransitRouterMulticastDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	action := "ModifyTransitRouterMulticastDomain"
	conn, err := client.NewCenClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TransitRouterMulticastDomainId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("transit_router_multicast_domain_name") {
		update = true
	}
	if v, ok := d.GetOk("transit_router_multicast_domain_name"); ok {
		request["TransitRouterMulticastDomainName"] = v
	}

	if d.HasChange("transit_router_multicast_domain_description") {
		update = true
	}
	if v, ok := d.GetOk("transit_router_multicast_domain_description"); ok {
		request["TransitRouterMulticastDomainDescription"] = v
	}

	if d.HasChange("options.0.igmpv2_support") {
		update = true
	}
	if v, ok := d.GetOk("options"); ok {
		jsonPathResult3, err := jsonpath.Get("$[0].igmpv2_support", v)
		if err == nil && jsonPathResult3 != "" {
			request["Options.Igmpv2Support"] = jsonPathResult3
		}
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), query, request, &runtime)
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
	}

	if d.HasChange("tags") {
		cenServiceV2 := CenServiceV2{client}
		if err := cenServiceV2.SetResourceTags(d, "TRANSITROUTERMULTICASTDOMAIN"); err != nil {
			return WrapError(err)
		}
	}
	return resourceAliCloudCenTransitRouterMulticastDomainRead(d, meta)
}

func resourceAliCloudCenTransitRouterMulticastDomainDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTransitRouterMulticastDomain"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewCenClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["TransitRouterMulticastDomainId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-09-12"), StringPointer("AK"), query, request, &runtime)
		request["ClientToken"] = buildClientToken(action)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
