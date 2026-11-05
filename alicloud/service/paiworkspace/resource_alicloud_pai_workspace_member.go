// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package paiworkspace

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourcePaiWorkspaceMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPaiWorkspaceMemberCreate,
		Read:   resourceAliCloudPaiWorkspaceMemberRead,
		Update: resourceAliCloudPaiWorkspaceMemberUpdate,
		Delete: resourceAliCloudPaiWorkspaceMemberDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"member_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"roles": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudPaiWorkspaceMemberCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	WorkspaceId := d.Get("workspace_id")
	action := fmt.Sprintf("/api/v1/workspaces/%s/members", WorkspaceId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	objectDataLocalMap := make(map[string]interface{})

	if v, ok := d.GetOk("user_id"); ok {
		objectDataLocalMap["UserId"] = v
	}

	if v, ok := d.GetOk("roles"); ok {
		roles1, _ := jsonpath.Get("$", v)
		if roles1 != nil && roles1 != "" {
			objectDataLocalMap["Roles"] = roles1.(*schema.Set).List()
		}
	}

	MembersMap := make([]interface{}, 0)
	MembersMap = append(MembersMap, objectDataLocalMap)
	request["Members"] = MembersMap

	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("AIWorkSpace", "2021-02-04", action, query, nil, body, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_pai_workspace_member", action, errs.AlibabaCloudSdkGoERROR)
	}

	MembersMemberIdVar, _ := jsonpath.Get("$.Members[0].MemberId", response)
	d.SetId(fmt.Sprintf("%v:%v", WorkspaceId, MembersMemberIdVar))

	return resourceAliCloudPaiWorkspaceMemberRead(d, meta)
}

func resourceAliCloudPaiWorkspaceMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	paiWorkspaceServiceV2 := NewService(client)

	objectRaw, err := paiWorkspaceServiceV2.DescribePaiWorkspaceMember(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pai_workspace_member DescribePaiWorkspaceMember Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["GmtCreateTime"])
	d.Set("user_id", objectRaw["UserId"])
	d.Set("member_id", objectRaw["MemberId"])

	rolesRaw := make([]interface{}, 0)
	if objectRaw["Roles"] != nil {
		rolesRaw = objectRaw["Roles"].([]interface{})
	}

	d.Set("roles", rolesRaw)

	parts := strings.Split(d.Id(), ":")
	d.Set("workspace_id", parts[0])

	return nil
}

func resourceAliCloudPaiWorkspaceMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	var err error

	if d.HasChange("roles") {
		oldEntry, newEntry := d.GetChange("roles")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			roles := removed.List()

			for _, item := range roles {
				parts := strings.Split(d.Id(), ":")
				WorkspaceId := parts[0]
				MemberId := parts[1]
				action := fmt.Sprintf("/api/v1/workspaces/%s/members/%s/roles/%s", WorkspaceId, MemberId, item)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				request["WorkspaceId"] = WorkspaceId
				request["MemberId"] = MemberId

				if v, ok := item.(string); ok {
					request["RoleName"] = v
				}
				body = request
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RoaDelete("AIWorkSpace", "2021-02-04", action, query, nil, nil, true)
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
		}

		if added.Len() > 0 {
			roles := added.List()

			for _, item := range roles {
				parts := strings.Split(d.Id(), ":")
				WorkspaceId := parts[0]
				MemberId := parts[1]
				action := fmt.Sprintf("/api/v1/workspaces/%s/members/%s/roles/%s", WorkspaceId, MemberId, item)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				request["WorkspaceId"] = WorkspaceId
				request["MemberId"] = MemberId

				if v, ok := item.(string); ok {
					request["RoleName"] = v
				}
				body = request
				wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RoaPost("AIWorkSpace", "2021-02-04", action, query, nil, body, true)
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
		}

	}
	return resourceAliCloudPaiWorkspaceMemberRead(d, meta)
}

func resourceAliCloudPaiWorkspaceMemberDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	WorkspaceId := parts[0]
	action := fmt.Sprintf("/api/v1/workspaces/%s/members", WorkspaceId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["MemberIds"] = helper.StringPointer(parts[1])

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("AIWorkSpace", "2021-02-04", action, query, nil, nil, true)

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
		if errs.IsExpectedErrors(err, []string{"100600017"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
