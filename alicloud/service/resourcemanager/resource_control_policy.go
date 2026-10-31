// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package resourcemanager

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceControlPolicy() *schema.Resource {
	return &schema.Resource{
		Create: ResourceControlPolicyCreate,
		Read:   ResourceControlPolicyRead,
		Update: ResourceControlPolicyUpdate,
		Delete: ResourceControlPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"control_policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"effect_scope": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_document": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"tags": tags.TagsSchema(),
		},
	}
}

func ResourceControlPolicyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateControlPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["PolicyDocument"] = d.Get("policy_document")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["PolicyName"] = d.Get("control_policy_name")
	request["EffectScope"] = d.Get("effect_scope")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_resource_manager_control_policy", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.ControlPolicy.PolicyId", response)
	d.SetId(fmt.Sprint(id))

	return ResourceControlPolicyUpdate(d, meta)
}

func ResourceControlPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerServiceV2 := ResourceManagerServiceV2{client}

	objectRaw, err := resourceManagerServiceV2.DescribeResourceManagerControlPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_manager_control_policy DescribeResourceManagerControlPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("control_policy_name", objectRaw["PolicyName"])
	d.Set("create_time", objectRaw["CreateDate"])
	d.Set("description", objectRaw["Description"])
	d.Set("effect_scope", objectRaw["EffectScope"])
	d.Set("policy_document", objectRaw["PolicyDocument"])

	objectRaw, err = resourceManagerServiceV2.DescribeControlPolicyListTagResources(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	tagsMaps := objectRaw["TagResources"]
	d.Set("tags", tags.TagsToMap(tagsMaps))

	return nil
}

func ResourceControlPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateControlPolicy"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["PolicyId"] = d.Id()

	if !d.IsNewResource() && d.HasChange("policy_document") {
		update = true
	}
	request["NewPolicyDocument"] = d.Get("policy_document")
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["NewDescription"] = d.Get("description")
	}

	if !d.IsNewResource() && d.HasChange("control_policy_name") {
		update = true
	}
	request["NewPolicyName"] = d.Get("control_policy_name")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		resourceManagerServiceV2 := ResourceManagerServiceV2{client}
		if err := resourceManagerServiceV2.SetResourceTags(d, "ControlPolicy"); err != nil {
			return errs.WrapError(err)
		}
	}
	return ResourceControlPolicyRead(d, meta)
}

func ResourceControlPolicyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteControlPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["PolicyId"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, query, request, true)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ConcurrentCallNotSupported"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.ControlPolicy"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
