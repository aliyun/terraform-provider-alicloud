package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
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

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"ids":             `[ "${alicloud_vpn_gateway.default.id}" ]`,
			"name_regex":      `"${alicloud_vpn_gateway.default.name}"`,
			"vpc_id":          `"${alicloud_vpn_gateway.default.vpc_id}"`,
			"status":          `"Active"`,
			"business_status": `"Normal"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnGatewaysDataSourceConfig(rand, map[string]string{
			"ids":             `[ "${alicloud_vpn_gateway.default.id}" ]`,
			"name_regex":      `"${alicloud_vpn_gateway.default.name}"`,
			"vpc_id":          `"${alicloud_vpn_gateway.default.vpc_id}"`,
			"status":          `"Active"`,
			"business_status": `"FinancialLocked"`,
		}),
	}

	vpnGatewaysCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, vpcIdConf, statusConf, businessStatusConf, allConf)
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
  description = "${var.name}"

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
