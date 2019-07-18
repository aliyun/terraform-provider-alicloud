package alicloud

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudMnsTopicSubscriptionDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_mns_topic_subscriptions.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccMNSTopicConfig-%d", rand),
		dataSourceMnsTopicSubscriptionConfigDependence)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"topic_name":  "${alicloud_mns_topic.default.name}",
			"name_prefix": "${alicloud_mns_topic_subscription.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"topic_name":  "${alicloud_mns_topic.default.name}",
			"name_prefix": "${alicloud_mns_topic_subscription.default.name}-fake",
		}),
	}

	var existMnsTopicSubscriptionMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                               "1",
			"names.0":                               fmt.Sprintf("tf-testAccMNSTopicConfig-%d", rand),
			"subscriptions.#":                       "1",
			"subscriptions.0.name":                  fmt.Sprintf("tf-testAccMNSTopicConfig-%d", rand),
			"subscriptions.0.id":                    fmt.Sprintf("tf-testAccMNSTopicConfig-%d", rand),
			"subscriptions.0.topic_name":            fmt.Sprintf("tf-testAccMNSTopicConfig-%d", rand),
			"subscriptions.0.endpoint":              "http://www.test.com/test",
			"subscriptions.0.notify_strategy":       "EXPONENTIAL_DECAY_RETRY",
			"subscriptions.0.notify_content_format": "SIMPLIFIED",
		}
	}

	var fakeMnsTopicSubscriptionMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":         "0",
			"subscriptions.#": "0",
		}
	}

	var mnsTopicSubscriptionCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMnsTopicSubscriptionMapFunc,
		fakeMapFunc:  fakeMnsTopicSubscriptionMapFunc,
	}

	mnsTopicSubscriptionCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

func dataSourceMnsTopicSubscriptionConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	    default = "%s"
	}

	resource "alicloud_mns_topic" "default"{
		name="${var.name}"
		maximum_message_size=12357
		logging_enabled=true
	}

	resource "alicloud_mns_topic_subscription" "default"{
		topic_name="${alicloud_mns_topic.default.name}"
		name="${var.name}"
		endpoint="http://www.test.com/test"
		notify_strategy="EXPONENTIAL_DECAY_RETRY"
		notify_content_format="SIMPLIFIED"
	}`, name)
}
