package kms

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudKmsKeyVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKmsKeyVersionCreate,
		Read:   resourceAlicloudKmsKeyVersionRead,
		Delete: resourceAlicloudKmsKeyVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"key_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key_version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudKmsKeyVersionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateKeyVersion"
	request := make(map[string]interface{})
	var err error
	request["KeyId"] = d.Get("key_id")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_kms_key_version", action, errs.AlibabaCloudSdkGoERROR)
	}
	responseKeyVersion := response["KeyVersion"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseKeyVersion["KeyId"], ":", responseKeyVersion["KeyVersionId"]))

	return resourceAlicloudKmsKeyVersionRead(d, meta)
}
func resourceAlicloudKmsKeyVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsService := KmsService{client}
	_, err := kmsService.DescribeKmsKeyVersion(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kms_key_version kmsService.DescribeKmsKeyVersion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("key_id", parts[0])
	d.Set("key_version_id", parts[1])
	return nil
}
func resourceAlicloudKmsKeyVersionDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudKmsKeyVersion. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
