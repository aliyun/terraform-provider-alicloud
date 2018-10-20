package alicloud

import (
	"fmt"
	"testing"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudVpnConnection_basic(t *testing.T) {
	var vpnConn vpc.DescribeVpnConnectionResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_vpn_connection.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpnConnectionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpnConnConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnConnectionExists("alicloud_vpn_connection.foo", &vpnConn),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "name", "tf-testAccVpnConnConfig_test1"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpn_connection.foo", "vpn_gateway_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpn_connection.foo", "customer_gateway_id"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_auth_alg", "md5"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_enc_alg", "des"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_version", "ikev2"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_mode", "main"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_lifetime", "86400"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.psk", "tf-testvpn2"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_pfs", "group1"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_remote_id", "testbob2"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_local_id", "testalice2"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ipsec_config.0.ipsec_pfs", "group5"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ipsec_config.0.ipsec_enc_alg", "des"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ipsec_config.0.ipsec_auth_alg", "md5"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ipsec_config.0.ipsec_lifetime", "8640"),

					testAccCheckVpnConnectionAttr("alicloud_vpn_connection.foo", &vpnConn, "config"),
				),
			},
		},
	})

}

func TestAccAlicloudVpnConnection_update(t *testing.T) {
	var vpnConn vpc.DescribeVpnConnectionResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpnConnectionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpnConnConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnConnectionExists("alicloud_vpn_connection.foo", &vpnConn),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "name", "tf-testAccVpnConnConfig_test1"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpn_connection.foo", "vpn_gateway_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpn_connection.foo", "customer_gateway_id"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_auth_alg", "md5"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_enc_alg", "des"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_version", "ikev2"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_mode", "main"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_lifetime", "86400"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.psk", "tf-testvpn2"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_pfs", "group1"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_remote_id", "testbob2"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_local_id", "testalice2"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ipsec_config.0.ipsec_pfs", "group5"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ipsec_config.0.ipsec_enc_alg", "des"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ipsec_config.0.ipsec_auth_alg", "md5"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ipsec_config.0.ipsec_lifetime", "8640"),
					testAccCheckVpnConnectionAttr("alicloud_vpn_connection.foo", &vpnConn, "config"),
				),
			},

			resource.TestStep{
				Config: testAccVpnConnConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnConnectionExists("alicloud_vpn_connection.foo", &vpnConn),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "name", "tf-testAccVpnConnConfig_test2"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpn_connection.foo", "vpn_gateway_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpn_connection.foo", "customer_gateway_id"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_auth_alg", "sha1"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_enc_alg", "3des"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_version", "ikev2"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_mode", "aggressive"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_lifetime", "8640"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.psk", "tf-testvpn1"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_pfs", "group2"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_remote_id", "testbob1"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config.0.ike_local_id", "testalice1"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ipsec_config.0.ipsec_pfs", "group2"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ipsec_config.0.ipsec_enc_alg", "aes"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ipsec_config.0.ipsec_auth_alg", "sha1"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ipsec_config.0.ipsec_lifetime", "86400"),
					testAccCheckVpnConnectionAttr("alicloud_vpn_connection.foo", &vpnConn, "update"),
				),
			},
		},
	})
}

const cfg_local_subnet = "172.16.0.0/24,172.16.1.0/24"
const cfg_remote_subnet = "10.0.0.0/24,10.0.1.0/24"
const update_local_subnet = "172.16.1.0/24,172.16.2.0/24"
const update_remote_subnet = "10.4.0.0/24,10.0.3.0/24"

func compareSubnet(astr string, bstr string) bool {
	aarry := strings.Split(astr, ",")
	barry := strings.Split(bstr, ",")
	if len(aarry) != len(barry) {
		return false
	}

	for _, item := range aarry {
		if !strings.Contains(bstr, item) {
			return false
		}
	}
	return true
}

func testAccCheckVpnConnectionAttr(n string, vpnConn *vpc.DescribeVpnConnectionResponse, step string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var local_subnet string
		var remote_subnet string
		if step == "config" {
			local_subnet = cfg_local_subnet
			remote_subnet = cfg_remote_subnet
		} else {
			local_subnet = update_local_subnet
			remote_subnet = update_remote_subnet
		}

		if !compareSubnet(vpnConn.LocalSubnet, local_subnet) {
			return fmt.Errorf("wrong local subnet, expect %s, get %s", local_subnet, vpnConn.LocalSubnet)
		}

		if !compareSubnet(vpnConn.RemoteSubnet, remote_subnet) {
			return fmt.Errorf("wrong remote subnet, expect %s, get %s", remote_subnet, vpnConn.RemoteSubnet)
		}

		return nil
	}
}

func testAccCheckVpnConnectionExists(n string, vpn *vpc.DescribeVpnConnectionResponse) resource.TestCheckFunc {
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
		instance, err := vpnGatewayService.DescribeVpnConnection(rs.Primary.ID)

		if err != nil {
			return err
		}

		*vpn = instance
		return nil
	}
}

func testAccCheckVpnConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vpn_connection" {
			continue
		}

		instance, err := vpnGatewayService.DescribeVpnConnection(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		return fmt.Errorf("VPN connection %s still exist", instance.VpnConnectionId)
	}
	return nil
}

const testAccVpnConnConfig = `
variable "name" {
	default = "tf-testAccVpnConnConfig_test1"
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
	vpc_id = "${alicloud_vpc.foo.id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_create_description"
}

resource "alicloud_vpn_customer_gateway" "foo" {
	name = "${var.name}"
	ip_address = "42.104.22.228"
	description = "testAccVpnCgwDesc"
}

resource "alicloud_vpn_connection" "foo" {
	name = "${var.name}"
	vpn_gateway_id = "${alicloud_vpn_gateway.foo.id}"
	customer_gateway_id = "${alicloud_vpn_customer_gateway.foo.id}"
	local_subnet = ["172.16.0.0/24", "172.16.1.0/24"]
	remote_subnet = ["10.0.0.0/24", "10.0.1.0/24"]
	effect_immediately = true
	ike_config = [{
        ike_auth_alg = "md5"
        ike_enc_alg = "des"
        ike_version = "ikev2"
        ike_mode = "main"
        ike_lifetime = 86400
        psk = "tf-testvpn2"
        ike_pfs = "group1"
        ike_remote_id = "testbob2"
        ike_local_id = "testalice2"
        }
    ]
	ipsec_config = [{
        ipsec_pfs = "group5"
        ipsec_enc_alg = "des"
        ipsec_auth_alg = "md5"
        ipsec_lifetime = 8640
    }]
}
`

const testAccVpnConnConfigUpdate = `
variable "name" {
	default = "tf-testAccVpnConnConfig_test2"
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
	vpc_id = "${alicloud_vpc.foo.id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "test_update_description"
}

resource "alicloud_vpn_customer_gateway" "foo" {
	name = "${var.name}"
	ip_address = "42.104.22.228"
	description = "testAccVpnCgwDesc"
}

resource "alicloud_vpn_connection" "foo" {
	name = "${var.name}"
	vpn_gateway_id = "${alicloud_vpn_gateway.foo.id}"
	customer_gateway_id = "${alicloud_vpn_customer_gateway.foo.id}"
	local_subnet = ["172.16.1.0/24", "172.16.2.0/24"]
	remote_subnet = ["10.4.0.0/24", "10.0.3.0/24"]
	effect_immediately = true
	ike_config = [{
        ike_auth_alg = "sha1"
        ike_enc_alg = "3des"
        ike_version = "ikev2"
        ike_mode = "aggressive"
        ike_lifetime = 8640
        psk = "tf-testvpn1"
        ike_pfs = "group2"
        ike_remote_id = "testbob1"
        ike_local_id = "testalice1"
        }
    ]
	ipsec_config = [{
        ipsec_pfs = "group2"
        ipsec_enc_alg = "aes"
        ipsec_auth_alg = "sha1"
        ipsec_lifetime = 86400
    }]
}
`
