// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudOssAccessPointPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssAccessPointPolicyCreate,
		Read:   resourceAliCloudOssAccessPointPolicyRead,
		Update: resourceAliCloudOssAccessPointPolicyUpdate,
		Delete: resourceAliCloudOssAccessPointPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_point_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
		},
	}
}

func resourceAliCloudOssAccessPointPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/?accessPointPolicy")
	var request string
	var response map[string]interface{}
	query := make(map[string]*string)
	body := ""
	hostMap := make(map[string]*string)
	var err error
	request = ""
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))
	query["x-oss-access-point-name"] = StringPointer(d.Get("access_point_name").(string))

	request = d.Get("policy").(string)
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Oss", jsonXmlParam("PUT", "2019-05-17", "PutAccessPointPolicy", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_access_point_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", d.Get("bucket").(string), d.Get("access_point_name").(string)))

	return resourceAliCloudOssAccessPointPolicyRead(d, meta)
}

func resourceAliCloudOssAccessPointPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssAccessPointPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_access_point_policy DescribeOssAccessPointPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("policy", convertMapToJsonStringIgnoreError(objectRaw))
	if len(parts) > 0 {
		d.Set("bucket", parts[0])
	}
	if len(parts) > 1 {
		d.Set("access_point_name", parts[1])
	}

	return nil
}

func resourceAliCloudOssAccessPointPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request string
	var response map[string]interface{}
	var query map[string]*string
	var body string
	update := false
	action := fmt.Sprintf("/?accessPointPolicy")
	var err error
	request = ""
	query = make(map[string]*string)
	body = ""
	hostMap := make(map[string]*string)
	parts := strings.Split(d.Id(), ":")
	hostMap["bucket"] = StringPointer(parts[0])
	if len(parts) > 1 {
		query["x-oss-access-point-name"] = StringPointer(parts[1])
	}
	if d.HasChange("policy") {
		update = true
	}
	request = d.Get("policy").(string)
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Oss", jsonXmlParam("PUT", "2019-05-17", "PutAccessPointPolicy", action), query, body, nil, hostMap, false)
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

	return resourceAliCloudOssAccessPointPolicyRead(d, meta)
}

func resourceAliCloudOssAccessPointPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/?accessPointPolicy")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	parts := strings.Split(d.Id(), ":")
	hostMap["bucket"] = StringPointer(parts[0])
	if len(parts) > 1 {
		query["x-oss-access-point-name"] = StringPointer(parts[1])
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Oss", jsonXmlParam("DELETE", "2019-05-17", "DeleteAccessPointPolicy", action), query, body, nil, hostMap, false)
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
		if NotFoundError(err) {
			return nil
		}
		if IsExpectedErrors(err, []string{"NoSuchAccessPointPolicy", "NoSuchAccessPoint"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
