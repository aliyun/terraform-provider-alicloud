package alicloud

import (
	"fmt"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudTagPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudTagAttachmentCreate,
		Read:   resourceAlicloudTagPolicyAttachmentRead,
		Update: resourceAlicloudTagPolicyAttachmentUpdate,
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
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"USER", "ROOT", "FOLDER", "ACCOUNT"}, false),
			},
			"target_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudTagAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AttachPolicy"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	conn, err := client.NewTagClient()
	if err != nil {
		return WrapError(err)
	}
	request["PolicyId"] = d.Get("policy_id")
	if v, ok := d.GetOk("target_type"); ok {
		request["TargetType"] = v
	}
	if v, ok := d.GetOk("target_id"); ok {
		request["TargetId"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-08-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

func resourceAlicloudTagPolicyAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	return WrapErrorf(Error("tag_policy_attachment not support modify operation"), DefaultErrorMsg, "policy_id", "Modify", AlibabaCloudSdkGoERROR)
}

func resourceAlicloudTagPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	tagService := TagService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	policyId, targetId, targetType := parts[0], parts[1], parts[2]
	object, err := tagService.DescribeTagPolicyAttachment(policyId)
	if err != nil {
		return WrapError(err)
	}

	for _, v := range object.([]interface{}) {
		m := v.(map[string]interface{})
		if m["TargetId"] == targetId && m["TargetType"] == targetType {
			d.Set("PolicyId", policyId)
			d.Set("TargetId", m["TargetId"])
			d.Set("TargetType", m["TargetType"])
			return nil
		}
	}
	return nil
}
func resourceAlicloudTagPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DetachPolicy"
	var response map[string]interface{}
	conn, err := client.NewTagClient()
	if err != nil {
		return WrapError(err)
	}
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
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-08-28"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
