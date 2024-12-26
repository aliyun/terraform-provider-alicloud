// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudOssBucketWorm() *schema.Resource {
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
	conn, err := client.NewOssClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("retention_period_in_days"); !IsNil(v) {
		objectDataLocalMap["RetentionPeriodInDays"] = v
		request["InitiateWormConfiguration"] = objectDataLocalMap
	}

	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.Execute(genXmlParam("InitiateBucketWorm", "POST", "2019-05-17", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_worm", action, AlibabaCloudSdkGoERROR)
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
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_worm DescribeOssBucketWorm Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["State"].(string) != target {
			if target == "Locked" {
				parts := strings.Split(d.Id(), ":")
				action := fmt.Sprintf("/")
				conn, err := client.NewOssClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				hostMap := make(map[string]*string)
				hostMap["bucket"] = StringPointer(parts[0])
				query["wormId"] = StringPointer(parts[1])

				body = request
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.Execute(genXmlParam("CompleteBucketWorm", "POST", "2019-05-17", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})
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
		}
	}

	parts := strings.Split(d.Id(), ":")
	action := fmt.Sprintf("/?wormExtend")
	conn, err := client.NewOssClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = StringPointer(parts[0])
	query["wormId"] = StringPointer(parts[1])

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
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.Execute(genXmlParam("ExtendBucketWorm", "POST", "2019-05-17", action), &openapi.OpenApiRequest{Query: query, Body: body, HostMap: hostMap}, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("retention_period_in_days"))}, d.Timeout(schema.TimeoutUpdate), 0, ossServiceV2.OssBucketWormStateRefreshFunc(d.Id(), "RetentionPeriodInDays", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
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
	conn, err := client.NewOssClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(parts[0])

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.Execute(genXmlParam("AbortBucketWorm", "DELETE", "2019-05-17", action), &openapi.OpenApiRequest{Query: query, Body: nil, HostMap: hostMap}, &util.RuntimeOptions{})

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
		if IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchWORMConfiguration"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
