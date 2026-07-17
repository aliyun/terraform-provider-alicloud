package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAliCloudApigGatewayDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_gateway.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_gateway.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_apig_gateway.default.gateway_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_apig_gateway.default.gateway_name}_fake"`,
		}),
	}

	gatewayNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_apig_gateway.default.id}"]`,
			"gateway_name": `"${alicloud_apig_gateway.default.gateway_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_apig_gateway.default.id}"]`,
			"gateway_name": `"${alicloud_apig_gateway.default.gateway_name}_fake"`,
		}),
	}

	resourceGroupConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_gateway.default.id}"]`,
			"resource_group_id": `"${alicloud_apig_gateway.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_gateway.default.id}"]`,
			"resource_group_id": `"${alicloud_apig_gateway.default.resource_group_id}_fake"`,
		}),
	}

	enableDetailsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_gateway.default.id}"]`,
			"enable_details": `"true"`,
		}),
		fakeConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_gateway.default.id}_fake"]`,
			"enable_details": `"true"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_gateway.default.id}"]`,
			"name_regex":        `"${alicloud_apig_gateway.default.gateway_name}"`,
			"gateway_name":      `"${alicloud_apig_gateway.default.gateway_name}"`,
			"resource_group_id": `"${alicloud_apig_gateway.default.resource_group_id}"`,
			"enable_details":    `"true"`,
		}),
		fakeConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_gateway.default.id}_fake"]`,
			"name_regex":        `"${alicloud_apig_gateway.default.gateway_name}_fake"`,
			"gateway_name":      `"${alicloud_apig_gateway.default.gateway_name}_fake"`,
			"resource_group_id": `"${alicloud_apig_gateway.default.resource_group_id}"`,
			"enable_details":    `"true"`,
		}),
	}

	AliCloudApigGatewaysCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, gatewayNameConf, resourceGroupConf, enableDetailsConf, allConf)
}

var existAliCloudApigGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#":                   "1",
		"gateways.0.id":                CHECKSET,
		"gateways.0.gateway_id":        CHECKSET,
		"gateways.0.gateway_name":      CHECKSET,
		"gateways.0.payment_type":      "PayAsYouGo",
		"gateways.0.spec":              CHECKSET,
		"gateways.0.status":            CHECKSET,
		"gateways.0.create_time":       CHECKSET,
		"gateways.0.resource_group_id": CHECKSET,
		"ids.#":                        "1",
		"names.#":                      "1",
	}
}

var fakeAliCloudApigGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#": "0",
		"ids.#":      "0",
		"names.#":    "0",
	}
}

var AliCloudApigGatewaysCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_apig_gateways.default",
	existMapFunc: existAliCloudApigGatewaysMapFunc,
	fakeMapFunc:  fakeAliCloudApigGatewaysMapFunc,
}

func testAccCheckAliCloudApigGatewaysDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testacc-apig-gateway-ds-%d"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "default" {
  gateway_name = var.name
  spec         = "apigw.small.x1"
  payment_type = "PayAsYouGo"
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  zone_config {
    select_option = "Auto"
  }
  network_access_config {
    type = "Intranet"
  }
  log_config {
    sls {
      enable = "false"
    }
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.1
}

data "alicloud_apig_gateways" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
