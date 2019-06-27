package alicloud

import (
	"fmt"
	"testing"

	"github.com/dxh031/ali_mns"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudMnsTopicSubscription_basic(t *testing.T) {

	var attr ali_mns.TopicAttribute
	var subscriptionAttr ali_mns.SubscriptionAttribute
	resourceId := "alicloud_mns_topic_subscription.subscription"

	rand := acctest.RandIntRange(10000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		Providers:     testAccProviders,
		IDRefreshName: resourceId,
		CheckDestroy:  testAccCheckMNSTopicSubscriptionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMNSTopicSubscriptionConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccMNSTopicExist("alicloud_mns_topic.topic", &attr),
					testAccMNSTopicSubscriptionExist("alicloud_mns_topic_subscription.subscription", &subscriptionAttr),
					resource.TestCheckResourceAttr("alicloud_mns_topic_subscription.subscription", "name", fmt.Sprintf("tf-testAccMNSTopicSubscriptionConfig-%d", rand)),
					resource.TestCheckResourceAttr("alicloud_mns_topic_subscription.subscription", "endpoint", "http://www.test.com/test"),
					resource.TestCheckResourceAttr("alicloud_mns_topic_subscription.subscription", "notify_content_format", "SIMPLIFIED"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{

				Config: testAccMNSTopicSubscriptionConfigUpdate(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccMNSTopicExist("alicloud_mns_topic.topic", &attr),
					testAccMNSTopicSubscriptionExist("alicloud_mns_topic_subscription.subscription", &subscriptionAttr),
					resource.TestCheckResourceAttr("alicloud_mns_topic_subscription.subscription", "name", fmt.Sprintf("tf-testAccMNSTopicSubscriptionConfig-%d", rand)),
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
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		mnsService := MnsService{}
		topicName, name := mnsService.GetTopicNameAndSubscriptionName(rs.Primary.ID)

		raw, err := client.WithMnsSubscriptionManagerByTopicName(topicName, func(subscriptionManager ali_mns.AliMNSTopic) (interface{}, error) {
			return subscriptionManager.GetSubscriptionAttributes(name)
		})
		if err != nil {
			return err
		}
		instance, _ := raw.(ali_mns.SubscriptionAttribute)
		if instance.SubscriptionName != name {
			return fmt.Errorf("mns subscription %s not found", n)
		}
		*attr = instance
		return nil
	}

}

func testAccCheckMNSTopicSubscriptionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	mnsService := MnsService{}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_mns_topic_subscription" {
			continue
		}
		topicName, name := mnsService.GetTopicNameAndSubscriptionName(rs.Primary.ID)
		_, err := client.WithMnsSubscriptionManagerByTopicName(topicName, func(subscriptionManager ali_mns.AliMNSTopic) (interface{}, error) {
			return subscriptionManager.GetSubscriptionAttributes(name)
		})
		if err != nil {
			if mnsService.TopicNotExistFunc(err) || mnsService.SubscriptionNotExistFunc(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("MNS topic subscription %s still exist", name)
	}

	return nil
}

func testAccMNSTopicSubscriptionConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccMNSTopicConfig"
	}
	variable "subscriptionName" {
		default = "tf-testAccMNSTopicSubscriptionConfig-%d"
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
	}`, rand)
}

func testAccMNSTopicSubscriptionConfigUpdate(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccMNSTopicConfig"
	}
	variable "subscriptionName" {
		default = "tf-testAccMNSTopicSubscriptionConfig-%d"
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
	}`, rand)
}
