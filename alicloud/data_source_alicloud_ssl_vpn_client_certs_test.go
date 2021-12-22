package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSslVpnClientCertsDataSourceBasic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithTime(t, []int{1})
	}

	idsConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslVpnVpnClientCertsConfig(rand, map[string]string{
			"ids": `["${alicloud_ssl_vpn_client_cert.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSslVpnVpnClientCertsConfig(rand, map[string]string{
			"ids": `["${alicloud_ssl_vpn_client_cert.default.id}_fake"]`,
		}),
	}

	nameRegexConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslVpnVpnClientCertsConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ssl_vpn_client_cert.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSslVpnVpnClientCertsConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ssl_vpn_client_cert.default.name}_fake"`,
		}),
	}

	vpnServerIdConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslVpnVpnClientCertsConfig(rand, map[string]string{
			"ssl_vpn_server_id": `"${alicloud_ssl_vpn_client_cert.default.ssl_vpn_server_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSslVpnVpnClientCertsConfig(rand, map[string]string{
			"ssl_vpn_server_id": `"${alicloud_ssl_vpn_client_cert.default.ssl_vpn_server_id}_fake"`,
		}),
	}

	allServerIdConfig := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSslVpnVpnClientCertsConfig(rand, map[string]string{
			"ids":               `["${alicloud_ssl_vpn_client_cert.default.id}"]`,
			"name_regex":        `"${alicloud_ssl_vpn_client_cert.default.name}"`,
			"ssl_vpn_server_id": `"${alicloud_ssl_vpn_client_cert.default.ssl_vpn_server_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSslVpnVpnClientCertsConfig(rand, map[string]string{
			"ids":               `["${alicloud_ssl_vpn_client_cert.default.id}"]`,
			"name_regex":        `"${alicloud_ssl_vpn_client_cert.default.name}"`,
			"ssl_vpn_server_id": `"${alicloud_ssl_vpn_client_cert.default.ssl_vpn_server_id}_fake"`,
		}),
	}

	sslVpnClientCertsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConfig, nameRegexConfig,
		vpnServerIdConfig, allServerIdConfig)
}

func testAccCheckAlicloudSslVpnVpnClientCertsConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSslVpnClientCertsDataResource%d"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	vpc_name = "${var.name}"
}

resource "alicloud_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "default" {
	name="${var.name}"
	vpn_gateway_id="${alicloud_vpn_gateway.default.id}"
	client_ip_pool="192.168.1.0/24"
	local_subnet="172.16.1.0/24"
	protocol="UDP"
	port="1194"
	cipher="AES-128-CBC"
	compress="false"
}

resource "alicloud_ssl_vpn_client_cert" "default" {
	name="${var.name}"
	ssl_vpn_server_id="${alicloud_ssl_vpn_server.default.id}"
}

data "alicloud_ssl_vpn_client_certs" "default" {
	%s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existSslVpnClientCertsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                     "1",
		"names.#":                   "1",
		"certs.#":                   "1",
		"certs.0.name":              fmt.Sprintf("tf-testAccSslVpnClientCertsDataResource%d", rand),
		"certs.0.ssl_vpn_server_id": CHECKSET,
		"certs.0.create_time":       CHECKSET,
		"certs.0.id":                CHECKSET,
		"certs.0.end_time":          CHECKSET,
		"certs.0.status":            "normal",
	}
}

var fakeSslVpnClientCertsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"names.#": "0",
		"certs.#": "0",
	}
}

var sslVpnClientCertsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ssl_vpn_client_certs.default",
	existMapFunc: existSslVpnClientCertsMapFunc,
	fakeMapFunc:  fakeSslVpnClientCertsMapFunc,
}
