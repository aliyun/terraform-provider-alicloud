package agentloop_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/agentloop"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Agentloop EvaluationTask. >>> Resource test cases, automatically generated.
// Case EvaluationTask全生命周期 11081
func TestAccAliCloudAgentloopEvaluationTask_basic11081(t *testing.T) {
	resourceId := "alicloud_agentloop_evaluation_task.default"
	rac := resourceAttrCheckInit(resourceId, AlicloudAgentloopEvaluationTaskMap11081, (*agentloop.Service).DescribeEvaluationTask)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccevaltask%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAgentloopEvaluationTaskBasicDependence11081)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
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
					// 首步保留启用的即时/持续策略，不配置可能异步变化的 status。
					// 仅 ID-only 比较忽略 status；Read 仍回填真实状态。
					// 后续策略禁用的重建步骤显式验证 Pending。
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
									"interval_unit":      "HOUR",
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
							"name":          "${alicloud_agentloop_evaluator.default.name}",
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
						"agent_space":                           name + "-as",
						"task_name":                             name,
						"task_mode":                             "batch",
						"data_type":                             "trace",
						"channel":                               "default",
						"description":                           "terraform自动化测试描述",
						"evaluators.#":                          "1",
						"evaluators.0.result_name":              "resulta",
						"evaluators.0.type":                     "AGENT",
						"evaluators.0.config.%":                 "1",
						"evaluators.0.variable_mapping.input":   "$.input",
						"config.%":                              "1",
						"config.config1":                        "value1",
						"tags.tag1":                             "value1",
						"run_strategies.#":                      "1",
						"run_strategies.0.continuous.0.enabled": "true",
						"run_strategies.0.continuous.0.interval_unit":  "HOUR",
						"run_strategies.0.continuous.0.interval_value": "6",
						"run_strategies.0.backfill.0.immediate":        "true",
						"run_strategies.0.backfill.0.start_time":       "1735689600000",
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
									"interval_unit":      "MINUTE",
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
							"name":          "${alicloud_agentloop_evaluator.default2.name}",
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
						"channel":                                      "console",
						"description":                                  "terraform自动化测试描述更新",
						"evaluators.#":                                 "1",
						"evaluators.0.result_name":                     "resultb",
						"evaluators.0.result_type":                     "score",
						"evaluators.0.variable_mapping.input":          "$.output",
						"config.config1":                               "value2",
						"tags.tag1":                                    "value2",
						"run_strategies.0.continuous.0.enabled":        "false",
						"run_strategies.0.continuous.0.interval_unit":  "MINUTE",
						"run_strategies.0.continuous.0.interval_value": "12",
						"run_strategies.0.backfill.0.enabled":          "false",
						"run_strategies.0.backfill.0.immediate":        "false",
						"run_strategies.0.backfill.0.start_time":       "1735862400000",
						"run_strategies.0.backfill.0.end_time":         "1735948800000",
					}),
				),
			},
			{
				// run_strategies is updated in place (no recreation); every
				// strategy stays disabled so the task remains Pending.
				Config: testAccConfig(map[string]interface{}{
					"run_strategies": []map[string]interface{}{
						{
							"continuous": []map[string]interface{}{
								{
									"enabled":            "false",
									"interval_unit":      "DAY",
									"interval_value":     "1",
									"data_delay_minutes": "15",
								},
							},
							"backfill": []map[string]interface{}{
								{
									"enabled":    "false",
									"immediate":  "false",
									"start_time": "1736035200000",
									"end_time":   "1736121600000",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"run_strategies.0.continuous.0.interval_unit":      "DAY",
						"run_strategies.0.continuous.0.interval_value":     "1",
						"run_strategies.0.continuous.0.data_delay_minutes": "15",
						"run_strategies.0.backfill.0.start_time":           "1736035200000",
						"run_strategies.0.backfill.0.end_time":             "1736121600000",
					}),
				),
			},
			{
				// TypeList order coverage baseline for evaluators: two distinct
				// entries referencing the two dependence evaluators.
				Config: testAccConfig(map[string]interface{}{
					"evaluators": []map[string]interface{}{
						{
							"name":          "${alicloud_agentloop_evaluator.default.name}",
							"result_name":   "resulta",
							"type":          "AGENT",
							"result_type":   "score",
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
							"name":          "${alicloud_agentloop_evaluator.default2.name}",
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
						"evaluators.#":                        "2",
						"evaluators.0.result_name":            "resulta",
						"evaluators.0.result_type":            "score",
						"evaluators.0.variable_mapping.input": "$.input",
						"evaluators.1.result_name":            "resultb",
						"evaluators.1.result_type":            "score",
						"evaluators.1.variable_mapping.input": "$.output",
					}),
				),
			},
			{
				// TypeList order coverage: a reorder-only plan against the
				// previous apply must produce a non-empty plan.
				Config: testAccConfig(map[string]interface{}{
					"evaluators": []map[string]interface{}{
						{
							"name":          "${alicloud_agentloop_evaluator.default2.name}",
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
							"name":          "${alicloud_agentloop_evaluator.default.name}",
							"result_name":   "resulta",
							"type":          "AGENT",
							"result_type":   "score",
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
				// empty post-apply plan: the read-back order follows the API,
				// which keeps the submitted order.
				Config: testAccConfig(map[string]interface{}{
					"evaluators": []map[string]interface{}{
						{
							"name":          "${alicloud_agentloop_evaluator.default2.name}",
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
							"name":          "${alicloud_agentloop_evaluator.default.name}",
							"result_name":   "resulta",
							"type":          "AGENT",
							"result_type":   "score",
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
						"evaluators.#":                        "2",
						"evaluators.0.result_name":            "resultb",
						"evaluators.0.result_type":            "score",
						"evaluators.0.variable_mapping.input": "$.output",
						"evaluators.1.result_name":            "resulta",
						"evaluators.1.result_type":            "score",
						"evaluators.1.variable_mapping.input": "$.input",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// config/evaluators/run_strategies are read back (backend-derived
				// keys are skipped, see the resource Read). status is ignored:
				// the task transitions Pending -> Running asynchronously, so
				// the imported state may differ from the post-apply state.
				ImportStateVerifyIgnore: []string{"status"},
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
