// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudThreatDetectionImageEventOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionImageEventOperationCreate,
		Read:   resourceAliCloudThreatDetectionImageEventOperationRead,
		Update: resourceAliCloudThreatDetectionImageEventOperationUpdate,
		Delete: resourceAliCloudThreatDetectionImageEventOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"conditions": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"event_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"event_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"event_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"sensitiveFile", "maliciousFile"}, false),
			},
			"operation_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"whitelist"}, false),
			},
			"scenarios": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudThreatDetectionImageEventOperationCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "AddImageEventOperation"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("event_type"); ok {
		request["EventType"] = v
	}
	if v, ok := d.GetOk("event_key"); ok {
		request["EventKey"] = v
	}
	if v, ok := d.GetOk("event_name"); ok {
		request["EventName"] = v
	}
	if v, ok := d.GetOk("operation_code"); ok {
		request["OperationCode"] = v
	}
	if v, ok := d.GetOk("conditions"); ok {
		request["Conditions"] = v
	}
	if v, ok := d.GetOk("scenarios"); ok {
		request["Scenarios"] = v
	}
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
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_image_event_operation", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.Id", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudThreatDetectionImageEventOperationRead(d, meta)
}

func resourceAliCloudThreatDetectionImageEventOperationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionImageEventOperation(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_image_event_operation DescribeThreatDetectionImageEventOperation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("conditions", objectRaw["Conditions"])
	d.Set("event_key", objectRaw["EventKey"])
	d.Set("event_name", objectRaw["EventName"])
	d.Set("event_type", objectRaw["EventType"])
	d.Set("operation_code", objectRaw["OperationCode"])
	d.Set("scenarios", objectRaw["Scenarios"])

	return nil
}

func resourceAliCloudThreatDetectionImageEventOperationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "UpdateImageEventOperation"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Id"] = d.Id()
	if d.HasChange("scenarios") {
		update = true
		request["Scenarios"] = d.Get("scenarios")
	}

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
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudThreatDetectionImageEventOperationRead(d, meta)
}

func resourceAliCloudThreatDetectionImageEventOperationDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteImageEventOperation"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["Id"] = d.Id()

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
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
