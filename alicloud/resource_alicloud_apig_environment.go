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

func resourceAliCloudApigEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigEnvironmentCreate,
		Read:   resourceAliCloudApigEnvironmentRead,
		Update: resourceAliCloudApigEnvironmentUpdate,
		Delete: resourceAliCloudApigEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudApigEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/v1/environments")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["name"] = d.Get("environment_name")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	request["gatewayId"] = d.Get("gateway_id")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["resourceGroupId"] = v
	}
	body = request
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_environment", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.environmentId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudApigEnvironmentRead(d, meta)
}

func resourceAliCloudApigEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := ApigServiceV2{client}

	objectRaw, err := apigServiceV2.DescribeApigEnvironment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_environment DescribeApigEnvironment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["description"] != nil {
		d.Set("description", objectRaw["description"])
	}
	if objectRaw["name"] != nil {
		d.Set("environment_name", objectRaw["name"])
	}
	if objectRaw["resourceGroupId"] != nil {
		d.Set("resource_group_id", objectRaw["resourceGroupId"])
	}

	gatewayInfo1RawObj, _ := jsonpath.Get("$.gatewayInfo", objectRaw)
	gatewayInfo1Raw := make(map[string]interface{})
	if gatewayInfo1RawObj != nil {
		gatewayInfo1Raw = gatewayInfo1RawObj.(map[string]interface{})
	}
	if gatewayInfo1Raw["gatewayId"] != nil {
		d.Set("gateway_id", gatewayInfo1Raw["gatewayId"])
	}

	return nil
}

func resourceAliCloudApigEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	environmentId := d.Id()
	action := fmt.Sprintf("/v1/environments/%s", environmentId)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["environmentId"] = d.Id()

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("APIG", "2024-03-27", action, query, nil, body, true)
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

	query["ResourceType"] = StringPointer("Environment")
	body = request
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

	d.Partial(false)
	return resourceAliCloudApigEnvironmentRead(d, meta)
}

func resourceAliCloudApigEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	environmentId := d.Id()
	action := fmt.Sprintf("/v1/environments/%s", environmentId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	request["environmentId"] = d.Id()

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
		if IsExpectedErrors(err, []string{"NotFound.EnvironmentNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
