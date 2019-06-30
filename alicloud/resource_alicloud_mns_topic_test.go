package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/dxh031/ali_mns"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_mns_topic", &resource.Sweeper{
		Name: "alicloud_mns_topic",
		F:    testSweepMnsTopics,
	})
}

func testSweepMnsTopics(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf-testacc",
	}

	var topicAttrs []ali_mns.TopicAttribute
	for _, namePrefix := range prefixes {
		for {
			var nextMaker string
			raw, err := client.WithMnsTopicManager(func(topicManager ali_mns.AliTopicManager) (interface{}, error) {
				return topicManager.ListTopicDetail(nextMaker, 1000, namePrefix)
			})
			if err != nil {
				return fmt.Errorf("get topicDetails  error: %#v", err)
			}
			topicDetails, _ := raw.(ali_mns.TopicDetails)
			for _, attr := range topicDetails.Attrs {
				topicAttrs = append(topicAttrs, attr)
			}
			nextMaker = topicDetails.NextMarker
			if nextMaker == "" {
				break
			}
		}
	}
	for _, topicAttr := range topicAttrs {
		name := topicAttr.TopicName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping mns topic : %s ", name)
			continue
		}
		log.Printf("[INFO] delete  mns topic : %s ", name)
		_, err := client.WithMnsTopicManager(func(topicManager ali_mns.AliTopicManager) (interface{}, error) {
			return nil, topicManager.DeleteTopic(name)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete mns topic (%s (%s)): %s", topicAttr.TopicName, topicAttr.TopicName, err)
		}
	}

	return nil
}

func TestAccAlicloudMnsTopic_basic(t *testing.T) {

	var attr ali_mns.TopicAttribute
	rand := acctest.RandIntRange(10000, 999999)
	resourceId := "alicloud_mns_topic.topic"

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		Providers:     testAccProviders,
		IDRefreshName: resourceId,
		CheckDestroy:  testAccCheckMNSTopicDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccMNSTopicConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccMNSTopicExist("alicloud_mns_topic.topic", &attr),
					resource.TestCheckResourceAttr("alicloud_mns_topic.topic", "name", fmt.Sprintf("tf-testAccMNSTopicConfig-%d", rand)),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMNSTopicConfigUpdate(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccMNSTopicExist("alicloud_mns_topic.topic", &attr),
					resource.TestCheckResourceAttr("alicloud_mns_topic.topic", "name", fmt.Sprintf("tf-testAccMNSTopicConfig-%d", rand)),
					resource.TestCheckResourceAttr("alicloud_mns_topic.topic", "maximum_message_size", "12357"),
					resource.TestCheckResourceAttr("alicloud_mns_topic.topic", "logging_enabled", "true"),
				),
			},
		},
	})
}

func testAccMNSTopicExist(n string, attr *ali_mns.TopicAttribute) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No MNSTopic ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		raw, err := client.WithMnsTopicManager(func(topicManager ali_mns.AliTopicManager) (interface{}, error) {
			return topicManager.GetTopicAttributes(rs.Primary.ID)
		})
		if err != nil {
			return err
		}
		instance, _ := raw.(ali_mns.TopicAttribute)
		if instance.TopicName != rs.Primary.ID {
			return fmt.Errorf("mns topic %s not found", n)
		}
		*attr = instance
		return nil
	}

}

func testAccCheckMNSTopicDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	mnsService := MnsService{}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_mns_topic" {
			continue
		}

		_, err := client.WithMnsTopicManager(func(topicManager ali_mns.AliTopicManager) (interface{}, error) {
			return topicManager.GetTopicAttributes(rs.Primary.ID)
		})
		if err != nil {
			if mnsService.TopicNotExistFunc(err) {
				continue
			}
			return err
		}

		return fmt.Errorf("MNS Topic %s still exist", rs.Primary.ID)
	}

	return nil
}

func testAccMNSTopicConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccMNSTopicConfig-%d"
	}
	resource "alicloud_mns_topic" "topic"{
		name="${var.name}"
	}`, rand)
}

func testAccMNSTopicConfigUpdate(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccMNSTopicConfig-%d"
	}
	resource "alicloud_mns_topic" "topic"{
		name="${var.name}"
		maximum_message_size=12357
		logging_enabled=true
	}`, rand)
}
