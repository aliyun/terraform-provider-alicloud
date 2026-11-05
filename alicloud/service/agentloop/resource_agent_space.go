package agentloop

import (
	"context"
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

func ResourceAgentSpace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAgentSpaceCreate,
		Read:   resourceAgentSpaceRead,
		Update: resourceAgentSpaceUpdate,
		Delete: resourceAgentSpaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: agentloopImportStatePassthrough(1),
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
	agentloopService := NewService(client)
	stateConf := helper.BuildStateConf([]string{"settling"}, []string{"settled"}, timeout, 5*time.Second, func() (interface{}, string, error) {
		objectRaw, err := agentloopService.DescribeAgentSpace(id)
		if err != nil {
			if errs.NotFoundError(err) {
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
	if _, err := stateConf.WaitForStateContext(context.Background()); err != nil {
		return errs.WrapError(err)
	}
	return nil
}

func resourceAgentSpaceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "/agentspace"
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
	request["trajectoryStoreEnabled"] = d.Get("trajectory_store_enabled")
	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("AgentLoop", "2026-05-20", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"AgentSpaceOperationInProgress"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_agentloop_agent_space", action, errs.AlibabaCloudSdkGoERROR)
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

	return resourceAgentSpaceRead(d, meta)
}

func resourceAgentSpaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	agentloopService := NewService(client)

	objectRaw, err := agentloopService.DescribeAgentSpace(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_agentloop_agent_space DescribeAgentSpace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["agentSpace"] != nil {
		if err := d.Set("agent_space", objectRaw["agentSpace"]); err != nil {
			return err
		}
	}
	if objectRaw["cmsWorkspace"] != nil {
		if err := d.Set("cms_workspace", objectRaw["cmsWorkspace"]); err != nil {
			return err
		}
	}
	if objectRaw["cmsWorkspaceBindType"] != nil {
		if err := d.Set("cms_workspace_bind_type", objectRaw["cmsWorkspaceBindType"]); err != nil {
			return err
		}
	}
	if objectRaw["createTime"] != nil {
		if err := d.Set("create_time", objectRaw["createTime"]); err != nil {
			return err
		}
	}
	if objectRaw["description"] != nil {
		if err := d.Set("description", objectRaw["description"]); err != nil {
			return err
		}
	}
	if objectRaw["regionId"] != nil {
		if err := d.Set("region_id", objectRaw["regionId"]); err != nil {
			return err
		}
	}
	if objectRaw["slsProject"] != nil {
		if err := d.Set("sls_project", objectRaw["slsProject"]); err != nil {
			return err
		}
	}
	if objectRaw["updateTime"] != nil {
		if err := d.Set("update_time", objectRaw["updateTime"]); err != nil {
			return err
		}
	}

	mseNamespaceRawObj, _ := jsonpath.Get("$.mseNamespace", objectRaw)
	if mseNamespaceRawObj != nil {
		if mseNamespaceRaw, ok := mseNamespaceRawObj.(map[string]interface{}); ok {
			if mseNamespaceRaw["namespaceId"] != nil {
				if err := d.Set("mse_namespace_id", mseNamespaceRaw["namespaceId"]); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func resourceAgentSpaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPut("AgentLoop", "2026-05-20", action, query, nil, body, true)
			if err != nil {
				if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"AgentSpaceOperationInProgress"}) {
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

	return resourceAgentSpaceRead(d, meta)
}

func resourceAgentSpaceDelete(d *schema.ResourceData, meta interface{}) error {

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
	agentloopService := NewService(client)
	objectRaw, err := agentloopService.DescribeAgentSpace(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	query["deleteSlsProject"] = helper.StringPointer("true")
	if fmt.Sprint(objectRaw["cmsWorkspaceBindType"]) == "AutoCreated" {
		query["deleteCmsWorkspace"] = helper.StringPointer("true")
	}
	if mseNamespace, ok := objectRaw["mseNamespace"].(map[string]interface{}); ok && fmt.Sprint(mseNamespace["bindType"]) == "AutoCreated" {
		query["deleteMseNamespace"] = helper.StringPointer("true")
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("AgentLoop", "2026-05-20", action, query, nil, nil, true)

		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"AgentSpaceOperationInProgress"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"AgentSpaceNotExist"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
