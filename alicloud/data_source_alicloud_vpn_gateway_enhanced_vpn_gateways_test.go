// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlicloudVpnGatewayEnhancedVpnGatewayDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"ap-southeast-3"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewayEnhancedVpnGatewaySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpn_gateway_enhanced_vpn_gateway.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewayEnhancedVpnGatewaySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpn_gateway_enhanced_vpn_gateway.default.id}_fake"]`,
		}),
	}

	VpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewayEnhancedVpnGatewaySourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_vpn_gateway_enhanced_vpn_gateway.default.id}"]`,
			"vpc_id": `"${alicloud_vpc.defaulttYTx5F.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewayEnhancedVpnGatewaySourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_vpn_gateway_enhanced_vpn_gateway.default.id}_fake"]`,
			"vpc_id": `"${alicloud_vpc.defaulttYTx5F.id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewayEnhancedVpnGatewaySourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_vpn_gateway_enhanced_vpn_gateway.default.id}"]`,
			"vpc_id": `"${alicloud_vpc.defaulttYTx5F.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewayEnhancedVpnGatewaySourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_vpn_gateway_enhanced_vpn_gateway.default.id}_fake"]`,
			"vpc_id": `"${alicloud_vpc.defaulttYTx5F.id}_fake"`,
		}),
	}

	VpnGatewayEnhancedVpnGatewayCheckInfo.dataSourceTestCheck(t, rand, idsConf, VpcIdConf, allConf)
}

var existVpnGatewayEnhancedVpnGatewayMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#":                              "1",
		"gateways.0.status":                       CHECKSET,
		"gateways.0.vpn_type":                     CHECKSET,
		"gateways.0.description":                  CHECKSET,
		"gateways.0.disaster_recovery_vswitch_id": CHECKSET,
		"gateways.0.vpn_gateway_name":             CHECKSET,
		"gateways.0.create_time":                  CHECKSET,
		"gateways.0.vswitch_id":                   CHECKSET,
		"gateways.0.gateway_type":                 CHECKSET,
		"gateways.0.auto_propagate":               CHECKSET,
		"gateways.0.vpn_instance_id":              CHECKSET,
		"gateways.0.vpc_id":                       CHECKSET,
		"gateways.0.network_type":                 CHECKSET,
		"gateways.0.tags.%":                       CHECKSET,
	}
}

var fakeVpnGatewayEnhancedVpnGatewayMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#": "0",
	}
}

var VpnGatewayEnhancedVpnGatewayCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_vpn_gateway_enhanced_vpn_gateways.default",
	existMapFunc: existVpnGatewayEnhancedVpnGatewayMapFunc,
	fakeMapFunc:  fakeVpnGatewayEnhancedVpnGatewayMapFunc,
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
