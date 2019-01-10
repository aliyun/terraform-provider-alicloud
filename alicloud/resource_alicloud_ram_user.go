package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/ram"
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRamName,
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

	args := ram.UserRequest{
		User: ram.User{
			UserName: d.Get("name").(string),
		},
	}

	raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
		return ramClient.CreateUser(args)
	})
	if err != nil {
		return fmt.Errorf("CreateUser got an error: %#v", err)
	}
	response, _ := raw.(ram.UserResponse)
	d.SetId(response.User.UserName)
	return resourceAlicloudRamUserUpdate(d, meta)
}

func resourceAlicloudRamUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	args := ram.UpdateUserRequest{
		UserName:    d.Id(),
		NewUserName: d.Id(),
	}
	attributeUpdate := false

	if d.HasChange("name") && !d.IsNewResource() {
		ov, nv := d.GetChange("name")
		args.UserName = ov.(string)
		args.NewUserName = nv.(string)
		d.SetPartial("name")
		attributeUpdate = true
	}

	if d.HasChange("display_name") {
		d.SetPartial("display_name")
		args.NewDisplayName = d.Get("display_name").(string)
		attributeUpdate = true
	}

	if d.HasChange("mobile") {
		d.SetPartial("mobile")
		args.NewMobilePhone = d.Get("mobile").(string)
		attributeUpdate = true
	}

	if d.HasChange("email") {
		d.SetPartial("email")
		args.NewEmail = d.Get("email").(string)
		attributeUpdate = true
	}

	if d.HasChange("comments") {
		d.SetPartial("comments")
		args.NewComments = d.Get("comments").(string)
		attributeUpdate = true
	}

	if attributeUpdate {
		_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.UpdateUser(args)
		})
		if err != nil {
			return fmt.Errorf("Update user got an error: %v", err)
		}
	}

	d.Partial(false)
	return resourceAlicloudRamUserRead(d, meta)
}

func resourceAlicloudRamUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := ram.UserQueryRequest{
		UserName: d.Id(),
	}

	raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
		return ramClient.GetUser(args)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("GetUser got an error: %#v", err)
	}
	response, _ := raw.(ram.UserResponse)
	user := response.User
	d.Set("name", user.UserName)
	d.Set("display_name", user.DisplayName)
	d.Set("mobile", user.MobilePhone)
	d.Set("email", user.Email)
	d.Set("comments", user.Comments)
	return nil
}

func resourceAlicloudRamUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	userName := d.Id()
	args := ram.UserQueryRequest{
		UserName: userName,
	}

	if d.Get("force").(bool) {
		// list and delete access keys for this user
		raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.ListAccessKeys(args)
		})
		if err != nil {
			return fmt.Errorf("Error listing access keys for User (%s) when trying to delete: %#v", d.Id(), err)
		}
		akResp, _ := raw.(ram.AccessKeyListResponse)
		if len(akResp.AccessKeys.AccessKey) > 0 {
			for _, v := range akResp.AccessKeys.AccessKey {
				_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
					return ramClient.DeleteAccessKey(ram.UpdateAccessKeyRequest{
						UserAccessKeyId: v.AccessKeyId,
						UserName:        userName,
					})
				})
				if err != nil && !RamEntityNotExist(err) {
					return fmt.Errorf("Error deleting access key %s: %#v", v.AccessKeyId, err)
				}
			}
		}

		// list and delete policies for this user
		raw, err = client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.ListPoliciesForUser(args)
		})
		if err != nil {
			return fmt.Errorf("Error listing policies for User (%s) when trying to delete: %#v", d.Id(), err)
		}
		policyResp, _ := raw.(ram.PolicyListResponse)
		if len(policyResp.Policies.Policy) > 0 {
			for _, v := range policyResp.Policies.Policy {
				_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
					return ramClient.DetachPolicyFromUser(ram.AttachPolicyRequest{
						PolicyRequest: ram.PolicyRequest{
							PolicyName: v.PolicyName,
							PolicyType: ram.Type(v.PolicyType),
						},
						UserName: userName,
					})
				})
				if err != nil && !RamEntityNotExist(err) {
					return fmt.Errorf("Error deleting policy %s: %#v", v.PolicyName, err)
				}
			}
		}

		// list and delete groups for this user
		raw, err = client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.ListGroupsForUser(args)
		})
		if err != nil {
			return fmt.Errorf("Error listing groups for User (%s) when trying to delete: %#v", d.Id(), err)
		}
		groupResp, _ := raw.(ram.GroupListResponse)
		if len(groupResp.Groups.Group) > 0 {
			for _, v := range groupResp.Groups.Group {
				_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
					return ramClient.RemoveUserFromGroup(ram.UserRelateGroupRequest{
						UserName:  userName,
						GroupName: v.GroupName,
					})
				})
				if err != nil && !RamEntityNotExist(err) {
					return fmt.Errorf("Error deleting group %s: %#v", v.GroupName, err)
				}
			}
		}

		// delete login profile for this user
		_, err = client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.DeleteLoginProfile(args)
		})
		if err != nil && !RamEntityNotExist(err) {
			return fmt.Errorf("Error deleting login profile for User (%s): %#v", d.Id(), err)
		}

		// unbind MFA device for this user
		_, err = client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.UnbindMFADevice(args)
		})
		if err != nil && !RamEntityNotExist(err) {
			return fmt.Errorf("Error deleting login profile for User (%s): %#v", d.Id(), err)
		}

	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.DeleteUser(args)
		})
		if err != nil {
			if IsExceptedError(err, DeleteConflictUserAccessKey) || IsExceptedError(err, DeleteConflictUserGroup) ||
				IsExceptedError(err, DeleteConflictUserPolicy) || IsExceptedError(err, DeleteConflictUserLoginProfile) ||
				IsExceptedError(err, DeleteConflictUserMFADevice) {
				return resource.RetryableError(fmt.Errorf("The user can not has any access keys or login profile or attached group or attached policies or attached mfa device while deleting the user.- you can set force with true to force delete the user."))
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting user %s: %#v, you can set force with true to force delete the user.", d.Id(), err))
		}
		return nil
	})
}
