package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudSslVpnServer_basic(t *testing.T) {
	var sslVpnServer vpc.SslVpnServer

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_ssl_vpn_server.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSslVpnServerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSslVpnServerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslVpnServerExists("alicloud_ssl_vpn_server.foo", &sslVpnServer),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "name", "testAccSslVpnServerConfig_create"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "client_ip_pool", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "local_subnet", "172.16.0.0/21"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "protocol", "UDP"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "port", "1194"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "cipher", "AES-128-CBC"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "compress", "false"),
				),
			},
		},
	})

}

func TestAccAlicloudSslVpnServer_update(t *testing.T) {
	var sslVpnServer vpc.SslVpnServer

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSslVpnServerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSslVpnServerConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslVpnServerExists("alicloud_ssl_vpn_server.foo", &sslVpnServer),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "name", "testAccSslVpnServerConfig_create"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "client_ip_pool", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "local_subnet", "172.16.0.0/21"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "protocol", "UDP"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "port", "1194"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "cipher", "AES-128-CBC"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "compress", "false"),
				),
			},
			resource.TestStep{
				Config: testAccSslVpnServerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslVpnServerExists("alicloud_ssl_vpn_server.foo", &sslVpnServer),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "name", "testAccSslVpnServerConfig_update"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "client_ip_pool", "192.168.10.0/24"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "local_subnet", "172.16.0.0/24"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "protocol", "UDP"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "port", "1194"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "cipher", "none"),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_server.foo", "compress", "false"),
				),
			},
		},
	})
}

func testAccCheckSslVpnServerExists(n string, vpn *vpc.SslVpnServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribeSslVpnServer(rs.Primary.ID)

		if err != nil {
			return err
		}

		*vpn = instance
		return nil
	}
}

func testAccCheckSslVpnServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ssl_vpn_server" {
			continue
		}

		instance, err := client.DescribeSslVpnServer(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		return fmt.Errorf("SSL VPN server %s still exist", instance.SslVpnServerId)
	}

	return nil
}

const testAccSslVpnServerConfig = `
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

resource "alicloud_vpn_gateway" "foo" {
	name = "testAccVpnConfig_create"
	vpc_id = "${alicloud_vpc.foo.id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
	name = "testAccSslVpnServerConfig_create"
	vpn_gateway_id = "${alicloud_vpn_gateway.foo.id}"
	client_ip_pool = "192.168.0.0/16"
	local_subnet = "172.16.0.0/21"
	protocol = "UDP"
	cipher = "AES-128-CBC"
	port = 1194
	compress = "false"
}
`

const testAccSslVpnServerConfigUpdate = `
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

resource "alicloud_vpn_gateway" "foo" {
	name = "testAccVpnConfig_update"
	vpc_id = "${alicloud_vpc.foo.id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_update_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
    name = "testAccSslVpnServerConfig_update"
    vpn_gateway_id = "${alicloud_vpn_gateway.foo.id}"
    client_ip_pool = "192.168.10.0/24"
    local_subnet = "172.16.0.0/24"
    protocol = "UDP"
    cipher = "none"
    port = 1194
    compress = "false"
}
`
