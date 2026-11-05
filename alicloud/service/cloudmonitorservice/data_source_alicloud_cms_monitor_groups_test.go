package cloudmonitorservice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCmsMonitorGroupsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_monitor_group.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_monitor_group.default.id}_fake"]`,
		}),
	}
	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_monitor_group.default.id}"]`,
			"tags": `{
				"Created" = "TF"
				"For" = "Acceptance-test"
		}`,
		}),
		FakeConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_monitor_group.default.id}"]`,
			"tags": `{
				"Created" = "TF-fake"
				"For" = "Acceptance-test"
			}`,
		}),
	}
	typeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_cms_monitor_group.default.id}"]`,
			"type": `"custom"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_cms_monitor_group.default.id}"]`,
			"type": `"kubernetes"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cms_monitor_group.default.monitor_group_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cms_monitor_group.default.monitor_group_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cms_monitor_group.default.id}"]`,
			"name_regex": `"${alicloud_cms_monitor_group.default.monitor_group_name}"`,
			"tags": `{
					"Created" = "TF"
					"For" = "Acceptance-test"
			}`,
			"type": `"custom"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsMonitorGroupsDataSourceName(rand, map[string]string{
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
			"groups.0.gmt_create":         acctest.CHECKSET,
			"groups.0.gmt_modified":       acctest.CHECKSET,
			"groups.0.monitor_group_name": acctest.CHECKSET,
			"groups.0.type":               "custom",
			"groups.0.contact_groups.#":   acctest.CHECKSET,
		}
	}
	var fakeAlicloudCmsMonitorGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCmsMonitorGroupsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cms_monitor_groups.default",
		ExistMapFunc: existAlicloudCmsMonitorGroupsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCmsMonitorGroupsDataSourceNameMapFunc,
	}
	alicloudCmsMonitorGroupsCheckInfo.DataSourceTestCheck(t, rand, idsConf, tagsConf, typeConf, nameRegexConf, allConf)
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
