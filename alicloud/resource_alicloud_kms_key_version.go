package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_key_version", action, AlibabaCloudSdkGoERROR)
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
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kms_key_version kmsService.DescribeKmsKeyVersion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("key_id", parts[0])
	d.Set("key_version_id", parts[1])
	return nil
}
func resourceAlicloudKmsKeyVersionDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudKmsKeyVersion. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
