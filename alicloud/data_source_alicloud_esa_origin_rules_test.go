package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAliCloudEsaOriginRulesDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_esa_origin_rules.default"
	name := fmt.Sprintf("tf-testAcc-EsaOriginRule%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEsaOriginRulesConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id": "${alicloud_esa_origin_rule.default.site_id}",
			"ids":     []string{"${alicloud_esa_origin_rule.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id": "${alicloud_esa_origin_rule.default.site_id}",
			"ids":     []string{"${alicloud_esa_origin_rule.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":    "${alicloud_esa_origin_rule.default.site_id}",
			"name_regex": "${alicloud_esa_origin_rule.default.rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":    "${alicloud_esa_origin_rule.default.site_id}",
			"name_regex": "${alicloud_esa_origin_rule.default.rule_name}_fake",
		}),
	}

	configIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_origin_rule.default.site_id}",
			"config_id": "${alicloud_esa_origin_rule.default.config_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_origin_rule.default.site_id}",
			"config_id": "${alicloud_esa_origin_rule.default.config_id}000",
		}),
	}

	configTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_origin_rule.default.site_id}",
			"ids":         []string{"${alicloud_esa_origin_rule.default.id}"},
			"config_type": "rule",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_origin_rule.default.site_id}",
			"ids":         []string{"${alicloud_esa_origin_rule.default.id}"},
			"config_type": "global",
		}),
	}

	ruleNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_origin_rule.default.site_id}",
			"rule_name": "${alicloud_esa_origin_rule.default.rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_origin_rule.default.site_id}",
			"rule_name": "${alicloud_esa_origin_rule.default.rule_name}_fake",
		}),
	}

	siteVersionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_origin_rule.default.site_id}",
			"ids":          []string{"${alicloud_esa_origin_rule.default.id}"},
			"site_version": "${alicloud_esa_origin_rule.default.site_version}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_origin_rule.default.site_id}",
			"ids":          []string{"${alicloud_esa_origin_rule.default.id}"},
			"site_version": "1",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_origin_rule.default.site_id}",
			"ids":          []string{"${alicloud_esa_origin_rule.default.id}"},
			"name_regex":   "${alicloud_esa_origin_rule.default.rule_name}",
			"config_id":    "${alicloud_esa_origin_rule.default.config_id}",
			"config_type":  "rule",
			"rule_name":    "${alicloud_esa_origin_rule.default.rule_name}",
			"site_version": "${alicloud_esa_origin_rule.default.site_version}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_origin_rule.default.site_id}",
			"ids":          []string{"${alicloud_esa_origin_rule.default.id}_fake"},
			"name_regex":   "${alicloud_esa_origin_rule.default.rule_name}_fake",
			"config_id":    "${alicloud_esa_origin_rule.default.config_id}000",
			"config_type":  "global",
			"rule_name":    "${alicloud_esa_origin_rule.default.rule_name}_fake",
			"site_version": "1",
		}),
	}

	var existAliCloudEsaOriginRulesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"rules.#":                         "1",
			"rules.0.id":                      CHECKSET,
			"rules.0.config_id":               CHECKSET,
			"rules.0.config_type":             CHECKSET,
			"rules.0.site_version":            CHECKSET,
			"rules.0.rule_enable":             CHECKSET,
			"rules.0.rule_name":               CHECKSET,
			"rules.0.rule":                    CHECKSET,
			"rules.0.sequence":                CHECKSET,
			"rules.0.origin_host":             CHECKSET,
			"rules.0.origin_scheme":           CHECKSET,
			"rules.0.origin_sni":              CHECKSET,
			"rules.0.origin_http_port":        CHECKSET,
			"rules.0.origin_https_port":       CHECKSET,
			"rules.0.origin_read_timeout":     CHECKSET,
			"rules.0.dns_record":              CHECKSET,
			"rules.0.origin_verify":           CHECKSET,
			"rules.0.origin_mtls":             CHECKSET,
			"rules.0.follow302_enable":        CHECKSET,
			"rules.0.follow302_max_tries":     CHECKSET,
			"rules.0.follow302_target_host":   CHECKSET,
			"rules.0.follow302_retain_header": CHECKSET,
			"rules.0.follow302_retain_args":   CHECKSET,
			"rules.0.range":                   CHECKSET,
			"rules.0.range_chunk_size":        CHECKSET,
		}
	}

	var fakeAliCloudEsaOriginRulesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"rules.#": "0",
		}
	}

	var aliCloudEsaOriginRulesInfo = dataSourceAttr{
		resourceId:   "data.alicloud_esa_origin_rules.default",
		existMapFunc: existAliCloudEsaOriginRulesMapFunc,
		fakeMapFunc:  fakeAliCloudEsaOriginRulesMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudEsaOriginRulesInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, configIdConf, configTypeConf, ruleNameConf, siteVersionConf, allConf)
}

func dataSourceEsaOriginRulesConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_origin_rule" "default" {
  site_id                 = data.alicloud_esa_sites.default.sites.0.id
  site_version            = 0
  rule_enable             = "on"
  rule_name               = var.name
  rule                    = "true"
  sequence                = 1
  origin_host             = "origin.example.com"
  origin_scheme           = "http"
  origin_sni              = "origin.example.com"
  origin_https_port       = "443"
  origin_http_port        = "8080"
  origin_read_timeout     = "30"
  dns_record              = "test.example.com"
  origin_verify           = "on"
  origin_mtls             = "on"
  follow302_enable        = "on"
  follow302_max_tries     = "3"
  follow302_target_host   = "redirect.example.com"
  follow302_retain_header = "on"
  follow302_retain_args   = "on"
  range                   = "on"
  range_chunk_size        = "1MB"
}
`, name)
}
