package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudResourceManagerPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudResourceManagerPolicyCreate,
		Read:   resourceAlicloudResourceManagerPolicyRead,
		Update: resourceAlicloudResourceManagerPolicyUpdate,
		Delete: resourceAlicloudResourceManagerPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"default_version": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'default_version' has been deprecated from provider version 1.90.0",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"policy_document": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudResourceManagerPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePolicy"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	request["PolicyDocument"] = d.Get("policy_document")
	request["PolicyName"] = d.Get("policy_name")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_manager_policy", action, AlibabaCloudSdkGoERROR)
	}
	responsePolicy := response["Policy"].(map[string]interface{})
	d.SetId(fmt.Sprint(responsePolicy["PolicyName"]))

	return resourceAlicloudResourceManagerPolicyUpdate(d, meta)
}
func resourceAlicloudResourceManagerPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	object, err := resourcemanagerService.DescribeResourceManagerPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_policy resourcemanagerService.DescribeResourceManagerPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("policy_name", d.Id())
	d.Set("default_version", object["DefaultVersion"])
	d.Set("description", object["Description"])
	d.Set("policy_type", object["PolicyType"])

	getPolicyVersionObject, err := resourcemanagerService.GetPolicyVersion(d.Id(), d)
	if err != nil {
		return WrapError(err)
	}
	d.Set("policy_document", getPolicyVersionObject["PolicyDocument"])
	return nil
}
func resourceAlicloudResourceManagerPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	if !d.IsNewResource() && d.HasChange("default_version") {
		request := map[string]interface{}{
			"PolicyName": d.Id(),
		}
		request["VersionId"] = d.Get("default_version")
		action := "SetDefaultPolicyVersion"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
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
	return resourceAlicloudResourceManagerPolicyRead(d, meta)
}
func resourceAlicloudResourceManagerPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePolicy"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"PolicyName": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
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
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
