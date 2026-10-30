package nas_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudNASAccessGroupDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100000, 999999)
	vpcTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_nas_access_group.default.access_group_name}"`,
			"access_group_type": `"${alicloud_nas_access_group.default.access_group_type}"`,
		}),
		FakeConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_nas_access_group.default.access_group_name}"`,
			"access_group_type": `"${alicloud_nas_access_group.default.access_group_type}_fake"`,
		}),
	}
	descriptionConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex":  `"${alicloud_nas_access_group.default.access_group_name}"`,
			"description": `"${alicloud_nas_access_group.default.description}"`,
		}),
		FakeConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex":  `"${alicloud_nas_access_group.default.access_group_name}"`,
			"description": `"${alicloud_nas_access_group.default.description}_fake"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_nas_access_group.default.access_group_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex": `"fake"`,
		}),
	}

	accessGroupNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_nas_access_group.default.access_group_name}"`,
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"file_system_type":  `"standard"`,
		}),
		FakeConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex":        `"fake"`,
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"file_system_type":  `"standard"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_nas_access_group.default.access_group_name}"`,
			"access_group_type": `"${alicloud_nas_access_group.default.access_group_type}"`,
			"description":       `"${alicloud_nas_access_group.default.description}"`,
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"file_system_type":  `"standard"`,
			"useutc_date_time":  `true`,
		}),
		FakeConfig: testAccCheckAlicloudAccessGroupDataSourceConfig(rand, map[string]string{
			"name_regex":        `"fake"`,
			"access_group_type": `"${alicloud_nas_access_group.default.access_group_type}"`,
			"description":       `"${alicloud_nas_access_group.default.description}"`,
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"file_system_type":  `"standard"`,
			"useutc_date_time":  `true`,
		}),
	}
	accessGroupCheckInfo.DataSourceTestCheck(t, rand, vpcTypeConf, descriptionConf, nameRegexConf, accessGroupNameConf, allConf)
}

func TestAccAliCloudNASAccessGroupDataSourceClassic(t *testing.T) {
	rand := tfacctest.RandIntRange(100000, 999999)
	classicTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_nas_access_group.default.access_group_name}"`,
			"access_group_type": `"${alicloud_nas_access_group.default.access_group_type}"`,
		}),
		FakeConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_nas_access_group.default.access_group_name}"`,
			"access_group_type": `"${alicloud_nas_access_group.default.access_group_type}_fake"`,
		}),
	}
	descriptionConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex":  `"${alicloud_nas_access_group.default.access_group_name}"`,
			"description": `"${alicloud_nas_access_group.default.description}"`,
		}),
		FakeConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex":  `"${alicloud_nas_access_group.default.access_group_name}"`,
			"description": `"${alicloud_nas_access_group.default.description}_fake"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex": `"${alicloud_nas_access_group.default.access_group_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex": `"fake"`,
		}),
	}

	accessGroupNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_nas_access_group.default.access_group_name}"`,
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"file_system_type":  `"standard"`,
		}),
		FakeConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex":        `"fake"`,
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"file_system_type":  `"standard"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_nas_access_group.default.access_group_name}"`,
			"access_group_type": `"${alicloud_nas_access_group.default.access_group_type}"`,
			"description":       `"${alicloud_nas_access_group.default.description}"`,
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"file_system_type":  `"standard"`,
			"useutc_date_time":  `true`,
		}),
		FakeConfig: testAccCheckAlicloudAccessGroupDataSourceClassicConfig(rand, map[string]string{
			"name_regex":        `"fake"`,
			"access_group_type": `"${alicloud_nas_access_group.default.access_group_type}"`,
			"description":       `"${alicloud_nas_access_group.default.description}"`,
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"file_system_type":  `"standard"`,
			"useutc_date_time":  `true`,
		}),
	}
	accessGroupCheckClassicInfo.DataSourceTestCheck(t, rand, classicTypeConf, descriptionConf, nameRegexConf, accessGroupNameConf, allConf)
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
		"groups.0.rule_count":         acctest.CHECKSET,
		"groups.0.access_group_type":  "Vpc",
		"groups.0.type":               "Vpc",
		"groups.0.description":        "tf-testAccAccessGroupsdatasource",
		"groups.0.access_group_name":  fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d", rand),
		"groups.0.id":                 fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d:standard", rand),
		"groups.0.mount_target_count": acctest.CHECKSET,
		"names.#":                     "1",
		"names.0":                     fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d", rand),
	}
}

var existAccessGroupClassicMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"groups.#":                    "1",
		"groups.0.rule_count":         acctest.CHECKSET,
		"groups.0.access_group_type":  "Classic",
		"groups.0.type":               "Classic",
		"groups.0.description":        "tf-testAccAccessGroupsdatasource",
		"groups.0.access_group_name":  fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d", rand),
		"groups.0.id":                 fmt.Sprintf("tf-testAccAccessGroupsdatasource-%d:standard", rand),
		"groups.0.mount_target_count": acctest.CHECKSET,
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

var accessGroupCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_nas_access_groups.default",
	ExistMapFunc: existAccessGroupMapCheck,
	FakeMapFunc:  fakeAccessGroupMapCheck,
}

var accessGroupCheckClassicInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_nas_access_groups.default",
	ExistMapFunc: existAccessGroupClassicMapCheck,
	FakeMapFunc:  fakeAccessGroupMapCheck,
}
