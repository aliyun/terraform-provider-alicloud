// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package vpngateway_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudVpnGatewayEnhancedVpnGatewayDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"ap-southeast-3"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnGatewayEnhancedVpnGatewaySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpn_gateway_enhanced_vpn_gateway.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudVpnGatewayEnhancedVpnGatewaySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpn_gateway_enhanced_vpn_gateway.default.id}_fake"]`,
		}),
	}

	VpcIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnGatewayEnhancedVpnGatewaySourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_vpn_gateway_enhanced_vpn_gateway.default.id}"]`,
			"vpc_id": `"${alicloud_vpc.defaulttYTx5F.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpnGatewayEnhancedVpnGatewaySourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_vpn_gateway_enhanced_vpn_gateway.default.id}_fake"]`,
			"vpc_id": `"${alicloud_vpc.defaulttYTx5F.id}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnGatewayEnhancedVpnGatewaySourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_vpn_gateway_enhanced_vpn_gateway.default.id}"]`,
			"vpc_id": `"${alicloud_vpc.defaulttYTx5F.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpnGatewayEnhancedVpnGatewaySourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_vpn_gateway_enhanced_vpn_gateway.default.id}_fake"]`,
			"vpc_id": `"${alicloud_vpc.defaulttYTx5F.id}_fake"`,
		}),
	}

	VpnGatewayEnhancedVpnGatewayCheckInfo.DataSourceTestCheck(t, rand, idsConf, VpcIdConf, allConf)
}

var existVpnGatewayEnhancedVpnGatewayMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#":                              "1",
		"gateways.0.status":                       acctest.CHECKSET,
		"gateways.0.vpn_type":                     acctest.CHECKSET,
		"gateways.0.description":                  acctest.CHECKSET,
		"gateways.0.disaster_recovery_vswitch_id": acctest.CHECKSET,
		"gateways.0.vpn_gateway_name":             acctest.CHECKSET,
		"gateways.0.create_time":                  acctest.CHECKSET,
		"gateways.0.vswitch_id":                   acctest.CHECKSET,
		"gateways.0.gateway_type":                 acctest.CHECKSET,
		"gateways.0.auto_propagate":               acctest.CHECKSET,
		"gateways.0.vpn_instance_id":              acctest.CHECKSET,
		"gateways.0.vpc_id":                       acctest.CHECKSET,
		"gateways.0.network_type":                 acctest.CHECKSET,
		"gateways.0.tags.%":                       acctest.CHECKSET,
	}
}

var fakeVpnGatewayEnhancedVpnGatewayMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#": "0",
	}
}

var VpnGatewayEnhancedVpnGatewayCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_vpn_gateway_enhanced_vpn_gateways.default",
	ExistMapFunc: existVpnGatewayEnhancedVpnGatewayMapFunc,
	FakeMapFunc:  fakeVpnGatewayEnhancedVpnGatewayMapFunc,
}

func testAccCheckAlicloudVpnGatewayEnhancedVpnGatewaySourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccVpnGatewayEnhancedVpnGateway%d"
}
variable "region" {
  default = "ap-southeast-3"
}

variable "zone2" {
  default = "ap-southeast-3a"
}

variable "zone1" {
  default = "ap-southeast-3b"
}

resource "alicloud_vpc" "defaulttYTx5F" {
  cidr_block = "192.168.0.0/16"
  is_default = false
}

resource "alicloud_vswitch" "defaultTRk7k3" {
  vpc_id     = alicloud_vpc.defaulttYTx5F.id
  zone_id    = var.zone1
  cidr_block = "192.168.10.0/24"
}

resource "alicloud_vswitch" "default23kGFr" {
  vpc_id     = alicloud_vpc.defaulttYTx5F.id
  zone_id    = var.zone2
  cidr_block = "192.168.20.0/24"
}



resource "alicloud_vpn_gateway_enhanced_vpn_gateway" "default" {
  vpn_type                     = "Normal"
  description                  = "default"
  disaster_recovery_vswitch_id = alicloud_vswitch.default23kGFr.id
  vpc_id                       = alicloud_vpc.defaulttYTx5F.id
  vpn_gateway_name             = "default"
  network_type                 = "public"
  vswitch_id                   = alicloud_vswitch.defaultTRk7k3.id
  gateway_type                 = "Enhanced.SiteToSite"
  auto_propagate               = false
  tags = {
    default_key1 = "default_value1"
    default_key2 = "default_value2"
    default_key3 = "default_value2"
  }
}

data "alicloud_vpn_gateway_enhanced_vpn_gateways" "default" {
  enable_details = true
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
