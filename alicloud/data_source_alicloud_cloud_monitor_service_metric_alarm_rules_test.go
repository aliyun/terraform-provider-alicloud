package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCloudMonitorServiceMetricAlarmRulesDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	resourceId := "data.alicloud_cloud_monitor_service_metric_alarm_rules.default"
	name := fmt.Sprintf("tf-testscc-cmsalarm%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceCloudMonitorServiceMetricAlarmRulesConfig0)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_cms_alarm.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_cms_alarm.default.id}_fake"},
		}),
	}
	dimensionsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"dimensions": `{\"instanceId\":\"` + "${data.alicloud_instances.default.ids.0}" + `\"}`,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"dimensions": `{\"instanceId\":\"i-1688\"}`,
		}),
	}
	metricNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_cms_alarm.default.id}"},
			"metric_name": "${alicloud_cms_alarm.default.metric}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"metric_name": "${alicloud_cms_alarm.default.metric}_fake",
		}),
	}
	namespaceConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":       []string{"${alicloud_cms_alarm.default.id}"},
			"namespace": "${alicloud_cms_alarm.default.project}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"namespace": "${alicloud_cms_alarm.default.project}_fake",
		}),
	}
	ruleNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"rule_name": "${alicloud_cms_alarm.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"rule_name": "${alicloud_cms_alarm.default.name}_fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_cms_alarm.default.id}"},
			"status": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_cms_alarm.default.id}_fake"},
			"status": "false",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_cms_alarm.default.id}"},
			"dimensions":  `{\"instanceId\":\"` + "${data.alicloud_instances.default.ids.0}" + `\"}`,
			"metric_name": "${alicloud_cms_alarm.default.metric}",
			"namespace":   "${alicloud_cms_alarm.default.project}",
			"rule_name":   "${alicloud_cms_alarm.default.name}",
			"status":      "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_cms_alarm.default.id}_fake"},
			"dimensions":  `{\"instanceId\":\"i-1688\"}`,
			"metric_name": "${alicloud_cms_alarm.default.metric}_fake",
			"namespace":   "${alicloud_cms_alarm.default.project}_fake",
			"rule_name":   "${alicloud_cms_alarm.default.name}_fake",
			"status":      "false",
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
	alicloudCloudMonitorServiceMetricAlarmRulesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, dimensionsConf, metricNameConf, namespaceConf, ruleNameConf, statusConf, allConf)
}

func dataSourceCloudMonitorServiceMetricAlarmRulesConfig0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
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
`, name)
}

func TestAccAliCloudCloudMonitorServiceMetricAlarmRulesDataSource_basic1(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	resourceId := "data.alicloud_cloud_monitor_service_metric_alarm_rules.default"
	name := fmt.Sprintf("tf-testscc-cmsalarm%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceCloudMonitorServiceMetricAlarmRulesConfig1)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_cms_alarm.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_cms_alarm.default.id}_fake"},
		}),
	}
	dimensionsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"dimensions": `{\"instanceId\":\"` + "${data.alicloud_instances.default.ids.0}" + `\"}`,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"dimensions": `{\"instanceId\":\"i-1688\"}`,
		}),
	}
	metricNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_cms_alarm.default.id}"},
			"metric_name": "${alicloud_cms_alarm.default.metric}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"metric_name": "${alicloud_cms_alarm.default.metric}_fake",
		}),
	}
	namespaceConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":       []string{"${alicloud_cms_alarm.default.id}"},
			"namespace": "${alicloud_cms_alarm.default.project}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"namespace": "${alicloud_cms_alarm.default.project}_fake",
		}),
	}
	ruleNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"rule_name": "${alicloud_cms_alarm.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"rule_name": "${alicloud_cms_alarm.default.name}_fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_cms_alarm.default.id}"},
			"status": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_cms_alarm.default.id}_fake"},
			"status": "false",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_cms_alarm.default.id}"},
			"dimensions":  `{\"instanceId\":\"` + "${data.alicloud_instances.default.ids.0}" + `\"}`,
			"metric_name": "${alicloud_cms_alarm.default.metric}",
			"namespace":   "${alicloud_cms_alarm.default.project}",
			"rule_name":   "${alicloud_cms_alarm.default.name}",
			"status":      "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_cms_alarm.default.id}_fake"},
			"dimensions":  `{\"instanceId\":\"i-1688\"}`,
			"metric_name": "${alicloud_cms_alarm.default.metric}_fake",
			"namespace":   "${alicloud_cms_alarm.default.project}_fake",
			"rule_name":   "${alicloud_cms_alarm.default.name}_fake",
			"status":      "false",
		}),
	}
	var existAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"rules.#":                              "1",
			"rules.0.id":                           CHECKSET,
			"rules.0.contact_groups":               CHECKSET,
			"rules.0.effective_interval":           CHECKSET,
			"rules.0.email_subject":                CHECKSET,
			"rules.0.metric_name":                  CHECKSET,
			"rules.0.namespace":                    CHECKSET,
			"rules.0.period":                       CHECKSET,
			"rules.0.resources":                    CHECKSET,
			"rules.0.rule_name":                    CHECKSET,
			"rules.0.silence_time":                 CHECKSET,
			"rules.0.source_type":                  CHECKSET,
			"rules.0.status":                       CHECKSET,
			"rules.0.webhook":                      CHECKSET,
			"rules.0.composite_expression.#":       "1",
			"rules.0.composite_expression.0.times": CHECKSET,
			"rules.0.composite_expression.0.expression_raw":                        CHECKSET,
			"rules.0.composite_expression.0.expression_list_join":                  CHECKSET,
			"rules.0.composite_expression.0.level":                                 CHECKSET,
			"rules.0.composite_expression.0.expression_list.#":                     "1",
			"rules.0.composite_expression.0.expression_list.0.metric_name":         CHECKSET,
			"rules.0.composite_expression.0.expression_list.0.comparison_operator": CHECKSET,
			"rules.0.composite_expression.0.expression_list.0.period":              CHECKSET,
			"rules.0.composite_expression.0.expression_list.0.statistics":          CHECKSET,
			"rules.0.composite_expression.0.expression_list.0.threshold":           CHECKSET,
			"rules.0.labels.#":       "1",
			"rules.0.labels.0.value": CHECKSET,
			"rules.0.labels.0.key":   CHECKSET,
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
	alicloudCloudMonitorServiceMetricAlarmRulesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, dimensionsConf, metricNameConf, namespaceConf, ruleNameConf, statusConf, allConf)
}

func dataSourceCloudMonitorServiceMetricAlarmRulesConfig1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
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
  		composite_expression {
    		level                = "CRITICAL"
    		times                = 1
    		expression_raw       = "$Average > 60"
    		expression_list_join = "&&"
    		expression_list {
      			metric_name         = "cpu_total"
      			comparison_operator = "=="
      			statistics          = "Maximum"
      			threshold           = 50
      			period              = 2
    		}
  		}
  		tags = {
    		Created = "TF"
  		}
	}
`, name)
}

func TestAccAliCloudCloudMonitorServiceMetricAlarmRulesDataSource_basic2(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	resourceId := "data.alicloud_cloud_monitor_service_metric_alarm_rules.default"
	name := fmt.Sprintf("tf-testscc-cmsalarm%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceCloudMonitorServiceMetricAlarmRulesConfig2)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_cms_alarm.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_cms_alarm.default.id}_fake"},
		}),
	}
	metricNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_cms_alarm.default.id}"},
			"metric_name": "${alicloud_cms_alarm.default.metric}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"metric_name": "${alicloud_cms_alarm.default.metric}_fake",
		}),
	}
	namespaceConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":       []string{"${alicloud_cms_alarm.default.id}"},
			"namespace": "${alicloud_cms_alarm.default.project}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"namespace": "${alicloud_cms_alarm.default.project}_fake",
		}),
	}
	ruleNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"rule_name": "${alicloud_cms_alarm.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"rule_name": "${alicloud_cms_alarm.default.name}_fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_cms_alarm.default.id}"},
			"status": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_cms_alarm.default.id}_fake"},
			"status": "false",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_cms_alarm.default.id}"},
			"metric_name": "${alicloud_cms_alarm.default.metric}",
			"namespace":   "${alicloud_cms_alarm.default.project}",
			"rule_name":   "${alicloud_cms_alarm.default.name}",
			"status":      "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alicloud_cms_alarm.default.id}_fake"},
			"metric_name": "${alicloud_cms_alarm.default.metric}_fake",
			"namespace":   "${alicloud_cms_alarm.default.project}_fake",
			"rule_name":   "${alicloud_cms_alarm.default.name}_fake",
			"status":      "false",
		}),
	}
	var existAliCloudCloudMonitorServiceMetricAlarmRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                    "1",
			"rules.#":                                  "1",
			"rules.0.id":                               CHECKSET,
			"rules.0.contact_groups":                   CHECKSET,
			"rules.0.dimensions":                       CHECKSET,
			"rules.0.effective_interval":               CHECKSET,
			"rules.0.email_subject":                    CHECKSET,
			"rules.0.metric_name":                      CHECKSET,
			"rules.0.namespace":                        CHECKSET,
			"rules.0.period":                           CHECKSET,
			"rules.0.rule_name":                        CHECKSET,
			"rules.0.silence_time":                     CHECKSET,
			"rules.0.source_type":                      CHECKSET,
			"rules.0.status":                           CHECKSET,
			"rules.0.webhook":                          CHECKSET,
			"rules.0.prometheus.#":                     "1",
			"rules.0.prometheus.0.prom_ql":             CHECKSET,
			"rules.0.prometheus.0.times":               CHECKSET,
			"rules.0.prometheus.0.level":               CHECKSET,
			"rules.0.prometheus.0.annotations.#":       "1",
			"rules.0.prometheus.0.annotations.0.value": CHECKSET,
			"rules.0.prometheus.0.annotations.0.key":   CHECKSET,
			"rules.0.labels.#":                         "1",
			"rules.0.labels.0.value":                   CHECKSET,
			"rules.0.labels.0.key":                     CHECKSET,
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

func dataSourceCloudMonitorServiceMetricAlarmRulesConfig2(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_instances" "default" {
  		status = "Running"
	}

	resource "alicloud_cms_alarm_contact_group" "default" {
  		alarm_contact_group_name = var.name
	}

	resource "alicloud_cms_alarm" "default" {
  		name               = var.name
  		project            = "acs_prometheus"
  		metric             = "AliyunEcs_cpu_total"
  		period             = 900
  		silence_time       = 300
  		webhook            = "https://www.aliyun.com"
  		enabled            = true
  		contact_groups     = [alicloud_cms_alarm_contact_group.default.alarm_contact_group_name]
  		effective_interval = "06:00-20:00"
  		prometheus {
    		prom_ql = var.name
    		times   = 1
    		level   = "Critical"
    		annotations = {
      			Created = "TF"
    		}
  		}
  		tags = {
    		Created = "TF"
  		}
	}
`, name)
}
