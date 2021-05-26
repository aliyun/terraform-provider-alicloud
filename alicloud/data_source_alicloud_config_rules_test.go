package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudConfigRulesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_config_rules.example"

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_rule.example.rule_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_rule.example.rule_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_config_rule.example.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_config_rule.example.id}_fake"]`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_config_rule.example.id}"]`,
			"status": `"ACTIVE"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_config_rule.example.id}"]`,
			"status": `"INACTIVE"`,
		}),
	}

	riskLevelConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_rule.example.rule_name}"`,
			"risk_level": `1`,
		}),
		fakeConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_rule.example.rule_name}"`,
			"risk_level": `2`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_rule.example.rule_name}"`,
			"ids":        `["${alicloud_config_rule.example.id}"]`,
			"risk_level": `1`,
			"status":     `"ACTIVE"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_config_rule.example.rule_name}_fake"`,
			"ids":        `["${alicloud_config_rule.example.id}_fake"]`,
			"risk_level": `2`,
			"status":     `"INACTIVE"`,
		}),
	}

	var existConfigRulesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#":                                    "1",
			"names.#":                                    "1",
			"ids.#":                                      "1",
			"rules.0.account_id":                         CHECKSET,
			"rules.0.config_rule_arn":                    CHECKSET,
			"rules.0.id":                                 CHECKSET,
			"rules.0.config_rule_id":                     CHECKSET,
			"rules.0.config_rule_state":                  CHECKSET,
			"rules.0.status":                             "ACTIVE",
			"rules.0.compliance.#":                       CHECKSET,
			"rules.0.description":                        fmt.Sprintf("tf-testAccConfigRule%d", rand),
			"rules.0.input_parameters.%":                 "1",
			"rules.0.modified_timestamp":                 "",
			"rules.0.risk_level":                         "1",
			"rules.0.rule_name":                          fmt.Sprintf("tf-testAccConfigRule%d", rand),
			"rules.0.event_source":                       "aliyun.config",
			"rules.0.source_maximum_execution_frequency": "",
			"rules.0.scope_compliance_resource_types.#":  "1",
			"rules.0.source_detail_message_type":         "ConfigurationItemChangeNotification",
			"rules.0.source_identifier":                  "ecs-instances-in-vpc",
			"rules.0.source_owner":                       "ALIYUN",
			"rules.0.tag_key_scope":                      "tfTest",
			"rules.0.tag_value_scope":                    "tfTest 123",
			"rules.0.resource_types_scope.#":             "1",
			"rules.0.resource_group_ids_scope":           CHECKSET,
			"rules.0.region_ids_scope":                   "cn-hangzhou",
			"rules.0.maximum_execution_frequency":        "",
			"rules.0.exclude_resource_ids_scope":         CHECKSET,
			"rules.0.config_rule_trigger_types":          "ConfigurationItemChangeNotification",
			"rules.0.compliance_pack_id":                 "",
		}
	}

	var fakeConfigRulesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#": "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var rolesRecordsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existConfigRulesRecordsMapFunc,
		fakeMapFunc:  fakeConfigRulesRecordsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.CloudConfigSupportedRegions)
	}

	rolesRecordsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, statusConf, riskLevelConf, allConf)

}

func testAccCheckAlicloudConfigRulesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccConfigRule%d"
}

data "alicloud_instances" "default"{}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_config_rule" "example" {
  rule_name                       = var.name
  description                     = var.name
  source_identifier               = "ecs-instances-in-vpc"
  source_owner                    = "ALIYUN"
  resource_types_scope 			  = ["ACS::ECS::Instance"]
  risk_level                      = 1
  config_rule_trigger_types       = "ConfigurationItemChangeNotification"
  tag_key_scope 				  = "tfTest"
  tag_value_scope 				  = "tfTest 123"
  resource_group_ids_scope 		  = data.alicloud_resource_manager_resource_groups.default.ids.0
  exclude_resource_ids_scope      = data.alicloud_instances.default.instances[0].id
  region_ids_scope 				  = "cn-hangzhou"
  input_parameters  = {
		vpcIds= data.alicloud_instances.default.instances[0].vpc_id
  }
}

data "alicloud_config_rules" "example"{
 enable_details = true
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
