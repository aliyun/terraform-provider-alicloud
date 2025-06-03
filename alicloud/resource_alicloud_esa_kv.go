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

func resourceAliCloudEsaKv() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaKvCreate,
		Read:   resourceAliCloudEsaKvRead,
		Update: resourceAliCloudEsaKvUpdate,
		Delete: resourceAliCloudEsaKvDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"expiration": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"expiration_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"isbase": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEsaKvCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	invalidCreate := false
	if _, ok := d.GetOk("url"); ok {
		invalidCreate = true
	}
	if !invalidCreate {

		action := "PutKv"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		if v, ok := d.GetOk("namespace"); ok {
			request["Namespace"] = v
		}
		if v, ok := d.GetOk("key"); ok {
			request["Key"] = v
		}
		request["RegionId"] = client.RegionId

		if v, ok := d.GetOkExists("expiration_ttl"); ok {
			request["ExpirationTtl"] = v
		}
		if v, ok := d.GetOkExists("expiration"); ok {
			request["Expiration"] = v
		}
		if v, ok := d.GetOkExists("isbase"); ok {
			request["Base64"] = v
		}
		request["Value"] = d.Get("value")
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_kv", action, AlibabaCloudSdkGoERROR)
		}

		d.SetId(fmt.Sprintf("%v:%v", request["Namespace"], request["Key"]))

	}

	invalidCreate = false
	if _, ok := d.GetOk("value"); ok {
		invalidCreate = true
	}
	if !invalidCreate {

		action := "PutKvWithHighCapacity"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		if v, ok := d.GetOk("namespace"); ok {
			request["Namespace"] = v
		}
		if v, ok := d.GetOk("key"); ok {
			request["Key"] = v
		}
		request["RegionId"] = client.RegionId

		request["Url"] = d.Get("url")
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_kv", action, AlibabaCloudSdkGoERROR)
		}

		d.SetId(fmt.Sprintf("%v:%v", request["Namespace"], request["Key"]))

	}

	return resourceAliCloudEsaKvRead(d, meta)
}

func resourceAliCloudEsaKvRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaKv(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_kv DescribeEsaKv Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("value", objectRaw["Value"])

	parts := strings.Split(d.Id(), ":")
	d.Set("namespace", parts[0])
	d.Set("key", parts[1])

	return nil
}

func resourceAliCloudEsaKvUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	enablePutKv := false
	checkValue00 := d.Get("value")
	if !(checkValue00 == "") {
		enablePutKv = true
	}
	parts := strings.Split(d.Id(), ":")
	action := "PutKv"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Namespace"] = parts[0]
	request["Key"] = parts[1]
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("expiration_ttl"); ok {
		request["ExpirationTtl"] = v
	}
	if v, ok := d.GetOk("expiration"); ok {
		request["Expiration"] = v
	}
	if v, ok := d.GetOk("isbase"); ok {
		request["Base64"] = v
	}
	if d.HasChange("value") {
		update = true
	}
	request["Value"] = d.Get("value")
	if update && enablePutKv {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
	update = false
	enablePutKvWithHighCapacity := false
	checkValue00 = d.Get("url")
	if !(checkValue00 == "") {
		enablePutKvWithHighCapacity = true
	}
	parts = strings.Split(d.Id(), ":")
	action = "PutKvWithHighCapacity"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Namespace"] = parts[0]
	request["Key"] = parts[1]
	request["RegionId"] = client.RegionId
	request["Url"] = d.Get("url")
	if enablePutKvWithHighCapacity {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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

	d.Partial(false)
	return resourceAliCloudEsaKvRead(d, meta)
}

func resourceAliCloudEsaKvDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteKv"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["Namespace"] = parts[0]
	query["Key"] = parts[1]
	query["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcGet("ESA", "2024-09-10", action, query, request)

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
