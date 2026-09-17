package alicloud

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Agentloop EvaluationTask. >>> Resource test cases, automatically generated.
// Case EvaluationTask全生命周期 11081
func TestAccAliCloudAgentloopEvaluationTask_basic11081(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_agentloop_evaluation_task.default"
	ra := resourceAttrInit(resourceId, AlicloudAgentloopEvaluationTaskMap11081)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AgentloopServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAgentloopEvaluationTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccevaltask%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAgentloopEvaluationTaskBasicDependence11081)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
					"task_name":   name,
					"task_mode":   "batch",
					"data_type":   "trace",
					"channel":     "default",
					"data_filter": `{\"query\":\"checkout-service\",\"maxRecords\":10}`,
					"description": "terraform自动化测试描述",
					// NOTE: no "status" in this step on purpose. With
					// backfill.enabled=true + immediate=true the task flips to
					// Failed within seconds in this no-trace-data env (probe34 A),
					// and a status wait here would fail on that terminal state.
					// Skipping the status attribute means the post-create Update
					// does not wait (HasChange("status")=false), so the Failed
					// transition is harmless. status is exercised in step 1 where
					// the recreated task (all strategies disabled) stays Pending.
					"config": map[string]interface{}{
						"config1": "value1",
					},
					"tags": map[string]interface{}{
						"tag1": "value1",
					},
					"run_strategies": []map[string]interface{}{
						{
							"continuous": []map[string]interface{}{
								{
									"enabled":            "true",
									"interval_unit":      "hours",
									"interval_value":     "6",
									"data_delay_minutes": "30",
								},
							},
							"backfill": []map[string]interface{}{
								{
									"enabled":    "true",
									"immediate":  "true",
									"start_time": "1735689600000",
									"end_time":   "1735776000000",
								},
							},
						},
					},
					"evaluators": []map[string]interface{}{
						{
							"name":          "evaluatora",
							"result_name":   "resulta",
							"type":          "AGENT",
							"evaluator_ref": "${alicloud_agentloop_evaluator.default.name}",
							"filters": map[string]interface{}{
								"filter1": "value1",
							},
							"config": map[string]interface{}{
								"config1": "value1",
							},
							"variable_mapping": map[string]interface{}{
								"input": "$.input",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"agent_space":  name + "-as",
						"task_name":    name,
						"task_mode":    "batch",
						"data_type":    "trace",
						"channel":      "default",
						"description":  "terraform自动化测试描述",
						"evaluators.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"channel":     "console",
					"data_filter": `{\"query\":\"order-service\",\"maxRecords\":20}`,
					"description": "terraform自动化测试描述更新",
					// The channel change recreates the task (ForceNew), and the
					// recreated task has every run strategy disabled below, so it
					// stays Pending (probe34 B). The status wait in Update then
					// succeeds immediately. Do NOT wait for "Running": without
					// trace data the task never leaves Pending (probe > 5min).
					"status": "Pending",
					"config": map[string]interface{}{
						"config1": "value2",
					},
					"tags": map[string]interface{}{
						"tag1": "value2",
					},
					"run_strategies": []map[string]interface{}{
						{
							"continuous": []map[string]interface{}{
								{
									"enabled":            "false",
									"interval_unit":      "minutes",
									"interval_value":     "12",
									"data_delay_minutes": "60",
								},
							},
							"backfill": []map[string]interface{}{
								{
									"enabled":    "false",
									"immediate":  "false",
									"start_time": "1735862400000",
									"end_time":   "1735948800000",
								},
							},
						},
					},
					"evaluators": []map[string]interface{}{
						{
							"name":          "evaluatorb",
							"result_name":   "resultb",
							"type":          "${alicloud_agentloop_evaluator.default2.type}",
							"result_type":   "score",
							"evaluator_ref": "${alicloud_agentloop_evaluator.default2.name}",
							"filters": map[string]interface{}{
								"filter1": "value2",
							},
							"config": map[string]interface{}{
								"config1": "value2",
							},
							"variable_mapping": map[string]interface{}{
								"input": "$.output",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"channel":      "console",
						"description":  "terraform自动化测试描述更新",
						"evaluators.#": "1",
					}),
				),
			},
			{
				// TypeList order coverage baseline for evaluators: two distinct
				// entries referencing the two dependence evaluators.
				Config: testAccConfig(map[string]interface{}{
					"evaluators": []map[string]interface{}{
						{
							"name":          "evaluatora",
							"result_name":   "resulta",
							"type":          "AGENT",
							"evaluator_ref": "${alicloud_agentloop_evaluator.default.name}",
							"filters": map[string]interface{}{
								"filter1": "value1",
							},
							"config": map[string]interface{}{
								"config1": "value1",
							},
							"variable_mapping": map[string]interface{}{
								"input": "$.input",
							},
						},
						{
							"name":          "evaluatorb",
							"result_name":   "resultb",
							"type":          "${alicloud_agentloop_evaluator.default2.type}",
							"result_type":   "score",
							"evaluator_ref": "${alicloud_agentloop_evaluator.default2.name}",
							"filters": map[string]interface{}{
								"filter1": "value2",
							},
							"config": map[string]interface{}{
								"config1": "value2",
							},
							"variable_mapping": map[string]interface{}{
								"input": "$.output",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"evaluators.#": "2",
					}),
				),
			},
			{
				// TypeList order coverage: a reorder-only plan against the
				// previous apply must produce a non-empty plan.
				Config: testAccConfig(map[string]interface{}{
					"evaluators": []map[string]interface{}{
						{
							"name":          "evaluatorb",
							"result_name":   "resultb",
							"type":          "${alicloud_agentloop_evaluator.default2.type}",
							"result_type":   "score",
							"evaluator_ref": "${alicloud_agentloop_evaluator.default2.name}",
							"filters": map[string]interface{}{
								"filter1": "value2",
							},
							"config": map[string]interface{}{
								"config1": "value2",
							},
							"variable_mapping": map[string]interface{}{
								"input": "$.output",
							},
						},
						{
							"name":          "evaluatora",
							"result_name":   "resulta",
							"type":          "AGENT",
							"evaluator_ref": "${alicloud_agentloop_evaluator.default.name}",
							"filters": map[string]interface{}{
								"filter1": "value1",
							},
							"config": map[string]interface{}{
								"config1": "value1",
							},
							"variable_mapping": map[string]interface{}{
								"input": "$.input",
							},
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Applying the reordered evaluators must then converge with an
				// empty post-apply plan. evaluators is never read back (the API
				// returns the resolved form), so convergence is guaranteed at
				// the state level.
				Config: testAccConfig(map[string]interface{}{
					"evaluators": []map[string]interface{}{
						{
							"name":          "evaluatorb",
							"result_name":   "resultb",
							"type":          "${alicloud_agentloop_evaluator.default2.type}",
							"result_type":   "score",
							"evaluator_ref": "${alicloud_agentloop_evaluator.default2.name}",
							"filters": map[string]interface{}{
								"filter1": "value2",
							},
							"config": map[string]interface{}{
								"config1": "value2",
							},
							"variable_mapping": map[string]interface{}{
								"input": "$.output",
							},
						},
						{
							"name":          "evaluatora",
							"result_name":   "resulta",
							"type":          "AGENT",
							"evaluator_ref": "${alicloud_agentloop_evaluator.default.name}",
							"filters": map[string]interface{}{
								"filter1": "value1",
							},
							"config": map[string]interface{}{
								"config1": "value1",
							},
							"variable_mapping": map[string]interface{}{
								"input": "$.input",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"evaluators.#": "2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// config/evaluators/run_strategies are not read back (the API
				// returns backend-resolved forms, see the resource Read), so the
				// imported state cannot contain them. status is ignored too:
				// the task transitions Pending -> Running asynchronously, so
				// the imported state may differ from the post-apply state.
				ImportStateVerifyIgnore: []string{"config", "evaluators", "status", "run_strategies"},
			},
		},
	})
}

var AlicloudAgentloopEvaluationTaskMap11081 = map[string]string{
	"task_id":    CHECKSET,
	"region_id":  CHECKSET,
	"created_at": CHECKSET,
	"status":     CHECKSET,
}

func AlicloudAgentloopEvaluationTaskBasicDependence11081(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_agentloop_agent_space" "default" {
    agent_space = "${var.name}-as"
}

resource "alicloud_agentloop_evaluator" "default" {
    agent_space = alicloud_agentloop_agent_space.default.agent_space
    name        = "${var.name}eval"
    type        = "AGENT"
    metric_name = "tf_test_metric"
    version     = "v1"
    config = {
        variables = jsonencode([{ name = "input", type = "text" }])
    }
}

resource "alicloud_agentloop_evaluator" "default2" {
    agent_space = alicloud_agentloop_agent_space.default.agent_space
    name        = "${var.name}eval2"
    type        = "AGENT"
    metric_name = "tf_test_metric"
    version     = "v1"
    config = {
        variables = jsonencode([{ name = "input", type = "text" }])
    }
}

`, name)
}

// Test Agentloop EvaluationTask. <<< Resource test cases, automatically generated.

// TestUnitAliCloudAgentloopEvaluationTaskStatusWaitSets pins the pending/fail
// sets derived for every waitable target status: non-target intermediate
// states stay pending, and the terminal target state itself must never be
// classified as a failure.
func TestUnitAliCloudAgentloopEvaluationTaskStatusWaitSets(t *testing.T) {
	cases := []struct {
		target  string
		pending []string
		fail    []string
	}{
		{"Pending", []string{"Scheduling", "Running"}, []string{"Failed", "Terminated", "Deleted"}},
		{"Scheduling", []string{"Pending", "Running"}, []string{"Failed", "Terminated", "Deleted"}},
		{"Running", []string{"Pending", "Scheduling"}, []string{"Failed", "Terminated", "Deleted"}},
		{"Completed", []string{"Pending", "Scheduling", "Running"}, []string{"Failed", "Terminated", "Deleted"}},
		{"Failed", []string{"Pending", "Scheduling", "Running"}, []string{"Terminated", "Deleted"}},
		{"Terminated", []string{"Pending", "Scheduling", "Running"}, []string{"Failed", "Deleted"}},
		{"Deleted", []string{"Pending", "Scheduling", "Running"}, []string{"Failed", "Terminated"}},
	}
	for _, tc := range cases {
		pending, fail := agentloopEvaluationTaskStatusWaitSets(tc.target)
		if !reflect.DeepEqual(pending, tc.pending) {
			t.Fatalf("target %s: pending = %v, expected %v", tc.target, pending, tc.pending)
		}
		if !reflect.DeepEqual(fail, tc.fail) {
			t.Fatalf("target %s: fail = %v, expected %v", tc.target, fail, tc.fail)
		}
	}
}

func TestAccAliCloudAgentloopJsonStringDiffSuppress(t *testing.T) {
	cases := []struct {
		name     string
		oldValue string
		newValue string
		want     bool
	}{
		{"both empty", "", "", true},
		{"old empty", "", "{}", false},
		{"new empty", "{}", "", false},
		{"identical", `{"a":1}`, `{"a":1}`, true},
		{"semantically equal", `{"a":1,"b":2}`, `{"b":2,"a":1}`, true},
		{"semantically different", `{"a":1}`, `{"a":2}`, false},
		{"invalid old", `not-json`, `{"a":1}`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := agentloopJsonStringDiffSuppress("data_filter", tc.oldValue, tc.newValue, nil); got != tc.want {
				t.Errorf("agentloopJsonStringDiffSuppress(_, %q, %q, _) = %v, want %v", tc.oldValue, tc.newValue, got, tc.want)
			}
		})
	}
}
