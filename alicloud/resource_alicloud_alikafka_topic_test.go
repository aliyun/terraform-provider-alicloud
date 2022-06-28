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
		return WrapErrorf(err, "error getting Alicloud client.")
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

func TestAccAlicloudAlikafkaTopic_basic(t *testing.T) {

	var v *alikafka.TopicVO
	resourceId := "alicloud_alikafka_topic.default"
	ra := resourceAttrInit(resourceId, alikafkaTopicBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkatopicbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaTopicConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.AlikafkaSupportedRegions)
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":   "${alicloud_alikafka_instance.default.id}",
					"topic":         "${var.name}",
					"local_topic":   "false",
					"compact_topic": "false",
					"partition_num": "6",
					"remark":        "alicloud_alikafka_topic_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topic":         fmt.Sprintf("tf-testacc-alikafkatopicbasic%v", rand),
						"local_topic":   "false",
						"compact_topic": "false",
						"partition_num": "6",
						"remark":        "alicloud_alikafka_topic_remark",
					}),
				),
			},

			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"partition_num": "9",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"partition_num": "9"}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"topic":  "tf-testacc-alicloud_alikafka_default_topic_change",
					"remark": "modified remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topic":  "tf-testacc-alicloud_alikafka_default_topic_change",
						"remark": "modified remark"}),
				),
			},

			// alicloud_alikafka_instance only support create post pay instance.
			// Post pay instance does not support create local or compact topic, so skip the following two test case temporarily.
			//{
			//	SkipFunc: shouldSkipLocalAndCompact("${alicloud_alikafka_instance.default.id}"),
			//	Config: testAccConfig(map[string]interface{}{
			//		"local_topic": "true",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"local_topic": "true",
			//		}),
			//	),
			//},

			//{
			//	SkipFunc: shouldSkipLocalAndCompact("${alicloud_alikafka_instance.default.id}"),
			//	Config: testAccConfig(map[string]interface{}{
			//		"compact_topic": "true",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"compact_topic": "true",
			//		}),
			//	),
			//},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "acceptance test",
						"Updated": "TF",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "3",
						"tags.Created": "TF",
						"tags.For":     "acceptance test",
						"tags.Updated": "TF",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"topic":         "${var.name}",
					"local_topic":   "false",
					"compact_topic": "false",
					"partition_num": "12",
					"remark":        "alicloud_alikafka_topic_remark",
					"tags":          REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topic":         fmt.Sprintf("tf-testacc-alikafkatopicbasic%v", rand),
						"local_topic":   "false",
						"compact_topic": "false",
						"partition_num": "12",
						"remark":        "alicloud_alikafka_topic_remark",
						"tags.%":        REMOVEKEY,
						"tags.Created":  REMOVEKEY,
						"tags.For":      REMOVEKEY,
						"tags.Updated":  REMOVEKEY,
					}),
				),
			},
		},
	})

}

func TestAccAlicloudAlikafkaTopic_multi(t *testing.T) {

	var v *alikafka.TopicVO
	resourceId := "alicloud_alikafka_topic.default.4"
	ra := resourceAttrInit(resourceId, alikafkaTopicBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkatopicbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaTopicConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.AlikafkaSupportedRegions)
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":         "5",
					"instance_id":   "${alicloud_alikafka_instance.default.id}",
					"topic":         "${var.name}-${count.index}",
					"local_topic":   "false",
					"compact_topic": "false",
					"partition_num": "6",
					"remark":        "alicloud_alikafka_topic_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topic":         fmt.Sprintf("tf-testacc-alikafkatopicbasic%v-4", rand),
						"local_topic":   "false",
						"compact_topic": "false",
						"partition_num": "6",
						"remark":        "alicloud_alikafka_topic_remark",
					}),
				),
			},
		},
	})

}

func resourceAlikafkaTopicConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

data "alicloud_vpcs" "default" {
 name_regex = "^default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_alikafka_instance" "default" {
  name = "${var.name}"
  topic_quota = "50"
  disk_type = "1"
  disk_size = "500"
  deploy_type = "5"
  io_max = "20"
  vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
  security_group = alicloud_security_group.default.id
}
`, name)
}

var alikafkaTopicBasicMap = map[string]string{
	"topic":         "${var.name}",
	"local_topic":   "false",
	"compact_topic": "false",
	"partition_num": "12",
	"remark":        "alicloud_alikafka_topic_remark",
}
