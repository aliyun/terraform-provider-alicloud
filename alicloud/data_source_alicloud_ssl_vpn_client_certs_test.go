package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSslVpnClientCertsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSslVpnClientCertsDataCfg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ssl_vpn_client_certs.foo"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_client_certs.foo", "ssl_vpn_client_certs.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_client_certs.foo", "ssl_vpn_client_certs.0.name", "tf-testAccSslVpnClientCertsDataResource"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_client_certs.foo", "ssl_vpn_client_certs.0.ssl_vpn_server_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_client_certs.foo", "ssl_vpn_client_certs.0.create_time"),
				),
			},
		},
	})
}

func TestAccAlicloudSslVpnClientCertsDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSslVpnClientCertsDataEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ssl_vpn_client_certs.foo"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_client_certs.foo", "ssl_vpn_client_certs.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_client_certs.foo", "ssl_vpn_client_certs.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_client_certs.foo", "ssl_vpn_client_certs.0.ssl_vpn_server_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_ssl_vpn_client_certs.foo", "ssl_vpn_client_certs.0.create_time"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSslVpnClientCertsDataCfg = `
variable "name" {
	default = "tf-testAccSslVpnClientCertsDataResource"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

resource "alicloud_vpn_gateway" "foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.foo.id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
	name="${var.name}"
	vpn_gateway_id="${alicloud_vpn_gateway.foo.id}"
	client_ip_pool="192.168.1.0/24"
	local_subnet="172.16.1.0/24"
	protocol="UDP"
	port="1194"
	cipher="AES-128-CBC"
	compress="false"
}

resource "alicloud_ssl_vpn_client_cert" "foo" {
	name="${var.name}"
	ssl_vpn_server_id="${alicloud_ssl_vpn_server.foo.id}"
}

data "alicloud_ssl_vpn_client_certs" "foo" {
	name_regex = "tf-testAcc*"
	ids = ["${alicloud_ssl_vpn_client_cert.foo.id}"]
}
`

const testAccCheckAlicloudSslVpnClientCertsDataEmpty = `
data "alicloud_ssl_vpn_client_certs" "foo" {
	name_regex = "^tf-testacc-fake-name"
}
`
