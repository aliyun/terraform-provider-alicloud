package alicloud

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEssScheduledtasksDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"scheduled_task_id": `"${alicloud_ess_scheduled_task.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"scheduled_task_id": `"${alicloud_ess_scheduled_task.default.id}_fake"`,
		}),
	}
	actionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"scheduled_action": `"${alicloud_ess_scheduled_task.default.scheduled_action}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"scheduled_action": `"${alicloud_ess_scheduled_task.default.scheduled_action}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ess_scheduled_task.default.scheduled_task_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ess_scheduled_task.default.scheduled_task_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ess_scheduled_task.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ess_scheduled_task.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"scheduled_action":  `"${alicloud_ess_scheduled_task.default.scheduled_action}"`,
			"ids":               `["${alicloud_ess_scheduled_task.default.id}"]`,
			"name_regex":        `"${alicloud_ess_scheduled_task.default.scheduled_task_name}"`,
			"scheduled_task_id": `"${alicloud_ess_scheduled_task.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEssScheduledTasksDataSourceConfig(rand, map[string]string{
			"scheduled_action":  `"${alicloud_ess_scheduled_task.default.scheduled_action}_fake"`,
			"ids":               `["${alicloud_ess_scheduled_task.default.id}"]`,
			"name_regex":        `"${alicloud_ess_scheduled_task.default.scheduled_task_name}"`,
			"scheduled_task_id": `"${alicloud_ess_scheduled_task.default.id}"`,
		}),
	}

	var existEssScheduledTasksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"tasks.#":                        "1",
			"tasks.0.name":                   fmt.Sprintf("tf-testAccDataSourceEssScheduledTasks-%d", rand),
			"tasks.0.id":                     CHECKSET,
			"tasks.0.scheduled_action":       CHECKSET,
			"tasks.0.launch_expiration_time": CHECKSET,
			"tasks.0.launch_time":            CHECKSET,
			"tasks.0.max_value":              CHECKSET,
			"tasks.0.min_value":              CHECKSET,
			"tasks.0.task_enabled":           CHECKSET,
		}
	}

	var fakeEssScheduledTasksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"tasks.#": "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var essScheduledTasksCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_scheduled_tasks.default",
		existMapFunc: existEssScheduledTasksMapFunc,
		fakeMapFunc:  fakeEssScheduledTasksMapFunc,
	}

	essScheduledTasksCheckInfo.dataSourceTestCheck(t, rand, idConf, actionConf, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudEssScheduledTasksDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	oneDay, _ := time.ParseDuration("24h")

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccDataSourceEssScheduledTasks-%d"
}

resource "alicloud_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${alicloud_vswitch.default.id}"]
}
resource "alicloud_ess_scaling_rule" "default" {
  scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
  adjustment_type  = "TotalCapacity"
  adjustment_value = 2
  cooldown         = 60
}

resource "alicloud_ess_scheduled_task" "default" {
  scheduled_action    = "${alicloud_ess_scaling_rule.default.ari}"
  launch_time         = "%s"
  scheduled_task_name = "${var.name}"
}

data "alicloud_ess_scheduled_tasks" "default"{
  %s
}
`, EcsInstanceCommonTestCase, rand, time.Now().Add(oneDay).Format("2006-01-02T15:04Z"), strings.Join(pairs, "\n  "))
	return config
}
