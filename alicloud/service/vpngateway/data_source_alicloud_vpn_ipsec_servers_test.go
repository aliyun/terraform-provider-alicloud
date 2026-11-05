package vpngateway_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudVPNIpsecServersDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpn_ipsec_server.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpn_ipsec_server.default.id}_fake"]`,
		}),
	}
	ipsecServerNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_vpn_ipsec_server.default.id}"]`,
			"ipsec_server_name": `"${alicloud_vpn_ipsec_server.default.ipsec_server_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_vpn_ipsec_server.default.id}"]`,
			"ipsec_server_name": `"${alicloud_vpn_ipsec_server.default.ipsec_server_name}_fake"`,
		}),
	}
	vpnGatewayIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_vpn_ipsec_server.default.id}"]`,
			"vpn_gateway_id": `"${alicloud_vpn_ipsec_server.default.vpn_gateway_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_vpn_ipsec_server.default.id}"]`,
			"vpn_gateway_id": `"${alicloud_vpn_ipsec_server.default.vpn_gateway_id}_fake"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_ipsec_server.default.ipsec_server_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vpn_ipsec_server.default.ipsec_server_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_vpn_ipsec_server.default.id}"]`,
			"ipsec_server_name": `"${alicloud_vpn_ipsec_server.default.ipsec_server_name}"`,
			"name_regex":        `"${alicloud_vpn_ipsec_server.default.ipsec_server_name}"`,
			"vpn_gateway_id":    `"${alicloud_vpn_ipsec_server.default.vpn_gateway_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpnIpsecServersDataSourceName(rand, map[string]string{
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
			"servers.0.psk":                       acctest.CHECKSET,
			"servers.0.psk_enabled":               "true",
			"servers.0.vpn_gateway_id":            acctest.CHECKSET,
			"servers.0.internet_ip":               acctest.CHECKSET,
			"servers.0.id":                        acctest.CHECKSET,
			"servers.0.idaas_instance_id":         "",
			"servers.0.create_time":               acctest.CHECKSET,
			"servers.0.ipsec_server_id":           acctest.CHECKSET,
			"servers.0.max_connections":           acctest.CHECKSET,
			"servers.0.multi_factor_auth_enabled": acctest.CHECKSET,
			"servers.0.online_client_count":       acctest.CHECKSET,
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
	var alicloudVpnIpsecServersCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_vpn_ipsec_servers.default",
		ExistMapFunc: existAlicloudVpnIpsecServersDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudVpnIpsecServersDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudVpnIpsecServersCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, ipsecServerNameConf, vpnGatewayIdConf, nameRegexConf, allConf)
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

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  name       = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/21"
  availability_zone = data.alicloud_zones.default.zones.0.id
  name              = var.name
}

resource "alicloud_vswitch" "default1" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.17.0.0/22"
  availability_zone = data.alicloud_zones.default.zones.1.id
  name              = var.name
}

resource "alicloud_vpn_gateway" "default" {
  name                         = var.name
  instance_charge_type         = "PostPaid"
  vpc_id                       = alicloud_vpc.default.id
  vswitch_id                   = alicloud_vswitch.default.id
  disaster_recovery_vswitch_id = alicloud_vswitch.default1.id
  bandwidth                    = "10"
  description                  = "test"
}

locals {
  vpn_gateway_id = alicloud_vpn_gateway.default.id
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
