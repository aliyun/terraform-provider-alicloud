package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudTagPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudTagPolicyAttachmentCreate,
		Read:   resourceAlicloudTagPolicyAttachmentRead,
		Delete: resourceAlicloudTagPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"USER", "ROOT", "FOLDER", "ACCOUNT"}, false),
			},
			"target_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceAlicloudTagPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AttachPolicy"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	var err error
	request["PolicyId"] = d.Get("policy_id")
	request["TargetType"] = d.Get("target_type")
	request["TargetId"] = d.Get("target_id")

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_tag_policy_attachment", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprintf("%v:%v:%v", request["PolicyId"], request["TargetId"], request["TargetType"]))
	return resourceAlicloudTagPolicyAttachmentRead(d, meta)
}

func resourceAlicloudTagPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	tagService := TagService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	policyId := parts[0]
	object, err := tagService.DescribeTagPolicyAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("policy_id", policyId)
	d.Set("target_id", object["TargetId"])
	d.Set("target_type", object["TargetType"])
	return nil
}
func resourceAlicloudTagPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DetachPolicy"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	policyId, targetId, targetType := parts[0], parts[1], parts[2]
	request["PolicyId"] = policyId
	request["TargetId"] = targetId
	request["TargetType"] = targetType
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
		if IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExists.Target"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
