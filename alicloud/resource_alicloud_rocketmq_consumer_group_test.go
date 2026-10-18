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

// testSweepRocketMqConsumerGroup sweeps test RocketMQ consumer groups across
// all instances in the region. It first lists instances via ListInstances
// (GET /instances), then for each instance lists its consumer groups via
// ListConsumerGroups (GET /instances/{instanceId}/consumerGroups) and deletes
// any consumer group whose consumerGroupId matches the standard test prefixes
// (or all of them under SWEEP_ALL=1), skipping entries protected by the
// do-not-delete markers. The instance sweeper declares this as a Dependency
// so consumer groups are removed before the parent instance is released.
func init() {
	resource.AddTestSweepers("alicloud_rocketmq_consumer_group", &resource.Sweeper{
		Name: "alicloud_rocketmq_consumer_group",
		F:    testSweepRocketMqConsumerGroup,
	})
}

func testSweepRocketMqConsumerGroup(region string) error {
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

	// List every RocketMQ instance in the region; consumer groups live
	// under instances, so we have to enumerate instances first.
	instanceAction := "/instances"
	instanceQuery := make(map[string]*string)
	instanceQuery["pageNumber"] = tea.String("1")
	instanceQuery["pageSize"] = tea.String("200")
	instanceResp, err := client.RoaGet("RocketMQ", "2022-08-01", instanceAction, instanceQuery, nil, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve rocketmq instances for consumer group sweep: %s", err)
		return nil
	}
	instanceRaw, err := jsonpath.Get("$.data.list", instanceResp)
	if err != nil {
		// When the instance list is empty, the RocketMQ ListInstances API
		// omits the "list" key entirely; treat as "no instances to sweep".
		log.Printf("[INFO] rocketmq instance list empty or parse failed for consumer group sweep: %s", err)
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

		// List consumer groups under this instance.
		action := fmt.Sprintf("/instances/%s/consumerGroups", instanceId)
		query := make(map[string]*string)
		query["pageNumber"] = tea.String("1")
		query["pageSize"] = tea.String("100")
		response, err := client.RoaGet("RocketMQ", "2022-08-01", action, query, nil, nil)
		if err != nil {
			log.Printf("[ERROR] Failed to list consumer groups for instance %s: %s", instanceId, err)
			continue
		}
		grpRaw, err := jsonpath.Get("$.data.list", response)
		if err != nil {
			// When the consumer-group list is empty, the API omits the
			// "list" key entirely; treat as "no consumer groups under this
			// instance" and move on to the next instance.
			log.Printf("[INFO] rocketmq consumer group list empty or parse failed for instance %s: %s", instanceId, err)
			continue
		}
		groups, _ := grpRaw.([]interface{})

		for _, g := range groups {
			item, ok := g.(map[string]interface{})
			if !ok {
				continue
			}
			name, _ := item["consumerGroupId"].(string)
			if name == "" {
				continue
			}
			// Name-based NODELETING/do-not-delete protection. Runs even
			// under SWEEP_ALL=1 so explicitly protected consumer groups
			// are never force-deleted.
			if protectedSweepVpcName(name) {
				log.Printf("[INFO] Skipping rocketmq consumer group (do-not-delete marker): %s/%s", instanceId, name)
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
					log.Printf("[INFO] Skipping rocketmq consumer group: %s/%s", instanceId, name)
					continue
				}
			}
			log.Printf("[INFO] delete rocketmq consumer group: %s/%s", instanceId, name)
			deleteAction := fmt.Sprintf("/instances/%s/consumerGroups/%s", instanceId, name)
			deleteQuery := make(map[string]*string)
			if _, err := client.RoaDelete("RocketMQ", "2022-08-01", deleteAction, deleteQuery, nil, nil, true); err != nil {
				log.Printf("[ERROR] Failed to delete rocketmq consumer group (%s/%s): %s", instanceId, name, err)
			}
		}
	}

	return nil
}

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
		ProviderFactories: testAccProviderFactory,
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
		ProviderFactories: testAccProviderFactory,
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
