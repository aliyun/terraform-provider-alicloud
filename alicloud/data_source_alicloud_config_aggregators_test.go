package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudConfigAggregatorsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigAggregatorsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_config_aggregator.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudConfigAggregatorsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_config_aggregator.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigAggregatorsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_aggregator.default.aggregator_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigAggregatorsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_aggregator.default.aggregator_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigAggregatorsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_config_aggregator.default.id}"]`,
			"status": `"Normal"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigAggregatorsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_config_aggregator.default.id}"]`,
			"status": `"Deleting"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigAggregatorsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_config_aggregator.default.id}"]`,
			"name_regex": `"${alicloud_config_aggregator.default.aggregator_name}"`,
			"status":     `"Normal"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigAggregatorsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_config_aggregator.default.id}_fake"]`,
			"name_regex": `"${alicloud_config_aggregator.default.aggregator_name}_fake"`,
			"status":     `"Deleting"`,
		}),
	}
	var existAlicloudConfigAggregatorsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"aggregators.#":                       "1",
			"aggregators.0.aggregator_accounts.#": "1",
			"aggregators.0.aggregator_name":       CHECKSET,
			"aggregators.0.description":           `tf-create-aggregator`,
			"aggregators.0.status":                `Normal`,
		}
	}
	var fakeAlicloudConfigAggregatorsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudConfigAggregatorsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_config_aggregators.default",
		existMapFunc: existAlicloudConfigAggregatorsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudConfigAggregatorsDataSourceNameMapFunc,
	}

	PreCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckEnterpriseAccountEnabled(t)
	}

	alicloudConfigAggregatorsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, PreCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudConfigAggregatorsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccAggregator-%d"
}

data "alicloud_resource_manager_accounts" "default" {
  status  = "CreateSuccess"
}

resource "alicloud_config_aggregator" "default" {
	aggregator_accounts {
	account_id   =  data.alicloud_resource_manager_accounts.default.accounts.0.account_id
	account_name =  data.alicloud_resource_manager_accounts.default.accounts.0.display_name
	account_type = "ResourceDirectory"
	}
	aggregator_name = var.name
	description = "tf-create-aggregator"
}

data "alicloud_config_aggregators" "default" {	
	enable_details = true
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
