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

func resourceAlicloudTagPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudTagPolicyCreate,
		Read:   resourceAlicloudTagPolicyRead,
		Update: resourceAlicloudTagPolicyUpdate,
		Delete: resourceAlicloudTagPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_content": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"user_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"USER", "RD"}, false),
				ForceNew:     true,
				Computed:     true,
			},
		},
	}
}

func resourceAlicloudTagPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePolicy"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	var err error
	if v, ok := d.GetOk("policy_desc"); ok {
		request["PolicyDesc"] = v
	}
	request["PolicyContent"] = d.Get("policy_content")
	request["PolicyName"] = d.Get("policy_name")
	if v, ok := d.GetOk("user_type"); ok {
		request["UserType"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_resource_tag_policy", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(response["PolicyId"]))
	return resourceAlicloudTagPolicyRead(d, meta)
}
func resourceAlicloudTagPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	tagService := TagService{client}
	object, err := tagService.DescribeTagPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_resource_tag_policy tagService.DescribeTagPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("policy_name", object["PolicyName"])
	d.Set("policy_desc", object["PolicyDesc"])
	d.Set("policy_content", object["PolicyContent"])
	d.Set("user_type", object["UserType"])
	return nil
}

func resourceAlicloudTagPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}

	request := map[string]interface{}{
		"PolicyId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	update := false
	if d.HasChange("policy_desc") {
		update = true
		if v, ok := d.GetOk("policy_desc"); ok {
			request["PolicyDesc"] = v
		}
	}
	if d.HasChange("policy_content") {
		request["PolicyContent"] = d.Get("policy_content")
		update = true
	}
	if d.HasChange("policy_name") {
		request["PolicyName"] = d.Get("policy_name")
		update = true
	}

	action := "ModifyPolicy"

	if update {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, false)
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
	return resourceAlicloudTagPolicyRead(d, meta)
}
func resourceAlicloudTagPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePolicy"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"PolicyId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, false)
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
