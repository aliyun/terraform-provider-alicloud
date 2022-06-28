package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCmsMonitorGroupsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_monitor_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_monitor_group.default.id}_fake"]`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_monitor_group.default.id}"]`,
			"tags": `{
				"Created" = "TF"
				"For" = "Acceptance-test"
		}`,
		}),
		fakeConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_monitor_group.default.id}"]`,
			"tags": `{
				"Created" = "TF-fake"
				"For" = "Acceptance-test"
			}`,
		}),
	}
	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_cms_monitor_group.default.id}"]`,
			"type": `"custom"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_cms_monitor_group.default.id}"]`,
			"type": `"kubernetes"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cms_monitor_group.default.monitor_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cms_monitor_group.default.monitor_group_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cms_monitor_group.default.id}"]`,
			"name_regex": `"${alicloud_cms_monitor_group.default.monitor_group_name}"`,
			"tags": `{
					"Created" = "TF"
					"For" = "Acceptance-test"
			}`,
			"type": `"custom"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cms_monitor_group.default.id}_fake"]`,
			"name_regex": `"${alicloud_cms_monitor_group.default.monitor_group_name}_fake"`,
			"tags": `{
				"Created" = "TF-fake"
				"For" = "Acceptance-test"
			}`,
			"type": `"kubernetes"`,
		}),
	}
	var existAlicloudCmsMonitorGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"groups.#":                    "1",
			"groups.0.gmt_create":         CHECKSET,
			"groups.0.gmt_modified":       CHECKSET,
			"groups.0.monitor_group_name": CHECKSET,
			"groups.0.type":               "custom",
			"groups.0.contact_groups.#":   CHECKSET,
		}
	}
	var fakeAlicloudCmsMonitorGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCmsMonitorGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_monitor_groups.default",
		existMapFunc: existAlicloudCmsMonitorGroupsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCmsMonitorGroupsDataSourceNameMapFunc,
	}
	alicloudCmsMonitorGroupsCheckInfo.dataSourceTestCheck(t, rand, idsConf, tagsConf, typeConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccMonitorGroup-%d"
}

resource "alicloud_cms_alarm_contact_group" "default" {
alarm_contact_group_name = var.name
}

resource "alicloud_cms_monitor_group" "default" {
monitor_group_name = var.name
contact_groups = ["${alicloud_cms_alarm_contact_group.default.alarm_contact_group_name}"]
tags = {
		Created = "TF"
		For = "Acceptance-test"
 }
}

data "alicloud_cms_monitor_groups" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
