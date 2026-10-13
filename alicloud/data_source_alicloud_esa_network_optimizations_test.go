package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAliCloudEsaNetworkOptimizationsDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_esa_network_optimizations.default"
	name := fmt.Sprintf("tf-testAcc-EsaNetworkOptimization%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEsaNetworkOptimizationsConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id": "${alicloud_esa_network_optimization.default.site_id}",
			"ids":     []string{"${alicloud_esa_network_optimization.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id": "${alicloud_esa_network_optimization.default.site_id}",
			"ids":     []string{"${alicloud_esa_network_optimization.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":    "${alicloud_esa_network_optimization.default.site_id}",
			"name_regex": "${alicloud_esa_network_optimization.default.rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":    "${alicloud_esa_network_optimization.default.site_id}",
			"name_regex": "${alicloud_esa_network_optimization.default.rule_name}_fake",
		}),
	}

	configIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_network_optimization.default.site_id}",
			"config_id": "${alicloud_esa_network_optimization.default.config_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_network_optimization.default.site_id}",
			"config_id": "${alicloud_esa_network_optimization.default.config_id}000",
		}),
	}

	configTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_network_optimization.default.site_id}",
			"ids":         []string{"${alicloud_esa_network_optimization.default.id}"},
			"config_type": "rule",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_network_optimization.default.site_id}",
			"ids":         []string{"${alicloud_esa_network_optimization.default.id}"},
			"config_type": "global",
		}),
	}

	ruleNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_network_optimization.default.site_id}",
			"rule_name": "${alicloud_esa_network_optimization.default.rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_network_optimization.default.site_id}",
			"rule_name": "${alicloud_esa_network_optimization.default.rule_name}_fake",
		}),
	}

	siteVersionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_network_optimization.default.site_id}",
			"ids":          []string{"${alicloud_esa_network_optimization.default.id}"},
			"site_version": "${alicloud_esa_network_optimization.default.site_version}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":      "${alicloud_esa_network_optimization.default.site_id}",
			"ids":          []string{"${alicloud_esa_network_optimization.default.id}"},
			"site_version": "1",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_network_optimization.default.site_id}",
			"ids":         []string{"${alicloud_esa_network_optimization.default.id}"},
			"name_regex":  "${alicloud_esa_network_optimization.default.rule_name}",
			"config_id":   "${alicloud_esa_network_optimization.default.config_id}",
			"config_type": "rule",
			"rule_name":   "${alicloud_esa_network_optimization.default.rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
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
			"optimizations.0.id":                  CHECKSET,
			"optimizations.0.config_id":           CHECKSET,
			"optimizations.0.config_type":         CHECKSET,
			"optimizations.0.site_version":        CHECKSET,
			"optimizations.0.rule_enable":         CHECKSET,
			"optimizations.0.rule_name":           CHECKSET,
			"optimizations.0.rule":                CHECKSET,
			"optimizations.0.sequence":            CHECKSET,
			"optimizations.0.smart_routing":       CHECKSET,
			"optimizations.0.grpc":                CHECKSET,
			"optimizations.0.http2_origin":        CHECKSET,
			"optimizations.0.websocket":           CHECKSET,
			"optimizations.0.upload_max_filesize": CHECKSET,
		}
	}

	var fakeAliCloudEsaNetworkOptimizationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           "0",
			"names.#":         "0",
			"optimizations.#": "0",
		}
	}

	var aliCloudEsaNetworkOptimizationsInfo = dataSourceAttr{
		resourceId:   "data.alicloud_esa_network_optimizations.default",
		existMapFunc: existAliCloudEsaNetworkOptimizationsMapFunc,
		fakeMapFunc:  fakeAliCloudEsaNetworkOptimizationsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudEsaNetworkOptimizationsInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, configIdConf, configTypeConf, ruleNameConf, siteVersionConf, allConf)
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
