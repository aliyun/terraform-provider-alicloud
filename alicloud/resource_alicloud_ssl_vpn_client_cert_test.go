package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudSslVpnClientCert_basic(t *testing.T) {
	var sslVpnClientCert vpc.DescribeSslVpnClientCertResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ssl_vpn_client_cert.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSslVpnClientCertDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSslVpnClientCertConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslVpnClientCertExists("alicloud_ssl_vpn_client_cert.foo", &sslVpnClientCert),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_client_cert.foo", "name", "test_create_client_cert"),
					resource.TestCheckResourceAttrSet(
						"alicloud_ssl_vpn_client_cert.foo", "ssl_vpn_server_id"),
				),
			},
		},
	})

}

func TestAccAlicloudSslVpnClientCert_update(t *testing.T) {
	var sslVpnClientCert vpc.DescribeSslVpnClientCertResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSslVpnClientCertDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSslVpnClientCertConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslVpnClientCertExists("alicloud_ssl_vpn_client_cert.foo", &sslVpnClientCert),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_client_cert.foo", "name", "test_create_client_cert"),
					resource.TestCheckResourceAttrSet(
						"alicloud_ssl_vpn_client_cert.foo", "ssl_vpn_server_id"),
				),
			},
			resource.TestStep{
				Config: testAccSslVpnClientCertConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslVpnClientCertExists("alicloud_ssl_vpn_client_cert.foo", &sslVpnClientCert),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_client_cert.foo", "name", "test_update_client_cert"),
					resource.TestCheckResourceAttrSet(
						"alicloud_ssl_vpn_client_cert.foo", "ssl_vpn_server_id"),
				),
			},
		},
	})
}

func testAccCheckSslVpnClientCertExists(n string, vpn *vpc.DescribeSslVpnClientCertResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribeSslVpnClientCert(rs.Primary.ID)

		if err != nil {
			return err
		}

		*vpn = instance
		return nil
	}
}

func testAccCheckSslVpnClientCertDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ssl_vpn_client_cert" {
			continue
		}

		// Try to find the VPN
		instance, err := client.DescribeSslVpnClientCert(rs.Primary.ID)

		if err != nil {
			if IsExceptedError(err, SslVpnClientCertNofFound) {
				continue
			}
			return err
		}

		if instance.SslVpnClientCertId != "" {
			return fmt.Errorf("VPN %s still exist", instance.SslVpnClientCertId)
		}
	}

	return nil
}

const testAccSslVpnClientCertConfig = `
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
`

const testAccSslVpnClientCertConfigUpdate = `
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
        name = "testAccVpnConfig_update"
        vpc_id = "${alicloud_vpc.foo.id}"
		bandwidth = "10"
        enable_ssl = true
        instance_charge_type = "postpaid"
        auto_pay = true
		description = "test_update_description"
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
    name = "test_update_client_cert"
}
`
