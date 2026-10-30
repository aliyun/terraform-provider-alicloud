package rocketmq_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rocketmq"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Rocketmq Topic. >>> Resource test cases, automatically generated.
// Case 4729
func TestAccAliCloudRocketmqTopic_basic4729(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_topic.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRocketmqTopicMap4729)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeTopic")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRocketmqTopicBasicDependence4729)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_rocketmq_instance.default.id}",
					"topic_name":   name,
					"message_type": "TRANSACTION",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  acctest.CHECKSET,
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
					"max_send_tps": "1500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_send_tps": "1500",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "222",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "222",
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

var AliCloudRocketmqTopicMap4729 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func AliCloudRocketmqTopicBasicDependence4729(name string) string {
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
  ip_whitelists     = ["192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16", "0.0.0.0/0"]
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
`, name)
}

// Case 4728
func TestAccAliCloudRocketmqTopic_basic4728(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_topic.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRocketmqTopicMap4728)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeTopic")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRocketmqTopicBasicDependence4728)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_rocketmq_instance.default.id}",
					"topic_name":   name,
					"message_type": "DELAY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  acctest.CHECKSET,
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
					"max_send_tps": "1500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_send_tps": "1500",
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

var AliCloudRocketmqTopicMap4728 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func AliCloudRocketmqTopicBasicDependence4728(name string) string {
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
  ip_whitelists     = ["192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16", "0.0.0.0/0"]
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
`, name)
}

// Case 4727
func TestAccAliCloudRocketmqTopic_basic4727(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_topic.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRocketmqTopicMap4727)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeTopic")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRocketmqTopicBasicDependence4727)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_rocketmq_instance.default.id}",
					"topic_name":   name,
					"message_type": "FIFO",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  acctest.CHECKSET,
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
					"max_send_tps": "1500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_send_tps": "1500",
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

var AliCloudRocketmqTopicMap4727 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func AliCloudRocketmqTopicBasicDependence4727(name string) string {
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
  ip_whitelists     = ["192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16", "0.0.0.0/0"]
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
`, name)
}

// Case 4416
func TestAccAliCloudRocketmqTopic_basic4416(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_topic.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRocketmqTopicMap4416)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeTopic")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRocketmqTopicBasicDependence4416)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_rocketmq_instance.default.id}",
					"topic_name":   name,
					"message_type": "NORMAL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  acctest.CHECKSET,
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
					"max_send_tps": "1500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_send_tps": "1500",
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

var AliCloudRocketmqTopicMap4416 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func AliCloudRocketmqTopicBasicDependence4416(name string) string {
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
  ip_whitelists     = ["192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16", "0.0.0.0/0"]
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
`, name)
}

// Case 4729  twin
func TestAccAliCloudRocketmqTopic_basic4729_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_topic.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRocketmqTopicMap4729)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeTopic")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRocketmqTopicBasicDependence4729)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_rocketmq_instance.default.id}",
					"message_type": "TRANSACTION",
					"topic_name":   name,
					"remark":       "1111",
					"max_send_tps": "1500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  acctest.CHECKSET,
						"message_type": "TRANSACTION",
						"topic_name":   name,
						"remark":       "1111",
						"max_send_tps": "1500",
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRocketmqTopicMap4728)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeTopic")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRocketmqTopicBasicDependence4728)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_rocketmq_instance.default.id}",
					"message_type": "DELAY",
					"topic_name":   name,
					"remark":       "1111",
					"max_send_tps": "1500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  acctest.CHECKSET,
						"message_type": "DELAY",
						"topic_name":   name,
						"remark":       "1111",
						"max_send_tps": "1500",
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRocketmqTopicMap4727)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeTopic")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRocketmqTopicBasicDependence4727)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_rocketmq_instance.default.id}",
					"message_type": "FIFO",
					"topic_name":   name,
					"remark":       "1111",
					"max_send_tps": "1500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  acctest.CHECKSET,
						"message_type": "FIFO",
						"topic_name":   name,
						"remark":       "1111",
						"max_send_tps": "1500",
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRocketmqTopicMap4416)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeTopic")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRocketmqTopicBasicDependence4416)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_rocketmq_instance.default.id}",
					"message_type": "NORMAL",
					"topic_name":   name,
					"remark":       "2222",
					"max_send_tps": "1500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  acctest.CHECKSET,
						"message_type": "NORMAL",
						"topic_name":   name,
						"remark":       "2222",
						"max_send_tps": "1500",
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

// Case supported fixture: representative lifecycle on a supported instance config
// (reuses the verified-passing Instance_basic4101 instance configuration, without
// the deprecated software.maintain_time / ip_whitelists fields that cause non-refresh
// plan drift). Validates Topic Create/Read/Update/Import/Delete.
func TestAccAliCloudRocketmqTopic_basicSupported(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_topic.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRocketmqTopicMapSupported)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rocketmq.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeTopic")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRocketmqTopicBasicDependenceSupported)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_rocketmq_instance.default.id}",
					"topic_name":   name,
					"message_type": "NORMAL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":  acctest.CHECKSET,
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudRocketmqTopicMapSupported = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
}

func AliCloudRocketmqTopicBasicDependenceSupported(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "createVPC" {
  description = "1111"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "createVSwitch" {
  description  = "1111"
  vpc_id       = alicloud_vpc.createVPC.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_rocketmq_instance" "default" {
  service_code    = "rmq"
  series_code     = "professional"
  sub_series_code = "cluster_ha"
  payment_type    = "PayAsYouGo"
  instance_name   = var.name

  product_info {
    auto_scaling           = "false"
    msg_process_spec       = "rmq.p2.4xlarge"
    send_receive_ratio     = "0.3"
    message_retention_time = "70"
  }

  network_info {
    vpc_info {
      vpc_id     = alicloud_vpc.createVPC.id
      vswitch_id = alicloud_vswitch.createVSwitch.id
    }

    internet_info {
      internet_spec      = "enable"
      flow_out_type      = "payByBandwidth"
      flow_out_bandwidth = "30"
    }
  }
}
`, name)
}
