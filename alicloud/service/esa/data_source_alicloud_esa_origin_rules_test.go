package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEsaOriginRulesDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_esa_origin_rules.default"
	name := fmt.Sprintf("tf-testAcc-EsaOriginRule%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceEsaOriginRulesConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id": "${alicloud_esa_origin_rule.default.site_id}",
			"ids":     []string{"${alicloud_esa_origin_rule.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id": "${alicloud_esa_origin_rule.default.site_id}",
			"ids":     []string{"${alicloud_esa_origin_rule.default.id}_fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":    "${alicloud_esa_origin_rule.default.site_id}",
			"name_regex": "${alicloud_esa_origin_rule.default.rule_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id":    "${alicloud_esa_origin_rule.default.site_id}",
			"name_regex": "${alicloud_esa_origin_rule.default.rule_name}_fake",
		}),
	}

	configIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_origin_rule.default.site_id}",
			"config_id": "${alicloud_esa_origin_rule.default.config_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_origin_rule.default.site_id}",
			"config_id": "${alicloud_esa_origin_rule.default.config_id}000",
		}),
	}

	configTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_origin_rule.default.site_id}",
			"ids":         []string{"${alicloud_esa_origin_rule.default.id}"},
			"config_type": "rule",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_origin_rule.default.site_id}",
			"ids":         []string{"${alicloud_esa_origin_rule.default.id}"},
			"config_type": "global",
		}),
	}

	ruleNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_origin_rule.default.site_id}",
			"rule_name": "${alicloud_esa_origin_rule.default.rule_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_origin_rule.default.site_id}",
			"rule_name": "${alicloud_esa_origin_rule.default.rule_name}_fake",
		}),
	}

	siteVersionConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_origin_rule.default.site_id}",
			"ids":          []string{"${alicloud_esa_origin_rule.default.id}"},
			"site_version": "${alicloud_esa_origin_rule.default.site_version}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_origin_rule.default.site_id}",
			"ids":          []string{"${alicloud_esa_origin_rule.default.id}"},
			"site_version": "1",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_origin_rule.default.site_id}",
			"ids":          []string{"${alicloud_esa_origin_rule.default.id}"},
			"name_regex":   "${alicloud_esa_origin_rule.default.rule_name}",
			"config_id":    "${alicloud_esa_origin_rule.default.config_id}",
			"config_type":  "rule",
			"rule_name":    "${alicloud_esa_origin_rule.default.rule_name}",
			"site_version": "${alicloud_esa_origin_rule.default.site_version}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"rules.0.id":                      acctest.CHECKSET,
			"rules.0.config_id":               acctest.CHECKSET,
			"rules.0.config_type":             acctest.CHECKSET,
			"rules.0.site_version":            acctest.CHECKSET,
			"rules.0.rule_enable":             acctest.CHECKSET,
			"rules.0.rule_name":               acctest.CHECKSET,
			"rules.0.rule":                    acctest.CHECKSET,
			"rules.0.sequence":                acctest.CHECKSET,
			"rules.0.origin_host":             acctest.CHECKSET,
			"rules.0.origin_scheme":           acctest.CHECKSET,
			"rules.0.origin_sni":              acctest.CHECKSET,
			"rules.0.origin_http_port":        acctest.CHECKSET,
			"rules.0.origin_https_port":       acctest.CHECKSET,
			"rules.0.origin_read_timeout":     acctest.CHECKSET,
			"rules.0.dns_record":              acctest.CHECKSET,
			"rules.0.origin_verify":           acctest.CHECKSET,
			"rules.0.origin_mtls":             acctest.CHECKSET,
			"rules.0.follow302_enable":        acctest.CHECKSET,
			"rules.0.follow302_max_tries":     acctest.CHECKSET,
			"rules.0.follow302_target_host":   acctest.CHECKSET,
			"rules.0.follow302_retain_header": acctest.CHECKSET,
			"rules.0.follow302_retain_args":   acctest.CHECKSET,
			"rules.0.range":                   acctest.CHECKSET,
			"rules.0.range_chunk_size":        acctest.CHECKSET,
		}
	}

	var fakeAliCloudEsaOriginRulesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"rules.#": "0",
		}
	}

	var aliCloudEsaOriginRulesInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_esa_origin_rules.default",
		ExistMapFunc: existAliCloudEsaOriginRulesMapFunc,
		FakeMapFunc:  fakeAliCloudEsaOriginRulesMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudEsaOriginRulesInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, configIdConf, configTypeConf, ruleNameConf, siteVersionConf, allConf)
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
