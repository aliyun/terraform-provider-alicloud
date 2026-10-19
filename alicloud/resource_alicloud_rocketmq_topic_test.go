package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// testSweepRocketMqTopic sweeps test RocketMQ topics across all instances in
// the region. It first lists instances via ListInstances (GET /instances),
// then for each instance lists its topics via ListTopics
// (GET /instances/{instanceId}/topics) and deletes any topic whose topicName
// matches the standard test prefixes (or all of them under SWEEP_ALL=1),
// skipping entries protected by the do-not-delete markers. The instance
// sweeper declares this as a Dependency so topics are removed before the
// parent instance is released.
func init() {
	resource.AddTestSweepers("alicloud_rocketmq_topic", &resource.Sweeper{
		Name: "alicloud_rocketmq_topic",
		F:    testSweepRocketMqTopic,
	})
}

func testSweepRocketMqTopic(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"testAcc",
	}

	// List every RocketMQ instance in the region; topics live under
	// instances, so we have to enumerate instances first.
	instanceAction := "/instances"
	instanceQuery := make(map[string]*string)
	instanceQuery["pageNumber"] = tea.String("1")
	instanceQuery["pageSize"] = tea.String("200")
	instanceResp, err := client.RoaGet("RocketMQ", "2022-08-01", instanceAction, instanceQuery, nil, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve rocketmq instances for topic sweep: %s", err)
		return nil
	}
	instanceRaw, err := jsonpath.Get("$.data.list", instanceResp)
	if err != nil {
		// When the instance list is empty, the RocketMQ ListInstances API
		// omits the "list" key entirely; treat as "no instances to sweep".
		log.Printf("[INFO] rocketmq instance list empty or parse failed for topic sweep: %s", err)
		return nil
	}
	instances, _ := instanceRaw.([]interface{})

	for _, inst := range instances {
		instanceItem, ok := inst.(map[string]interface{})
		if !ok {
			continue
		}
		instanceId, _ := instanceItem["instanceId"].(string)
		if instanceId == "" {
			continue
		}

		// List topics under this instance.
		action := fmt.Sprintf("/instances/%s/topics", instanceId)
		query := make(map[string]*string)
		query["pageNumber"] = tea.String("1")
		query["pageSize"] = tea.String("100")
		response, err := client.RoaGet("RocketMQ", "2022-08-01", action, query, nil, nil)
		if err != nil {
			log.Printf("[ERROR] Failed to list topics for instance %s: %s", instanceId, err)
			continue
		}
		topicRaw, err := jsonpath.Get("$.data.list", response)
		if err != nil {
			// When the topic list is empty, the API omits the "list" key
			// entirely; treat as "no topics under this instance" and move
			// on to the next instance.
			log.Printf("[INFO] rocketmq topic list empty or parse failed for instance %s: %s", instanceId, err)
			continue
		}
		topics, _ := topicRaw.([]interface{})

		for _, tp := range topics {
			item, ok := tp.(map[string]interface{})
			if !ok {
				continue
			}
			name, _ := item["topicName"].(string)
			if name == "" {
				continue
			}
			// Name-based NODELETING/do-not-delete protection. Runs even
			// under SWEEP_ALL=1 so explicitly protected topics are never
			// force-deleted.
			if protectedSweepVpcName(name) {
				log.Printf("[INFO] Skipping rocketmq topic (do-not-delete marker): %s/%s", instanceId, name)
				continue
			}
			skip := true
			if !sweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
				if skip {
					log.Printf("[INFO] Skipping rocketmq topic: %s/%s", instanceId, name)
					continue
				}
			}
			log.Printf("[INFO] delete rocketmq topic: %s/%s", instanceId, name)
			deleteAction := fmt.Sprintf("/instances/%s/topics/%s", instanceId, name)
			deleteQuery := make(map[string]*string)
			if _, err := client.RoaDelete("RocketMQ", "2022-08-01", deleteAction, deleteQuery, nil, nil, true); err != nil {
				log.Printf("[ERROR] Failed to delete rocketmq topic (%s/%s): %s", instanceId, name, err)
			}
		}
	}

	return nil
}

// Test Rocketmq Topic. >>> Resource test cases, automatically generated.
// Case 4729
func TestAccAliCloudRocketmqTopic_basic4729(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rocketmq_topic.default"
	ra := resourceAttrInit(resourceId, AliCloudRocketmqTopicMap4729)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqTopicBasicDependence4729)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_rocketmq_instance.default.id}",
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
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudRocketmqTopicMap4728)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqTopicBasicDependence4728)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_rocketmq_instance.default.id}",
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
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudRocketmqTopicMap4727)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqTopicBasicDependence4727)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_rocketmq_instance.default.id}",
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
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudRocketmqTopicMap4416)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqTopicBasicDependence4416)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":  "${alicloud_rocketmq_instance.default.id}",
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
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudRocketmqTopicMap4729)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqTopicBasicDependence4729)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"instance_id":  CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudRocketmqTopicMap4728)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqTopicBasicDependence4728)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"instance_id":  CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudRocketmqTopicMap4727)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqTopicBasicDependence4727)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"instance_id":  CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudRocketmqTopicMap4416)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RocketmqServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRocketmqTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqtopic%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRocketmqTopicBasicDependence4416)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"instance_id":  CHECKSET,
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
