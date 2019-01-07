package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/ram"
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

	args := ram.GroupRequest{
		Group: ram.Group{
			GroupName: d.Get("name").(string),
		},
	}

	raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
		return ramClient.CreateGroup(args)
	})
	if err != nil {
		return fmt.Errorf("CreateGroup got an error: %#v", err)
	}
	response, _ := raw.(ram.GroupResponse)
	d.SetId(response.Group.GroupName)
	return resourceAlicloudRamGroupUpdate(d, meta)
}

func resourceAlicloudRamGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	args := ram.GroupUpdateRequest{
		GroupName: d.Id(),
	}
	attributeUpdate := false

	if d.HasChange("name") && !d.IsNewResource() {
		ov, nv := d.GetChange("name")
		args.GroupName = ov.(string)
		args.NewGroupName = nv.(string)
		d.SetPartial("name")
		attributeUpdate = true
	}

	if d.HasChange("comments") {
		d.SetPartial("comments")
		args.NewComments = d.Get("comments").(string)
		attributeUpdate = true
	}

	if attributeUpdate {
		_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.UpdateGroup(args)
		})
		if err != nil {
			return fmt.Errorf("UpdateGroup got an error: %v", err)
		}
	}

	d.Partial(false)
	return resourceAlicloudRamGroupRead(d, meta)
}

func resourceAlicloudRamGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := ram.GroupQueryRequest{
		GroupName: d.Id(),
	}

	raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
		return ramClient.GetGroup(args)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("GetGroup got an error: %#v", err)
	}
	response, _ := raw.(ram.GroupResponse)
	group := response.Group
	d.Set("name", group.GroupName)
	d.Set("comments", group.Comments)
	return nil
}

func resourceAlicloudRamGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := ram.GroupQueryRequest{
		GroupName: d.Id(),
	}

	if d.Get("force").(bool) {
		// list and delete users which in this group
		raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.ListUsersForGroup(args)
		})
		if err != nil {
			return fmt.Errorf("Error while listing users for group %s: %#v", d.Id(), err)
		}
		listUserResp, _ := raw.(ram.ListUserResponse)
		users := listUserResp.Users.User
		if len(users) > 0 {
			for _, v := range users {
				_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
					return ramClient.RemoveUserFromGroup(ram.UserRelateGroupRequest{
						UserName:  v.UserName,
						GroupName: args.GroupName,
					})
				})
				if err != nil && !RamEntityNotExist(err) {
					return fmt.Errorf("Error while deleting user %s from group %s: %#v", v.UserName, d.Id(), err)
				}
			}
		}

		// list and detach policies which attach this group
		raw, err = client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.ListPoliciesForGroup(args)
		})
		if err != nil {
			return fmt.Errorf("Error while listing policies for group %s: %#v", d.Id(), err)
		}
		listPolicyResp, _ := raw.(ram.PolicyListResponse)
		policies := listPolicyResp.Policies.Policy
		if len(policies) > 0 {
			for _, v := range policies {
				_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
					return ramClient.DetachPolicyFromGroup(ram.AttachPolicyToGroupRequest{
						PolicyRequest: ram.PolicyRequest{
							PolicyType: ram.Type(v.PolicyType),
							PolicyName: v.PolicyName,
						},
						GroupName: args.GroupName,
					})
				})
				if err != nil && !RamEntityNotExist(err) {
					return fmt.Errorf("Error while detaching policy %s from group %s: %#v", v.PolicyName, d.Id(), err)
				}
			}
		}
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.DeleteGroup(args)
		})
		if err != nil {
			if IsExceptedError(err, DeleteConflictGroupUser) || IsExceptedError(err, DeleteConflictGroupPolicy) {
				return resource.RetryableError(fmt.Errorf("The group can not has any user member or any attached policy while deleting the group.- you can set force with true to force delete the group."))
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting group %s: %#v, you can set force with true to force delete the group.", d.Id(), err))
		}
		return nil
	})
}
