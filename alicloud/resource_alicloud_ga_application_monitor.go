// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudGaApplicationMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGaApplicationMonitorCreate,
		Read:   resourceAliCloudGaApplicationMonitorRead,
		Update: resourceAliCloudGaApplicationMonitorUpdate,
		Delete: resourceAliCloudGaApplicationMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accelerator_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"detect_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"detect_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"detect_times": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"options_json": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"silence_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliCloudGaApplicationMonitorCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateApplicationMonitor"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewGaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["AcceleratorId"] = d.Get("accelerator_id")
	request["ListenerId"] = d.Get("listener_id")
	request["TaskName"] = d.Get("task_name")
	request["Address"] = d.Get("address")
	if v, ok := d.GetOk("options_json"); ok {
		request["OptionsJson"] = v
	}
	if v, ok := d.GetOkExists("detect_enable"); ok {
		request["DetectEnable"] = v
	}
	if v, ok := d.GetOk("detect_threshold"); ok {
		request["DetectThreshold"] = v
	}
	if v, ok := d.GetOk("detect_times"); ok {
		request["DetectTimes"] = v
	}
	if v, ok := d.GetOk("silence_time"); ok {
		request["SilenceTime"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_application_monitor", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TaskId"]))

	return resourceAliCloudGaApplicationMonitorUpdate(d, meta)
}

func resourceAliCloudGaApplicationMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaServiceV2 := GaServiceV2{client}

	objectRaw, err := gaServiceV2.DescribeGaApplicationMonitor(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_application_monitor DescribeGaApplicationMonitor Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("accelerator_id", objectRaw["AcceleratorId"])
	d.Set("address", objectRaw["Address"])
	d.Set("detect_enable", objectRaw["DetectEnable"])
	d.Set("detect_threshold", objectRaw["DetectThreshold"])
	d.Set("detect_times", objectRaw["DetectTimes"])
	d.Set("listener_id", objectRaw["ListenerId"])
	d.Set("options_json", objectRaw["OptionsJson"])
	d.Set("silence_time", objectRaw["SilenceTime"])
	d.Set("status", objectRaw["State"])
	d.Set("task_name", objectRaw["TaskName"])

	return nil
}

func resourceAliCloudGaApplicationMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	action := "UpdateApplicationMonitor"
	conn, err := client.NewGaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["TaskId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("listener_id") {
		update = true
	}
	request["ListenerId"] = d.Get("listener_id")
	if !d.IsNewResource() && d.HasChange("task_name") {
		update = true
	}
	request["TaskName"] = d.Get("task_name")
	if !d.IsNewResource() && d.HasChange("address") {
		update = true
	}
	request["Address"] = d.Get("address")
	if !d.IsNewResource() && d.HasChange("options_json") {
		update = true
		request["OptionsJson"] = d.Get("options_json")
	}

	if !d.IsNewResource() && d.HasChange("detect_threshold") {
		update = true
		request["DetectThreshold"] = d.Get("detect_threshold")
	}

	if !d.IsNewResource() && d.HasChange("detect_times") {
		update = true
		request["DetectTimes"] = d.Get("detect_times")
	}

	if !d.IsNewResource() && d.HasChange("silence_time") {
		update = true
		request["SilenceTime"] = d.Get("silence_time")
	}

	if !d.IsNewResource() && d.HasChange("detect_enable") {
		update = true
		request["DetectEnable"] = d.Get("detect_enable")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			request["ClientToken"] = buildClientToken(action)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("status") {
		client := meta.(*connectivity.AliyunClient)
		gaServiceV2 := GaServiceV2{client}
		object, err := gaServiceV2.DescribeGaApplicationMonitor(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["State"].(string) != target {
			if target == "Disable" {
				action = "DisableApplicationMonitor"
				conn, err = client.NewGaClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				request["TaskId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					request["ClientToken"] = buildClientToken(action)

					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, response, request)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
			if target == "Enable" {
				action = "EnableApplicationMonitor"
				conn, err = client.NewGaClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				request["TaskId"] = d.Id()
				request["RegionId"] = client.RegionId
				request["ClientToken"] = buildClientToken(action)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
					request["ClientToken"] = buildClientToken(action)

					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(action, response, request)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

			}
		}
	}

	return resourceAliCloudGaApplicationMonitorRead(d, meta)
}

func resourceAliCloudGaApplicationMonitorDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteApplicationMonitor"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewGaClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["TaskId"] = d.Id()
	request["RegionId"] = client.RegionId

	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
