package agentloop_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/agentloop"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Agentloop Pipeline. >>> Resource test cases, automatically generated.
// Case Pipeline全生命周期 11086
func TestAccAliCloudAgentloopPipeline_basic11086(t *testing.T) {
	resourceId := "alicloud_agentloop_pipeline.default"
	var v map[string]interface{}
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAgentloopPipelineMap11086)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return agentloop.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePipeline")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccpipeline%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAgentloopPipelineBasicDependence11086)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"agent_space":   "${alicloud_agentloop_agent_space.default.agent_space}",
					"pipeline_name": name,
					"description":   "terraform自动化测试描述",
					"pipeline": []map[string]interface{}{
						{
							"nodes": []map[string]interface{}{
								{
									"id":   "node1",
									"type": "where",
									"parameters": map[string]interface{}{
										"filter": "input IS NOT NULL",
									},
								},
							},
						},
					},
					"source": []map[string]interface{}{
						{
							"type": "dataset",
							"input_fields": []map[string]interface{}{
								{
									"name": "input",
									"type": "text",
								},
							},
							"dataset": []map[string]interface{}{
								{
									"dataset": "${alicloud_agentloop_dataset.default.dataset_name}",
									"filter":  "status = 'pending'",
								},
							},
						},
					},
					"sink": []map[string]interface{}{
						{
							"type": "dataset",
							"dataset": []map[string]interface{}{
								{
									"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
									"dataset":     "${alicloud_agentloop_dataset.default2.dataset_name}",
								},
							},
						},
					},
					"execute_policy": []map[string]interface{}{
						{
							"mode": "RunOnce",
							"run_once": []map[string]interface{}{
								{
									"from_time": "1735660800",
									"to_time":   "1735747200",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"agent_space":      name + "-as",
						"pipeline_name":    name,
						"description":      "terraform自动化测试描述",
						"pipeline.#":       "1",
						"source.#":         "1",
						"sink.#":           "1",
						"execute_policy.#": "1",
					}),
				),
			},
			{
				// description, pipeline and source are all mutable. Updating them
				// together exercises the real PUT path without replacing the sink
				// or execute policy, which remain immutable.
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform自动化测试描述更新",
					"pipeline": []map[string]interface{}{
						{
							"nodes": []map[string]interface{}{
								{
									"id":   "node1",
									"type": "where",
									"parameters": map[string]interface{}{
										"filter": "input != ''",
									},
								},
							},
						},
					},
					"source": []map[string]interface{}{
						{
							"type": "dataset",
							"input_fields": []map[string]interface{}{
								{
									"name": "input",
									"type": "text",
								},
							},
							"dataset": []map[string]interface{}{
								{
									"dataset": "${alicloud_agentloop_dataset.default.dataset_name}",
									"filter":  "status = 'ready'",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                          "terraform自动化测试描述更新",
						"pipeline.0.nodes.0.parameters.filter": "input != ''",
						"source.0.dataset.0.filter":            "status = 'ready'",
						"sink.#":                               "1",
						"execute_policy.#":                     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source": []map[string]interface{}{
						{
							"type": "logstore",
							"input_fields": []map[string]interface{}{
								{
									"name": "input",
									"type": "text",
								},
							},
							"logstore": []map[string]interface{}{
								{
									"project": "${alicloud_log_project.default.project_name}",
									// Reference the index resource (instead of the
									// logstore directly) so the pipeline is created
									// only after the required SLS index exists.
									"logstore": "${alicloud_log_store_index.default.logstore}",
									// Project only the declared input field so the
									// route output schema matches the target dataset
									// schema ("SELECT *" also returns logstore
									// metadata columns and is rejected).
									"query": "* | SELECT input",
								},
							},
						},
					},
					"sink": []map[string]interface{}{
						{
							"type": "condition",
							"condition": []map[string]interface{}{
								{
									"match_mode": "all",
									"routes": []map[string]interface{}{
										{
											"id":         "routerefund",
											"expression": "* | where input IS NOT NULL",
											"sink": []map[string]interface{}{
												{
													"type": "dataset",
													"dataset": []map[string]interface{}{
														{
															"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
															"dataset":     "${alicloud_agentloop_dataset.default.dataset_name}",
														},
													},
												},
											},
										},
									},
									"default_sink": []map[string]interface{}{
										{
											"type": "dataset",
											"dataset": []map[string]interface{}{
												{
													"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
													"dataset":     "${alicloud_agentloop_dataset.default2.dataset_name}",
												},
											},
										},
									},
								},
							},
						},
					},
					"execute_policy": []map[string]interface{}{
						{
							"mode": "Scheduled",
							"scheduled": []map[string]interface{}{
								{
									"interval":  "1h",
									"from_time": "1924905600000",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source.0.dataset.0.filter": acctest.REMOVEKEY,
						"source.#":                  "1",
						"sink.#":                    "1",
						"execute_policy.#":          "1",
					}),
				),
			},
			{
				// TypeList order coverage baseline: two distinct entries for
				// pipeline.nodes, source.input_fields and sink.condition.routes.
				// The backend requires condition routes to target distinct
				// datasets (probe37), so the routes aim at default4/default5,
				// the default_sink at default6 and the source reads default3.
				// probe37 also verified the API echoes all three lists in the
				// submitted order, so reordered updates converge.
				Config: testAccConfig(map[string]interface{}{
					"pipeline": []map[string]interface{}{
						{
							"nodes": []map[string]interface{}{
								{
									"id":   "node1",
									"type": "where",
									"parameters": map[string]interface{}{
										"filter": "input IS NOT NULL",
									},
								},
								{
									"id":   "node2",
									"type": "where",
									"parameters": map[string]interface{}{
										"filter": "output != ''",
									},
								},
							},
						},
					},
					"source": []map[string]interface{}{
						{
							"type": "dataset",
							"input_fields": []map[string]interface{}{
								{
									"name": "input",
									"type": "text",
								},
								{
									"name": "output",
									"type": "text",
								},
							},
							"dataset": []map[string]interface{}{
								{
									"dataset": "${alicloud_agentloop_dataset.default3.dataset_name}",
									"filter":  "status = 'pending'",
								},
							},
						},
					},
					"sink": []map[string]interface{}{
						{
							"type": "condition",
							"condition": []map[string]interface{}{
								{
									"match_mode": "all",
									"routes": []map[string]interface{}{
										{
											"id":         "routealpha",
											"expression": "* | where input IS NOT NULL",
											"sink": []map[string]interface{}{
												{
													"type": "dataset",
													"dataset": []map[string]interface{}{
														{
															"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
															"dataset":     "${alicloud_agentloop_dataset.default4.dataset_name}",
														},
													},
												},
											},
										},
										{
											"id":         "routebeta",
											"expression": "* | where output != ''",
											"sink": []map[string]interface{}{
												{
													"type": "dataset",
													"dataset": []map[string]interface{}{
														{
															"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
															"dataset":     "${alicloud_agentloop_dataset.default5.dataset_name}",
														},
													},
												},
											},
										},
									},
									"default_sink": []map[string]interface{}{
										{
											"type": "dataset",
											"dataset": []map[string]interface{}{
												{
													"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
													"dataset":     "${alicloud_agentloop_dataset.default6.dataset_name}",
												},
											},
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pipeline.0.nodes.0.parameters.filter": "input IS NOT NULL",
						"source.0.dataset.0.filter":            "status = 'pending'",
						"pipeline.#":                           "1",
						"source.#":                             "1",
						"sink.#":                               "1",
					}),
				),
			},
			{
				// TypeList order coverage (pipeline.nodes): a reorder-only plan
				// against the previous apply must produce a non-empty plan.
				Config: testAccConfig(map[string]interface{}{
					"pipeline": []map[string]interface{}{
						{
							"nodes": []map[string]interface{}{
								{
									"id":   "node2",
									"type": "where",
									"parameters": map[string]interface{}{
										"filter": "output != ''",
									},
								},
								{
									"id":   "node1",
									"type": "where",
									"parameters": map[string]interface{}{
										"filter": "input IS NOT NULL",
									},
								},
							},
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Applying the reordered nodes updates the pipeline in place and must
				// converge with an empty post-apply plan (API preserves order).
				Config: testAccConfig(map[string]interface{}{
					"pipeline": []map[string]interface{}{
						{
							"nodes": []map[string]interface{}{
								{
									"id":   "node2",
									"type": "where",
									"parameters": map[string]interface{}{
										"filter": "output != ''",
									},
								},
								{
									"id":   "node1",
									"type": "where",
									"parameters": map[string]interface{}{
										"filter": "input IS NOT NULL",
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pipeline.0.nodes.0.parameters.filter": "output != ''",
						"pipeline.#":                           "1",
					}),
				),
			},
			{
				// TypeList order coverage (source.input_fields): reorder-only plan.
				Config: testAccConfig(map[string]interface{}{
					"source": []map[string]interface{}{
						{
							"type": "dataset",
							"input_fields": []map[string]interface{}{
								{
									"name": "output",
									"type": "text",
								},
								{
									"name": "input",
									"type": "text",
								},
							},
							"dataset": []map[string]interface{}{
								{
									"dataset": "${alicloud_agentloop_dataset.default3.dataset_name}",
									"filter":  "status = 'pending'",
								},
							},
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Applying the reordered input_fields must converge as well.
				Config: testAccConfig(map[string]interface{}{
					"source": []map[string]interface{}{
						{
							"type": "dataset",
							"input_fields": []map[string]interface{}{
								{
									"name": "output",
									"type": "text",
								},
								{
									"name": "input",
									"type": "text",
								},
							},
							"dataset": []map[string]interface{}{
								{
									"dataset": "${alicloud_agentloop_dataset.default3.dataset_name}",
									"filter":  "status = 'pending'",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source.#": "1",
					}),
				),
			},
			{
				// TypeList order coverage (sink.condition.routes): reorder-only plan.
				Config: testAccConfig(map[string]interface{}{
					"sink": []map[string]interface{}{
						{
							"type": "condition",
							"condition": []map[string]interface{}{
								{
									"match_mode": "all",
									"routes": []map[string]interface{}{
										{
											"id":         "routebeta",
											"expression": "* | where output != ''",
											"sink": []map[string]interface{}{
												{
													"type": "dataset",
													"dataset": []map[string]interface{}{
														{
															"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
															"dataset":     "${alicloud_agentloop_dataset.default5.dataset_name}",
														},
													},
												},
											},
										},
										{
											"id":         "routealpha",
											"expression": "* | where input IS NOT NULL",
											"sink": []map[string]interface{}{
												{
													"type": "dataset",
													"dataset": []map[string]interface{}{
														{
															"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
															"dataset":     "${alicloud_agentloop_dataset.default4.dataset_name}",
														},
													},
												},
											},
										},
									},
									"default_sink": []map[string]interface{}{
										{
											"type": "dataset",
											"dataset": []map[string]interface{}{
												{
													"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
													"dataset":     "${alicloud_agentloop_dataset.default6.dataset_name}",
												},
											},
										},
									},
								},
							},
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Applying the reordered routes must converge as well. The final
				// state keeps all three lists reversed and is fully verified by
				// the import step below (ImportStateVerifyIgnore is empty).
				Config: testAccConfig(map[string]interface{}{
					"sink": []map[string]interface{}{
						{
							"type": "condition",
							"condition": []map[string]interface{}{
								{
									"match_mode": "all",
									"routes": []map[string]interface{}{
										{
											"id":         "routebeta",
											"expression": "* | where output != ''",
											"sink": []map[string]interface{}{
												{
													"type": "dataset",
													"dataset": []map[string]interface{}{
														{
															"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
															"dataset":     "${alicloud_agentloop_dataset.default5.dataset_name}",
														},
													},
												},
											},
										},
										{
											"id":         "routealpha",
											"expression": "* | where input IS NOT NULL",
											"sink": []map[string]interface{}{
												{
													"type": "dataset",
													"dataset": []map[string]interface{}{
														{
															"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
															"dataset":     "${alicloud_agentloop_dataset.default4.dataset_name}",
														},
													},
												},
											},
										},
									},
									"default_sink": []map[string]interface{}{
										{
											"type": "dataset",
											"dataset": []map[string]interface{}{
												{
													"agent_space": "${alicloud_agentloop_agent_space.default.agent_space}",
													"dataset":     "${alicloud_agentloop_dataset.default6.dataset_name}",
												},
											},
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sink.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudAgentloopPipelineMap11086 = map[string]string{
	"create_time": acctest.CHECKSET,
	"update_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func AlicloudAgentloopPipelineBasicDependence11086(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_agentloop_agent_space" "default" {
  agent_space = "${var.name}-as"
}

resource "alicloud_agentloop_dataset" "default" {
  agent_space  = alicloud_agentloop_agent_space.default.agent_space
  dataset_name = "${var.name}ds1"
  schema = {
    input = "{\"chn\":true,\"type\":\"text\"}"
  }
}

resource "alicloud_agentloop_dataset" "default2" {
  agent_space  = alicloud_agentloop_agent_space.default.agent_space
  dataset_name = "${var.name}ds2"
  schema = {
    input = "{\"chn\":true,\"type\":\"text\"}"
  }
}

# default3-default6 back the TypeList order coverage steps: the source reads
# default3, the two condition routes target default4/default5 (the backend
# requires distinct route targets), and the default_sink targets default6.
# Their schemas carry both input and output so the two-field input_fields and
# the route outputs match.
resource "alicloud_agentloop_dataset" "default3" {
  agent_space  = alicloud_agentloop_agent_space.default.agent_space
  dataset_name = "${var.name}ds3"
  schema = {
    input  = "{\"chn\":true,\"type\":\"text\"}"
    output = "{\"chn\":false,\"type\":\"text\"}"
  }
}

resource "alicloud_agentloop_dataset" "default4" {
  agent_space  = alicloud_agentloop_agent_space.default.agent_space
  dataset_name = "${var.name}ds4"
  schema = {
    input  = "{\"chn\":true,\"type\":\"text\"}"
    output = "{\"chn\":false,\"type\":\"text\"}"
  }
}

resource "alicloud_agentloop_dataset" "default5" {
  agent_space  = alicloud_agentloop_agent_space.default.agent_space
  dataset_name = "${var.name}ds5"
  schema = {
    input  = "{\"chn\":true,\"type\":\"text\"}"
    output = "{\"chn\":false,\"type\":\"text\"}"
  }
}

resource "alicloud_agentloop_dataset" "default6" {
  agent_space  = alicloud_agentloop_agent_space.default.agent_space
  dataset_name = "${var.name}ds6"
  schema = {
    input  = "{\"chn\":true,\"type\":\"text\"}"
    output = "{\"chn\":false,\"type\":\"text\"}"
  }
}

resource "alicloud_log_project" "default" {
  project_name = "${var.name}-log"
}

resource "alicloud_log_store" "default" {
  project_name  = alicloud_log_project.default.project_name
  logstore_name = "${var.name}-store"
}

# The pipeline logstore source runs its query against SLS, which requires an
# index on the logstore (otherwise the backend fails with a 500 SLS
# IndexConfigNotExist error). The SPL route expressions reference the "input"
# field, so a field index for it must exist as well (otherwise pipeline
# planning fails with "field not found in complete input schema").
resource "alicloud_log_store_index" "default" {
  project  = alicloud_log_project.default.project_name
  logstore = alicloud_log_store.default.logstore_name
  full_text {
    token = " #$^*\r\n\t"
  }
  field_search {
    name  = "input"
    type  = "text"
    token = " #$^*\r\n\t"
  }
}`, name)
}

// Test Agentloop Pipeline. <<< Resource test cases, automatically generated.
