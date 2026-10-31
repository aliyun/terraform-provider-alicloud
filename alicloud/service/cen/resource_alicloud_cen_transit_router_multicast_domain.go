// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cen

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"tags": tags.TagsSchema(),
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
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)
	request["TransitRouterId"] = d.Get("transit_router_id")

	if v, ok := d.GetOk("transit_router_multicast_domain_name"); ok {
		request["TransitRouterMulticastDomainName"] = v
	}
	if v, ok := d.GetOk("transit_router_multicast_domain_description"); ok {
		request["TransitRouterMulticastDomainDescription"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tagsMap := tags.ConvertTags(v.(map[string]interface{}))
		request = tags.ExpandTagsToMap(request, tagsMap)
	}

	if v, ok := d.GetOk("options"); ok {
		jsonPathResult3, err := jsonpath.Get("$[0].igmpv2_support", v)
		if err == nil && jsonPathResult3 != "" {
			request["Options.Igmpv2Support"] = jsonPathResult3
		}
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cen_transit_router_multicast_domain", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TransitRouterMulticastDomainId"]))

	return resourceAliCloudCenTransitRouterMulticastDomainRead(d, meta)
}

func resourceAliCloudCenTransitRouterMulticastDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenServiceV2 := NewCenServiceV2(client)

	objectRaw, err := cenServiceV2.DescribeCenTransitRouterMulticastDomain(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_multicast_domain DescribeCenTransitRouterMulticastDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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
	d.Set("tags", tags.TagsToMap(tagsMaps))

	return nil
}

func resourceAliCloudCenTransitRouterMulticastDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	action := "ModifyTransitRouterMulticastDomain"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TransitRouterMulticastDomainId"] = d.Id()

	request["ClientToken"] = helper.BuildClientToken(action)
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		cenServiceV2 := NewCenServiceV2(client)
		if err := cenServiceV2.SetResourceTags(d, "TRANSITROUTERMULTICASTDOMAIN"); err != nil {
			return errs.WrapError(err)
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
	var err error
	request = make(map[string]interface{})
	request["TransitRouterMulticastDomainId"] = d.Id()

	request["ClientToken"] = helper.BuildClientToken(action)

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		request["ClientToken"] = helper.BuildClientToken(action)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"Operation.Blocking", "IncorrectStatus.MulticastDomain", "InvalidOperation.MulticastGroupExist"}) || helper.NeedRetry(err) {
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

	return nil
}
