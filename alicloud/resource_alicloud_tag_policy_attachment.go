package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceTagPolicyAttachment() *schema.Resource {
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
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
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

	//d.SetId(fmt.Sprint(request["PolicyId"], ":", request["TargetId"], ":", request["TargetType"]))
	d.SetId(fmt.Sprint(request["PolicyId"]))
	return resourceAlicloudTagPolicyAttachmentRead(d, meta)
}

func resourceAlicloudTagPolicyAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	return WrapErrorf(Error("tag_policy_attachment not support modify operation"), DefaultErrorMsg, "policy_id", "Modify", AlibabaCloudSdkGoERROR)
}

func resourceAlicloudTagPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	tagService := TagService{client}

	object, err := tagService.DescribeTagPolicyAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_tag_policy_attachment tagService.DescribeTagPolicyAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	//parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("PolicyId", d.Id())
	d.Set("TargetId", object["TargetId"])
	d.Set("TargetType", object["TargetType"])
	return nil
}
func resourceAlicloudTagPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	//if len(strings.Split(d.Id(), ":")) != 3 {
	//	d.SetId(fmt.Sprintf("%v:%v", d.Id(), d.Get("policy_id").(string)))
	//}
	//parts, err := ParseResourceId(d.Id(), 3)
	//if err != nil {
	//	return WrapError(err)
	//}
	action := "DetachPolicy"
	var response map[string]interface{}
	conn, err := client.NewTagClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{}
	request["PolicyId"] = d.Id()
	if v, ok := d.GetOk("target_type"); ok {
		request["TargetType"] = v
	}
	if v, ok := d.GetOk("target_id"); ok {
		request["TargetId"] = v
	}
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
