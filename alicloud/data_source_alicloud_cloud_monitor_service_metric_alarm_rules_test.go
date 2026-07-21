// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudMonitorServiceMetricAlarmRuleDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudMonitorServiceMetricAlarmRuleSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cloud_monitor_service_metric_alarm_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudMonitorServiceMetricAlarmRuleSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cloud_monitor_service_metric_alarm_rule.default.id}_fake"]`,
		}),
	}

	NamespaceConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudMonitorServiceMetricAlarmRuleSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cloud_monitor_service_metric_alarm_rule.default.id}"]`,
			"namespace": `"acs_drds"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudMonitorServiceMetricAlarmRuleSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cloud_monitor_service_metric_alarm_rule.default.id}_fake"]`,
			"namespace": `"acs_drds_fake"`,
		}),
	}
	MetricAlarmRuleIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudMonitorServiceMetricAlarmRuleSourceConfig(rand, map[string]string{
			"ids":                  `["${alicloud_cloud_monitor_service_metric_alarm_rule.default.id}"]`,
			"metric_alarm_rule_id": `"SystemDefault_acs_drds_IOPSUsageOfDo"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudMonitorServiceMetricAlarmRuleSourceConfig(rand, map[string]string{
			"ids":                  `["${alicloud_cloud_monitor_service_metric_alarm_rule.default.id}_fake"]`,
			"metric_alarm_rule_id": `"SystemDefault_acs_drds_IOPSUsageOfDo_fake"`,
		}),
	}
	StatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudMonitorServiceMetricAlarmRuleSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cloud_monitor_service_metric_alarm_rule.default.id}"]`,
			"status": `true`,
		}),
		fakeConfig: testAccCheckAlicloudCloudMonitorServiceMetricAlarmRuleSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cloud_monitor_service_metric_alarm_rule.default.id}_fake"]`,
			"status": `false`,
		}),
	}
	MetricNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudMonitorServiceMetricAlarmRuleSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_cloud_monitor_service_metric_alarm_rule.default.id}"]`,
			"metric_name": `"IOPSUsageOfDN"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudMonitorServiceMetricAlarmRuleSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_cloud_monitor_service_metric_alarm_rule.default.id}_fake"]`,
			"metric_name": `"IOPSUsageOfDN_fake"`,
		}),
	}
	RuleNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudMonitorServiceMetricAlarmRuleSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cloud_monitor_service_metric_alarm_rule.default.id}"]`,
			"rule_name": `"SystemDefault_acs_drds_IOPSUsageOfDN"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudMonitorServiceMetricAlarmRuleSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cloud_monitor_service_metric_alarm_rule.default.id}_fake"]`,
			"rule_name": `"SystemDefault_acs_drds_IOPSUsageOfDN_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudMonitorServiceMetricAlarmRuleSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cloud_monitor_service_metric_alarm_rule.default.id}"]`,
			"namespace": `"acs_drds"`,

			"metric_alarm_rule_id": `"SystemDefault_acs_drds_IOPSUsageOfDo"`,

			"status": `true`,

			"metric_name": `"IOPSUsageOfDN"`,

			"rule_name": `"SystemDefault_acs_drds_IOPSUsageOfDN"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudMonitorServiceMetricAlarmRuleSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cloud_monitor_service_metric_alarm_rule.default.id}_fake"]`,
			"namespace": `"acs_drds_fake"`,

			"metric_alarm_rule_id": `"SystemDefault_acs_drds_IOPSUsageOfDo_fake"`,

			"status": `false`,

			"metric_name": `"IOPSUsageOfDN_fake"`,

			"rule_name": `"SystemDefault_acs_drds_IOPSUsageOfDN_fake"`,
		}),
	}

	CloudMonitorServiceMetricAlarmRuleCheckInfo.dataSourceTestCheck(t, rand, idsConf, NamespaceConf, MetricAlarmRuleIdConf, StatusConf, MetricNameConf, RuleNameConf, allConf)
}

var existCloudMonitorServiceMetricAlarmRuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"rules.#":                        "1",
		"rules.0.source_type":            CHECKSET,
		"rules.0.prometheus.#":           CHECKSET,
		"rules.0.metric_name":            CHECKSET,
		"rules.0.email_subject":          CHECKSET,
		"rules.0.composite_expression.#": CHECKSET,
		"rules.0.rule_name":              CHECKSET,
		"rules.0.status":                 CHECKSET,
		"rules.0.silence_time":           CHECKSET,
		"rules.0.contact_groups":         CHECKSET,
		"rules.0.send_ok":                "true",
		"rules.0.metric_alarm_rule_id":   CHECKSET,
		"rules.0.period":                 CHECKSET,
		"rules.0.labels.#":               CHECKSET,
		"rules.0.effective_interval":     CHECKSET,
		"rules.0.no_data_policy":         "KEEP_LAST_STATE",
		"rules.0.namespace":              CHECKSET,
		"rules.0.escalations.#":          CHECKSET,
		"rules.0.resources":              CHECKSET,
	}
}

var fakeCloudMonitorServiceMetricAlarmRuleMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"rules.#": "0",
	}
}

var CloudMonitorServiceMetricAlarmRuleCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cloud_monitor_service_metric_alarm_rules.default",
	existMapFunc: existCloudMonitorServiceMetricAlarmRuleMapFunc,
	fakeMapFunc:  fakeCloudMonitorServiceMetricAlarmRuleMapFunc,
}

func testAccCheckAlicloudCloudMonitorServiceMetricAlarmRuleSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCloudMonitorServiceMetricAlarmRule%d"
}


resource "alicloud_cloud_monitor_service_metric_alarm_rule" "default" {
  status               = true
  send_ok              = true
  contact_groups       = "云账号报警联系人"
  silence_time         = "86400"
  metric_alarm_rule_id = "SystemDefault_acs_drds_IOPSUsageOfDo"
  period               = "60"
  effective_interval   = "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7"
  no_data_policy       = "KEEP_LAST_STATE"
  namespace            = "acs_drds"
  metric_name          = "IOPSUsageOfDN"
  escalations {
    critical {
      comparison_operator = "GreaterThanThreshold"
      times               = "5"
      statistics          = "Average"
      threshold           = "80"
    }
  }
  resources = "[{\"resource\":\"_ALL\"}]"
  rule_name = "SystemDefault_acs_drds_IOPSUsageOfDN"
}

data "alicloud_cloud_monitor_service_metric_alarm_rules" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
