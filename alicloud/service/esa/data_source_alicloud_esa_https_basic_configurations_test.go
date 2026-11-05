package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEsaHttpsBasicConfigurationsDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_esa_https_basic_configurations.default"
	name := fmt.Sprintf("tf-testAcc-EsaHttpsBasicConfiguration%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceEsaHttpsBasicConfigurationsConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id": "${alicloud_esa_https_basic_configuration.default.site_id}",
			"ids":     []string{"${alicloud_esa_https_basic_configuration.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id": "${alicloud_esa_https_basic_configuration.default.site_id}",
			"ids":     []string{"${alicloud_esa_https_basic_configuration.default.id}_fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":    "${alicloud_esa_https_basic_configuration.default.site_id}",
			"name_regex": "${alicloud_esa_https_basic_configuration.default.rule_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id":    "${alicloud_esa_https_basic_configuration.default.site_id}",
			"name_regex": "${alicloud_esa_https_basic_configuration.default.rule_name}_fake",
		}),
	}

	configIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_https_basic_configuration.default.site_id}",
			"config_id": "${alicloud_esa_https_basic_configuration.default.config_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_https_basic_configuration.default.site_id}",
			"config_id": "${alicloud_esa_https_basic_configuration.default.config_id}000",
		}),
	}

	configTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_https_basic_configuration.default.site_id}",
			"ids":         []string{"${alicloud_esa_https_basic_configuration.default.id}"},
			"config_type": "rule",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_https_basic_configuration.default.site_id}",
			"ids":         []string{"${alicloud_esa_https_basic_configuration.default.id}"},
			"config_type": "global",
		}),
	}

	ruleNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_https_basic_configuration.default.site_id}",
			"rule_name": "${alicloud_esa_https_basic_configuration.default.rule_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_https_basic_configuration.default.site_id}",
			"rule_name": "${alicloud_esa_https_basic_configuration.default.rule_name}_fake",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_https_basic_configuration.default.site_id}",
			"ids":         []string{"${alicloud_esa_https_basic_configuration.default.id}"},
			"name_regex":  "${alicloud_esa_https_basic_configuration.default.rule_name}",
			"config_id":   "${alicloud_esa_https_basic_configuration.default.config_id}",
			"config_type": "rule",
			"rule_name":   "${alicloud_esa_https_basic_configuration.default.rule_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_https_basic_configuration.default.site_id}",
			"ids":         []string{"${alicloud_esa_https_basic_configuration.default.id}_fake"},
			"name_regex":  "${alicloud_esa_https_basic_configuration.default.rule_name}_fake",
			"config_id":   "${alicloud_esa_https_basic_configuration.default.config_id}000",
			"config_type": "global",
			"rule_name":   "${alicloud_esa_https_basic_configuration.default.rule_name}_fake",
		}),
	}

	var existAliCloudEsaHttpsBasicConfigurationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"names.#":                            "1",
			"configurations.#":                   "1",
			"configurations.0.id":                acctest.CHECKSET,
			"configurations.0.config_id":         acctest.CHECKSET,
			"configurations.0.config_type":       acctest.CHECKSET,
			"configurations.0.rule_enable":       acctest.CHECKSET,
			"configurations.0.rule_name":         acctest.CHECKSET,
			"configurations.0.rule":              acctest.CHECKSET,
			"configurations.0.sequence":          acctest.CHECKSET,
			"configurations.0.ciphersuite":       acctest.CHECKSET,
			"configurations.0.ciphersuite_group": acctest.CHECKSET,
			"configurations.0.ocsp_stapling":     acctest.CHECKSET,
			"configurations.0.http2":             acctest.CHECKSET,
			"configurations.0.http3":             acctest.CHECKSET,
			"configurations.0.https":             acctest.CHECKSET,
			"configurations.0.tls10":             acctest.CHECKSET,
			"configurations.0.tls11":             acctest.CHECKSET,
			"configurations.0.tls12":             acctest.CHECKSET,
			"configurations.0.tls13":             acctest.CHECKSET,
		}
	}

	var fakeAliCloudEsaHttpsBasicConfigurationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            "0",
			"names.#":          "0",
			"configurations.#": "0",
		}
	}

	var aliCloudEsaHttpsBasicConfigurationsInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_esa_https_basic_configurations.default",
		ExistMapFunc: existAliCloudEsaHttpsBasicConfigurationsMapFunc,
		FakeMapFunc:  fakeAliCloudEsaHttpsBasicConfigurationsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudEsaHttpsBasicConfigurationsInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, configIdConf, configTypeConf, ruleNameConf, allConf)
}

func dataSourceEsaHttpsBasicConfigurationsConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

data "alicloud_esa_sites" "default" {
 plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_https_basic_configuration" "default" {
 site_id           = data.alicloud_esa_sites.default.sites.0.id
 rule_enable       = "on"
 rule_name         = var.name
 rule              = "true"
 ciphersuite       = "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256"
 ciphersuite_group = "all"
 tls10             = "on"
 tls11             = "on"
 tls12             = "on"
 tls13             = "on"
 ocsp_stapling     = "on"
 http2             = "on"
 http3             = "on"
 https             = "on"
 sequence          = 1
}
`, name)
}
