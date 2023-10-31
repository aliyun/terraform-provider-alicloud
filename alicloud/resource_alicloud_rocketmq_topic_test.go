package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Rocketmq Topic. >>> Resource test cases, automatically generated.
// Case 4729
func TestAccAliCloudRocketmqTopic_basic4729(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_topic.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqTopicMap4729)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqTopicBasicDependence4729)
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
					"instance_id":  "${alicloud_rocketmq_instance.createInstance.id}",
					"topic_name":   name,
					"message_type": "TRANSACTION",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  CHECKSET,
						"topic_name":   name,
						"message_type": "TRANSACTION",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "1111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "1111",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_rocketmq_instance.createInstance.id}",
					"message_type": "TRANSACTION",
					"topic_name":   name + "_update",
					"remark":       "1111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  CHECKSET,
						"message_type": "TRANSACTION",
						"topic_name":   name + "_update",
						"remark":       "1111",
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

var AlicloudRocketmqTopicMap4729 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudRocketmqTopicBasicDependence4729(name string) string {
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
  description  = "111"
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

// Case 4728
func TestAccAliCloudRocketmqTopic_basic4728(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_topic.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqTopicMap4728)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqTopicBasicDependence4728)
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
					"instance_id":  "${alicloud_rocketmq_instance.createInstance.id}",
					"topic_name":   name,
					"message_type": "DELAY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  CHECKSET,
						"topic_name":   name,
						"message_type": "DELAY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "1111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "1111",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_rocketmq_instance.createInstance.id}",
					"message_type": "DELAY",
					"topic_name":   name + "_update",
					"remark":       "1111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  CHECKSET,
						"message_type": "DELAY",
						"topic_name":   name + "_update",
						"remark":       "1111",
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

var AlicloudRocketmqTopicMap4728 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudRocketmqTopicBasicDependence4728(name string) string {
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
  description  = "111"
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

// Case 4727
func TestAccAliCloudRocketmqTopic_basic4727(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_topic.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqTopicMap4727)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqTopicBasicDependence4727)
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
					"instance_id":  "${alicloud_rocketmq_instance.createInstance.id}",
					"topic_name":   name,
					"message_type": "FIFO",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  CHECKSET,
						"topic_name":   name,
						"message_type": "FIFO",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "1111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "1111",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_rocketmq_instance.createInstance.id}",
					"message_type": "FIFO",
					"topic_name":   name + "_update",
					"remark":       "1111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  CHECKSET,
						"message_type": "FIFO",
						"topic_name":   name + "_update",
						"remark":       "1111",
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

var AlicloudRocketmqTopicMap4727 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudRocketmqTopicBasicDependence4727(name string) string {
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
  description  = "111"
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

// Case 4416
func TestAccAliCloudRocketmqTopic_basic4416(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_topic.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqTopicMap4416)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqTopicBasicDependence4416)
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
					"instance_id":  "${alicloud_rocketmq_instance.createInstance.id}",
					"topic_name":   name,
					"message_type": "NORMAL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  CHECKSET,
						"topic_name":   name,
						"message_type": "NORMAL",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "1111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "1111",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "2222",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "2222",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "1111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "1111",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "2222",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "2222",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_rocketmq_instance.createInstance.id}",
					"message_type": "NORMAL",
					"topic_name":   name + "_update",
					"remark":       "1111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  CHECKSET,
						"message_type": "NORMAL",
						"topic_name":   name + "_update",
						"remark":       "1111",
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

var AlicloudRocketmqTopicMap4416 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudRocketmqTopicBasicDependence4416(name string) string {
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
  description  = "111"
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

// Case 4729  twin
func TestAccAliCloudRocketmqTopic_basic4729_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_topic.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqTopicMap4729)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqTopicBasicDependence4729)
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
					"instance_id":  "${alicloud_rocketmq_instance.createInstance.id}",
					"message_type": "TRANSACTION",
					"topic_name":   name,
					"remark":       "1111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  CHECKSET,
						"message_type": "TRANSACTION",
						"topic_name":   name,
						"remark":       "1111",
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

// Case 4728  twin
func TestAccAliCloudRocketmqTopic_basic4728_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_topic.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqTopicMap4728)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqTopicBasicDependence4728)
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
					"instance_id":  "${alicloud_rocketmq_instance.createInstance.id}",
					"message_type": "DELAY",
					"topic_name":   name,
					"remark":       "1111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  CHECKSET,
						"message_type": "DELAY",
						"topic_name":   name,
						"remark":       "1111",
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

// Case 4727  twin
func TestAccAliCloudRocketmqTopic_basic4727_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_topic.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqTopicMap4727)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqTopicBasicDependence4727)
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
					"instance_id":  "${alicloud_rocketmq_instance.createInstance.id}",
					"message_type": "FIFO",
					"topic_name":   name,
					"remark":       "1111",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  CHECKSET,
						"message_type": "FIFO",
						"topic_name":   name,
						"remark":       "1111",
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

// Case 4416  twin
func TestAccAliCloudRocketmqTopic_basic4416_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_topic.default"
	ra := resourceAttrInit(resourceId, AlicloudRocketmqTopicMap4416)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRocketmqTopicBasicDependence4416)
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
					"instance_id":  "${alicloud_rocketmq_instance.createInstance.id}",
					"message_type": "NORMAL",
					"topic_name":   name,
					"remark":       "2222",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  CHECKSET,
						"message_type": "NORMAL",
						"topic_name":   name,
						"remark":       "2222",
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

// Test Rocketmq Topic. <<< Resource test cases, automatically generated.
