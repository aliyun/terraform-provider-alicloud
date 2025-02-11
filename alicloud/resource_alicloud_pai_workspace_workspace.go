// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudPaiWorkspaceWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPaiWorkspaceWorkspaceCreate,
		Read:   resourceAliCloudPaiWorkspaceWorkspaceRead,
		Update: resourceAliCloudPaiWorkspaceWorkspaceUpdate,
		Delete: resourceAliCloudPaiWorkspaceWorkspaceDelete,
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
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"env_types": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudPaiWorkspaceWorkspaceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/api/v1/workspaces")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["WorkspaceName"] = d.Get("workspace_name")
	request["Description"] = d.Get("description")
	if v, ok := d.GetOk("display_name"); ok {
		request["DisplayName"] = v
	}
	if v, ok := d.GetOk("env_types"); ok {
		envTypesMaps := v.([]interface{})
		request["EnvTypes"] = envTypesMaps
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("AIWorkSpace", "2021-02-04", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pai_workspace_workspace", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["WorkspaceId"]))

	paiWorkspaceServiceV2 := PaiWorkspaceServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"ENABLED"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, paiWorkspaceServiceV2.PaiWorkspaceWorkspaceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudPaiWorkspaceWorkspaceRead(d, meta)
}

func resourceAliCloudPaiWorkspaceWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	paiWorkspaceServiceV2 := PaiWorkspaceServiceV2{client}

	objectRaw, err := paiWorkspaceServiceV2.DescribePaiWorkspaceWorkspace(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pai_workspace_workspace DescribePaiWorkspaceWorkspace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
	if objectRaw["Status"] != nil {
		d.Set("status", objectRaw["Status"])
	}
	if objectRaw["WorkspaceName"] != nil {
		d.Set("workspace_name", objectRaw["WorkspaceName"])
	}

	envTypes1Raw := make([]interface{}, 0)
	if objectRaw["EnvTypes"] != nil {
		envTypes1Raw = objectRaw["EnvTypes"].([]interface{})
	}

	d.Set("env_types", envTypes1Raw)

	return nil
}

func resourceAliCloudPaiWorkspaceWorkspaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	WorkspaceId := d.Id()
	action := fmt.Sprintf("/api/v1/workspaces/%s", WorkspaceId)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["WorkspaceId"] = d.Id()

	if d.HasChange("display_name") {
		update = true
	}
	if v, ok := d.GetOk("display_name"); ok || d.HasChange("display_name") {
		request["DisplayName"] = v
	}
	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("AIWorkSpace", "2021-02-04", action, query, nil, body, true)
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

	return resourceAliCloudPaiWorkspaceWorkspaceRead(d, meta)
}

func resourceAliCloudPaiWorkspaceWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	WorkspaceId := d.Id()
	action := fmt.Sprintf("/api/v1/workspaces/%s", WorkspaceId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	request["WorkspaceId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("AIWorkSpace", "2021-02-04", action, query, nil, nil, true)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	paiWorkspaceServiceV2 := PaiWorkspaceServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, paiWorkspaceServiceV2.PaiWorkspaceWorkspaceStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
