package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudEssAlarmsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	scalingGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_alarm.default.scaling_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_alarm.default.scaling_group_id}_fake"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ess_alarm.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ess_alarm.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ess_alarm.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ess_alarm.default.id}_fake"]`,
		}),
	}
	isEnableConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"enable": `"${alicloud_ess_alarm.default.IsEnable}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"enable": `"${alicloud_ess_alarm.default.IsEnable}_fake"`,
		}),
	}
	metricTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"Metric_type": `"${alicloud_ess_alarm.default.Metric_type}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"Metric_type": `"${alicloud_ess_alarm.default.Metric_type}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_scaling_configuration.default.scaling_group_id}"`,
			"ids":              `["${alicloud_ess_alarm.default.id}"]`,
			"name_regex":       `"${alicloud_ess_alarm.default.alarm_name}"`,
			"IsEnable":         `"${alicloud_ess_alarm.default.IsEnable}"`,
			"Metric_type":      `"${alicloud_ess_alarm.default.Metric_type}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_alarm.default.scaling_group_id}"`,
			"ids":              `["${alicloud_ess_alarm.default.id}_fake"]`,
			"name_regex":       `"${alicloud_ess_alarm.default.alarm_name}"`,
			"IsEnable":         `"${alicloud_ess_alarm.default.IsEnable}"`,
			"Metric_type":      `"${alicloud_ess_alarm.default.Metric_type}"`,
		}),
	}

	var existEssAlarmsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"alarms.#":                        "1",
			"alarms.0.name":                   fmt.Sprintf("tf-testAccDataSourceEssAlarms-%d", rand),
			"alarms.0.scaling_group_id":       CHECKSET,
			"alarms.0.description":            CHECKSET,
			"alarms.0.enable":                 CHECKSET,
			"alarms.0.alarm_actions":          CHECKSET,
			"alarms.0.metric_type":            CHECKSET,
			"alarms.0.metric_name":            CHECKSET,
			"alarms.0.period":                 CHECKSET,
			"alarms.0.statistics":             CHECKSET,
			"alarms.0.threshold":              CHECKSET,
			"alarms.0.comparison_operator":    CHECKSET,
			"alarms.0.evaluation_count":       CHECKSET,
			"alarms.0.cloud_monitor_group_id": CHECKSET,
			"alarms.0.dimensions":             CHECKSET,
			"alarms.0.state":                  CHECKSET,
		}
	}

	var fakeEssAlarmsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"alarms.#": "0",
			"ids.#":    "0",
			"names.#":  "0",
		}
	}

	var essAlarmsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_alarms.default",
		existMapFunc: existEssAlarmsMapFunc,
		fakeMapFunc:  fakeEssAlarmsMapFunc,
	}

	essAlarmsCheckInfo.dataSourceTestCheck(t, rand, scalingGroupIdConf, nameRegexConf, isEnableConf, metricTypeConf, idsConf, allConf)
}

func testAccCheckAlicloudEssAlarmsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccDataSourceEssAlarms-%d"
}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${alicloud_vswitch.default.id}"]
}

resource "alicloud_ess_alarm" "default"{
	scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
	name = "${var.name}"
	metric_type= "system"
    state= OK
}

data "alicloud_ess_alarms" "default"{
  %s
}
`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
