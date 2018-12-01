package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudMnsTopicSubscriptionDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudMNSTopicSubscriptionDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_mns_topic_subscriptions.subscriptions"),
					resource.TestCheckResourceAttr("data.alicloud_mns_topic_subscriptions.subscriptions", "subscriptions.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_mns_topic_subscriptions.subscriptions", "subscriptions.0.name", "tf-testAccMNSTopicSubscriptionConfig"),
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
				Config: testAccCheckAlicloudMNSTopicSubscriptionDataSourceEmpty,
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

const testAccCheckAlicloudMNSTopicSubscriptionDataSourceConfig = `
data "alicloud_mns_topic_subscriptions" "subscriptions" {
  topic_name="${alicloud_mns_topic.topic.name}"
  name_prefix = "${alicloud_mns_topic_subscription.subscription.name}"
}

resource "alicloud_mns_topic" "topic"{
	name="tf-testAccMNSTopicConfig1"
	maximum_message_size=12357
	logging_enabled=true
}

resource "alicloud_mns_topic_subscription" "subscription"{
	topic_name="${alicloud_mns_topic.topic.name}"
	name="tf-testAccMNSTopicSubscriptionConfig"
	endpoint="http://www.test.com/test"
	notify_strategy="EXPONENTIAL_DECAY_RETRY"
	notify_content_format="SIMPLIFIED"
}`

const testAccCheckAlicloudMNSTopicSubscriptionDataSourceEmpty = `
data "alicloud_mns_topic_subscriptions" "subscriptions" {
  topic_name="${alicloud_mns_topic.topic.name}"
  name_prefix = "tf-testacc-fake-name"
}

resource "alicloud_mns_topic" "topic"{
	name="tf-testAccMNSTopicConfig1"
	maximum_message_size=12357
	logging_enabled=true
}
`
