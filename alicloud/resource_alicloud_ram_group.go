package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"comments": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudRamGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramSercvice := RamService{client}
	request := map[string]interface{}{
		"RegionId":  client.RegionId,
		"GroupName": d.Get("name").(string),
	}
	if v, ok := d.GetOk("comments"); ok {
		request["Comments"] = v.(string)
	}

	action := "CreateGroup"
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_group", action, AlibabaCloudSdkGoERROR)
	}

	groupName, _ := jsonpath.Get("$.Group.GroupName", response)
	d.SetId(fmt.Sprint(groupName))
	err = ramSercvice.WaitForRamGroup(d.Id(), Normal, DefaultTimeout)
	if err != nil {
		return WrapError(err)
	}
	return resourceAlicloudRamGroupRead(d, meta)
}

func resourceAlicloudRamGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateUpdateGroupRequest()
	request.RegionId = client.RegionId
	request.GroupName = d.Id()

	if d.HasChange("comments") {
		request.NewComments = d.Get("comments").(string)

		var err error
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
				return ramClient.UpdateGroup(request)
			})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
		})

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudRamGroupRead(d, meta)
}

func resourceAlicloudRamGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ramService := RamService{client}

	object, err := ramService.DescribeRamGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	group := object.Group
	d.Set("name", group.GroupName)
	d.Set("comments", group.Comments)
	return nil

}

func resourceAlicloudRamGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := &RamService{client}
	request := ram.CreateListUsersForGroupRequest()
	request.RegionId = client.RegionId
	request.GroupName = d.Id()
	users := []ram.User{}

	if d.Get("force").(bool) {
		// list and delete users which in this group
		var raw interface{}
		var err error

		for {
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
				raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
					return ramClient.ListUsersForGroup(request)
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

			listUserResponse, _ := raw.(*ram.ListUsersForGroupResponse)
			users = append(users, listUserResponse.Users.User...)
			if !listUserResponse.IsTruncated {
				break
			}
			request.Marker = listUserResponse.Marker

		}

		if len(users) > 0 {
			for _, v := range users {
				request := ram.CreateRemoveUserFromGroupRequest()
				request.UserName = v.UserName
				request.GroupName = d.Id()

				wait := incrementalWait(3*time.Second, 3*time.Second)
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

		// list and detach policies which attach this group
		request := ram.CreateListPoliciesForGroupRequest()
		request.RegionId = client.RegionId
		request.GroupName = d.Id()

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
			raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
				return ramClient.ListPoliciesForGroup(request)
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
		listPolicyResponse, _ := raw.(*ram.ListPoliciesForGroupResponse)
		policies := listPolicyResponse.Policies.Policy
		if len(policies) > 0 {
			for _, v := range policies {
				request := ram.CreateDetachPolicyFromGroupRequest()
				request.RegionId = client.RegionId
				request.PolicyType = v.PolicyType
				request.PolicyName = v.PolicyName
				request.GroupName = d.Id()

				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
					raw, err = client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
						return ramClient.DetachPolicyFromGroup(request)
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
	}

	deleteGroupRequest := ram.CreateDeleteGroupRequest()
	deleteGroupRequest.RegionId = client.RegionId
	deleteGroupRequest.GroupName = d.Id()

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.DeleteGroup(deleteGroupRequest)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"DeleteConflict.Group.User", "DeleteConflict.Group.Policy"}) || NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(deleteGroupRequest.GetActionName(), raw, deleteGroupRequest.RpcRequest, deleteGroupRequest)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExist.Group"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), deleteGroupRequest.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(ramService.WaitForRamGroup(d.Id(), Deleted, DefaultTimeout))
}
