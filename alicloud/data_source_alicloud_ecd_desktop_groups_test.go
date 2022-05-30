package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"strings"
	"testing"
)

func TestAccAlicloudEcdDesktopGroupsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.ECDSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_desktop_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_desktop_group.default.id}_fake"]`,
		}),
	}
	desktopGroupNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
			"ids":                `["${alicloud_ecd_desktop_group.default.id}"]`,
			"desktop_group_name": `"${alicloud_ecd_desktop_group.default.desktop_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
			"ids":                `["${alicloud_ecd_desktop_group.default.id}"]`,
			"desktop_group_name": `"${alicloud_ecd_desktop_group.default.desktop_group_name}_fake"`,
		}),
	}

	excludedEndUserIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_ecd_desktop_group.default.id}"]`,
			"excluded_end_user_ids": `"${alicloud_ecd_desktop_group.default.excluded_end_user_ids}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_ecd_desktop_group.default.id}"]`,
			"excluded_end_user_ids": `"${alicloud_ecd_desktop_group.default.excluded_end_user_ids}_fake"`,
		}),
	}

	officeSiteIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_ecd_desktop_group.default.id}"]`,
			"office_site_id": `"${alicloud_ecd_desktop_group.default.office_site_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_ecd_desktop_group.default.id}"]`,
			"office_site_id": `"${alicloud_ecd_desktop_group.default.office_site_id}_fake"`,
		}),
	}
	ownTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_ecd_desktop_group.default.id}"]`,
			"own_type": `"${alicloud_ecd_desktop_group.default.own_type}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_ecd_desktop_group.default.id}"]`,
			"own_type": `"${alicloud_ecd_desktop_group.default.own_type}_fake"`,
		}),
	}
	periodConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecd_desktop_group.default.id}"]`,
			"period": `"${alicloud_ecd_desktop_group.default.period}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ecd_desktop_group.default.id}"]`,
			"period": `"${alicloud_ecd_desktop_group.default.period}_fake"`,
		}),
	}
	periodUnitConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_ecd_desktop_group.default.id}"]`,
			"period_unit": `"${alicloud_ecd_desktop_group.default.period_unit}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_ecd_desktop_group.default.id}"]`,
			"period_unit": `"${alicloud_ecd_desktop_group.default.period_unit}_fake"`,
		}),
	}
	//periodUnitConf := dataSourceTestAccConfig{
	//	existConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
	//		"ids":    `["${alicloud_ecd_desktop_group.default.id}"]`,
	//		"status": `"${alicloud_ecd_desktop_group.default.status}"`,
	//	}),
	//	fakeConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
	//		"ids":    `["${alicloud_ecd_desktop_group.default.id}"]`,
	//		"status": `"${alicloud_ecd_desktop_group.default.status}_fake"`,
	//	}),
	//}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_desktop_group.default.desktop_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_desktop_group.default.desktop_group_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
			"desktop_group_name":    `"${alicloud_ecd_desktop_group.default.desktop_group_name}"`,
			"excluded_end_user_ids": `"${alicloud_ecd_desktop_group.default.excluded_end_user_ids}"`,
			"ids":                   `["${alicloud_ecd_desktop_group.default.id}"]`,
			"name_regex":            `"${alicloud_ecd_desktop_group.default.desktop_group_name}"`,
			"office_site_id":        `"${alicloud_ecd_desktop_group.default.office_site_id}"`,
			"own_type":              `"${alicloud_ecd_desktop_group.default.own_type}"`,
			"period":                `"${alicloud_ecd_desktop_group.default.period}"`,
			"period_unit":           `"${alicloud_ecd_desktop_group.default.period_unit}"`,
			"status":                `"${alicloud_ecd_desktop_group.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand, map[string]string{
			"desktop_group_name":    `"${alicloud_ecd_desktop_group.default.desktop_group_name}_fake"`,
			"excluded_end_user_ids": `"${alicloud_ecd_desktop_group.default.excluded_end_user_ids}_fake"`,
			"ids":                   `["${alicloud_ecd_desktop_group.default.id}_fake"]`,
			"name_regex":            `"${alicloud_ecd_desktop_group.default.desktop_group_name}_fake"`,
			"office_site_id":        `"${alicloud_ecd_desktop_group.default.office_site_id}_fake"`,
			"own_type":              `"${alicloud_ecd_desktop_group.default.own_type}_fake"`,
			"period":                `"${alicloud_ecd_desktop_group.default.period}_fake"`,
			"period_unit":           `"${alicloud_ecd_desktop_group.default.period_unit}_fake"`,
			"status":                `"${alicloud_ecd_desktop_group.default.status}_fake"`,
		}),
	}
	var existAlicloudEcdDesktopGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"groups.#":                    "1",
			"groups.0.allow_auto_setup":   `406`,
			"groups.0.allow_buffer_count": `452`,
			"groups.0.bundle_id":          `tf-testAcc-TSnZo`,
			"groups.0.comments":           `tf-testAcc-lxamC`,
			"groups.0.desktop_group_name": CHECKSET,
			"groups.0.directory_id":       `tf-testAcc-FJgX3`,
			"groups.0.end_user_ids":       `<nil>`,
			"groups.0.keep_duration":      `440`,
			"groups.0.max_desktops_count": `776`,
			"groups.0.min_desktops_count": `984`,
			"groups.0.office_site_id":     `tf-testAcc-ZYq3e`,
			"groups.0.policy_group_id":    `tf-testAcc-OafWP`,
			"groups.0.scale_strategy_id":  `tf-testAcc-ZAaKk`,
		}
	}
	var fakeAlicloudEcdDesktopGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEcdDesktopGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecd_desktop_groups.default",
		existMapFunc: existAlicloudEcdDesktopGroupsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcdDesktopGroupsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEcdDesktopGroupsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, desktopGroupNameConf, excludedEndUserIdsConf, officeSiteIdConf, ownTypeConf, periodConf, periodUnitConf, periodUnitConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudEcdDesktopGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccDesktopGroup-%d"
}

resource "alicloud_ecd_desktop_group" "default" {
allow_auto_setup = "406"
allow_buffer_count = "452"
bundle_id = "tf-testAcc-TSnZo"
comments = "tf-testAcc-lxamC"
desktop_group_name = var.name
directory_id = "tf-testAcc-FJgX3"
end_user_ids = "<nil>"
keep_duration = "440"
max_desktops_count = "776"
min_desktops_count = "984"
office_site_id = "tf-testAcc-ZYq3e"
policy_group_id = "tf-testAcc-OafWP"
scale_strategy_id = "tf-testAcc-ZAaKk"
}

data "alicloud_ecd_desktop_groups" "default" {	
	enable_details = true
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
