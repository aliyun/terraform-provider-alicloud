package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudNASAccessGroupDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100000, 999999)
	vpcTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_nas_access_group.default.access_group_name}"`,
			"access_group_type": `"${alicloud_nas_access_group.default.access_group_type}"`,
		}),
		fakeConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_nas_access_group.default.access_group_name}"`,
			"access_group_type": `"${alicloud_nas_access_group.default.access_group_type}_fake"`,
		}),
	}
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex":  `"${alicloud_nas_access_group.default.access_group_name}"`,
			"description": `"${alicloud_nas_access_group.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex":  `"${alicloud_nas_access_group.default.access_group_name}"`,
			"description": `"${alicloud_nas_access_group.default.description}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_nas_access_group.default.access_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex": `"fake"`,
		}),
	}

	accessGroupNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_nas_access_group.default.access_group_name}"`,
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"file_system_type":  `"standard"`,
		}),
		fakeConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex":        `"fake"`,
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"file_system_type":  `"standard"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_nas_access_group.default.access_group_name}"`,
			"access_group_type": `"${alicloud_nas_access_group.default.access_group_type}"`,
			"description":       `"${alicloud_nas_access_group.default.description}"`,
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"file_system_type":  `"standard"`,
			"useutc_date_time":  `true`,
		}),
		fakeConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex":        `"fake"`,
			"access_group_type": `"${alicloud_nas_access_group.default.access_group_type}"`,
			"description":       `"${alicloud_nas_access_group.default.description}"`,
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"file_system_type":  `"standard"`,
			"useutc_date_time":  `true`,
		}),
	}
	accessGroupCheckInfo.dataSourceTestCheck(t, rand, vpcTypeConf, descriptionConf, nameRegexConf, accessGroupNameConf, allConf)
}

func TestAccAlicloudNASAccessGroupDataSourceClassic(t *testing.T) {
	rand := acctest.RandIntRange(100000, 999999)
	classicTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_nas_access_group.default.access_group_name}"`,
			"access_group_type": `"${alicloud_nas_access_group.default.access_group_type}"`,
		}),
		fakeConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_nas_access_group.default.access_group_name}"`,
			"access_group_type": `"${alicloud_nas_access_group.default.access_group_type}_fake"`,
		}),
	}
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex":  `"${alicloud_nas_access_group.default.access_group_name}"`,
			"description": `"${alicloud_nas_access_group.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex":  `"${alicloud_nas_access_group.default.access_group_name}"`,
			"description": `"${alicloud_nas_access_group.default.description}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex": `"${alicloud_nas_access_group.default.access_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex": `"fake"`,
		}),
	}

	accessGroupNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_nas_access_group.default.access_group_name}"`,
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"file_system_type":  `"standard"`,
		}),
		fakeConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex":        `"fake"`,
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"file_system_type":  `"standard"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_nas_access_group.default.access_group_name}"`,
			"access_group_type": `"${alicloud_nas_access_group.default.access_group_type}"`,
			"description":       `"${alicloud_nas_access_group.default.description}"`,
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"file_system_type":  `"standard"`,
			"useutc_date_time":  `true`,
		}),
		fakeConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex":        `"fake"`,
			"access_group_type": `"${alicloud_nas_access_group.default.access_group_type}"`,
			"description":       `"${alicloud_nas_access_group.default.description}"`,
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"file_system_type":  `"standard"`,
			"useutc_date_time":  `true`,
		}),
	}
	accessGroupCheckClassicInfo.dataSourceTestCheck(t, rand, classicTypeConf, descriptionConf, nameRegexConf, accessGroupNameConf, allConf)
}

func testAccCheckAlicloudAccessGroupDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
		default = "tf-testAccAccessGroupsdatasource-%d"
}
resource "alicloud_nas_access_group" "default" {
		access_group_name = "${var.name}"
		access_group_type = "Vpc"
		description = "tf-testAccAccessGroupsdatasource"
}
data "alicloud_nas_access_groups" "default" {
		%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

func testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
		default = "tf-testAccAccessGroupsdatasource-%d"
}
resource "alicloud_nas_access_group" "default" {
		access_group_name = "${var.name}"
		access_group_type = "Classic"
		description = "tf-testAccAccessGroupsdatasource"
}
data "alicloud_nas_access_groups" "default" {
		%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existAccessGroupMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"groups.#":                    "1",
		"groups.0.rule_count":         CHECKSET,
		"groups.0.access_group_type":  "Vpc",
		"groups.0.type":               "Vpc",
		"groups.0.description":        "tf-testAccAccessGroupsdatasource",
		"groups.0.access_group_name":  fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d", rand),
		"groups.0.id":                 fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d:standard", rand),
		"groups.0.mount_target_count": CHECKSET,
		"names.#":                     "1",
		"names.0":                     fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d", rand),
	}
}

var existAccessGroupClassicMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"groups.#":                    "1",
		"groups.0.rule_count":         CHECKSET,
		"groups.0.access_group_type":  "Classic",
		"groups.0.type":               "Classic",
		"groups.0.description":        "tf-testAccAccessGroupsdatasource",
		"groups.0.access_group_name":  fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d", rand),
		"groups.0.id":                 fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d:standard", rand),
		"groups.0.mount_target_count": CHECKSET,
		"names.#":                     "1",
		"names.0":                     fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d", rand),
	}
}

var fakeAccessGroupMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"groups.#": "0",
		"names.#":  "0",
	}
}

var accessGroupCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_nas_access_groups.default",
	existMapFunc: existAccessGroupMapCheck,
	fakeMapFunc:  fakeAccessGroupMapCheck,
}

var accessGroupCheckClassicInfo = dataSourceAttr{
	resourceId:   "data.alicloud_nas_access_groups.default",
	existMapFunc: existAccessGroupClassicMapCheck,
	fakeMapFunc:  fakeAccessGroupMapCheck,
}
