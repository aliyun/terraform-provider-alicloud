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

func resourceAliCloudOssBucketUserDefinedLogFields() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudOssBucketUserDefinedLogFieldsCreate,
		Read:   resourceAliCloudOssBucketUserDefinedLogFieldsRead,
		Update: resourceAliCloudOssBucketUserDefinedLogFieldsUpdate,
		Delete: resourceAliCloudOssBucketUserDefinedLogFieldsDelete,
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
			"header_set": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"param_set": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudOssBucketUserDefinedLogFieldsCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/?userDefinedLogFieldsConfig")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["bucket"] = StringPointer(d.Get("bucket").(string))

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("header_set"); !IsNil(v) {
		headerSet := make(map[string]interface{})
		nodeNative, _ := jsonpath.Get("$", v)
		if nodeNative != nil && nodeNative != "" {
			headerSet["header"] = nodeNative.(*schema.Set).List()
		}

		objectDataLocalMap["HeaderSet"] = headerSet
	}

	if v := d.Get("param_set"); !IsNil(v) {
		paramSet := make(map[string]interface{})
		nodeNative1, _ := jsonpath.Get("$", v)
		if nodeNative1 != nil && nodeNative1 != "" {
			paramSet["parameter"] = nodeNative1.(*schema.Set).List()
		}

		objectDataLocalMap["ParamSet"] = paramSet
	}

	request["UserDefinedLogFieldsConfiguration"] = objectDataLocalMap
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutUserDefinedLogFieldsConfig", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oss_bucket_user_defined_log_fields", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(*hostMap["bucket"]))

	return resourceAliCloudOssBucketUserDefinedLogFieldsRead(d, meta)
}

func resourceAliCloudOssBucketUserDefinedLogFieldsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ossServiceV2 := OssServiceV2{client}

	objectRaw, err := ossServiceV2.DescribeOssBucketUserDefinedLogFields(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oss_bucket_user_defined_log_fields DescribeOssBucketUserDefinedLogFields Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	header1RawObj, _ := jsonpath.Get("$.UserDefinedLogFieldsConfiguration.HeaderSet.header[*]", objectRaw)
	header1Raw := make([]interface{}, 0)
	if header1RawObj != nil {
		header1Raw = header1RawObj.([]interface{})
	}
	if len(header1Raw) == 0 {
		header1RawObj, _ := jsonpath.Get("$.UserDefinedLogFieldsConfiguration.HeaderSet.header", objectRaw)
		if header1RawObj != nil && fmt.Sprint(header1RawObj) != "" {
			header1Raw = append(header1Raw, header1RawObj)
		}
	}
	if len(header1Raw) > 0 {
		d.Set("header_set", header1Raw)
	}
	parameter1RawObj, _ := jsonpath.Get("$.UserDefinedLogFieldsConfiguration.ParamSet.parameter[*]", objectRaw)
	parameter1Raw := make([]interface{}, 0)
	if parameter1RawObj != nil {
		parameter1Raw = parameter1RawObj.([]interface{})
	}
	if len(parameter1Raw) == 0 {
		parameter1RawObj, _ := jsonpath.Get("$.UserDefinedLogFieldsConfiguration.ParamSet.parameter", objectRaw)
		if parameter1RawObj != nil && fmt.Sprint(parameter1RawObj) != "" {
			parameter1Raw = append(parameter1Raw, parameter1RawObj)
		}
	}

	if len(parameter1Raw) > 0 {
		d.Set("param_set", parameter1Raw)
	}

	d.Set("bucket", d.Id())

	return nil
}

func resourceAliCloudOssBucketUserDefinedLogFieldsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	action := fmt.Sprintf("/?userDefinedLogFieldsConfig")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["bucket"] = StringPointer(d.Id())
	objectDataLocalMap := make(map[string]interface{})

	if d.HasChange("header_set") {
		update = true
	}
	if v := d.Get("header_set"); v != nil {
		headerSet := make(map[string]interface{})
		nodeNative, _ := jsonpath.Get("$", d.Get("header_set"))
		if nodeNative != nil && nodeNative != "" {
			headerSet["header"] = nodeNative.(*schema.Set).List()
		}

		objectDataLocalMap["HeaderSet"] = headerSet
	}

	if d.HasChange("param_set") {
		update = true
	}
	if v := d.Get("param_set"); v != nil {
		paramSet := make(map[string]interface{})
		nodeNative1, _ := jsonpath.Get("$", d.Get("param_set"))
		if nodeNative1 != nil && nodeNative1 != "" {
			paramSet["parameter"] = nodeNative1.(*schema.Set).List()
		}

		objectDataLocalMap["ParamSet"] = paramSet
	}

	request["UserDefinedLogFieldsConfiguration"] = objectDataLocalMap
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Oss", xmlParam("PUT", "2019-05-17", "PutUserDefinedLogFieldsConfig", action), query, body, nil, hostMap, false)
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

	return resourceAliCloudOssBucketUserDefinedLogFieldsRead(d, meta)
}

func resourceAliCloudOssBucketUserDefinedLogFieldsDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/?userDefinedLogFieldsConfig")
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
		response, err = client.Do("Oss", xmlParam("DELETE", "2019-05-17", "DeleteUserDefinedLogFieldsConfig", action), query, body, nil, hostMap, false)
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
