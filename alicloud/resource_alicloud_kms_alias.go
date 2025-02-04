package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudKmsAlias() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKmsAliasCreate,
		Read:   resourceAlicloudKmsAliasRead,
		Update: resourceAlicloudKmsAliasUpdate,
		Delete: resourceAlicloudKmsAliasDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"alias_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudKmsAliasCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAlias"
	request := make(map[string]interface{})
	var err error
	request["AliasName"] = d.Get("alias_name")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_alias", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["AliasName"]))

	return resourceAlicloudKmsAliasRead(d, meta)
}
func resourceAlicloudKmsAliasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsService := KmsService{client}
	object, err := kmsService.DescribeKmsAlias(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kms_alias kmsService.DescribeKmsAlias Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("alias_name", d.Id())
	d.Set("key_id", object["KeyId"])
	return nil
}
func resourceAlicloudKmsAliasUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	if d.HasChange("key_id") {
		request := map[string]interface{}{
			"AliasName": d.Id(),
		}
		request["KeyId"] = d.Get("key_id")
		action := "UpdateAlias"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudKmsAliasRead(d, meta)
}
func resourceAlicloudKmsAliasDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAlias"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"AliasName": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
