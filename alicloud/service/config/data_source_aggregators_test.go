package config_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudConfigAggregatorsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigAggregatorsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_config_aggregator.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudConfigAggregatorsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_config_aggregator.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigAggregatorsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_aggregator.default.aggregator_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudConfigAggregatorsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_aggregator.default.aggregator_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigAggregatorsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_config_aggregator.default.id}"]`,
			"status": `"Normal"`,
		}),
		FakeConfig: testAccCheckAlicloudConfigAggregatorsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_config_aggregator.default.id}"]`,
			"status": `"Deleting"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigAggregatorsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_config_aggregator.default.id}"]`,
			"name_regex": `"${alicloud_config_aggregator.default.aggregator_name}"`,
			"status":     `"Normal"`,
		}),
		FakeConfig: testAccCheckAlicloudConfigAggregatorsDataSourceName(rand, map[string]string{
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
			"aggregators.0.aggregator_name":       acctest.CHECKSET,
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
	var alicloudConfigAggregatorsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_config_aggregators.default",
		ExistMapFunc: existAlicloudConfigAggregatorsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudConfigAggregatorsDataSourceNameMapFunc,
	}

	PreCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckEnterpriseAccountEnabled(t)
	}

	alicloudConfigAggregatorsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, PreCheck, idsConf, nameRegexConf, statusConf, allConf)
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
