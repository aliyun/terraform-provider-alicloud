package alicloud

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGpdbDbInstancePlansDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbDbInstancePlansDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_gpdb_db_instance_plan.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbDbInstancePlansDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_gpdb_db_instance_plan.default.id}_fake"]`,
		}),
	}
	planScheduleTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbDbInstancePlansDataSourceName(rand, map[string]string{
			"ids":                `["${alicloud_gpdb_db_instance_plan.default.id}"]`,
			"plan_schedule_type": `"${alicloud_gpdb_db_instance_plan.default.plan_schedule_type}"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbDbInstancePlansDataSourceName(rand, map[string]string{
			"ids":                `["${alicloud_gpdb_db_instance_plan.default.id}"]`,
			"plan_schedule_type": `"Postpone"`,
		}),
	}
	planTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbDbInstancePlansDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_gpdb_db_instance_plan.default.id}"]`,
			"plan_type": `"${alicloud_gpdb_db_instance_plan.default.plan_type}"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbDbInstancePlansDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_gpdb_db_instance_plan.default.id}"]`,
			"plan_type": `"Resize"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbDbInstancePlansDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_gpdb_db_instance_plan.default.db_instance_plan_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbDbInstancePlansDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_gpdb_db_instance_plan.default.db_instance_plan_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbDbInstancePlansDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_gpdb_db_instance_plan.default.id}"]`,
			"status": `"${alicloud_gpdb_db_instance_plan.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbDbInstancePlansDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_gpdb_db_instance_plan.default.id}"]`,
			"status": `"cancel"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGpdbDbInstancePlansDataSourceName(rand, map[string]string{
			"ids":                `["${alicloud_gpdb_db_instance_plan.default.id}"]`,
			"name_regex":         `"${alicloud_gpdb_db_instance_plan.default.db_instance_plan_name}"`,
			"plan_schedule_type": `"${alicloud_gpdb_db_instance_plan.default.plan_schedule_type}"`,
			"plan_type":          `"${alicloud_gpdb_db_instance_plan.default.plan_type}"`,
			"status":             `"${alicloud_gpdb_db_instance_plan.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudGpdbDbInstancePlansDataSourceName(rand, map[string]string{
			"ids":                `["${alicloud_gpdb_db_instance_plan.default.id}_fake"]`,
			"name_regex":         `"${alicloud_gpdb_db_instance_plan.default.db_instance_plan_name}_fake"`,
			"plan_schedule_type": `"Postpone"`,
			"plan_type":          `"Resize"`,
			"status":             `"cancel"`,
		}),
	}
	var existAlicloudGpdbDbInstancePlansDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"plans.#":                        "1",
			"plans.0.id":                     CHECKSET,
			"plans.0.plan_id":                CHECKSET,
			"plans.0.status":                 CHECKSET,
			"plans.0.plan_desc":              fmt.Sprintf("tf-testAccDBInstancePlan-%d", rand),
			"plans.0.plan_end_date":          CHECKSET,
			"plans.0.db_instance_plan_name":  fmt.Sprintf("tf-testAccDBInstancePlan-%d", rand),
			"plans.0.plan_schedule_type":     "Regular",
			"plans.0.plan_start_date":        CHECKSET,
			"plans.0.plan_type":              "PauseResume",
			"plans.0.plan_config.#":          "1",
			"plans.0.plan_config.0.resume.#": "1",
			"plans.0.plan_config.0.resume.0.plan_cron_time": "0 0 0 1/1 * ? ",
			"plans.0.plan_config.0.pause.#":                 "1",
			"plans.0.plan_config.0.pause.0.plan_cron_time":  "0 0 10 1/1 * ? ",
		}
	}
	var fakeAlicloudGpdbDbInstancePlansDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudGpdbDbInstancePlansCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_gpdb_db_instance_plans.default",
		existMapFunc: existAlicloudGpdbDbInstancePlansDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudGpdbDbInstancePlansDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudGpdbDbInstancePlansCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, planScheduleTypeConf, planTypeConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudGpdbDbInstancePlansDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	planStartDate := time.Now().Format("2006-01-02T15:04:05Z")
	planEndDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")
	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccDBInstancePlan-%d"
}

data "alicloud_gpdb_instances" "default" {	
	name_regex = "default-NODELETING"
}

resource "alicloud_gpdb_db_instance_plan" "default" {
	db_instance_plan_name =           "${var.name}"
	plan_desc =           "${var.name}"
	plan_type =           "PauseResume"
	plan_schedule_type =  "Regular"
	plan_start_date =     "%s"
	plan_end_date =       "%s"
	plan_config {
		resume {
				plan_cron_time =  "0 0 0 1/1 * ? "
			}
		pause {
				plan_cron_time =  "0 0 10 1/1 * ? "
			}
		}
	db_instance_id = data.alicloud_gpdb_instances.default.ids.0
}

data "alicloud_gpdb_db_instance_plans" "default" {	
	db_instance_id = data.alicloud_gpdb_instances.default.ids.0
	%s
}
`, rand, planStartDate, planEndDate, strings.Join(pairs, " \n "))
	return config
}
