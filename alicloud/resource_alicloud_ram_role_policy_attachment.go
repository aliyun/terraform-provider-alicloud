package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudRamRolePolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamRolePolicyAttachmentCreate,
		Read:   resourceAlicloudRamRolePolicyAttachmentRead,
		//Update: resourceAlicloudRamRolePolicyAttachmentUpdate,
		Delete: resourceAlicloudRamRolePolicyAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamName,
			},
			"policy_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamPolicyName,
			},
			"policy_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validatePolicyType,
			},
		},
	}
}

func resourceAlicloudRamRolePolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	args := ram.AttachPolicyToRoleRequest{
		PolicyRequest: ram.PolicyRequest{
			PolicyName: d.Get("policy_name").(string),
			PolicyType: ram.Type(d.Get("policy_type").(string)),
		},
		RoleName: d.Get("role_name").(string),
	}

	_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
		return ramClient.AttachPolicyToRole(args)
	})
	if err != nil {
		return fmt.Errorf("AttachPolicyToRole got an error: %#v", err)
	}
	d.SetId("role" + args.PolicyName + string(args.PolicyType) + args.RoleName)

	return resourceAlicloudRamRolePolicyAttachmentRead(d, meta)
}

func resourceAlicloudRamRolePolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := ram.RoleQueryRequest{
		RoleName: d.Get("role_name").(string),
	}

	raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
		return ramClient.ListPoliciesForRole(args)
	})
	if err != nil {
		return fmt.Errorf("Get list policies for role got an error: %v", err)
	}
	response, _ := raw.(ram.PolicyListResponse)
	if len(response.Policies.Policy) > 0 {
		for _, v := range response.Policies.Policy {
			if v.PolicyName == d.Get("policy_name").(string) && v.PolicyType == d.Get("policy_type").(string) {
				d.Set("role_name", args.RoleName)
				d.Set("policy_name", v.PolicyName)
				d.Set("policy_type", v.PolicyType)
				return nil
			}
		}
	}

	d.SetId("")
	return nil
}

func resourceAlicloudRamRolePolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := ram.AttachPolicyToRoleRequest{
		PolicyRequest: ram.PolicyRequest{
			PolicyName: d.Get("policy_name").(string),
			PolicyType: ram.Type(d.Get("policy_type").(string)),
		},
		RoleName: d.Get("role_name").(string),
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.DetachPolicyFromRole(args)
		})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting role policy attachment: %#v", err))
		}

		raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.ListPoliciesForRole(ram.RoleQueryRequest{RoleName: args.RoleName})
		})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		response, _ := raw.(ram.PolicyListResponse)

		if len(response.Policies.Policy) < 1 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("Error deleting role policy attachment - trying again while it is deleted."))
	})
}
