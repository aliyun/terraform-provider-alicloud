package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudEsaWafRuleSetsDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_esa_waf_rulesets.default"
	name := fmt.Sprintf("tf-testAcc-EsaWafRuleSet%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEsaWafRuleSetsConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_waf_rule.default.site_id}",
			"phase":        "http_custom",
			"site_version": "0",
			"ids":          []string{"${alicloud_esa_waf_ruleset.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_waf_rule.default.site_id}",
			"phase":        "http_custom",
			"site_version": "0",
			"ids":          []string{"${alicloud_esa_waf_ruleset.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_waf_rule.default.site_id}",
			"phase":        "http_custom",
			"site_version": "0",
			"name_regex":   "${alicloud_esa_waf_ruleset.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_waf_rule.default.site_id}",
			"phase":        "http_custom",
			"site_version": "0",
			"name_regex":   "${alicloud_esa_waf_ruleset.default.name}_fake",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_waf_rule.default.site_id}",
			"phase":        "http_custom",
			"site_version": "0",
			"ids":          []string{"${alicloud_esa_waf_ruleset.default.id}"},
			"status":       "on",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_waf_rule.default.site_id}",
			"phase":        "http_custom",
			"site_version": "0",
			"ids":          []string{"${alicloud_esa_waf_ruleset.default.id}_fake"},
			"status":       "off",
		}),
	}

	anyQueryArgsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_waf_rule.default.site_id}",
			"phase":        "http_custom",
			"site_version": "0",
			"query_args": []map[string]interface{}{
				{
					"any_like": "${alicloud_esa_waf_rule.default.ruleset_id}",
					"order_by": "id",
					"desc":     "true",
				},
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_waf_rule.default.site_id}",
			"phase":        "http_custom",
			"site_version": "0",
			"query_args": []map[string]interface{}{
				{
					"any_like": "${alicloud_esa_waf_rule.default.ruleset_id}_fake",
					"order_by": "id",
					"desc":     "false",
				},
			},
		}),
	}

	nameQueryArgsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_waf_rule.default.site_id}",
			"phase":        "http_custom",
			"site_version": "0",
			"query_args": []map[string]interface{}{
				{
					"name_like": "${alicloud_esa_waf_ruleset.default.name}",
					"order_by":  "name",
					"desc":      "true",
				},
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_waf_rule.default.site_id}",
			"phase":        "http_custom",
			"site_version": "0",
			"query_args": []map[string]interface{}{
				{
					"name_like": "${alicloud_esa_waf_ruleset.default.name}_fake",
					"order_by":  "name",
					"desc":      "false",
				},
			},
		}),
	}

	allQueryArgsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_waf_rule.default.site_id}",
			"phase":        "http_custom",
			"site_version": "0",
			"query_args": []map[string]interface{}{
				{
					"any_like":  "${alicloud_esa_waf_rule.default.ruleset_id}",
					"name_like": "${alicloud_esa_waf_ruleset.default.name}",
					"order_by":  "id",
					"desc":      "true",
				},
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_waf_rule.default.site_id}",
			"phase":        "http_custom",
			"site_version": "0",
			"query_args": []map[string]interface{}{
				{
					"any_like":  "${alicloud_esa_waf_rule.default.ruleset_id}_fake",
					"name_like": "${alicloud_esa_waf_ruleset.default.name}_fake",
					"order_by":  "id",
					"desc":      "false",
				},
			},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_waf_rule.default.site_id}",
			"phase":        "http_custom",
			"site_version": "0",
			"ids":          []string{"${alicloud_esa_waf_ruleset.default.id}"},
			"name_regex":   "${alicloud_esa_waf_ruleset.default.name}",
			"status":       "on",
			"query_args": []map[string]interface{}{
				{
					"any_like":  "${alicloud_esa_waf_rule.default.ruleset_id}",
					"name_like": "${alicloud_esa_waf_ruleset.default.name}",
					"order_by":  "id",
					"desc":      "true",
				},
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_waf_rule.default.site_id}",
			"phase":        "http_custom",
			"site_version": "0",
			"ids":          []string{"${alicloud_esa_waf_ruleset.default.id}_fake"},
			"name_regex":   "${alicloud_esa_waf_ruleset.default.name}_fake",
			"status":       "off",
			"query_args": []map[string]interface{}{
				{
					"any_like":  "${alicloud_esa_waf_rule.default.ruleset_id}_fake",
					"name_like": "${alicloud_esa_waf_ruleset.default.name}_fake",
					"order_by":  "id",
					"desc":      "false",
				},
			},
		}),
	}

	var existAliCloudEsaWafRuleSetsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":              "1",
			"names.#":            "1",
			"sets.#":             "1",
			"sets.0.id":          CHECKSET,
			"sets.0.ruleset_id":  CHECKSET,
			"sets.0.phase":       CHECKSET,
			"sets.0.name":        CHECKSET,
			"sets.0.status":      CHECKSET,
			"sets.0.update_time": CHECKSET,
			"sets.0.types.#":     CHECKSET,
			//"sets.0.target":      CHECKSET,
			//"sets.0.fields.#":    CHECKSET,
		}
	}

	var fakeAliCloudEsaWafRuleSetsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"sets.#":  "0",
		}
	}

	var aliCloudEsaWafRuleSetsInfo = dataSourceAttr{
		resourceId:   "data.alicloud_esa_waf_rulesets.default",
		existMapFunc: existAliCloudEsaWafRuleSetsMapFunc,
		fakeMapFunc:  fakeAliCloudEsaWafRuleSetsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudEsaWafRuleSetsInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, anyQueryArgsConf, nameQueryArgsConf, allQueryArgsConf, allConf)
}

func dataSourceEsaWafRuleSetsConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_waf_ruleset" "default" {
  site_id      = data.alicloud_esa_sites.default.sites.0.site_id
  phase        = "http_custom"
  site_version = "0"
  name         = var.name
}

resource "alicloud_esa_waf_rule" "default" {
  site_id      = data.alicloud_esa_sites.default.sites.0.site_id
  ruleset_id   = alicloud_esa_waf_ruleset.default.ruleset_id
  phase        = "http_custom"
  site_version = "0"
  config {
    status     = "on"
    action     = "deny"
    expression = "(http.host in {\"123.example.top\"})"
    actions {
      response {
        id   = "0"
        code = "403"
      }
    }
  }
}
`, name)
}
