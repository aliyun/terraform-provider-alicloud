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

	configRuleStatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_config_rule.example.rule_name}"`,
			"config_rule_state": `"ACTIVE"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_config_rule.example.rule_name}"`,
			"config_rule_state": `"INACTIVE"`,
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
			"name_regex":        `"${alicloud_config_rule.example.rule_name}"`,
			"ids":               `["${alicloud_config_rule.example.id}"]`,
			"config_rule_state": `"ACTIVE"`,
			"risk_level":        `1`,
		}),
		fakeConfig: testAccCheckAlicloudConfigRulesSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_config_rule.example.rule_name}_fake"`,
			"ids":               `["${alicloud_config_rule.example.id}_fake"]`,
			"config_rule_state": `"INACTIVE"`,
			"risk_level":        `2`,
		}),
	}

	var existConfigRulesRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#":                               "1",
			"names.#":                               "1",
			"ids.#":                                 "1",
			"rules.0.account_id":                    CHECKSET,
			"rules.0.config_rule_arn":               CHECKSET,
			"rules.0.id":                            CHECKSET,
			"rules.0.config_rule_id":                CHECKSET,
			"rules.0.config_rule_state":             "ACTIVE",
			"rules.0.create_timestamp":              CHECKSET,
			"rules.0.description":                   "",
			"rules.0.input_parameters.%":            "0",
			"rules.0.modified_timestamp":            CHECKSET,
			"rules.0.risk_level":                    "1",
			"rules.0.rule_name":                     fmt.Sprintf("tf-testAccConfigRule%d", rand),
			"rules.0.source_details.#":              "1",
			"rules.0.source_details.0.event_source": "aliyun.config",
			"rules.0.source_details.0.maximum_execution_frequency": "",
			"rules.0.source_details.0.message_type":                "ConfigurationItemChangeNotification",
			"rules.0.source_identifier":                            "ecs-instances-in-vpc",
			"rules.0.source_owner":                                 "ALIYUN",
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

	rolesRecordsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, configRuleStatusConf, riskLevelConf, allConf)

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
resource "alicloud_config_rule" "example" {
  rule_name                       = "${var.name}"
  source_identifier               = "ecs-instances-in-vpc"
  source_owner                    = "ALIYUN"
  scope_compliance_resource_types = ["ACS::ECS::Instance"]
  risk_level                         = 1
  source_detail_message_type         = "ConfigurationItemChangeNotification"
}
data "alicloud_config_rules" "example"{
 enable_details = true
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
