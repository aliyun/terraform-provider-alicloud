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
	ra := resourceAttrInit(resourceId, AlicloudRocketmqConsumerGroupMap4419)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqConsumerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqconsumergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqConsumerGroupBasicDependence4419)
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
					"instance_id":         "${alicloud_rocketmq_instance.createInstance.id}",
					"consume_retry_policy": []map[string]interface{}{
						{
							"max_retry_times": "10",
							"retry_policy":    "DefaultRetryPolicy",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_order_type": "Concurrently",
						"consumer_group_id":   "pop-test-group",
						"instance_id":         CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_order_type": "Concurrently",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_order_type": "Concurrently",
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
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delivery_order_type": "Orderly",
					"consume_retry_policy": []map[string]interface{}{
						{
							"max_retry_times": "5",
							"retry_policy":    "FixedRetryPolicy",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_order_type": "Orderly",
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
					"delivery_order_type": "Concurrently",
					"consume_retry_policy": []map[string]interface{}{
						{
							"max_retry_times": "10",
							"retry_policy":    "DefaultRetryPolicy",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delivery_order_type": "Concurrently",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"consumer_group_id": "pop-test-group",
					"instance_id":       "${alicloud_rocketmq_instance.createInstance.id}",
					"consume_retry_policy": []map[string]interface{}{
						{
							"max_retry_times": "10",
							"retry_policy":    "DefaultRetryPolicy",
						},
					},
					"delivery_order_type": "Concurrently",
					"remark":              "123321",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"consumer_group_id":   "pop-test-group",
						"instance_id":         CHECKSET,
						"delivery_order_type": "Concurrently",
						"remark":              "123321",
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

var AlicloudRocketmqConsumerGroupMap4419 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudRocketmqConsumerGroupBasicDependence4419(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "createVpc" {
  description = "111"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name

}

resource "alicloud_vswitch" "createVswitch" {
  description  = "1111"
  vpc_id       = alicloud_vpc.createVpc.id
  zone_id      = data.alicloud_zones.default.zones.0.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name

}

resource "alicloud_rocketmq_instance" "createInstance" {
  auto_renew_period = "1"
  product_info {
    msg_process_spec       = "rmq.p2.4xlarge"
    send_receive_ratio     = 0.3
    message_retention_time = "70"
  }
  network_info {
    vpc_info {
      vpc_id     = alicloud_vpc.createVpc.id
      vswitch_id = alicloud_vswitch.createVswitch.id
    }
    internet_info {
      internet_spec      = "enable"
      flow_out_type      = "payByBandwidth"
      flow_out_bandwidth = "30"
    }
  }
  period          = "1"
  sub_series_code = "cluster_ha"
  remark          = "自动化测试购买使用11"
  instance_name   = var.name

  service_code = "rmq"
  series_code  = "professional"
  payment_type = "PayAsYouGo"
  period_unit = "Month"
}


`, name)
}

// Case 4419  twin
func TestAccAliCloudRocketmqConsumerGroup_basic4419_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_consumer_group.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqConsumerGroupMap4419)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqConsumerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqconsumergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqConsumerGroupBasicDependence4419)
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
					"instance_id":       "${alicloud_rocketmq_instance.createInstance.id}",
					"consume_retry_policy": []map[string]interface{}{
						{
							"max_retry_times": "10",
							"retry_policy":    "DefaultRetryPolicy",
						},
					},
					"delivery_order_type": "Concurrently",
					"remark":              "123321",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"consumer_group_id":   "pop-test-group",
						"instance_id":         CHECKSET,
						"delivery_order_type": "Concurrently",
						"remark":              "123321",
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
