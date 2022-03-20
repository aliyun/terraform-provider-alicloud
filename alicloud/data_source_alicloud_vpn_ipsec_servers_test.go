package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPNIpsecServersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpn_ipsec_server.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpn_ipsec_server.default.id}_fake"]`,
		}),
	}
	ipsecServerNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_vpn_ipsec_server.default.id}"]`,
			"ipsec_server_name": `"${alicloud_vpn_ipsec_server.default.ipsec_server_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_vpn_ipsec_server.default.id}"]`,
			"ipsec_server_name": `"${alicloud_vpn_ipsec_server.default.ipsec_server_name}_fake"`,
		}),
	}
	vpnGatewayIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_vpn_ipsec_server.default.id}"]`,
			"vpn_gateway_id": `"${alicloud_vpn_ipsec_server.default.vpn_gateway_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_vpn_ipsec_server.default.id}"]`,
			"vpn_gateway_id": `"${alicloud_vpn_ipsec_server.default.vpn_gateway_id}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_ipsec_server.default.ipsec_server_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_ipsec_server.default.ipsec_server_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_vpn_ipsec_server.default.id}"]`,
			"ipsec_server_name": `"${alicloud_vpn_ipsec_server.default.ipsec_server_name}"`,
			"name_regex":        `"${alicloud_vpn_ipsec_server.default.ipsec_server_name}"`,
			"vpn_gateway_id":    `"${alicloud_vpn_ipsec_server.default.vpn_gateway_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_vpn_ipsec_server.default.id}_fake"]`,
			"ipsec_server_name": `"${alicloud_vpn_ipsec_server.default.ipsec_server_name}_fake"`,
			"name_regex":        `"${alicloud_vpn_ipsec_server.default.ipsec_server_name}_fake"`,
			"vpn_gateway_id":    `"${alicloud_vpn_ipsec_server.default.vpn_gateway_id}_fake"`,
		}),
	}
	var existAlicloudVpnIpsecServersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"servers.#":                           "1",
			"servers.0.client_ip_pool":            "10.0.0.0/24",
			"servers.0.effect_immediately":        "true",
			"servers.0.ipsec_server_name":         fmt.Sprintf("tf-testAccIpsecServer-%d", rand),
			"servers.0.local_subnet":              "192.168.0.0/24",
			"servers.0.psk":                       CHECKSET,
			"servers.0.psk_enabled":               "true",
			"servers.0.vpn_gateway_id":            CHECKSET,
			"servers.0.internet_ip":               CHECKSET,
			"servers.0.id":                        CHECKSET,
			"servers.0.idaas_instance_id":         "",
			"servers.0.create_time":               CHECKSET,
			"servers.0.ipsec_server_id":           CHECKSET,
			"servers.0.max_connections":           CHECKSET,
			"servers.0.multi_factor_auth_enabled": CHECKSET,
			"servers.0.online_client_count":       CHECKSET,
			"servers.0.ike_config.#":              "1",
			"servers.0.ipsec_config.#":            "1",
		}
	}
	var fakeAlicloudVpnIpsecServersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudVpnIpsecServersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_vpn_ipsec_servers.default",
		existMapFunc: existAlicloudVpnIpsecServersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudVpnIpsecServersDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudVpnIpsecServersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, ipsecServerNameConf, vpnGatewayIdConf, nameRegexConf, allConf)
}

func testAccCheckAlicloudVpnIpsecServersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccIpsecServer-%d"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_zones.default.zones.0.id
}
locals {
  vswitch_id =  data.alicloud_vswitches.default.ids[0]
}

data "alicloud_vpn_gateways" "default" {
  vpc_id       = data.alicloud_vpcs.default.ids.0
  enable_ipsec = true
}

locals {
  vpn_gateway_id = data.alicloud_vpn_gateways.default.ids.0
}

resource "alicloud_vpn_ipsec_server" "default" {
  client_ip_pool     = "10.0.0.0/24"
  effect_immediately = true
  ipsec_server_name  = var.name
  local_subnet       = "192.168.0.0/24"
  psk_enabled        = true
  vpn_gateway_id     = local.vpn_gateway_id
}

data "alicloud_vpn_ipsec_servers" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
