// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudWafv3DefenseRuleDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafv3DefenseRuleSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_wafv3_defense_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudWafv3DefenseRuleSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_wafv3_defense_rule.default.id}_fake"]`,
		}),
	}

	DefenseTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafv3DefenseRuleSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_wafv3_defense_rule.default.id}"]`,
			"defense_type": `"template"`,
		}),
		fakeConfig: testAccCheckAlicloudWafv3DefenseRuleSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_wafv3_defense_rule.default.id}_fake"]`,
			"defense_type": `"resource"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafv3DefenseRuleSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_wafv3_defense_rule.default.id}"]`,
			"defense_type": `"template"`,
		}),
		fakeConfig: testAccCheckAlicloudWafv3DefenseRuleSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_wafv3_defense_rule.default.id}_fake"]`,
			"defense_type": `"resource"`,
		}),
	}

	Wafv3DefenseRuleCheckInfo.dataSourceTestCheck(t, rand, idsConf, DefenseTypeConf, allConf)
}

var existWafv3DefenseRuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"rules.#":                "1",
		"rules.0.defense_origin": CHECKSET,
		"rules.0.config.#":       CHECKSET,
		"rules.0.rule_id":        CHECKSET,
		"rules.0.gmt_modified":   CHECKSET,
		"rules.0.defense_type":   CHECKSET,
		"rules.0.defense_scene":  CHECKSET,
		"rules.0.rule_status":    CHECKSET,
		"rules.0.template_id":    CHECKSET,
		"rules.0.rule_name":      CHECKSET,
	}
}

var fakeWafv3DefenseRuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"rules.#": "0",
	}
}

var Wafv3DefenseRuleCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_wafv3_defense_rules.default",
	existMapFunc: existWafv3DefenseRuleMapFunc,
	fakeMapFunc:  fakeWafv3DefenseRuleMapFunc,
}

func testAccCheckAlicloudWafv3DefenseRuleSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccWafv3DefenseRule%d"
}
variable "region_id" {
  default = "cn-hangzhou"
}

data "alicloud_wafv3_instances" "default" {
}

resource "alicloud_wafv3_defense_template" "defaultfIoHt5" {
  status                = "1"
  description           = "testCreate"
  instance_id           = data.alicloud_wafv3_instances.default.ids.0
  defense_template_name = "tf-tpl-${var.name}"
  template_origin       = "custom"
  defense_scene         = "custom_acl"
  template_type         = "user_custom"
}

resource "alicloud_wafv3_address_book" "default9dtEmt" {
  description       = "test"
  instance_id       = data.alicloud_wafv3_instances.default.ids.0
  address_book_name = "tf-ab-${var.name}"
  address_list      = ["100.100.100.100/32", "101.101.101.101/32", "102.102.102.102/32"]
  address_book_type = "ip"
}

resource "alicloud_wafv3_defense_rule" "default" {
  defense_origin = "custom"
  instance_id    = data.alicloud_wafv3_instances.default.ids.0
  config {
    rule_action = "block"
    conditions {
      op_value = "contain"
      values   = "abc"
      key      = "URL"
    }
    conditions {
      op_value = "contain"
      values   = "abc"
      key      = "URLPath"
    }
    conditions {
      op_value = "contain"
      values   = "1.1.1.2"
      key      = "IP"
    }
    conditions {
      key      = "IP"
      op_value = "in-list"
      values   = alicloud_wafv3_address_book.default9dtEmt.address_book_id
    }
    cc_status = "0"
    cc_effect = "service"
    rate_limit {
      target    = "remote_addr"
      interval  = "16"
      threshold = "204"
      ttl       = "68"
      status {
        code  = "414"
        count = "333"
      }
      sub_key = "testky1"
    }
    gray_status = "1"
    gray_config {
      gray_target = "remote_addr"
      gray_rate   = "80"
    }
    time_config {
      time_scope = "period"
      time_zone  = "8"
      time_periods {
        start = "1760174804000"
        end   = "1760175804000"
      }
      time_periods {
        start = "1760171804000"
        end   = "1760172804000"
      }
      time_periods {
        start = "1760176804000"
        end   = "1760177804000"
      }
      time_periods {
        start = "1760178804000"
        end   = "1760179804000"
      }
      time_periods {
        start = "1760170804000"
        end   = "1760171804000"
      }
    }
  }
  defense_scene = "custom_acl"
  rule_status   = "1"
  defense_type  = "template"
  template_id   = alicloud_wafv3_defense_template.defaultfIoHt5.defense_template_id
  rule_name     = "custom_acl-create"
}

data "alicloud_wafv3_defense_rules" "default" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  %s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
