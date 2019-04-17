package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudRamUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamUserCreate,
		Read:   resourceAlicloudRamUserRead,
		Update: resourceAlicloudRamUserUpdate,
		Delete: resourceAlicloudRamUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mobile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"comments": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateComment,
			},
		},
	}
}

func resourceAlicloudRamUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateCreateUserRequest()
	request.UserName = d.Get("name").(string)

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.CreateUser(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "ram_user", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*ram.CreateUserResponse)

	d.SetId(response.User.UserId)
	return resourceAlicloudRamUserUpdate(d, meta)
}

func resourceAlicloudRamUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	request := ram.CreateUpdateUserRequest()
	request.UserName = d.Get("name").(string)
	request.NewUserName = d.Get("name").(string)

	attributeUpdate := false

	if d.HasChange("name") && !d.IsNewResource() {
		ov, nv := d.GetChange("name")
		request.UserName = ov.(string)
		request.NewUserName = nv.(string)
		d.SetPartial("name")
		attributeUpdate = true
	}

	if d.HasChange("display_name") {
		d.SetPartial("display_name")
		request.NewDisplayName = d.Get("display_name").(string)
		attributeUpdate = true
	}

	if d.HasChange("mobile") {
		d.SetPartial("mobile")
		request.NewMobilePhone = d.Get("mobile").(string)
		attributeUpdate = true
	}

	if d.HasChange("email") {
		d.SetPartial("email")
		request.NewEmail = d.Get("email").(string)
		attributeUpdate = true
	}

	if d.HasChange("comments") {
		d.SetPartial("comments")
		request.NewComments = d.Get("comments").(string)
		attributeUpdate = true
	}

	if attributeUpdate {
		_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.UpdateUser(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	d.Partial(false)
	return resourceAlicloudRamUserRead(d, meta)
}

func resourceAlicloudRamUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ramService := &RamService{client: client}
	user, err := ramService.GetSpecifiedUser(d.Id())
	if err != nil {
		if IsExceptedError(err, NotFound) {
			return nil
		}
		return WrapError(err)
	}

	request := ram.CreateGetUserRequest()
	request.UserName = user.UserName
	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetUser(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	resp, _ := raw.(*ram.GetUserResponse)
	user = &resp.User
	d.SetId(user.UserId)
	d.Set("name", user.UserName)
	d.Set("display_name", user.DisplayName)
	d.Set("mobile", user.MobilePhone)
	d.Set("email", user.Email)
	d.Set("comments", user.Comments)
	return nil
}

func resourceAlicloudRamUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ramService := &RamService{client: client}
	user, err := ramService.GetSpecifiedUser(d.Id())
	if err != nil {
		if IsExceptedError(err, NotFound) {
			return nil
		}
		return WrapError(err)
	}

	userName := user.UserName
	request := ram.CreateListAccessKeysRequest()
	request.UserName = userName

	if d.Get("force").(bool) {
		// list and delete access keys for this user
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListAccessKeys(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		akResp, _ := raw.(*ram.ListAccessKeysResponse)
		if len(akResp.AccessKeys.AccessKey) > 0 {
			for _, v := range akResp.AccessKeys.AccessKey {
				request := ram.CreateDeleteAccessKeyRequest()
				request.UserAccessKeyId = v.AccessKeyId
				request.UserName = userName
				_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.DeleteAccessKey(request)
				})
				if err != nil && !RamEntityNotExist(err) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
			}
		}

		// list and delete policies for this user
		request := ram.CreateListPoliciesForUserRequest()
		request.UserName = userName
		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForUser(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		policyResp, _ := raw.(*ram.ListPoliciesForUserResponse)
		if len(policyResp.Policies.Policy) > 0 {
			for _, v := range policyResp.Policies.Policy {
				request := ram.CreateDetachPolicyFromUserRequest()
				request.PolicyName = v.PolicyName
				request.PolicyType = v.PolicyType
				request.UserName = userName
				_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.DetachPolicyFromUser(request)
				})
				if err != nil && !RamEntityNotExist(err) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
			}
		}

		// list and delete groups for this user
		request1 := ram.CreateListGroupsForUserRequest()
		request1.UserName = userName
		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListGroupsForUser(request1)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request1.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		groupResp, _ := raw.(*ram.ListGroupsForUserResponse)
		if len(groupResp.Groups.Group) > 0 {
			for _, v := range groupResp.Groups.Group {
				request := ram.CreateRemoveUserFromGroupRequest()
				request.UserName = userName
				request.GroupName = v.GroupName
				_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.RemoveUserFromGroup(request)
				})
				if err != nil && !RamEntityNotExist(err) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
			}
		}

		// delete login profile for this user
		request2 := ram.CreateDeleteLoginProfileRequest()
		request2.UserName = userName
		_, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.DeleteLoginProfile(request2)
		})
		if err != nil && !RamEntityNotExist(err) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request2.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		// unbind MFA device for this user
		request3 := ram.CreateUnbindMFADeviceRequest()
		request3.UserName = userName
		_, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.UnbindMFADevice(request3)
		})
		if err != nil && !RamEntityNotExist(err) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request3.GetActionName(), AlibabaCloudSdkGoERROR)
		}

	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			request := ram.CreateDeleteUserRequest()
			request.UserName = userName
			return ramClient.DeleteUser(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{DeleteConflictUserAccessKey, DeleteConflictUserGroup, DeleteConflictUserPolicy, DeleteConflictUserLoginProfile, DeleteConflictUserMFADevice}) {
				return resource.RetryableError(WrapError(Error("The user can not has any access keys or login profile or attached group or attached policies or attached mfa device while deleting the user.- you can set force with true to force delete the user.")))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		return nil
	})
}
