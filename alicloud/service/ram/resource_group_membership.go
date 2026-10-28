package ram

import (
	"time"

	ramsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceGroupMembership() *schema.Resource {
	return &schema.Resource{
		Create: ResourceGroupMembershipCreate,
		Read:   ResourceGroupMembershipRead,
		Update: ResourceGroupMembershipUpdate,
		Delete: ResourceGroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		DeprecationMessage: "This resource has been deprecated since v1.267.0 and will be removed in the future. Please use 'alicloud_ram_user_group_attachment' instead.",

		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_names": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
		},
	}
}

func ResourceGroupMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	group := d.Get("group_name").(string)
	users := convert.ExpandStringList(d.Get("user_names").(*schema.Set).List())

	err := addUsersToGroup(d, client, users, group, false)
	if err != nil {
		return errs.WrapError(err)
	}

	d.SetId(group)

	return ResourceGroupMembershipRead(d, meta)
}

func ResourceGroupMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	if d.HasChange("user_names") {
		o, n := d.GetChange("user_names")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}
		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)

		remove := convert.ExpandStringList(oldSet.Difference(newSet).List())
		add := convert.ExpandStringList(newSet.Difference(oldSet).List())
		group := d.Id()

		if err := removeUsersFromGroup(d, client, remove, group, true); err != nil {
			return errs.WrapError(err)
		}

		if err := addUsersToGroup(d, client, add, group, true); err != nil {
			return errs.WrapError(err)
		}
	}

	d.Partial(false)
	return ResourceGroupMembershipRead(d, meta)
}

func ResourceGroupMembershipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}

	object, err := ramService.DescribeRamGroupMembership(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	var users []string
	for _, v := range object.Users.User {
		users = append(users, v.UserName)
	}
	d.Set("group_name", d.Id())
	if err := d.Set("user_names", users); err != nil {
		return errs.WrapError(err)
	}
	return nil
}

func ResourceGroupMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	users := convert.ExpandStringList(d.Get("user_names").(*schema.Set).List())
	group := d.Id()

	if err := removeUsersFromGroup(d, client, users, group, false); err != nil {
		return errs.WrapError(err)
	}
	return errs.WrapError(ramService.WaitForRamGroupMembership(d.Id(), Deleted, defaultTimeout))
}

func addUsersToGroup(d *schema.ResourceData, client *connectivity.AliyunClient, users []string, group string, isUpdate bool) error {
	for _, u := range users {
		request := ramsdk.CreateAddUserToGroupRequest()
		request.RegionId = client.RegionId
		request.UserName = u
		request.GroupName = group

		timeout := schema.TimeoutCreate
		if isUpdate {
			timeout = schema.TimeoutUpdate
		}

		var err error
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(timeout)), func() *retry.RetryError {
			raw, err := client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
				return ramClient.AddUserToGroup(request)
			})
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, u, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}

	}
	return nil
}

func removeUsersFromGroup(d *schema.ResourceData, client *connectivity.AliyunClient, users []string, group string, isUpdate bool) error {
	for _, u := range users {
		request := ramsdk.CreateRemoveUserFromGroupRequest()
		request.RegionId = client.RegionId
		request.UserName = u
		request.GroupName = group

		timeout := schema.TimeoutDelete
		if isUpdate {
			timeout = schema.TimeoutUpdate
		}

		var err error
		wait := helper.IncrementalWait(3*time.Second, 20*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(timeout)), func() *retry.RetryError {
			raw, err := client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
				return ramClient.RemoveUserFromGroup(request)
			})
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})

		if err != nil && !errs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, u, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		}

	}
	return nil
}
