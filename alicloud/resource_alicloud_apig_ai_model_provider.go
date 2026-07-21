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

func resourceAliCloudApigAiModelProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigAiModelProviderCreate,
		Read:   resourceAliCloudApigAiModelProviderRead,
		Update: resourceAliCloudApigAiModelProviderUpdate,
		Delete: resourceAliCloudApigAiModelProviderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"model_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"model_provider": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudApigAiModelProviderCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/v1/ai-model-providers")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["provider"] = d.Get("model_provider")
	request["gatewayId"] = d.Get("gateway_id")
	request["displayName"] = d.Get("display_name")
	if v, ok := d.GetOk("service_ids"); ok {
		serviceIdsMapsArray := convertToInterfaceArray(v)

		request["serviceIds"] = serviceIdsMapsArray
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_ai_model_provider", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.modelProviderId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudApigAiModelProviderRead(d, meta)
}

func resourceAliCloudApigAiModelProviderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := ApigServiceV2{client}

	objectRaw, err := apigServiceV2.DescribeApigAiModelProvider(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_ai_model_provider DescribeApigAiModelProvider Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("model_count", objectRaw["modelCount"])
	d.Set("source", objectRaw["source"])
	d.Set("update_time", objectRaw["updateTime"])
	d.Set("display_name", objectRaw["displayName"])
	d.Set("gateway_id", objectRaw["gatewayId"])
	d.Set("model_provider", objectRaw["provider"])

	serviceIds := make([]string, 0)
	if boundServicesRaw, ok := objectRaw["boundServices"].([]interface{}); ok {
		for _, item := range boundServicesRaw {
			if svc, ok := item.(map[string]interface{}); ok {
				if id, ok := svc["serviceId"].(string); ok && id != "" {
					serviceIds = append(serviceIds, id)
				}
			}
		}
	}
	d.Set("service_ids", serviceIds)

	return nil
}

func resourceAliCloudApigAiModelProviderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	modelProviderId := d.Id()
	action := fmt.Sprintf("/v1/ai-model-providers/%s", modelProviderId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	request["displayName"] = d.Get("display_name")
	if d.HasChange("display_name") {
		update = true
	}
	if d.HasChange("service_ids") {
		update = true
		if v, ok := d.GetOk("service_ids"); ok || d.HasChange("service_ids") {
			serviceIdsMapsArray := convertToInterfaceArray(v)

			request["serviceIds"] = serviceIdsMapsArray
		}
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

	return resourceAliCloudApigAiModelProviderRead(d, meta)
}

func resourceAliCloudApigAiModelProviderDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	modelProviderId := d.Id()
	action := fmt.Sprintf("/v1/ai-model-providers/%s", modelProviderId)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
