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

func resourceAliCloudThreatDetectionFileUploadLimit() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudThreatDetectionFileUploadLimitCreate,
		Read:   resourceAliCloudThreatDetectionFileUploadLimitRead,
		Update: resourceAliCloudThreatDetectionFileUploadLimitUpdate,
		Delete: resourceAliCloudThreatDetectionFileUploadLimitDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"limit": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntBetween(100, 10000),
			},
		},
	}
}

func resourceAliCloudThreatDetectionFileUploadLimitCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateFileUploadLimit"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["Limit"] = d.Get("limit")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_threat_detection_file_upload_limit", action, AlibabaCloudSdkGoERROR)
	}

	accountId, err := meta.(*connectivity.AliyunClient).AccountId()
	d.SetId(fmt.Sprintf(accountId))

	return resourceAliCloudThreatDetectionFileUploadLimitRead(d, meta)
}

func resourceAliCloudThreatDetectionFileUploadLimitRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}

	objectRaw, err := threatDetectionServiceV2.DescribeThreatDetectionFileUploadLimit(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_threat_detection_file_upload_limit DescribeThreatDetectionFileUploadLimit Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("limit", objectRaw["Limit"])

	return nil
}

func resourceAliCloudThreatDetectionFileUploadLimitUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "UpdateFileUploadLimit"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	if d.HasChange("limit") {
		update = true
	}
	request["Limit"] = d.Get("limit")
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

	return resourceAliCloudThreatDetectionFileUploadLimitRead(d, meta)
}

func resourceAliCloudThreatDetectionFileUploadLimitDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource File Upload Limit. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
