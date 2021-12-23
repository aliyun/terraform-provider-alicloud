package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCsgGatewayDataSource(t *testing.T) {
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
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCsgGatewayDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_cloud_storage_gateway_gateway.default.id}"]`,
			"storage_bundle_id": `"${alicloud_cloud_storage_gateway_gateway.default.storage_bundle_id}"`,
			"name_regex":        `"${alicloud_cloud_storage_gateway_gateway.default.gateway_name}"`,
			"status":            `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudCsgGatewayDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_cloud_storage_gateway_gateway.default.id}_fake"]`,
			"storage_bundle_id": `"${alicloud_cloud_storage_gateway_gateway.default.storage_bundle_id}"`,
			"name_regex":        `"${alicloud_cloud_storage_gateway_gateway.default.gateway_name}_fake"`,
			"status":            `"Unknown"`,
		}),
	}
	var existAlicloudCsgGatewayDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               CHECKSET,
			"names.#":                             CHECKSET,
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
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.CsgSupportRegions)
	}
	alicloudCsgGatewayCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
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

resource "alicloud_vpc" "vpc" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

data "alicloud_zones" "default"{
  available_resource_creation = "VSwitch"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.vpc.id
  cidr_block        = "172.16.0.0/21"
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alicloud_cloud_storage_gateway_storage_bundle" "example" {
  storage_bundle_name = var.name
}

resource "alicloud_cloud_storage_gateway_gateway" "default" {
  description = "tf-acctestDesalone"
  gateway_class = "Standard"
  type = "File"
  payment_type = "PayAsYouGo"
  vswitch_id = alicloud_vswitch.default.id
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
