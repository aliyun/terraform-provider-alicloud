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

func resourceAliCloudEsaRoutineRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaRoutineRouteCreate,
		Read:   resourceAliCloudEsaRoutineRouteRead,
		Update: resourceAliCloudEsaRoutineRouteUpdate,
		Delete: resourceAliCloudEsaRoutineRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bypass": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"fallback": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"route_enable": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"route_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"routine_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sequence": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudEsaRoutineRouteCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateRoutineRoute"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}
	if v, ok := d.GetOk("routine_name"); ok {
		request["RoutineName"] = v
	}

	if v, ok := d.GetOkExists("sequence"); ok {
		request["Sequence"] = v
	}
	if v, ok := d.GetOk("route_name"); ok {
		request["RouteName"] = v
	}
	if v, ok := d.GetOk("route_enable"); ok {
		request["RouteEnable"] = v
	}
	if v, ok := d.GetOk("bypass"); ok {
		request["Bypass"] = v
	}
	if v, ok := d.GetOk("fallback"); ok {
		request["Fallback"] = v
	}
	if v, ok := d.GetOk("rule"); ok {
		request["Rule"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_routine_route", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["SiteId"], request["RoutineName"], response["ConfigId"]))

	return resourceAliCloudEsaRoutineRouteRead(d, meta)
}

func resourceAliCloudEsaRoutineRouteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaRoutineRoute(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_routine_route DescribeEsaRoutineRoute Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("bypass", objectRaw["Bypass"])
	d.Set("fallback", objectRaw["Fallback"])
	d.Set("route_enable", objectRaw["RouteEnable"])
	d.Set("route_name", objectRaw["RouteName"])
	d.Set("rule", objectRaw["Rule"])
	d.Set("sequence", objectRaw["Sequence"])
	d.Set("config_id", objectRaw["ConfigId"])
	d.Set("routine_name", objectRaw["RoutineName"])

	parts := strings.Split(d.Id(), ":")
	d.Set("site_id", fmt.Sprint(parts[0]))

	return nil
}

func resourceAliCloudEsaRoutineRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateRoutineRoute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ConfigId"] = parts[2]
	request["SiteId"] = parts[0]
	request["RoutineName"] = parts[1]

	if d.HasChange("sequence") {
		update = true
		request["Sequence"] = d.Get("sequence")
	}

	if d.HasChange("route_name") {
		update = true
		request["RouteName"] = d.Get("route_name")
	}

	if d.HasChange("route_enable") {
		update = true
		request["RouteEnable"] = d.Get("route_enable")
	}

	if d.HasChange("bypass") {
		update = true
		request["Bypass"] = d.Get("bypass")
	}

	if d.HasChange("fallback") {
		update = true
		request["Fallback"] = d.Get("fallback")
	}

	if d.HasChange("rule") {
		update = true
		request["Rule"] = d.Get("rule")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
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

	return resourceAliCloudEsaRoutineRouteRead(d, meta)
}

func resourceAliCloudEsaRoutineRouteDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteRoutineRoute"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ConfigId"] = parts[2]
	request["SiteId"] = parts[0]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
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
