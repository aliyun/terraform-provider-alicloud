package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"time"
)

func resourceAlicloudRamGroupPolicyAtatchment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamGroupPolicyAttachmentCreate,
		Read:   resourceAlicloudRamGroupPolicyAttachmentRead,
		Delete: resourceAlicloudRamGroupPolicyAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"group_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamGroupName,
			},
			"policy_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamPolicyName,
			},
			"policy_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validatePolicyType,
			},
		},
	}
}

func resourceAlicloudRamGroupPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.AttachPolicyToGroupRequest{
		PolicyRequest: ram.PolicyRequest{
			PolicyName: d.Get("policy_name").(string),
			PolicyType: ram.Type(d.Get("policy_type").(string)),
		},
		GroupName: d.Get("group_name").(string),
	}

	if _, err := conn.AttachPolicyToGroup(args); err != nil {
		return fmt.Errorf("AttachPolicyToGroup got an error: %#v", err)
	}
	d.SetId("group" + args.PolicyName + string(args.PolicyType) + args.GroupName)

	return resourceAlicloudRamGroupPolicyAttachmentRead(d, meta)
}

func resourceAlicloudRamGroupPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.GroupQueryRequest{
		GroupName: d.Get("group_name").(string),
	}

	response, err := conn.ListPoliciesForGroup(args)
	if err != nil {
		if RamEntityNotExist(err) {
			d.SetId("")
		}
		return fmt.Errorf("Get list policies for group got an error: %#v", err)
	}

	if len(response.Policies.Policy) > 0 {
		for _, v := range response.Policies.Policy {
			if v.PolicyName == d.Get("policy_name").(string) && v.PolicyType == d.Get("policy_type").(string) {
				d.Set("group_name", args.GroupName)
				d.Set("policy_name", v.PolicyName)
				d.Set("policy_type", v.PolicyType)
				return nil
			}
		}
	}

	d.SetId("")
	return nil
}

func resourceAlicloudRamGroupPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.AttachPolicyToGroupRequest{
		PolicyRequest: ram.PolicyRequest{
			PolicyName: d.Get("policy_name").(string),
			PolicyType: ram.Type(d.Get("policy_type").(string)),
		},
		GroupName: d.Get("group_name").(string),
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.DetachPolicyFromGroup(args); err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting group policy attachment: %#v", err))
		}

		response, err := conn.ListPoliciesForGroup(ram.GroupQueryRequest{GroupName: args.GroupName})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}

			return resource.NonRetryableError(err)
		}

		if len(response.Policies.Policy) < 1 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("Error deleting group policy attachment - trying again while it is deleted."))
	})
}
