package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMessageServiceSubscriptionsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMessageServiceSubscriptionsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_message_service_subscription.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMessageServiceSubscriptionsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_message_service_subscription.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMessageServiceSubscriptionsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_message_service_subscription.default.subscription_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMessageServiceSubscriptionsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_message_service_subscription.default.subscription_name}_fake"`,
		}),
	}
	subscriptionNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMessageServiceSubscriptionsDataSourceName(rand, map[string]string{
			"subscription_name": `"${alicloud_message_service_subscription.default.subscription_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMessageServiceSubscriptionsDataSourceName(rand, map[string]string{
			"subscription_name": `"${alicloud_message_service_subscription.default.subscription_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMessageServiceSubscriptionsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_message_service_subscription.default.id}"]`,
			"name_regex":        `"${alicloud_message_service_subscription.default.subscription_name}"`,
			"subscription_name": `"${alicloud_message_service_subscription.default.subscription_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMessageServiceSubscriptionsDataSourceName(rand, map[string]string{
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
			"subscriptions.0.id":                    CHECKSET,
			"subscriptions.0.topic_name":            CHECKSET,
			"subscriptions.0.subscription_name":     CHECKSET,
			"subscriptions.0.endpoint":              "http://www.test.com/test",
			"subscriptions.0.filter_tag":            "tf-test",
			"subscriptions.0.notify_content_format": "JSON",
			"subscriptions.0.notify_strategy":       "EXPONENTIAL_DECAY_RETRY",
			"subscriptions.0.topic_owner":           CHECKSET,
			// subscription_url is intentionally not asserted: neither
			// ListSubscriptionByTopic nor GetSubscriptionAttributes returns
			// SubscriptionURL, so the attribute is always empty. The schema field is
			// kept for backward compatibility.
			"subscriptions.0.last_modify_time":                                   CHECKSET,
			"subscriptions.0.create_time":                                        CHECKSET,
			"subscriptions.0.tenant_rate_limit_policy.#":                         "1",
			"subscriptions.0.tenant_rate_limit_policy.0.enabled":                 "true",
			"subscriptions.0.tenant_rate_limit_policy.0.max_receives_per_second": "50",
			"subscriptions.0.dlq_policy.#":                                       "1",
			"subscriptions.0.dlq_policy.0.enabled":                               "false",
		}
	}
	var fakeAlicloudMessageServiceSubscriptionsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           "0",
			"names.#":         "0",
			"subscriptions.#": "0",
		}
	}
	var alicloudMessageServiceSubscriptionsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_message_service_subscriptions.default",
		existMapFunc: existAlicloudMessageServiceSubscriptionsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudMessageServiceSubscriptionsDataSourceNameMapFunc,
	}
	alicloudMessageServiceSubscriptionsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, subscriptionNameConf, allConf)
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

  		tenant_rate_limit_policy {
  			enabled                 = true
  			max_receives_per_second = 50
  		}
	}

	data "alicloud_message_service_subscriptions" "default" {
		topic_name     = alicloud_message_service_subscription.default.topic_name
		enable_details = true
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
