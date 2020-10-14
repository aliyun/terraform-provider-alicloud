package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCmsAlarmContactGroupGroups_basic(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsAlarmContactGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cms_alarm_contact_group.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsAlarmContactGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cms_alarm_contact_group.default.id}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsAlarmContactGroupsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_alarm_contact_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsAlarmContactGroupsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_alarm_contact_group.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsAlarmContactGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cms_alarm_contact_group.default.id}"`,
			"ids":        `["${alicloud_cms_alarm_contact_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsAlarmContactGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cms_alarm_contact_group.default.id}_fake"`,
			"ids":        `["${alicloud_cms_alarm_contact_group.default.id}_fake"]`,
		}),
	}

	var existCmsAlarmContactGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"groups.#":                          "1",
			"groups.0.id":                       CHECKSET,
			"groups.0.alarm_contact_group_name": CHECKSET,
			"groups.0.describe":                 "For Test",
			"groups.0.enable_subscribed":        "true",
		}
	}

	var fakeCmsAlarmContactGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var cmsAlarmContactGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_alarm_contact_groups.default",
		existMapFunc: existCmsAlarmContactGroupsMapFunc,
		fakeMapFunc:  fakeCmsAlarmContactGroupsMapFunc,
	}

	cmsAlarmContactGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudCmsAlarmContactGroupsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
		variable "name" {
			default = "tf-testAccCmsAlarmContactGroupBisic-%d"
		}
		resource "alicloud_cms_alarm_contact_group" "default" {
		  alarm_contact_group_name = var.name
		  describe                 = "For Test"
		  enable_subscribed        = true
		}

		data "alicloud_cms_alarm_contact_groups" "default" {
		  %s
		}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
