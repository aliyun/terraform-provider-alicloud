package alicloud

import (
	"fmt"
	"time"

	"strings"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudRamUserPolicyAtatchment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamUserPolicyAttachmentCreate,
		Read:   resourceAlicloudRamUserPolicyAttachmentRead,
		Delete: resourceAlicloudRamUserPolicyAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"user_name": {
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

func resourceAlicloudRamUserPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := ram.AttachPolicyRequest{
		PolicyRequest: ram.PolicyRequest{
			PolicyName: d.Get("policy_name").(string),
			PolicyType: ram.Type(d.Get("policy_type").(string)),
		},
		UserName: d.Get("user_name").(string),
	}

	_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
		return ramClient.AttachPolicyToUser(args)
	})
	if err != nil {
		return fmt.Errorf("AttachPolicyToUser got an error: %#v", err)
	}

	d.SetId(fmt.Sprintf("%s%s%s%s%s", args.UserName, COLON_SEPARATED, args.PolicyName, COLON_SEPARATED, args.PolicyType))
	return resourceAlicloudRamUserPolicyAttachmentRead(d, meta)
}

func resourceAlicloudRamUserPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	// In order to be compatible with previous Id (before 1.9.6) which format to user<policuy_name><policy_type><user_name>
	id := fmt.Sprintf("%s%s%s%s%s", d.Get("user_name").(string), COLON_SEPARATED, d.Get("policy_name").(string), COLON_SEPARATED, d.Get("policy_type").(string))

	if d.Id() != id {
		d.SetId(id)
	}

	split := strings.Split(d.Id(), COLON_SEPARATED)
	args := ram.UserQueryRequest{
		UserName: split[0],
	}

	raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
		return ramClient.ListPoliciesForUser(args)
	})
	if err != nil {
		return fmt.Errorf("Get list policies for user got an error: %#v", err)
	}
	response, _ := raw.(ram.PolicyListResponse)
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
	client := meta.(*connectivity.AliyunClient)

	// In order to be compatible with previous Id (before 1.9.6) which format to user<policuy_name><policy_type><user_name>
	id := fmt.Sprintf("%s%s%s%s%s", d.Get("user_name").(string), COLON_SEPARATED, d.Get("policy_name").(string), COLON_SEPARATED, d.Get("policy_type").(string))

	if d.Id() != id {
		d.SetId(id)
	}

	split := strings.Split(d.Id(), COLON_SEPARATED)

	args := ram.AttachPolicyRequest{
		PolicyRequest: ram.PolicyRequest{
			PolicyName: split[1],
			PolicyType: ram.Type(split[2]),
		},
		UserName: split[0],
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.DetachPolicyFromUser(args)
		})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting user policy attachment: %#v", err))
		}

		raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.ListPoliciesForUser(ram.UserQueryRequest{UserName: args.UserName})
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
		return resource.RetryableError(fmt.Errorf("Error deleting user policy attachment - trying again while it is deleted."))
	})
}
