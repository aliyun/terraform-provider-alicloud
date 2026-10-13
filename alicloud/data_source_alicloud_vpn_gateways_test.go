package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlicloudVPNGatewaysDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	preCheck := func() {
		testAccPreCheck(t)
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_vpn_gateway.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_vpn_gateway.default.id}_fake" ]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_gateway.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_gateway.default.name}_fake"`,
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_gateway.default.name}"`,
			"vpc_id":     `"${alicloud_vpn_gateway.default.vpc_id}"`,
		}),

		fakeConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_gateway.default.name}"`,
			"vpc_id":     `"${alicloud_vpn_gateway.default.vpc_id}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_gateway.default.name}"`,
			"status":     `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_gateway.default.name}"`,
			"status":     `"Init"`,
		}),
	}

	businessStatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex":      `"${alicloud_vpn_gateway.default.name}"`,
			"business_status": `"Normal"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex":      `"${alicloud_vpn_gateway.default.name}"`,
			"business_status": `"FinancialLocked"`,
		}),
	}

	gatewayTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex":   `"${alicloud_vpn_gateway.default.name}"`,
			"gateway_type": `"Traditional"`,
		}),
		// The default gateway is a Traditional type, so filtering by the
		// Enhanced.SiteToSite type must return no matches.
		fakeConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex":   `"${alicloud_vpn_gateway.default.name}"`,
			"gateway_type": `"Enhanced.SiteToSite"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"ids":             `[ "${alicloud_vpn_gateway.default.id}" ]`,
			"name_regex":      `"${alicloud_vpn_gateway.default.name}"`,
			"vpc_id":          `"${alicloud_vpn_gateway.default.vpc_id}"`,
			"status":          `"Active"`,
			"business_status": `"Normal"`,
			"gateway_type":    `"Traditional"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"ids":             `[ "${alicloud_vpn_gateway.default.id}" ]`,
			"name_regex":      `"${alicloud_vpn_gateway.default.name}"`,
			"vpc_id":          `"${alicloud_vpn_gateway.default.vpc_id}"`,
			"status":          `"Active"`,
			"business_status": `"FinancialLocked"`,
			"gateway_type":    `"Traditional"`,
		}),
	}

	vpnGatewaysCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, vpcIdConf, statusConf, businessStatusConf, gatewayTypeConf, allConf)
}

func testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccVpcGatewayConfig%d"
}

variable "spec" {
  default = "5"
}

data "alicloud_vpn_gateway_zones" "default" {
  spec = "5M"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_vpn_gateway_zones.default.ids.0
}

resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 1)
  zone_id      = data.alicloud_vpn_gateway_zones.default.ids.0
  vswitch_name = var.name
}

data "alicloud_vswitches" "default2" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_vpn_gateway_zones.default.ids.1
}

resource "alicloud_vswitch" "vswitch2" {
  count        = length(data.alicloud_vswitches.default2.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id      = data.alicloud_vpn_gateway_zones.default.ids.1
  vswitch_name = var.name
}

locals {
  vswitch_id  = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  vswitch_id2 = length(data.alicloud_vswitches.default2.ids) > 0 ? data.alicloud_vswitches.default2.ids[0] : concat(alicloud_vswitch.vswitch2.*.id, [""])[0]
}

resource "alicloud_vpn_gateway" "default" {
  vpn_type                     = "Normal"
  disaster_recovery_vswitch_id = local.vswitch_id2
  vpn_gateway_name             = "${var.name}"
  description                  = "${var.name}"

  vswitch_id   = local.vswitch_id
  auto_pay     = true
  vpc_id       = data.alicloud_vpcs.default.ids.0
  network_type = "public"
  payment_type = "Subscription"
  enable_ipsec = true
  bandwidth    = var.spec
  tags = {
    Created = "tfTestAcc0"
    For     = "Tftestacc 0"
  }
}

data "alicloud_vpn_gateways" "default" {
	%s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existVpnGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#":                      "1",
		"ids.#":                           "1",
		"names.#":                         "1",
		"gateways.0.id":                   CHECKSET,
		"gateways.0.vpc_id":               CHECKSET,
		"gateways.0.internet_ip":          CHECKSET,
		"gateways.0.create_time":          CHECKSET,
		"gateways.0.end_time":             CHECKSET,
		"gateways.0.name":                 fmt.Sprintf("tf-testAccVpcGatewayConfig%d", rand),
		"gateways.0.specification":        "5M",
		"gateways.0.description":          fmt.Sprintf("tf-testAccVpcGatewayConfig%d", rand),
		"gateways.0.enable_ssl":           "disable",
		"gateways.0.enable_ipsec":         "enable",
		"gateways.0.status":               "Active",
		"gateways.0.business_status":      "Normal",
		"gateways.0.instance_charge_type": string(PrePaid),
		"gateways.0.ssl_connections":      "0",
		"gateways.0.network_type":         "public",
		"gateways.0.tags.%":               "2",
	}
}

var fakeVpnGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":      "0",
		"names.#":    "0",
		"gateways.#": "0",
	}
}

var vpnGatewaysCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_vpn_gateways.default",
	existMapFunc: existVpnGatewaysMapFunc,
	fakeMapFunc:  fakeVpnGatewaysMapFunc,
}

// TestAccAlicloudVPNGatewaysDataSourceEnhanced covers reading an enhanced (Enhance.SiteToSite)
// VPN gateway through the alicloud_vpn_gateways data source. Enhanced gateways are returned by
// DescribeVpnGateways without a ChargeType/Status field, which used to panic the provider
// (interface conversion: interface {} is nil, not string). It also exercises the new
// gateway_type query argument.
func TestAccAlicloudVPNGatewaysDataSourceEnhanced(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"ap-southeast-3"})
	rand := acctest.RandIntRange(1000000, 9999999)
	preCheck := func() {
		testAccPreCheck(t)
	}

	gatewayTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewaysEnhancedDataSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_vpn_gateway_enhanced_vpn_gateway.default.id}"]`,
			"gateway_type": `"Enhanced.SiteToSite"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewaysEnhancedDataSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_vpn_gateway_enhanced_vpn_gateway.default.id}"]`,
			"gateway_type": `"Traditional"`,
		}),
	}

	vpnGatewaysEnhancedCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, gatewayTypeConf)
}

func testAccCheckAlicloudVpnGatewaysEnhancedDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccVpnGatewaysEnhanced%d"
}

variable "zone1" {
  default = "ap-southeast-3b"
}

variable "zone2" {
  default = "ap-southeast-3a"
}

resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/16"
  is_default = false
}

resource "alicloud_vswitch" "default1" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = var.zone1
  cidr_block = "192.168.10.0/24"
}

resource "alicloud_vswitch" "default2" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = var.zone2
  cidr_block = "192.168.20.0/24"
}

resource "alicloud_vpn_gateway_enhanced_vpn_gateway" "default" {
  vpn_type                     = "Normal"
  description                  = var.name
  disaster_recovery_vswitch_id = alicloud_vswitch.default2.id
  vpc_id                       = alicloud_vpc.default.id
  vpn_gateway_name             = var.name
  network_type                 = "public"
  vswitch_id                   = alicloud_vswitch.default1.id
  gateway_type                 = "Enhanced.SiteToSite"
  auto_propagate               = false
}

data "alicloud_vpn_gateways" "default" {
	%s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existVpnGatewaysEnhancedMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#":              "1",
		"ids.#":                   "1",
		"gateways.0.id":           CHECKSET,
		"gateways.0.vpc_id":       CHECKSET,
		"gateways.0.status":       CHECKSET,
		"gateways.0.network_type": CHECKSET,
	}
}

var fakeVpnGatewaysEnhancedMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":      "0",
		"gateways.#": "0",
	}
}

var vpnGatewaysEnhancedCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_vpn_gateways.default",
	existMapFunc: existVpnGatewaysEnhancedMapFunc,
	fakeMapFunc:  fakeVpnGatewaysEnhancedMapFunc,
}
