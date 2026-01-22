package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test EventBridge EventSource. >>> Resource test cases, automatically generated.
// Case test-event-http 11700
func TestAccAliCloudEventBridgeEventSourceV2_basic11700(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_source_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventSourceV2Map11700)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceventbridge%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeEventSourceV2BasicDependence11700)
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
					"event_bus_name":    "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":    CHECKSET,
						"event_source_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_http_event_parameters": []map[string]interface{}{
						{
							"type":            "HTTP",
							"security_config": "referer",
							"method":          []string{"GET", "POST", "DELETE"},
							"referer":         []string{"www.aliyun.com", "www.alicloud.com", "www.taobao.com"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_http_event_parameters.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_http_event_parameters": []map[string]interface{}{
						{
							"type":            "HTTP",
							"security_config": "ip",
							"method":          []string{"GET", "DELETE", "POST", "PUT"},
							"ip":              []string{"192.168.1.1/12", "192.168.1.2/12", "192.168.1.3/12", "192.168.1.4/12"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_http_event_parameters.#": "1",
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
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"linked_external_source": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

func TestAccAliCloudEventBridgeEventSourceV2_basic11700_referer_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_source_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventSourceV2Map11700)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceventbridge%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeEventSourceV2BasicDependence11700)
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
					"event_bus_name":    "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name": name,
					"source_http_event_parameters": []map[string]interface{}{
						{
							"type":            "HTTP",
							"security_config": "referer",
							"method":          []string{"GET", "POST", "DELETE"},
							"referer":         []string{"www.aliyun.com", "www.alicloud.com", "www.taobao.com"},
						},
					},
					"description":            name,
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":                 CHECKSET,
						"event_source_name":              name,
						"source_http_event_parameters.#": "1",
						"description":                    name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

func TestAccAliCloudEventBridgeEventSourceV2_basic11700_ip_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_source_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventSourceV2Map11700)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceventbridge%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeEventSourceV2BasicDependence11700)
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
					"event_bus_name":    "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name": name,
					"source_http_event_parameters": []map[string]interface{}{
						{
							"type":            "HTTP",
							"security_config": "ip",
							"method":          []string{"GET", "DELETE", "POST", "PUT"},
							"ip":              []string{"192.168.1.1/12", "192.168.1.2/12", "192.168.1.3/12", "192.168.1.4/12"},
						},
					},
					"description":            name,
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":                 CHECKSET,
						"event_source_name":              name,
						"source_http_event_parameters.#": "1",
						"description":                    name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

var AliCloudEventBridgeEventSourceV2Map11700 = map[string]string{}

func AliCloudEventBridgeEventSourceV2BasicDependence11700(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_event_bridge_event_bus" "default" {
  		event_bus_name = var.name
	}
`, name)
}

// Case test-event-kafka 11692
func TestAccAliCloudEventBridgeEventSourceV2_basic11692(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_source_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventSourceV2Map11692)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceventbridge%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeEventSourceV2BasicDependence11692)
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
					"event_bus_name":    "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":    CHECKSET,
						"event_source_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_kafka_parameters": []map[string]interface{}{
						{
							"instance_id":       "${data.alicloud_alikafka_instances.default.instances.0.id}",
							"consumer_group":    "${data.alicloud_alikafka_consumer_groups.default.groups.0.id}",
							"topic":             "${data.alicloud_alikafka_topics.default.topics.0.id}",
							"offset_reset":      "latest",
							"region_id":         "cn-hangzhou",
							"network":           "Default",
							"security_group_id": "${data.alicloud_alikafka_instances.default.instances.0.security_group}",
							"vpc_id":            "${data.alicloud_alikafka_instances.default.instances.0.vpc_id}",
							"vswitch_ids":       "${data.alicloud_alikafka_instances.default.instances.0.vswitch_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_kafka_parameters.#": "1",
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
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"linked_external_source": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

func TestAccAliCloudEventBridgeEventSourceV2_basic11692_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_source_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventSourceV2Map11692)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceventbridge%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeEventSourceV2BasicDependence11692)
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
					"event_bus_name":    "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name": name,
					"source_kafka_parameters": []map[string]interface{}{
						{
							"instance_id":       "${data.alicloud_alikafka_instances.default.instances.0.id}",
							"consumer_group":    "${data.alicloud_alikafka_consumer_groups.default.groups.0.id}",
							"topic":             "${data.alicloud_alikafka_topics.default.topics.0.id}",
							"offset_reset":      "latest",
							"region_id":         "cn-hangzhou",
							"network":           "Default",
							"security_group_id": "${data.alicloud_alikafka_instances.default.instances.0.security_group}",
							"vpc_id":            "${data.alicloud_alikafka_instances.default.instances.0.vpc_id}",
							"vswitch_ids":       "${data.alicloud_alikafka_instances.default.instances.0.vswitch_id}",
						},
					},
					"description":            name,
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":            CHECKSET,
						"event_source_name":         name,
						"source_kafka_parameters.#": "1",
						"description":               name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

var AliCloudEventBridgeEventSourceV2Map11692 = map[string]string{}

func AliCloudEventBridgeEventSourceV2BasicDependence11692(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_alikafka_instances" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_alikafka_consumer_groups" "default" {
  		instance_id = data.alicloud_alikafka_instances.default.instances.0.id
	}

	data "alicloud_alikafka_topics" "default" {
  		instance_id = data.alicloud_alikafka_instances.default.instances.0.id
	}

	resource "alicloud_event_bridge_event_bus" "default" {
  		event_bus_name = var.name
	}
`, name)
}

// Case test-event-source-mns 11678
func TestAccAliCloudEventBridgeEventSourceV2_basic11678(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_source_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventSourceV2Map11678)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceventbridge%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeEventSourceV2BasicDependence11678)
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
					"event_bus_name":    "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":    CHECKSET,
						"event_source_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_mns_parameters": []map[string]interface{}{
						{
							"region_id":        "cn-hangzhou",
							"queue_name":       "${alicloud_message_service_queue.default.queue_name}",
							"is_base64_decode": "true",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_mns_parameters.#": "1",
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
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"linked_external_source": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

func TestAccAliCloudEventBridgeEventSourceV2_basic11678_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_source_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventSourceV2Map11678)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceventbridge%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeEventSourceV2BasicDependence11678)
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
					"event_bus_name":    "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name": name,
					"source_mns_parameters": []map[string]interface{}{
						{
							"region_id":        "cn-hangzhou",
							"queue_name":       "${alicloud_message_service_queue.default.queue_name}",
							"is_base64_decode": "true",
						},
					},
					"description":            name,
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":          CHECKSET,
						"event_source_name":       name,
						"source_mns_parameters.#": "1",
						"description":             name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

var AliCloudEventBridgeEventSourceV2Map11678 = map[string]string{}

func AliCloudEventBridgeEventSourceV2BasicDependence11678(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_event_bridge_event_bus" "default" {
  		event_bus_name = var.name
	}

	resource "alicloud_message_service_queue" "default" {
  		queue_name = var.name
	}
`, name)
}

// Case test-event-rabbitmq 11691
func TestAccAliCloudEventBridgeEventSourceV2_basic11691(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_source_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventSourceV2Map11691)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceventbridge%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeEventSourceV2BasicDependence11691)
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
					"event_bus_name":    "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":    CHECKSET,
						"event_source_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_rabbit_mq_parameters": []map[string]interface{}{
						{
							"region_id":         "cn-hangzhou",
							"instance_id":       "${alicloud_amqp_queue.default.instance_id}",
							"virtual_host_name": "${alicloud_amqp_queue.default.virtual_host_name}",
							"queue_name":        "${alicloud_amqp_queue.default.queue_name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_rabbit_mq_parameters.#": "1",
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
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"linked_external_source": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

func TestAccAliCloudEventBridgeEventSourceV2_basic11691_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_source_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventSourceV2Map11691)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceventbridge%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeEventSourceV2BasicDependence11691)
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
					"event_bus_name":    "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name": name,
					"source_rabbit_mq_parameters": []map[string]interface{}{
						{
							"region_id":         "cn-hangzhou",
							"instance_id":       "${alicloud_amqp_queue.default.instance_id}",
							"virtual_host_name": "${alicloud_amqp_queue.default.virtual_host_name}",
							"queue_name":        "${alicloud_amqp_queue.default.queue_name}",
						},
					},
					"description":            name,
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":                CHECKSET,
						"event_source_name":             name,
						"source_rabbit_mq_parameters.#": "1",
						"description":                   name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

var AliCloudEventBridgeEventSourceV2Map11691 = map[string]string{}

func AliCloudEventBridgeEventSourceV2BasicDependence11691(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_amqp_instances" "default" {
		status = "SERVING"
	}

	resource "alicloud_event_bridge_event_bus" "default" {
  		event_bus_name = var.name
	}

	resource "alicloud_amqp_virtual_host" "default" {
		instance_id       = data.alicloud_amqp_instances.default.ids.0
		virtual_host_name = var.name
	}

	resource "alicloud_amqp_queue" "default" {
  		instance_id       = alicloud_amqp_virtual_host.default.instance_id
  		virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
  		queue_name        = var.name
  		auto_delete_state = true
	}
`, name)
}

// Case test-event-rocketmq 11693
func TestAccAliCloudEventBridgeEventSourceV2_basic11693(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_source_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventSourceV2Map11693)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceventbridge%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeEventSourceV2BasicDependence11693)
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
					"event_bus_name":    "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":    CHECKSET,
						"event_source_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_rocketmq_parameters": []map[string]interface{}{
						{
							"region_id":                  "cn-hangzhou",
							"instance_id":                "test",
							"topic":                      "test",
							"tag":                        "test",
							"offset":                     "test",
							"timestamp":                  "1",
							"group_id":                   "test",
							"instance_type":              "test",
							"instance_vpc_id":            "test",
							"instance_vswitch_ids":       "test",
							"instance_security_group_id": "test",
							"auth_type":                  "test",
							"instance_endpoint":          "test",
							"instance_username":          "test",
							"instance_password":          "test",
							"instance_network":           "test",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_rocketmq_parameters.#": "1",
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
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"linked_external_source": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

func TestAccAliCloudEventBridgeEventSourceV2_basic11693_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_source_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventSourceV2Map11693)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceventbridge%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeEventSourceV2BasicDependence11693)
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
					"event_bus_name":    "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name": name,
					"source_rocketmq_parameters": []map[string]interface{}{
						{
							"region_id":                  "cn-hangzhou",
							"instance_id":                "test",
							"topic":                      "test",
							"tag":                        "test",
							"offset":                     "test",
							"timestamp":                  "1",
							"group_id":                   "test",
							"instance_type":              "test",
							"instance_vpc_id":            "test",
							"instance_vswitch_ids":       "test",
							"instance_security_group_id": "test",
							"auth_type":                  "test",
							"instance_endpoint":          "test",
							"instance_username":          "test",
							"instance_password":          "test",
							"instance_network":           "test",
						},
					},
					"description":            name,
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":               CHECKSET,
						"event_source_name":            name,
						"source_rocketmq_parameters.#": "1",
						"description":                  name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

var AliCloudEventBridgeEventSourceV2Map11693 = map[string]string{}

func AliCloudEventBridgeEventSourceV2BasicDependence11693(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_event_bridge_event_bus" "default" {
  		event_bus_name = var.name
	}
`, name)
}

// Case test-event-sls 11694
func TestAccAliCloudEventBridgeEventSourceV2_basic11694(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_source_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventSourceV2Map11694)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceventbridge%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeEventSourceV2BasicDependence11694)
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
					"event_bus_name":    "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":    CHECKSET,
						"event_source_name": name,
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
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"linked_external_source": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

func TestAccAliCloudEventBridgeEventSourceV2_basic11694_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_source_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventSourceV2Map11694)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceventbridge%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeEventSourceV2BasicDependence11694)
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
					"event_bus_name":    "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name": name,
					"source_sls_parameters": []map[string]interface{}{
						{
							"project":          "${alicloud_log_store.default.project_name}",
							"log_store":        "${alicloud_log_store.default.logstore_name}",
							"consume_position": "begin",
							"role_name":        "${alicloud_ram_role_policy_attachment.default.role_name}",
						},
					},
					"description":            name,
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":          CHECKSET,
						"event_source_name":       name,
						"source_sls_parameters.#": "1",
						"description":             name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

var AliCloudEventBridgeEventSourceV2Map11694 = map[string]string{}

func AliCloudEventBridgeEventSourceV2BasicDependence11694(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_event_bridge_event_bus" "default" {
  		event_bus_name = var.name
	}

	resource "alicloud_log_project" "default" {
  		project_name = var.name
	}

	resource "alicloud_log_store" "default" {
  		project_name  = alicloud_log_project.default.project_name
  		logstore_name = var.name
	}

	resource "alicloud_ram_role" "default" {
  		role_name                   = var.name
  		assume_role_policy_document = <<EOF
		{
			"Statement": [
				{
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": {
						"Service": [
							"eventbridge.aliyuncs.com"
						]
					}
				}
		  	],
			"Version": "1"
		}
	  	EOF
		force    = true
	}

	resource "alicloud_ram_role_policy_attachment" "default" {
  		role_name   = alicloud_ram_role.default.role_name
  		policy_name = "AliyunLogFullAccess"
  		policy_type = "System"
	}
`, name)
}

// Case test-event-Scheduled 11695
func TestAccAliCloudEventBridgeEventSourceV2_basic11695(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_source_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventSourceV2Map11695)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceventbridge%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeEventSourceV2BasicDependence11695)
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
					"event_bus_name":    "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":    CHECKSET,
						"event_source_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_scheduled_event_parameters": []map[string]interface{}{
						{
							"schedule":  "10 * * * * *",
							"time_zone": "GMT+8:00",
							"user_data": "{\\\"a\\\":\\\"b\\\"}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_scheduled_event_parameters.#": "1",
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
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"linked_external_source": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

func TestAccAliCloudEventBridgeEventSourceV2_basic11695_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_source_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventSourceV2Map11695)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceventbridge%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeEventSourceV2BasicDependence11695)
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
					"event_bus_name":    "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name": name,
					"source_scheduled_event_parameters": []map[string]interface{}{
						{
							"schedule":  "10 * * * * *",
							"time_zone": "GMT+8:00",
							"user_data": "{\\\"a\\\":\\\"b\\\"}",
						},
					},
					"description":            name,
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":                      CHECKSET,
						"event_source_name":                   name,
						"source_scheduled_event_parameters.#": "1",
						"description":                         name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

var AliCloudEventBridgeEventSourceV2Map11695 = map[string]string{}

func AliCloudEventBridgeEventSourceV2BasicDependence11695(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_event_bridge_event_bus" "default" {
  		event_bus_name = var.name
	}
`, name)
}

// Case test-event-oss 11707
func TestAccAliCloudEventBridgeEventSourceV2_basic11707(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_source_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventSourceV2Map11707)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceventbridge%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeEventSourceV2BasicDependence11707)
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
					"event_bus_name":    "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name": name,
					"source_oss_event_parameters": []map[string]interface{}{
						{
							"event_types": []string{
								"oss:ObjectCreated:AppendObject", "oss:ObjectCreated:CompleteMultipartUpload", "oss:ObjectCreated:CopyObject"},
							"sts_role_arn": "acs:ram::${data.alicloud_account.default.id}:role/${alicloud_ram_role_policy_attachment.AliyunEventBridgePutEventsPolicy.role_name}",
							"match_rules": [][]map[string]interface{}{
								{
									{
										"match_state": "true",
										"suffix":      "t",
										"prefix":      "acs:oss:cn-hangzhou:${data.alicloud_account.default.id}:${data.alicloud_oss_buckets.default.buckets.0.name}/",
										"name":        "",
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":                CHECKSET,
						"event_source_name":             name,
						"source_oss_event_parameters.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source_oss_event_parameters": []map[string]interface{}{
						{
							"event_types": []string{
								"oss:ObjectCreated:AppendObject"},
							"sts_role_arn": "acs:ram::${data.alicloud_account.default.id}:role/${alicloud_ram_role_policy_attachment.update.role_name}",
							"match_rules": [][]map[string]interface{}{
								{
									{
										"match_state": "true",
										"suffix":      "",
										"prefix":      "",
										"name":        "acs:oss:cn-hangzhou:${data.alicloud_account.default.id}:${data.alicloud_oss_buckets.default.buckets.0.name}/t1",
									},
									{
										"match_state": "true",
										"suffix":      "",
										"prefix":      "acs:oss:cn-hangzhou:${data.alicloud_account.default.id}:${data.alicloud_oss_buckets.default.buckets.0.name}/t2",
										"name":        "",
									},
									{
										"match_state": "true",
										"suffix":      "t4",
										"prefix":      "acs:oss:cn-hangzhou:${data.alicloud_account.default.id}:${data.alicloud_oss_buckets.default.buckets.0.name}/t3",
										"name":        "",
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_oss_event_parameters.#": "1",
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
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"linked_external_source": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

func TestAccAliCloudEventBridgeEventSourceV2_basic11707_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_event_bridge_event_source_v2.default"
	ra := resourceAttrInit(resourceId, AliCloudEventBridgeEventSourceV2Map11707)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EventBridgeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEventBridgeEventSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceventbridge%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEventBridgeEventSourceV2BasicDependence11707)
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
					"event_bus_name":    "${alicloud_event_bridge_event_bus.default.event_bus_name}",
					"event_source_name": name,
					"source_oss_event_parameters": []map[string]interface{}{
						{
							"event_types": []string{
								"oss:ObjectCreated:AppendObject", "oss:ObjectCreated:CopyObject"},
							"sts_role_arn": "acs:ram::${data.alicloud_account.default.id}:role/${alicloud_ram_role_policy_attachment.AliyunEventBridgePutEventsPolicy.role_name}",
							"match_rules": [][]map[string]interface{}{
								{
									{
										"match_state": "true",
										"suffix":      "",
										"prefix":      "",
										"name":        "acs:oss:cn-hangzhou:${data.alicloud_account.default.id}:${data.alicloud_oss_buckets.default.buckets.0.name}/t1",
									},
									{
										"match_state": "true",
										"suffix":      "",
										"prefix":      "acs:oss:cn-hangzhou:${data.alicloud_account.default.id}:${data.alicloud_oss_buckets.default.buckets.0.name}/t2",
										"name":        "",
									},
									{
										"match_state": "true",
										"suffix":      "t4",
										"prefix":      "acs:oss:cn-hangzhou:${data.alicloud_account.default.id}:${data.alicloud_oss_buckets.default.buckets.0.name}/t3",
										"name":        "",
									},
								},
							},
						},
					},
					"description":            name,
					"linked_external_source": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"event_bus_name":                CHECKSET,
						"event_source_name":             name,
						"source_oss_event_parameters.#": "1",
						"description":                   name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linked_external_source"},
			},
		},
	})
}

var AliCloudEventBridgeEventSourceV2Map11707 = map[string]string{}

func AliCloudEventBridgeEventSourceV2BasicDependence11707(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_account" "default" {
	}

	data "alicloud_oss_buckets" "default" {
	}

	resource "alicloud_event_bridge_event_bus" "default" {
  		event_bus_name = var.name
	}

	resource "alicloud_ram_role" "default" {
  		role_name                   = var.name
  		assume_role_policy_document = <<EOF
		{
			"Statement": [
				{
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": {
						"Service": [
							"sendevent-mns.eventbridge.aliyuncs.com"
						]
					}
				}
		  	],
			"Version": "1"
		}
	  	EOF
		force    = true
	}

	resource "alicloud_ram_role" "update" {
  		role_name                   = "${var.name}update"
  		assume_role_policy_document = <<EOF
		{
			"Statement": [
				{
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": {
						"Service": [
							"sendevent-mns.eventbridge.aliyuncs.com"
						]
					}
				}
		  	],
			"Version": "1"
		}
	  	EOF
		force    = true
	}

	resource "alicloud_ram_role_policy_attachment" "AliyunMNSFullAccess" {
  		role_name   = alicloud_ram_role.default.role_name
  		policy_name = "AliyunMNSFullAccess"
  		policy_type = "System"
	}

	resource "alicloud_ram_role_policy_attachment" "AliyunEventBridgePutEventsPolicy" {
  		role_name   = alicloud_ram_role_policy_attachment.AliyunMNSFullAccess.role_name
  		policy_name = "AliyunEventBridgePutEventsPolicy"
  		policy_type = "System"
	}

	resource "alicloud_ram_role_policy_attachment" "default" {
  		role_name   = alicloud_ram_role.update.role_name
  		policy_name = "AliyunMNSFullAccess"
  		policy_type = "System"
	}

	resource "alicloud_ram_role_policy_attachment" "update" {
  		role_name   = alicloud_ram_role_policy_attachment.default.role_name
  		policy_name = "AliyunEventBridgePutEventsPolicy"
  		policy_type = "System"
	}
`, name)
}

// Test EventBridge EventSource. <<< Resource test cases, automatically generated.
