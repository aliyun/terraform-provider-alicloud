// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAgentloopOptimizeTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAgentloopOptimizeTaskCreate,
		Read:   resourceAliCloudAgentloopOptimizeTaskRead,
		Delete: resourceAliCloudAgentloopOptimizeTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"agent_space": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"optimize_task_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Enabled", "Disabled"}, false),
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudAgentloopOptimizeTaskCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	agentSpace := d.Get("agent_space")
	action := fmt.Sprintf("/agentspace/%s/optimizetasks", agentSpace)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	// CreateOptimizeTask requires a clientToken for idempotency; generate one per Create call.
	clientToken := resource.UniqueId()
	query["clientToken"] = &clientToken
	if v, ok := d.GetOk("optimize_task_name"); ok {
		request["optimizeTaskName"] = v
	}

	if v, ok := d.GetOk("config"); ok {
		request["config"] = v
	}
	request["status"] = d.Get("status")
	request["type"] = d.Get("type")
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("AgentLoop", "2026-05-20", action, query, nil, body, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict."}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_agentloop_optimize_task", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", agentSpace, request["optimizeTaskName"]))

	return resourceAliCloudAgentloopOptimizeTaskRead(d, meta)
}

func resourceAliCloudAgentloopOptimizeTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	agentloopServiceV2 := AgentloopServiceV2{client}

	objectRaw, err := agentloopServiceV2.DescribeAgentloopOptimizeTask(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_agentloop_optimize_task DescribeAgentloopOptimizeTask Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("config", objectRaw["config"])
	d.Set("create_time", objectRaw["createTime"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("status", objectRaw["status"])
	d.Set("type", objectRaw["type"])
	d.Set("update_time", objectRaw["updateTime"])
	d.Set("optimize_task_name", objectRaw["optimizeTaskName"])

	parts := strings.Split(d.Id(), ":")
	d.Set("agent_space", parts[0])

	return nil
}

func resourceAliCloudAgentloopOptimizeTaskDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	agentSpace := parts[0]
	optimizeTaskName := parts[1]
	action := fmt.Sprintf("/agentspace/%s/optimizetasks/%s", agentSpace, optimizeTaskName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("AgentLoop", "2026-05-20", action, query, nil, nil, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict."}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"OptimizeTaskNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
