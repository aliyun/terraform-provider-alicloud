package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dxh031/ali_mns"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceAlicloudMNSTopicSubscription_basic(t *testing.T) {

	var attr ali_mns.TopicAttribute

	var subscriptionAttr ali_mns.SubscriptionAttribute

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// 配置 资源销毁结果检查函数
		CheckDestroy: testAccCheckMNSTopicSubscriptionDestroy,

		Steps: []resource.TestStep{
			{
				// 配置 配置内容
				Config: testAccMNSTopicSubscriptionConfig,
				// 配置 验证函数
				Check: resource.ComposeTestCheckFunc(
					// 验证资源topic
					testAccMNSTopicExist("alicloud_mns_topic.topic", &attr),
					//验证subscription
					testAccMNSTopicSubscriptionExist("alicloud_mns_topic_subscription.subscription", &subscriptionAttr),
					// 验证资源属性，能匹配到，肯定就是创建成功了
					resource.TestCheckResourceAttr("alicloud_mns_topic_subscription.subscription", "name", "tf-testAccMNSTopicSubscriptionConfig"),
				),
			},
			{
				// 配置 配置内容
				Config: testAccMNSTopicSubscriptionConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccMNSTopicExist("alicloud_mns_topic.topic", &attr),
					//验证subscription
					testAccMNSTopicSubscriptionExist("alicloud_mns_topic_subscription.subscription", &subscriptionAttr),
					// 验证修改后的属性值，如果能匹配到，肯定就是修改成功了
					resource.TestCheckResourceAttr("alicloud_mns_topic_subscription.subscription", "name", "tf-testAccMNSTopicSubscriptionConfig"),
					resource.TestCheckResourceAttr("alicloud_mns_topic_subscription.subscription", "notify_strategy", "EXPONENTIAL_DECAY_RETRY"),
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
		mnsClient, err := client.Mnsconn()
		if err != nil {
			return fmt.Errorf(" creating alicoudMNSTopic  error: %#v", err)
		}
		arr := strings.Split(rs.Primary.ID, "#")

		subscriptionManager := ali_mns.NewMNSTopic(arr[0], *mnsClient)
		instance, err := subscriptionManager.GetSubscriptionAttributes(arr[1])

		if err != nil {
			return err
		}
		if instance.SubscriptionName != arr[1] {
			return fmt.Errorf("mns subscription:%snot found", n)
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
		mnsClient, err := client.Mnsconn()
		if err != nil {
			return fmt.Errorf(" creating MNS Topic  error: %#v", err)
		}
		arr := strings.Split(rs.Primary.ID, "#")

		subscriptionManager := ali_mns.NewMNSTopic(arr[0], *mnsClient)

		if _, err := subscriptionManager.GetSubscriptionAttributes(arr[1]); err != nil {
			if strings.Contains(err.Error(), "SubscriptionNotExist") {
				continue
			}
			return err
		}

		return fmt.Errorf("MNS topic subscription %s still exist", arr[1])
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
