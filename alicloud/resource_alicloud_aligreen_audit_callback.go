// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAligreenAuditCallback() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAligreenAuditCallbackCreate,
		Read:   resourceAliCloudAligreenAuditCallbackRead,
		Update: resourceAliCloudAligreenAuditCallbackUpdate,
		Delete: resourceAliCloudAligreenAuditCallbackDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"audit_callback_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"callback_suggestions": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"callback_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"crypt_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliCloudAligreenAuditCallbackCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAuditCallback"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["Name"] = d.Get("audit_callback_name")

	request["Url"] = d.Get("url")
	request["CryptType"] = d.Get("crypt_type")
	jsonPathResult2, err := jsonpath.Get("$", d.Get("callback_types"))
	if err == nil {
		request["CallbackTypes"] = convertListToJsonString(jsonPathResult2.([]interface{}))
	}

	jsonPathResult3, err := jsonpath.Get("$", d.Get("callback_suggestions"))
	if err == nil {
		request["CallbackSuggestions"] = convertListToJsonString(jsonPathResult3.([]interface{}))
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Green", "2017-08-23", action, query, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_aligreen_audit_callback", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(query["Name"]))

	return resourceAliCloudAligreenAuditCallbackUpdate(d, meta)
}

func resourceAliCloudAligreenAuditCallbackRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	aligreenServiceV2 := AligreenServiceV2{client}

	objectRaw, err := aligreenServiceV2.DescribeAligreenAuditCallback(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_aligreen_audit_callback DescribeAligreenAuditCallback Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["CryptType"] != nil {
		d.Set("crypt_type", objectRaw["CryptType"])
	}
	if objectRaw["Url"] != nil {
		d.Set("url", objectRaw["Url"])
	}
	if objectRaw["Name"] != nil {
		d.Set("audit_callback_name", objectRaw["Name"])
	}

	callbackSuggestions1Raw := make([]interface{}, 0)
	if objectRaw["CallbackSuggestions"] != nil {
		callbackSuggestions1Raw = objectRaw["CallbackSuggestions"].([]interface{})
	}

	d.Set("callback_suggestions", callbackSuggestions1Raw)
	callbackTypes1Raw := make([]interface{}, 0)
	if objectRaw["CallbackTypes"] != nil {
		callbackTypes1Raw = objectRaw["CallbackTypes"].([]interface{})
	}

	d.Set("callback_types", callbackTypes1Raw)

	d.Set("audit_callback_name", d.Id())

	return nil
}

func resourceAliCloudAligreenAuditCallbackUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "ModifyAuditCallback"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Name"] = d.Id()

	if !d.IsNewResource() && d.HasChange("url") {
		update = true
	}
	request["Url"] = d.Get("url")
	if !d.IsNewResource() && d.HasChange("crypt_type") {
		update = true
	}
	request["CryptType"] = d.Get("crypt_type")
	if !d.IsNewResource() && d.HasChange("callback_types") {
		update = true
	}
	jsonPathResult2, err := jsonpath.Get("$", d.Get("callback_types"))
	if err == nil {
		request["CallbackTypes"] = convertListToJsonString(jsonPathResult2.([]interface{}))
	}

	if !d.IsNewResource() && d.HasChange("callback_suggestions") {
		update = true
	}
	jsonPathResult3, err := jsonpath.Get("$", d.Get("callback_suggestions"))
	if err == nil {
		request["CallbackSuggestions"] = convertListToJsonString(jsonPathResult3.([]interface{}))
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Green", "2017-08-23", action, query, request, false)
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

	return resourceAliCloudAligreenAuditCallbackRead(d, meta)
}

func resourceAliCloudAligreenAuditCallbackDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAuditCallback"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["Name"] = d.Id()
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Green", "2017-08-23", action, query, request, false)

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
