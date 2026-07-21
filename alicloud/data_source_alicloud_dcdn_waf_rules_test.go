package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudDcdnWafRuleDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDcdnWafRuleSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_dcdn_waf_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDcdnWafRuleSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_dcdn_waf_rule.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDcdnWafRuleSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_dcdn_waf_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDcdnWafRuleSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_dcdn_waf_rule.default.id}_fake"]`,
		}),
	}

	DcdnWafRuleCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existDcdnWafRuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                   "1",
		"waf_rules.#":                             "1",
		"waf_rules.0.id":                          CHECKSET,
		"waf_rules.0.policy_id":                   CHECKSET,
		"waf_rules.0.rule_name":                   CHECKSET,
		"waf_rules.0.conditions.#":                "2",
		"waf_rules.0.status":                      "on",
		"waf_rules.0.cc_status":                   "on",
		"waf_rules.0.action":                      "monitor",
		"waf_rules.0.effect":                      "rule",
		"waf_rules.0.rate_limit.#":                "1",
		"waf_rules.0.rate_limit.0.target":         "IP",
		"waf_rules.0.rate_limit.0.interval":       "5",
		"waf_rules.0.rate_limit.0.threshold":      "5",
		"waf_rules.0.rate_limit.0.ttl":            "1800",
		"waf_rules.0.rate_limit.0.status.#":       "1",
		"waf_rules.0.rate_limit.0.status.0.code":  "200",
		"waf_rules.0.rate_limit.0.status.0.ratio": "60",
	}
}

var fakeDcdnWafRuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":       "0",
		"waf_rules.#": "0",
	}
}

var DcdnWafRuleCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_dcdn_waf_rules.default",
	existMapFunc: existDcdnWafRuleMapFunc,
	fakeMapFunc:  fakeDcdnWafRuleMapFunc,
}

func testAccCheckAlicloudDcdnWafRuleSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf_testAccDcdnWafRule%d"
}

resource "alicloud_dcdn_waf_policy" "default" {
	defense_scene = "custom_acl"
	policy_name = var.name
	policy_type = "custom"
	status = "on"
}

resource "alicloud_dcdn_waf_rule" "default" {
  policy_id = alicloud_dcdn_waf_policy.default.id
  rule_name = var.name
  conditions {
    key      = "URI"
    op_value = "ne"
    values   = "/login.php"
  }
  conditions {
    key      = "Header"
    sub_key  = "a"
    op_value = "eq"
    values   = "b"
  }
  status    = "on"
  cc_status = "on"
  action    = "monitor"
  effect    = "rule"
  rate_limit {
    target    = "IP"
    interval  = "5"
    threshold = "5"
    ttl       = "1800"
    status {
      code  = "200"
      ratio = "60"
    }
  }
}

data "alicloud_dcdn_waf_rules" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
