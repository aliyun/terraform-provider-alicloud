// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudAgentloopAgentSpace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAgentloopAgentSpaceCreate,
		Read:   resourceAliCloudAgentloopAgentSpaceRead,
		Update: resourceAliCloudAgentloopAgentSpaceUpdate,
		Delete: resourceAliCloudAgentloopAgentSpaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"agent_space": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cms_workspace": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cms_workspace_bind_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mse_namespace_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sls_project": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trajectory_store_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// agentloopAgentSpaceWaitBindTypeSettled waits until the CMS workspace / MSE
// namespace bind types reported by the API match the expected values (an
// empty expectation skips that check). The backend assigns the bind type
// asynchronously: right after creation an auto-created binding may still
// read "UserSelected" for a short window, and right after a cms_workspace
// rebind it may still read the stale type. Delete relies on the final
// bindType to decide the cascade-delete flags, so callers must wait for the
// bind types to settle before returning; otherwise an immediate destroy (or
// a cleanup triggered by a later provisioning failure) could misread the
// stale bind type, skip the cascade deletion and leak the auto-created
// workspace/namespace, or send a delete flag the backend rejects.
func agentloopAgentSpaceWaitBindTypeSettled(client *connectivity.AliyunClient, id string, timeout time.Duration, expectedCmsBindType, expectedMseBindType string) error {
	agentloopServiceV2 := AgentloopServiceV2{client}
	stateConf := BuildStateConf([]string{"settling"}, []string{"settled"}, timeout, 5*time.Second, func() (interface{}, string, error) {
		objectRaw, err := agentloopServiceV2.DescribeAgentloopAgentSpace(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "settling", nil
			}
			return nil, "", err
		}
		if expectedCmsBindType != "" && fmt.Sprint(objectRaw["cmsWorkspaceBindType"]) != expectedCmsBindType {
			return objectRaw, "settling", nil
		}
		if expectedMseBindType != "" {
			mseNamespace, ok := objectRaw["mseNamespace"].(map[string]interface{})
			if !ok || fmt.Sprint(mseNamespace["bindType"]) != expectedMseBindType {
				return objectRaw, "settling", nil
			}
		}
		return objectRaw, "settled", nil
	})
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAliCloudAgentloopAgentSpaceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/agentspace")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["agentSpace"] = d.Get("agent_space")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("cms_workspace"); ok {
		request["cmsWorkspace"] = v
	}
	if v, ok := d.GetOk("mse_namespace_id"); ok {
		request["mseNamespaceId"] = v
	}
	if v, ok := d.GetOkExists("trajectory_store_enabled"); ok {
		request["trajectoryStoreEnabled"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("AgentLoop", "2026-05-20", action, query, nil, body, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"AgentSpaceOperationInProgress"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_agentloop_agent_space", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["agentSpace"]))

	// The expectations are derived from the configuration here, which is
	// reliable because this wait runs before the first Read: afterwards the
	// Optional+Computed mse_namespace_id is backfilled into the state and
	// d.GetOk can no longer distinguish a user-supplied id from the
	// auto-created one.
	expectedCmsBindType := "AutoCreated"
	if _, ok := d.GetOk("cms_workspace"); ok {
		expectedCmsBindType = "UserSelected"
	}
	expectedMseBindType := "AutoCreated"
	if _, ok := d.GetOk("mse_namespace_id"); ok {
		expectedMseBindType = "UserSelected"
	}
	if err := agentloopAgentSpaceWaitBindTypeSettled(client, d.Id(), d.Timeout(schema.TimeoutCreate), expectedCmsBindType, expectedMseBindType); err != nil {
		return err
	}

	return resourceAliCloudAgentloopAgentSpaceRead(d, meta)
}

func resourceAliCloudAgentloopAgentSpaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	agentloopServiceV2 := AgentloopServiceV2{client}

	objectRaw, err := agentloopServiceV2.DescribeAgentloopAgentSpace(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_agentloop_agent_space DescribeAgentloopAgentSpace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["agentSpace"] != nil {
		d.Set("agent_space", objectRaw["agentSpace"])
	}
	if objectRaw["cmsWorkspace"] != nil {
		d.Set("cms_workspace", objectRaw["cmsWorkspace"])
	}
	if objectRaw["cmsWorkspaceBindType"] != nil {
		d.Set("cms_workspace_bind_type", objectRaw["cmsWorkspaceBindType"])
	}
	if objectRaw["createTime"] != nil {
		d.Set("create_time", objectRaw["createTime"])
	}
	if objectRaw["description"] != nil {
		d.Set("description", objectRaw["description"])
	}
	if objectRaw["regionId"] != nil {
		d.Set("region_id", objectRaw["regionId"])
	}
	if objectRaw["slsProject"] != nil {
		d.Set("sls_project", objectRaw["slsProject"])
	}
	if objectRaw["updateTime"] != nil {
		d.Set("update_time", objectRaw["updateTime"])
	}

	mseNamespaceRawObj, _ := jsonpath.Get("$.mseNamespace", objectRaw)
	if mseNamespaceRawObj != nil {
		if mseNamespaceRaw, ok := mseNamespaceRawObj.(map[string]interface{}); ok {
			if mseNamespaceRaw["namespaceId"] != nil {
				d.Set("mse_namespace_id", mseNamespaceRaw["namespaceId"])
			}
		}
	}

	return nil
}

func resourceAliCloudAgentloopAgentSpaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	agentSpace := d.Id()
	action := fmt.Sprintf("/agentspace/%s", agentSpace)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	if d.HasChange("cms_workspace") {
		update = true
		// Only submit cmsWorkspace on an actual change. Read backfills the
		// auto-created workspace name into this Optional+Computed attribute, so
		// an unconditional send would turn every unrelated update (e.g.
		// description) into a rebind request, flipping the bind type to
		// UserSelected and leaking the auto-created workspace on destroy.
		request["cmsWorkspace"] = d.Get("cms_workspace")
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("AgentLoop", "2026-05-20", action, query, nil, body, true)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"AgentSpaceOperationInProgress"}) {
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
		// Rebinding the CMS workspace flips its bind type asynchronously;
		// wait for the settled value for the same reason as in Create. Only
		// the CMS bind type is checked: mse_namespace_id is ForceNew (it
		// cannot change in an update and the Create-time wait already
		// settled it), and its Optional+Computed state value has been
		// backfilled by Read, so d.GetOk can no longer tell whether the id
		// was user-supplied or auto-created.
		if d.HasChange("cms_workspace") {
			expectedCmsBindType := "AutoCreated"
			if v, ok := d.GetOk("cms_workspace"); ok && v.(string) != "" {
				expectedCmsBindType = "UserSelected"
			}
			if err := agentloopAgentSpaceWaitBindTypeSettled(client, d.Id(), d.Timeout(schema.TimeoutUpdate), expectedCmsBindType, ""); err != nil {
				return err
			}
		}
	}

	d.Partial(false)
	return resourceAliCloudAgentloopAgentSpaceRead(d, meta)
}

func resourceAliCloudAgentloopAgentSpaceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	agentSpace := d.Id()
	action := fmt.Sprintf("/agentspace/%s", agentSpace)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	request["agentSpace"] = d.Id()
	// Cascade-delete the resources auto-provisioned together with the
	// AgentSpace; otherwise every AgentSpace leaks an SLS project, a CMS
	// workspace and an MSE namespace, which eventually exhausts the account
	// SLS project quota (150). The SLS project is always auto-created, but
	// the CMS workspace / MSE namespace may be user-selected at creation time
	// and the backend then rejects the corresponding delete flag with
	// InvalidParameter ("only supported for auto-created ..."). Describe the
	// AgentSpace first and only request cascade deletion for resources whose
	// bindType is AutoCreated.
	agentloopService := AgentloopServiceV2{client: client}
	objectRaw, err := agentloopService.DescribeAgentloopAgentSpace(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	query["deleteSlsProject"] = StringPointer("true")
	if fmt.Sprint(objectRaw["cmsWorkspaceBindType"]) == "AutoCreated" {
		query["deleteCmsWorkspace"] = StringPointer("true")
	}
	if mseNamespace, ok := objectRaw["mseNamespace"].(map[string]interface{}); ok && fmt.Sprint(mseNamespace["bindType"]) == "AutoCreated" {
		query["deleteMseNamespace"] = StringPointer("true")
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("AgentLoop", "2026-05-20", action, query, nil, nil, true)

		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"AgentSpaceOperationInProgress"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"AgentSpaceNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
