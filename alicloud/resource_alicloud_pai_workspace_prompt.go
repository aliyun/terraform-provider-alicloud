// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudPaiWorkspacePrompt() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPaiWorkspacePromptCreate,
		Read:   resourceAliCloudPaiWorkspacePromptRead,
		Update: resourceAliCloudPaiWorkspacePromptUpdate,
		Delete: resourceAliCloudPaiWorkspacePromptDelete,
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
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"PRIVATE", "PUBLIC"}, false),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"framework_content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"framework_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"modify_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prompt_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prompt_name": {
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

func resourceAliCloudPaiWorkspacePromptCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/api/v1/prompts")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("workspace_id"); ok {
		request["WorkspaceId"] = v
	}
	request["PromptName"] = d.Get("prompt_name")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("accessibility"); ok {
		request["Accessibility"] = v
	}
	if v, ok := d.GetOk("framework_type"); ok {
		request["FrameworkType"] = v
	}
	if v, ok := d.GetOk("framework_content"); ok {
		request["FrameworkContent"] = v
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_pai_workspace_prompt", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.PromptId", response)
	d.SetId(fmt.Sprintf("%v:%v", request["WorkspaceId"], id))

	return resourceAliCloudPaiWorkspacePromptRead(d, meta)
}

func resourceAliCloudPaiWorkspacePromptRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	paiWorkspaceServiceV2 := PaiWorkspaceServiceV2{client}

	objectRaw, err := paiWorkspaceServiceV2.DescribePaiWorkspacePrompt(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pai_workspace_prompt DescribePaiWorkspacePrompt Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("accessibility", objectRaw["Accessibility"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("framework_content", objectRaw["FrameworkContent"])
	d.Set("framework_type", objectRaw["FrameworkType"])
	d.Set("modify_time", objectRaw["ModifyTime"])
	d.Set("prompt_name", objectRaw["PromptName"])

	parts := strings.Split(d.Id(), ":")
	d.Set("workspace_id", parts[0])
	d.Set("prompt_id", parts[1])

	return nil
}

func resourceAliCloudPaiWorkspacePromptUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	var err error

	parts := strings.Split(d.Id(), ":")
	PromptId := parts[1]
	action := fmt.Sprintf("/api/v1/prompts/%s", PromptId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["WorkspaceId"] = parts[0]

	if d.HasChange("framework_type") {
		update = true
	}
	if v, ok := d.GetOk("framework_type"); ok || d.HasChange("framework_type") {
		request["FrameworkType"] = v
	}
	if d.HasChange("framework_content") {
		update = true
	}
	if v, ok := d.GetOk("framework_content"); ok || d.HasChange("framework_content") {
		request["FrameworkContent"] = v
	}
	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["Description"] = v
	}
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

	return resourceAliCloudPaiWorkspacePromptRead(d, meta)
}

func resourceAliCloudPaiWorkspacePromptDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	PromptId := parts[1]
	action := fmt.Sprintf("/api/v1/prompts/%s", PromptId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["WorkspaceId"] = StringPointer(parts[0])

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

	return nil
}
