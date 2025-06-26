// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudThreatDetectionCycleTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionCycleTaskCreate,
		Read:   resourceAliCloudThreatDetectionCycleTaskRead,
		Update: resourceAliCloudThreatDetectionCycleTaskUpdate,
		Delete: resourceAliCloudThreatDetectionCycleTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"enable": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"first_date_str": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"interval_period": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"param": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period_unit": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_end_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"target_start_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"task_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudThreatDetectionCycleTaskCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateCycleTask"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["PeriodUnit"] = d.Get("period_unit")
	request["TaskType"] = d.Get("task_type")
	request["TargetEndTime"] = d.Get("target_end_time")
	if v, ok := d.GetOk("param"); ok {
		request["Param"] = v
	}
	request["TargetStartTime"] = d.Get("target_start_time")
	request["FirstDateStr"] = d.Get("first_date_str")
	if v, ok := d.GetOk("source"); ok {
		request["Source"] = v
	}
	request["Enable"] = d.Get("enable")
	request["IntervalPeriod"] = d.Get("interval_period")
	request["TaskName"] = d.Get("task_name")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_cycle_task", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ConfigId"]))

	return resourceAliCloudThreatDetectionCycleTaskRead(d, meta)
}

func resourceAliCloudThreatDetectionCycleTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionCycleTask(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_cycle_task DescribeThreatDetectionCycleTask Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("enable", objectRaw["Enable"])
	d.Set("first_date_str", objectRaw["FirstDateStr"])
	d.Set("interval_period", objectRaw["IntervalPeriod"])
	d.Set("param", objectRaw["Param"])
	d.Set("period_unit", objectRaw["PeriodUnit"])
	d.Set("target_end_time", objectRaw["TargetEndTime"])
	d.Set("target_start_time", objectRaw["TargetStartTime"])
	d.Set("task_name", objectRaw["TaskName"])
	d.Set("task_type", objectRaw["TaskType"])

	return nil
}

func resourceAliCloudThreatDetectionCycleTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyCycleTask"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ConfigId"] = d.Id()

	if d.HasChange("period_unit") {
		update = true
	}
	request["PeriodUnit"] = d.Get("period_unit")
	if d.HasChange("target_end_time") {
		update = true
	}
	request["TargetEndTime"] = d.Get("target_end_time")
	if d.HasChange("param") {
		update = true
		request["Param"] = d.Get("param")
	}

	if d.HasChange("target_start_time") {
		update = true
	}
	request["TargetStartTime"] = d.Get("target_start_time")
	if d.HasChange("first_date_str") {
		update = true
	}
	request["FirstDateStr"] = d.Get("first_date_str")
	if d.HasChange("enable") {
		update = true
	}
	request["Enable"] = d.Get("enable")
	if d.HasChange("interval_period") {
		update = true
	}
	request["IntervalPeriod"] = d.Get("interval_period")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
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

	return resourceAliCloudThreatDetectionCycleTaskRead(d, meta)
}

func resourceAliCloudThreatDetectionCycleTaskDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCycleTask"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ConfigId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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
