package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/dxh031/ali_mns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_mns_topic", &resource.Sweeper{
		Name: "alicloud_mns_topic",
		F:    testSweepMnsTopics,
	})
}

func testSweepMnsTopics(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	conn := client.(*AliyunClient)

	topicManager, err := conn.MnsTopicManager()
	if err != nil {
		return fmt.Errorf(" Creating alicoudMNSTopicManager  error: %#v", err)
	}

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var topicAttrs []ali_mns.TopicAttribute
	for _, namePrefix := range prefixes {
		for {
			var nextMaker string
			topicDetails, err1 := topicManager.ListTopicDetail(nextMaker, 1000, namePrefix)
			if err1 != nil {
				return fmt.Errorf("get topicDetails  error: %#v", err)
			}
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
		err = topicManager.DeleteTopic(name)
		if err != nil {
			log.Printf("[ERROR] Failed to delete mns topic (%s (%s)): %s", topicAttr.TopicName, topicAttr.TopicName, err)
		}
	}

	return nil
}

func TestAccResourceAlicloudMNSTopic_basic(t *testing.T) {

	var attr ali_mns.TopicAttribute

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMNSTopicDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccMNSTopicConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccMNSTopicExist("alicloud_mns_topic.topic", &attr),
					resource.TestCheckResourceAttr("alicloud_mns_topic.topic", "name", "tf-testAccMNSTopicConfig"),
				),
			},
			{
				Config: testAccMNSTopicConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccMNSTopicExist("alicloud_mns_topic.topic", &attr),
					resource.TestCheckResourceAttr("alicloud_mns_topic.topic", "name", "tf-testAccMNSTopicConfig"),
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

		client := testAccProvider.Meta().(*AliyunClient)

		topicManager, err := client.MnsTopicManager()
		if err != nil {
			return fmt.Errorf(" Creating alicoudMNSTopicManager  error: %#v", err)
		}
		instance, err := topicManager.GetTopicAttributes(rs.Primary.ID)
		if err != nil {
			return err
		}
		if instance.TopicName != rs.Primary.ID {
			return fmt.Errorf("mns topic %s not found", n)
		}
		*attr = instance
		return nil
	}

}

func testAccCheckMNSTopicDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_mns_topic" {
			continue
		}
		topicManager, err := client.MnsTopicManager()
		if err != nil {
			return fmt.Errorf(" Creating alicoudMNSTopicManager  error: %#v", err)
		}

		if _, err := topicManager.GetTopicAttributes(rs.Primary.ID); err != nil {
			if TopicNotExistFunc(err) {
				continue
			}
			return err
		}

		return fmt.Errorf("MNS Topic %s still exist", rs.Primary.ID)
	}

	return nil
}

const testAccMNSTopicConfig = `
variable "name" {
	default = "tf-testAccMNSTopicConfig"
}
resource "alicloud_mns_topic" "topic"{
	name="${var.name}"
}`

const testAccMNSTopicConfigUpdate = `
variable "name" {
	default = "tf-testAccMNSTopicConfig"
}
resource "alicloud_mns_topic" "topic"{
	name="${var.name}"
	maximum_message_size=12357
	logging_enabled=true
}`
