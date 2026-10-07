// Package alicloud. This file is hand-written from the Cms 2024-03-30 OpenAPI definition.
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCmsOncallSchedulesDataSource_nameRegex(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsOncallSchedulesDataSourceName(rand, map[string]string{
			"oncall_schedule_name_regex": `"${alicloud_cms_oncall_schedule.default.oncall_schedule_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsOncallSchedulesDataSourceName(rand, map[string]string{
			"oncall_schedule_name_regex": `"${alicloud_cms_oncall_schedule.default.oncall_schedule_name}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsOncallSchedulesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_oncall_schedule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsOncallSchedulesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_oncall_schedule.default.id}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsOncallSchedulesDataSourceName(rand, map[string]string{
			"oncall_schedule_name_regex": `"${alicloud_cms_oncall_schedule.default.oncall_schedule_name}"`,
			"ids":                        `["${alicloud_cms_oncall_schedule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsOncallSchedulesDataSourceName(rand, map[string]string{
			"oncall_schedule_name_regex": `"${alicloud_cms_oncall_schedule.default.oncall_schedule_name}_fake"`,
			"ids":                        `["${alicloud_cms_oncall_schedule.default.id}_fake"]`,
		}),
	}
	var existAlicloudCmsOncallSchedulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"schedules.#":                      "1",
			"schedules.0.oncall_schedule_id":   CHECKSET,
			"schedules.0.oncall_schedule_name": CHECKSET,
		}
	}
	var fakeAlicloudCmsOncallSchedulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"schedules.#": "0",
		}
	}
	var alicloudCmsOncallSchedulesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_oncall_schedules.default",
		existMapFunc: existAlicloudCmsOncallSchedulesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCmsOncallSchedulesDataSourceNameMapFunc,
	}
	alicloudCmsOncallSchedulesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudCmsOncallSchedulesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tftestacccmsoncall%d"
}

resource "alicloud_cms_oncall_schedule" "default" {
  oncall_schedule_name = var.name
  source               = var.name
}

data "alicloud_cms_oncall_schedules" "default" {
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
