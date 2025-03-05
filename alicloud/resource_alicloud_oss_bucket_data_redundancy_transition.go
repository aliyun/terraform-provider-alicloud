// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudOssBucketDataRedundancyTransition() *schema.Resource {
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
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))

	query["x-oss-target-redundancy-type"] = StringPointer("ZRS")
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("POST", "2019-05-17", "CreateBucketDataRedundancyTransition", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_data_redundancy_transition", action, AlibabaCloudSdkGoERROR)
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
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_data_redundancy_transition DescribeOssBucketDataRedundancyTransition Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
	hostMap["bucket"] = StringPointer(parts[0])
	query["x-oss-redundancy-transition-taskid"] = StringPointer(parts[1])

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("DELETE", "2019-05-17", "DeleteBucketDataRedundancyTransition", action), query, body, nil, hostMap, false)
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
