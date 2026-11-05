package cloudmonitorservice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCmsEventRulesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsEventRulesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_event_rule.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCmsEventRulesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_event_rule.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsEventRulesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cms_event_rule.default.rule_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsEventRulesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cms_event_rule.default.rule_name}_fake"`,
		}),
	}
	namePrefixConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsEventRulesDataSourceName(rand, map[string]string{
			"name_prefix": `"${alicloud_cms_event_rule.default.rule_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsEventRulesDataSourceName(rand, map[string]string{
			"name_prefix": `"${alicloud_cms_event_rule.default.rule_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsEventRulesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cms_event_rule.default.id}"]`,
			"status": `"ENABLED"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsEventRulesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cms_event_rule.default.id}_fake"]`,
			"status": `"DISABLED"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsEventRulesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_cms_event_rule.default.id}"]`,
			"name_regex":  `"${alicloud_cms_event_rule.default.rule_name}"`,
			"name_prefix": `"tf-testAcc"`,
			"status":      `"ENABLED"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsEventRulesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_cms_event_rule.default.id}_fake"]`,
			"name_regex":  `"${alicloud_cms_event_rule.default.rule_name}_fake"`,
			"name_prefix": `"tf-testAcc_fake"`,
			"status":      `"DISABLED"`,
		}),
	}
	var existAlicloudCmsEventRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"names.#":                 "1",
			"rules.#":                 "1",
			"rules.0.id":              acctest.CHECKSET,
			"rules.0.event_rule_name": acctest.CHECKSET,
			"rules.0.description":     acctest.CHECKSET,
			"rules.0.event_type":      acctest.CHECKSET,
			"rules.0.group_id":        acctest.CHECKSET,
			"rules.0.silence_time":    "100",
			"rules.0.status":          "ENABLED",
			"rules.0.event_pattern.#": "1",
			"rules.0.event_pattern.0.event_type_list.#": "1",
			"rules.0.event_pattern.0.level_list.#":      "1",
			"rules.0.event_pattern.0.name_list.#":       "1",
			"rules.0.event_pattern.0.product":           acctest.CHECKSET,
			"rules.0.event_pattern.0.sql_filter":        acctest.CHECKSET,
			"rules.0.event_pattern.0.keyword_filter.#":  acctest.NOSET,
		}
	}
	var fakeAlicloudCmsEventRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"rules.#": "0",
		}
	}
	var alicloudCmsEventRulesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cms_event_rules.default",
		ExistMapFunc: existAlicloudCmsEventRulesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCmsEventRulesDataSourceNameMapFunc,
	}
	alicloudCmsEventRulesCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, namePrefixConf, statusConf, allConf)
}

func testAccCheckAlicloudCmsEventRulesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccCmsEventRuleBisic-%d"
	}
	
	resource "alicloud_cms_monitor_group" "default" {
		monitor_group_name = var.name
	}
	
	resource "alicloud_cms_event_rule" "default" {
		rule_name    = var.name
		group_id     = "${alicloud_cms_monitor_group.default.id}"
		description  = "tf-testAcc"
		status       = "ENABLED"
		silence_time = 100
		event_pattern {
		product         = "ecs"
		event_type_list = ["StatusNotification"]
		level_list      = ["CRITICAL"]
		name_list       = ["test"]
		sql_filter      = "test"
		}
	}
	
	data "alicloud_cms_event_rules" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
