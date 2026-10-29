// Package mns_test.
package mns_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudMnsTopicSubscriptionDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_mns_topic_subscriptions.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccMNSTopicConfig-%d", rand),
		dataSourceMnsTopicSubscriptionConfigDependence)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"topic_name":  "${alicloud_mns_topic.default.name}",
			"name_prefix": "${alicloud_mns_topic_subscription.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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

	var mnsTopicSubscriptionCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existMnsTopicSubscriptionMapFunc,
		FakeMapFunc:  fakeMnsTopicSubscriptionMapFunc,
	}

	mnsTopicSubscriptionCheckInfo.DataSourceTestCheck(t, rand, allConf)
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
