// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package dataworks

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceProjectMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectMemberCreate,
		Read:   resourceProjectMemberRead,
		Update: resourceProjectMemberUpdate,
		Delete: resourceProjectMemberDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"roles": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"code": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceProjectMemberCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateProjectMember"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ProjectId"] = d.Get("project_id")
	request["UserId"] = d.Get("user_id")
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("roles"); ok {
		localData, err := jsonpath.Get("$[*].code", v)
		if err != nil {
			return errs.WrapError(err)
		}
		localDataArray := localData.([]interface{})
		localDataJson, err := json.Marshal(localDataArray)
		if err != nil {
			return errs.WrapError(err)
		}
		request["RoleCodes"] = string(localDataJson)
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_data_works_project_member", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["ProjectId"], request["UserId"]))

	return resourceProjectMemberUpdate(d, meta)
}

func resourceProjectMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataWorksServiceV2 := NewService(client)

	objectRaw, err := dataWorksServiceV2.DescribeDataWorksProjectMember(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_data_works_project_member DescribeDataWorksProjectMember Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["ProjectId"] != nil {
		d.Set("project_id", objectRaw["ProjectId"])
	}
	if objectRaw["UserId"] != nil {
		d.Set("user_id", objectRaw["UserId"])
	}

	roles1Raw := objectRaw["Roles"]
	rolesMaps := make([]map[string]interface{}, 0)
	if roles1Raw != nil {
		for _, rolesChild1Raw := range roles1Raw.([]interface{}) {
			rolesMap := make(map[string]interface{})
			rolesChild1Raw := rolesChild1Raw.(map[string]interface{})
			rolesMap["code"] = rolesChild1Raw["Code"]
			rolesMap["name"] = rolesChild1Raw["Name"]
			rolesMap["type"] = rolesChild1Raw["Type"]

			rolesMaps = append(rolesMaps, rolesMap)
		}
	}
	if objectRaw["Roles"] != nil {
		if err := d.Set("roles", rolesMaps); err != nil {
			return err
		}
	}

	return nil
}

func resourceProjectMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}

	if !d.IsNewResource() && d.HasChange("roles") {
		oldEntry, newEntry := d.GetChange("roles")
		removed := oldEntry
		added := newEntry

		if len(removed.([]interface{})) > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "RevokeMemberProjectRoles"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ProjectId"] = parts[0]
			request["UserId"] = parts[1]
			request["RegionId"] = client.RegionId
			localData := removed.([]interface{})
			rolesMaps := make([]interface{}, 0)
			for _, rolesChild1Raw := range localData {
				rolesChild1Raw := rolesChild1Raw.(map[string]interface{})
				rolesMaps = append(rolesMaps, rolesChild1Raw["code"])
			}
			localData = rolesMaps
			localDataArray := localData
			localDataJson, err := json.Marshal(localDataArray)
			if err != nil {
				return errs.WrapError(err)
			}
			request["RoleCodes"] = string(localDataJson)

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
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

		if len(added.([]interface{})) > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "GrantMemberProjectRoles"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ProjectId"] = parts[0]
			request["UserId"] = parts[1]
			request["RegionId"] = client.RegionId
			localData := added.([]interface{})
			rolesMaps := make([]interface{}, 0)
			for _, rolesChild1Raw := range localData {
				rolesChild1Raw := rolesChild1Raw.(map[string]interface{})
				rolesMaps = append(rolesMaps, rolesChild1Raw["code"])
			}
			localData = rolesMaps
			localDataArray := localData
			localDataJson, err := json.Marshal(localDataArray)
			if err != nil {
				return errs.WrapError(err)
			}
			request["RoleCodes"] = string(localDataJson)

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)
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
	return resourceProjectMemberRead(d, meta)
}

func resourceProjectMemberDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteProjectMember"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ProjectId"] = parts[0]
	request["UserId"] = parts[1]
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("dataworks-public", "2024-05-18", action, query, request, true)

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
		// A member that does not exist any more leaves the deletion nothing to do, and 400010
		// means the UserId no longer belongs to the tenant of the workspace.
		if errs.IsExpectedErrors(err, []string{"100002001", "1101080166", "400010"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
