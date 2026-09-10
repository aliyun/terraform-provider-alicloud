package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRdsParameterGroupsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRdsParameterGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_rds_parameter_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudRdsParameterGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_rds_parameter_group.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRdsParameterGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_rds_parameter_group.default.parameter_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudRdsParameterGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_rds_parameter_group.default.parameter_group_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRdsParameterGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_rds_parameter_group.default.id}"]`,
			"name_regex": `"${alicloud_rds_parameter_group.default.parameter_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudRdsParameterGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_rds_parameter_group.default.id}_fake"]`,
			"name_regex": `"${alicloud_rds_parameter_group.default.parameter_group_name}_fake"`,
		}),
	}
	var existAlicloudRdsParameterGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"groups.#":                      "1",
			"groups.0.engine":               `mysql`,
			"groups.0.engine_version":       `5.7`,
			"groups.0.parameter_group_desc": `test`,
			"groups.0.parameter_group_name": fmt.Sprintf("tftestAccParameterGroup%d", rand),
			"groups.0.force_restart":        "1",
			"groups.0.param_counts":         "2",
			"groups.0.param_detail.#":       "2",
			"groups.0.parameter_group_id":   CHECKSET,
			"groups.0.parameter_group_type": "1",
		}
	}
	var fakeAlicloudRdsParameterGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudRdsParameterGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_rds_parameter_groups.default",
		existMapFunc: existAlicloudRdsParameterGroupsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudRdsParameterGroupsDataSourceNameMapFunc,
	}
	alicloudRdsParameterGroupsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudRdsParameterGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tftestAccParameterGroup%d"
}

resource "alicloud_rds_parameter_group" "default" {
  engine = "mysql"
  engine_version = "5.7"
  param_detail{
    param_name = "back_log"
    param_value = "4000"
  }
  param_detail{
    param_name = "wait_timeout"
    param_value = "86460"
  }
  parameter_group_desc = "test"
  parameter_group_name = var.name
}

data "alicloud_rds_parameter_groups" "default" {	
	enable_details = true
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
