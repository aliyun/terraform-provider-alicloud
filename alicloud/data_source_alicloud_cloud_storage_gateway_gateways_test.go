package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudStorageGatewayGatewayDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCsgGatewayDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_cloud_storage_gateway_gateway.default.id}"]`,
			"storage_bundle_id": `"${alicloud_cloud_storage_gateway_gateway.default.storage_bundle_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCsgGatewayDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_cloud_storage_gateway_gateway.default.id}_fake"]`,
			"storage_bundle_id": `"${alicloud_cloud_storage_gateway_gateway.default.storage_bundle_id}"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCsgGatewayDataSourceName(rand, map[string]string{
			"name_regex":        `"${alicloud_cloud_storage_gateway_gateway.default.gateway_name}"`,
			"storage_bundle_id": `"${alicloud_cloud_storage_gateway_gateway.default.storage_bundle_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCsgGatewayDataSourceName(rand, map[string]string{
			"name_regex":        `"${alicloud_cloud_storage_gateway_gateway.default.gateway_name}_fake"`,
			"storage_bundle_id": `"${alicloud_cloud_storage_gateway_gateway.default.storage_bundle_id}"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCsgGatewayDataSourceName(rand, map[string]string{
			"storage_bundle_id": `"${alicloud_cloud_storage_gateway_gateway.default.storage_bundle_id}"`,
			"status":            `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudCsgGatewayDataSourceName(rand, map[string]string{
			"storage_bundle_id": `"${alicloud_cloud_storage_gateway_gateway.default.storage_bundle_id}"`,
			"status":            `"Unknown"`,
		}),
	}
	pagingConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCsgGatewayDataSourceName(rand, map[string]string{
			"storage_bundle_id": `"${alicloud_cloud_storage_gateway_gateway.default.storage_bundle_id}"`,
			"page_number":       `1`,
		}),
		fakeConfig: testAccCheckAlicloudCsgGatewayDataSourceName(rand, map[string]string{
			"storage_bundle_id": `"${alicloud_cloud_storage_gateway_gateway.default.storage_bundle_id}"`,
			"page_number":       `2`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCsgGatewayDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_cloud_storage_gateway_gateway.default.id}"]`,
			"storage_bundle_id": `"${alicloud_cloud_storage_gateway_gateway.default.storage_bundle_id}"`,
			"name_regex":        `"${alicloud_cloud_storage_gateway_gateway.default.gateway_name}"`,
			"status":            `"Running"`,
			"page_number":       `1`,
		}),
		fakeConfig: testAccCheckAlicloudCsgGatewayDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_cloud_storage_gateway_gateway.default.id}_fake"]`,
			"storage_bundle_id": `"${alicloud_cloud_storage_gateway_gateway.default.storage_bundle_id}"`,
			"name_regex":        `"${alicloud_cloud_storage_gateway_gateway.default.gateway_name}_fake"`,
			"status":            `"Unknown"`,
			"page_number":       `2`,
		}),
	}
	var existAlicloudCsgGatewayDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               CHECKSET,
			"names.#":                             CHECKSET,
			"total_count":                         CHECKSET,
			"gateways.#":                          CHECKSET,
			"gateways.0.gateway_class":            "Standard",
			"gateways.0.storage_bundle_id":        CHECKSET,
			"gateways.0.gateway_name":             CHECKSET,
			"gateways.0.location":                 "Cloud",
			"gateways.0.type":                     "File",
			"gateways.0.payment_type":             "PayAsYouGo",
			"gateways.0.public_network_bandwidth": "10",
			"gateways.0.status":                   CHECKSET,
			"gateways.0.vswitch_id":               CHECKSET,
		}
	}
	var fakeAlicloudCsgGatewayDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCsgGatewayCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cloud_storage_gateway_gateways.default",
		existMapFunc: existAlicloudCsgGatewayDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCsgGatewayDataSourceNameMapFunc,
	}
	alicloudCsgGatewayCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusConf, pagingConf, allConf)
}
func testAccCheckAlicloudCsgGatewayDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccCsgName-%d"
}

data "alicloud_zones" "default"{
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones[0].id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "example" {
  storage_bundle_name = var.name
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
  description = "tf-acctestDesalone"
  gateway_class = "Standard"
  type = "File"
  payment_type = "PayAsYouGo"
  vswitch_id = local.vswitch_id
  release_after_expiration = true
  public_network_bandwidth = 10
  storage_bundle_id = alicloud_cloud_storage_gateway_storage_bundle.example.id
  location = "Cloud"
  gateway_name = var.name
}


data "alicloud_cloud_storage_gateway_gateways" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
