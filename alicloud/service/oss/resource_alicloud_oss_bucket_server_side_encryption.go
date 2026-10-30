package oss

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAliCloudOssBucketServerSideEncryption() *schema.Resource {
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
				ValidateFunc: validators.StringInSlice([]string{"SM4"}, false),
			},
			"kms_master_key_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sse_algorithm": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validators.StringInSlice([]string{"KMS", "AES256", "SM4"}, false),
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
	hostMap["bucket"] = helper.StringPointer(d.Get("bucket").(string))

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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("PUT", "2019-05-17", "PutBucketEncryption", action), query, body, nil, hostMap, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_oss_bucket_server_side_encryption", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(*hostMap["bucket"]))

	ossServiceV2 := OssServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ossServiceV2.OssBucketServerSideEncryptionStateRefreshFunc(d.Id(), "#SSEAlgorithm", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudOssBucketServerSideEncryptionRead(d, meta)
}

func resourceAliCloudOssBucketServerSideEncryptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketServerSideEncryption(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_server_side_encryption DescribeOssBucketServerSideEncryption Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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

	var err error
	action := fmt.Sprintf("/?encryption")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(d.Id())

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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.Do("Oss", helper.XmlParam("PUT", "2019-05-17", "PutBucketEncryption", action), query, body, nil, hostMap, false)
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
		ossServiceV2 := OssServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"#CHECKSET"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ossServiceV2.OssBucketServerSideEncryptionStateRefreshFunc(d.Id(), "#SSEAlgorithm", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
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
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = helper.StringPointer(d.Id())

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("DELETE", "2019-05-17", "DeleteBucketEncryption", action), query, nil, nil, hostMap, false)

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
		if errs.IsExpectedErrors(err, []string{"404"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
