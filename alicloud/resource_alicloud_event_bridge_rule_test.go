package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudEventBridgeRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudEventBridgeRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventbridgeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventbridgerule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEventBridgeRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EventBridgeSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"event_bus_name": "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"rule_name":      "${var.name}",
					"filter_pattern": `{\"source\":[\"crmabc.newsletter\"],\"type\":[\"UserSignUp\", \"UserLogin\"]}`,
					"targets": []map[string]interface{}{
						{
							"endpoint":  "${local.mns_endpoint_a}",
							"target_id": "tf-test1",
							"type":      "acs.mns.queue",
							"param_list": []map[string]interface{}{
								{
									"form":         "CONSTANT",
									"resource_key": "queue",
									"value":        "tf-testaccEbRule",
								},
								{
									"form":         "ORIGINAL",
									"resource_key": "Body",
								},
								{
									"form":         "CONSTANT",
									"resource_key": "IsBase64Encode",
									"value":        "true",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name":      name,
						"event_bus_name": name,
						"filter_pattern": "{\"source\":[\"crmabc.newsletter\"],\"type\":[\"UserSignUp\", \"UserLogin\"]}",
						"targets.#":      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"targets": []map[string]interface{}{
						{
							"endpoint":  "${local.mns_endpoint_b}",
							"target_id": "tf-test1",
							"type":      "acs.mns.queue",
							"param_list": []map[string]interface{}{
								{
									"form":         "CONSTANT",
									"resource_key": "queue",
									"value":        "tf-testaccEbRule",
								},
								{
									"form":         "JSONPATH",
									"resource_key": "Body",
									"value":        "$.data.name",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"targets.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"targets": []map[string]interface{}{
						{
							"endpoint":  "${local.mns_endpoint_b}",
							"target_id": "tf-test1",
							"type":      "acs.mns.queue",
							"param_list": []map[string]interface{}{
								{
									"form":         "CONSTANT",
									"resource_key": "queue",
									"value":        "tf-testaccEbRule",
								},
								{
									"form":         "CONSTANT",
									"resource_key": "Body",
									"value":        "tf-testAcc",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"targets.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"targets": []map[string]interface{}{
						{
							"endpoint":  "${local.mns_endpoint_b}",
							"target_id": "tf-test1",
							"type":      "acs.mns.queue",
							"param_list": []map[string]interface{}{
								{
									"form":         "CONSTANT",
									"resource_key": "queue",
									"value":        "tf-testaccEbRule",
								},
								{
									"form":         "TEMPLATE",
									"resource_key": "Body",
									"template":     "This is $${v1}",
									"value":        `{\n \"v1\":\"$.source\" \n}`,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"targets.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"filter_pattern": `{\"source\":[\"crmabc.newsletter\"],\"type\":[\"UserSignUp\"]}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"filter_pattern": "{\"source\":[\"crmabc.newsletter\"],\"type\":[\"UserSignUp\"]}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "DISABLE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "DISABLE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "ENABLE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "ENABLE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":    name,
					"filter_pattern": `{\"source\":[\"crmabc.newsletter\"],\"type\":[\"UserSignUp\", \"UserLogin\"]}`,
					"targets": []map[string]interface{}{
						{
							"endpoint":  "${local.mns_endpoint_a}",
							"target_id": "tf-test1",
							"type":      "acs.mns.queue",
							"param_list": []map[string]interface{}{
								{
									"form":         "CONSTANT",
									"resource_key": "queue",
									"value":        "tf-testaccEbRule",
								},
								{
									"form":         "ORIGINAL",
									"resource_key": "Body",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":    name,
						"filter_pattern": "{\"source\":[\"crmabc.newsletter\"],\"type\":[\"UserSignUp\", \"UserLogin\"]}",
						"targets.#":      "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"targets": []map[string]interface{}{
						{
							"endpoint":            "http://test.com",
							"target_id":           "tf-test1",
							"type":                "http",
							"push_retry_strategy": "BACKOFF_RETRY",
							"dead_letter_queue": []map[string]interface{}{
								{
									"arn": "acs:mns:" + "${data.alicloud_regions.default.regions.0.id}" + ":" + "${data.alicloud_account.default.id}" + ":/queues/rule-deadletterqueue",
								},
							},
							"param_list": []map[string]interface{}{
								{
									"form":         "CONSTANT",
									"resource_key": "url",
									"value":        "http://test.com",
								},
								{
									"form":         "ORIGINAL",
									"resource_key": "Body",
								},
								{
									"form":         "CONSTANT",
									"resource_key": "Network",
									"value":        "PublicNetwork",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"targets.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"targets": []map[string]interface{}{
						{
							"endpoint":            "http://tftest.com",
							"type":                "http",
							"target_id":           "tf-test1",
							"push_retry_strategy": "EXPONENTIAL_DECAY_RETRY",
							"dead_letter_queue": []map[string]interface{}{
								{
									"arn": "acs:mq:" + "${data.alicloud_regions.default.regions.0.id}" + ":" + "${data.alicloud_account.default.id}" + ":/instances/myinstance/topic/mytopic",
								},
							},
							"param_list": []map[string]interface{}{
								{
									"form":         "CONSTANT",
									"resource_key": "url",
									"value":        "http://tftest.com",
								},
								{
									"form":         "ORIGINAL",
									"resource_key": "Body",
								},
								{
									"form":         "CONSTANT",
									"resource_key": "Network",
									"value":        "PublicNetwork",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"targets.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"targets": []map[string]interface{}{
						{
							"endpoint":            "${local.fnf_endpoint}",
							"type":                "acs.fnf",
							"target_id":           "tf-test1",
							"push_retry_strategy": "BACKOFF_RETRY",
							"dead_letter_queue": []map[string]interface{}{
								{
									"arn": "acs:mq:" + "${data.alicloud_regions.default.regions.0.id}" + ":" + "${data.alicloud_account.default.id}" + ":/instances/myinstance/topic/mytopic",
								},
							},
							"param_list": []map[string]interface{}{
								{
									"form":         "CONSTANT",
									"resource_key": "FlowName",
									"value":        "demoFlow",
								},
								{
									"form":         "JSONPATH",
									"resource_key": "Input",
									"value":        "$.data.name",
								},
								{
									"form":         "CONSTANT",
									"resource_key": "RoleName",
									"value":        "roleToEB",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"targets.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudEventBridgeRuleMap0 = map[string]string{
	"event_bus_name": CHECKSET,
	"rule_name":      CHECKSET,
	"targets.#":      "1",
	"description":    "",
	"status":         CHECKSET,
	"filter_pattern": "",
}

func AlicloudEventBridgeRuleBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%[1]s"
	}

	data "alicloud_account" "default" {
	}

	data "alicloud_regions" "default" {
  		current = true
	}

	resource "alicloud_event_bridge_event_bus" "default" {
  		event_bus_name = var.name
	}

	resource "alicloud_mns_queue" "queue1" {
  		name = var.name
	}

	resource "alicloud_mns_queue" "queue2" {
  		name = format("%%schange", var.name)
	}

	locals {
  		mns_endpoint_a = format("acs:mns:%[2]s:%%s:queues/%%s", data.alicloud_account.default.id, alicloud_mns_queue.queue1.name)
  		mns_endpoint_b = format("acs:mns:%[2]s:%%s:queues/%%s", data.alicloud_account.default.id, alicloud_mns_queue.queue2.name)
  		fnf_endpoint   = format("acs:fnf:%[2]s:%%s:flow/$${flow}", data.alicloud_account.default.id)
	}
`, name, defaultRegionToTest)
}
