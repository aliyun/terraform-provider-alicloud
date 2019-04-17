package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudRamGroupPolicyAtatchment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamGroupPolicyAttachmentCreate,
		Read:   resourceAlicloudRamGroupPolicyAttachmentRead,
		Delete: resourceAlicloudRamGroupPolicyAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"group_name": {
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

func resourceAlicloudRamGroupPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateAttachPolicyToGroupRequest()
	request.PolicyType = d.Get("policy_type").(string)
	request.PolicyName = d.Get("policy_name").(string)
	request.GroupName = d.Get("group_name").(string)

	_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.AttachPolicyToGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "ram_group_policy_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	d.SetId("group" + request.PolicyName + string(request.PolicyType) + request.GroupName)

	return resourceAlicloudRamGroupPolicyAttachmentRead(d, meta)
}

func resourceAlicloudRamGroupPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateListPoliciesForGroupRequest()
	request.GroupName = d.Get("group_name").(string)

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListPoliciesForGroup(request)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*ram.ListPoliciesForGroupResponse)
	if len(response.Policies.Policy) > 0 {
		for _, v := range response.Policies.Policy {
			if v.PolicyName == d.Get("policy_name").(string) && v.PolicyType == d.Get("policy_type").(string) {
				d.Set("group_name", request.GroupName)
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
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateDetachPolicyFromGroupRequest()
	request.PolicyName = d.Get("policy_name").(string)
	request.PolicyType = d.Get("policy_type").(string)
	request.GroupName = d.Get("group_name").(string)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.DetachPolicyFromGroup(request)
		})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			request := ram.CreateListPoliciesForGroupRequest()
			request.GroupName = d.Get("group_name").(string)
			return ramClient.ListPoliciesForGroup(request)
		})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}

			return resource.NonRetryableError(WrapError(err))
		}
		response, _ := raw.(*ram.ListPoliciesForGroupResponse)
		if len(response.Policies.Policy) < 1 {
			return nil
		}
		return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR))
	})
}
