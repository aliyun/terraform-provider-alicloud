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
	resource.AddTestSweepers("alicloud_alikafka_sasl_acl", &resource.Sweeper{
		Name: "alicloud_alikafka_sasl_acl",
		F:    testSweepAlikafkaSaslAcl,
	})
}

func testSweepAlikafkaSaslAcl(region string) error {
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
	instanceListReq.RegionId = defaultRegionToTest

	raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
		return alikafkaClient.GetInstanceList(instanceListReq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve alikafka instance in service list: %s", err)
	}

	instanceListResp, _ := raw.(*alikafka.GetInstanceListResponse)

	for _, v := range instanceListResp.InstanceList.InstanceVO {

		if v.ServiceStatus == 10 {
			log.Printf("[INFO] Skipping released alikafka instance id: %s ", v.InstanceId)
			continue
		}

		// Control the request rate.
		time.Sleep(time.Duration(400) * time.Millisecond)

		// Query users to delete
		userListReq := alikafka.CreateDescribeSaslUsersRequest()
		userListReq.InstanceId = v.InstanceId
		userListReq.RegionId = defaultRegionToTest

		saslUserRaw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DescribeSaslUsers(userListReq)
		})

		if err != nil {
			log.Printf("[ERROR] Failed to retrieve alikafka sasl users on instance (%s): %s", v.InstanceId, err)
			continue
		}

		saslUserListResp, _ := saslUserRaw.(*alikafka.DescribeSaslUsersResponse)
		var usersToDelete []string
		for _, saslUser := range saslUserListResp.SaslUserList.SaslUserVO {
			name := saslUser.Username
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping alikafka sasl username: %s ", name)
				continue
			}
			usersToDelete = append(usersToDelete, name)
		}
		if len(usersToDelete) == 0 {
			log.Printf("[INFO] Skipping by no users in alikafka instance id: %s ", v.InstanceId)
			continue
		}

		// Query All topic resource
		topicListReq := alikafka.CreateGetTopicListRequest()
		topicListReq.InstanceId = v.InstanceId
		topicListReq.RegionId = defaultRegionToTest

		topicRaw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.GetTopicList(topicListReq)
		})

		if err != nil {
			log.Printf("[ERROR] Failed to retrieve alikafka topics on instance (%s): %s", v.InstanceId, err)
			continue
		}
		topicListResp, _ := topicRaw.(*alikafka.GetTopicListResponse)
		topics := topicListResp.TopicList.TopicVO

		// Query all consumer groups
		consumerListReq := alikafka.CreateGetConsumerListRequest()
		consumerListReq.InstanceId = v.InstanceId
		consumerListReq.RegionId = defaultRegionToTest

		consumerRaw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.GetConsumerList(consumerListReq)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve alikafka consumer groups on instance (%s): %s", v.InstanceId, err)
			continue
		}
		consumerListResp, _ := consumerRaw.(*alikafka.GetConsumerListResponse)
		consumers := consumerListResp.ConsumerList.ConsumerVO

		// If there is no resource, skip
		if len(topics) == 0 && len(consumers) == 0 {
			log.Printf("[INFO] Skipping by no topics and consumers in alikafka instance id: %s ", v.InstanceId)
			continue
		}

		for _, username := range usersToDelete {

			for _, topic := range topics {

				deleteAcl(alikafkaService, v.InstanceId, username, "Topic", topic.Topic)
			}

			for _, consumer := range consumers {

				deleteAcl(alikafkaService, v.InstanceId, username, "Group", consumer.ConsumerId)
			}
		}
	}

	return nil
}

func deleteAcl(alikafkaService AlikafkaService, instanceId string, username string, resourceType string, resourceName string) {

	// Control the sasl username delete rate
	time.Sleep(time.Duration(400) * time.Millisecond)

	describeAclReq := alikafka.CreateDescribeAclsRequest()
	describeAclReq.InstanceId = instanceId
	describeAclReq.Username = username
	describeAclReq.RegionId = defaultRegionToTest
	describeAclReq.AclResourceName = resourceName
	describeAclReq.AclResourceType = resourceType

	raw, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
		return alikafkaClient.DescribeAcls(describeAclReq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to query alikafka acl username (%s), resourceType (%s), "+
			"resourceName (%s) instanceId (%s): %s", username, resourceType, resourceName, instanceId, err)
	}
	aclListResp, _ := raw.(*alikafka.DescribeAclsResponse)

	for _, kafkaAcl := range aclListResp.KafkaAclList.KafkaAclVO {

		if kafkaAcl.Username != username {
			continue
		}
		log.Printf("[INFO] delete alikafka acl: %s, ", kafkaAcl)
		deleteAclReq := alikafka.CreateDeleteAclRequest()
		deleteAclReq.RegionId = defaultRegionToTest
		deleteAclReq.InstanceId = instanceId
		deleteAclReq.Username = username
		deleteAclReq.AclResourceType = kafkaAcl.AclResourceType
		deleteAclReq.AclResourceName = kafkaAcl.AclResourceName
		deleteAclReq.AclResourcePatternType = kafkaAcl.AclResourcePatternType
		deleteAclReq.AclOperationType = kafkaAcl.AclOperationType

		_, err := alikafkaService.client.WithAlikafkaClient(func(alikafkaClient *alikafka.Client) (interface{}, error) {
			return alikafkaClient.DeleteAcl(deleteAclReq)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete alikafka acl acl (%s): %s", kafkaAcl, err)
		}
	}
}

func TestAccAlicloudAlikafkaSaslAcl_basic(t *testing.T) {

	var v *alikafka.KafkaAclVO
	resourceId := "alicloud_alikafka_sasl_acl.default"
	ra := resourceAttrInit(resourceId, alikafkaSaslAclBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkasaslaclbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaSaslAclConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAlikafkaAclEnable(t)
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
					"instance_id":               "${alicloud_alikafka_instance.default.id}",
					"username":                  "${alicloud_alikafka_sasl_user.default.username}",
					"acl_resource_type":         "Topic",
					"acl_resource_name":         "${alicloud_alikafka_topic.default.topic}",
					"acl_resource_pattern_type": "LITERAL",
					"acl_operation_type":        "Write",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"username":          fmt.Sprintf("tf-testacc-alikafkasaslaclbasic%v", rand),
						"acl_resource_name": fmt.Sprintf("tf-testacc-alikafkasaslaclbasic%v", rand),
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
					"acl_resource_pattern_type": "PREFIXED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_resource_pattern_type": "PREFIXED"}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"acl_operation_type": "Read",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_operation_type": "Read"}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"acl_resource_type": "Group",
					"acl_resource_name": "${alicloud_alikafka_consumer_group.default.consumer_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_resource_type": "Group",
						"acl_resource_name": fmt.Sprintf("tf-testacc-alikafkasaslaclbasic%v", rand)}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"username":           "${alicloud_alikafka_sasl_user.default.username}",
					"acl_resource_type":  "Topic",
					"acl_resource_name":  "${alicloud_alikafka_topic.default.topic}",
					"acl_operation_type": "Write",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"username":           fmt.Sprintf("tf-testacc-alikafkasaslaclbasic%v", rand),
						"acl_resource_type":  "Topic",
						"acl_resource_name":  fmt.Sprintf("tf-testacc-alikafkasaslaclbasic%v", rand),
						"acl_operation_type": "Write"}),
				),
			},
		},
	})

}

func TestAccAlicloudAlikafkaSaslAcl_multi(t *testing.T) {

	var v *alikafka.KafkaAclVO
	resourceId := "alicloud_alikafka_sasl_acl.default.1"
	ra := resourceAttrInit(resourceId, alikafkaSaslAclBasicMap)
	serviceFunc := func() interface{} {
		return &AlikafkaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-alikafkasaslaclbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlikafkaSaslAclConfigDependenceForMulti)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithAlikafkaAclEnable(t)
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
					"count":                     "2",
					"instance_id":               "${alicloud_alikafka_instance.default.id}",
					"username":                  "${alicloud_alikafka_sasl_user.default.username}",
					"acl_resource_type":         "Topic",
					"acl_resource_name":         "${alicloud_alikafka_topic.default.topic}",
					"acl_resource_pattern_type": "LITERAL",
					"acl_operation_type":        "${element(var.operation, count.index)}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"username":           fmt.Sprintf("tf-testacc-alikafkasaslaclbasic%v", rand),
						"acl_resource_name":  fmt.Sprintf("tf-testacc-alikafkasaslaclbasic%v", rand),
						"acl_operation_type": "Read",
					}),
				),
			},
		},
	})

}

func resourceAlikafkaSaslAclConfigDependence(name string) string {
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

resource "alicloud_alikafka_topic" "default" {
  instance_id = "${alicloud_alikafka_instance.default.id}"
  topic = "${var.name}"
  remark = "topic-remark"
}

resource "alicloud_alikafka_consumer_group" "default" {
  instance_id = "${alicloud_alikafka_instance.default.id}"
  consumer_id = "${var.name}"
}

resource "alicloud_alikafka_sasl_user" "default" {
  instance_id = "${alicloud_alikafka_instance.default.id}"
  username = "${var.name}"
  password = "password"
}
`, name)
}

func resourceAlikafkaSaslAclConfigDependenceForMulti(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

variable "operation" {
  default = ["Write", "Read"]
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

resource "alicloud_alikafka_topic" "default" {
  instance_id = "${alicloud_alikafka_instance.default.id}"
  topic = "${var.name}"
  remark = "topic-remark"
}

resource "alicloud_alikafka_sasl_user" "default" {
  instance_id = "${alicloud_alikafka_instance.default.id}"
  username = "${var.name}"
  password = "password"
}
`, name)
}

var alikafkaSaslAclBasicMap = map[string]string{
	"username":                  "${var.name}",
	"acl_resource_type":         "Topic",
	"acl_resource_name":         "${var.name}",
	"acl_resource_pattern_type": "LITERAL",
	"acl_operation_type":        "Write",
}
