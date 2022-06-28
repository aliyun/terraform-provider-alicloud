package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSddpRuleDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sddp_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sddp_rule.default.id}_fake"]`,
		}),
	}
	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"name": `"${alicloud_sddp_rule.default.rule_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"name": `"${alicloud_sddp_rule.default.rule_name}_fake"`,
		}),
	}
	nameregexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_sddp_rule.default.rule_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_sddp_rule.default.rule_name}_fake"`,
		}),
	}
	risklevelidConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"risk_level_id": `"${alicloud_sddp_rule.default.risk_level_id}"`,
			"ids":           `["${alicloud_sddp_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"risk_level_id": `"0"`,
			"ids":           `["${alicloud_sddp_rule.default.id}"]`,
		}),
	}
	customtypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"custom_type": `"${alicloud_sddp_rule.default.custom_type}"`,
			"ids":         `["${alicloud_sddp_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"custom_type": `"2"`,
			"ids":         `["${alicloud_sddp_rule.default.id}"]`,
		}),
	}
	productidConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"product_id": `"${alicloud_sddp_rule.default.product_id}"`,
			"ids":        `["${alicloud_sddp_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"product_id": `"2"`,
			"ids":        `["${alicloud_sddp_rule.default.id}"]`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"status": `"${alicloud_sddp_rule.default.status}"`,
			"ids":    `["${alicloud_sddp_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"status": `"0"`,
			"ids":    `["${alicloud_sddp_rule.default.id}"]`,
		}),
	}
	ruletypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"rule_type": `"${alicloud_sddp_rule.default.rule_type}"`,
			"ids":       `["${alicloud_sddp_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"rule_type": `"3"`,
			"ids":       `["${alicloud_sddp_rule.default.id}"]`,
		}),
	}

	categoryConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"category": `"${alicloud_sddp_rule.default.category}"`,
			"ids":      `["${alicloud_sddp_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"category": `"2"`,
			"ids":      `["${alicloud_sddp_rule.default.id}"]`,
		}),
	}

	warnevelConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"warn_level": `"${alicloud_sddp_rule.default.warn_level}"`,
			"ids":        `["${alicloud_sddp_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"warn_level": `"1"`,
			"ids":        `["${alicloud_sddp_rule.default.id}"]`,
		}),
	}

	contentCategoryConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"content_category": `"${alicloud_sddp_rule.default.content_category}"`,
			"ids":              `["${alicloud_sddp_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"content_category": `"1"`,
			"ids":              `["${alicloud_sddp_rule.default.id}"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_sddp_rule.default.id}"]`,
			"name":          `"${alicloud_sddp_rule.default.rule_name}"`,
			"name_regex":    `"${alicloud_sddp_rule.default.rule_name}"`,
			"risk_level_id": `"${alicloud_sddp_rule.default.risk_level_id}"`,
			"custom_type":   `"${alicloud_sddp_rule.default.custom_type}"`,
			"product_id":    `"${alicloud_sddp_rule.default.product_id}"`,
			"status":        `"${alicloud_sddp_rule.default.status}"`,
			"rule_type":     `"${alicloud_sddp_rule.default.rule_type}"`,
			"category":      `"${alicloud_sddp_rule.default.category}"`,
			"warn_level":    `"${alicloud_sddp_rule.default.warn_level}"`,
		}),
		fakeConfig: testAccCheckAlicloudSddpRuleDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_sddp_rule.default.id}_fake"]`,
			"name":          `"${alicloud_sddp_rule.default.rule_name}_fake"`,
			"name_regex":    `"${alicloud_sddp_rule.default.rule_name}_fake"`,
			"risk_level_id": `"0"`,
			"custom_type":   `"2"`,
			"product_id":    `"2"`,
			"status":        `"0"`,
			"rule_type":     `"3"`,
			"category":      `"2"`,
			"warn_level":    `"1"`,
		}),
	}
	var existAlicloudSaeNamespaceDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "1",
			"names.#":               "1",
			"rules.#":               "1",
			"rules.0.risk_level_id": "4",
			"rules.0.status":        "1",
			"rules.0.warn_level":    "3",
			"rules.0.name":          fmt.Sprintf("tf-testAccSddpRule-%d", rand),
			"rules.0.category":      "0",
			"rules.0.product_id":    "5",
		}
	}
	var fakeAlicloudSaeNamespaceDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudSaeNamespaceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_sddp_rules.default",
		existMapFunc: existAlicloudSaeNamespaceDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudSaeNamespaceDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.SddpSupportRegions)
	}
	alicloudSaeNamespaceCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameConf, nameregexConf, risklevelidConf, customtypeConf, productidConf, productidConf, statusConf, ruletypeConf, categoryConf, warnevelConf, contentCategoryConf, allConf)

}
func testAccCheckAlicloudSddpRuleDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccSddpRule-%d"
}

resource "alicloud_sddp_rule" "default" {
  category=  "0"
  content=   var.name
  rule_name= var.name
  risk_level_id = "4"
  warn_level = "3"
  product_code = "RDS"
  product_id = "5"
  
}

data "alicloud_sddp_rules" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
