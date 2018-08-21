package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
						"alicloud_vpn_connection.foo", "name", "tf-vco_test1"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpn_connection.foo", "vpn_gateway_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpn_connection.foo", "customer_gateway_id"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "local_subnet", "172.16.0.0/16"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "remote_subnet", "10.0.0.0/8"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config", create_ike_config),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ipsec_config", create_ipsec_config),
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
						"alicloud_vpn_connection.foo", "name", "tf-vco_test1"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpn_connection.foo", "vpn_gateway_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpn_connection.foo", "customer_gateway_id"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "local_subnet", "172.16.0.0/16"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "remote_subnet", "10.0.0.0/8"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config", create_ike_config),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ipsec_config", create_ipsec_config),
				),
			},

			resource.TestStep{
				Config: testAccVpnConnConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnConnectionExists("alicloud_vpn_connection.foo", &vpnConn),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "name", "tf-vco_test1"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpn_connection.foo", "vpn_gateway_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpn_connection.foo", "customer_gateway_id"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "local_subnet", "172.16.1.0/24"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "remote_subnet", "192.16.0.0/16"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ike_config", update_ike_config),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_connection.foo", "ipsec_config", update_ipsec_config),
				),
			},
		},
	})
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

		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribeVpnConnection(rs.Primary.ID)

		if err != nil {
			return err
		}

		*vpn = instance
		return nil
	}
}

func testAccCheckVpnConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vpn_connection" {
			continue
		}

		instance, err := client.DescribeVpnConnection(rs.Primary.ID)

		if err != nil {
			if IsExceptedError(err, VpnConnNotFound) {
				continue
			}
			return err
		}

		if instance.VpnConnectionId != "" {
			return fmt.Errorf("VPN connection %s still exist", instance.VpnConnectionId)
		}
	}

	return nil
}

const create_ike_config string = "{\"IkeAuthAlg\":\"sha1\",\"IkeEncAlg\":\"aes\",\"IkeVersion\":\"ikev2\",\"IkeMode\":\"main\",\"IkeLifetime\":3600,\"Psk\":\"tf-testvpn1\",\"IkePfs\":\"group2\"}"
const create_ipsec_config string = "{\"IpsecPfs\":\"group2\",\"IpsecEncAlg\":\"aes\",\"IpsecAuthAlg\":\"sha1\",\"IpsecLifetime\":7200}"
const update_ike_config string = "{\"IkeAuthAlg\":\"sha1\",\"IkeEncAlg\":\"aes\",\"IkeVersion\":\"ikev2\",\"IkeMode\":\"aggressive\",\"IkeLifetime\":86400,\"Psk\":\"tf-testvpn1\",\"IkePfs\":\"group2\"}"
const update_ipsec_config string = "{\"IpsecPfs\":\"group2\",\"IpsecEncAlg\":\"aes\",\"IpsecAuthAlg\":\"sha1\",\"IpsecLifetime\":86400}"

const testAccVpnConnConfig = `
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

resource "alicloud_vpn_customer_gateway" "foo" {
	name = "testAccVpnCgwName"
	ip_address = "42.104.22.228"
	description = "testAccVpnCgwDesc"
}

resource "alicloud_vpn_connection" "foo" {
	name = "tf-vco_test1"
	vpn_gateway_id = "${alicloud_vpn.foo.id}"
	customer_gateway_id = "${alicloud_vpn_customer_gateway.foo.id}"
	local_subnet = "172.16.0.0/16"
	remote_subnet = "10.0.0.0/8"
	effect_immediately = true
	ike_config = "{\"IkeAuthAlg\":\"sha1\",\"IkeEncAlg\":\"aes\",\"IkeVersion\":\"ikev2\",\"IkeMode\":\"main\",\"IkeLifetime\":3600,\"Psk\":\"tf-testvpn1\",\"IkePfs\":\"group2\"}"
	ipsec_config = "{\"IpsecPfs\":\"group2\",\"IpsecEncAlg\":\"aes\",\"IpsecAuthAlg\":\"sha1\",\"IpsecLifetime\":7200}"
}
`

const testAccVpnConnConfigUpdate = `
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

resource "alicloud_vpn_customer_gateway" "foo" {
	name = "testAccVpnCgwName"
	ip_address = "42.104.22.228"
	description = "testAccVpnCgwDesc"
}

resource "alicloud_vpn_connection" "foo" {
	name = "tf-vco_test1"
	vpn_gateway_id = "${alicloud_vpn.foo.id}"
	customer_gateway_id = "${alicloud_vpn_customer_gateway.foo.id}"
	local_subnet = "172.16.1.0/24"
	remote_subnet = "192.16.0.0/16"
	effect_immediately = true
	ike_config = "{\"IkeAuthAlg\":\"sha1\",\"IkeEncAlg\":\"aes\",\"IkeVersion\":\"ikev2\",\"IkeMode\":\"aggressive\",\"IkeLifetime\":86400,\"Psk\":\"tf-testvpn1\",\"IkePfs\":\"group2\"}"
	ipsec_config = "{\"IpsecPfs\":\"group2\",\"IpsecEncAlg\":\"aes\",\"IpsecAuthAlg\":\"sha1\",\"IpsecLifetime\":86400}"
}
`
