package alicloud

import (
	"time"

	"github.com/alibabacloud-go/tea/tea"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
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
			"comments": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringLenBetween(0, 128),
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudRamUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}

	request := ram.CreateCreateUserRequest()
	request.RegionId = client.RegionId
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	var err error
	var raw interface{}
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.CreateUser(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_user", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

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
	var response map[string]interface{}
	var err error
	update := false
	request := map[string]interface{}{
		"UserName":    d.Get("name"),
		"NewUserName": d.Get("name"),
	}

	if !d.IsNewResource() && d.HasChange("name") {
		update = true

		oldName, newName := d.GetChange("name")
		request["UserName"] = oldName.(string)
		request["NewUserName"] = newName.(string)
	}

	if d.HasChange("display_name") {
		update = true

		if v, ok := d.GetOk("display_name"); ok {
			request["NewDisplayName"] = v
		} else {
			request["NewDisplayName"] = tea.String("")
		}
	}

	if d.HasChange("mobile") {
		update = true
		if v, ok := d.GetOk("mobile"); ok {
			request["NewMobilePhone"] = v
		} else {
			request["NewMobilePhone"] = tea.String("")
		}
	}

	if d.HasChange("email") {
		update = true

		if v, ok := d.GetOk("email"); ok {
			request["NewEmail"] = v
		} else {
			request["NewEmail"] = tea.String("")
		}
	}

	if d.HasChange("comments") {
		update = true

		if v, ok := d.GetOk("comments"); ok {
			request["NewComments"] = v
		} else {
			request["NewComments"] = tea.String("")
		}
	}

	if update {
		action := "UpdateUser"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Ram", "2015-05-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
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

	var raw interface{}
	var wait func()

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
	request.RegionId = client.RegionId
	request.UserName = userName

	if d.Get("force").(bool) {
		// list and delete access keys for this user
		wait = incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
			raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
				return ramClient.ListAccessKeys(request)
			})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		listAccessKeysResponse, _ := raw.(*ram.ListAccessKeysResponse)
		if len(listAccessKeysResponse.AccessKeys.AccessKey) > 0 {
			for _, v := range listAccessKeysResponse.AccessKeys.AccessKey {
				request := ram.CreateDeleteAccessKeyRequest()
				request.RegionId = client.RegionId
				request.UserAccessKeyId = v.AccessKeyId
				request.UserName = userName

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait = incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
					raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
						return ramClient.DeleteAccessKey(request)
					})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})

				if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist"}) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			}
		}

		// list and delete policies for this user
		request := ram.CreateListPoliciesForUserRequest()
		request.RegionId = client.RegionId
		request.UserName = userName

		wait = incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
			raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
				return ramClient.ListPoliciesForUser(request)
			})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		listPoliciesForUserResponse, _ := raw.(*ram.ListPoliciesForUserResponse)
		if len(listPoliciesForUserResponse.Policies.Policy) > 0 {
			for _, v := range listPoliciesForUserResponse.Policies.Policy {
				request := ram.CreateDetachPolicyFromUserRequest()
				request.RegionId = client.RegionId
				request.PolicyName = v.PolicyName
				request.PolicyType = v.PolicyType
				request.UserName = userName

				wait = incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
					raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
						return ramClient.DetachPolicyFromUser(request)
					})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})

				if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist"}) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			}
		}

		// list and delete groups for this user
		listGroupsForUserRequest := ram.CreateListGroupsForUserRequest()
		listGroupsForUserRequest.RegionId = client.RegionId
		listGroupsForUserRequest.UserName = userName

		wait = incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
			raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
				return ramClient.ListGroupsForUser(listGroupsForUserRequest)
			})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), listGroupsForUserRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(listGroupsForUserRequest.GetActionName(), raw, listGroupsForUserRequest.RpcRequest, listGroupsForUserRequest)

		listGroupsForUserResponse, _ := raw.(*ram.ListGroupsForUserResponse)
		if len(listGroupsForUserResponse.Groups.Group) > 0 {
			for _, v := range listGroupsForUserResponse.Groups.Group {
				request := ram.CreateRemoveUserFromGroupRequest()
				request.RegionId = client.RegionId
				request.UserName = userName
				request.GroupName = v.GroupName

				wait = incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
					raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
						return ramClient.RemoveUserFromGroup(request)
					})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist"}) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			}
		}

		// delete login profile for this user
		deleteLoginProfileRequest := ram.CreateDeleteLoginProfileRequest()
		deleteLoginProfileRequest.RegionId = client.RegionId
		deleteLoginProfileRequest.UserName = userName

		wait = incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
			raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
				return ramClient.DeleteLoginProfile(deleteLoginProfileRequest)
			})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist.User.LoginProfile"}) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), deleteLoginProfileRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(deleteLoginProfileRequest.GetActionName(), raw)
		// unbind MFA device for this user
		unbindMFADeviceRequest := ram.CreateUnbindMFADeviceRequest()
		unbindMFADeviceRequest.UserName = userName

		wait = incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
			raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
				return ramClient.UnbindMFADevice(unbindMFADeviceRequest)
			})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist", "EntityNotExist.User.MFADevice"}) {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), unbindMFADeviceRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(unbindMFADeviceRequest.GetActionName(), raw, deleteLoginProfileRequest.RpcRequest, deleteLoginProfileRequest)
	}
	deleteUserRequest := ram.CreateDeleteUserRequest()
	deleteUserRequest.RegionId = client.RegionId

	wait = incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			deleteUserRequest.UserName = userName
			return ramClient.DeleteUser(deleteUserRequest)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"DeleteConflict.User.AccessKey", "DeleteConflict.User.Group", "DeleteConflict.User.Policy", "DeleteConflict.User.LoginProfile", "DeleteConflict.User.MFADevice"}) {
			return WrapError(Error("The user can not be deleted if he has any access keys, login profile, groups, policies, or MFA device attached. You can force the deletion of the user by setting force equals true."))
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(deleteUserRequest.GetActionName(), raw, deleteUserRequest.RpcRequest, deleteUserRequest)
	return WrapError(ramService.WaitForRamUser(d.Id(), Deleted, DefaultTimeout))
}
