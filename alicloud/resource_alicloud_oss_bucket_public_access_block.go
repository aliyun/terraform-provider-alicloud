// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
)

func resourceAliCloudOssBucketPublicAccessBlock() *schema.Resource {
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
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))

	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("block_public_access"); v != nil {
		objectDataLocalMap["BlockPublicAccess"] = v
		request["PublicAccessBlockConfiguration"] = objectDataLocalMap
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Oss", genXmlParam("PUT", "2019-05-17", "PutBucketPublicAccessBlock", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_public_access_block", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(*hostMap["bucket"]))

	return resourceAliCloudOssBucketPublicAccessBlockRead(d, meta)
}

func resourceAliCloudOssBucketPublicAccessBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketPublicAccessBlock(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_public_access_block DescribeOssBucketPublicAccessBlock Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	publicAccessBlockConfiguration1RawObj, _ := jsonpath.Get("$.PublicAccessBlockConfiguration", objectRaw)
	publicAccessBlockConfiguration1Raw := make(map[string]interface{})
	if publicAccessBlockConfiguration1RawObj != nil {
		publicAccessBlockConfiguration1Raw = publicAccessBlockConfiguration1RawObj.(map[string]interface{})
	}
	if len(publicAccessBlockConfiguration1Raw) > 0 {
		d.Set("block_public_access", formatBool(publicAccessBlockConfiguration1Raw["BlockPublicAccess"]))
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
	hostMap["bucket"] = StringPointer(d.Id())
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Oss", genXmlParam("PUT", "2019-05-17", "PutBucketPublicAccessBlock", action), query, body, nil, hostMap, false)
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
	hostMap["bucket"] = StringPointer(d.Id())

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Oss", genXmlParam("DELETE", "2019-05-17", "DeleteBucketPublicAccessBlock", action), query, body, nil, hostMap, false)
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
