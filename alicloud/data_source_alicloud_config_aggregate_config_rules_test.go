package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudConfigAggregateConfigRulesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigAggregateConfigRulesDataSourceName(rand, map[string]string{
			"ids": `[split(":",alicloud_config_aggregate_config_rule.default.id)[1]]`,
		}),
		fakeConfig: testAccCheckAlicloudConfigAggregateConfigRulesDataSourceName(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}
	aggregateConfigRuleNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigAggregateConfigRulesDataSourceName(rand, map[string]string{
			"ids":                        `[split(":",alicloud_config_aggregate_config_rule.default.id)[1]]`,
			"aggregate_config_rule_name": `"${alicloud_config_aggregate_config_rule.default.aggregate_config_rule_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigAggregateConfigRulesDataSourceName(rand, map[string]string{
			"ids":                        `[split(":",alicloud_config_aggregate_config_rule.default.id)[1]]`,
			"aggregate_config_rule_name": `"${alicloud_config_aggregate_config_rule.default.aggregate_config_rule_name}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigAggregateConfigRulesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_aggregate_config_rule.default.aggregate_config_rule_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigAggregateConfigRulesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_aggregate_config_rule.default.aggregate_config_rule_name}_fake"`,
		}),
	}
	riskLevelConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigAggregateConfigRulesDataSourceName(rand, map[string]string{
			"ids":        `[split(":",alicloud_config_aggregate_config_rule.default.id)[1]]`,
			"risk_level": `"1"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigAggregateConfigRulesDataSourceName(rand, map[string]string{
			"ids":        `[split(":",alicloud_config_aggregate_config_rule.default.id)[1]]`,
			"risk_level": `"2"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigAggregateConfigRulesDataSourceName(rand, map[string]string{
			"ids":    `[split(":",alicloud_config_aggregate_config_rule.default.id)[1]]`,
			"status": `"ACTIVE"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigAggregateConfigRulesDataSourceName(rand, map[string]string{
			"ids":    `[split(":",alicloud_config_aggregate_config_rule.default.id)[1]]`,
			"status": `"INACTIVE"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigAggregateConfigRulesDataSourceName(rand, map[string]string{
			"ids":                        `[split(":",alicloud_config_aggregate_config_rule.default.id)[1]]`,
			"aggregate_config_rule_name": `"${alicloud_config_aggregate_config_rule.default.aggregate_config_rule_name}"`,
			"status":                     `"ACTIVE"`,
			"name_regex":                 `"${alicloud_config_aggregate_config_rule.default.aggregate_config_rule_name}"`,
			"risk_level":                 `"1"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigAggregateConfigRulesDataSourceName(rand, map[string]string{
			"ids":                        `["fake"]`,
			"status":                     `"INACTIVE"`,
			"aggregate_config_rule_name": `"${alicloud_config_aggregate_config_rule.default.aggregate_config_rule_name}_fake"`,
			"name_regex":                 `"${alicloud_config_aggregate_config_rule.default.aggregate_config_rule_name}_fake"`,
			"risk_level":                 `"2"`,
		}),
	}
	var existAlicloudConfigAggregateConfigRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"rules.#":                             "1",
			"rules.0.account_id":                  CHECKSET,
			"rules.0.compliance.#":                "1",
			"rules.0.compliance_pack_id":          "",
			"rules.0.aggregator_id":               CHECKSET,
			"rules.0.config_rule_arn":             CHECKSET,
			"rules.0.id":                          CHECKSET,
			"rules.0.config_rule_id":              CHECKSET,
			"rules.0.status":                      "ACTIVE",
			"rules.0.aggregate_config_rule_name":  CHECKSET,
			"rules.0.config_rule_trigger_types":   `ConfigurationItemChangeNotification`,
			"rules.0.description":                 fmt.Sprintf("tf-testAccAggregateConfigRule-%d", rand),
			"rules.0.source_identifier":           "ecs-cpu-min-count-limit",
			"rules.0.source_owner":                "ALIYUN",
			"rules.0.event_source":                CHECKSET,
			"rules.0.region_ids_scope":            "cn-hangzhou",
			"rules.0.risk_level":                  "1",
			"rules.0.exclude_resource_ids_scope":  CHECKSET,
			"rules.0.resource_types_scope.#":      "1",
			"rules.0.maximum_execution_frequency": "",
			"rules.0.tag_key_scope":               "tFTest",
			"rules.0.tag_value_scope":             "forTF 123",
			"rules.0.modified_timestamp":          CHECKSET,
			"rules.0.resource_group_ids_scope":    CHECKSET,
		}
	}
	var fakeAlicloudConfigAggregateConfigRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudConfigAggregateConfigRulesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_config_aggregate_config_rules.default",
		existMapFunc: existAlicloudConfigAggregateConfigRulesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudConfigAggregateConfigRulesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckEnterpriseAccountEnabled(t)
	}

	alicloudConfigAggregateConfigRulesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, aggregateConfigRuleNameConf, nameRegexConf, riskLevelConf, statusConf, allConf)
}
func testAccCheckAlicloudConfigAggregateConfigRulesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccAggregateConfigRule-%d"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_instances" "default" {}

data "alicloud_config_aggregators" "default" {}

resource "alicloud_config_aggregate_config_rule" "default" {
  aggregator_id              = data.alicloud_config_aggregators.default.ids.0
  aggregate_config_rule_name = var.name
  source_owner               = "ALIYUN"
  source_identifier    		= "ecs-cpu-min-count-limit"
  config_rule_trigger_types = "ConfigurationItemChangeNotification"
  resource_types_scope      = ["ACS::ECS::Instance"]
  risk_level                = 1
  description                = var.name
  exclude_resource_ids_scope = data.alicloud_instances.default.ids.0
  input_parameters = {
    cpuCount = "4",
  }
  region_ids_scope         = "cn-hangzhou"
  resource_group_ids_scope = data.alicloud_resource_manager_resource_groups.default.ids.0
  tag_key_scope            = "tFTest"
  tag_value_scope          = "forTF 123"
}

data "alicloud_config_aggregate_config_rules" "default" {	
  	aggregator_id  = alicloud_config_aggregate_config_rule.default.aggregator_id
	enable_details = true
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
