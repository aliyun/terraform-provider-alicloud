package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea-rpc/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_event_bridge_event_streaming",
		&resource.Sweeper{
			Name: "alicloud_event_bridge_event_streaming",
			F:    testSweepEventBridgeEventStreaming,
		})
}

func testSweepEventBridgeEventStreaming(region string) error {
	if testSweepPreCheckWithRegions(region, true, connectivity.EventBridgeSupportRegions) {
		log.Printf("[INFO] Skipping Event Bridge Event Streaming unsupported region: %s", region)
		return nil
	}

	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	aliyunClient := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	action := "ListEventStreaming"
	request := map[string]interface{}{}
	request["MaxResults"] = PageSizeLarge

	var response map[string]interface{}
	conn, err := aliyunClient.NewEventbridgeClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
		return nil
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.Data.EventStreamings", response)
		// 确认 ：formatInt(response["Total"]) != 0 是否保留
		if formatInt(response["Total"]) != 0 && err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Data.EventStreamings", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["EventStreamingName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Event Bridge Event Streaming: %s", item["EventStreamingName"].(string))
				continue
			}
			action := "DeleteEventStreaming"
			request := map[string]interface{}{
				"EventStreamingName": item["EventStreamingName"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Event Bridge Event Streaming (%s): %s", item["EventStreamingName"].(string), err)
			}
			log.Printf("[INFO] Delete Event Bridge Event Streaming success: %s ", item["EventStreamingName"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudEventBridgeEventStreaming_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_streaming.default"
	checkoutSupportedRegions(t, true, connectivity.EventBridgeSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudEventBridgeEventStreamingMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventbridgeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventStreaming")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%seventbridgeeventstreaming%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEventBridgeEventStreamingBasicDependence0)
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
					"event_streaming_name": "${var.name}",
					"description":          "${var.name}",
					"source": []map[string]interface{}{
						{
							"source_mns_parameters": []map[string]interface{}{
								{
									"queue_name":       "test",
									"is_base64_decode": "true",
								},
							},
						},
					},
					"sink": []map[string]interface{}{
						{
							"sink_mns_parameters": []map[string]interface{}{
								{
									"queue_name": []map[string]interface{}{
										{
											"value": "test",
											"form":  "CONSTANT",
										},
									},
									"body": []map[string]interface{}{
										{
											"value": "$.data",
											"form":  "JSONPATH",
										},
									},
									"is_base64_encode": []map[string]interface{}{
										{
											"value": "true",
											"form":  "CONSTANT",
										},
									},
								},
							},
						},
					},
					"run_options": []map[string]interface{}{
						{
							"errors_tolerance": "ALL",
							"retry_strategy": []map[string]interface{}{
								{
									"push_retry_strategy": "BACKOFF_RETRY",
								},
							},
						},
					},
					"filter_pattern": "{}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_streaming_name":             name,
						"description":                      name,
						"source.#":                         "1",
						"source.0.source_mns_parameters.#": "1",
						"source.0.source_mns_parameters.0.queue_name":       "test",
						"source.0.source_mns_parameters.0.is_base64_decode": "true",
						"sink.#":                       "1",
						"sink.0.sink_mns_parameters.#": "1",
						"sink.0.sink_mns_parameters.0.queue_name.#":             "1",
						"sink.0.sink_mns_parameters.0.queue_name.0.value":       "test",
						"sink.0.sink_mns_parameters.0.queue_name.0.form":        "CONSTANT",
						"sink.0.sink_mns_parameters.0.body.#":                   "1",
						"sink.0.sink_mns_parameters.0.body.0.value":             "$.data",
						"sink.0.sink_mns_parameters.0.body.0.form":              "JSONPATH",
						"sink.0.sink_mns_parameters.0.is_base64_encode.#":       "1",
						"sink.0.sink_mns_parameters.0.is_base64_encode.0.value": "true",
						"sink.0.sink_mns_parameters.0.is_base64_encode.0.form":  "CONSTANT",
						"run_options.#":                                      "1",
						"run_options.0.errors_tolerance":                     "ALL",
						"run_options.0.retry_strategy.#":                     "1",
						"run_options.0.retry_strategy.0.push_retry_strategy": "BACKOFF_RETRY",
						"filter_pattern":                                     "{}",
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

var AlicloudEventBridgeEventStreamingMap0 = map[string]string{
	"event_streaming_name": CHECKSET,
	"run_options.#":        CHECKSET,
	"sink.#":               CHECKSET,
	"source.#":             CHECKSET,
	"status":               CHECKSET,
}

func AlicloudEventBridgeEventStreamingBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

func TestAccAlicloudEventBridgeEventStreaming_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_streaming.default"
	checkoutSupportedRegions(t, true, connectivity.EventBridgeSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudEventBridgeEventStreamingMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventbridgeService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventStreaming")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%seventbridgeeventstreaming%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEventBridgeEventStreamingBasicDependence0)
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
					"event_streaming_name": "${var.name}",
					"description":          "${var.name}",
					"filter_pattern":       "{}",
					"source": []map[string]interface{}{
						{
							"source_mns_parameters": []map[string]interface{}{
								{
									"queue_name":       "test",
									"is_base64_decode": "true",
								},
							},
						},
					},
					"sink": []map[string]interface{}{
						{
							"sink_mns_parameters": []map[string]interface{}{
								{
									"queue_name": []map[string]interface{}{
										{
											"value": "test",
											"form":  "CONSTANT",
										},
									},
									"body": []map[string]interface{}{
										{
											"value": "$.data",
											"form":  "JSONPATH",
										},
									},
									"is_base64_encode": []map[string]interface{}{
										{
											"value": "true",
											"form":  "CONSTANT",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_streaming_name": name,
						"description":          name,
						"filter_pattern":       "{}",
						"source.#":             "1",
						"sink.#":               "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"filter_pattern": "{}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"filter_pattern": "{}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source": []map[string]interface{}{
						{
							"source_mns_parameters": []map[string]interface{}{
								{
									"queue_name":       "test1",
									"is_base64_decode": "true",
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
				Config: testAccConfig(map[string]interface{}{
					"sink": []map[string]interface{}{
						{
							"sink_mns_parameters": []map[string]interface{}{
								{
									"queue_name": []map[string]interface{}{
										{
											"value": "test1",
											"form":  "CONSTANT",
										},
									},
									"body": []map[string]interface{}{
										{
											"value": "$.data",
											"form":  "JSONPATH",
										},
									},
									"is_base64_encode": []map[string]interface{}{
										{
											"value": "true",
											"form":  "CONSTANT",
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
				Config: testAccConfig(map[string]interface{}{
					"run_options": []map[string]interface{}{
						{
							"errors_tolerance": "ALL",
							"retry_strategy": []map[string]interface{}{
								{
									"push_retry_strategy": "BACKOFF_RETRY",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"run_options.#": "1",
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

func TestUnitAlicloudEventBridgeEventStreaming(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	dInit, _ := schema.InternalMap(p["alicloud_event_bridge_event_streaming"].Schema).Data(nil, nil)
	dExisted, _ := schema.InternalMap(p["alicloud_event_bridge_event_streaming"].Schema).Data(nil, nil)
	dInit.MarkNewResource()
	attributes := map[string]interface{}{
		"event_streaming_name": "CreateEventBridgeEventStreamingValue",
		"description":          "CreateEventBridgeEventStreamingValue",
		"source": []map[string]interface{}{
			{
				"source_mns_parameters": []map[string]interface{}{
					{
						"region_id":        "CreateEventBridgeEventStreamingValue",
						"queue_name":       "CreateEventBridgeEventStreamingValue",
						"is_base64_decode": true,
					},
				},
				"source_rabbit_mq_parameters": []map[string]interface{}{
					{
						"region_id":         "CreateEventBridgeEventStreamingValue",
						"instance_id":       "CreateEventBridgeEventStreamingValue",
						"virtual_host_name": "CreateEventBridgeEventStreamingValue",
						"queue_name":        "CreateEventBridgeEventStreamingValue",
					},
				},
				"source_kafka_parameters": []map[string]interface{}{
					{
						"region_id":         "CreateEventBridgeEventStreamingValue",
						"instance_id":       "CreateEventBridgeEventStreamingValue",
						"topic":             "CreateEventBridgeEventStreamingValue",
						"consumer_group":    "CreateEventBridgeEventStreamingValue",
						"offset_reset":      "CreateEventBridgeEventStreamingValue",
						"network":           "CreateEventBridgeEventStreamingValue",
						"vpc_id":            "CreateEventBridgeEventStreamingValue",
						"vswitch_ids":       "CreateEventBridgeEventStreamingValue",
						"security_group_id": "CreateEventBridgeEventStreamingValue",
					},
				},
				"source_rocket_mq_parameters": []map[string]interface{}{
					{
						"region_id":   "CreateEventBridgeEventStreamingValue",
						"instance_id": "CreateEventBridgeEventStreamingValue",
						"topic":       "CreateEventBridgeEventStreamingValue",
						"tag":         "CreateEventBridgeEventStreamingValue",
						"offset":      "CreateEventBridgeEventStreamingValue",
						"group_id":    "CreateEventBridgeEventStreamingValue",
						"timestamp":   1636597951964,
					},
				},
				"source_mqtt_parameters": []map[string]interface{}{
					{
						"region_id":   "CreateEventBridgeEventStreamingValue",
						"instance_id": "CreateEventBridgeEventStreamingValue",
						"topic":       "CreateEventBridgeEventStreamingValue",
					},
				},
				"source_dts_parameters": []map[string]interface{}{
					{
						"task_id":          "CreateEventBridgeEventStreamingValue",
						"broker_url":       "CreateEventBridgeEventStreamingValue",
						"topic":            "CreateEventBridgeEventStreamingValue",
						"sid":              "CreateEventBridgeEventStreamingValue",
						"username":         "CreateEventBridgeEventStreamingValue",
						"password":         "CreateEventBridgeEventStreamingValue",
						"init_check_point": "CreateEventBridgeEventStreamingValue",
					},
				},
				"source_sls_parameters": []map[string]interface{}{
					{
						"project":          "CreateEventBridgeEventStreamingValue",
						"log_store":        "CreateEventBridgeEventStreamingValue",
						"consume_position": "CreateEventBridgeEventStreamingValue",
						"role_name":        "CreateEventBridgeEventStreamingValue",
					},
				},
			},
		},
		"sink": []map[string]interface{}{
			{
				"sink_mns_parameters": []map[string]interface{}{
					{
						"queue_name": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"body": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"is_base64_encode": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
					},
				},
				"sink_rabbit_mq_parameters": []map[string]interface{}{
					{
						"instance_id": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"virtual_host_name": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"target_type": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"exchange": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"routing_key": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"queue_name": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"body": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"message_id": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"properties": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
					},
				},
				"sink_kafka_parameters": []map[string]interface{}{
					{
						"instance_id": []interface{}{
							map[string]interface{}{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"topic": []interface{}{
							map[string]interface{}{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"acks": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"key": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"value": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"sasl_user": []interface{}{
							map[string]interface{}{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
					},
				},
				"sink_fc_parameters": []map[string]interface{}{
					{
						"service_name": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"function_name": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"qualifier": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"invocation_type": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"body": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
					},
				},
				"sink_rocket_mq_parameters": []map[string]interface{}{
					{
						"instance_id": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"topic": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"body": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"properties": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"keys": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"tags": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
					},
				},
				"sink_sls_parameters": []map[string]interface{}{
					{
						"project": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"log_store": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"topic": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"body": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"role_name": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
					},
				},
			},
		},
		"run_options": []map[string]interface{}{
			{
				"maximum_tasks":    10,
				"errors_tolerance": "CreateEventBridgeEventStreamingValue",
				"retry_strategy": []map[string]interface{}{
					{
						"push_retry_strategy":          "CreateEventBridgeEventStreamingValue",
						"maximum_event_age_in_seconds": 10,
						"maximum_retry_attempts":       10,
					},
				},
				"dead_letter_queue": []map[string]interface{}{
					{
						"arn": "CreateEventBridgeEventStreamingValue",
					},
				},
				"batch_window": []map[string]interface{}{
					{
						"count_based_window": 10,
						"time_based_window":  10,
					},
				},
			},
		},
		"filter_pattern": "CreateEventBridgeEventStreamingValue",
	}
	for key, value := range attributes {
		err := dInit.Set(key, value)
		assert.Nil(t, err)
		err = dExisted.Set(key, value)
		assert.Nil(t, err)
		if err != nil {
			log.Printf("[ERROR] the field %s setting error", key)
		}
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		// GetEventStreaming
		"Data": map[string]interface{}{
			"Status":        "READY",
			"FilterPattern": "CreateEventBridgeEventStreamingValue",
			"Description":   "CreateEventBridgeEventStreamingValue",
			"Sink": map[string]interface{}{
				"SinkMNSParameters": map[string]interface{}{
					"IsBase64Encode": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"QueueName": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Body": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
				},
				"SinkFcParameters": map[string]interface{}{
					"Body": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"FunctionName": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"InvocationType": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Qualifier": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"ServiceName": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
				},
				"SinkKafkaParameters": map[string]interface{}{
					"Acks": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"InstanceId": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Key": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"SaslUser": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Topic": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Value": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
				},
				"SinkRabbitMQParameters": map[string]interface{}{
					"Exchange": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"InstanceId": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Body": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"MessageId": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Properties": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"QueueName": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"RoutingKey": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"TargetType": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"VirtualHostName": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
				},
				"SinkRocketMQParameters": map[string]interface{}{
					"Body": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"InstanceId": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Keys": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Properties": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Topic": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Tags": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
				},
				"SinkSLSParameters": map[string]interface{}{
					"Body": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"LogStore": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Project": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"RoleName": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Topic": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
				},
			},
			"EventStreamingName": "CreateEventBridgeEventStreamingValue",
			"Source": map[string]interface{}{
				"SourceMNSParameters": map[string]interface{}{
					"IsBase64Decode": true,
					"RegionId":       "CreateEventBridgeEventStreamingValue",
					"QueueName":      "CreateEventBridgeEventStreamingValue",
				},
				"SourceDTSParameters": map[string]interface{}{
					"BrokerUrl":      "CreateEventBridgeEventStreamingValue",
					"InitCheckPoint": "CreateEventBridgeEventStreamingValue",
					"Password":       "CreateEventBridgeEventStreamingValue",
					"Sid":            "CreateEventBridgeEventStreamingValue",
					"TaskId":         "CreateEventBridgeEventStreamingValue",
					"Topic":          "CreateEventBridgeEventStreamingValue",
					"Username":       "CreateEventBridgeEventStreamingValue",
				},
				"SourceKafkaParameters": map[string]interface{}{
					"ConsumerGroup":   "CreateEventBridgeEventStreamingValue",
					"InstanceId":      "CreateEventBridgeEventStreamingValue",
					"Network":         "CreateEventBridgeEventStreamingValue",
					"OffsetReset":     "CreateEventBridgeEventStreamingValue",
					"RegionId":        "CreateEventBridgeEventStreamingValue",
					"SecurityGroupId": "CreateEventBridgeEventStreamingValue",
					"Topic":           "CreateEventBridgeEventStreamingValue",
					"VSwitchIds":      "CreateEventBridgeEventStreamingValue",
					"VpcId":           "CreateEventBridgeEventStreamingValue",
				},
				"SourceMQTTParameters": map[string]interface{}{
					"InstanceId": "CreateEventBridgeEventStreamingValue",
					"RegionId":   "CreateEventBridgeEventStreamingValue",
					"Topic":      "CreateEventBridgeEventStreamingValue",
				},
				"SourceRabbitMQParameters": map[string]interface{}{
					"InstanceId":      "CreateEventBridgeEventStreamingValue",
					"RegionId":        "CreateEventBridgeEventStreamingValue",
					"QueueName":       "CreateEventBridgeEventStreamingValue",
					"VirtualHostName": "CreateEventBridgeEventStreamingValue",
				},
				"SourceRocketMQParameters": map[string]interface{}{
					"GroupID":    "CreateEventBridgeEventStreamingValue",
					"InstanceId": "CreateEventBridgeEventStreamingValue",
					"Offset":     "CreateEventBridgeEventStreamingValue",
					"RegionId":   "CreateEventBridgeEventStreamingValue",
					"Tag":        "CreateEventBridgeEventStreamingValue",
					"Timestamp":  1636597951964,
					"Topic":      "CreateEventBridgeEventStreamingValue",
				},
				"SourceSLSParameters": map[string]interface{}{
					"ConsumePosition": "CreateEventBridgeEventStreamingValue",
					"LogStore":        "CreateEventBridgeEventStreamingValue",
					"Project":         "CreateEventBridgeEventStreamingValue",
					"RoleName":        "CreateEventBridgeEventStreamingValue",
				},
			},
			"RunOptions": map[string]interface{}{
				"ErrorsTolerance": "CreateEventBridgeEventStreamingValue",
				"MaximumTasks":    10,
				"RetryStrategy": map[string]interface{}{
					"PushRetryStrategy":        "CreateEventBridgeEventStreamingValue",
					"MaximumRetryAttempts":     10,
					"MaximumEventAgeInSeconds": 10,
				},
				"BatchWindow": map[string]interface{}{
					"CountBasedWindow": 10,
					"TimeBasedWindow":  10,
				},
				"DeadLetterQueue": map[string]interface{}{
					"Arn": "CreateEventBridgeEventStreamingValue",
				},
			},
		},
		"Code":    "Success",
		"Success": true,
	}
	CreateMockResponse := map[string]interface{}{
		"Success": true,
	}
	ReadMockResponseDiff := map[string]interface{}{}
	failedResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, &tea.SDKError{
			Code:       String(errorCode),
			Data:       String(errorCode),
			Message:    String(errorCode),
			StatusCode: tea.Int(400),
		}
	}
	notFoundResponseMock := func(errorCode string) (map[string]interface{}, error) {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_event_bridge_event_streaming", errorCode))
	}
	successResponseMock := func(operationMockResponse map[string]interface{}) (map[string]interface{}, error) {
		if len(operationMockResponse) > 0 {
			mapMerge(ReadMockResponse, operationMockResponse)
		}
		return ReadMockResponse, nil
	}

	// Create
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEventbridgeClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudEventBridgeEventStreamingCreate(dInit, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	ReadMockResponseDiff = map[string]interface{}{}
	errorCodes := []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1 // a counter used to cover retry scenario; the same below
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "CreateEventStreaming" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						successResponseMock(ReadMockResponseDiff)
						return CreateMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEventBridgeEventStreamingCreate(dInit, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_event_bridge_event_streaming"].Schema).Data(dInit.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dInit.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Update
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEventbridgeClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudEventBridgeEventStreamingUpdate(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	// UpdateEventStreaming
	attributesDiff := map[string]interface{}{
		"description": "UpdateEventBridgeEventStreamingValue",
		"source": []map[string]interface{}{
			{
				"source_mns_parameters": []map[string]interface{}{
					{
						"region_id":        "CreateEventBridgeEventStreamingValue",
						"queue_name":       "UpdateEventBridgeEventStreamingValue",
						"is_base64_decode": true,
					},
				},
				"source_rabbit_mq_parameters": []map[string]interface{}{
					{
						"region_id":         "CreateEventBridgeEventStreamingValue",
						"instance_id":       "CreateEventBridgeEventStreamingValue",
						"virtual_host_name": "CreateEventBridgeEventStreamingValue",
						"queue_name":        "CreateEventBridgeEventStreamingValue",
					},
				},
				"source_kafka_parameters": []map[string]interface{}{
					{
						"region_id":         "CreateEventBridgeEventStreamingValue",
						"instance_id":       "CreateEventBridgeEventStreamingValue",
						"topic":             "CreateEventBridgeEventStreamingValue",
						"consumer_group":    "CreateEventBridgeEventStreamingValue",
						"offset_reset":      "CreateEventBridgeEventStreamingValue",
						"network":           "CreateEventBridgeEventStreamingValue",
						"vpc_id":            "CreateEventBridgeEventStreamingValue",
						"vswitch_ids":       "CreateEventBridgeEventStreamingValue",
						"security_group_id": "CreateEventBridgeEventStreamingValue",
					},
				},
				"source_rocket_mq_parameters": []map[string]interface{}{
					{
						"region_id":   "CreateEventBridgeEventStreamingValue",
						"instance_id": "CreateEventBridgeEventStreamingValue",
						"topic":       "CreateEventBridgeEventStreamingValue",
						"tag":         "CreateEventBridgeEventStreamingValue",
						"offset":      "CreateEventBridgeEventStreamingValue",
						"group_id":    "CreateEventBridgeEventStreamingValue",
						"timestamp":   1636597951964,
					},
				},
				"source_mqtt_parameters": []map[string]interface{}{
					{
						"region_id":   "CreateEventBridgeEventStreamingValue",
						"instance_id": "CreateEventBridgeEventStreamingValue",
						"topic":       "CreateEventBridgeEventStreamingValue",
					},
				},
				"source_dts_parameters": []map[string]interface{}{
					{
						"task_id":          "CreateEventBridgeEventStreamingValue",
						"broker_url":       "CreateEventBridgeEventStreamingValue",
						"topic":            "CreateEventBridgeEventStreamingValue",
						"sid":              "CreateEventBridgeEventStreamingValue",
						"username":         "CreateEventBridgeEventStreamingValue",
						"password":         "CreateEventBridgeEventStreamingValue",
						"init_check_point": "CreateEventBridgeEventStreamingValue",
					},
				},
				"source_sls_parameters": []map[string]interface{}{
					{
						"project":          "CreateEventBridgeEventStreamingValue",
						"log_store":        "CreateEventBridgeEventStreamingValue",
						"consume_position": "CreateEventBridgeEventStreamingValue",
						"role_name":        "CreateEventBridgeEventStreamingValue",
					},
				},
			},
		},
		"sink": []map[string]interface{}{
			{
				"sink_mns_parameters": []map[string]interface{}{
					{
						"queue_name": []map[string]interface{}{
							{
								"value":    "UpdateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"body": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"is_base64_encode": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
					},
				},
				"sink_rabbit_mq_parameters": []map[string]interface{}{
					{
						"instance_id": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"virtual_host_name": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"target_type": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"exchange": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"routing_key": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"queue_name": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"body": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"message_id": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"properties": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
					},
				},
				"sink_kafka_parameters": []map[string]interface{}{
					{
						"instance_id": []interface{}{
							map[string]interface{}{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"topic": []interface{}{
							map[string]interface{}{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"acks": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"key": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"value": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"sasl_user": []interface{}{
							map[string]interface{}{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
					},
				},
				"sink_fc_parameters": []map[string]interface{}{
					{
						"service_name": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"function_name": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"qualifier": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"invocation_type": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"body": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
					},
				},
				"sink_rocket_mq_parameters": []map[string]interface{}{
					{
						"instance_id": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"topic": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"body": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"properties": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"keys": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"tags": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
					},
				},
				"sink_sls_parameters": []map[string]interface{}{
					{
						"project": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"log_store": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"topic": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"body": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
						"role_name": []map[string]interface{}{
							{
								"value":    "CreateEventBridgeEventStreamingValue",
								"form":     "CreateEventBridgeEventStreamingValue",
								"template": "CreateEventBridgeEventStreamingValue",
							},
						},
					},
				},
			},
		},
		"run_options": []map[string]interface{}{
			{
				"maximum_tasks":    10,
				"errors_tolerance": "UpdateEventBridgeEventStreamingValue",
				"retry_strategy": []map[string]interface{}{
					{
						"push_retry_strategy":          "CreateEventBridgeEventStreamingValue",
						"maximum_event_age_in_seconds": 10,
						"maximum_retry_attempts":       10,
					},
				},
				"dead_letter_queue": []map[string]interface{}{
					{
						"arn": "CreateEventBridgeEventStreamingValue",
					},
				},
				"batch_window": []map[string]interface{}{
					{
						"count_based_window": 10,
						"time_based_window":  10,
					},
				},
			},
		},
		"filter_pattern": "UpdateEventBridgeEventStreamingValue",
	}
	diff, err := newInstanceDiff("alicloud_event_bridge_event_streaming", attributes, attributesDiff, dInit.State())
	if err != nil {
		t.Error(err)
	}
	dExisted, _ = schema.InternalMap(p["alicloud_event_bridge_event_streaming"].Schema).Data(dInit.State(), diff)
	ReadMockResponseDiff = map[string]interface{}{
		// GetEventStreaming Response
		"Data": map[string]interface{}{
			"FilterPattern": "UpdateEventBridgeEventStreamingValue",
			"Description":   "UpdateEventBridgeEventStreamingValue",
			"Sink": map[string]interface{}{
				"SinkMNSParameters": map[string]interface{}{
					"IsBase64Encode": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"QueueName": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "UpdateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Body": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
				},
				"SinkFcParameters": map[string]interface{}{
					"Body": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"FunctionName": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"InvocationType": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Qualifier": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"ServiceName": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
				},
				"SinkKafkaParameters": map[string]interface{}{
					"Acks": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"InstanceId": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Key": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"SaslUser": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Topic": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Value": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
				},
				"SinkRabbitMQParameters": map[string]interface{}{
					"Exchange": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"InstanceId": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Body": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"MessageId": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Properties": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"QueueName": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"RoutingKey": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"TargetType": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"VirtualHostName": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
				},
				"SinkRocketMQParameters": map[string]interface{}{
					"Body": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"InstanceId": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Keys": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Properties": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Topic": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Tags": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
				},
				"SinkSLSParameters": map[string]interface{}{
					"Body": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"LogStore": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Project": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"RoleName": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
					"Topic": map[string]interface{}{
						"Form":     "CreateEventBridgeEventStreamingValue",
						"Value":    "CreateEventBridgeEventStreamingValue",
						"Template": "CreateEventBridgeEventStreamingValue",
					},
				},
			},
			"Source": map[string]interface{}{
				"SourceMNSParameters": map[string]interface{}{
					"IsBase64Decode": true,
					"RegionId":       "CreateEventBridgeEventStreamingValue",
					"QueueName":      "UpdateEventBridgeEventStreamingValue",
				},
				"SourceDTSParameters": map[string]interface{}{
					"BrokerUrl":      "CreateEventBridgeEventStreamingValue",
					"InitCheckPoint": "CreateEventBridgeEventStreamingValue",
					"Password":       "CreateEventBridgeEventStreamingValue",
					"Sid":            "CreateEventBridgeEventStreamingValue",
					"TaskId":         "CreateEventBridgeEventStreamingValue",
					"Topic":          "CreateEventBridgeEventStreamingValue",
					"Username":       "CreateEventBridgeEventStreamingValue",
				},
				"SourceKafkaParameters": map[string]interface{}{
					"ConsumerGroup":   "CreateEventBridgeEventStreamingValue",
					"InstanceId":      "CreateEventBridgeEventStreamingValue",
					"Network":         "CreateEventBridgeEventStreamingValue",
					"OffsetReset":     "CreateEventBridgeEventStreamingValue",
					"RegionId":        "CreateEventBridgeEventStreamingValue",
					"SecurityGroupId": "CreateEventBridgeEventStreamingValue",
					"Topic":           "CreateEventBridgeEventStreamingValue",
					"VSwitchIds":      "CreateEventBridgeEventStreamingValue",
					"VpcId":           "CreateEventBridgeEventStreamingValue",
				},
				"SourceMQTTParameters": map[string]interface{}{
					"InstanceId": "CreateEventBridgeEventStreamingValue",
					"RegionId":   "CreateEventBridgeEventStreamingValue",
					"Topic":      "CreateEventBridgeEventStreamingValue",
				},
				"SourceRabbitMQParameters": map[string]interface{}{
					"InstanceId":      "CreateEventBridgeEventStreamingValue",
					"RegionId":        "CreateEventBridgeEventStreamingValue",
					"QueueName":       "CreateEventBridgeEventStreamingValue",
					"VirtualHostName": "CreateEventBridgeEventStreamingValue",
				},
				"SourceRocketMQParameters": map[string]interface{}{
					"GroupID":    "CreateEventBridgeEventStreamingValue",
					"InstanceId": "CreateEventBridgeEventStreamingValue",
					"Offset":     "CreateEventBridgeEventStreamingValue",
					"RegionId":   "CreateEventBridgeEventStreamingValue",
					"Tag":        "CreateEventBridgeEventStreamingValue",
					"Timestamp":  1636597951964,
					"Topic":      "CreateEventBridgeEventStreamingValue",
				},
				"SourceSLSParameters": map[string]interface{}{
					"ConsumePosition": "CreateEventBridgeEventStreamingValue",
					"LogStore":        "CreateEventBridgeEventStreamingValue",
					"Project":         "CreateEventBridgeEventStreamingValue",
					"RoleName":        "CreateEventBridgeEventStreamingValue",
				},
			},
			"RunOptions": map[string]interface{}{
				"ErrorsTolerance": "UpdateEventBridgeEventStreamingValue",
				"MaximumTasks":    10,
				"RetryStrategy": map[string]interface{}{
					"PushRetryStrategy":        "CreateEventBridgeEventStreamingValue",
					"MaximumRetryAttempts":     10,
					"MaximumEventAgeInSeconds": 10,
				},
				"BatchWindow": map[string]interface{}{
					"CountBasedWindow": 10,
					"TimeBasedWindow":  10,
				},
				"DeadLetterQueue": map[string]interface{}{
					"Arn": "CreateEventBridgeEventStreamingValue",
				},
			},
		},
	}
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "UpdateEventStreaming" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if retryIndex >= len(errorCodes)-1 {
						return successResponseMock(ReadMockResponseDiff)
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEventBridgeEventStreamingUpdate(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		default:
			assert.Nil(t, err)
			dCompare, _ := schema.InternalMap(p["alicloud_event_bridge_event_streaming"].Schema).Data(dExisted.State(), nil)
			for key, value := range attributes {
				dCompare.Set(key, value)
			}
			assert.Equal(t, dCompare.State().Attributes, dExisted.State().Attributes)
		}
		if retryIndex >= len(errorCodes)-1 {
			break
		}
	}

	// Read
	errorCodes = []string{"NonRetryableError", "Throttling", "nil", "{}"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "GetEventStreaming" {
				switch errorCode {
				case "{}":
					return notFoundResponseMock(errorCode)
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEventBridgeEventStreamingRead(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "{}":
			assert.Nil(t, err)
		}
	}

	// Delete
	patches = gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewEventbridgeClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
		return nil, &tea.SDKError{
			Code:    String("loadEndpoint error"),
			Data:    String("loadEndpoint error"),
			Message: String("loadEndpoint error"),
		}
	})
	err = resourceAlicloudEventBridgeEventStreamingDelete(dExisted, rawClient)
	patches.Reset()
	assert.NotNil(t, err)
	errorCodes = []string{"NonRetryableError", "Throttling", "nil"}
	for index, errorCode := range errorCodes {
		retryIndex := index - 1
		patches = gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, action *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if *action == "DeleteEventStreaming" {
				switch errorCode {
				case "NonRetryableError":
					return failedResponseMock(errorCode)
				default:
					retryIndex++
					if errorCodes[retryIndex] == "nil" {
						ReadMockResponse = map[string]interface{}{}
						return ReadMockResponse, nil
					}
					return failedResponseMock(errorCodes[retryIndex])
				}
			}
			return ReadMockResponse, nil
		})
		err := resourceAlicloudEventBridgeEventStreamingDelete(dExisted, rawClient)
		patches.Reset()
		switch errorCode {
		case "NonRetryableError":
			assert.NotNil(t, err)
		case "nil":
			assert.Nil(t, err)
		}
	}
}
