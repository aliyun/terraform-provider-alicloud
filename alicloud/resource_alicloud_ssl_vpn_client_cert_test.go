package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudSslVpnClientCert_basic(t *testing.T) {
	var sslVpnClientCert vpc.DescribeSslVpnClientCertResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},

		// module name
		IDRefreshName: "alicloud_ssl_vpn_client_cert.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSslVpnClientCertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSslVpnClientCertConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslVpnClientCertExists("alicloud_ssl_vpn_client_cert.foo", &sslVpnClientCert),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_client_cert.foo", "name", "tf-testAccSslVpnClientCertConfig"),
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
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSslVpnClientCertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSslVpnClientCertConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslVpnClientCertExists("alicloud_ssl_vpn_client_cert.foo", &sslVpnClientCert),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_client_cert.foo", "name", "tf-testAccSslVpnClientCertConfig"),
					resource.TestCheckResourceAttrSet(
						"alicloud_ssl_vpn_client_cert.foo", "ssl_vpn_server_id"),
				),
			},
			{
				Config: testAccSslVpnClientCertConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSslVpnClientCertExists("alicloud_ssl_vpn_client_cert.foo", &sslVpnClientCert),
					resource.TestCheckResourceAttr(
						"alicloud_ssl_vpn_client_cert.foo", "name", "tf-testAccSslVpnClientCertUpdate"),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpnGatewayService := VpnGatewayService{client}
		instance, err := vpnGatewayService.DescribeSslVpnClientCert(rs.Primary.ID)

		if err != nil {
			return err
		}

		*vpn = instance
		return nil
	}
}

func testAccCheckSslVpnClientCertDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ssl_vpn_client_cert" {
			continue
		}

		instance, err := vpnGatewayService.DescribeSslVpnClientCert(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if instance.SslVpnClientCertId != "" {
			return fmt.Errorf("Ssl VPN client cert %s still exist", instance.SslVpnClientCertId)
		}
	}

	return nil
}

const testAccSslVpnClientCertConfig = `
variable "name" {
	default = "tf-testAccSslVpnClientCertConfig"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
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

resource "alicloud_vpn_gateway" "foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.foo.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
	name = "${var.name}"
	vpn_gateway_id = "${alicloud_vpn_gateway.foo.id}"
	client_ip_pool = "192.168.0.0/16"
	local_subnet = "172.16.0.0/21"
	protocol = "UDP"
	cipher = "AES-128-CBC"
	port = "1194"
	compress = "false"
}

resource "alicloud_ssl_vpn_client_cert" "foo" {
	ssl_vpn_server_id = "${alicloud_ssl_vpn_server.foo.id}"
	name = "${var.name}"
}
`

const testAccSslVpnClientCertConfigUpdate = `
variable "name" {
	default = "tf-testAccSslVpnClientCertUpdate"
}

resource "alicloud_vpc" "foo" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
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

resource "alicloud_vpn_gateway" "foo" {
	name = "${var.name}"
	vpc_id = "${alicloud_vswitch.foo.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_update_description"
}

resource "alicloud_ssl_vpn_server" "foo" {
    name = "${var.name}"
    vpn_gateway_id = "${alicloud_vpn_gateway.foo.id}"
    client_ip_pool = "192.168.0.0/16"
    local_subnet = "172.16.0.0/21"
    protocol = "UDP"
    cipher = "AES-128-CBC"
    port = "1194"
    compress = "false"
}

resource "alicloud_ssl_vpn_client_cert" "foo" {
    ssl_vpn_server_id = "${alicloud_ssl_vpn_server.foo.id}"
    name = "${var.name}"
}
`
