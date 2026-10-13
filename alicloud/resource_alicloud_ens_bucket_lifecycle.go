// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudEnsBucketLifecycle() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEnsBucketLifecycleCreate,
		Read:   resourceAliCloudEnsBucketLifecycleRead,
		Update: resourceAliCloudEnsBucketLifecycleUpdate,
		Delete: resourceAliCloudEnsBucketLifecycleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expiration_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"allow_same_action_overlap": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Enabled", "Disabled"}, false),
			},
		},
	}
}

func resourceAliCloudEnsBucketLifecycleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "PutBucketLifecycle"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["BucketName"] = d.Get("bucket_name")
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("prefix"); ok {
		request["Prefix"] = v
	}
	if v, ok := d.GetOk("expiration_days"); ok {
		request["ExpirationDays"] = v
	}
	if v, ok := d.GetOkExists("allow_same_action_overlap"); ok && v.(bool) {
		request["AllowSameActionOverlap"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ens_bucket_lifecycle", action, AlibabaCloudSdkGoERROR)
	}

	ruleId := ""
	if v, ok := response["RuleId"]; ok && v != nil {
		ruleId = fmt.Sprint(v)
	}
	if ruleId == "" {
		return WrapError(Error("failed to get RuleId from PutBucketLifecycle response"))
	}
	d.SetId(fmt.Sprintf("%s:%s", d.Get("bucket_name"), ruleId))

	return resourceAliCloudEnsBucketLifecycleRead(d, meta)
}

func resourceAliCloudEnsBucketLifecycleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ensServiceV2 := EnsServiceV2{client}

	objectRaw, err := ensServiceV2.DescribeEnsBucketLifecycle(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ens_bucket_lifecycle DescribeEnsBucketLifecycle Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("bucket_name", parts[0])
	d.Set("rule_id", objectRaw["ID"])
	d.Set("prefix", objectRaw["Prefix"])
	d.Set("status", objectRaw["Status"])
	if expiration, ok := objectRaw["Expiration"].(map[string]interface{}); ok {
		if days, ok := expiration["Days"]; ok && days != nil {
			if daysStr := fmt.Sprint(days); daysStr != "" && daysStr != "<nil>" {
				if n, e := strconv.Atoi(daysStr); e == nil {
					d.Set("expiration_days", n)
				}
			}
		}
	}

	return nil
}

func resourceAliCloudEnsBucketLifecycleUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "PutBucketLifecycle"
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["BucketName"] = parts[0]
	request["RuleId"] = parts[1]
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("prefix"); ok {
		request["Prefix"] = v
	}
	if v, ok := d.GetOk("expiration_days"); ok {
		request["ExpirationDays"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return resourceAliCloudEnsBucketLifecycleRead(d, meta)
}

func resourceAliCloudEnsBucketLifecycleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteBucketLifecycle"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["BucketName"] = parts[0]
	query["RuleId"] = parts[1]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, []string{"NoSuchBucket", "NoSuchLifecycle"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
