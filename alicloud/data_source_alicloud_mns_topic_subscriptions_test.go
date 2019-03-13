package alicloud

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudMnsTopicSubscriptionDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudMNSTopicSubscriptionDataSourceConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_mns_topic_subscriptions.subscriptions"),
					resource.TestCheckResourceAttr("data.alicloud_mns_topic_subscriptions.subscriptions", "subscriptions.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_mns_topic_subscriptions.subscriptions", "subscriptions.0.name", fmt.Sprintf("tf-testAccMNSTopicSubscriptionConfig-%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_mns_topic_subscriptions.subscriptions", "subscriptions.0.endpoint", "http://www.test.com/test"),
					resource.TestCheckResourceAttr("data.alicloud_mns_topic_subscriptions.subscriptions", "subscriptions.0.notify_strategy", "EXPONENTIAL_DECAY_RETRY"),
					resource.TestCheckResourceAttr("data.alicloud_mns_topic_subscriptions.subscriptions", "subscriptions.0.notify_content_format", "SIMPLIFIED"),
				),
			},
		},
	})
}

func TestAccAlicloudMnsTopicSubscriptionDataSourceEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudMNSTopicSubscriptionDataSourceEmpty(acctest.RandIntRange(10000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_mns_topic_subscriptions.subscriptions"),
					resource.TestCheckResourceAttr("data.alicloud_mns_topic_subscriptions.subscriptions", "subscriptions.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_mns_topic_subscriptions.subscriptions", "subscriptions.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_mns_topic_subscriptions.subscriptions", "subscriptions.0.endpoint"),
					resource.TestCheckNoResourceAttr("data.alicloud_mns_topic_subscriptions.subscriptions", "subscriptions.0.notify_strategy"),
					resource.TestCheckNoResourceAttr("data.alicloud_mns_topic_subscriptions.subscriptions", "subscriptions.0.notify_content_format"),
				),
			},
		},
	})
}

func testAccCheckAlicloudMNSTopicSubscriptionDataSourceConfig(rand int) string {
	return fmt.Sprintf(`
	data "alicloud_mns_topic_subscriptions" "subscriptions" {
	  topic_name="${alicloud_mns_topic.topic.name}"
	  name_prefix = "${alicloud_mns_topic_subscription.subscription.name}"
	}

	resource "alicloud_mns_topic" "topic"{
		name="tf-testAccMNSTopicConfig-%d"
		maximum_message_size=12357
		logging_enabled=true
	}

	resource "alicloud_mns_topic_subscription" "subscription"{
		topic_name="${alicloud_mns_topic.topic.name}"
		name="tf-testAccMNSTopicSubscriptionConfig-%d"
		endpoint="http://www.test.com/test"
		notify_strategy="EXPONENTIAL_DECAY_RETRY"
		notify_content_format="SIMPLIFIED"
	}`, rand, rand)
}

func testAccCheckAlicloudMNSTopicSubscriptionDataSourceEmpty(rand int) string {
	return fmt.Sprintf(`
	data "alicloud_mns_topic_subscriptions" "subscriptions" {
	  topic_name="${alicloud_mns_topic.topic.name}"
	  name_prefix = "tf-testacc-fake-name"
	}

	resource "alicloud_mns_topic" "topic"{
		name="tf-testAccMNSTopicConfig-%d"
		maximum_message_size=12357
		logging_enabled=true
	}
	`, rand)
}
