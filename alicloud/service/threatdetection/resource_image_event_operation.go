// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package threatdetection

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Required: true,
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
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"sensitiveFile", "maliciousFile", "buildRisk"}, false),
			},
			"note": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"operation_code": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"whitelist"}, false),
			},
			"scenarios": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"agentless", "default"}, false),
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

	if v, ok := d.GetOk("event_name"); ok {
		request["EventName"] = v
	}
	if v, ok := d.GetOk("event_key"); ok {
		request["EventKey"] = v
	}
	request["OperationCode"] = d.Get("operation_code")
	request["EventType"] = d.Get("event_type")
	if v, ok := d.GetOk("note"); ok {
		request["Note"] = v
	}
	request["Conditions"] = d.Get("conditions")
	if v, ok := d.GetOk("source"); ok {
		request["Source"] = v
	}
	if v, ok := d.GetOk("scenarios"); ok {
		request["Scenarios"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_threat_detection_image_event_operation", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.Id", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudThreatDetectionImageEventOperationRead(d, meta)
}

func resourceAliCloudThreatDetectionImageEventOperationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionImageEventOperation(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_image_event_operation DescribeThreatDetectionImageEventOperation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("conditions", objectRaw["Conditions"])
	d.Set("event_key", objectRaw["EventKey"])
	d.Set("event_name", objectRaw["EventName"])
	d.Set("event_type", objectRaw["EventType"])
	d.Set("note", objectRaw["Note"])
	d.Set("operation_code", objectRaw["OperationCode"])
	d.Set("scenarios", objectRaw["Scenarios"])
	d.Set("source", objectRaw["Source"])

	return nil
}

func resourceAliCloudThreatDetectionImageEventOperationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateImageEventOperation"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Id"] = d.Id()

	if d.HasChange("note") {
		update = true
	}
	if v, ok := d.GetOk("note"); ok {
		request["Note"] = v
	}

	if d.HasChange("conditions") {
		update = true
	}
	request["Conditions"] = d.Get("conditions")

	if d.HasChange("scenarios") {
		update = true
	}
	if v, ok := d.GetOk("scenarios"); ok {
		request["Scenarios"] = v
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
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
	request["Id"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
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
