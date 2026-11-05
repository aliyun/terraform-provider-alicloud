package mnsopen_test

import (
	"fmt"
	"strings"
	"testing"

	tfacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudMessageServiceSubscriptionsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMessageServiceSubscriptionsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_message_service_subscription.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudMessageServiceSubscriptionsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_message_service_subscription.default.id}_fake"]`,
		}),
	}
	nameRegexConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMessageServiceSubscriptionsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_message_service_subscription.default.subscription_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudMessageServiceSubscriptionsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_message_service_subscription.default.subscription_name}_fake"`,
		}),
	}
	subscriptionNameConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMessageServiceSubscriptionsDataSourceName(rand, map[string]string{
			"subscription_name": `"${alicloud_message_service_subscription.default.subscription_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudMessageServiceSubscriptionsDataSourceName(rand, map[string]string{
			"subscription_name": `"${alicloud_message_service_subscription.default.subscription_name}_fake"`,
		}),
	}
	allConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMessageServiceSubscriptionsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_message_service_subscription.default.id}"]`,
			"name_regex":        `"${alicloud_message_service_subscription.default.subscription_name}"`,
			"subscription_name": `"${alicloud_message_service_subscription.default.subscription_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudMessageServiceSubscriptionsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_message_service_subscription.default.id}_fake"]`,
			"name_regex":        `"${alicloud_message_service_subscription.default.subscription_name}_fake"`,
			"subscription_name": `"${alicloud_message_service_subscription.default.subscription_name}_fake"`,
		}),
	}
	var existAlicloudMessageServiceSubscriptionsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 "1",
			"names.#":                               "1",
			"subscriptions.#":                       "1",
			"subscriptions.0.id":                    tfacctest.CHECKSET,
			"subscriptions.0.topic_name":            tfacctest.CHECKSET,
			"subscriptions.0.subscription_name":     tfacctest.CHECKSET,
			"subscriptions.0.endpoint":              "http://www.test.com/test",
			"subscriptions.0.filter_tag":            "tf-test",
			"subscriptions.0.notify_content_format": "JSON",
			"subscriptions.0.notify_strategy":       "EXPONENTIAL_DECAY_RETRY",
			"subscriptions.0.topic_owner":           tfacctest.CHECKSET,
			"subscriptions.0.subscription_url":      tfacctest.CHECKSET,
			"subscriptions.0.last_modify_time":      tfacctest.CHECKSET,
			"subscriptions.0.create_time":           tfacctest.CHECKSET,
		}
	}
	var fakeAlicloudMessageServiceSubscriptionsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           "0",
			"names.#":         "0",
			"subscriptions.#": "0",
		}
	}
	var alicloudMessageServiceSubscriptionsCheckInfo = tfacctest.DataSourceAttr{
		ResourceId:   "data.alicloud_message_service_subscriptions.default",
		ExistMapFunc: existAlicloudMessageServiceSubscriptionsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudMessageServiceSubscriptionsDataSourceNameMapFunc,
	}
	alicloudMessageServiceSubscriptionsCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, subscriptionNameConf, allConf)
}

func testAccCheckAlicloudMessageServiceSubscriptionsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  	default = "tf-testAccMNSSubscription-%d"
	}

	resource "alicloud_message_service_topic" "default" {
  		topic_name       = var.name
  		max_message_size = 12357
  		logging_enabled  = true
	}

	resource "alicloud_message_service_subscription" "default" {
  		topic_name            = alicloud_message_service_topic.default.topic_name
  		subscription_name     = var.name
  		endpoint              = "http://www.test.com/test"
  		push_type             = "http"
  		filter_tag            = "tf-test"
  		notify_content_format = "JSON"
  		notify_strategy       = "EXPONENTIAL_DECAY_RETRY"
	}
	
	data "alicloud_message_service_subscriptions" "default" {
		topic_name = alicloud_message_service_subscription.default.topic_name
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
