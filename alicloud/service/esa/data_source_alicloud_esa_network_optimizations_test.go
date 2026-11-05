package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEsaNetworkOptimizationsDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_esa_network_optimizations.default"
	name := fmt.Sprintf("tf-testAcc-EsaNetworkOptimization%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceEsaNetworkOptimizationsConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id": "${alicloud_esa_network_optimization.default.site_id}",
			"ids":     []string{"${alicloud_esa_network_optimization.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id": "${alicloud_esa_network_optimization.default.site_id}",
			"ids":     []string{"${alicloud_esa_network_optimization.default.id}_fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":    "${alicloud_esa_network_optimization.default.site_id}",
			"name_regex": "${alicloud_esa_network_optimization.default.rule_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id":    "${alicloud_esa_network_optimization.default.site_id}",
			"name_regex": "${alicloud_esa_network_optimization.default.rule_name}_fake",
		}),
	}

	configIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_network_optimization.default.site_id}",
			"config_id": "${alicloud_esa_network_optimization.default.config_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_network_optimization.default.site_id}",
			"config_id": "${alicloud_esa_network_optimization.default.config_id}000",
		}),
	}

	configTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_network_optimization.default.site_id}",
			"ids":         []string{"${alicloud_esa_network_optimization.default.id}"},
			"config_type": "rule",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_network_optimization.default.site_id}",
			"ids":         []string{"${alicloud_esa_network_optimization.default.id}"},
			"config_type": "global",
		}),
	}

	ruleNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_network_optimization.default.site_id}",
			"rule_name": "${alicloud_esa_network_optimization.default.rule_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_network_optimization.default.site_id}",
			"rule_name": "${alicloud_esa_network_optimization.default.rule_name}_fake",
		}),
	}

	siteVersionConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_network_optimization.default.site_id}",
			"ids":          []string{"${alicloud_esa_network_optimization.default.id}"},
			"site_version": "${alicloud_esa_network_optimization.default.site_version}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_network_optimization.default.site_id}",
			"ids":          []string{"${alicloud_esa_network_optimization.default.id}"},
			"site_version": "1",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_network_optimization.default.site_id}",
			"ids":         []string{"${alicloud_esa_network_optimization.default.id}"},
			"name_regex":  "${alicloud_esa_network_optimization.default.rule_name}",
			"config_id":   "${alicloud_esa_network_optimization.default.config_id}",
			"config_type": "rule",
			"rule_name":   "${alicloud_esa_network_optimization.default.rule_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_network_optimization.default.site_id}",
			"ids":         []string{"${alicloud_esa_network_optimization.default.id}_fake"},
			"name_regex":  "${alicloud_esa_network_optimization.default.rule_name}_fake",
			"config_id":   "${alicloud_esa_network_optimization.default.config_id}000",
			"config_type": "global",
			"rule_name":   "${alicloud_esa_network_optimization.default.rule_name}_fake",
		}),
	}

	var existAliCloudEsaNetworkOptimizationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"optimizations.#":                     "1",
			"optimizations.0.id":                  acctest.CHECKSET,
			"optimizations.0.config_id":           acctest.CHECKSET,
			"optimizations.0.config_type":         acctest.CHECKSET,
			"optimizations.0.site_version":        acctest.CHECKSET,
			"optimizations.0.rule_enable":         acctest.CHECKSET,
			"optimizations.0.rule_name":           acctest.CHECKSET,
			"optimizations.0.rule":                acctest.CHECKSET,
			"optimizations.0.sequence":            acctest.CHECKSET,
			"optimizations.0.smart_routing":       acctest.CHECKSET,
			"optimizations.0.grpc":                acctest.CHECKSET,
			"optimizations.0.http2_origin":        acctest.CHECKSET,
			"optimizations.0.websocket":           acctest.CHECKSET,
			"optimizations.0.upload_max_filesize": acctest.CHECKSET,
		}
	}

	var fakeAliCloudEsaNetworkOptimizationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           "0",
			"names.#":         "0",
			"optimizations.#": "0",
		}
	}

	var aliCloudEsaNetworkOptimizationsInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_esa_network_optimizations.default",
		ExistMapFunc: existAliCloudEsaNetworkOptimizationsMapFunc,
		FakeMapFunc:  fakeAliCloudEsaNetworkOptimizationsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudEsaNetworkOptimizationsInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, configIdConf, configTypeConf, ruleNameConf, siteVersionConf, allConf)
}

func dataSourceEsaNetworkOptimizationsConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_network_optimization" "default" {
  site_id             = data.alicloud_esa_sites.default.sites.0.id
  site_version        = 0
  rule_enable         = "on"
  rule_name           = var.name
  rule                = "true"
  sequence            = 1
  smart_routing       = "on"
  websocket           = "on"
  http2_origin        = "on"
  grpc                = "on"
  upload_max_filesize = "100"
}
`, name)
}
