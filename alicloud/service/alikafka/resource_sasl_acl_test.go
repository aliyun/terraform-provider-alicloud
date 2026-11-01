package alikafka_test

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	alikafkasdk "github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	alikafka "github.com/aliyun/terraform-provider-alicloud/alicloud/service/alikafka"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_alikafka_sasl_acl", &resource.Sweeper{
		Name: "alicloud_alikafka_sasl_acl",
		F:    testSweepAlikafkaSaslAcl,
	})
}

func testSweepAlikafkaSaslAcl(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return errs.WrapErrorf(err, "error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	instanceListReq := alikafkasdk.CreateGetInstanceListRequest()
	instanceListReq.RegionId = acctest.DefaultRegionToTest

	raw, err := client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
		return alikafkaClient.GetInstanceList(instanceListReq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve alikafka instance in service list: %s", err)
	}

	instanceListResp, _ := raw.(*alikafkasdk.GetInstanceListResponse)

	for _, v := range instanceListResp.InstanceList.InstanceVO {

		if v.ServiceStatus == 10 {
			log.Printf("[INFO] Skipping released alikafka instance id: %s ", v.InstanceId)
			continue
		}

		// Control the request rate.
		time.Sleep(time.Duration(400) * time.Millisecond)

		// Query users to delete
		userListReq := alikafkasdk.CreateDescribeSaslUsersRequest()
		userListReq.InstanceId = v.InstanceId
		userListReq.RegionId = acctest.DefaultRegionToTest

		saslUserRaw, err := client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
			return alikafkaClient.DescribeSaslUsers(userListReq)
		})

		if err != nil {
			log.Printf("[ERROR] Failed to retrieve alikafka sasl users on instance (%s): %s", v.InstanceId, err)
			continue
		}

		saslUserListResp, _ := saslUserRaw.(*alikafkasdk.DescribeSaslUsersResponse)
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
		topicListReq := alikafkasdk.CreateGetTopicListRequest()
		topicListReq.InstanceId = v.InstanceId
		topicListReq.RegionId = acctest.DefaultRegionToTest

		topicRaw, err := client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
			return alikafkaClient.GetTopicList(topicListReq)
		})

		if err != nil {
			log.Printf("[ERROR] Failed to retrieve alikafka topics on instance (%s): %s", v.InstanceId, err)
			continue
		}
		topicListResp, _ := topicRaw.(*alikafkasdk.GetTopicListResponse)
		topics := topicListResp.TopicList.TopicVO

		// Query all consumer groups
		consumerListReq := alikafkasdk.CreateGetConsumerListRequest()
		consumerListReq.InstanceId = v.InstanceId
		consumerListReq.RegionId = acctest.DefaultRegionToTest

		consumerRaw, err := client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
			return alikafkaClient.GetConsumerList(consumerListReq)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve alikafka consumer groups on instance (%s): %s", v.InstanceId, err)
			continue
		}
		consumerListResp, _ := consumerRaw.(*alikafkasdk.GetConsumerListResponse)
		consumers := consumerListResp.ConsumerList.ConsumerVO

		// If there is no resource, skip
		if len(topics) == 0 && len(consumers) == 0 {
			log.Printf("[INFO] Skipping by no topics and consumers in alikafka instance id: %s ", v.InstanceId)
			continue
		}

		for _, username := range usersToDelete {

			for _, topic := range topics {

				deleteAcl(client, v.InstanceId, username, "Topic", topic.Topic)
			}

			for _, consumer := range consumers {

				deleteAcl(client, v.InstanceId, username, "Group", consumer.ConsumerId)
			}
		}
	}

	return nil
}

func deleteAcl(client *connectivity.AliyunClient, instanceId string, username string, resourceType string, resourceName string) {

	// Control the sasl username delete rate
	time.Sleep(time.Duration(400) * time.Millisecond)

	describeAclReq := alikafkasdk.CreateDescribeAclsRequest()
	describeAclReq.InstanceId = instanceId
	describeAclReq.Username = username
	describeAclReq.RegionId = acctest.DefaultRegionToTest
	describeAclReq.AclResourceName = resourceName
	describeAclReq.AclResourceType = resourceType

	raw, err := client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
		return alikafkaClient.DescribeAcls(describeAclReq)
	})
	if err != nil {
		log.Printf("[ERROR] Failed to query alikafka acl username (%s), resourceType (%s), "+
			"resourceName (%s) instanceId (%s): %s", username, resourceType, resourceName, instanceId, err)
	}
	aclListResp, _ := raw.(*alikafkasdk.DescribeAclsResponse)

	for _, kafkaAcl := range aclListResp.KafkaAclList.KafkaAclVO {

		if kafkaAcl.Username != username {
			continue
		}
		log.Printf("[INFO] delete alikafka acl: %s, ", kafkaAcl)
		deleteAclReq := alikafkasdk.CreateDeleteAclRequest()
		deleteAclReq.RegionId = acctest.DefaultRegionToTest
		deleteAclReq.InstanceId = instanceId
		deleteAclReq.Username = username
		deleteAclReq.AclResourceType = kafkaAcl.AclResourceType
		deleteAclReq.AclResourceName = kafkaAcl.AclResourceName
		deleteAclReq.AclResourcePatternType = kafkaAcl.AclResourcePatternType
		deleteAclReq.AclOperationType = kafkaAcl.AclOperationType

		_, err := client.WithAlikafkaClient(func(alikafkaClient *alikafkasdk.Client) (interface{}, error) {
			return alikafkaClient.DeleteAcl(deleteAclReq)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete alikafka acl acl (%s): %s", kafkaAcl, err)
		}
	}
}

// Test Alikafka SaslAcl. >>> Resource test cases, automatically generated.
// Case v3 创建ACL生命周期 12413
func TestAccAliCloudAlikafkaSaslAcl_basic12413(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_sasl_acl.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAlikafkaSaslAclMap12413)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alikafka.NewAlikafkaServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlikafkaSaslAcl")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalikafka%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAlikafkaSaslAclBasicDependence12413)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_operation_type":        "Write",
					"acl_resource_name":         "${alicloud_alikafka_topic.default.topic}",
					"acl_resource_pattern_type": "LITERAL",
					"acl_resource_type":         "Topic",
					"instance_id":               "${alicloud_alikafka_sasl_user.default.instance_id}",
					"username":                  "${alicloud_alikafka_sasl_user.default.username}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_operation_type":        "Write",
						"acl_resource_name":         acctest.CHECKSET,
						"acl_resource_pattern_type": "LITERAL",
						"acl_resource_type":         "Topic",
						"instance_id":               acctest.CHECKSET,
						"username":                  acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl_operation_types"},
			},
		},
	})
}

func TestAccAliCloudAlikafkaSaslAcl_basic12413_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_sasl_acl.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAlikafkaSaslAclMap12413)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alikafka.NewAlikafkaServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlikafkaSaslAcl")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalikafka%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAlikafkaSaslAclBasicDependence12413)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_operation_type":        "Write",
					"acl_resource_name":         "${alicloud_alikafka_topic.default.topic}",
					"acl_resource_pattern_type": "LITERAL",
					"acl_resource_type":         "Topic",
					"instance_id":               "${alicloud_alikafka_sasl_user.default.instance_id}",
					"username":                  "${alicloud_alikafka_sasl_user.default.username}",
					"acl_permission_type":       "DENY",
					"host":                      "127.0.0.1",
					"acl_operation_types":       "Write",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_operation_type":        "Write",
						"acl_resource_name":         acctest.CHECKSET,
						"acl_resource_pattern_type": "LITERAL",
						"acl_resource_type":         "Topic",
						"instance_id":               acctest.CHECKSET,
						"username":                  acctest.CHECKSET,
						"acl_permission_type":       "DENY",
						"host":                      "127.0.0.1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl_operation_types"},
			},
		},
	})
}

var AliCloudAlikafkaSaslAclMap12413 = map[string]string{
	"acl_permission_type": acctest.CHECKSET,
	"host":                acctest.CHECKSET,
}

func AliCloudAlikafkaSaslAclBasicDependence12413(name string) string {
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
  		deploy_type     = "4"
  		instance_type   = "alikafka_serverless"
  		vswitch_id      = data.alicloud_vswitches.default.vswitches.0.id
  		spec_type       = "normal"
  		service_version = "3.3.1"
  		security_group  = alicloud_security_group.default.id
  		config          = "{\"enable.acl\":\"true\"}"
  		serverless_config {
    		reserved_publish_capacity   = 60
    		reserved_subscribe_capacity = 60
  		}
	}

	resource "alicloud_alikafka_sasl_user" "default" {
  		instance_id = alicloud_alikafka_instance.default.id
  		username    = var.name
  		type        = "scram"
  		password    = "YourPassword1234!"
	}

	resource "alicloud_alikafka_topic" "default" {
  		instance_id = alicloud_alikafka_instance.default.id
  		topic       = var.name
  		remark      = var.name
	}
`, name)
}

// Case v2 创建ACL生命周期 12414
func TestAccAliCloudAlikafkaSaslAcl_basic12414(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_sasl_acl.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAlikafkaSaslAclMap12414)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alikafka.NewAlikafkaServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlikafkaSaslAcl")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalikafka%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAlikafkaSaslAclBasicDependence12414)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"acl_operation_type":        "Write",
					"acl_resource_name":         "${alicloud_alikafka_topic.default.topic}",
					"acl_resource_pattern_type": "LITERAL",
					"acl_resource_type":         "Topic",
					"instance_id":               "${alicloud_alikafka_sasl_user.default.instance_id}",
					"username":                  "${alicloud_alikafka_sasl_user.default.username}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_operation_type":        "Write",
						"acl_resource_name":         acctest.CHECKSET,
						"acl_resource_pattern_type": "LITERAL",
						"acl_resource_type":         "Topic",
						"instance_id":               acctest.CHECKSET,
						"username":                  acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"acl_operation_types"},
			},
		},
	})
}

func TestAccAliCloudAlikafkaSaslAcl_multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alikafka_sasl_acl.default.2"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudAlikafkaSaslAclMap12414)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return alikafka.NewAlikafkaServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAlikafkaSaslAcl")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccalikafka%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudAlikafkaSaslAclBasicDependence12414)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":                     "3",
					"acl_operation_type":        "${element(var.operation, count.index)}",
					"acl_resource_name":         "${alicloud_alikafka_topic.default.topic}",
					"acl_resource_pattern_type": "LITERAL",
					"acl_resource_type":         "Topic",
					"instance_id":               "${alicloud_alikafka_sasl_user.default.instance_id}",
					"username":                  "${alicloud_alikafka_sasl_user.default.username}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_operation_type":        "Describe",
						"acl_resource_name":         acctest.CHECKSET,
						"acl_resource_pattern_type": "LITERAL",
						"acl_resource_type":         "Topic",
						"instance_id":               acctest.CHECKSET,
						"username":                  acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

var AliCloudAlikafkaSaslAclMap12414 = map[string]string{
	"host": acctest.CHECKSET,
}

func AliCloudAlikafkaSaslAclBasicDependence12414(name string) string {
	return fmt.Sprintf(`
	variable "name" {
    	default = "%s"
	}

	variable "operation" {
  		default = ["Write", "Read", "Describe"]
	}

	data "alicloud_alikafka_instances" "default" {
	}

	resource "alicloud_alikafka_sasl_user" "default" {
  		instance_id = data.alicloud_alikafka_instances.default.instances.0.id
  		username    = var.name
  		password    = "YourPassword1234!"
	}

	resource "alicloud_alikafka_topic" "default" {
  		instance_id = data.alicloud_alikafka_instances.default.instances.0.id
  		topic       = var.name
  		remark      = var.name
	}
`, name)
}

// Test Alikafka SaslAcl. <<< Resource test cases, automatically generated.
