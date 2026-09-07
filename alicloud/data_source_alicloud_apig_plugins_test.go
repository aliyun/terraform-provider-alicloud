// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudApigPluginDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigPluginSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_plugin.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudApigPluginSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_plugin.default.id}_fake"]`,
		}),
	}

	PluginClassIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigPluginSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_apig_plugin.default.id}"]`,
			"plugin_class_id": `"pls-crpqb35lhtgo800k2m86"`,
		}),
		fakeConfig: testAccCheckAlicloudApigPluginSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_apig_plugin.default.id}_fake"]`,
			"plugin_class_id": `"pls-crpqb35lhtgo800k2m86_fake"`,
		}),
	}
	GatewayIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigPluginSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_apig_plugin.default.id}"]`,
			"gateway_id": `"${alicloud_apig_gateway.plugin_gateway_pre.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudApigPluginSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_apig_plugin.default.id}_fake"]`,
			"gateway_id": `"${alicloud_apig_gateway.plugin_gateway_pre.id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigPluginSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_apig_plugin.default.id}"]`,
			"plugin_class_id": `"pls-crpqb35lhtgo800k2m86"`,

			"gateway_id": `"${alicloud_apig_gateway.plugin_gateway_pre.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudApigPluginSourceConfig(rand, map[string]string{
			"ids":             `["${alicloud_apig_plugin.default.id}_fake"]`,
			"plugin_class_id": `"pls-crpqb35lhtgo800k2m86_fake"`,

			"gateway_id": `"${alicloud_apig_gateway.plugin_gateway_pre.id}_fake"`,
		}),
	}

	ApigPluginCheckInfo.dataSourceTestCheck(t, rand, idsConf, PluginClassIdConf, GatewayIdConf, allConf)
}

var existApigPluginMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"plugins.#":                   "1",
		"plugins.0.plugin_class_name": CHECKSET,
		"plugins.0.gateway_name":      CHECKSET,
		"plugins.0.plugin_id":         CHECKSET,
	}
}

var fakeApigPluginMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"plugins.#": "0",
	}
}

var ApigPluginCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_apig_plugins.default",
	existMapFunc: existApigPluginMapFunc,
	fakeMapFunc:  fakeApigPluginMapFunc,
}

func testAccCheckAlicloudApigPluginSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccApigPlugin%d"
}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "plugin_gateway_pre" {
  network_access_config {
    type = "Internet"
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  zone_config {
    select_option = "Auto"
  }
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  gateway_type = "API"
  payment_type = "PayAsYouGo"
  gateway_name = var.name
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = true
    }
  }
}



resource "alicloud_apig_plugin" "default" {
  plugin_class_id = "pls-crpqb35lhtgo800k2m86"
  gateway_id      = alicloud_apig_gateway.plugin_gateway_pre.id
}

data "alicloud_apig_plugins" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
