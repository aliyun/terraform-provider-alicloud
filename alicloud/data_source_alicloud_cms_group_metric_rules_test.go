package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCmsGroupMetricRulesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsGroupMetricRulesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_group_metric_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsGroupMetricRulesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_group_metric_rule.default.id}_fake"]`,
		}),
	}
	groupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsGroupMetricRulesDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_cms_group_metric_rule.default.id}"]`,
			"group_id": `"${alicloud_cms_group_metric_rule.default.group_id}"`,
		}),
	}
	groupMetricRuleNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsGroupMetricRulesDataSourceName(rand, map[string]string{
			"ids":                    `["${alicloud_cms_group_metric_rule.default.id}"]`,
			"group_metric_rule_name": `"${alicloud_cms_group_metric_rule.default.group_metric_rule_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsGroupMetricRulesDataSourceName(rand, map[string]string{
			"ids":                    `["${alicloud_cms_group_metric_rule.default.id}"]`,
			"group_metric_rule_name": `"${alicloud_cms_group_metric_rule.default.group_metric_rule_name}_fake"`,
		}),
	}
	metricNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsGroupMetricRulesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_cms_group_metric_rule.default.id}"]`,
			"metric_name": `"${alicloud_cms_group_metric_rule.default.metric_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsGroupMetricRulesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_cms_group_metric_rule.default.id}"]`,
			"metric_name": `"${alicloud_cms_group_metric_rule.default.metric_name}_fake"`,
		}),
	}
	namespaceConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsGroupMetricRulesDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_cms_group_metric_rule.default.id}"]`,
			"namespace": `"${alicloud_cms_group_metric_rule.default.namespace}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsGroupMetricRulesDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_cms_group_metric_rule.default.id}"]`,
			"namespace": `"${alicloud_cms_group_metric_rule.default.namespace}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsGroupMetricRulesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cms_group_metric_rule.default.id}"]`,
			"status": `"${alicloud_cms_group_metric_rule.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsGroupMetricRulesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cms_group_metric_rule.default.id}"]`,
			"status": `"ALARM"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsGroupMetricRulesDataSourceName(rand, map[string]string{
			"group_id":               `"${alicloud_cms_group_metric_rule.default.group_id}"`,
			"group_metric_rule_name": `"${alicloud_cms_group_metric_rule.default.group_metric_rule_name}"`,
			"ids":                    `["${alicloud_cms_group_metric_rule.default.id}"]`,
			"metric_name":            `"${alicloud_cms_group_metric_rule.default.metric_name}"`,
			"namespace":              `"${alicloud_cms_group_metric_rule.default.namespace}"`,
			"status":                 `"${alicloud_cms_group_metric_rule.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsGroupMetricRulesDataSourceName(rand, map[string]string{
			"group_metric_rule_name": `"${alicloud_cms_group_metric_rule.default.group_metric_rule_name}_fake"`,
			"ids":                    `["${alicloud_cms_group_metric_rule.default.id}_fake"]`,
			"metric_name":            `"${alicloud_cms_group_metric_rule.default.metric_name}_fake"`,
			"namespace":              `"${alicloud_cms_group_metric_rule.default.namespace}_fake"`,
			"status":                 `"ALARM"`,
		}),
	}
	var existAlicloudCmsGroupMetricRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"rules.#":                        "1",
			"rules.0.contact_groups":         CHECKSET,
			"rules.0.email_subject":          "tf-testacc-rule-name-warning",
			"rules.0.group_id":               CHECKSET,
			"rules.0.group_metric_rule_name": CHECKSET,
			"rules.0.metric_name":            "cpu_total",
			"rules.0.namespace":              "acs_ecs_dashboard",
			"rules.0.period":                 "60",
			"rules.0.webhook":                "http://www.aliyun.com",
			"rules.0.status":                 "OK",
			"rules.0.rule_id":                "4a9a8978-a9cc-55ca-aa7c-530ccd91ae57",
			"rules.0.escalations.0.warn.0.comparison_operator": "GreaterThanOrEqualToThreshold",
			"rules.0.escalations.0.warn.0.statistics":          "Average",
			"rules.0.escalations.0.warn.0.threshold":           "90",
			"rules.0.escalations.0.warn.0.times":               "3",
		}
	}
	var fakeAlicloudCmsGroupMetricRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"rules.#": "0",
		}
	}
	var alicloudCmsGroupMetricRulesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_group_metric_rules.default",
		existMapFunc: existAlicloudCmsGroupMetricRulesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCmsGroupMetricRulesDataSourceNameMapFunc,
	}
	alicloudCmsGroupMetricRulesCheckInfo.dataSourceTestCheck(t, rand, idsConf, groupIdConf, groupMetricRuleNameConf, metricNameConf, namespaceConf, statusConf, allConf)
}
func testAccCheckAlicloudCmsGroupMetricRulesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

		variable "name" {
			default = "tf-testAccCmsAlarmContactGroupBisic-%d"
		}

		resource "alicloud_cms_alarm_contact_group" "this" {
		  alarm_contact_group_name = "${var.name}" 
		  describe = "tf-testacc" 
		  contacts = ["zhangsan","lisi","lll"] 
		}
		resource "alicloud_cms_monitor_group" "default" {
			monitor_group_name = var.name
			contact_groups = [alicloud_cms_alarm_contact_group.this.id]
		}
		resource "alicloud_cms_group_metric_rule" "default" {
		  contact_groups         = alicloud_cms_alarm_contact_group.this.alarm_contact_group_name
		  group_id               = "${alicloud_cms_monitor_group.default.id}" 
		  rule_id                = "4a9a8978-a9cc-55ca-aa7c-530ccd91ae57"
		
		  category               = "ecs"
		  namespace              = "acs_ecs_dashboard"
		  metric_name            = "cpu_total"
		  period                 = "60"
		
		  group_metric_rule_name = "${var.name}" 
		  email_subject          = "tf-testacc-rule-name-warning"
		  interval               = "3600"
		  silence_time           = 86400
		  effective_interval     = ""
		  no_effective_interval  = "00:00-05:30"
		  webhook                = "http://www.aliyun.com"
		  escalations {
			warn {
			  comparison_operator = "GreaterThanOrEqualToThreshold"
			  statistics          = "Average"
			  threshold           = "90"
			  times               = 3
			}
		  }
		  depends_on = [alicloud_cms_alarm_contact_group.this]
		}

		data "alicloud_cms_group_metric_rules" "default" {
			  %s
		}
`, rand, strings.Join(pairs, " \n "))
	return config
}
