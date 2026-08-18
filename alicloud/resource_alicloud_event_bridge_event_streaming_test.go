package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudEventBridgeEventStreaming_regionPreCheck(t *testing.T) {
	t.Setenv("ALICLOUD_ACCESS_KEY", "test-access-key")
	t.Setenv("ALICLOUD_SECRET_KEY", "test-secret-key")
	t.Setenv("ALICLOUD_REGION", "")

	region := testAccEventBridgeEventStreamingRegion(t)
	if region == "" {
		t.Fatal("testAccEventBridgeEventStreamingRegion returned an empty region")
	}
	if got := os.Getenv("ALICLOUD_REGION"); region != got {
		t.Fatalf("testAccEventBridgeEventStreamingRegion = %q, ALICLOUD_REGION = %q", region, got)
	}
}

func testAccEventBridgeEventStreamingRegion(t *testing.T) string {
	testAccPreCheckWithRegions(t, true, connectivity.EventBridgeSupportRegions)
	return os.Getenv("ALICLOUD_REGION")
}

// Case 1: full lifecycle of an event streaming, covering source/filter_pattern/sink/
// run_options/transforms updates, status transitions and tags management.
func TestAccAliCloudEventBridgeEventStreaming_basic0(t *testing.T) {
	var v map[string]interface{}
	region := testAccEventBridgeEventStreamingRegion(t)
	resourceId := "alicloud_event_bridge_event_streaming.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventStreamingMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventStreaming")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventstreaming%d", region, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return AliCloudEventBridgeEventStreamingBasicDependence0(name, region)
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
					"event_streaming_name": name,
					"source":               "${local.source_json}",
					"filter_pattern":       `{}`,
					"sink":                 "${local.sink_json}",
					"run_options":          "${local.run_options_json}",
					"transforms":           "${local.transforms_json}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_streaming_name": name,
						"filter_pattern":       `{}`,
						"source":               CHECKSET,
						"sink":                 CHECKSET,
						"run_options":          CHECKSET,
						"transforms":           CHECKSET,
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
					"source":         "${local.source_json_update}",
					"filter_pattern": `{\"source\":[\"acs.mns\"]}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source":         CHECKSET,
						"filter_pattern": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sink":        "${local.sink_json_update}",
					"run_options": "${local.run_options_json_update}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sink":        CHECKSET,
						"run_options": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transforms": "${local.transforms_json_update}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transforms": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "RUNNING",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "RUNNING",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "PAUSED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "PAUSED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"created": "true",
						"purpose": "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.created": "true",
						"tags.purpose": "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"created": "false",
						"purpose": "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.created": "false",
						"tags.purpose": "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"purpose": "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "1",
						"tags.created": REMOVEKEY,
						"tags.purpose": "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.purpose": REMOVEKEY,
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

// Case 2: create an event streaming with all attributes set, including the
// expected status RUNNING and tags, then verify import.
func TestAccAliCloudEventBridgeEventStreaming_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	region := testAccEventBridgeEventStreamingRegion(t)
	resourceId := "alicloud_event_bridge_event_streaming.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventStreamingMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventStreaming")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seventstreaming%d", region, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return AliCloudEventBridgeEventStreamingBasicDependence0(name, region)
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
					"event_streaming_name": name,
					"description":          name,
					"source":               "${local.source_json}",
					"filter_pattern":       `{}`,
					"sink":                 "${local.sink_json}",
					"run_options":          "${local.run_options_json}",
					"transforms":           "${local.transforms_json}",
					"status":               "RUNNING",
					"tags": map[string]string{
						"created": "true",
						"purpose": "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_streaming_name": name,
						"description":          name,
						"filter_pattern":       `{}`,
						"source":               CHECKSET,
						"sink":                 CHECKSET,
						"run_options":          CHECKSET,
						"transforms":           CHECKSET,
						"status":               "RUNNING",
						"tags.%":               "2",
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

var AliCloudEventBridgeEventStreamingMap0 = map[string]string{
	"status": CHECKSET,
}

func AliCloudEventBridgeEventStreamingBasicDependence0(name, region string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%[1]s"
	}

	data "alicloud_account" "default" {
	}

	resource "alicloud_message_service_queue" "source" {
		queue_name = "${var.name}-source"
	}

	resource "alicloud_message_service_queue" "sink" {
		queue_name = "${var.name}-sink"
	}

	locals {
		source_json = format("{\"SourceMNSParameters\":{\"RegionId\":\"%[2]s\",\"QueueName\":\"%%s\",\"IsBase64Decode\":true}}", alicloud_message_service_queue.source.queue_name)
		source_json_update = format("{\"SourceMNSParameters\":{\"RegionId\":\"%[2]s\",\"QueueName\":\"%%s\",\"IsBase64Decode\":false}}", alicloud_message_service_queue.source.queue_name)
		sink_json = format("{\"SinkMNSParameters\":{\"QueueName\":{\"Value\":\"%%s\",\"Form\":\"CONSTANT\"},\"Body\":{\"Value\":\"$.data\",\"Form\":\"JSONPATH\"},\"IsBase64Encode\":{\"Value\":\"true\",\"Form\":\"CONSTANT\"}}}", alicloud_message_service_queue.sink.queue_name)
		sink_json_update = format("{\"SinkMNSParameters\":{\"QueueName\":{\"Value\":\"%%s\",\"Form\":\"CONSTANT\"},\"Body\":{\"Value\":\"$.data.messageBody\",\"Form\":\"JSONPATH\"},\"IsBase64Encode\":{\"Value\":\"false\",\"Form\":\"CONSTANT\"}}}", alicloud_message_service_queue.sink.queue_name)
		run_options_json = "{\"ErrorsTolerance\":\"ALL\",\"MaximumTasks\":1,\"BatchWindow\":{\"CountBasedWindow\":1,\"TimeBasedWindow\":10},\"RetryStrategy\":{\"MaximumEventAgeInSeconds\":600,\"MaximumRetryAttempts\":3,\"PushRetryStrategy\":\"BACKOFF_RETRY\"}}"
		run_options_json_update = "{\"ErrorsTolerance\":\"NONE\",\"MaximumTasks\":2,\"BatchWindow\":{\"CountBasedWindow\":10,\"TimeBasedWindow\":15},\"RetryStrategy\":{\"MaximumEventAgeInSeconds\":1200,\"MaximumRetryAttempts\":2,\"PushRetryStrategy\":\"EXPONENTIAL_DECAY_RETRY\"}}"
		transforms_json = format("[{\"Arn\":\"acs:fc:%[2]s:%%s:functions/%[1]s\"}]", data.alicloud_account.default.id)
		transforms_json_update = format("[{\"Arn\":\"acs:fc:%[2]s:%%s:functions/%[1]s-second\"}]", data.alicloud_account.default.id)
	}
`, name, region)
}
