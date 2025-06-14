// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test MessageService EventRule. >>> Resource test cases, automatically generated.
// Case rule-topic 10919
func TestAccAliCloudMessageServiceEventRule_basic10919(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_message_service_event_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudMessageServiceEventRuleMap10919)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MessageServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMessageServiceEventRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmessageservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMessageServiceEventRuleBasicDependence10919)
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
					"endpoints": []map[string]interface{}{
						{
							"endpoint_type":  "topic",
							"endpoint_value": "${alicloud_message_service_queue.CreateQueue.id}",
						},
					},
					"rule_name": "${var.rule_name}",
					"event_types": []string{
						"ObjectCreated:PutObject"},
					"match_rules": [][]map[string]interface{}{
						{
							{
								"match_state": "true",
								"name":        "aaac",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoints.#":   "1",
						"rule_name":     CHECKSET,
						"event_types.#": "1",
						"match_rules.#": "1",
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

var AlicloudMessageServiceEventRuleMap10919 = map[string]string{}

func AlicloudMessageServiceEventRuleBasicDependence10919(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "queue_name" {
  default = "tf-test-topic2queue"
}

variable "rule_name" {
  default = "testRule-topic-1"
}

variable "topic_name" {
  default = "tf-test-topic2queue"
}

resource "alicloud_message_service_topic" "CreateTopic" {
  maximum_message_size = "65536"
  name                 = var.topic_name
  logging_enabled      = false
}

resource "alicloud_message_service_queue" "CreateQueue" {
  delay_seconds            = "2"
  polling_wait_seconds     = "2"
  message_retention_period = "566"
  maximum_message_size     = "1123"
  visibility_timeout       = "30"
  queue_name               = var.queue_name
  logging_enabled          = false
}

resource "alicloud_message_service_subscription" "CreateSub" {
  push_type             = "queue"
  notify_strategy       = "BACKOFF_RETRY"
  notify_content_format = "SIMPLIFIED"
  subscription_name     = "RDK-test-sub"
  filter_tag            = "important"
  topic_name            = alicloud_message_service_topic.CreateTopic.name
  endpoint              = alicloud_message_service_queue.CreateQueue.id
}


`, name)
}

// Case rule-http 10837
func TestAccAliCloudMessageServiceEventRule_basic10837(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_message_service_event_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudMessageServiceEventRuleMap10837)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MessageServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMessageServiceEventRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmessageservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMessageServiceEventRuleBasicDependence10837)
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
					"endpoints": []map[string]interface{}{
						{
							"endpoint_type":  "http",
							"endpoint_value": "http://www.baidu.com",
						},
					},
					"rule_name": "${var.rule_name}",
					"event_types": []string{
						"ObjectCreated:PutObject"},
					"match_rules": [][]map[string]interface{}{
						{
							{
								"name":        "zxc",
								"match_state": "true",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoints.#":   "1",
						"rule_name":     CHECKSET,
						"event_types.#": "1",
						"match_rules.#": "1",
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

var AlicloudMessageServiceEventRuleMap10837 = map[string]string{}

func AlicloudMessageServiceEventRuleBasicDependence10837(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "rule_name" {
  default = "testRule-http"
}


`, name)
}

// Case rule-queue-Prefix-Suffix 10917
func TestAccAliCloudMessageServiceEventRule_basic10917(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_message_service_event_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudMessageServiceEventRuleMap10917)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MessageServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMessageServiceEventRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmessageservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMessageServiceEventRuleBasicDependence10917)
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
					"endpoints": []map[string]interface{}{
						{
							"endpoint_type":  "queue",
							"endpoint_value": "${alicloud_message_service_queue.CreateQueue.id}",
						},
					},
					"rule_name": "${var.rule_name}",
					"event_types": []string{
						"ObjectCreated:PutObject"},
					"match_rules": [][]map[string]interface{}{
						{
							{
								"match_state": "true",
								"prefix":      "aaa",
								"suffix":      "ccc",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoints.#":   "1",
						"rule_name":     CHECKSET,
						"event_types.#": "1",
						"match_rules.#": "1",
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

var AlicloudMessageServiceEventRuleMap10917 = map[string]string{}

func AlicloudMessageServiceEventRuleBasicDependence10917(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "rule_name" {
  default = "testRule-queue"
}

variable "queue_name" {
  default = "tf-test-queue"
}

resource "alicloud_message_service_queue" "CreateQueue" {
  delay_seconds            = "2"
  polling_wait_seconds     = "2"
  message_retention_period = "566"
  maximum_message_size     = "1123"
  visibility_timeout       = "30"
  queue_name               = var.queue_name
  logging_enabled          = false
}


`, name)
}

// Test MessageService EventRule. <<< Resource test cases, automatically generated.
