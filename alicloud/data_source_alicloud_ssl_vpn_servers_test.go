package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudSslVpnServersDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSslVpnServersDataCfg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ssl_vpn_servers.foo"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "ssl_vpn_servers.0.name", "testAccSslVpnServerConfig_create"),
					resource.TestCheckResourceAttrSet("data.alicloud_ssl_vpn_servers.foo", "ssl_vpn_servers.0.vpn_gateway_id"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "ssl_vpn_servers.0.client_ip_pool", "192.168.10.0/24"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "ssl_vpn_servers.0.local_subnet", "172.16.0.0/24"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "ssl_vpn_servers.0.proto", "UDP"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "ssl_vpn_servers.0.cipher", "AES-192-CBC"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "ssl_vpn_servers.0.port", "1194"),
					resource.TestCheckResourceAttr("data.alicloud_ssl_vpn_servers.foo", "ssl_vpn_servers.0.compress", "false"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSslVpnServersDataCfg = `
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
    name = "testAccSslVpnServerConfig_create"
    vpn_gateway_id = "${alicloud_vpn.foo.id}"
    client_ip_pool = "192.168.10.0/24"
    local_subnet = "172.16.0.0/24"
    proto = "UDP"
    cipher = "AES-192-CBC"
    port = "1194"
    compress = "false"
}

data "alicloud_ssl_vpn_servers" "foo" {
	ssl_vpn_server_id = "${alicloud_ssl_vpn_server.foo.id}"
	output_file = "/tmp/sslVpnServers"
}


`
