package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudSslVpnServersDataSourceBasic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	PreCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithAccountSiteType(t, IntlSite)
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslVpnServerConfig(rand, map[string]string{
			"ids": `[ alicloud_ssl_vpn_server.default.id ]`,
		}),
		fakeConfig: testAccCheckAlicloudSslVpnServerConfig(rand, map[string]string{
			"ids": `[ "${alicloud_ssl_vpn_server.default.id}_fake" ]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslVpnServerConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccSslVpnServersDataResource%d"`, rand),
		}),
		fakeConfig: testAccCheckAlicloudSslVpnServerConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccSslVpnServersDataResource%d_fake"`, rand),
		}),
	}

	vpnGatewayIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslVpnServerConfig(rand, map[string]string{
			"vpn_gateway_id": `alicloud_vpn_gateway.default.id`,
		}),
		fakeConfig: testAccCheckAlicloudSslVpnServerConfig(rand, map[string]string{
			"vpn_gateway_id": `"${alicloud_vpn_gateway.default.id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslVpnServerConfig(rand, map[string]string{
			"ids":            `[ alicloud_ssl_vpn_server.default.id ]`,
			"name_regex":     fmt.Sprintf(`"tf-testAccSslVpnServersDataResource%d"`, rand),
			"vpn_gateway_id": `alicloud_vpn_gateway.default.id`,
		}),
		fakeConfig: testAccCheckAlicloudSslVpnServerConfig(rand, map[string]string{
			"ids":            `[ alicloud_ssl_vpn_server.default.id ]`,
			"name_regex":     fmt.Sprintf(`"tf-testAccSslVpnServersDataResource%d"`, rand),
			"vpn_gateway_id": `"${alicloud_vpn_gateway.default.id}_fake"`,
		}),
	}

	sslVpnServerCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, PreCheck, idsConf, nameRegexConf, vpnGatewayIdConf, allConf)
}

func testAccCheckAlicloudSslVpnServerConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslVpnServersDataResource%d"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = var.name
}

resource "alicloud_vswitch" "default" {
	vpc_id = alicloud_vpc.default.id
	cidr_block = "172.16.0.0/21"
	availability_zone = data.alicloud_zones.default.zones.0.id
	name = var.name
}

resource "alicloud_vpn_gateway" "default" {
	name = var.name
	vpc_id = alicloud_vpc.default.id
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
}

resource "alicloud_ssl_vpn_server" "default" {
	name=var.name
	vpn_gateway_id=alicloud_vpn_gateway.default.id
	client_ip_pool="192.168.1.0/24"
	local_subnet="172.16.1.0/24"
	protocol="UDP"
	port="1194"
	cipher="AES-128-CBC"
	compress="false"
}

data "alicloud_ssl_vpn_servers" "default" {
	%s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existSslVpnServersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"servers.#":                 "1",
		"ids.#":                     "1",
		"names.#":                   "1",
		"servers.0.name":            fmt.Sprintf("tf-testAccSslVpnServersDataResource%d", rand),
		"servers.0.vpn_gateway_id":  CHECKSET,
		"servers.0.id":              CHECKSET,
		"servers.0.create_time":     CHECKSET,
		"servers.0.compress":        "false",
		"servers.0.cipher":          "AES-128-CBC",
		"servers.0.port":            "1194",
		"servers.0.proto":           "UDP",
		"servers.0.local_subnet":    "172.16.1.0/24",
		"servers.0.client_ip_pool":  "192.168.1.0/24",
		"servers.0.internet_ip":     CHECKSET,
		"servers.0.max_connections": "5",
		"servers.0.connections":     "0",
	}
}

var fakeSslVpnServersMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"servers.#": "0",
		"ids.#":     "0",
		"names.#":   "0",
	}
}

var sslVpnServerCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ssl_vpn_servers.default",
	existMapFunc: existSslVpnServersMapFunc,
	fakeMapFunc:  fakeSslVpnServersMapFunc,
}
