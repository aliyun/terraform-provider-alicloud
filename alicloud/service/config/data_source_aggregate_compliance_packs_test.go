package config_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudConfigAggregateCompliancePacksDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigAggregateCompliancePacksDataSourceName(rand, map[string]string{
			"ids": `[split(":",alicloud_config_aggregate_compliance_pack.default.id)[1]]`,
		}),
		FakeConfig: testAccCheckAlicloudConfigAggregateCompliancePacksDataSourceName(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigAggregateCompliancePacksDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_aggregate_compliance_pack.default.aggregate_compliance_pack_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudConfigAggregateCompliancePacksDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_aggregate_compliance_pack.default.aggregate_compliance_pack_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigAggregateCompliancePacksDataSourceName(rand, map[string]string{
			"ids":    `[split(":",alicloud_config_aggregate_compliance_pack.default.id)[1]]`,
			"status": `"ACTIVE"`,
		}),
		FakeConfig: testAccCheckAlicloudConfigAggregateCompliancePacksDataSourceName(rand, map[string]string{
			"ids":    `[split(":",alicloud_config_aggregate_compliance_pack.default.id)[1]]`,
			"status": `"INACTIVE"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudConfigAggregateCompliancePacksDataSourceName(rand, map[string]string{
			"ids":        `[split(":",alicloud_config_aggregate_compliance_pack.default.id)[1]]`,
			"name_regex": `"${alicloud_config_aggregate_compliance_pack.default.aggregate_compliance_pack_name}"`,
			"status":     `"ACTIVE"`,
		}),
		FakeConfig: testAccCheckAlicloudConfigAggregateCompliancePacksDataSourceName(rand, map[string]string{
			"ids":        `["fake"]`,
			"name_regex": `"${alicloud_config_aggregate_compliance_pack.default.aggregate_compliance_pack_name}_fake"`,
			"status":     `"INACTIVE"`,
		}),
	}
	var existAlicloudConfigAggregateCompliancePacksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"packs.#":                                "1",
			"packs.0.account_id":                     acctest.CHECKSET,
			"packs.0.id":                             acctest.CHECKSET,
			"packs.0.aggregator_compliance_pack_id":  acctest.CHECKSET,
			"packs.0.aggregate_compliance_pack_name": fmt.Sprintf("tf-testAccAggregateCompliancePack-%d", rand),
			"packs.0.compliance_pack_template_id":    "ct-3d20ff4e06a30027f76e",
			"packs.0.config_rules.#":                 "1",
			"packs.0.description":                    fmt.Sprintf("tf-testAccAggregateCompliancePack-%d", rand),
			"packs.0.risk_level":                     "1",
			"packs.0.status":                         "ACTIVE",
		}
	}
	var fakeAlicloudConfigAggregateCompliancePacksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudConfigAggregateCompliancePacksCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_config_aggregate_compliance_packs.default",
		ExistMapFunc: existAlicloudConfigAggregateCompliancePacksDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudConfigAggregateCompliancePacksDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckEnterpriseAccountEnabled(t)
	}
	alicloudConfigAggregateCompliancePacksCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudConfigAggregateCompliancePacksDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccAggregateCompliancePack-%d"
}

data "alicloud_resource_manager_accounts" "default" {
  status = "CreateSuccess"
}
resource "alicloud_config_aggregator" "default" {
  aggregator_accounts {
    account_id   = data.alicloud_resource_manager_accounts.default.accounts.1.account_id
    account_name = data.alicloud_resource_manager_accounts.default.accounts.1.display_name
    account_type = "ResourceDirectory"
  }
  aggregator_name = var.name
  description     = var.name
}

resource "alicloud_config_aggregate_compliance_pack" "default" {
  aggregate_compliance_pack_name  = var.name
  aggregator_id               = alicloud_config_aggregator.default.id
  compliance_pack_template_id = "ct-3d20ff4e06a30027f76e"
  description                 = var.name
  risk_level                  = "1"
  config_rules {
    managed_rule_identifier = "ecs-snapshot-retention-days"
    config_rule_parameters {
      parameter_name  = "days"
      parameter_value = "7"
    }
  }
}

data "alicloud_config_aggregate_compliance_packs" "default" {
	aggregator_id = alicloud_config_aggregate_compliance_pack.default.aggregator_id
	enable_details = true
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
