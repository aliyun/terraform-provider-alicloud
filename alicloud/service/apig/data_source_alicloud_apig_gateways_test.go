package apig_test

import (
	"fmt"
	"strings"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudApigGatewayDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_gateway.default.id}"]`,
		}),
		FakeConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_gateway.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_apig_gateway.default.gateway_name}"`,
		}),
		FakeConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_apig_gateway.default.gateway_name}_fake"`,
		}),
	}

	gatewayNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_apig_gateway.default.id}"]`,
			"gateway_name": `"${alicloud_apig_gateway.default.gateway_name}"`,
		}),
		FakeConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_apig_gateway.default.id}"]`,
			"gateway_name": `"${alicloud_apig_gateway.default.gateway_name}_fake"`,
		}),
	}

	resourceGroupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_gateway.default.id}"]`,
			"resource_group_id": `"${alicloud_apig_gateway.default.resource_group_id}"`,
		}),
		FakeConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_gateway.default.id}"]`,
			"resource_group_id": `"${alicloud_apig_gateway.default.resource_group_id}_fake"`,
		}),
	}

	enableDetailsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_gateway.default.id}"]`,
			"enable_details": `"true"`,
		}),
		FakeConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_apig_gateway.default.id}_fake"]`,
			"enable_details": `"true"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_gateway.default.id}"]`,
			"name_regex":        `"${alicloud_apig_gateway.default.gateway_name}"`,
			"gateway_name":      `"${alicloud_apig_gateway.default.gateway_name}"`,
			"resource_group_id": `"${alicloud_apig_gateway.default.resource_group_id}"`,
			"enable_details":    `"true"`,
		}),
		FakeConfig: testAccCheckAliCloudApigGatewaysDataSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_gateway.default.id}_fake"]`,
			"name_regex":        `"${alicloud_apig_gateway.default.gateway_name}_fake"`,
			"gateway_name":      `"${alicloud_apig_gateway.default.gateway_name}_fake"`,
			"resource_group_id": `"${alicloud_apig_gateway.default.resource_group_id}"`,
			"enable_details":    `"true"`,
		}),
	}

	AliCloudApigGatewaysCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, gatewayNameConf, resourceGroupConf, enableDetailsConf, allConf)
}

var existAliCloudApigGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#":                   "1",
		"gateways.0.id":                acctest.CHECKSET,
		"gateways.0.gateway_id":        acctest.CHECKSET,
		"gateways.0.gateway_name":      acctest.CHECKSET,
		"gateways.0.payment_type":      "PayAsYouGo",
		"gateways.0.spec":              acctest.CHECKSET,
		"gateways.0.status":            acctest.CHECKSET,
		"gateways.0.create_time":       acctest.CHECKSET,
		"gateways.0.resource_group_id": acctest.CHECKSET,
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

var AliCloudApigGatewaysCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_apig_gateways.default",
	ExistMapFunc: existAliCloudApigGatewaysMapFunc,
	FakeMapFunc:  fakeAliCloudApigGatewaysMapFunc,
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
