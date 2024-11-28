// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudPaiWorkspaceCodeSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPaiWorkspaceCodeSourceCreate,
		Read:   resourceAliCloudPaiWorkspaceCodeSourceRead,
		Update: resourceAliCloudPaiWorkspaceCodeSourceUpdate,
		Delete: resourceAliCloudPaiWorkspaceCodeSourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accessibility": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"PRIVATE", "PUBLIC"}, false),
			},
			"code_branch": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"code_commit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"code_repo": {
				Type:     schema.TypeString,
				Required: true,
			},
			"code_repo_access_token": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"code_repo_user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mount_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudPaiWorkspaceCodeSourceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/api/v1/codesources")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	conn, err := client.NewPaiworkspaceClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["WorkspaceId"] = d.Get("workspace_id")
	request["DisplayName"] = d.Get("display_name")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["CodeRepo"] = d.Get("code_repo")
	if v, ok := d.GetOk("code_repo_user_name"); ok {
		request["CodeRepoUserName"] = v
	}
	if v, ok := d.GetOk("code_repo_access_token"); ok {
		request["CodeRepoAccessToken"] = v
	}
	if v, ok := d.GetOk("code_branch"); ok {
		request["CodeBranch"] = v
	}
	request["MountPath"] = d.Get("mount_path")
	request["Accessibility"] = d.Get("accessibility")
	if v, ok := d.GetOk("code_commit"); ok {
		request["CodeCommit"] = v
	}
	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2021-02-04"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pai_workspace_code_source", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.body.CodeSourceId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudPaiWorkspaceCodeSourceRead(d, meta)
}

func resourceAliCloudPaiWorkspaceCodeSourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	paiWorkspaceServiceV2 := PaiWorkspaceServiceV2{client}

	objectRaw, err := paiWorkspaceServiceV2.DescribePaiWorkspaceCodeSource(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pai_workspace_code_source DescribePaiWorkspaceCodeSource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Accessibility"] != nil {
		d.Set("accessibility", objectRaw["Accessibility"])
	}
	if objectRaw["CodeBranch"] != nil {
		d.Set("code_branch", objectRaw["CodeBranch"])
	}
	if objectRaw["CodeCommit"] != nil {
		d.Set("code_commit", objectRaw["CodeCommit"])
	}
	if objectRaw["CodeRepo"] != nil {
		d.Set("code_repo", objectRaw["CodeRepo"])
	}
	if objectRaw["CodeRepoAccessToken"] != nil {
		d.Set("code_repo_access_token", objectRaw["CodeRepoAccessToken"])
	}
	if objectRaw["CodeRepoUserName"] != nil {
		d.Set("code_repo_user_name", objectRaw["CodeRepoUserName"])
	}
	if objectRaw["GmtCreateTime"] != nil {
		d.Set("create_time", objectRaw["GmtCreateTime"])
	}
	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["DisplayName"] != nil {
		d.Set("display_name", objectRaw["DisplayName"])
	}
	if objectRaw["MountPath"] != nil {
		d.Set("mount_path", objectRaw["MountPath"])
	}
	if objectRaw["WorkspaceId"] != nil {
		d.Set("workspace_id", objectRaw["WorkspaceId"])
	}

	return nil
}

func resourceAliCloudPaiWorkspaceCodeSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	if d.HasChange("accessibility") {
		paiWorkspaceServiceV2 := PaiWorkspaceServiceV2{client}
		object, err := paiWorkspaceServiceV2.DescribePaiWorkspaceCodeSource(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("accessibility").(string)
		if object["Accessibility"].(string) != target {
			if target == "PUBLIC" {
				CodeSourceId := d.Id()
				action := fmt.Sprintf("/api/v1/codesources/%s/publish", CodeSourceId)
				conn, err := client.NewPaiworkspaceClient()
				if err != nil {
					return WrapError(err)
				}
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				request["CodeSourceId"] = d.Id()

				body = request
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer("2021-02-04"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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
	}

	CodeSourceId := d.Id()
	action := fmt.Sprintf("/api/v1/codesources/%s", CodeSourceId)
	conn, err := client.NewPaiworkspaceClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["CodeSourceId"] = d.Id()

	if d.HasChange("display_name") {
		update = true
	}
	request["DisplayName"] = d.Get("display_name")
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["Description"] = v
	}
	if d.HasChange("code_repo") {
		update = true
	}
	request["CodeRepo"] = d.Get("code_repo")
	if d.HasChange("code_branch") {
		update = true
	}
	if v, ok := d.GetOk("code_branch"); ok || d.HasChange("code_branch") {
		request["CodeBranch"] = v
	}
	if d.HasChange("code_commit") {
		update = true
	}
	if v, ok := d.GetOk("code_commit"); ok || d.HasChange("code_commit") {
		request["CodeCommit"] = v
	}
	if d.HasChange("code_repo_user_name") {
		update = true
	}
	if v, ok := d.GetOk("code_repo_user_name"); ok || d.HasChange("code_repo_user_name") {
		request["CodeRepoUserName"] = v
	}
	if d.HasChange("code_repo_access_token") {
		update = true
	}
	if v, ok := d.GetOk("code_repo_access_token"); ok || d.HasChange("code_repo_access_token") {
		request["CodeRepoAccessToken"] = v
	}
	if d.HasChange("mount_path") {
		update = true
	}
	request["MountPath"] = d.Get("mount_path")
	body = request
	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2021-02-04"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), query, nil, body, &runtime)
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

	return resourceAliCloudPaiWorkspaceCodeSourceRead(d, meta)
}

func resourceAliCloudPaiWorkspaceCodeSourceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	CodeSourceId := d.Id()
	action := fmt.Sprintf("/api/v1/codesources/%s", CodeSourceId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	conn, err := client.NewPaiworkspaceClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["CodeSourceId"] = d.Id()

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2021-02-04"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, nil, &runtime)

		if err != nil {
			if IsExpectedErrors(err, []string{"201400004"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"201400002"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
