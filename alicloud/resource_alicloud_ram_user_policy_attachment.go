package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
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

func resourceAlicloudRamUserPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateAttachPolicyToUserRequest()
	request.UserName = d.Get("user_name").(string)
	request.PolicyName = d.Get("policy_name").(string)
	request.PolicyType = d.Get("policy_type").(string)

	_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.AttachPolicyToUser(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "ram_user_policy_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s%s%s", request.UserName, COLON_SEPARATED, request.PolicyName, COLON_SEPARATED, request.PolicyType))
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
	request := ram.CreateListPoliciesForUserRequest()
	request.UserName = split[0]

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListPoliciesForUser(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*ram.ListPoliciesForUserResponse)
	if len(response.Policies.Policy) > 0 {
		for _, v := range response.Policies.Policy {
			if v.PolicyName == d.Get("policy_name").(string) && v.PolicyType == d.Get("policy_type").(string) {
				d.Set("user_name", request.UserName)
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

	request := ram.CreateDetachPolicyFromUserRequest()
	request.UserName = split[0]
	request.PolicyName = split[1]
	request.PolicyType = split[2]

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.DetachPolicyFromUser(request)
		})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			request := ram.CreateListPoliciesForUserRequest()
			request.UserName = split[0]
			return ramClient.ListPoliciesForUser(request)
		})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}

			return resource.NonRetryableError(WrapError(err))
		}
		response, _ := raw.(*ram.ListPoliciesForUserResponse)
		if len(response.Policies.Policy) < 1 {
			return nil
		}
		return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR))
	})
}
