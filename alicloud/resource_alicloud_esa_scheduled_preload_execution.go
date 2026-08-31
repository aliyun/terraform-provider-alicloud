// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEsaScheduledPreloadExecution() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaScheduledPreloadExecutionCreate,
		Read:   resourceAliCloudEsaScheduledPreloadExecutionRead,
		Update: resourceAliCloudEsaScheduledPreloadExecutionUpdate,
		Delete: resourceAliCloudEsaScheduledPreloadExecutionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"interval": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"scheduled_preload_execution_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scheduled_preload_job_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"slice_len": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEsaScheduledPreloadExecutionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateScheduledPreloadExecutions"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("scheduled_preload_job_id"); ok {
		request["Id"] = v
	}
	request["RegionId"] = client.RegionId

	objectDataLocalMap := make(map[string]interface{})

	if v, ok := d.GetOkExists("interval"); ok {
		objectDataLocalMap["Interval"] = v
	}

	if v, ok := d.GetOkExists("start_time"); ok {
		objectDataLocalMap["StartTime"] = v
	}

	if v, ok := d.GetOkExists("slice_len"); ok {
		objectDataLocalMap["SliceLen"] = v
	}

	if v, ok := d.GetOkExists("end_time"); ok {
		objectDataLocalMap["EndTime"] = v
	}

	ExecutionsMap := make([]interface{}, 0)
	ExecutionsMap = append(ExecutionsMap, objectDataLocalMap)
	objectDataLocalMapJson, err := json.Marshal(ExecutionsMap)
	if err != nil {
		return WrapError(err)
	}
	request["Executions"] = string(objectDataLocalMapJson)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_scheduled_preload_execution", action, AlibabaCloudSdkGoERROR)
	}

	SuccessExecutionsJobIdVar, _ := jsonpath.Get("$.SuccessExecutions[0].JobId", response)
	SuccessExecutionsIdVar, _ := jsonpath.Get("$.SuccessExecutions[0].Id", response)
	d.SetId(fmt.Sprintf("%v:%v", SuccessExecutionsJobIdVar, SuccessExecutionsIdVar))

	return resourceAliCloudEsaScheduledPreloadExecutionRead(d, meta)
}

func resourceAliCloudEsaScheduledPreloadExecutionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaScheduledPreloadExecution(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_scheduled_preload_execution DescribeEsaScheduledPreloadExecution Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("end_time", objectRaw["EndTime"])
	d.Set("interval", objectRaw["Interval"])
	d.Set("slice_len", objectRaw["SliceLen"])
	d.Set("start_time", objectRaw["StartTime"])
	d.Set("status", objectRaw["Status"])
	d.Set("scheduled_preload_execution_id", objectRaw["Id"])
	d.Set("scheduled_preload_job_id", objectRaw["JobId"])

	return nil
}

func resourceAliCloudEsaScheduledPreloadExecutionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateScheduledPreloadExecution"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Id"] = parts[1]
	request["RegionId"] = client.RegionId
	if d.HasChange("interval") {
		update = true
	}
	request["Interval"] = d.Get("interval")
	if d.HasChange("start_time") {
		update = true
		request["StartTime"] = d.Get("start_time")
	}

	if d.HasChange("slice_len") {
		update = true
	}
	request["SliceLen"] = d.Get("slice_len")
	if d.HasChange("end_time") {
		update = true
		request["EndTime"] = d.Get("end_time")
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

	return resourceAliCloudEsaScheduledPreloadExecutionRead(d, meta)
}

func resourceAliCloudEsaScheduledPreloadExecutionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteScheduledPreloadExecution"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Id"] = parts[1]
	request["RegionId"] = client.RegionId

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
