package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCmsDynamicTagGroups_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.CmsDynamicTagGroupSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsDynamicTagGroupsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_dynamic_tag_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsDynamicTagGroupsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cms_dynamic_tag_group.default.id}_fake"]`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsDynamicTagGroupsDataSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cms_dynamic_tag_group.default.id}"]`,
			"status": `"FINISH"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsDynamicTagGroupsDataSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cms_dynamic_tag_group.default.id}"]`,
			"status": `"RUNNING"`,
		}),
	}

	tagKeyConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsDynamicTagGroupsDataSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_cms_dynamic_tag_group.default.id}"]`,
			"tag_key": `"${alicloud_cms_dynamic_tag_group.default.tag_key}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsDynamicTagGroupsDataSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_cms_dynamic_tag_group.default.id}"]`,
			"tag_key": `"${alicloud_cms_dynamic_tag_group.default.tag_key}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsDynamicTagGroupsDataSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_cms_dynamic_tag_group.default.id}"]`,
			"status":  `"FINISH"`,
			"tag_key": `"${alicloud_cms_dynamic_tag_group.default.tag_key}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsDynamicTagGroupsDataSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_cms_dynamic_tag_group.default.id}_fake"]`,
			"status":  `"RUNNING"`,
			"tag_key": `"${alicloud_cms_dynamic_tag_group.default.tag_key}_fake"`,
		}),
	}

	var existCmsDynamicTagGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"groups.#":                               "1",
			"groups.0.id":                            CHECKSET,
			"groups.0.dynamic_tag_rule_id":           CHECKSET,
			"groups.0.status":                        CHECKSET,
			"groups.0.tag_key":                       CHECKSET,
			"groups.0.match_express_filter_relation": CHECKSET,
			"groups.0.match_express.0.tag_value":     CHECKSET,
			"groups.0.match_express.0.tag_value_match_function": CHECKSET,
		}
	}

	var fakeCmsDynamicTagGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}

	var cmsDynamicTagGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_dynamic_tag_groups.default",
		existMapFunc: existCmsDynamicTagGroupsMapFunc,
		fakeMapFunc:  fakeCmsDynamicTagGroupsMapFunc,
	}

	cmsDynamicTagGroupsCheckInfo.dataSourceTestCheck(t, rand, idsConf, statusConf, tagKeyConf, allConf)
}

func testAccCheckAlicloudCmsDynamicTagGroupsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
		variable "name" {
			default = " tf-testacc%d"
		}
		resource "alicloud_cms_alarm_contact_group" "default" {
		  alarm_contact_group_name = var.name
		  describe                 = "For Test"
		  enable_subscribed        = true
		}
		resource "alicloud_cms_dynamic_tag_group" "default" {
		  	contact_group_list  = [alicloud_cms_alarm_contact_group.default.id]
		  	tag_key = "appgroup"
		  	match_express{
				tag_value = "landingzone"
                tag_value_match_function = "all"
			}
		}

		data "alicloud_cms_dynamic_tag_groups" "default" {
		  %s
		}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
