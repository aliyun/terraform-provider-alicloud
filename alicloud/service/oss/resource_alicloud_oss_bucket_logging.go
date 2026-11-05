package oss

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

func ResourceAliCloudOssBucketLogging() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketLoggingCreate,
		Read:   resourceAliCloudOssBucketLoggingRead,
		Update: resourceAliCloudOssBucketLoggingUpdate,
		Delete: resourceAliCloudOssBucketLoggingDelete,
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
			"logging_role": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudOssBucketLoggingCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?logging")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = helper.StringPointer(d.Get("bucket").(string))

	dataList := make(map[string]interface{})

	if v := d.Get("target_bucket"); !helper.IsNil(v) {
		loggingEnabled := make(map[string]interface{})
		loggingEnabled["TargetPrefix"] = d.Get("target_prefix")
		loggingEnabled["TargetBucket"] = d.Get("target_bucket")
		loggingEnabled["LoggingRole"] = d.Get("logging_role")
		dataList["LoggingEnabled"] = loggingEnabled
	}

	request["BucketLoggingStatus"] = dataList

	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("PUT", "2019-05-17", "PutBucketLogging", action), query, body, nil, hostMap, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_oss_bucket_logging", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(*hostMap["bucket"]))

	ossServiceV2 := OssServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("target_bucket"))}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ossServiceV2.OssBucketLoggingStateRefreshFunc(d.Id(), "TargetBucket", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudOssBucketLoggingRead(d, meta)
}

func resourceAliCloudOssBucketLoggingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketLogging(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_logging DescribeOssBucketLogging Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("logging_role", objectRaw["LoggingRole"])
	d.Set("target_bucket", objectRaw["TargetBucket"])
	d.Set("target_prefix", objectRaw["TargetPrefix"])

	d.Set("bucket", d.Id())

	return nil
}

func resourceAliCloudOssBucketLoggingUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	action := fmt.Sprintf("/?logging")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(d.Id())

	dataList := make(map[string]interface{})

	if d.HasChanges("target_bucket", "target_prefix", "logging_role") {
		update = true
	}
	if v := d.Get("target_bucket"); v != nil {
		loggingEnabled := make(map[string]interface{})
		if v, ok := d.GetOk("target_bucket"); ok {
			loggingEnabled["TargetBucket"] = v
		}
		if v, ok := d.GetOk("target_prefix"); ok {
			loggingEnabled["TargetPrefix"] = v
		}
		if v, ok := d.GetOk("logging_role"); ok {
			loggingEnabled["LoggingRole"] = v
		}

		dataList["LoggingEnabled"] = loggingEnabled
	}

	request["BucketLoggingStatus"] = dataList

	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.Do("Oss", helper.XmlParam("PUT", "2019-05-17", "PutBucketLogging", action), query, body, nil, hostMap, false)
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
		stateConf := helper.BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("target_bucket"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ossServiceV2.OssBucketLoggingStateRefreshFunc(d.Id(), "TargetBucket", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceAliCloudOssBucketLoggingRead(d, meta)
}

func resourceAliCloudOssBucketLoggingDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/?logging")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = helper.StringPointer(d.Id())

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("DELETE", "2019-05-17", "DeleteBucketLogging", action), query, nil, nil, hostMap, false)
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
