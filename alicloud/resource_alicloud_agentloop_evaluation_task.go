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

func resourceAliCloudAgentloopEvaluationTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAgentloopEvaluationTaskCreate,
		Read:   resourceAliCloudAgentloopEvaluationTaskRead,
		Update: resourceAliCloudAgentloopEvaluationTaskUpdate,
		Delete: resourceAliCloudAgentloopEvaluationTaskDelete,
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
			"channel": {
				Type:     schema.TypeString,
				Optional: true,
				// The UpdateEvaluationTask spec marks channel @visibility("Private"),
				// i.e. the public update silently ignores it (probe-verified: PUT
				// channel returns 200 but the stored value never changes), so the
				// attribute can only be set at creation time.
				ForceNew: true,
			},
			"config": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"evaluators": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"result_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"filters": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"config": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"variable_mapping": {
							Type:     schema.TypeMap,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"result_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"evaluator_ref": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"run_strategies": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"continuous": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"interval_unit": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"data_delay_minutes": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"interval_value": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"backfill": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_time": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"immediate": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"start_time": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Pending", "Running", "Completed", "Scheduling", "Failed", "Terminated"}, false),
			},
			"tags": tagsSchema(),
			"task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudAgentloopEvaluationTaskCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	agentSpace := d.Get("agent_space")
	action := fmt.Sprintf("/api/v1/evaluation-task/%s", agentSpace)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("channel"); ok {
		request["channel"] = v
	}
	if v, ok := d.GetOk("config"); ok {
		request["config"] = v
	}
	runStrategies := make(map[string]interface{})

	if v := d.Get("run_strategies"); !IsNil(v) {
		backfill := make(map[string]interface{})
		endTime1, _ := jsonpath.Get("$[0].backfill[0].end_time", v)
		if endTime1 != nil && endTime1 != "" {
			backfill["endTime"] = endTime1
		}
		startTime1, _ := jsonpath.Get("$[0].backfill[0].start_time", v)
		if startTime1 != nil && startTime1 != "" {
			backfill["startTime"] = startTime1
		}
		immediate1, _ := jsonpath.Get("$[0].backfill[0].immediate", v)
		if immediate1 != nil && immediate1 != "" {
			backfill["immediate"] = immediate1
		}
		enabled1, _ := jsonpath.Get("$[0].backfill[0].enabled", v)
		if enabled1 != nil && enabled1 != "" {
			backfill["enabled"] = enabled1
		}

		if len(backfill) > 0 {
			runStrategies["backfill"] = backfill
		}
		continuous := make(map[string]interface{})
		intervalValue1, _ := jsonpath.Get("$[0].continuous[0].interval_value", v)
		if intervalValue1 != nil && intervalValue1 != "" {
			continuous["intervalValue"] = intervalValue1
		}
		dataDelayMinutes1, _ := jsonpath.Get("$[0].continuous[0].data_delay_minutes", v)
		if dataDelayMinutes1 != nil && dataDelayMinutes1 != "" {
			continuous["dataDelayMinutes"] = dataDelayMinutes1
		}
		enabled3, _ := jsonpath.Get("$[0].continuous[0].enabled", v)
		if enabled3 != nil && enabled3 != "" {
			continuous["enabled"] = enabled3
		}
		intervalUnit1, _ := jsonpath.Get("$[0].continuous[0].interval_unit", v)
		if intervalUnit1 != nil && intervalUnit1 != "" {
			continuous["intervalUnit"] = intervalUnit1
		}

		if len(continuous) > 0 {
			runStrategies["continuous"] = continuous
		}

		if len(runStrategies) > 0 {
			request["runStrategies"] = runStrategies
		}
	}

	if v, ok := d.GetOk("evaluators"); ok {
		evaluatorsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["evaluatorRef"] = dataLoopTmp["evaluator_ref"]
			dataLoopMap["variableMapping"] = dataLoopTmp["variable_mapping"]
			dataLoopMap["resultName"] = dataLoopTmp["result_name"]
			dataLoopMap["filters"] = dataLoopTmp["filters"]
			dataLoopMap["type"] = dataLoopTmp["type"]
			dataLoopMap["config"] = dataLoopTmp["config"]
			dataLoopMap["name"] = dataLoopTmp["name"]
			dataLoopMap["resultType"] = dataLoopTmp["result_type"]
			evaluatorsMapsArray = append(evaluatorsMapsArray, dataLoopMap)
		}
		request["evaluators"] = evaluatorsMapsArray
	}

	if v, ok := d.GetOk("data_filter"); ok {
		request["dataFilter"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		request["tags"] = v
	}
	request["taskName"] = d.Get("task_name")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("data_type"); ok {
		request["dataType"] = v
	}
	if v, ok := d.GetOk("task_mode"); ok {
		request["taskMode"] = v
	}
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_agentloop_evaluation_task", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", agentSpace, response["taskId"]))

	return resourceAliCloudAgentloopEvaluationTaskUpdate(d, meta)
}

func resourceAliCloudAgentloopEvaluationTaskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	agentloopServiceV2 := AgentloopServiceV2{client}

	objectRaw, err := agentloopServiceV2.DescribeAgentloopEvaluationTask(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_agentloop_evaluation_task DescribeAgentloopEvaluationTask Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("channel", objectRaw["channel"])
	// config is intentionally NOT read back: GetEvaluationTask returns the
	// RESOLVED effective config, which the backend enriches with data-source
	// pointers it derived itself (dataScope/project/storeName/traceFormat).
	// Writing that back would produce a perpetual diff against the user
	// configuration, so the configured value is kept in state.
	d.Set("created_at", objectRaw["createdAt"])
	d.Set("data_filter", objectRaw["dataFilter"])
	d.Set("data_type", objectRaw["dataType"])
	d.Set("description", objectRaw["description"])
	// GetEvaluationTask does not return regionId. The task lives in the region of
	// the client endpoint (agentloop.<region>.aliyuncs.com), so fall back to the
	// client region when the API omits the field.
	if v, ok := objectRaw["regionId"]; ok && fmt.Sprint(v) != "" {
		d.Set("region_id", v)
	} else {
		d.Set("region_id", client.RegionId)
	}
	d.Set("status", objectRaw["status"])
	d.Set("task_mode", objectRaw["taskMode"])
	d.Set("task_name", objectRaw["taskName"])
	d.Set("task_id", objectRaw["taskId"])
	if objectRaw["tags"] != nil {
		d.Set("tags", objectRaw["tags"])
	}

	// evaluators is intentionally NOT read back: GetEvaluationTask resolves
	// each evaluator entry (name is replaced by the referenced evaluator's
	// real name, and the resolved version is injected into config), so the
	// response never matches the configured value and would cause a
	// perpetual diff. The configured value is kept in state.

	parts := strings.Split(d.Id(), ":")
	d.Set("agent_space", parts[0])

	return nil
}

func resourceAliCloudAgentloopEvaluationTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	agentSpace := parts[0]
	taskId := parts[1]
	action := fmt.Sprintf("/api/v1/evaluation-task/%s/%s", agentSpace, taskId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	// NOTE: channel is not sent here. It is @visibility("Private") in the
	// UpdateEvaluationTask spec and the public update silently ignores it
	// (probe-verified), so the schema marks it ForceNew instead.
	if !d.IsNewResource() && d.HasChange("config") {
		update = true
	}
	if v, ok := d.GetOk("config"); ok || d.HasChange("config") {
		request["config"] = v
	}

	if !d.IsNewResource() && d.HasChange("evaluators") {
		update = true
	}
	if v, ok := d.GetOk("evaluators"); ok || d.HasChange("evaluators") {
		evaluatorsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["evaluatorRef"] = dataLoopTmp["evaluator_ref"]
			dataLoopMap["variableMapping"] = dataLoopTmp["variable_mapping"]
			dataLoopMap["resultName"] = dataLoopTmp["result_name"]
			dataLoopMap["filters"] = dataLoopTmp["filters"]
			dataLoopMap["type"] = dataLoopTmp["type"]
			dataLoopMap["config"] = dataLoopTmp["config"]
			dataLoopMap["name"] = dataLoopTmp["name"]
			dataLoopMap["resultType"] = dataLoopTmp["result_type"]
			evaluatorsMapsArray = append(evaluatorsMapsArray, dataLoopMap)
		}
		request["evaluators"] = evaluatorsMapsArray
	}

	if !d.IsNewResource() && d.HasChange("data_filter") {
		update = true
	}
	if v, ok := d.GetOk("data_filter"); ok || d.HasChange("data_filter") {
		request["dataFilter"] = v
	}
	if !d.IsNewResource() && d.HasChange("tags") {
		update = true
	}
	if v, ok := d.GetOk("tags"); ok || d.HasChange("tags") {
		request["tags"] = v
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	if d.HasChange("status") {
		update = true
	}
	if v, ok := d.GetOk("status"); ok || d.HasChange("status") {
		request["status"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("AgentLoop", "2026-05-20", action, query, nil, body, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	// status transitions are applied asynchronously by the backend (e.g.
	// Pending -> Running). When the configuration carries an explicit status,
	// wait until the backend reaches it so the post-apply state matches the
	// configuration (the transition can lag the update response by tens of
	// seconds, which otherwise causes a non-empty plan right after apply).
	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			targetStatus := fmt.Sprint(v)
			pendingStates, failStates := agentloopEvaluationTaskStatusWaitSets(targetStatus)
			agentloopServiceV2 := AgentloopServiceV2{client}
			stateConf := BuildStateConf(pendingStates, []string{targetStatus}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, agentloopServiceV2.AgentloopEvaluationTaskStateRefreshFunc(d.Id(), "status", failStates))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapError(err)
			}
		}
	}

	d.Partial(false)
	return resourceAliCloudAgentloopEvaluationTaskRead(d, meta)
}

// agentloopEvaluationTaskStatusWaitSets derives the pending and fail state
// sets for waiting on an EvaluationTask status transition to targetStatus.
// Every other non-terminal state is a legal intermediate (e.g. Running on
// the way to Completed), and a terminal state requested as the target must
// not be treated as a failure when reached.
func agentloopEvaluationTaskStatusWaitSets(targetStatus string) (pendingStates []string, failStates []string) {
	pendingStates = make([]string, 0)
	for _, s := range []string{"Pending", "Scheduling", "Running"} {
		if s != targetStatus {
			pendingStates = append(pendingStates, s)
		}
	}
	failStates = make([]string, 0)
	for _, s := range []string{"Failed", "Terminated", "Deleted"} {
		if s != targetStatus {
			failStates = append(failStates, s)
		}
	}
	return pendingStates, failStates
}

func resourceAliCloudAgentloopEvaluationTaskDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	agentSpace := parts[0]
	taskId := parts[1]
	action := fmt.Sprintf("/api/v1/evaluation-task/%s/%s", agentSpace, taskId)
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
		if IsExpectedErrors(err, []string{"EvaluationTaskNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
