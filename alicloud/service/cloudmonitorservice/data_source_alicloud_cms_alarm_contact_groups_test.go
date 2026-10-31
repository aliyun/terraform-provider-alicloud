package cloudmonitorservice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCmsAlarmContactGroupGroups_basic(t *testing.T) {
	rand := tfacctest.RandInt()
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsAlarmContactGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cms_alarm_contact_group.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsAlarmContactGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cms_alarm_contact_group.default.id}_fake"`,
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsAlarmContactGroupsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_alarm_contact_group.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCmsAlarmContactGroupsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_alarm_contact_group.default.id}_fake"]`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsAlarmContactGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cms_alarm_contact_group.default.id}"`,
			"ids":        `["${alicloud_cms_alarm_contact_group.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCmsAlarmContactGroupsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cms_alarm_contact_group.default.id}_fake"`,
			"ids":        `["${alicloud_cms_alarm_contact_group.default.id}_fake"]`,
		}),
	}

	var existCmsAlarmContactGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"groups.#":                          "1",
			"groups.0.id":                       acctest.CHECKSET,
			"groups.0.alarm_contact_group_name": acctest.CHECKSET,
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

	var cmsAlarmContactGroupsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cms_alarm_contact_groups.default",
		ExistMapFunc: existCmsAlarmContactGroupsMapFunc,
		FakeMapFunc:  fakeCmsAlarmContactGroupsMapFunc,
	}

	cmsAlarmContactGroupsCheckInfo.DataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
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
