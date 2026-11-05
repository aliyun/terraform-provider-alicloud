package ram

import (
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	ramsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceUser() *schema.Resource {
	return &schema.Resource{
		Create: ResourceUserCreate,
		Read:   ResourceUserRead,
		Update: ResourceUserUpdate,
		Delete: ResourceUserDelete,
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
				ValidateFunc: validators.StringLenBetween(0, 128),
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": tags.TagsSchema(),
		},
	}
}

func ResourceUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}

	request := ramsdk.CreateCreateUserRequest()
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

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	var err error
	var raw interface{}
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		raw, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
			return ramClient.CreateUser(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ram_user", request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*ramsdk.CreateUserResponse)

	d.SetId(response.User.UserId)

	err = ramService.WaitForRamUser(d.Id(), Normal, defaultTimeout)
	if err != nil {
		return errs.WrapError(err)
	}

	// Tags on a RAM user are managed via the IMS (2019-08-15) tagging API
	// (ResourceType=user), because the legacy RAM (2015-05-01) TagResources
	// API does not enumerate "user" as a taggable resource type.
	if _, ok := d.GetOk("tags"); ok {
		imsService := RamServiceV2{client}
		if err := imsService.SetImsUserTags(d); err != nil {
			return errs.WrapError(err)
		}
	}

	return ResourceUserRead(d, meta)
}

func ResourceUserUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("Ram", "2015-05-01", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("tags") {
		imsService := RamServiceV2{client}
		if err := imsService.SetImsUserTags(d); err != nil {
			return errs.WrapError(err)
		}
	}

	return ResourceUserRead(d, meta)
}

func ResourceUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ramService := &RamService{client: client}
	object, err := ramService.DescribeRamUser(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.SetId(object.UserId)
	d.Set("name", object.UserName)
	d.Set("display_name", object.DisplayName)
	d.Set("mobile", object.MobilePhone)
	d.Set("email", object.Email)
	d.Set("comments", object.Comments)

	// Tags on a RAM user are read via the IMS (2019-08-15) ListTagResources
	// API with ResourceType=user and the RAM user id as the resource id.
	imsService := RamServiceV2{client}
	tagsRaw, err := imsService.DescribeImsUserTags(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}
	if tagsRaw != nil {
		// IMS (2019-08-15) ListTagResources nests the tag list under
		// TagResources.TagResource, unlike the legacy RAM (2015-05-01)
		// API where TagResources is the list directly. Align with the
		// canonical ims tag parse used by cdn_domain_new,
		// common_bandwidth_package and ddoscoo_instance.
		tagsMaps, _ := jsonpath.Get("$.TagResources.TagResource", tagsRaw)
		d.Set("tags", tags.TagsToMap(tagsMaps))
	}

	return nil
}

func ResourceUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var raw interface{}
	var wait func()

	ramService := &RamService{client: client}
	object, err := ramService.DescribeRamUser(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapError(err)
	}

	userName := object.UserName

	request := ramsdk.CreateListAccessKeysRequest()
	request.RegionId = client.RegionId
	request.UserName = userName

	if d.Get("force").(bool) {
		// list and delete access keys for this user
		wait = helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
			raw, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
				return ramClient.ListAccessKeys(request)
			})
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

		listAccessKeysResponse, _ := raw.(*ramsdk.ListAccessKeysResponse)
		if len(listAccessKeysResponse.AccessKeys.AccessKey) > 0 {
			for _, v := range listAccessKeysResponse.AccessKeys.AccessKey {
				request := ramsdk.CreateDeleteAccessKeyRequest()
				request.RegionId = client.RegionId
				request.UserAccessKeyId = v.AccessKeyId
				request.UserName = userName

				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait = helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
					raw, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
						return ramClient.DeleteAccessKey(request)
					})
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})

				if err != nil && !errs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
				}
				logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
			}
		}

		// list and delete policies for this user
		request := ramsdk.CreateListPoliciesForUserRequest()
		request.RegionId = client.RegionId
		request.UserName = userName

		wait = helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
			raw, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
				return ramClient.ListPoliciesForUser(request)
			})
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

		listPoliciesForUserResponse, _ := raw.(*ramsdk.ListPoliciesForUserResponse)
		if len(listPoliciesForUserResponse.Policies.Policy) > 0 {
			for _, v := range listPoliciesForUserResponse.Policies.Policy {
				request := ramsdk.CreateDetachPolicyFromUserRequest()
				request.RegionId = client.RegionId
				request.PolicyName = v.PolicyName
				request.PolicyType = v.PolicyType
				request.UserName = userName

				wait = helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
					raw, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
						return ramClient.DetachPolicyFromUser(request)
					})
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})

				if err != nil && !errs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
				}
				logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
			}
		}

		// list and delete groups for this user
		listGroupsForUserRequest := ramsdk.CreateListGroupsForUserRequest()
		listGroupsForUserRequest.RegionId = client.RegionId
		listGroupsForUserRequest.UserName = userName

		wait = helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
			raw, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
				return ramClient.ListGroupsForUser(listGroupsForUserRequest)
			})
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), listGroupsForUserRequest.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(listGroupsForUserRequest.GetActionName(), raw, listGroupsForUserRequest.RpcRequest, listGroupsForUserRequest)

		listGroupsForUserResponse, _ := raw.(*ramsdk.ListGroupsForUserResponse)
		if len(listGroupsForUserResponse.Groups.Group) > 0 {
			for _, v := range listGroupsForUserResponse.Groups.Group {
				request := ramsdk.CreateRemoveUserFromGroupRequest()
				request.RegionId = client.RegionId
				request.UserName = userName
				request.GroupName = v.GroupName

				wait = helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
					raw, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
						return ramClient.RemoveUserFromGroup(request)
					})
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				if err != nil && !errs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
				}
				logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
			}
		}

		// delete login profile for this user
		deleteLoginProfileRequest := ramsdk.CreateDeleteLoginProfileRequest()
		deleteLoginProfileRequest.RegionId = client.RegionId
		deleteLoginProfileRequest.UserName = userName

		wait = helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
			raw, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
				return ramClient.DeleteLoginProfile(deleteLoginProfileRequest)
			})
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil && !errs.IsExpectedErrors(err, []string{"EntityNotExist.User.LoginProfile"}) {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), deleteLoginProfileRequest.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(deleteLoginProfileRequest.GetActionName(), raw)
		// unbind MFA device for this user
		unbindMFADeviceRequest := ramsdk.CreateUnbindMFADeviceRequest()
		unbindMFADeviceRequest.UserName = userName

		wait = helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
			raw, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
				return ramClient.UnbindMFADevice(unbindMFADeviceRequest)
			})
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})

		if err != nil && !errs.IsExpectedErrors(err, []string{"EntityNotExist", "EntityNotExist.User.MFADevice"}) {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), unbindMFADeviceRequest.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(unbindMFADeviceRequest.GetActionName(), raw, deleteLoginProfileRequest.RpcRequest, deleteLoginProfileRequest)
	}
	deleteUserRequest := ramsdk.CreateDeleteUserRequest()
	deleteUserRequest.RegionId = client.RegionId

	wait = helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		raw, err = client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
			deleteUserRequest.UserName = userName
			return ramClient.DeleteUser(deleteUserRequest)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"DeleteConflict.User.AccessKey", "DeleteConflict.User.Group", "DeleteConflict.User.Policy", "DeleteConflict.User.LoginProfile", "DeleteConflict.User.MFADevice"}) {
			return errs.WrapError(errs.Error("The user can not be deleted if he has any access keys, login profile, groups, policies, or MFA device attached. You can force the deletion of the user by setting force equals true."))
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(deleteUserRequest.GetActionName(), raw, deleteUserRequest.RpcRequest, deleteUserRequest)
	return errs.WrapError(ramService.WaitForRamUser(d.Id(), Deleted, defaultTimeout))
}
