package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCdnDomainConfigsDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_cdn_domain_configs.default"
	name := fmt.Sprintf("tf-testacc%s%d.alicloud-provider.cn", defaultRegionToTest, rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceCdnDomainConfigsConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_cdn_domain_config.default.domain_name}",
			"ids":         []string{"${alicloud_cdn_domain_config.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_cdn_domain_config.default.domain_name}",
			"ids":         []string{"${alicloud_cdn_domain_config.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_cdn_domain_config.default.domain_name}",
			"name_regex":  "${alicloud_cdn_domain_config.default.function_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_cdn_domain_config.default.domain_name}",
			"name_regex":  "${alicloud_cdn_domain_config.default.function_name}_fake",
		}),
	}

	functionNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name":   "${alicloud_cdn_domain_config.default.domain_name}",
			"function_name": "${alicloud_cdn_domain_config.default.function_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name":   "${alicloud_cdn_domain_config.default.domain_name}",
			"ids":           []string{"${alicloud_cdn_domain_config.default.id}_fake"},
			"function_name": "${alicloud_cdn_domain_config.default.function_name}",
		}),
	}

	configIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_cdn_domain_config.default.domain_name}",
			"config_id":   "${alicloud_cdn_domain_config.default.config_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_cdn_domain_config.default.domain_name}",
			"config_id":   "${alicloud_cdn_domain_config.default.config_id}0",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_cdn_domain_config.default.domain_name}",
			"ids":         []string{"${alicloud_cdn_domain_config.default.id}"},
			"status":      "success",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name": "${alicloud_cdn_domain_config.default.domain_name}",
			"ids":         []string{"${alicloud_cdn_domain_config.default.id}"},
			"status":      "failed",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"domain_name":   "${alicloud_cdn_domain_config.default.domain_name}",
			"ids":           []string{"${alicloud_cdn_domain_config.default.id}"},
			"name_regex":    "${alicloud_cdn_domain_config.default.function_name}",
			"function_name": "${alicloud_cdn_domain_config.default.function_name}",
			"config_id":     "${alicloud_cdn_domain_config.default.config_id}",
			"status":        "success",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"domain_name":   "${alicloud_cdn_domain_config.default.domain_name}",
			"ids":           []string{"${alicloud_cdn_domain_config.default.id}_fake"},
			"name_regex":    "${alicloud_cdn_domain_config.default.function_name}_fake",
			"function_name": "${alicloud_cdn_domain_config.default.function_name}",
			"config_id":     "${alicloud_cdn_domain_config.default.config_id}0",
			"status":        "failed",
		}),
	}

	var existAliCloudCdnDomainConfigsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"configs.#":                           "1",
			"configs.0.id":                        CHECKSET,
			"configs.0.function_name":             CHECKSET,
			"configs.0.config_id":                 CHECKSET,
			"configs.0.parent_id":                 CHECKSET,
			"configs.0.status":                    CHECKSET,
			"configs.0.function_args.#":           CHECKSET,
			"configs.0.function_args.0.arg_name":  CHECKSET,
			"configs.0.function_args.0.arg_value": CHECKSET,
		}
	}

	var fakeAliCloudCdnDomainConfigsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"configs.#": "0",
		}
	}

	var aliCloudCdnDomainConfigsInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cdn_domain_configs.default",
		existMapFunc: existAliCloudCdnDomainConfigsMapFunc,
		fakeMapFunc:  fakeAliCloudCdnDomainConfigsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudCdnDomainConfigsInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, functionNameConf, configIdConf, statusConf, allConf)
}

func dataSourceCdnDomainConfigsConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

resource "alicloud_cdn_domain_new" "default" {
  domain_name = var.name
  cdn_type    = "web"
  scope       = "overseas"
  sources {
    content  = "www.aliyuntest.com"
    type     = "domain"
    priority = 20
    port     = 80
    weight   = 10
  }
}

resource "alicloud_cdn_domain_config" "default" {
  domain_name   = alicloud_cdn_domain_new.default.domain_name
  function_name = "condition"
  function_args {
    arg_name  = "rule"
    arg_value = "{\"match\":{\"logic\":\"and\",\"criteria\":[{\"matchType\":\"clientipVer\",\"matchObject\":\"CONNECTING_IP\",\"matchOperator\":\"equals\",\"matchValue\":\"v6\",\"negate\":false}]},\"name\":\"example\",\"status\":\"enable\"}"
  }
}
`, name)
}
