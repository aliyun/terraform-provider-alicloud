package resourcemanager

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourcePolicyVersion() *schema.Resource {
	return &schema.Resource{
		Create: ResourcePolicyVersionCreate,
		Read:   ResourcePolicyVersionRead,
		Update: ResourcePolicyVersionUpdate,
		Delete: ResourcePolicyVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"is_default_version": {
				Type:       schema.TypeBool,
				Optional:   true,
				Default:    false,
				Deprecated: "Field 'is_default_version' has been deprecated from provider version 1.90.0",
			},
			"policy_document": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func ResourcePolicyVersionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePolicyVersion"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOkExists("is_default_version"); ok {
		request["SetAsDefault"] = v
	}

	request["PolicyDocument"] = d.Get("policy_document")
	request["PolicyName"] = d.Get("policy_name")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_resource_manager_policy_version", action, errs.AlibabaCloudSdkGoERROR)
	}
	responsePolicyVersion := response["PolicyVersion"].(map[string]interface{})
	d.SetId(fmt.Sprint(request["PolicyName"], ":", responsePolicyVersion["VersionId"]))

	return ResourcePolicyVersionRead(d, meta)
}
func ResourcePolicyVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	object, err := resourcemanagerService.DescribeResourceManagerPolicyVersion(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_policy_version resourcemanagerService.DescribeResourceManagerPolicyVersion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("policy_name", parts[0])
	d.Set("is_default_version", object["IsDefaultVersion"])
	d.Set("policy_document", object["PolicyDocument"])
	return nil
}
func ResourcePolicyVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourcemanagerService := ResourcemanagerService{client}
	var err error
	var response map[string]interface{}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	if d.HasChange("is_default_version") {
		object, err := resourcemanagerService.DescribeResourceManagerPolicyVersion(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}
		target := strconv.FormatBool(d.Get("is_default_version").(bool))
		if strconv.FormatBool(object["IsDefaultVersion"].(bool)) != target {
			if target == "true" {
				request := map[string]interface{}{
					"PolicyName": parts[0],
					"VersionId":  parts[1],
				}
				action := "SetDefaultPolicyVersion"
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					logs.AddDebug(action, response, request)
					return nil
				})
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
			}
		}
	}
	return ResourcePolicyVersionRead(d, meta)
}
func ResourcePolicyVersionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DeletePolicyVersion"
	var response map[string]interface{}
	request := map[string]interface{}{
		"PolicyName": parts[0],
		"VersionId":  parts[1],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExist.Policy.Version"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
