// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDataWorksProjectMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDataWorksProjectMemberCreate,
		Read:   resourceAliCloudDataWorksProjectMemberRead,
		Update: resourceAliCloudDataWorksProjectMemberUpdate,
		Delete: resourceAliCloudDataWorksProjectMemberDelete,
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

func resourceAliCloudDataWorksProjectMemberCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateProjectMember"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ProjectId"] = d.Get("project_id")
	request["UserId"] = d.Get("user_id")
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("roles"); ok {
		localData, err := jsonpath.Get("$[*].code", v)
		if err != nil {
			return WrapError(err)
		}
		localDataArray := localData.([]interface{})
		localDataJson, err := json.Marshal(localDataArray)
		if err != nil {
			return WrapError(err)
		}
		request["RoleCodes"] = string(localDataJson)
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2024-05-18"), StringPointer("AK"), query, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_data_works_project_member", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["ProjectId"], request["UserId"]))

	return resourceAliCloudDataWorksProjectMemberUpdate(d, meta)
}

func resourceAliCloudDataWorksProjectMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataWorksServiceV2 := DataWorksServiceV2{client}

	objectRaw, err := dataWorksServiceV2.DescribeDataWorksProjectMember(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_data_works_project_member DescribeDataWorksProjectMember Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

func resourceAliCloudDataWorksProjectMemberUpdate(d *schema.ResourceData, meta interface{}) error {
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
			conn, err := client.NewDataworkspublicClient()
			if err != nil {
				return WrapError(err)
			}
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
				return WrapError(err)
			}
			request["RoleCodes"] = string(localDataJson)

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2024-05-18"), StringPointer("AK"), query, request, &runtime)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}

		if len(added.([]interface{})) > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "GrantMemberProjectRoles"
			conn, err := client.NewDataworkspublicClient()
			if err != nil {
				return WrapError(err)
			}
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
				return WrapError(err)
			}
			request["RoleCodes"] = string(localDataJson)

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2024-05-18"), StringPointer("AK"), query, request, &runtime)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, request)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}
	}
	return resourceAliCloudDataWorksProjectMemberRead(d, meta)
}

func resourceAliCloudDataWorksProjectMemberDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteProjectMember"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["ProjectId"] = parts[0]
	request["UserId"] = parts[1]
	request["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2024-05-18"), StringPointer("AK"), query, request, &runtime)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"100002001", "1101080166"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
