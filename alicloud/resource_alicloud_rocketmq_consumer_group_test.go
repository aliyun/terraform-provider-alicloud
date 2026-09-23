package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Rocketmq ConsumerGroup. >>> Resource test cases, automatically generated.
// Case 4419
func TestAccAliCloudRocketmqConsumerGroup_basic4419(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_consumer_group.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqConsumerGroupMap4419)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqConsumerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqconsumergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqConsumerGroupBasicDependence4419)
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
					"delivery_order_type": "Concurrently",
					"consumer_group_id":   "pop-test-group",
					"instance_id":         "${alicloud_rocketmq_instance.default.id}",
					"consume_retry_policy": []map[string]interface{}{
						{
							"retry_policy": "DefaultRetryPolicy",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_order_type":    "Concurrently",
						"consumer_group_id":      "pop-test-group",
						"instance_id":            CHECKSET,
						"consume_retry_policy.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "123321",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "123321",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"consume_retry_policy": []map[string]interface{}{
						{
							"max_retry_times": "10",
							"retry_policy":    "DefaultRetryPolicy",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"consume_retry_policy.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"consume_retry_policy": []map[string]interface{}{
						{
							"max_retry_times":          "10",
							"retry_policy":             "DefaultRetryPolicy",
							"dead_letter_target_topic": "${alicloud_rocketmq_topic.defaultoHnhFz.topic_name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"consume_retry_policy.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_order_type": "Orderly",
					"consume_retry_policy": []map[string]interface{}{
						{
							"max_retry_times":          "5",
							"retry_policy":             "FixedRetryPolicy",
							"dead_letter_target_topic": "${alicloud_rocketmq_topic.default2j1J7Q.topic_name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_order_type":    "Orderly",
						"consume_retry_policy.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "333",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "333",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_receive_tps": "1500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_receive_tps": "1500",
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

var AliCloudRocketmqConsumerGroupMap4419 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AliCloudRocketmqConsumerGroupBasicDependence4419(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "createVPC" {
  description = "example"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "createVSwitch" {
  description  = "example"
  vpc_id       = alicloud_vpc.createVPC.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_rocketmq_instance" "default" {
  product_info {
    msg_process_spec       = "rmq.u2.10xlarge"
    send_receive_ratio     = "0.3"
    message_retention_time = "70"
  }
  service_code      = "rmq"
  payment_type      = "PayAsYouGo"
  instance_name     = var.name
  sub_series_code   = "cluster_ha"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  remark            = "example"
  ip_whitelists     = ["192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"]
  software {
    maintain_time = "02:00-06:00"
  }
  tags = {
    Created = "TF"
    For     = "example"
  }
  series_code = "ultimate"
  network_info {
    vpc_info {
      vpc_id = alicloud_vpc.createVPC.id
      vswitches {
        vswitch_id = alicloud_vswitch.createVSwitch.id
      }
    }
    internet_info {
      internet_spec      = "enable"
      flow_out_type      = "payByBandwidth"
      flow_out_bandwidth = "30"
    }
  }
}

resource "alicloud_rocketmq_topic" "defaultoHnhFz" {
  instance_id  = alicloud_rocketmq_instance.default.id
  message_type = "NORMAL"
  topic_name   = "${var.name}-1"
}

resource "alicloud_rocketmq_topic" "default2j1J7Q" {
  instance_id  = alicloud_rocketmq_instance.default.id
  message_type = "NORMAL"
  topic_name   = "${var.name}-2"
}
`, name)
}

// Case 4419  twin
func TestAccAliCloudRocketmqConsumerGroup_basic4419_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_consumer_group.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqConsumerGroupMap4419)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqConsumerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqconsumergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqConsumerGroupBasicDependence4419)
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
					"consumer_group_id": "pop-test-group",
					"instance_id":       "${alicloud_rocketmq_instance.default.id}",
					"consume_retry_policy": []map[string]interface{}{
						{
							"dead_letter_target_topic": "${alicloud_rocketmq_topic.defaultoHnhFz.topic_name}",
							"max_retry_times":          "10",
							"retry_policy":             "DefaultRetryPolicy",
						},
					},
					"delivery_order_type": "Concurrently",
					"max_receive_tps":     "1500",
					"remark":              "123321",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"consumer_group_id":      "pop-test-group",
						"instance_id":            CHECKSET,
						"consume_retry_policy.#": "1",
						"delivery_order_type":    "Concurrently",
						"max_receive_tps":        "1500",
						"remark":                 "123321",
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

// Test Rocketmq ConsumerGroup. <<< Resource test cases, automatically generated.
