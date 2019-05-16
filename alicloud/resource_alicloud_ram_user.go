package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
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
	ramService := RamService{client}

	request := ram.CreateCreateUserRequest()
	request.UserName = d.Get("name").(string)
	if v, ok := d.GetOk("display_name"); ok {
		request.DisplayName = v.(string)
	}
	if v, ok := d.GetOk("mobile"); ok {
		request.MobilePhone = v.(string)
	}
	if v, ok := d.GetOk("email"); ok {
		request.Email = v.(string)
	}
	if v, ok := d.GetOk("comments"); ok {
		request.Comments = v.(string)
	}

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.CreateUser(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_user", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*ram.CreateUserResponse)

	d.SetId(response.User.UserId)

	err = ramService.WaitForRamUser(d.Id(), Normal, DefaultTimeout)
	if err != nil {
		return WrapError(err)
	}
	return resourceAlicloudRamUserRead(d, meta)
}

func resourceAlicloudRamUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateUpdateUserRequest()
	request.UserName = d.Get("name").(string)
	request.NewUserName = d.Get("name").(string)

	update := false

	if d.HasChange("name") && !d.IsNewResource() {
		ov, nv := d.GetChange("name")
		request.UserName = ov.(string)
		request.NewUserName = nv.(string)
		update = true
	}

	if d.HasChange("display_name") {
		request.NewDisplayName = d.Get("display_name").(string)
		update = true
	}

	if d.HasChange("mobile") {
		request.NewMobilePhone = d.Get("mobile").(string)
		update = true
	}

	if d.HasChange("email") {
		request.NewEmail = d.Get("email").(string)
		update = true
	}

	if d.HasChange("comments") {
		request.NewComments = d.Get("comments").(string)
		update = true
	}

	if update {

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.UpdateUser(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
	}

	return resourceAlicloudRamUserRead(d, meta)
}

func resourceAlicloudRamUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ramService := &RamService{client: client}
	object, err := ramService.DescribeRamUser(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.SetId(object.UserId)
	d.Set("name", object.UserName)
	d.Set("display_name", object.DisplayName)
	d.Set("mobile", object.MobilePhone)
	d.Set("email", object.Email)
	d.Set("comments", object.Comments)
	return nil
}

func resourceAlicloudRamUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ramService := &RamService{client: client}
	object, err := ramService.DescribeRamUser(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	userName := object.UserName
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
		addDebug(request.GetActionName(), raw)
		listAccessKeysResponse, _ := raw.(*ram.ListAccessKeysResponse)
		if len(listAccessKeysResponse.AccessKeys.AccessKey) > 0 {
			for _, v := range listAccessKeysResponse.AccessKeys.AccessKey {
				request := ram.CreateDeleteAccessKeyRequest()
				request.UserAccessKeyId = v.AccessKeyId
				request.UserName = userName
				raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.DeleteAccessKey(request)
				})
				if err != nil && !RamEntityNotExist(err) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw)
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
		addDebug(request.GetActionName(), raw)
		listPoliciesForUserResponse, _ := raw.(*ram.ListPoliciesForUserResponse)
		if len(listPoliciesForUserResponse.Policies.Policy) > 0 {
			for _, v := range listPoliciesForUserResponse.Policies.Policy {
				request := ram.CreateDetachPolicyFromUserRequest()
				request.PolicyName = v.PolicyName
				request.PolicyType = v.PolicyType
				request.UserName = userName
				raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.DetachPolicyFromUser(request)
				})
				if err != nil && !RamEntityNotExist(err) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw)
			}
		}

		// list and delete groups for this user
		listGroupsForUserRequest := ram.CreateListGroupsForUserRequest()
		listGroupsForUserRequest.UserName = userName
		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListGroupsForUser(listGroupsForUserRequest)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), listGroupsForUserRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(listGroupsForUserRequest.GetActionName(), raw)
		listGroupsForUserResponse, _ := raw.(*ram.ListGroupsForUserResponse)
		if len(listGroupsForUserResponse.Groups.Group) > 0 {
			for _, v := range listGroupsForUserResponse.Groups.Group {
				request := ram.CreateRemoveUserFromGroupRequest()
				request.UserName = userName
				request.GroupName = v.GroupName
				raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.RemoveUserFromGroup(request)
				})
				if err != nil && !RamEntityNotExist(err) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw)
			}
		}

		// delete login profile for this user
		deleteLoginProfileRequest := ram.CreateDeleteLoginProfileRequest()
		deleteLoginProfileRequest.UserName = userName
		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.DeleteLoginProfile(deleteLoginProfileRequest)
		})
		if err != nil && !RamEntityNotExist(err) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), deleteLoginProfileRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(deleteLoginProfileRequest.GetActionName(), raw)
		// unbind MFA device for this user
		unbindMFADeviceRequest := ram.CreateUnbindMFADeviceRequest()
		unbindMFADeviceRequest.UserName = userName
		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.UnbindMFADevice(unbindMFADeviceRequest)
		})
		if err != nil && !RamEntityNotExist(err) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), unbindMFADeviceRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(unbindMFADeviceRequest.GetActionName(), raw)
	}
	deleteUserRequest := ram.CreateDeleteUserRequest()
	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		deleteUserRequest.UserName = userName
		return ramClient.DeleteUser(deleteUserRequest)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{DeleteConflictUserAccessKey, DeleteConflictUserGroup, DeleteConflictUserPolicy, DeleteConflictUserLoginProfile, DeleteConflictUserMFADevice}) {
			return WrapError(Error("The user can not has any access keys or login profile or attached group or attached policies or attached mfa device while deleting the user.- you can set force with true to force delete the user."))
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(deleteUserRequest.GetActionName(), raw)
	return WrapError(ramService.WaitForRamUser(d.Id(), Deleted, DefaultTimeout))
}
