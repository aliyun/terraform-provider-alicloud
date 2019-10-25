package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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

		wait := incrementalWait(3*time.Second, 5*time.Second)
		var raw interface{}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err = alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
				return alikafkaClient.GetTopicList(request)
			})
			if err != nil {
				if IsExceptedError(err, AlikafkaThrottlingUser) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			return nil
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
					"topic":         "${var.topic}",
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
					"topic":         "${var.topic}",
					"local_topic":   "false",
					"compact_topic": "false",
					"partition_num": "12",
					"remark":        "alicloud_alikafka_topic_remark",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"topic":         fmt.Sprintf("tf-testacc-alikafkatopicbasic%v", rand),
						"local_topic":   "false",
						"compact_topic": "false",
						"partition_num": "12",
						"remark":        "alicloud_alikafka_topic_remark",
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
					"topic":         "${var.topic}-${count.index}",
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

type skipLocalAndCompactFunc func() (bool, error)

func shouldSkipLocalAndCompact(instanceId string) skipLocalAndCompactFunc {

	return func() (bool, error) {

		rawClient, err := sharedClientForRegion(defaultRegionToTest)
		if err != nil {
			return false, err
		}
		client := rawClient.(*connectivity.AliyunClient)
		alikafkaService := AlikafkaService{client}

		instance, err := alikafkaService.DescribeAlikafkaInstance(instanceId)
		if err != nil {
			return false, err
		}

		supportLocalAndCompactTopic := false
		for _, v := range instance.UpgradeServiceDetailInfo.UpgradeServiceDetailInfoVO {
			if v.Current2OpenSourceVersion >= "2." {
				supportLocalAndCompactTopic = true
				break
			}
		}
		return !supportLocalAndCompactTopic, nil
	}
}

func resourceAlikafkaTopicConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "topic" {
 			default = "%v"
		}

		data "alicloud_zones" "default" {
			available_resource_creation= "VSwitch"
		}
		resource "alicloud_vpc" "default" {
		  cidr_block = "172.16.0.0/12"
		}
		
		resource "alicloud_vswitch" "default" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.0.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		}

		resource "alicloud_alikafka_instance" "default" {
          name = "tf-testacc-alikafkainstance"
		  topic_quota = "50"
		  disk_type = "1"
		  disk_size = "500"
		  deploy_type = "5"
		  io_max = "20"
          vswitch_id = "${alicloud_vswitch.default.id}"
		}
		`, name)
}

var alikafkaTopicBasicMap = map[string]string{
	"topic":         "${var.topic}",
	"local_topic":   "false",
	"compact_topic": "false",
	"partition_num": "12",
	"remark":        "alicloud_alikafka_topic_remark",
}
