package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudEsaErrorPagesRedirectRulesDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_esa_error_pages_redirect_rules.default"
	name := fmt.Sprintf("tf-testAcc-EsaErrorPagesRedirectRule%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEsaErrorPagesRedirectRulesConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id": "${alicloud_esa_error_pages_redirect_rule.default.site_id}",
			"ids":     []string{"${alicloud_esa_error_pages_redirect_rule.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id": "${alicloud_esa_error_pages_redirect_rule.default.site_id}",
			"ids":     []string{"${alicloud_esa_error_pages_redirect_rule.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":    "${alicloud_esa_error_pages_redirect_rule.default.site_id}",
			"name_regex": "${alicloud_esa_error_pages_redirect_rule.default.rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":    "${alicloud_esa_error_pages_redirect_rule.default.site_id}",
			"name_regex": "${alicloud_esa_error_pages_redirect_rule.default.rule_name}_fake",
		}),
	}

	ruleNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_error_pages_redirect_rule.default.site_id}",
			"rule_name": "${alicloud_esa_error_pages_redirect_rule.default.rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_error_pages_redirect_rule.default.site_id}",
			"rule_name": "${alicloud_esa_error_pages_redirect_rule.default.rule_name}_fake",
		}),
	}

	configTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_error_pages_redirect_rule.default.site_id}",
			"ids":         []string{"${alicloud_esa_error_pages_redirect_rule.default.id}"},
			"config_type": "rule",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_error_pages_redirect_rule.default.site_id}",
			"ids":         []string{"${alicloud_esa_error_pages_redirect_rule.default.id}"},
			"config_type": "global",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_error_pages_redirect_rule.default.site_id}",
			"ids":         []string{"${alicloud_esa_error_pages_redirect_rule.default.id}"},
			"name_regex":  "${alicloud_esa_error_pages_redirect_rule.default.rule_name}",
			"config_type": "rule",
			"rule_name":   "${alicloud_esa_error_pages_redirect_rule.default.rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_error_pages_redirect_rule.default.site_id}",
			"ids":         []string{"${alicloud_esa_error_pages_redirect_rule.default.id}_fake"},
			"name_regex":  "${alicloud_esa_error_pages_redirect_rule.default.rule_name}_fake",
			"config_type": "global",
			"rule_name":   "${alicloud_esa_error_pages_redirect_rule.default.rule_name}_fake",
		}),
	}

	var existAliCloudEsaErrorPagesRedirectRulesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"rules.#":                        "1",
			"rules.0.id":                     CHECKSET,
			"rules.0.config_id":              CHECKSET,
			"rules.0.config_type":            CHECKSET,
			"rules.0.site_version":           CHECKSET,
			"rules.0.rule_enable":            CHECKSET,
			"rules.0.rule_name":              CHECKSET,
			"rules.0.rule":                   CHECKSET,
			"rules.0.sequence":               CHECKSET,
			"rules.0.error_pages_redirect.#": CHECKSET,
		}
	}

	var fakeAliCloudEsaErrorPagesRedirectRulesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"rules.#": "0",
		}
	}

	var aliCloudEsaErrorPagesRedirectRulesInfo = dataSourceAttr{
		resourceId:   "data.alicloud_esa_error_pages_redirect_rules.default",
		existMapFunc: existAliCloudEsaErrorPagesRedirectRulesMapFunc,
		fakeMapFunc:  fakeAliCloudEsaErrorPagesRedirectRulesMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudEsaErrorPagesRedirectRulesInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, ruleNameConf, configTypeConf, allConf)
}

func dataSourceEsaErrorPagesRedirectRulesConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_error_pages_redirect_rule" "default" {
  site_id      = data.alicloud_esa_sites.default.sites.0.id
  rule_enable  = "off"
  rule         = "(http.host eq \"video.example.com\")"
  sequence     = 1
  site_version = 0
  rule_name    = var.name
  error_pages_redirect {
    target_url  = "https://example.com/foo/bar"
    status_code = "500"
  }
  error_pages_redirect {
    target_url  = "https://example.com/foo"
    status_code = "400"
  }
}
`, name)
}
