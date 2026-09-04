package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func init() {
	resource.AddTestSweepers("alicloud_alikafka_topic", &resource.Sweeper{
		Name: "alicloud_alikafka_topic",
		F:    testSweepAlikafkaTopic,
	})
}

func testSweepAlikafkaTopic(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	alikafkaService := AlikafkaService{client}

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	instanceListReq := alikafka.CreateGetInstanceListRequest()
	instanceListReq.RegionId = region

	raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
		return alikafkaClient.GetInstanceList(instanceListReq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve alikafka instance in service list: %s", err)
	}

	instanceListResp, _ := raw.(*alikafka.GetInstanceListResponse)

	var instanceIds []string
	for _, v := range instanceListResp.InstanceList.InstanceVO {
		instanceIds = append(instanceIds, v.InstanceId)
	}

	for _, instanceId := range instanceIds {

		// Control the topic list request rate.
		time.Sleep(time.Duration(400) * time.Millisecond)

		request := alikafka.CreateGetTopicListRequest()
		request.InstanceId = instanceId
		request.RegionId = defaultRegionToTest
		raw, err = alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.GetTopicList(request)
		})

		if err != nil {
			log.Printf("[ERROR] Failed to retrieve alikafka topics on instance (%s): %s", instanceId, err)
			continue
		}

		topicListResp, _ := raw.(*alikafka.GetTopicListResponse)
		topics := topicListResp.TopicList.TopicVO

		for _, v := range topics {
			name := v.Topic
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping alikafka topic: %s ", name)
				continue
			}
			log.Printf("[INFO] delete alikafka topic: %s ", name)

			// Control the topic delete rate
			time.Sleep(time.Duration(400) * time.Millisecond)

			request := alikafka.CreateDeleteTopicRequest()
			request.InstanceId = instanceId
			request.Topic = v.Topic
			request.RegionId = defaultRegionToTest

			_, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
				return alikafkaClient.DeleteTopic(request)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete alikafka topic (%s): %s", name, err)
			}
		}
	}

	return nil
}

// Test Alikafka Topic. >>> Resource test cases, automatically generated.
// Case topic全生命周期 10065
func TestAccAliCloudAlikafkaTopic_basic10065(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_topic.default"
	ra := resourceAttrInit(resourceId, AliCloudAlikafkaTopicMap10065)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlikafkaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlikafkaTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalikafka%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAlikafkaTopicBasicDependence10065)
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
					"instance_id": "${alicloud_alikafka_instance.default.id}",
					"topic":       name,
					"remark":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"topic":       name,
						"remark":      name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"partition_num": "18",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"partition_num": "18",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
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
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

func TestAccAliCloudAlikafkaTopic_basic10065_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_topic.default"
	ra := resourceAttrInit(resourceId, AliCloudAlikafkaTopicMap10065)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlikafkaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlikafkaTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalikafka%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAlikafkaTopicBasicDependence10065)
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
					"instance_id":   "${alicloud_alikafka_instance.default.id}",
					"topic":         name,
					"remark":        name,
					"local_topic":   "false",
					"compact_topic": "false",
					"partition_num": "6",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":   CHECKSET,
						"topic":         name,
						"remark":        name,
						"local_topic":   "false",
						"compact_topic": "false",
						"partition_num": "6",
						"tags.%":        "2",
						"tags.Created":  "TF",
						"tags.For":      "Test",
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

func TestAccAliCloudAlikafkaTopic_basic10066(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_topic.default"
	ra := resourceAttrInit(resourceId, AliCloudAlikafkaTopicMap10065)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlikafkaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlikafkaTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalikafka%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAlikafkaTopicBasicDependence10065)
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
					"instance_id": "${alicloud_alikafka_instance.default.id}",
					"topic":       name,
					"remark":      name,
					"local_topic": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"topic":       name,
						"remark":      name,
						"local_topic": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"partition_num": "18",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"partition_num": "18",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configs": `{\"message.format.version\":\"2.2.0\",\"max.message.bytes\":\"10485760\",\"min.insync.replicas\":\"1\",\"replication-factor\":\"3\",\"retention.ms\":\"3600000\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configs": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configs": `{\"message.format.version\":\"2.0\",\"max.message.bytes\":\"10485760\",\"min.insync.replicas\":\"2\",\"replication-factor\":\"3\",\"retention.ms\":\"10800000\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configs": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
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
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

func TestAccAliCloudAlikafkaTopic_basic10066_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_topic.default"
	ra := resourceAttrInit(resourceId, AliCloudAlikafkaTopicMap10065)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlikafkaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlikafkaTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalikafka%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAlikafkaTopicBasicDependence10065)
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
					"instance_id":   "${alicloud_alikafka_instance.default.id}",
					"topic":         name,
					"remark":        name,
					"local_topic":   "true",
					"compact_topic": "true",
					"configs":       `{\"message.format.version\":\"2.2.0\",\"max.message.bytes\":\"10485760\",\"min.insync.replicas\":\"1\",\"replication-factor\":\"2\",\"retention.ms\":\"3600000\"}`,
					"partition_num": "18",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":   CHECKSET,
						"topic":         name,
						"remark":        name,
						"local_topic":   "true",
						"compact_topic": "true",
						"configs":       CHECKSET,
						"partition_num": "18",
						"tags.%":        "2",
						"tags.Created":  "TF",
						"tags.For":      "Test",
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

func TestAccAliCloudAlikafkaTopic_multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_topic.default.5"
	ra := resourceAttrInit(resourceId, AliCloudAlikafkaTopicMap10065)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlikafkaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlikafkaTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalikafka%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAlikafkaTopicBasicDependence10065)
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
					"count":         "6",
					"instance_id":   "${alicloud_alikafka_instance.default.id}",
					"topic":         name + "-${count.index}",
					"remark":        name + "-${count.index}",
					"local_topic":   "true",
					"compact_topic": "true",
					"configs":       `{\"message.format.version\":\"2.2.0\",\"max.message.bytes\":\"10485760\",\"min.insync.replicas\":\"1\",\"replication-factor\":\"2\",\"retention.ms\":\"3600000\"}`,
					"partition_num": "18",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":   CHECKSET,
						"topic":         name + fmt.Sprint(-5),
						"remark":        name + fmt.Sprint(-5),
						"local_topic":   "true",
						"compact_topic": "true",
						"configs":       CHECKSET,
						"partition_num": "18",
						"tags.%":        "2",
						"tags.Created":  "TF",
						"tags.For":      "Test",
					}),
				),
			},
		},
	})
}

var AliCloudAlikafkaTopicMap10065 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AliCloudAlikafkaTopicBasicDependence10065(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_resource_creation = "VSwitch"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "172.16.0.0/12"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "172.16.0.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_alikafka_instance" "default" {
  		name            = var.name
  		partition_num   = 50
  		disk_type       = "1"
  		disk_size       = "500"
  		deploy_type     = "5"
  		io_max          = "20"
  		spec_type       = "professional"
  		service_version = "2.2.0"
  		vswitch_id      = alicloud_vswitch.default.id
  		security_group  = alicloud_security_group.default.id
  		config          = "{\"enable.acl\":\"true\"}"
	}
`, name)
}

// Test Alikafka Topic. <<< Resource test cases, automatically generated.

// TestUnitAlikafkaTopicConfigsSchemaDiffSuppress drives the diff through the resource schema so
// that the configs attribute stays wired to the subset comparator. Reproduces the case where a
// plan never converges because the instance injects config keys the configuration never set.
func TestUnitAlikafkaTopicConfigsSchemaDiffSuppress(t *testing.T) {
	suppress := resourceAliCloudAlikafkaTopic().Schema["configs"].DiffSuppressFunc
	if suppress == nil {
		t.Fatal("configs has no DiffSuppressFunc")
	}

	stateConfigs := `{"cloud.native.topic.type":"true","replication-factor":"3","max.message.bytes":"1048576","retention.hours":"72"}`

	testCases := []struct {
		name        string
		old         string
		new         string
		expected    bool
		description string
	}{
		{
			name:        "ConvergesOnServerInjectedKeys",
			old:         stateConfigs,
			new:         `{"max.message.bytes":"1048576","retention.hours":"72"}`,
			expected:    true,
			description: "apply then plan must converge when the server reports extra keys",
		},
		{
			name:        "StillPlansRealChange",
			old:         stateConfigs,
			new:         `{"max.message.bytes":"10485760","retention.hours":"72"}`,
			expected:    false,
			description: "changing a declared key must still be planned",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceAliCloudAlikafkaTopic().Schema, map[string]interface{}{
				"configs": tc.old,
			})
			if result := suppress("configs", tc.old, tc.new, d); result != tc.expected {
				t.Errorf("%s: expected %v got %v", tc.description, tc.expected, result)
			}
		})
	}
}

var AliCloudAlikafkaTopicServerlessConfigsMap = map[string]string{
	"instance_id": CHECKSET,
	"create_time": CHECKSET,
	"status":      CHECKSET,
}

func AliCloudAlikafkaTopicServerlessConfigsDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_vswitches" "default" {
  		zone_id = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		vpc_id = data.alicloud_vswitches.default.vswitches.0.vpc_id
	}

	resource "alicloud_alikafka_instance" "default" {
  		name            = "tf-testAcc-${var.name}"
  		deploy_type     = "4"
  		instance_type   = "alikafka_serverless"
  		vswitch_id      = data.alicloud_vswitches.default.vswitches.0.id
  		spec_type       = "normal"
  		service_version = "3.3.1"
  		security_group  = alicloud_security_group.default.id
  		serverless_config {
    		reserved_publish_capacity   = 60
    		reserved_subscribe_capacity = 60
  		}
	}
`, name)
}

// TestAccAliCloudAlikafkaTopic_serverlessConfigs covers the configs attribute on a serverless
// instance, where the service reports configuration keys that were never declared
// (cloud.native.topic.type and replication-factor). The defect this case guards against is a plan
// that never converges rather than an apply that fails, so every applied configuration is followed
// by a plan-only step: those steps refresh from the service first and then require an empty plan.
func TestAccAliCloudAlikafkaTopic_serverlessConfigs(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_topic.default"
	ra := resourceAttrInit(resourceId, AliCloudAlikafkaTopicServerlessConfigsMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlikafkaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlikafkaTopic")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalikafka%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAlikafkaTopicServerlessConfigsDependence)

	// Keys accepted by the topic configuration update on a serverless instance. retention.ms is
	// deliberately absent: it is only accepted on a reserved instance with the Local storage engine.
	declaredConfigs := `{\"retention.hours\":\"72\",\"max.message.bytes\":\"1048576\",\"message.timestamp.type\":\"LogAppendTime\"}`
	updatedConfigs := `{\"retention.hours\":\"72\",\"max.message.bytes\":\"10485760\",\"message.timestamp.type\":\"LogAppendTime\"}`
	// The same configuration with a numeric value, which is what jsonencode produces for a number
	// while the service keeps reporting the value as a string.
	updatedConfigsAsNumber := `{\"retention.hours\":\"72\",\"max.message.bytes\":10485760,\"message.timestamp.type\":\"LogAppendTime\"}`

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
					"instance_id": "${alicloud_alikafka_instance.default.id}",
					"topic":       name,
					"remark":      name,
					"configs":     declaredConfigs,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id": CHECKSET,
						"topic":       name,
						"remark":      name,
						"configs":     REGEXMATCH + `"max\.message\.bytes":\s*"?1048576"?[,}]`,
					}),
				),
			},
			{
				// The reported configuration is a superset of the declared one, so re-planning the
				// configuration that was just applied must still produce an empty plan.
				Config: testAccConfig(map[string]interface{}{
					"configs": declaredConfigs,
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				// A real change on a declared key must still be planned and applied.
				Config: testAccConfig(map[string]interface{}{
					"configs": updatedConfigs,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"configs": REGEXMATCH + `"max\.message\.bytes":\s*"?10485760"?[,}]`,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configs": updatedConfigs,
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				// Values are compared through their literal text, so declaring the same value as a
				// number must not produce a diff either.
				Config: testAccConfig(map[string]interface{}{
					"configs": updatedConfigsAsNumber,
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
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
