// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package ram

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: ResourceGroupCreate,
		Read:   ResourceGroupRead,
		Update: ResourceGroupUpdate,
		Delete: ResourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"comments": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "Field `name` has been deprecated from provider version 1.245.0. New field `group_name` instead.",
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func ResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("group_name"); ok {
		request["GroupName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["GroupName"] = v
	}

	if v, ok := d.GetOk("comments"); ok {
		request["Comments"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_ram_group", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Group.GroupName", response)
	d.SetId(fmt.Sprint(id))

	return ResourceGroupRead(d, meta)
}

func ResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramServiceV2 := RamServiceV2{client}

	objectRaw, err := ramServiceV2.DescribeRamGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ram_group DescribeRamGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("comments", objectRaw["Comments"])
	d.Set("create_time", objectRaw["CreateDate"])
	d.Set("group_name", objectRaw["GroupName"])
	d.Set("name", objectRaw["GroupName"])

	return nil
}

func ResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["GroupName"] = d.Id()

	if d.HasChange("comments") {
		update = true
	}
	if v, ok := d.GetOk("comments"); ok {
		request["NewComments"] = v
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)
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

	return ResourceGroupRead(d, meta)
}

func ResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error

	if d.Get("force").(bool) {
		// list and delete users which in this group
		action := "ListUsersForGroup"
		userNames := make([]string, 0)
		listUsersForGroupReq := map[string]interface{}{
			"GroupName": d.Id(),
		}

		for {
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
				response, err = client.RpcPost("Ram", "2015-05-01", action, query, listUsersForGroupReq, true)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, listUsersForGroupReq)

			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

			resp, err := jsonpath.Get("$.Users.User", response)
			if err != nil {
				return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Users.User", response)
			}

			result, _ := resp.([]interface{})
			for _, v := range result {
				item := v.(map[string]interface{})
				userNames = append(userNames, fmt.Sprint(item["UserName"]))
			}

			if !response["IsTruncated"].(bool) {
				break
			}

			listUsersForGroupReq["Marker"] = response["Marker"]
		}

		if len(userNames) > 0 {
			for _, userName := range userNames {
				action = "RemoveUserFromGroup"
				removeUserFromGroupReq := map[string]interface{}{
					"GroupName": d.Id(),
					"UserName":  userName,
				}

				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
					response, err = client.RpcPost("Ram", "2015-05-01", action, query, removeUserFromGroupReq, true)

					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, removeUserFromGroupReq)

				if err != nil && !errs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
			}
		}

		// list and detach policies which attach this group
		action = "ListPoliciesForGroup"
		listPoliciesForGroupReq := map[string]interface{}{
			"GroupName": d.Id(),
		}

		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
			response, err = client.RpcPost("Ram", "2015-05-01", action, query, listPoliciesForGroupReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, listPoliciesForGroupReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Policies.Policy", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Policies.Policy", response)
		}

		result, _ := resp.([]interface{})

		if len(result) > 0 {
			for _, v := range result {
				action = "DetachPolicyFromGroup"
				detachPolicyFromGroupReq := map[string]interface{}{
					"GroupName":  d.Id(),
					"PolicyName": v.(map[string]interface{})["PolicyName"],
					"PolicyType": v.(map[string]interface{})["PolicyType"],
				}

				err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
					response, err = client.RpcPost("Ram", "2015-05-01", action, query, detachPolicyFromGroupReq, true)

					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, detachPolicyFromGroupReq)

				if err != nil && !errs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
			}
		}
	}

	action := "DeleteGroup"
	request = make(map[string]interface{})
	request["GroupName"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, query, request, true)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"DeleteConflict.Group.User", "DeleteConflict.Group.Policy"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.Group"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
