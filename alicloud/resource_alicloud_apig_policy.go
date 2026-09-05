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

func resourceAliCloudApigPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigPolicyCreate,
		Read:   resourceAliCloudApigPolicyRead,
		Update: resourceAliCloudApigPolicyUpdate,
		Delete: resourceAliCloudApigPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"attach_resource_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"attach_resource_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"policy_class_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_class_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_config": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudApigPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/v1/policies")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("environment_id"); ok {
		request["environmentId"] = v
	}
	request["attachResourceType"] = d.Get("attach_resource_type")
	if v, ok := d.GetOk("gateway_id"); ok {
		request["gatewayId"] = v
	}
	if v, ok := d.GetOk("policy_name"); ok {
		request["name"] = v
	}
	if v, ok := d.GetOk("attach_resource_ids"); ok {
		attachResourceIdsMapsArray := convertToInterfaceArray(v)

		request["attachResourceIds"] = attachResourceIdsMapsArray
	}

	request["className"] = d.Get("policy_class_name")
	request["config"] = d.Get("policy_config")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_policy", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.policyId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudApigPolicyRead(d, meta)
}

func resourceAliCloudApigPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := ApigServiceV2{client}

	objectRaw, err := apigServiceV2.DescribeApigPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_policy DescribeApigPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("policy_class_id", objectRaw["classId"])
	d.Set("policy_class_name", objectRaw["className"])
	d.Set("policy_config", objectRaw["config"])
	d.Set("policy_name", objectRaw["name"])

	return nil
}

func resourceAliCloudApigPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	policyId := d.Id()
	action := fmt.Sprintf("/v1/policies/%s", policyId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if v, ok := d.GetOk("environment_id"); ok {
		request["environmentId"] = v
	}
	request["attachResourceType"] = d.Get("attach_resource_type")
	if d.HasChange("policy_name") {
		update = true
	}
	if v, ok := d.GetOk("policy_name"); ok || d.HasChange("policy_name") {
		request["name"] = v
	}
	if v, ok := d.GetOk("attach_resource_ids"); ok {
		attachResourceIdsMapsArray := convertToInterfaceArray(v)

		request["attachResourceIds"] = attachResourceIdsMapsArray
	}

	if v, ok := d.GetOk("gateway_id"); ok {
		request["gatewayId"] = v
	}
	if d.HasChange("policy_config") {
		update = true
	}
	request["config"] = d.Get("policy_config")
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

	return resourceAliCloudApigPolicyRead(d, meta)
}

func resourceAliCloudApigPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	policyId := d.Id()
	action := fmt.Sprintf("/v1/policies/%s", policyId)
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
