package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudEsaHttpsBasicConfigurationsDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_esa_https_basic_configurations.default"
	name := fmt.Sprintf("tf-testAcc-EsaHttpsBasicConfiguration%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEsaHttpsBasicConfigurationsConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id": "${alicloud_esa_https_basic_configuration.default.site_id}",
			"ids":     []string{"${alicloud_esa_https_basic_configuration.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id": "${alicloud_esa_https_basic_configuration.default.site_id}",
			"ids":     []string{"${alicloud_esa_https_basic_configuration.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":    "${alicloud_esa_https_basic_configuration.default.site_id}",
			"name_regex": "${alicloud_esa_https_basic_configuration.default.rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":    "${alicloud_esa_https_basic_configuration.default.site_id}",
			"name_regex": "${alicloud_esa_https_basic_configuration.default.rule_name}_fake",
		}),
	}

	configIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_https_basic_configuration.default.site_id}",
			"config_id": "${alicloud_esa_https_basic_configuration.default.config_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_https_basic_configuration.default.site_id}",
			"config_id": "${alicloud_esa_https_basic_configuration.default.config_id}000",
		}),
	}

	configTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_https_basic_configuration.default.site_id}",
			"ids":         []string{"${alicloud_esa_https_basic_configuration.default.id}"},
			"config_type": "rule",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_https_basic_configuration.default.site_id}",
			"ids":         []string{"${alicloud_esa_https_basic_configuration.default.id}"},
			"config_type": "global",
		}),
	}

	ruleNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_https_basic_configuration.default.site_id}",
			"rule_name": "${alicloud_esa_https_basic_configuration.default.rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"site_id":   "${alicloud_esa_https_basic_configuration.default.site_id}",
			"rule_name": "${alicloud_esa_https_basic_configuration.default.rule_name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"site_id":     "${alicloud_esa_https_basic_configuration.default.site_id}",
			"ids":         []string{"${alicloud_esa_https_basic_configuration.default.id}"},
			"name_regex":  "${alicloud_esa_https_basic_configuration.default.rule_name}",
			"config_id":   "${alicloud_esa_https_basic_configuration.default.config_id}",
			"config_type": "rule",
			"rule_name":   "${alicloud_esa_https_basic_configuration.default.rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
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
			"configurations.0.id":                CHECKSET,
			"configurations.0.config_id":         CHECKSET,
			"configurations.0.config_type":       CHECKSET,
			"configurations.0.rule_enable":       CHECKSET,
			"configurations.0.rule_name":         CHECKSET,
			"configurations.0.rule":              CHECKSET,
			"configurations.0.sequence":          CHECKSET,
			"configurations.0.ciphersuite":       CHECKSET,
			"configurations.0.ciphersuite_group": CHECKSET,
			"configurations.0.ocsp_stapling":     CHECKSET,
			"configurations.0.http2":             CHECKSET,
			"configurations.0.http3":             CHECKSET,
			"configurations.0.https":             CHECKSET,
			"configurations.0.tls10":             CHECKSET,
			"configurations.0.tls11":             CHECKSET,
			"configurations.0.tls12":             CHECKSET,
			"configurations.0.tls13":             CHECKSET,
		}
	}

	var fakeAliCloudEsaHttpsBasicConfigurationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            "0",
			"names.#":          "0",
			"configurations.#": "0",
		}
	}

	var aliCloudEsaHttpsBasicConfigurationsInfo = dataSourceAttr{
		resourceId:   "data.alicloud_esa_https_basic_configurations.default",
		existMapFunc: existAliCloudEsaHttpsBasicConfigurationsMapFunc,
		fakeMapFunc:  fakeAliCloudEsaHttpsBasicConfigurationsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudEsaHttpsBasicConfigurationsInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, configIdConf, configTypeConf, ruleNameConf, allConf)
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
