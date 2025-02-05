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

func resourceAliCloudThreatDetectionMaliciousFileWhitelistConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionMaliciousFileWhitelistConfigCreate,
		Read:   resourceAliCloudThreatDetectionMaliciousFileWhitelistConfigRead,
		Update: resourceAliCloudThreatDetectionMaliciousFileWhitelistConfigUpdate,
		Delete: resourceAliCloudThreatDetectionMaliciousFileWhitelistConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"event_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"field": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"field_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"target_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudThreatDetectionMaliciousFileWhitelistConfigCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateMaliciousFileWhitelistConfig"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("event_name"); ok {
		request["EventName"] = v
	}
	if v, ok := d.GetOk("target_type"); ok {
		request["TargetType"] = v
	}
	if v, ok := d.GetOk("operator"); ok {
		request["Operator"] = v
	}
	if v, ok := d.GetOk("field"); ok {
		request["Field"] = v
	}
	if v, ok := d.GetOk("target_value"); ok {
		request["TargetValue"] = v
	}
	if v, ok := d.GetOk("field_value"); ok {
		request["FieldValue"] = v
	}
	if v, ok := d.GetOk("source"); ok {
		request["Source"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_malicious_file_whitelist_config", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.Id", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudThreatDetectionMaliciousFileWhitelistConfigRead(d, meta)
}

func resourceAliCloudThreatDetectionMaliciousFileWhitelistConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionMaliciousFileWhitelistConfig(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_malicious_file_whitelist_config DescribeThreatDetectionMaliciousFileWhitelistConfig Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("event_name", objectRaw["EventName"])
	d.Set("field", objectRaw["Field"])
	d.Set("field_value", objectRaw["FieldValue"])
	d.Set("operator", objectRaw["Operator"])
	d.Set("source", objectRaw["Source"])
	d.Set("target_type", objectRaw["TargetType"])
	d.Set("target_value", objectRaw["TargetValue"])

	return nil
}

func resourceAliCloudThreatDetectionMaliciousFileWhitelistConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "UpdateMaliciousFileWhitelistConfig"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["ConfigId"] = d.Id()
	if d.HasChange("event_name") {
		update = true
	}
	request["EventName"] = d.Get("event_name")

	if d.HasChange("target_type") {
		update = true
	}
	request["TargetType"] = d.Get("target_type")

	if d.HasChange("operator") {
		update = true
	}
	request["Operator"] = d.Get("operator")

	if d.HasChange("field") {
		update = true
	}
	request["Field"] = d.Get("field")

	if d.HasChange("target_value") {
		update = true
	}
	request["TargetValue"] = d.Get("target_value")

	if d.HasChange("field_value") {
		update = true
	}
	request["FieldValue"] = d.Get("field_value")

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

	return resourceAliCloudThreatDetectionMaliciousFileWhitelistConfigRead(d, meta)
}

func resourceAliCloudThreatDetectionMaliciousFileWhitelistConfigDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteMaliciousFileWhitelistConfig"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["ConfigId"] = d.Id()

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
