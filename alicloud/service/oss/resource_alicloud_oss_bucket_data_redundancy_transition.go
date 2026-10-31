// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package oss

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAliCloudOssBucketDataRedundancyTransition() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketDataRedundancyTransitionCreate,
		Read:   resourceAliCloudOssBucketDataRedundancyTransitionRead,
		Delete: resourceAliCloudOssBucketDataRedundancyTransitionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudOssBucketDataRedundancyTransitionCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?redundancyTransition")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = helper.StringPointer(d.Get("bucket").(string))

	query["x-oss-target-redundancy-type"] = helper.StringPointer("ZRS")
	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("POST", "2019-05-17", "CreateBucketDataRedundancyTransition", action), query, body, nil, hostMap, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_oss_bucket_data_redundancy_transition", action, errs.AlibabaCloudSdkGoERROR)
	}

	BucketDataRedundancyTransitionTaskId, _ := jsonpath.Get("$.BucketDataRedundancyTransition.TaskId", response)
	d.SetId(fmt.Sprintf("%v:%v", *hostMap["bucket"], BucketDataRedundancyTransitionTaskId))

	return resourceAliCloudOssBucketDataRedundancyTransitionRead(d, meta)
}

func resourceAliCloudOssBucketDataRedundancyTransitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketDataRedundancyTransition(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_data_redundancy_transition DescribeOssBucketDataRedundancyTransition Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	bucketDataRedundancyTransition1RawObj, _ := jsonpath.Get("$.BucketDataRedundancyTransition", objectRaw)
	bucketDataRedundancyTransition1Raw := make(map[string]interface{})
	if bucketDataRedundancyTransition1RawObj != nil {
		bucketDataRedundancyTransition1Raw = bucketDataRedundancyTransition1RawObj.(map[string]interface{})
	}
	if len(bucketDataRedundancyTransition1Raw) > 0 {
		d.Set("create_time", bucketDataRedundancyTransition1Raw["CreateTime"])
		d.Set("status", bucketDataRedundancyTransition1Raw["Status"])
	}
	parts := strings.Split(d.Id(), ":")
	d.Set("bucket", parts[0])
	d.Set("task_id", parts[1])

	return nil
}

func resourceAliCloudOssBucketDataRedundancyTransitionDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := fmt.Sprintf("/?redundancyTransition")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = helper.StringPointer(parts[0])
	query["x-oss-redundancy-transition-taskid"] = helper.StringPointer(parts[1])

	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("DELETE", "2019-05-17", "DeleteBucketDataRedundancyTransition", action), query, body, nil, hostMap, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
