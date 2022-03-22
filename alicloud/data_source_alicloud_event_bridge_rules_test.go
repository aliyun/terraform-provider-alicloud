package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEventBridgeRulesDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.EventBridgeSupportRegions)
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeRulesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_event_bridge_rule.default.rule_name}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeRulesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_event_bridge_rule.default.rule_name}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeRulesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_event_bridge_rule.default.rule_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeRulesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_event_bridge_rule.default.rule_name}_fake"`,
		}),
	}
	namePrefixConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeRulesDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_event_bridge_rule.default.rule_name}"]`,
			"rule_name_prefix": `"${alicloud_event_bridge_rule.default.rule_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeRulesDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_event_bridge_rule.default.rule_name}"]`,
			"rule_name_prefix": `"tf-testAcc_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeRulesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_event_bridge_rule.default.rule_name}"]`,
			"status": `"ENABLE"`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeRulesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_event_bridge_rule.default.rule_name}"]`,
			"status": `"DISABLE"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeRulesDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_event_bridge_rule.default.rule_name}"]`,
			"name_regex":       `"${alicloud_event_bridge_rule.default.rule_name}"`,
			"rule_name_prefix": `"tf-testAcc"`,
			"status":           `"ENABLE"`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeRulesDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_event_bridge_rule.default.rule_name}_fake"]`,
			"name_regex":       `"${alicloud_event_bridge_rule.default.rule_name}_fake"`,
			"rule_name_prefix": `"tf-testAcc_fake"`,
			"status":           `"DISABLE"`,
		}),
	}
	var existAlicloudEventBridgeRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "1",
			"names.#":                "1",
			"rules.#":                "1",
			"rules.0.description":    fmt.Sprintf("tf-testAccRules-%d", rand),
			"rules.0.event_bus_name": fmt.Sprintf("tf-testAccRules-%d", rand),
			"rules.0.filter_pattern": "{\"source\":[\"crmabc.newsletter\"],\"type\":[\"UserSignUp\", \"UserLogin\"]}",
			"rules.0.id":             fmt.Sprintf("tf-testAccRules-%d", rand),
			"rules.0.rule_name":      fmt.Sprintf("tf-testAccRules-%d", rand),
			"rules.0.status":         "ENABLE",
			"rules.0.targets.#":      "1",
		}
	}
	var fakeAlicloudEventBridgeRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"rules.#": "0",
		}
	}
	var alicloudEventBridgeRulesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_event_bridge_rules.default",
		existMapFunc: existAlicloudEventBridgeRulesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEventBridgeRulesDataSourceNameMapFunc,
	}
	alicloudEventBridgeRulesCheckInfo.dataSourceTestCheck(t, rand, idsConf, namePrefixConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudEventBridgeRulesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccRules-%d"
}

resource "alicloud_event_bridge_event_bus" "default" {
	event_bus_name = var.name
}

data "alicloud_account" "default" {}

locals {
  mns_endpoint = format("acs:mns:%s:%%s:queues/%%s", data.alicloud_account.default.id, alicloud_mns_queue.queue1.name)
}

resource "alicloud_mns_queue" "queue1" {
  name = var.name
}

resource "alicloud_event_bridge_rule" "default" {
  rule_name      = var.name
  event_bus_name = alicloud_event_bridge_event_bus.default.id
  description    = var.name
  filter_pattern = "{\"source\":[\"crmabc.newsletter\"],\"type\":[\"UserSignUp\", \"UserLogin\"]}"
  targets {
    target_id           = "tf-test"
    endpoint            = local.mns_endpoint
    type                = "acs.mns.queue"
    param_list {
      resource_key = "queue"
      form         = "CONSTANT"
      value        = "tf-testaccEbRule"
    }
    param_list {
      resource_key = "Body"
      form         = "ORIGINAL"
    }
  }
}

data "alicloud_event_bridge_rules" "default" {
	event_bus_name = alicloud_event_bridge_event_bus.default.event_bus_name
	%s
}
`, rand, os.Getenv("ALICLOUD_REGION"), strings.Join(pairs, " \n "))
	return config
}
