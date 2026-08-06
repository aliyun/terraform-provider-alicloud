package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCmsSubscriptionsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCmsSubscriptionsSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_subscription.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudCmsSubscriptionsSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_subscription.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCmsSubscriptionsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cms_subscription.default.subscription_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudCmsSubscriptionsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cms_subscription.default.subscription_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCmsSubscriptionsSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_cms_subscription.default.id}"]`,
			"name_regex":  `"${alicloud_cms_subscription.default.subscription_name}"`,
			"output_file": `"./tf-testacc-cms-subscriptions.txt"`,
		}),
		fakeConfig: testAccCheckAliCloudCmsSubscriptionsSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cms_subscription.default.id}"]`,
			"workspace": `"tf-testacc-nonexistent-workspace"`,
		}),
	}
	CmsSubscriptionsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

var existCmsSubscriptionsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                     "1",
		"names.#":                                   "1",
		"subscriptions.#":                           "1",
		"subscriptions.0.subscription_id":           CHECKSET,
		"subscriptions.0.subscription_name":         CHECKSET,
		"subscriptions.0.notify_strategy_id":        CHECKSET,
		"subscriptions.0.enable":                    CHECKSET,
		"subscriptions.0.create_time":               CHECKSET,
		"subscriptions.0.update_time":               CHECKSET,
		"subscriptions.0.filter_setting.#":          "1",
		"subscriptions.0.filter_setting.0.relation": CHECKSET,
		"subscriptions.0.pushing_setting.#":         CHECKSET,
		"subscriptions.0.agent_config.#":            CHECKSET,
	}
}

var fakeCmsSubscriptionsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":           "0",
		"names.#":         "0",
		"subscriptions.#": "0",
	}
}

var CmsSubscriptionsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cms_subscriptions.default",
	existMapFunc: existCmsSubscriptionsMapFunc,
	fakeMapFunc:  fakeCmsSubscriptionsMapFunc,
}

func testAccCheckAliCloudCmsSubscriptionsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCmsSubscription%d"
}

resource "alicloud_cms_notify_strategy" "default" {
  notify_strategy_name = var.name
  grouping_setting {
    grouping_keys = ["severity"]
    period_min    = 5
    silence_sec   = 300
    times         = 1
  }
  routes {
    channels {
      channel_type = "DING"
      receivers    = ["CONTACT"]
    }
    severities = ["CRITICAL"]
  }
}

resource "alicloud_cms_subscription" "default" {
  subscription_name  = var.name
  notify_strategy_id = alicloud_cms_notify_strategy.default.id
  filter_setting {
    conditions {
      field = "severity"
      op    = "EQ"
      value = "CRITICAL"
    }
    relation = "AND"
  }
}

data "alicloud_cms_subscriptions" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
