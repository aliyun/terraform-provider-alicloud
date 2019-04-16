package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudRamGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamGroupCreate,
		Read:   resourceAlicloudRamGroupRead,
		Update: resourceAlicloudRamGroupUpdate,
		Delete: resourceAlicloudRamGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamGroupName,
			},
			"comments": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateComment,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAlicloudRamGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateCreateGroupRequest()
	request.GroupName = d.Get("name").(string)

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.CreateGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "ram_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*ram.CreateGroupResponse)
	d.SetId(response.Group.GroupName)
	return resourceAlicloudRamGroupUpdate(d, meta)
}

func resourceAlicloudRamGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateUpdateGroupRequest()
	request.GroupName = d.Id()

	if d.HasChange("comments") {
		request.NewComments = d.Get("comments").(string)
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.UpdateGroup(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
	}

	return resourceAlicloudRamGroupRead(d, meta)
}

func resourceAlicloudRamGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateGetGroupRequest()
	request.GroupName = d.Id()

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetGroup(request)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*ram.GetGroupResponse)
	group := response.Group
	d.Set("name", group.GroupName)
	d.Set("comments", group.Comments)
	return nil
}

func resourceAlicloudRamGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateListUsersForGroupRequest()
	request.GroupName = d.Id()

	if d.Get("force").(bool) {
		// list and delete users which in this group
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListUsersForGroup(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		listUserResp, _ := raw.(*ram.ListUsersForGroupResponse)
		users := listUserResp.Users.User
		if len(users) > 0 {
			for _, v := range users {
				request := ram.CreateRemoveUserFromGroupRequest()
				request.UserName = v.UserName
				request.GroupName = d.Id()
				raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.RemoveUserFromGroup(request)
				})
				if err != nil && !RamEntityNotExist(err) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw)
			}
		}

		// list and detach policies which attach this group
		request := ram.CreateListPoliciesForGroupRequest()
		request.GroupName = d.Id()
		raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForGroup(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		listPolicyResp, _ := raw.(*ram.ListPoliciesForGroupResponse)
		policies := listPolicyResp.Policies.Policy
		if len(policies) > 0 {
			for _, v := range policies {
				request := ram.CreateDetachPolicyFromGroupRequest()
				request.PolicyType = v.PolicyType
				request.PolicyName = v.PolicyName
				request.GroupName = d.Id()
				raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.DetachPolicyFromGroup(request)
				})
				if err != nil && !RamEntityNotExist(err) {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				addDebug(request.GetActionName(), raw)
			}
		}
	}

	request1 := ram.CreateDeleteGroupRequest()
	request1.GroupName = d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.DeleteGroup(request1)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{DeleteConflictGroupUser, DeleteConflictGroupPolicy}) {
				return resource.RetryableError(WrapError(Error("The group can not has any user member or any attached policy while deleting the group.- you can set force with true to force delete the group.")))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request1.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		addDebug(request1.GetActionName(), raw)
		return nil
	})
}
