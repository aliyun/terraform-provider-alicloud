package mnsopen_test

import (
	"fmt"
	"testing"

	tfacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	mnsopen "github.com/aliyun/terraform-provider-alicloud/alicloud/service/mnsopen"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudMessageServiceEventRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_message_service_event_rule.default"
	ra := tfacctest.ResourceAttrInit(resourceId, AlicloudMessageServiceEventRuleMap0)
	rc := tfacctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mnsopen.NewMessageServiceServiceV2(tfacctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMessageServiceEventRule")
	rac := tfacctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmessageservice%d", rand)
	testAccConfig := tfacctest.ResourceTestAccConfig(resourceId, name, AlicloudMessageServiceEventRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tfacctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
			tfacctest.PreCheck(t)
		},
		ProviderFactories: tfacctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint": []map[string]interface{}{
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
								"name":        "acs:oss:eu-central-1:1511928242963727:accccx",
								"prefix":      "",
								"suffix":      "",
							},
						},
						{
							{
								"match_state": "true",
								"prefix":      "acs:oss:eu-central-1:1511928242963727:home",
								"suffix":      ".avi",
								"name":        "",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint.#":    "1",
						"rule_name":     tfacctest.CHECKSET,
						"event_types.#": "1",
						"match_rules.#": "2",
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

var AlicloudMessageServiceEventRuleMap0 = map[string]string{}

func AlicloudMessageServiceEventRuleBasicDependence0(name string) string {
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

// Test MessageService EventRule. >>> Resource test cases, automatically generated.
// Case rule-topic 10919
func TestAccAliCloudMessageServiceEventRule_basic10919(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_message_service_event_rule.default"
	ra := tfacctest.ResourceAttrInit(resourceId, AlicloudMessageServiceEventRuleMap10919)
	rc := tfacctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mnsopen.NewMessageServiceServiceV2(tfacctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMessageServiceEventRule")
	rac := tfacctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmessageservice%d", rand)
	testAccConfig := tfacctest.ResourceTestAccConfig(resourceId, name, AlicloudMessageServiceEventRuleBasicDependence10919)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tfacctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			tfacctest.PreCheck(t)
		},
		ProviderFactories: tfacctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"endpoint": []map[string]interface{}{
						{
							"endpoint_type":  "topic",
							"endpoint_value": "${alicloud_message_service_subscription.CreateSub.topic_name}",
						},
					},
					"rule_name": "${var.rule_name}",
					"event_types": []string{
						"ObjectCreated:PutObject"},
					"match_rules": [][]map[string]interface{}{
						{
							{
								"match_state": "true",
								"name":        "acs:oss:cn-hangzhou:1511928242963727:accccx",
								"prefix":      "",
								"suffix":      "",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"endpoint.#":    "1",
						"rule_name":     tfacctest.CHECKSET,
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
  max_message_size = "65536"
  topic_name                 = var.topic_name
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
  topic_name            = alicloud_message_service_topic.CreateTopic.topic_name
  endpoint              = format("acs:mns:cn-hangzhou:1511928242963727:/queues/%%s", alicloud_message_service_queue.CreateQueue.id)
}


`, name)
}
