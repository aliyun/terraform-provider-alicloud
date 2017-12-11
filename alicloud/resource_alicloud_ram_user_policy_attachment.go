package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"time"
)

func resourceAlicloudRamUserPolicyAtatchment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamUserPolicyAttachmentCreate,
		Read:   resourceAlicloudRamUserPolicyAttachmentRead,
		Delete: resourceAlicloudRamUserPolicyAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"user_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamName,
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

func resourceAlicloudRamUserPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.AttachPolicyRequest{
		PolicyRequest: ram.PolicyRequest{
			PolicyName: d.Get("policy_name").(string),
			PolicyType: ram.Type(d.Get("policy_type").(string)),
		},
		UserName: d.Get("user_name").(string),
	}

	if _, err := conn.AttachPolicyToUser(args); err != nil {
		return fmt.Errorf("AttachPolicyToUser got an error: %#v", err)
	}

	d.SetId("user" + args.PolicyName + string(args.PolicyType) + args.UserName)
	return resourceAlicloudRamUserPolicyAttachmentRead(d, meta)
}

func resourceAlicloudRamUserPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.UserQueryRequest{
		UserName: d.Get("user_name").(string),
	}

	response, err := conn.ListPoliciesForUser(args)
	if err != nil {
		return fmt.Errorf("Get list policies for user got an error: %#v", err)
	}

	if len(response.Policies.Policy) > 0 {
		for _, v := range response.Policies.Policy {
			if v.PolicyName == d.Get("policy_name").(string) && v.PolicyType == d.Get("policy_type").(string) {
				d.Set("user_name", args.UserName)
				d.Set("policy_name", v.PolicyName)
				d.Set("policy_type", v.PolicyType)
				return nil
			}
		}
	}

	d.SetId("")
	return nil
}

func resourceAlicloudRamUserPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.AttachPolicyRequest{
		PolicyRequest: ram.PolicyRequest{
			PolicyName: d.Get("policy_name").(string),
			PolicyType: ram.Type(d.Get("policy_type").(string)),
		},
		UserName: d.Get("user_name").(string),
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.DetachPolicyFromUser(args); err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting user policy attachment: %#v", err))
		}

		response, err := conn.ListPoliciesForUser(ram.UserQueryRequest{UserName: args.UserName})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}

			return resource.NonRetryableError(err)
		}

		if len(response.Policies.Policy) < 1 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("Error deleting user policy attachment - trying again while it is deleted."))
	})
}
