package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSslVpnClientCertsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSslVpnClientCertsDataCfg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ssl_vpn_client_certs.foo"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_client_certs.foo", "ssl_vpn_client_cert_keys.0.ssl_vpn_client_cert_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_client_certs.foo", "ssl_vpn_client_cert_keys.0.ssl_vpn_server_id"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_client_certs.foo", "ssl_vpn_client_cert_keys.0.status", "normal"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_client_certs.foo", "ssl_vpn_client_cert_keys.0.create_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_client_certs.foo", "ssl_vpn_client_cert_keys.0.end_time"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_client_certs.foo", "ssl_vpn_client_cert_keys.0.name", "test_create_client_cert"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSslVpnClientCertsDataCfg = `
resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "testAccVpcConfig"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "testAccVswitchConfig"
}

resource "alicloud_vpn" "foo" {
        name = "testAccVpnConfig_create"
        vpc_id = "${alicloud_vpc.foo.id}"
		bandwidth = "10"
        enable_ssl = true
        instance_charge_type = "postpaid"
        auto_pay = true
		description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
    name = "testAccSslVpnServerConfig"
    vpn_gateway_id = "${alicloud_vpn.foo.id}"
    client_ip_pool = "192.168.0.0/16"
    local_subnet = "172.16.0.0/21"
    proto = "UDP"
    cipher = "AES-128-CBC"
    port = "1194"
    compress = "false"
}

resource "alicloud_ssl_vpn_client_cert" "foo" {
    ssl_vpn_server_id = "${alicloud_ssl_vpn_server.foo.id}"
    name = "test_create_client_cert"
}

data "alicloud_ssl_vpn_client_certs" "foo" {
	ssl_vpn_client_cert_id = "${alicloud_ssl_vpn_client_cert.foo.id}"
	output_file = "/tmp/vpnClientCerts"
}
`
