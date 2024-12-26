// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudApigHttpApi() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigHttpApiCreate,
		Read:   resourceAliCloudApigHttpApiRead,
		Update: resourceAliCloudApigHttpApiUpdate,
		Delete: resourceAliCloudApigHttpApiDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"base_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"http_api_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocols": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudApigHttpApiCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/v1/http-apis")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	conn, err := client.NewApigClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query["RegionId"] = StringPointer(client.RegionId)

	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("protocols"); ok {
		protocolsMapsArray := v.([]interface{})
		request["protocols"] = protocolsMapsArray
	}

	if v, ok := d.GetOk("base_path"); ok {
		request["basePath"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		request["type"] = v
	}
	request["name"] = d.Get("http_api_name")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2024-03-27"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_http_api", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.body.data.httpApiId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudApigHttpApiRead(d, meta)
}

func resourceAliCloudApigHttpApiRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := ApigServiceV2{client}

	objectRaw, err := apigServiceV2.DescribeApigHttpApi(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_http_api DescribeApigHttpApi Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["basePath"] != nil {
		d.Set("base_path", objectRaw["basePath"])
	}
	if objectRaw["description"] != nil {
		d.Set("description", objectRaw["description"])
	}
	if objectRaw["name"] != nil {
		d.Set("http_api_name", objectRaw["name"])
	}
	if objectRaw["resourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["resourceGroupId"])
	}
	if objectRaw["type"] != nil {
		d.Set("type", objectRaw["type"])
	}

	protocols1Raw := make([]interface{}, 0)
	if objectRaw["protocols"] != nil {
		protocols1Raw = objectRaw["protocols"].([]interface{})
	}

	d.Set("protocols", protocols1Raw)

	return nil
}

func resourceAliCloudApigHttpApiUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	httpApiId := d.Id()
	action := fmt.Sprintf("/v1/http-apis/%s", httpApiId)
	conn, err := client.NewApigClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["httpApiId"] = d.Id()

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	if d.HasChange("protocols") {
		update = true
	}
	if v, ok := d.GetOk("protocols"); ok || d.HasChange("protocols") {
		protocolsMapsArray := v.([]interface{})
		request["protocols"] = protocolsMapsArray
	}

	if d.HasChange("base_path") {
		update = true
	}
	request["basePath"] = d.Get("base_path")
	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2024-03-27"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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
	action = fmt.Sprintf("/move-resource-group")
	conn, err = client.NewApigClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["ResourceId"] = StringPointer(d.Id())
	query["RegionId"] = StringPointer(client.RegionId)
	if d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		query["ResourceGroupId"] = StringPointer(v.(string))
	}

	query["ResourceType"] = StringPointer("HttpApi")
	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2024-03-27"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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
	return resourceAliCloudApigHttpApiRead(d, meta)
}

func resourceAliCloudApigHttpApiDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	httpApiId := d.Id()
	action := fmt.Sprintf("/v1/http-apis/%s", httpApiId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	conn, err := client.NewApigClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["httpApiId"] = d.Id()

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2024-03-27"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, nil, &runtime)

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
		if IsExpectedErrors(err, []string{"Error.DatabaseError.RecordNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
