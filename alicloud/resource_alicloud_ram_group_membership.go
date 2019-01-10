package alicloud

import (
	"fmt"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudRamGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamGroupMembershipCreate,
		Read:   resourceAlicloudRamGroupMembershipRead,
		Update: resourceAlicloudRamGroupMembershipUpdate,
		Delete: resourceAlicloudRamGroupMembershipDelete,

		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRamGroupName,
			},
			"user_names": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateRamName,
				},
				Set: schema.HashString,
			},
		},
	}
}

func resourceAlicloudRamGroupMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	group := d.Get("group_name").(string)
	users := expandStringList(d.Get("user_names").(*schema.Set).List())

	err := addUsersToGroup(client, users, group)
	if err != nil {
		return fmt.Errorf("AddUserToGroup got an error: %#v", err)
	}

	d.SetId(group)

	return resourceAlicloudRamGroupMembershipUpdate(d, meta)
}

func resourceAlicloudRamGroupMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	if d.HasChange("user_names") && !d.IsNewResource() {
		d.SetPartial("user_names")
		o, n := d.GetChange("user_names")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}
		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)

		remove := expandStringList(oldSet.Difference(newSet).List())
		add := expandStringList(newSet.Difference(oldSet).List())
		group := d.Id()

		if err := removeUsersFromGroup(client, remove, group); err != nil {
			return fmt.Errorf("removeUsersFromGroup got an error: %#v", err)
		}

		if err := addUsersToGroup(client, add, group); err != nil {
			return fmt.Errorf("addUsersToGroup got an error: %#v", err)
		}
	}

	d.Partial(false)
	return resourceAlicloudRamGroupMembershipRead(d, meta)
}

func resourceAlicloudRamGroupMembershipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := ram.GroupQueryRequest{
		GroupName: d.Id(),
	}

	raw, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
		return ramClient.ListUsersForGroup(args)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("ListUsersForGroup got an error: %#v", err)
	}
	response, _ := raw.(ram.ListUserResponse)
	var users []string
	if len(response.Users.User) > 0 {
		for _, v := range response.Users.User {
			users = append(users, v.UserName)
		}
	}

	d.Set("group_name", args.GroupName)
	if err := d.Set("user_names", users); err != nil {
		return fmt.Errorf("Error setting user list from group membership (%s), error: %#v", args.GroupName, err)
	}

	return nil
}

func resourceAlicloudRamGroupMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	users := expandStringList(d.Get("user_names").(*schema.Set).List())
	group := d.Id()

	if err := removeUsersFromGroup(client, users, group); err != nil {
		return fmt.Errorf("removeUsersFromGroup got an error: %#v", err)
	}

	return nil
}

func addUsersToGroup(client *connectivity.AliyunClient, users []string, group string) error {
	for _, u := range users {
		_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.AddUserToGroup(ram.UserRelateGroupRequest{
				UserName:  u,
				GroupName: group,
			})
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func removeUsersFromGroup(client *connectivity.AliyunClient, users []string, group string) error {
	for _, u := range users {
		_, err := client.WithRamClient(func(ramClient ram.RamClientInterface) (interface{}, error) {
			return ramClient.RemoveUserFromGroup(ram.UserRelateGroupRequest{
				UserName:  u,
				GroupName: group,
			})
		})

		if err != nil && !RamEntityNotExist(err) {
			return err
		}
	}
	return nil
}
