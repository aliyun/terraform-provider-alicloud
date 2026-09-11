// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudApigOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigOperationCreate,
		Read:   resourceAliCloudApigOperationRead,
		Update: resourceAliCloudApigOperationUpdate,
		Delete: resourceAliCloudApigOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"http_api_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"operation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operation_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"method": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"GET", "POST", "PUT", "DELETE", "HEAD", "PATCH", "OPTIONS",
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mock": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"response_code": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"response_content": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudApigOperationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	httpApiId := d.Get("http_api_id").(string)
	action := fmt.Sprintf("/v1/http-apis/%s/operations", httpApiId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	operation := make(map[string]interface{})
	operation["name"] = d.Get("operation_name")
	operation["path"] = d.Get("path")
	operation["method"] = d.Get("method")
	if v, ok := d.GetOk("description"); ok {
		operation["description"] = v
	}
	if v, ok := d.GetOk("mock"); ok {
		if mockList, ok := v.([]interface{}); ok && len(mockList) > 0 {
			mock := make(map[string]interface{})
			if mockData, ok := mockList[0].(map[string]interface{}); ok {
				mock["enable"] = mockData["enable"]
				mock["responseCode"] = mockData["response_code"]
				mock["responseContent"] = mockData["response_content"]
			}
			operation["mock"] = mock
		}
	}
	body["operations"] = []interface{}{operation}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_operation", action, AlibabaCloudSdkGoERROR)
	}

	operationIdVar, _ := jsonpath.Get("$.data.operations[0].operationId", response)
	if operationIdVar == nil || fmt.Sprint(operationIdVar) == "" {
		return WrapError(fmt.Errorf("failed to create apig operation: no operationId returned in response"))
	}
	d.SetId(fmt.Sprintf("%v:%v", httpApiId, operationIdVar))

	return resourceAliCloudApigOperationRead(d, meta)
}

func resourceAliCloudApigOperationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := ApigServiceV2{client}

	objectRaw, err := apigServiceV2.DescribeApigOperation(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_operation DescribeApigOperation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("operation_id", objectRaw["operationId"])
	d.Set("operation_name", objectRaw["name"])
	d.Set("path", objectRaw["path"])
	d.Set("method", objectRaw["method"])
	d.Set("description", objectRaw["description"])
	d.Set("create_time", objectRaw["createTimestamp"])

	mockMaps := make([]map[string]interface{}, 0)
	mockMap := make(map[string]interface{})
	mockRaw := make(map[string]interface{})
	if objectRaw["mock"] != nil {
		if mr, ok := objectRaw["mock"].(map[string]interface{}); ok {
			mockRaw = mr
		}
	}
	if len(mockRaw) > 0 {
		mockMap["enable"] = mockRaw["enable"]
		mockMap["response_code"] = mockRaw["responseCode"]
		mockMap["response_content"] = mockRaw["responseContent"]
		mockMaps = append(mockMaps, mockMap)
	}
	if err := d.Set("mock", mockMaps); err != nil {
		return WrapError(err)
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("http_api_id", parts[0])

	return nil
}

func resourceAliCloudApigOperationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	httpApiId := parts[0]
	operationId := parts[1]
	action := fmt.Sprintf("/v1/http-apis/%s/operations/%s", httpApiId, operationId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if !d.IsNewResource() && d.HasChange("operation_name") {
		update = true
	}
	if !d.IsNewResource() && d.HasChange("path") {
		update = true
	}
	if !d.IsNewResource() && d.HasChange("method") {
		update = true
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if !d.IsNewResource() && d.HasChange("mock") {
		update = true
	}

	operation := make(map[string]interface{})
	operation["name"] = d.Get("operation_name")
	operation["path"] = d.Get("path")
	operation["method"] = d.Get("method")
	if v, ok := d.GetOk("description"); ok {
		operation["description"] = v
	}
	if v, ok := d.GetOk("mock"); ok {
		if mockList, ok := v.([]interface{}); ok && len(mockList) > 0 {
			mock := make(map[string]interface{})
			if mockData, ok := mockList[0].(map[string]interface{}); ok {
				mock["enable"] = mockData["enable"]
				mock["responseCode"] = mockData["response_code"]
				mock["responseContent"] = mockData["response_content"]
			}
			operation["mock"] = mock
		}
	}
	body["operation"] = operation

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
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

	return resourceAliCloudApigOperationRead(d, meta)
}

func resourceAliCloudApigOperationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	httpApiId := parts[0]
	operationId := parts[1]
	action := fmt.Sprintf("/v1/http-apis/%s/operations/%s", httpApiId, operationId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)
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
		if NotFoundError(err) || IsExpectedErrors(err, []string{"DatabaseError.RecordNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
