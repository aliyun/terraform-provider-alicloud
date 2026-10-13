// Package alicloud
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudModelstudioWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudModelstudioWorkspaceCreate,
		Read:   resourceAliCloudModelstudioWorkspaceRead,
		Update: resourceAliCloudModelstudioWorkspaceUpdate,
		Delete: resourceAliCloudModelstudioWorkspaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"api_host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_site": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringLenBetween(1, 30),
			},
		},
	}
}

func resourceAliCloudModelstudioWorkspaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/modelstudio/workspaces"
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error

	query["workspaceName"] = StringPointer(d.Get("workspace_name").(string))
	if v, ok := d.GetOk("service_site"); ok {
		query["serviceSite"] = StringPointer(v.(string))
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("modelstudio", "2026-02-10", action, query, nil, nil, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, query)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_modelstudio_workspace", action, AlibabaCloudSdkGoERROR)
	}

	id, idErr := extractModelstudioWorkspaceId(response)
	if idErr != nil {
		return WrapError(idErr)
	}
	d.SetId(id)

	return resourceAliCloudModelstudioWorkspaceRead(d, meta)
}

func resourceAliCloudModelstudioWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	modelstudioService := ModelstudioService{client}

	object, err := modelstudioService.DescribeModelstudioWorkspace(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_modelstudio_workspace DescribeModelstudioWorkspace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("workspace_id", object["workspaceId"])
	d.Set("workspace_name", object["workspaceName"])
	d.Set("region_id", object["region"])
	d.Set("api_host", object["apiHost"])
	d.Set("service_site", object["serviceSite"])
	d.Set("create_time", object["gmtCreate"])

	return nil
}

func resourceAliCloudModelstudioWorkspaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/modelstudio/workspaces/%s", d.Id())
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error

	d.Partial(true)

	if d.HasChange("workspace_name") {
		query["workspaceName"] = StringPointer(d.Get("workspace_name").(string))
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("modelstudio", "2026-02-10", action, query, nil, nil, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, query)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("workspace_name")
	}

	d.Partial(false)
	return resourceAliCloudModelstudioWorkspaceRead(d, meta)
}

func resourceAliCloudModelstudioWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/modelstudio/workspaces/%s", d.Id())
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("modelstudio", "2026-02-10", action, query, nil, nil, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, query)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

// extractModelstudioWorkspaceId extracts the workspace ID from the create-workspace
// response. The gateway response may expose the workspace object under a top-level
// "workspace" key (per the CloudSpec mapping rootMapping) or via the backend nested
// path data.DataV2.data.data. Both paths are tried defensively.
func extractModelstudioWorkspaceId(response map[string]interface{}) (string, error) {
	if response == nil {
		return "", fmt.Errorf("empty response from create-workspace")
	}

	// Primary: top-level "workspace" key (CloudSpec mapping rootMapping $.workspace).
	if raw, ok := response["workspace"]; ok && raw != nil {
		if m, ok := raw.(map[string]interface{}); ok {
			if id := fmt.Sprint(m["workspaceId"]); id != "" && id != "<nil>" {
				return id, nil
			}
		}
	}

	// Fallback: backend nested path data.DataV2.data.data.
	if v, perr := jsonpath.Get("$.data.DataV2.data.data", response); perr == nil && v != nil {
		if m, ok := v.(map[string]interface{}); ok {
			if id := fmt.Sprint(m["workspaceId"]); id != "" && id != "<nil>" {
				return id, nil
			}
		}
	}

	return "", fmt.Errorf("failed to extract workspaceId from create-workspace response: %#v", response)
}

type ModelstudioService struct {
	client *connectivity.AliyunClient
}

// DescribeModelstudioWorkspace fetches a single workspace by ID via the public
// list-workspaces API (get-workspace is not published). The response workspaces
// array is searched for a matching workspaceId; an empty result is treated as
// NotFound.
func (s *ModelstudioService) DescribeModelstudioWorkspace(id string) (map[string]interface{}, error) {
	action := "/modelstudio/workspaces"
	query := make(map[string]*string)
	query["workspaceId"] = StringPointer(id)

	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
		response, err = s.client.RoaGet("modelstudio", "2026-02-10", action, query, nil, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, query)
	if err != nil {
		return nil, WrapError(err)
	}

	workspaces := extractModelstudioWorkspacesList(response)
	for _, ws := range workspaces {
		if m, ok := ws.(map[string]interface{}); ok {
			if fmt.Sprint(m["workspaceId"]) == id {
				return m, nil
			}
		}
	}

	return nil, WrapErrorf(NotFoundErr("ModelstudioWorkspace", id), NotFoundMsg, response)
}

// extractModelstudioWorkspacesList extracts the workspaces array from the
// list-workspaces response, trying the top-level "workspaces" key (CloudSpec
// mapping rootMapping $.workspaces) and the backend nested path defensively.
func extractModelstudioWorkspacesList(response map[string]interface{}) []interface{} {
	if response == nil {
		return nil
	}

	// Primary: top-level "workspaces" key.
	if raw, ok := response["workspaces"]; ok && raw != nil {
		if arr, ok := raw.([]interface{}); ok {
			return arr
		}
	}

	// Fallback: backend nested path data.DataV2.data.data.
	if v, perr := jsonpath.Get("$.data.DataV2.data.data", response); perr == nil && v != nil {
		if arr, ok := v.([]interface{}); ok {
			return arr
		}
		// The backend may return a single object instead of an array when filtered.
		if m, ok := v.(map[string]interface{}); ok {
			return []interface{}{m}
		}
	}

	return nil
}
