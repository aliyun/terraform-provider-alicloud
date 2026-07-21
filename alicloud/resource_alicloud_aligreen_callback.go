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

func resourceAliCloudAligreenCallback() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAligreenCallbackCreate,
		Read:   resourceAliCloudAligreenCallbackRead,
		Update: resourceAliCloudAligreenCallbackUpdate,
		Delete: resourceAliCloudAligreenCallbackDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"callback_name": {
				Type:     schema.TypeString,
				Required: true,
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
			"callback_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"crypt_type": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudAligreenCallbackCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateCallback"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["Name"] = d.Get("callback_name")
	request["CallbackUrl"] = d.Get("callback_url")
	if v, ok := d.GetOk("crypt_type"); ok {
		request["CryptType"] = v
	}
	jsonPathResult3, err := jsonpath.Get("$", d.Get("callback_types"))
	if err == nil {
		request["CallbackTypes"] = convertListToJsonString(jsonPathResult3.([]interface{}))
	}

	jsonPathResult4, err := jsonpath.Get("$", d.Get("callback_suggestions"))
	if err == nil {
		request["CallbackSuggestions"] = convertListToJsonString(jsonPathResult4.([]interface{}))
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_aligreen_callback", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return resourceAliCloudAligreenCallbackRead(d, meta)
}

func resourceAliCloudAligreenCallbackRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	aligreenServiceV2 := AligreenServiceV2{client}

	objectRaw, err := aligreenServiceV2.DescribeAligreenCallback(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_aligreen_callback DescribeAligreenCallback Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Name"] != nil {
		d.Set("callback_name", objectRaw["Name"])
	}
	if objectRaw["CallbackUrl"] != nil {
		d.Set("callback_url", objectRaw["CallbackUrl"])
	}
	if objectRaw["CreateTime"] != nil {
		d.Set("create_time", objectRaw["CreateTime"])
	}
	if objectRaw["CryptType"] != nil {
		d.Set("crypt_type", objectRaw["CryptType"])
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

	return nil
}

func resourceAliCloudAligreenCallbackUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "UpdateCallback"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["Id"] = d.Id()

	if d.HasChange("callback_name") {
		update = true
	}
	request["Name"] = d.Get("callback_name")
	if d.HasChange("callback_url") {
		update = true
	}
	request["CallbackUrl"] = d.Get("callback_url")
	if d.HasChange("crypt_type") {
		update = true
		request["CryptType"] = d.Get("crypt_type")
	}

	if d.HasChange("callback_types") {
		update = true
	}
	jsonPathResult3, err := jsonpath.Get("$", d.Get("callback_types"))
	if err == nil {
		request["CallbackTypes"] = convertListToJsonString(jsonPathResult3.([]interface{}))
	}

	if d.HasChange("callback_suggestions") {
		update = true
	}
	jsonPathResult4, err := jsonpath.Get("$", d.Get("callback_suggestions"))
	if err == nil {
		request["CallbackSuggestions"] = convertListToJsonString(jsonPathResult4.([]interface{}))
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

	return resourceAliCloudAligreenCallbackRead(d, meta)
}

func resourceAliCloudAligreenCallbackDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteCallback"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["Id"] = d.Id()

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
