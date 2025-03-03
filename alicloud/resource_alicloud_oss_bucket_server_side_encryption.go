// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

func resourceAliCloudOssBucketServerSideEncryption() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketServerSideEncryptionCreate,
		Read:   resourceAliCloudOssBucketServerSideEncryptionRead,
		Update: resourceAliCloudOssBucketServerSideEncryptionUpdate,
		Delete: resourceAliCloudOssBucketServerSideEncryptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"kms_data_encryption": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"SM4"}, true),
			},
			"kms_master_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sse_algorithm": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"KMS", "AES256", "SM4"}, true),
			},
		},
	}
}

func resourceAliCloudOssBucketServerSideEncryptionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?encryption")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))

	objectDataLocalMap := make(map[string]interface{})
	applyServerSideEncryptionByDefault := make(map[string]interface{})

	if v, ok := d.GetOk("sse_algorithm"); ok {
		applyServerSideEncryptionByDefault["SSEAlgorithm"] = v
	}

	if v, ok := d.GetOk("kms_data_encryption"); ok {
		applyServerSideEncryptionByDefault["KMSDataEncryption"] = v
	}

	if v, ok := d.GetOk("kms_master_key_id"); ok {
		applyServerSideEncryptionByDefault["KMSMasterKeyID"] = v
	}

	objectDataLocalMap["ApplyServerSideEncryptionByDefault"] = applyServerSideEncryptionByDefault

	request["ServerSideEncryptionRule"] = objectDataLocalMap
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Oss", genXmlParam("PUT", "2019-05-17", "PutBucketEncryption", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_server_side_encryption", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(*hostMap["bucket"]))

	return resourceAliCloudOssBucketServerSideEncryptionRead(d, meta)
}

func resourceAliCloudOssBucketServerSideEncryptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketServerSideEncryption(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_server_side_encryption DescribeOssBucketServerSideEncryption Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("kms_data_encryption", objectRaw["KMSDataEncryption"])
	d.Set("kms_master_key_id", objectRaw["KMSMasterKeyID"])
	d.Set("sse_algorithm", objectRaw["SSEAlgorithm"])

	d.Set("bucket", d.Id())

	return nil
}

func resourceAliCloudOssBucketServerSideEncryptionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	action := fmt.Sprintf("/?encryption")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = StringPointer(d.Id())
	objectDataLocalMap := make(map[string]interface{})
	applyServerSideEncryptionByDefault := make(map[string]interface{})
	if d.HasChange("sse_algorithm") {
		update = true
	}
	if v, ok := d.GetOk("sse_algorithm"); ok {
		applyServerSideEncryptionByDefault["SSEAlgorithm"] = v
	}

	if d.HasChange("kms_data_encryption") {
		update = true
	}
	if v, ok := d.GetOk("kms_data_encryption"); ok {
		applyServerSideEncryptionByDefault["KMSDataEncryption"] = v
	}

	if d.HasChange("kms_master_key_id") {
		update = true
	}
	if v, ok := d.GetOk("kms_master_key_id"); ok {
		applyServerSideEncryptionByDefault["KMSMasterKeyID"] = v
	}
	objectDataLocalMap["ApplyServerSideEncryptionByDefault"] = applyServerSideEncryptionByDefault
	request["ServerSideEncryptionRule"] = objectDataLocalMap
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Oss", genXmlParam("PUT", "2019-05-17", "PutBucketEncryption", action), query, body, nil, hostMap, false)
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

	return resourceAliCloudOssBucketServerSideEncryptionRead(d, meta)
}

func resourceAliCloudOssBucketServerSideEncryptionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/?encryption")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Id())

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Oss", genXmlParam("DELETE", "2019-05-17", "DeleteBucketEncryption", action), query, body, nil, hostMap, false)
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
		if IsExpectedErrors(err, []string{"404"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
