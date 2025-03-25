// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCenTransitRouterCidr() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCenTransitRouterCidrCreate,
		Read:   resourceAliCloudCenTransitRouterCidrRead,
		Update: resourceAliCloudCenTransitRouterCidrUpdate,
		Delete: resourceAliCloudCenTransitRouterCidrDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"publish_cidr_route": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"transit_router_cidr_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transit_router_cidr_name": {
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

func resourceAliCloudCenTransitRouterCidrCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateTransitRouterCidr"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("transit_router_id"); ok {
		request["TransitRouterId"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOkExists("publish_cidr_route"); ok {
		request["PublishCidrRoute"] = v
	}
	request["Cidr"] = d.Get("cidr")
	if v, ok := d.GetOk("transit_router_cidr_name"); ok {
		request["Name"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User", "IncorrectStatus.Status"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_transit_router_cidr", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["TransitRouterId"], response["TransitRouterCidrId"]))

	return resourceAliCloudCenTransitRouterCidrRead(d, meta)
}

func resourceAliCloudCenTransitRouterCidrRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenServiceV2 := CenServiceV2{client}

	objectRaw, err := cenServiceV2.DescribeCenTransitRouterCidr(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_transit_router_cidr DescribeCenTransitRouterCidr Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cidr", objectRaw["Cidr"])
	d.Set("description", objectRaw["Description"])
	d.Set("publish_cidr_route", objectRaw["PublishCidrRoute"])
	d.Set("transit_router_cidr_name", objectRaw["Name"])
	d.Set("transit_router_cidr_id", objectRaw["TransitRouterCidrId"])
	d.Set("transit_router_id", objectRaw["TransitRouterId"])

	return nil
}

func resourceAliCloudCenTransitRouterCidrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyTransitRouterCidr"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TransitRouterCidrId"] = parts[1]
	request["TransitRouterId"] = parts[0]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("publish_cidr_route") {
		update = true
		request["PublishCidrRoute"] = d.Get("publish_cidr_route")
	}

	if d.HasChange("cidr") {
		update = true
	}
	request["Cidr"] = d.Get("cidr")
	if d.HasChange("transit_router_cidr_name") {
		update = true
		request["Name"] = d.Get("transit_router_cidr_name")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User", "IncorrectStatus.Status"}) || NeedRetry(err) {
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

	return resourceAliCloudCenTransitRouterCidrRead(d, meta)
}

func resourceAliCloudCenTransitRouterCidrDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteTransitRouterCidr"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["TransitRouterCidrId"] = parts[1]
	request["TransitRouterId"] = parts[0]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Cbn", "2017-09-12", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking", "Throttling.User", "IncorrectStatus.Status", "OperationFailed.CidrBlockAllocated"}) || NeedRetry(err) {
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
