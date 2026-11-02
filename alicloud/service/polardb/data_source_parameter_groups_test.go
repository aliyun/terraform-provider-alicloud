package polardb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudPolarDBParameterGroupsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPolarDBParameterGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_polardb_parameter_group.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudPolarDBParameterGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_polardb_parameter_group.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPolarDBParameterGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_polardb_parameter_group.default.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudPolarDBParameterGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_polardb_parameter_group.default.name}_fake"`,
		}),
	}
	dbTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPolarDBParameterGroupsDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_polardb_parameter_group.default.id}"]`,
			"db_type": `"${alicloud_polardb_parameter_group.default.db_type}"`,
		}),
		FakeConfig: testAccCheckAlicloudPolarDBParameterGroupsDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_polardb_parameter_group.default.id}_fake"]`,
			"db_type": `"${alicloud_polardb_parameter_group.default.db_type}"`,
		}),
	}
	dbVersionConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPolarDBParameterGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_polardb_parameter_group.default.id}"]`,
			"db_version": `"${alicloud_polardb_parameter_group.default.db_version}"`,
		}),
		FakeConfig: testAccCheckAlicloudPolarDBParameterGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_polardb_parameter_group.default.id}_fake"]`,
			"db_version": `"${alicloud_polardb_parameter_group.default.db_version}"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPolarDBParameterGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_polardb_parameter_group.default.id}"]`,
			"name_regex": `"${alicloud_polardb_parameter_group.default.name}"`,
			"db_type":    `"MySQL"`,
			"db_version": `"8.0"`,
		}),
		FakeConfig: testAccCheckAlicloudPolarDBParameterGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_polardb_parameter_group.default.id}_fake"]`,
			"name_regex": `"${alicloud_polardb_parameter_group.default.name}_fake"`,
			"db_type":    `"MySQL"`,
			"db_version": `"8.0"`,
		}),
	}
	var existAlicloudPolarDBParameterGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"groups.#":                      "1",
			"groups.0.id":                   CHECKSET,
			"groups.0.parameter_group_id":   CHECKSET,
			"groups.0.parameter_group_name": CHECKSET,
			"groups.0.db_type":              "MySQL",
			"groups.0.db_version":           "8.0",
			"groups.0.parameter_group_desc": CHECKSET,
			"groups.0.parameter_group_type": CHECKSET,
			"groups.0.parameter_counts":     CHECKSET,
			"groups.0.force_restart":        CHECKSET,
			"groups.0.create_time":          CHECKSET,
		}
	}
	var fakeAlicloudPolarDBParameterGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"groups.#": "0",
		}
	}
	var alicloudPolarDBParameterGroupsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_polardb_parameter_groups.default",
		ExistMapFunc: existAlicloudPolarDBParameterGroupsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudPolarDBParameterGroupsDataSourceNameMapFunc,
	}
	alicloudPolarDBParameterGroupsCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, dbTypeConf, dbVersionConf, allConf)
}

func testAccCheckAlicloudPolarDBParameterGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
		default = "tf_testAcc_%d"
	}
	
	resource "alicloud_polardb_parameter_group" "default" {
		name       = var.name
		db_type    = "MySQL"
		db_version = "8.0"
		parameters {
			param_name  = "wait_timeout"
			param_value = "86400"
		}
		parameters {
			param_name  = "innodb_old_blocks_time"
			param_value = "1000"
		}
		description = var.name
	}

	data "alicloud_polardb_parameter_groups" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
