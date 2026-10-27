package tag

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourcePolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: ResourcePolicyAttachmentCreate,
		Read:   ResourcePolicyAttachmentRead,
		Delete: ResourcePolicyAttachmentDelete,
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

func ResourcePolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AttachPolicy"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	var err error
	request["PolicyId"] = d.Get("policy_id")
	request["TargetType"] = d.Get("target_type")
	request["TargetId"] = d.Get("target_id")

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_tag_policy_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprintf("%v:%v:%v", request["PolicyId"], request["TargetId"], request["TargetType"]))
	return ResourcePolicyAttachmentRead(d, meta)
}

func ResourcePolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	tagService := NewTagService(client)
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	policyId := parts[0]
	object, err := tagService.DescribeTagPolicyAttachment(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("policy_id", policyId)
	d.Set("target_id", object["TargetId"])
	d.Set("target_type", object["TargetType"])
	return nil
}
func ResourcePolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DetachPolicy"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	policyId, targetId, targetType := parts[0], parts[1], parts[2]
	request["PolicyId"] = policyId
	request["TargetId"] = targetId
	request["TargetType"] = targetType
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Tag", "2018-08-28", action, nil, request, false)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Policy", "EntityNotExists.Target"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
