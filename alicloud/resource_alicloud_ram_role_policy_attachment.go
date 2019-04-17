package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
	request := ram.CreateAttachPolicyToRoleRequest()
	request.RoleName = d.Get("role_name").(string)
	request.PolicyType = d.Get("policy_type").(string)
	request.PolicyName = d.Get("policy_name").(string)

	_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.AttachPolicyToRole(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "ram_role_policy_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	d.SetId("role" + request.PolicyName + string(request.PolicyType) + request.RoleName)

	return resourceAlicloudRamRolePolicyAttachmentRead(d, meta)
}

func resourceAlicloudRamRolePolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateListPoliciesForRoleRequest()
	request.RoleName = d.Get("role_name").(string)

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListPoliciesForRole(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*ram.ListPoliciesForRoleResponse)
	if len(response.Policies.Policy) > 0 {
		for _, v := range response.Policies.Policy {
			if v.PolicyName == d.Get("policy_name").(string) && v.PolicyType == d.Get("policy_type").(string) {
				d.Set("role_name", request.RoleName)
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

	request := ram.CreateDetachPolicyFromRoleRequest()
	request.RoleName = d.Get("role_name").(string)
	request.PolicyType = d.Get("policy_type").(string)
	request.PolicyName = d.Get("policy_name").(string)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.DetachPolicyFromRole(request)
		})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			request := ram.CreateListPoliciesForRoleRequest()
			request.RoleName = d.Get("role_name").(string)
			return ramClient.ListPoliciesForRole(request)
		})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		response, _ := raw.(*ram.ListPoliciesForRoleResponse)

		if len(response.Policies.Policy) < 1 {
			return nil
		}
		return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR))
	})
}
