package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
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
	scalingGroupIdAndNameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_alarm.default.scaling_group_id}"`,
			"name_regex":       `"${alicloud_ess_alarm.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_alarm.default.scaling_group_id}"`,
			"name_regex":       `"${alicloud_ess_alarm.default.name}_fake"`,
		}),
	}
	scalingGroupIdAndIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_alarm.default.scaling_group_id}"`,
			"ids":              `["${alicloud_ess_alarm.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_alarm.default.scaling_group_id}"`,
			"ids":              `["${alicloud_ess_alarm.default.id}_fake"]`,
		}),
	}
	idsAndNameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ess_alarm.default.id}"]`,
			"name_regex": `"${alicloud_ess_alarm.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ess_alarm.default.id}"]`,
			"name_regex": `"${alicloud_ess_alarm.default.name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_alarm.default.scaling_group_id}"`,
			"ids":              `["${alicloud_ess_alarm.default.id}"]`,
			"name_regex":       `"${alicloud_ess_alarm.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssAlarmsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_alarm.default.scaling_group_id}_fake"`,
			"ids":              `["${alicloud_ess_alarm.default.id}_fake"]`,
			"name_regex":       `"${alicloud_ess_alarm.default.name}_fake"`,
		}),
	}

	var existEssAlarmsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"alarms.#":                        "1",
			"alarms.0.name":                   "tf-testAccDataSourceEssAlarms",
			"alarms.0.scaling_group_id":       CHECKSET,
			"alarms.0.description":            "",
			"alarms.0.metric_type":            "system",
			"alarms.0.metric_name":            CHECKSET,
			"alarms.0.period":                 CHECKSET,
			"alarms.0.statistics":             CHECKSET,
			"alarms.0.threshold":              CHECKSET,
			"alarms.0.comparison_operator":    CHECKSET,
			"alarms.0.evaluation_count":       CHECKSET,
			"alarms.0.cloud_monitor_group_id": CHECKSET,
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

	essAlarmsCheckInfo.dataSourceTestCheck(t, -1, scalingGroupIdConf, nameRegexConf, idsConf, scalingGroupIdAndIdsConf, scalingGroupIdAndNameRegexConf, idsAndNameRegexConf, allConf)
}

func testAccCheckAlicloudEssAlarmsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccDataSourceEssAlarms"
}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = var.name
	vswitch_ids = ["${alicloud_vswitch.default.id}"]
}

resource "alicloud_ess_scaling_rule" "default"{
	scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
	scaling_rule_name = var.name
	adjustment_type = "QuantityChangeInCapacity"
	adjustment_value = 1
}
resource "alicloud_ess_alarm" "default"{
	scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
	name = var.name
	metric_type= "system"
    metric_name= "CpuUtilization"
    threshold = 200.3
    alarm_actions       = ["${alicloud_ess_scaling_rule.default.ari}"]
}

data "alicloud_ess_alarms" "default"{
  %s
}
`, EcsInstanceCommonTestCase, strings.Join(pairs, "\n  "))
	return config
}
