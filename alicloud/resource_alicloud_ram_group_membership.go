package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
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
		return WrapError(err)
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
			return WrapError(err)
		}

		if err := addUsersToGroup(client, add, group); err != nil {
			return WrapError(err)
		}
	}

	d.Partial(false)
	return resourceAlicloudRamGroupMembershipRead(d, meta)
}

func resourceAlicloudRamGroupMembershipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateListUsersForGroupRequest()
	request.GroupName = d.Id()

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListUsersForGroup(request)
	})
	if err != nil {
		if RamEntityNotExist(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*ram.ListUsersForGroupResponse)
	var users []string
	if len(response.Users.User) > 0 {
		for _, v := range response.Users.User {
			users = append(users, v.UserName)
		}
	}

	d.Set("group_name", request.GroupName)
	if err := d.Set("user_names", users); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAlicloudRamGroupMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	users := expandStringList(d.Get("user_names").(*schema.Set).List())
	group := d.Id()

	if err := removeUsersFromGroup(client, users, group); err != nil {
		return WrapError(err)
	}

	return nil
}

func addUsersToGroup(client *connectivity.AliyunClient, users []string, group string) error {
	for _, u := range users {
		request := ram.CreateAddUserToGroupRequest()
		request.UserName = u
		request.GroupName = group
		_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.AddUserToGroup(request)
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, u, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}

func removeUsersFromGroup(client *connectivity.AliyunClient, users []string, group string) error {
	for _, u := range users {
		request := ram.CreateRemoveUserFromGroupRequest()
		request.UserName = u
		request.GroupName = group
		_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.RemoveUserFromGroup(request)
		})

		if err != nil && !RamEntityNotExist(err) {
			return WrapErrorf(err, DefaultErrorMsg, u, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}
