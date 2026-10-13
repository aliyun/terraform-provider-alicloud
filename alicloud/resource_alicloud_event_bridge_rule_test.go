package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAliCloudEventBridgeRule_regionPreCheck(t *testing.T) {
	t.Setenv("ALICLOUD_ACCESS_KEY", "test-access-key")
	t.Setenv("ALICLOUD_SECRET_KEY", "test-secret-key")
	t.Setenv("ALICLOUD_REGION", "")

	region := testAccEventBridgeRuleRegion(t)
	if region == "" {
		t.Fatal("testAccEventBridgeRuleRegion returned an empty region")
	}
	if got := os.Getenv("ALICLOUD_REGION"); region != got {
		t.Fatalf("testAccEventBridgeRuleRegion = %q, ALICLOUD_REGION = %q", region, got)
	}
}

func testAccEventBridgeRuleRegion(t *testing.T) string {
	testAccPreCheckWithRegions(t, true, connectivity.EventBridgeSupportRegions)
	return os.Getenv("ALICLOUD_REGION")
}

func TestAccAliCloudEventBridgeRule_basic0(t *testing.T) {
	var v map[string]interface{}
	region := testAccEventBridgeRuleRegion(t)
	resourceId := "alicloud_event_bridge_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventbridgeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventbridgerule%d", region, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return AliCloudEventBridgeRuleBasicDependence0(name, region)
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"event_bus_name": "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"rule_name":      name,
					"filter_pattern": `{\"source\":[\"crmabc.newsletter\"],\"type\":[\"UserSignUp\", \"UserLogin\"]}`,
					"targets": []map[string]interface{}{
						{
							"target_id": name,
							"type":      "acs.mns.queue",
							"endpoint":  "${local.mns_endpoint}",
							"param_list": []map[string]interface{}{
								{
									"resource_key": "Body",
									"form":         "ORIGINAL",
								},
								{
									"resource_key": "queue",
									"form":         "CONSTANT",
									"value":        name,
								},
								{
									"resource_key": "IsBase64Encode",
									"form":         "CONSTANT",
									"value":        "true",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name":      CHECKSET,
						"event_bus_name": name,
						"filter_pattern": "{\"source\":[\"crmabc.newsletter\"],\"type\":[\"UserSignUp\", \"UserLogin\"]}",
						"targets.#":      "1",
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
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
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
					"targets": []map[string]interface{}{
						{
							"target_id":           name,
							"type":                "acs.fnf",
							"endpoint":            "${local.fnf_endpoint}",
							"push_retry_strategy": "BACKOFF_RETRY",
							"param_list": []map[string]interface{}{
								{
									"resource_key": "Input",
									"form":         "JSONPATH",
									"value":        "$.data.name",
								},
								{
									"resource_key": "FlowName",
									"form":         "CONSTANT",
									"value":        "demoFlow",
								},
								{
									"resource_key": "RoleName",
									"form":         "CONSTANT",
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
				Config: testAccConfig(map[string]interface{}{
					"targets": []map[string]interface{}{
						{
							"target_id":           name,
							"type":                "http",
							"endpoint":            "http://www.aliyun.com",
							"push_retry_strategy": "EXPONENTIAL_DECAY_RETRY",
							"dead_letter_queue": []map[string]interface{}{
								{
									"arn": "${local.mns_endpoint}",
								},
							},
							"param_list": []map[string]interface{}{
								{
									"resource_key": "Body",
									"form":         "TEMPLATE",
									"template":     "This is $${v1}",
									"value":        `{\n \"v1\":\"$.source\" \n}`,
								},
								{
									"resource_key": "url",
									"form":         "CONSTANT",
									"value":        "http://www.aliyun.com",
								},
								{
									"resource_key": "Network",
									"form":         "CONSTANT",
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudEventBridgeRule_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	region := testAccEventBridgeRuleRegion(t)
	resourceId := "alicloud_event_bridge_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventbridgeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventbridgerule%d", region, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return AliCloudEventBridgeRuleBasicDependence0(name, region)
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"event_bus_name": "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"rule_name":      name,
					"filter_pattern": `{\"source\":[\"crmabc.newsletter\"],\"type\":[\"UserSignUp\", \"UserLogin\"]}`,
					"description":    name,
					"status":         "ENABLE",
					"targets": []map[string]interface{}{
						{
							"target_id":           name,
							"type":                "http",
							"endpoint":            "http://www.aliyun.com",
							"push_retry_strategy": "EXPONENTIAL_DECAY_RETRY",
							"dead_letter_queue": []map[string]interface{}{
								{
									"arn": "${local.mns_endpoint}",
								},
							},
							"param_list": []map[string]interface{}{
								{
									"resource_key": "Body",
									"form":         "TEMPLATE",
									"template":     "This is $${v1}",
									"value":        `{\n \"v1\":\"$.source\" \n}`,
								},
								{
									"resource_key": "url",
									"form":         "CONSTANT",
									"value":        "http://www.aliyun.com",
								},
								{
									"resource_key": "Network",
									"form":         "CONSTANT",
									"value":        "PublicNetwork",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name":      CHECKSET,
						"event_bus_name": name,
						"filter_pattern": "{\"source\":[\"crmabc.newsletter\"],\"type\":[\"UserSignUp\", \"UserLogin\"]}",
						"description":    name,
						"status":         "ENABLE",
						"targets.#":      "1",
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

var AliCloudEventBridgeRuleMap0 = map[string]string{
	"status": CHECKSET,
}

func AliCloudEventBridgeRuleBasicDependence0(name, region string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%[1]s"
	}

	data "alicloud_account" "default" {
	}

	resource "alicloud_event_bridge_event_bus" "default" {
  		event_bus_name = var.name
	}

	resource "alicloud_mns_queue" "default" {
  		name = var.name
	}

	locals {
  		mns_endpoint = format("acs:mns:%[2]s:%%s:queues/%%s", data.alicloud_account.default.id, alicloud_mns_queue.default.name)
  		fnf_endpoint   = format("acs:fnf:%[2]s:%%s:flow/$${flow}", data.alicloud_account.default.id)
	}
`, name, region)
}

func TestAccAliCloudEventBridgeRule_apacheKafkaSelf(t *testing.T) {
	var v map[string]interface{}
	region := testAccEventBridgeRuleRegion(t)
	resourceId := "alicloud_event_bridge_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventbridgeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventbridgerule%d", region, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return AliCloudEventBridgeRuleBasicDependence0(name, region)
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"event_bus_name": "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"rule_name":      name,
					"filter_pattern": `{\"source\":[\"crmabc.newsletter\"],\"type\":[\"UserSignUp\", \"UserLogin\"]}`,
					"targets": []map[string]interface{}{
						{
							"target_id": name,
							"type":      "ApacheKafkaSelf",
							"endpoint":  "192.168.0.1:9092",
							"param_list": []map[string]interface{}{
								{
									"resource_key": "Bootstraps",
									"form":         "CONSTANT",
									"value":        "192.168.0.1:9092",
								},
								{
									"resource_key": "Topic",
									"form":         "CONSTANT",
									"value":        name,
								},
								{
									"resource_key": "NetworkType",
									"form":         "CONSTANT",
									"value":        "PublicNetwork",
								},
								{
									"resource_key": "SecurityProtocol",
									"form":         "CONSTANT",
									"value":        "SASL_PLAINTEXT",
								},
								{
									"resource_key": "SaslMechanism",
									"form":         "CONSTANT",
									"value":        "PLAIN",
								},
								{
									"resource_key": "SaslUser",
									"form":         "CONSTANT",
									"value":        "kafkaUser",
								},
								{
									"resource_key": "SaslPassword",
									"form":         "CONSTANT",
									"value":        "kafkaPassword",
								},
								{
									"resource_key": "Key",
									"form":         "CONSTANT",
									"value":        "eventKey",
								},
								{
									"resource_key": "Value",
									"form":         "ORIGINAL",
								},
								{
									"resource_key": "Acks",
									"form":         "CONSTANT",
									"value":        "1",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name":      CHECKSET,
						"event_bus_name": name,
						"filter_pattern": "{\"source\":[\"crmabc.newsletter\"],\"type\":[\"UserSignUp\", \"UserLogin\"]}",
						"targets.#":      "1",
						"targets.0.type": "ApacheKafkaSelf",
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
