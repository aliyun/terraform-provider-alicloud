// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package oss

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAliCloudOssBucketWorm() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketWormCreate,
		Read:   resourceAliCloudOssBucketWormRead,
		Update: resourceAliCloudOssBucketWormUpdate,
		Delete: resourceAliCloudOssBucketWormDelete,
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
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"retention_period_in_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"worm_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudOssBucketWormCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?worm")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = helper.StringPointer(d.Get("bucket").(string))

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("retention_period_in_days"); !helper.IsNil(v) {
		objectDataLocalMap["RetentionPeriodInDays"] = v
		request["InitiateWormConfiguration"] = objectDataLocalMap
	}

	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("POST", "2019-05-17", "InitiateBucketWorm", action), query, body, nil, hostMap, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_oss_bucket_worm", action, errs.AlibabaCloudSdkGoERROR)
	}

	xOssWormIdVar, _ := response["headers"].(map[string]interface{})["x-oss-worm-id"]
	d.SetId(fmt.Sprintf("%v:%v", *hostMap["bucket"], xOssWormIdVar))

	return resourceAliCloudOssBucketWormUpdate(d, meta)
}

func resourceAliCloudOssBucketWormRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketWorm(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_worm DescribeOssBucketWorm Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["CreationDate"] != nil {
		d.Set("create_time", objectRaw["CreationDate"])
	}
	if objectRaw["RetentionPeriodInDays"] != nil {
		d.Set("retention_period_in_days", objectRaw["RetentionPeriodInDays"])
	}
	if objectRaw["State"] != nil {
		d.Set("status", objectRaw["State"])
	}
	if objectRaw["WormId"] != nil {
		d.Set("worm_id", objectRaw["WormId"])
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("bucket", parts[0])

	return nil
}

func resourceAliCloudOssBucketWormUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	if d.HasChange("status") {
		ossServiceV2 := OssServiceV2{client}
		object, err := ossServiceV2.DescribeOssBucketWorm(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}

		target := d.Get("status").(string)
		if object["State"].(string) != target {
			if target == "Locked" {
				parts := strings.Split(d.Id(), ":")
				action := fmt.Sprintf("/")
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				hostMap := make(map[string]*string)
				hostMap["bucket"] = helper.StringPointer(parts[0])
				query["wormId"] = helper.StringPointer(parts[1])

				body = request
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.Do("Oss", helper.XmlParam("POST", "2019-05-17", "CompleteBucketWorm", action), query, body, nil, hostMap, false)
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

			}
		}
	}

	parts := strings.Split(d.Id(), ":")
	action := fmt.Sprintf("/?wormExtend")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(parts[0])
	query["wormId"] = helper.StringPointer(parts[1])

	if !d.IsNewResource() && d.HasChange("retention_period_in_days") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("retention_period_in_days"); v != nil {
		objectDataLocalMap["RetentionPeriodInDays"] = d.Get("retention_period_in_days")
		request["ExtendWormConfiguration"] = objectDataLocalMap
	}

	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.Do("Oss", helper.XmlParam("POST", "2019-05-17", "ExtendBucketWorm", action), query, body, nil, hostMap, false)
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
		stateConf := helper.BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("retention_period_in_days"))}, d.Timeout(schema.TimeoutUpdate), 0, ossServiceV2.OssBucketWormStateRefreshFunc(d.Id(), "RetentionPeriodInDays", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceAliCloudOssBucketWormRead(d, meta)
}

func resourceAliCloudOssBucketWormDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("status"); ok {
		if v == "Locked" {
			log.Printf("[WARN] Cannot destroy resource alicloud_oss_bucket_worm which status valued Locked. Terraform will remove this resource from the state file, however resources may remain.")
			return nil
		}
	}
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := fmt.Sprintf("/?worm")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = helper.StringPointer(parts[0])

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("DELETE", "2019-05-17", "AbortBucketWorm", action), query, request, nil, hostMap, false)
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
		if errs.IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchWORMConfiguration"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
