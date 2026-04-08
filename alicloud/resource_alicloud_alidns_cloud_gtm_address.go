// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAlidnsCloudGtmAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAlidnsCloudGtmAddressCreate,
		Read:   resourceAliCloudAlidnsCloudGtmAddressRead,
		Update: resourceAliCloudAlidnsCloudGtmAddressUpdate,
		Delete: resourceAliCloudAlidnsCloudGtmAddressDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"available_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"auto", "manual"}, false),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_status": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"enable", "disable"}, false),
			},
			"health_judgement": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"any_ok", "all_ok", "p30_ok", "p50_ok", "p70_ok"}, false),
			},
			"health_tasks": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"manual_available_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"available", "unavailable"}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"IPv4", "IPv6", "domain"}, false),
			},
		},
	}
}

func resourceAliCloudAlidnsCloudGtmAddressCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateCloudGtmAddress"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("health_tasks"); ok {
		healthTasksMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Port"] = dataLoopTmp["port"]
			dataLoopMap["TemplateId"] = dataLoopTmp["template_id"]
			healthTasksMapsArray = append(healthTasksMapsArray, dataLoopMap)
		}
		healthTasksMapsJson, err := json.Marshal(healthTasksMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["HealthTasks"] = string(healthTasksMapsJson)
	}

	request["HealthJudgement"] = d.Get("health_judgement")
	request["Name"] = d.Get("name")
	request["Address"] = d.Get("address")
	request["EnableStatus"] = d.Get("enable_status")
	if v, ok := d.GetOk("manual_available_status"); ok {
		request["ManualAvailableStatus"] = v
	}
	if v, ok := d.GetOk("remark"); ok {
		request["Remark"] = v
	}
	request["AvailableMode"] = d.Get("available_mode")
	request["Type"] = d.Get("type")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alidns_cloud_gtm_address", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AddressId"]))

	return resourceAliCloudAlidnsCloudGtmAddressRead(d, meta)
}

func resourceAliCloudAlidnsCloudGtmAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsServiceV2 := AlidnsServiceV2{client}

	objectRaw, err := alidnsServiceV2.DescribeAlidnsCloudGtmAddress(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alidns_cloud_gtm_address DescribeAlidnsCloudGtmAddress Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("address", objectRaw["Address"])
	d.Set("available_mode", objectRaw["AvailableMode"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("enable_status", objectRaw["EnableStatus"])
	d.Set("health_judgement", objectRaw["HealthJudgement"])
	d.Set("manual_available_status", objectRaw["ManualAvailableStatus"])
	d.Set("name", objectRaw["Name"])
	d.Set("remark", objectRaw["Remark"])
	d.Set("type", objectRaw["Type"])

	healthTaskRaw, _ := jsonpath.Get("$.HealthTasks.HealthTask", objectRaw)
	healthTasksMaps := make([]map[string]interface{}, 0)
	if healthTaskRaw != nil {
		for _, healthTaskChildRaw := range convertToInterfaceArray(healthTaskRaw) {
			healthTasksMap := make(map[string]interface{})
			healthTaskChildRaw := healthTaskChildRaw.(map[string]interface{})
			healthTasksMap["port"] = healthTaskChildRaw["Port"]
			healthTasksMap["template_id"] = healthTaskChildRaw["TemplateId"]

			healthTasksMaps = append(healthTasksMaps, healthTasksMap)
		}
	}
	if err := d.Set("health_tasks", healthTasksMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudAlidnsCloudGtmAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "UpdateCloudGtmAddress"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AddressId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("health_tasks") {
		update = true
		if v, ok := d.GetOk("health_tasks"); ok || d.HasChange("health_tasks") {
			healthTasksMapsArray := make([]interface{}, 0)
			for _, dataLoop := range convertToInterfaceArray(v) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Port"] = dataLoopTmp["port"]
				dataLoopMap["TemplateId"] = dataLoopTmp["template_id"]
				healthTasksMapsArray = append(healthTasksMapsArray, dataLoopMap)
			}
			healthTasksMapsJson, err := json.Marshal(healthTasksMapsArray)
			if err != nil {
				return WrapError(err)
			}
			request["HealthTasks"] = string(healthTasksMapsJson)
		}
	}

	if d.HasChange("health_judgement") {
		update = true
	}
	request["HealthJudgement"] = d.Get("health_judgement")
	if d.HasChange("name") {
		update = true
	}
	request["Name"] = d.Get("name")
	if d.HasChange("address") {
		update = true
	}
	request["Address"] = d.Get("address")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
	update = false
	action = "UpdateCloudGtmAddressEnableStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AddressId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("enable_status") {
		update = true
	}
	request["EnableStatus"] = d.Get("enable_status")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
	update = false
	action = "UpdateCloudGtmAddressManualAvailableStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AddressId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("manual_available_status") {
		update = true
		request["ManualAvailableStatus"] = d.Get("manual_available_status")
	}

	if d.HasChange("available_mode") {
		update = true
	}
	request["AvailableMode"] = d.Get("available_mode")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
	update = false
	action = "UpdateCloudGtmAddressRemark"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AddressId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("remark") {
		update = true
		request["Remark"] = d.Get("remark")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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

	d.Partial(false)
	return resourceAliCloudAlidnsCloudGtmAddressRead(d, meta)
}

func resourceAliCloudAlidnsCloudGtmAddressDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCloudGtmAddress"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["AddressId"] = d.Id()

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Alidns", "2015-01-09", action, query, request, true)
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
