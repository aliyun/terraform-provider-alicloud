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

func resourceAliCloudAgentLoopAgentSpace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAgentLoopAgentSpaceCreate,
		Read:   resourceAliCloudAgentLoopAgentSpaceRead,
		Update: resourceAliCloudAgentLoopAgentSpaceUpdate,
		Delete: resourceAliCloudAgentLoopAgentSpaceDelete,
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
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mse_workspace": {
				Type:     schema.TypeString,
				Optional: true,
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
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delete_cms_workspace": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"delete_mse_namespace": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"delete_sls_project": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudAgentLoopAgentSpaceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/agentspace")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("agent_space"); ok {
		request["agentSpace"] = v
	}
	query["RegionId"] = StringPointer(client.RegionId)

	if v, ok := d.GetOk("cms_workspace"); ok {
		request["cmsWorkspace"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("mse_workspace"); ok {
		request["mseWorkspace"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("AgentLoop", "2026-05-20", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_agentloop_agent_space", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["agentSpace"]))

	return resourceAliCloudAgentLoopAgentSpaceRead(d, meta)
}

func resourceAliCloudAgentLoopAgentSpaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	agentLoopServiceV2 := AgentLoopServiceV2{client}

	objectRaw, err := agentLoopServiceV2.DescribeAgentLoopAgentSpace(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_agentloop_agent_space DescribeAgentLoopAgentSpace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cms_workspace", objectRaw["cmsWorkspace"])
	d.Set("create_time", objectRaw["createTime"])
	d.Set("description", objectRaw["description"])
	d.Set("mse_workspace", objectRaw["mseWorkspace"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("sls_project", objectRaw["slsProject"])
	d.Set("update_time", objectRaw["updateTime"])
	d.Set("agent_space", objectRaw["agentSpace"])

	return nil
}

func resourceAliCloudAgentLoopAgentSpaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	agentSpace := d.Id()
	action := fmt.Sprintf("/agentspace/%s", agentSpace)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("AgentLoop", "2026-05-20", action, query, nil, body, true)
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

	return resourceAliCloudAgentLoopAgentSpaceRead(d, meta)
}

func resourceAliCloudAgentLoopAgentSpaceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	agentSpace := d.Id()
	action := fmt.Sprintf("/agentspace/%s", agentSpace)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("delete_cms_workspace"); ok && v.(bool) {
		query["deleteCmsWorkspace"] = StringPointer("true")
	}
	if v, ok := d.GetOk("delete_mse_namespace"); ok && v.(bool) {
		query["deleteMseNamespace"] = StringPointer("true")
	}
	if v, ok := d.GetOk("delete_sls_project"); ok && v.(bool) {
		query["deleteSlsProject"] = StringPointer("true")
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("AgentLoop", "2026-05-20", action, query, nil, nil, true)
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
