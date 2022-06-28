package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudConfigCompliancePacksDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigCompliancePacksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_config_compliance_pack.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudConfigCompliancePacksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_config_compliance_pack.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigCompliancePacksDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_compliance_pack.default.compliance_pack_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigCompliancePacksDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_config_compliance_pack.default.compliance_pack_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigCompliancePacksDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_config_compliance_pack.default.id}"]`,
			"status": `"ACTIVE"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigCompliancePacksDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_config_compliance_pack.default.id}"]`,
			"status": `"INACTIVE"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigCompliancePacksDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_config_compliance_pack.default.id}"]`,
			"name_regex": `"${alicloud_config_compliance_pack.default.compliance_pack_name}"`,
			"status":     `"ACTIVE"`,
		}),
		fakeConfig: testAccCheckAlicloudConfigCompliancePacksDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_config_compliance_pack.default.id}_fake"]`,
			"name_regex": `"${alicloud_config_compliance_pack.default.compliance_pack_name}_fake"`,
			"status":     `"INACTIVE"`,
		}),
	}
	var existAlicloudConfigCompliancePacksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"packs.#":                             "1",
			"packs.0.account_id":                  CHECKSET,
			"packs.0.id":                          CHECKSET,
			"packs.0.compliance_pack_id":          CHECKSET,
			"packs.0.compliance_pack_name":        fmt.Sprintf("tf-testAccCompliancePack-%d", rand),
			"packs.0.compliance_pack_template_id": "ct-3d20ff4e06a30027f76e",
			"packs.0.config_rules.#":              "1",
			"packs.0.description":                 fmt.Sprintf("tf-testAccCompliancePack-%d", rand),
			"packs.0.risk_level":                  "1",
			"packs.0.status":                      "ACTIVE",
		}
	}
	var fakeAlicloudConfigCompliancePacksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudConfigCompliancePacksCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_config_compliance_packs.default",
		existMapFunc: existAlicloudConfigCompliancePacksDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudConfigCompliancePacksDataSourceNameMapFunc,
	}
	alicloudConfigCompliancePacksCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudConfigCompliancePacksDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccCompliancePack-%d"
}

resource "alicloud_config_compliance_pack" "default" {
  compliance_pack_name        = var.name
  compliance_pack_template_id = "ct-3d20ff4e06a30027f76e"
  description                 = var.name
  risk_level                  = "1"
  config_rules {
    managed_rule_identifier = "ecs-instance-expired-check"
    config_rule_parameters {
      parameter_name  = "days"
      parameter_value = "60"
    }
  }
}

data "alicloud_config_compliance_packs" "default" {	
	enable_details = true
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
