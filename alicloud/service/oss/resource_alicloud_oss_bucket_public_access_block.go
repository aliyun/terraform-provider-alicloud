// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package oss

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAliCloudOssBucketPublicAccessBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketPublicAccessBlockCreate,
		Read:   resourceAliCloudOssBucketPublicAccessBlockRead,
		Update: resourceAliCloudOssBucketPublicAccessBlockUpdate,
		Delete: resourceAliCloudOssBucketPublicAccessBlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"block_public_access": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudOssBucketPublicAccessBlockCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?publicAccessBlock")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = helper.StringPointer(d.Get("bucket").(string))

	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("block_public_access"); v != nil {
		objectDataLocalMap["BlockPublicAccess"] = v
		request["PublicAccessBlockConfiguration"] = objectDataLocalMap
	}

	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("PUT", "2019-05-17", "PutBucketPublicAccessBlock", action), query, body, nil, hostMap, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_oss_bucket_public_access_block", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(*hostMap["bucket"]))

	return resourceAliCloudOssBucketPublicAccessBlockRead(d, meta)
}

func resourceAliCloudOssBucketPublicAccessBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketPublicAccessBlock(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_public_access_block DescribeOssBucketPublicAccessBlock Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	publicAccessBlockConfiguration1RawObj, _ := jsonpath.Get("$.PublicAccessBlockConfiguration", objectRaw)
	publicAccessBlockConfiguration1Raw := make(map[string]interface{})
	if publicAccessBlockConfiguration1RawObj != nil {
		publicAccessBlockConfiguration1Raw = publicAccessBlockConfiguration1RawObj.(map[string]interface{})
	}
	if len(publicAccessBlockConfiguration1Raw) > 0 {
		block_public_accessVal, err := convert.FormatBool(publicAccessBlockConfiguration1Raw["BlockPublicAccess"])
		if err != nil {
			return err
		}
		d.Set("block_public_access", block_public_accessVal)
	}
	d.Set("bucket", d.Id())

	return nil
}

func resourceAliCloudOssBucketPublicAccessBlockUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	action := fmt.Sprintf("/?publicAccessBlock")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = helper.StringPointer(d.Id())
	if d.HasChange("block_public_access") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("block_public_access"); v != nil {
		objectDataLocalMap["BlockPublicAccess"] = d.Get("block_public_access")
		request["PublicAccessBlockConfiguration"] = objectDataLocalMap
	}

	body = request
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.Do("Oss", helper.XmlParam("PUT", "2019-05-17", "PutBucketPublicAccessBlock", action), query, body, nil, hostMap, false)
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
	}

	return resourceAliCloudOssBucketPublicAccessBlockRead(d, meta)
}

func resourceAliCloudOssBucketPublicAccessBlockDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/?publicAccessBlock")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = helper.StringPointer(d.Id())

	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.Do("Oss", helper.XmlParam("DELETE", "2019-05-17", "DeleteBucketPublicAccessBlock", action), query, body, nil, hostMap, false)
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
