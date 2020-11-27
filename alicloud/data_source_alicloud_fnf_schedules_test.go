package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudFnfSchedulesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFnfSchedulesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_fnf_schedule.default.schedule_name}"]`,
		}),
		fakeConfig: testAccCheckAlicloudFnfSchedulesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_fnf_schedule.default.schedule_name}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFnfSchedulesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_fnf_schedule.default.schedule_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudFnfSchedulesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_fnf_schedule.default.schedule_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFnfSchedulesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_fnf_schedule.default.schedule_name}"]`,
			"name_regex": `"${alicloud_fnf_schedule.default.schedule_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudFnfSchedulesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_fnf_schedule.default.schedule_name}_fake"]`,
			"name_regex": `"${alicloud_fnf_schedule.default.schedule_name}_fake"`,
		}),
	}
	var existAlicloudFnfSchedulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"schedules.#":                 "1",
			"schedules.0.cron_expression": `30 9 * * * *`,
			"schedules.0.description":     `tf-testaccFnFSchedule983041`,
			"schedules.0.enable":          `true`,
			"schedules.0.payload":         `{"tf-test": "test success"}`,
			"schedules.0.schedule_name":   CHECKSET,
		}
	}
	var fakeAlicloudFnfSchedulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudFnfSchedulesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_fnf_schedules.default",
		existMapFunc: existAlicloudFnfSchedulesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudFnfSchedulesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.FnfSupportRegions)
	}
	alicloudFnfSchedulesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudFnfSchedulesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccSchedule-%d"
}

resource "alicloud_fnf_flow" "default" {
definition= "version: v1beta1\ntype: flow\nsteps:\n  - type: pass\n    name: helloworld"
description= "tf-testaccFnFFlow983041"
name = var.name
type= "FDL"
}

resource "alicloud_fnf_schedule" "default" {
cron_expression = "30 9 * * * *"
description = "tf-testaccFnFSchedule983041"
enable = "true"
flow_name = "${alicloud_fnf_flow.default.name}"
payload = "{\"tf-test\": \"test success\"}"
schedule_name = var.name
}

data "alicloud_fnf_schedules" "default" {	
	flow_name = "${alicloud_fnf_flow.default.name}" 
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
