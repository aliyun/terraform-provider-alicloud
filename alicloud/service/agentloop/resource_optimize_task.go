package agentloop

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceOptimizeTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceOptimizeTaskCreate,
		Read:   resourceOptimizeTaskRead,
		Update: resourceOptimizeTaskUpdate,
		Delete: resourceOptimizeTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: agentloopImportStatePassthrough(2),
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
			"config": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: agentloopJsonStringDiffSuppress,
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
				ValidateFunc: validation.StringInSlice([]string{"Enabled", "Disabled"}, false),
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

func resourceOptimizeTaskCreate(d *schema.ResourceData, meta interface{}) error {

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
	clientToken := id.UniqueId()
	query["clientToken"] = &clientToken
	if v, ok := d.GetOk("optimize_task_name"); ok {
		request["optimizeTaskName"] = v
	}

	if v, ok := d.GetOk("config"); ok {
		var configObj interface{}
		if err := json.Unmarshal([]byte(v.(string)), &configObj); err != nil {
			return fmt.Errorf("config must be valid JSON: %s", err)
		}
		request["config"] = configObj
	}
	request["status"] = d.Get("status")
	request["type"] = d.Get("type")
	body = request
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RoaPost("AgentLoop", "2026-05-20", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"Conflict."}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_agentloop_optimize_task", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", agentSpace, request["optimizeTaskName"]))

	return resourceOptimizeTaskRead(d, meta)
}

func resourceOptimizeTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	agentloopService := NewService(client)

	objectRaw, err := agentloopService.DescribeOptimizeTask(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_agentloop_optimize_task DescribeOptimizeTask Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	if objectRaw["config"] != nil {
		configBytes, _ := json.Marshal(objectRaw["config"])
		if err := d.Set("config", string(configBytes)); err != nil {
			return err
		}
	}
	if err := d.Set("create_time", objectRaw["createTime"]); err != nil {
		return err
	}
	if err := d.Set("region_id", objectRaw["regionId"]); err != nil {
		return err
	}
	if err := d.Set("status", objectRaw["status"]); err != nil {
		return err
	}
	if err := d.Set("type", objectRaw["type"]); err != nil {
		return err
	}
	if err := d.Set("update_time", objectRaw["updateTime"]); err != nil {
		return err
	}
	if err := d.Set("optimize_task_name", objectRaw["optimizeTaskName"]); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	if err := d.Set("agent_space", parts[0]); err != nil {
		return err
	}

	return nil
}

func resourceOptimizeTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	agentSpace := parts[0]
	optimizeTaskName := parts[1]
	action := fmt.Sprintf("/agentspace/%s/optimizetasks/%s", agentSpace, optimizeTaskName)
	request := make(map[string]interface{})
	query := make(map[string]*string)
	update := false

	if d.HasChange("config") {
		update = true
		if v, ok := d.GetOk("config"); ok {
			var configObj interface{}
			if err := json.Unmarshal([]byte(v.(string)), &configObj); err != nil {
				return fmt.Errorf("config must be valid JSON: %s", err)
			}
			request["config"] = configObj
		}
	}
	if d.HasChange("status") {
		update = true
		request["status"] = d.Get("status")
	}
	if update {
		clientToken := id.UniqueId()
		query["clientToken"] = &clientToken
		var response map[string]interface{}
		var err error
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RoaPut("AgentLoop", "2026-05-20", action, query, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"Conflict."}) {
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

	return resourceOptimizeTaskRead(d, meta)
}

func resourceOptimizeTaskDelete(d *schema.ResourceData, meta interface{}) error {

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

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RoaDelete("AgentLoop", "2026-05-20", action, query, nil, nil, true)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"Conflict."}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"OptimizeTaskNotExist"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
