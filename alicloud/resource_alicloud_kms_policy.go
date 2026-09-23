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

func resourceAliCloudKmsPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudKmsPolicyCreate,
		Read:   resourceAliCloudKmsPolicyRead,
		Update: resourceAliCloudKmsPolicyUpdate,
		Delete: resourceAliCloudKmsPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_control_rules": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"kms_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permissions": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudKmsPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreatePolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["Name"] = d.Get("policy_name")

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["KmsInstance"] = d.Get("kms_instance_id")
	jsonPathResult2, err := jsonpath.Get("$", d.Get("permissions"))
	if err != nil {
		return WrapError(err)
	}
	request["Permissions"], _ = convertArrayObjectToJsonString(jsonPathResult2)

	jsonPathResult3, err := jsonpath.Get("$", d.Get("resources"))
	if err != nil {
		return WrapError(err)
	}
	request["Resources"], _ = convertArrayObjectToJsonString(jsonPathResult3)

	request["AccessControlRules"] = d.Get("access_control_rules")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_policy", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Name"]))

	return resourceAliCloudKmsPolicyRead(d, meta)
}

func resourceAliCloudKmsPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsServiceV2 := KmsServiceV2{client}

	objectRaw, err := kmsServiceV2.DescribeKmsPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_kms_policy DescribeKmsPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("access_control_rules", objectRaw["AccessControlRules"])
	d.Set("description", objectRaw["Description"])
	d.Set("kms_instance_id", objectRaw["KmsInstance"])
	d.Set("policy_name", objectRaw["Name"])
	permissions1Raw := make([]interface{}, 0)
	if objectRaw["Permissions"] != nil {
		permissions1Raw = objectRaw["Permissions"].([]interface{})
	}

	d.Set("permissions", permissions1Raw)
	resources1Raw := make([]interface{}, 0)
	if objectRaw["Resources"] != nil {
		resources1Raw = objectRaw["Resources"].([]interface{})
	}

	d.Set("resources", resources1Raw)

	return nil
}

func resourceAliCloudKmsPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	action := "UpdatePolicy"
	var err error
	request = make(map[string]interface{})
	request["Name"] = d.Id()
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("permissions") {
		update = true
	}
	jsonPathResult1, err := jsonpath.Get("$", d.Get("permissions"))
	if err != nil {
		return WrapError(err)
	}
	request["Permissions"], _ = convertArrayObjectToJsonString(jsonPathResult1)

	if d.HasChange("resources") {
		update = true
	}
	jsonPathResult2, err := jsonpath.Get("$", d.Get("resources"))
	if err != nil {
		return WrapError(err)
	}
	request["Resources"], _ = convertArrayObjectToJsonString(jsonPathResult2)

	if d.HasChange("access_control_rules") {
		update = true
	}
	request["AccessControlRules"] = d.Get("access_control_rules")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)

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

	return resourceAliCloudKmsPolicyRead(d, meta)
}

func resourceAliCloudKmsPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeletePolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	request["Name"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)

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
		if IsExpectedErrors(err, []string{"Forbidden.ResourceNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
