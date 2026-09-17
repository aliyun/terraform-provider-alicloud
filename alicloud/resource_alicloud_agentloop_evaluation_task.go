// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudAgentloopEvaluationTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAgentloopEvaluationTaskCreate,
		Read:   resourceAliCloudAgentloopEvaluationTaskRead,
		Update: resourceAliCloudAgentloopEvaluationTaskUpdate,
		Delete: resourceAliCloudAgentloopEvaluationTaskDelete,
		Importer: &schema.ResourceImporter{
			State: agentloopImportStatePassthrough(2),
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
				// The API defaults channel to "default" when it is not sent, so the
				// attribute is Computed as well: without it the backfilled value would
				// diff against an empty configuration and force a recreation.
				Computed: true,
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
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: agentloopJsonStringDiffSuppress,
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
						// In reference mode (evaluator_ref set) the API falls back to
						// the referenced evaluator for name/result_name and returns
						// the effective type, hence Optional+Computed.
						"result_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"AGENT", "CODE", "LLM"}, false),
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
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"result_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"score", "binary", "text"}, false),
						},
						"evaluator_ref": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: StringInSlice([]string{"HOUR", "MINUTE", "DAY"}, false),
									},
									// The API treats an omitted enabled as true.
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"data_delay_minutes": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
									},
									"interval_value": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntAtLeast(1),
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
									// The API treats an omitted enabled as true.
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
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
	if runStrategies := expandAgentloopEvaluationTaskRunStrategies(d.Get("run_strategies")); len(runStrategies) > 0 {
		request["runStrategies"] = runStrategies
	}

	if v, ok := d.GetOk("evaluators"); ok {
		evaluators, err := expandAgentloopEvaluationTaskEvaluators(v)
		if err != nil {
			return WrapError(err)
		}
		request["evaluators"] = evaluators
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

	// channel/dataType/taskMode/taskName are ForceNew attributes: when the API
	// omits one of them, keep the state value instead of clearing it, otherwise
	// the cleared state would force a destroy/create cycle on every plan.
	if objectRaw["channel"] != nil {
		if err := d.Set("channel", objectRaw["channel"]); err != nil {
			return err
		}
	}
	// GetEvaluationTask returns the effective config, which the backend enriches
	// with data-source pointers it derived itself. Those derived keys are only
	// tracked when the user manages them; every other key is read back.
	configuredConfig, _ := d.Get("config").(map[string]interface{})
	if err := d.Set("config", flattenAgentloopEvaluationTaskStringMap(objectRaw["config"], configuredConfig, agentloopEvaluationTaskDerivedConfigKeys)); err != nil {
		return err
	}
	if err := d.Set("created_at", objectRaw["createdAt"]); err != nil {
		return err
	}
	if objectRaw["dataFilter"] != nil {
		if err := d.Set("data_filter", objectRaw["dataFilter"]); err != nil {
			return err
		}
	}
	if objectRaw["dataType"] != nil {
		if err := d.Set("data_type", objectRaw["dataType"]); err != nil {
			return err
		}
	}
	if err := d.Set("description", objectRaw["description"]); err != nil {
		return err
	}
	// GetEvaluationTask does not return regionId. The task lives in the region of
	// the client endpoint (agentloop.<region>.aliyuncs.com), so fall back to the
	// client region when the API omits the field.
	if v, ok := objectRaw["regionId"]; ok && fmt.Sprint(v) != "" {
		if err := d.Set("region_id", v); err != nil {
			return err
		}
	} else {
		if err := d.Set("region_id", client.RegionId); err != nil {
			return err
		}
	}
	if err := d.Set("status", objectRaw["status"]); err != nil {
		return err
	}
	if objectRaw["taskMode"] != nil {
		if err := d.Set("task_mode", objectRaw["taskMode"]); err != nil {
			return err
		}
	}
	if objectRaw["taskName"] != nil {
		if err := d.Set("task_name", objectRaw["taskName"]); err != nil {
			return err
		}
	}
	if err := d.Set("task_id", objectRaw["taskId"]); err != nil {
		return err
	}
	if objectRaw["tags"] != nil {
		if err := d.Set("tags", objectRaw["tags"]); err != nil {
			return err
		}
	}

	if err := d.Set("evaluators", flattenAgentloopEvaluationTaskEvaluators(objectRaw["evaluators"], d.Get("evaluators"))); err != nil {
		return err
	}
	if err := d.Set("run_strategies", flattenAgentloopEvaluationTaskRunStrategies(objectRaw["runStrategyConfig"])); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	if err := d.Set("agent_space", parts[0]); err != nil {
		return err
	}

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
	// Only changed attributes are sent: the API replaces config wholesale
	// (dropping the derived data-source keys) and status carries the
	// backend-managed lifecycle value, so neither may be echoed back on an
	// unrelated update.
	if !d.IsNewResource() && d.HasChange("config") {
		update = true
		request["config"] = d.Get("config")
	}

	if !d.IsNewResource() && d.HasChange("evaluators") {
		update = true
		evaluators, err := expandAgentloopEvaluationTaskEvaluators(d.Get("evaluators"))
		if err != nil {
			return WrapError(err)
		}
		request["evaluators"] = evaluators
	}

	// runStrategies is replaced as a whole; an empty object clears it.
	if !d.IsNewResource() && d.HasChange("run_strategies") {
		update = true
		request["runStrategies"] = expandAgentloopEvaluationTaskRunStrategies(d.Get("run_strategies"))
	}

	if !d.IsNewResource() && d.HasChange("data_filter") {
		update = true
		request["dataFilter"] = d.Get("data_filter")
	}
	if !d.IsNewResource() && d.HasChange("tags") {
		update = true
		request["tags"] = d.Get("tags")
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["description"] = d.Get("description")
	}
	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			update = true
			request["status"] = v
		}
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

// agentloopJsonStringDiffSuppress suppresses diffs between two strings that
// carry semantically equal JSON documents. GetEvaluationTask returns dataFilter
// as a JSON string whose formatting (key order, whitespace) is chosen by the
// backend, so a plain byte comparison would flag a perpetual diff.
func agentloopJsonStringDiffSuppress(k, oldValue, newValue string, d *schema.ResourceData) bool {
	if oldValue == "" || newValue == "" {
		return oldValue == newValue
	}
	return agentloopJsonStringsEquivalent(oldValue, newValue)
}

// agentloopEvaluationTaskDerivedConfigKeys are injected into the task config
// by the backend (data source resolved from the AgentSpace).
var agentloopEvaluationTaskDerivedConfigKeys = map[string]bool{
	"dataScope":   true,
	"project":     true,
	"storeName":   true,
	"traceFormat": true,
}

// agentloopEvaluationTaskDerivedEvaluatorConfigKeys are injected into each
// evaluator config by the backend (resolved evaluator version).
var agentloopEvaluationTaskDerivedEvaluatorConfigKeys = map[string]bool{
	"version": true,
}

// flattenAgentloopEvaluationTaskStringMap converts an API map into the
// TypeMap-of-strings shape. Keys listed in derived are skipped unless the
// user manages them (present in configured); every other key is read back so
// external drift is detected. Non-string values are serialized to JSON, and
// the configured text is kept when it is semantically equal.
func flattenAgentloopEvaluationTaskStringMap(raw interface{}, configured map[string]interface{}, derived map[string]bool) map[string]interface{} {
	result := make(map[string]interface{})
	rawMap, ok := raw.(map[string]interface{})
	if !ok {
		return result
	}
	for k, v := range rawMap {
		if v == nil {
			continue
		}
		if _, managed := configured[k]; !managed && derived[k] {
			continue
		}
		if s, ok := v.(string); ok {
			result[k] = s
			continue
		}
		bs, err := json.Marshal(v)
		if err != nil {
			continue
		}
		flat := string(bs)
		if orig, ok := configured[k].(string); ok && agentloopJsonStringsEquivalent(orig, flat) {
			flat = orig
		}
		result[k] = flat
	}
	return result
}

// agentloopEvaluationTaskEvaluatorIdentity returns the stable identity of an
// evaluator entry: the referenced evaluator, or the inline name.
func agentloopEvaluationTaskEvaluatorIdentity(ref, name interface{}) string {
	if s, ok := ref.(string); ok && s != "" {
		return "ref:" + s
	}
	if s, ok := name.(string); ok && s != "" {
		return "name:" + s
	}
	return ""
}

func agentloopStringValue(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// flattenAgentloopEvaluationTaskEvaluators reads the resolved evaluators back.
// Each returned entry is paired with the configured entry at the same index
// when their identity matches, otherwise with the first unused configured
// entry of the same identity. For referenced evaluators the backend replaces
// name with the evaluator_ref value, so a configured name is kept in that
// case; the injected config.version is only tracked when configured.
func flattenAgentloopEvaluationTaskEvaluators(raw interface{}, configuredRaw interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	apiList, ok := raw.([]interface{})
	if !ok {
		return result
	}
	configured := make([]map[string]interface{}, 0)
	if l, ok := configuredRaw.([]interface{}); ok {
		for _, item := range l {
			m, _ := item.(map[string]interface{})
			configured = append(configured, m)
		}
	}
	used := make([]bool, len(configured))
	for i, item := range apiList {
		apiItem, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		id := agentloopEvaluationTaskEvaluatorIdentity(apiItem["evaluatorRef"], apiItem["name"])
		var conf map[string]interface{}
		if i < len(configured) && !used[i] && configured[i] != nil && agentloopEvaluationTaskEvaluatorIdentity(configured[i]["evaluator_ref"], configured[i]["name"]) == id {
			conf, used[i] = configured[i], true
		} else {
			for j, c := range configured {
				if !used[j] && c != nil && agentloopEvaluationTaskEvaluatorIdentity(c["evaluator_ref"], c["name"]) == id {
					conf, used[j] = c, true
					break
				}
			}
		}
		confConfig, _ := conf["config"].(map[string]interface{})
		confFilters, _ := conf["filters"].(map[string]interface{})
		confMapping, _ := conf["variable_mapping"].(map[string]interface{})

		ref := agentloopStringValue(apiItem["evaluatorRef"])
		name := agentloopStringValue(apiItem["name"])
		if confName := agentloopStringValue(conf["name"]); confName != "" && ref != "" && name == ref {
			name = confName
		}
		result = append(result, map[string]interface{}{
			"evaluator_ref":    ref,
			"name":             name,
			"type":             agentloopStringValue(apiItem["type"]),
			"result_name":      agentloopStringValue(apiItem["resultName"]),
			"result_type":      agentloopStringValue(apiItem["resultType"]),
			"config":           flattenAgentloopEvaluationTaskStringMap(apiItem["config"], confConfig, agentloopEvaluationTaskDerivedEvaluatorConfigKeys),
			"filters":          flattenAgentloopEvaluationTaskStringMap(apiItem["filters"], confFilters, nil),
			"variable_mapping": flattenAgentloopEvaluationTaskStringMap(apiItem["variableMapping"], confMapping, nil),
		})
	}
	return result
}

// expandAgentloopEvaluationTaskEvaluators builds the evaluators request body.
// Empty strings and maps are omitted so the API applies its own fallback
// (reference mode, default result type) instead of receiving "" or {}.
// Inline evaluators (no evaluator_ref) must carry name and result_type, LLM
// and AGENT inline evaluators also need variable_mapping, and CODE evaluators
// can only be referenced.
func expandAgentloopEvaluationTaskEvaluators(v interface{}) ([]interface{}, error) {
	result := make([]interface{}, 0)
	for i, dataLoop := range convertToInterfaceArray(v) {
		dataLoopTmp, ok := dataLoop.(map[string]interface{})
		if !ok {
			continue
		}
		dataLoopMap := make(map[string]interface{})
		for tfKey, apiKey := range map[string]string{
			"evaluator_ref": "evaluatorRef",
			"name":          "name",
			"type":          "type",
			"result_name":   "resultName",
			"result_type":   "resultType",
		} {
			if s := agentloopStringValue(dataLoopTmp[tfKey]); s != "" {
				dataLoopMap[apiKey] = s
			}
		}
		for tfKey, apiKey := range map[string]string{
			"variable_mapping": "variableMapping",
			"filters":          "filters",
			"config":           "config",
		} {
			if m, ok := dataLoopTmp[tfKey].(map[string]interface{}); ok && len(m) > 0 {
				dataLoopMap[apiKey] = m
			}
		}
		if _, isRef := dataLoopMap["evaluatorRef"]; !isRef {
			if _, ok := dataLoopMap["name"]; !ok {
				return nil, fmt.Errorf("evaluators.%d: name is required when evaluator_ref is not set", i)
			}
			if _, ok := dataLoopMap["resultType"]; !ok {
				return nil, fmt.Errorf("evaluators.%d: result_type is required when evaluator_ref is not set", i)
			}
			switch dataLoopMap["type"] {
			case "CODE":
				return nil, fmt.Errorf("evaluators.%d: CODE evaluators must be referenced via evaluator_ref", i)
			default:
				if _, ok := dataLoopMap["variableMapping"]; !ok {
					return nil, fmt.Errorf("evaluators.%d: variable_mapping is required for inline LLM/AGENT evaluators", i)
				}
			}
		}
		result = append(result, dataLoopMap)
	}
	return result, nil
}

// expandAgentloopEvaluationTaskRunStrategies builds the runStrategies body.
// enabled and immediate are always sent (enabled defaults to true, matching
// the API's omission semantics); zero interval values and zero timestamps are
// omitted rather than sent as invalid 0 values.
func expandAgentloopEvaluationTaskRunStrategies(v interface{}) map[string]interface{} {
	runStrategies := make(map[string]interface{})
	list, _ := v.([]interface{})
	if len(list) == 0 {
		return runStrategies
	}
	item, _ := list[0].(map[string]interface{})
	if l, ok := item["continuous"].([]interface{}); ok && len(l) > 0 {
		if c, ok := l[0].(map[string]interface{}); ok {
			continuous := make(map[string]interface{})
			if b, ok := c["enabled"].(bool); ok {
				continuous["enabled"] = b
			}
			if s := agentloopStringValue(c["interval_unit"]); s != "" {
				continuous["intervalUnit"] = s
			}
			if n, ok := c["interval_value"].(int); ok && n > 0 {
				continuous["intervalValue"] = n
			}
			if n, ok := c["data_delay_minutes"].(int); ok {
				continuous["dataDelayMinutes"] = n
			}
			runStrategies["continuous"] = continuous
		}
	}
	if l, ok := item["backfill"].([]interface{}); ok && len(l) > 0 {
		if b, ok := l[0].(map[string]interface{}); ok {
			backfill := make(map[string]interface{})
			if v, ok := b["enabled"].(bool); ok {
				backfill["enabled"] = v
			}
			if v, ok := b["immediate"].(bool); ok {
				backfill["immediate"] = v
			}
			if n, ok := b["start_time"].(int); ok && n > 0 {
				backfill["startTime"] = n
			}
			if n, ok := b["end_time"].(int); ok && n > 0 {
				backfill["endTime"] = n
			}
			runStrategies["backfill"] = backfill
		}
	}
	return runStrategies
}

// flattenAgentloopEvaluationTaskRunStrategies reads runStrategyConfig back.
// The API echoes what was sent; an omitted enabled means enabled (true) and
// an omitted immediate means false.
func flattenAgentloopEvaluationTaskRunStrategies(raw interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	rawMap, ok := raw.(map[string]interface{})
	if !ok || len(rawMap) == 0 {
		return result
	}
	item := make(map[string]interface{})
	continuousMaps := make([]map[string]interface{}, 0)
	if c, ok := rawMap["continuous"].(map[string]interface{}); ok && len(c) > 0 {
		continuous := map[string]interface{}{
			"interval_unit": agentloopStringValue(c["intervalUnit"]),
			"enabled":       true,
		}
		if v, ok := c["enabled"].(bool); ok {
			continuous["enabled"] = v
		}
		if v := c["intervalValue"]; v != nil {
			continuous["interval_value"] = formatInt(v)
		}
		if v := c["dataDelayMinutes"]; v != nil {
			continuous["data_delay_minutes"] = formatInt(v)
		}
		continuousMaps = append(continuousMaps, continuous)
	}
	item["continuous"] = continuousMaps
	backfillMaps := make([]map[string]interface{}, 0)
	if b, ok := rawMap["backfill"].(map[string]interface{}); ok && len(b) > 0 {
		backfill := map[string]interface{}{
			"enabled":   true,
			"immediate": false,
		}
		if v, ok := b["enabled"].(bool); ok {
			backfill["enabled"] = v
		}
		if v, ok := b["immediate"].(bool); ok {
			backfill["immediate"] = v
		}
		if v := b["startTime"]; v != nil {
			backfill["start_time"] = formatInt(v)
		}
		if v := b["endTime"]; v != nil {
			backfill["end_time"] = formatInt(v)
		}
		backfillMaps = append(backfillMaps, backfill)
	}
	item["backfill"] = backfillMaps
	if len(continuousMaps) == 0 && len(backfillMaps) == 0 {
		return result
	}
	return append(result, item)
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
