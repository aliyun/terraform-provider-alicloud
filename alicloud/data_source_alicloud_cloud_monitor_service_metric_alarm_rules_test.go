package alicloud

import (
	"fmt"
	"strings"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCloudMonitorServiceMetricAlarmRulesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_alarm.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_alarm.default.id}_fake"]`,
		}),
	}
	metricNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_cms_alarm.default.id}"]`,
			"metric_name": `"${alicloud_cms_alarm.default.metric}"`,
		}),
		fakeConfig: testAccCheckAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceName(rand, map[string]string{
			"metric_name": `"${alicloud_cms_alarm.default.metric}_fake"`,
		}),
	}
	namespaceConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_cms_alarm.default.id}"]`,
			"namespace": `"${alicloud_cms_alarm.default.project}"`,
		}),
		fakeConfig: testAccCheckAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceName(rand, map[string]string{
			"namespace": `"${alicloud_cms_alarm.default.project}_fake"`,
		}),
	}
	ruleNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceName(rand, map[string]string{
			"rule_name": `"${alicloud_cms_alarm.default.name}"`,
		}),
		fakeConfig: testAccCheckAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceName(rand, map[string]string{
			"rule_name": `"${alicloud_cms_alarm.default.name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cms_alarm.default.id}"]`,
			"status": `"true"`,
		}),
		fakeConfig: testAccCheckAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceName(rand, map[string]string{
			"status": `"false"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_cms_alarm.default.id}"]`,
			"metric_name": `"${alicloud_cms_alarm.default.metric}"`,
			"namespace":   `"${alicloud_cms_alarm.default.project}"`,
			"rule_name":   `"${alicloud_cms_alarm.default.name}"`,
			"status":      `"true"`,
		}),
		fakeConfig: testAccCheckAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_cms_alarm.default.id}_fake"]`,
			"metric_name": `"${alicloud_cms_alarm.default.metric}_fake"`,
			"namespace":   `"${alicloud_cms_alarm.default.project}_fake"`,
			"rule_name":   `"${alicloud_cms_alarm.default.name}_fake"`,
			"status":      `"false"`,
		}),
	}
	var existAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"rules.#":                          "1",
			"rules.0.id":                       CHECKSET,
			"rules.0.contact_groups":           CHECKSET,
			"rules.0.effective_interval":       CHECKSET,
			"rules.0.email_subject":            CHECKSET,
			"rules.0.metric_name":              CHECKSET,
			"rules.0.namespace":                CHECKSET,
			"rules.0.period":                   CHECKSET,
			"rules.0.resources":                CHECKSET,
			"rules.0.rule_name":                CHECKSET,
			"rules.0.silence_time":             CHECKSET,
			"rules.0.source_type":              CHECKSET,
			"rules.0.status":                   CHECKSET,
			"rules.0.webhook":                  CHECKSET,
			"rules.0.escalations.#":            "1",
			"rules.0.escalations.0.critical.#": "1",
			"rules.0.escalations.0.critical.0.comparison_operator": CHECKSET,
			"rules.0.escalations.0.critical.0.times":               CHECKSET,
			"rules.0.escalations.0.critical.0.statistics":          CHECKSET,
			"rules.0.escalations.0.critical.0.threshold":           CHECKSET,
			"rules.0.escalations.0.info.#":                         "1",
			"rules.0.escalations.0.info.0.comparison_operator":     CHECKSET,
			"rules.0.escalations.0.info.0.times":                   CHECKSET,
			"rules.0.escalations.0.info.0.statistics":              CHECKSET,
			"rules.0.escalations.0.info.0.threshold":               CHECKSET,
			"rules.0.escalations.0.warn.#":                         "1",
			"rules.0.escalations.0.warn.0.comparison_operator":     CHECKSET,
			"rules.0.escalations.0.warn.0.times":                   CHECKSET,
			"rules.0.escalations.0.warn.0.statistics":              CHECKSET,
			"rules.0.escalations.0.warn.0.threshold":               CHECKSET,
			"rules.0.labels.#":                                     "1",
			"rules.0.labels.0.value":                               CHECKSET,
			"rules.0.labels.0.key":                                 CHECKSET,
		}
	}
	var fakeAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"rules.#": "0",
		}
	}
	var alicloudCloudMonitorServiceMetricAlarmRulesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cloud_monitor_service_metric_alarm_rules.default",
		existMapFunc: existAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudCloudMonitorServiceMetricAlarmRulesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, metricNameConf, namespaceConf, ruleNameConf, statusConf, allConf)
}

func testAccCheckAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testacc-car-%d"
	}

	data "alicloud_instances" "default" {
  		status = "Running"
	}

	resource "alicloud_cms_alarm_contact_group" "default" {
  		alarm_contact_group_name = var.name
	}

	resource "alicloud_cms_alarm" "default" {
  		name               = var.name
  		project            = "acs_ecs_dashboard"
  		metric             = "disk_writebytes"
  		period             = 900
  		silence_time       = 300
  		webhook            = "https://www.aliyun.com"
  		enabled            = true
  		contact_groups     = [alicloud_cms_alarm_contact_group.default.alarm_contact_group_name]
  		effective_interval = "06:00-20:00"
  		metric_dimensions  = <<EOF
  [
    {
      "instanceId": "${data.alicloud_instances.default.ids.0}",
      "device": "/dev/vda1"
    }
  ]
  EOF
  		escalations_critical {
    		statistics          = "Average"
    		comparison_operator = "<="
    		threshold           = 90
    		times               = 1
  		}
  		escalations_info {
			statistics          = "Minimum"
			comparison_operator = "!="
			threshold           = 20
   		 	times               = 3
  		}
  		escalations_warn {
    		statistics          = "Average"
    		comparison_operator = "=="
    		threshold           = 30
    		times               = 5
  		}
  		tags = {
    		Created = "TF"
  		}
	}

	data "alicloud_cloud_monitor_service_metric_alarm_rules" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
