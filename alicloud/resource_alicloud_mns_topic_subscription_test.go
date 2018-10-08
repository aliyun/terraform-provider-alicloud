package alicloud

import (
	"fmt"
	"testing"

	"github.com/dxh031/ali_mns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceAlicloudMNSTopicSubscription_basic(t *testing.T) {

	var attr ali_mns.TopicAttribute

	var subscriptionAttr ali_mns.SubscriptionAttribute

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMNSTopicSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMNSTopicSubscriptionConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccMNSTopicExist("alicloud_mns_topic.topic", &attr),
					testAccMNSTopicSubscriptionExist("alicloud_mns_topic_subscription.subscription", &subscriptionAttr),
					resource.TestCheckResourceAttr("alicloud_mns_topic_subscription.subscription", "name", "tf-testAccMNSTopicSubscriptionConfig"),
					resource.TestCheckResourceAttr("alicloud_mns_topic_subscription.subscription", "endpoint", "http://www.test.com/test"),
					resource.TestCheckResourceAttr("alicloud_mns_topic_subscription.subscription", "notify_content_format", "SIMPLIFIED"),
				),
			},
			{

				Config: testAccMNSTopicSubscriptionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccMNSTopicExist("alicloud_mns_topic.topic", &attr),
					testAccMNSTopicSubscriptionExist("alicloud_mns_topic_subscription.subscription", &subscriptionAttr),
					resource.TestCheckResourceAttr("alicloud_mns_topic_subscription.subscription", "name", "tf-testAccMNSTopicSubscriptionConfig"),
					resource.TestCheckResourceAttr("alicloud_mns_topic_subscription.subscription", "notify_strategy", "EXPONENTIAL_DECAY_RETRY"),
					resource.TestCheckResourceAttr("alicloud_mns_topic_subscription.subscription", "endpoint", "http://www.test.com/test"),
					resource.TestCheckResourceAttr("alicloud_mns_topic_subscription.subscription", "notify_content_format", "SIMPLIFIED"),
				),
			},
		},
	})
}

func testAccMNSTopicSubscriptionExist(n string, attr *ali_mns.SubscriptionAttribute) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No MNSTopicSubscription ID is set")
		}
		client := testAccProvider.Meta().(*AliyunClient)
		topicName, name := GetTopicNameAndSubscriptionName(rs.Primary.ID)
		subscriptionManager, err := client.MnsSubscriptionManager(topicName)
		if err != nil {
			return fmt.Errorf("Creating mns subscription client  error: %#v", err)
		}
		instance, err := subscriptionManager.GetSubscriptionAttributes(name)

		if err != nil {
			return err
		}
		if instance.SubscriptionName != name {
			return fmt.Errorf("mns subscription %s not found", n)
		}
		*attr = instance
		return nil
	}

}

func testAccCheckMNSTopicSubscriptionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_mns_topic_subscription" {
			continue
		}
		topicName, name := GetTopicNameAndSubscriptionName(rs.Primary.ID)
		subscriptionManager, err := client.MnsSubscriptionManager(topicName)
		if err != nil {
			return fmt.Errorf("Creating mns subscription client  error: %#v", err)
		}
		if _, err := subscriptionManager.GetSubscriptionAttributes(name); err != nil {
			if SubscriptionNotExistFunc(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("MNS topic subscription %s still exist", name)
	}

	return nil
}

const testAccMNSTopicSubscriptionConfig = `variable "name" {
	default = "tf-testAccMNSTopicConfig"
}
variable "subscriptionName" {
	default = "tf-testAccMNSTopicSubscriptionConfig"
}
resource "alicloud_mns_topic" "topic"{
	name="${var.name}"
}
resource "alicloud_mns_topic_subscription" "subscription"{
	topic_name="${alicloud_mns_topic.topic.name}"
	name="${var.subscriptionName}"
	endpoint="http://www.test.com/test"
	notify_strategy="BACKOFF_RETRY"
	notify_content_format="SIMPLIFIED"
}`

const testAccMNSTopicSubscriptionConfigUpdate = `variable "name" {
	default = "tf-testAccMNSTopicConfig"
}
variable "subscriptionName" {
	default = "tf-testAccMNSTopicSubscriptionConfig"
}
resource "alicloud_mns_topic" "topic"{
	name="${var.name}"
	maximum_message_size=12357
	logging_enabled=true
}
resource "alicloud_mns_topic_subscription" "subscription"{
	topic_name="${alicloud_mns_topic.topic.name}"
	name="${var.subscriptionName}"
	endpoint="http://www.test.com/test"
	notify_strategy="EXPONENTIAL_DECAY_RETRY"
	notify_content_format="SIMPLIFIED"
}`
