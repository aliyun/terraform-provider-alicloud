// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudOssBucketLogging() *schema.Resource {
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
	conn, err := client.NewOssClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))

	objectDataLocalMap := make(map[string]interface{})

	loggingEnabled := make(map[string]interface{})
	if v, ok := d.GetOk("target_prefix"); ok {
		loggingEnabled["TargetPrefix"] = v
	}
	loggingEnabled["TargetBucket"] = d.Get("target_bucket")
	objectDataLocalMap["LoggingEnabled"] = loggingEnabled

	request["BucketLoggingStatus"] = objectDataLocalMap
	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.Execute(genXmlParam("PutBucketLogging", "PUT", "2019-05-17", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_logging", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(*hostMap["bucket"]))

	ossServiceV2 := OssServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("target_bucket"))}, d.Timeout(schema.TimeoutCreate), 0, ossServiceV2.OssBucketLoggingStateRefreshFunc(d.Id(), "TargetBucket", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudOssBucketLoggingRead(d, meta)
}

func resourceAliCloudOssBucketLoggingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketLogging(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_logging DescribeOssBucketLogging Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["TargetBucket"] != nil {
		d.Set("target_bucket", objectRaw["TargetBucket"])
	}
	if objectRaw["TargetPrefix"] != nil {
		d.Set("target_prefix", objectRaw["TargetPrefix"])
	}

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

	action := fmt.Sprintf("/?logging")
	conn, err := client.NewOssClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = StringPointer(d.Id())

	objectDataLocalMap := make(map[string]interface{})
	loggingEnabled := make(map[string]interface{})
	if d.HasChange("target_bucket") {
		update = true
	}
	loggingEnabled["TargetBucket"] = d.Get("target_bucket")

	if d.HasChange("target_prefix") {
		update = true
	}
	if v, ok := d.GetOk("target_prefix"); ok {
		loggingEnabled["TargetPrefix"] = v
	}
	objectDataLocalMap["LoggingEnabled"] = loggingEnabled

	request["BucketLoggingStatus"] = objectDataLocalMap
	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.Execute(genXmlParam("PutBucketLogging", "PUT", "2019-05-17", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})
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
		ossServiceV2 := OssServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("target_bucket"))}, d.Timeout(schema.TimeoutUpdate), 0, ossServiceV2.OssBucketLoggingStateRefreshFunc(d.Id(), "TargetBucket", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
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
	conn, err := client.NewOssClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Id())

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.Execute(genXmlParam("DeleteBucketLogging", "DELETE", "2019-05-17", action), &openapi.OpenApiRequest{Query: query, Body: nil, HostMap: hostMap}, &util.RuntimeOptions{})

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
