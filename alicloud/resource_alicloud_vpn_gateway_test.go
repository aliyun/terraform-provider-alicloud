package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudVpnGateway_basic(t *testing.T) {
	var vpn vpc.DescribeVpnGatewayResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_vpn_gateway.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckVpnGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpnConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayExists("alicloud_vpn_gateway.foo", &vpn),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "name", "testAccVpnConfig_create"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpn_gateway.foo", "vpc_id"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "specification", "10M"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "enable_ssl", "true"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "enable_ipsec", "true"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "description", "test_create_description"),
				),
			},
		},
	})

}

func TestAccAlicloudVpnGateway_update(t *testing.T) {
	var vpn vpc.DescribeVpnGatewayResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpnGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpnConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayExists("alicloud_vpn_gateway.foo", &vpn),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "name", "testAccVpnConfig_create"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpn_gateway.foo", "vpc_id"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "specification", "10M"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "enable_ssl", "true"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "enable_ipsec", "true"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "description", "test_create_description"),
				),
			},
			resource.TestStep{
				Config: testAccVpnConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpnGatewayExists("alicloud_vpn_gateway.foo", &vpn),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "name", "testAccVpnConfig_update"),
					resource.TestCheckResourceAttrSet(
						"alicloud_vpn_gateway.foo", "vpc_id"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "specification", "10M"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "enable_ssl", "true"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "enable_ipsec", "true"),
					resource.TestCheckResourceAttr(
						"alicloud_vpn_gateway.foo", "description", "test_update_description"),
				),
			},
		},
	})
}

func testAccCheckVpnGatewayExists(n string, vpn *vpc.DescribeVpnGatewayResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPN ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribeVpnGateway(rs.Primary.ID)

		if err != nil {
			return err
		}

		*vpn = instance
		return nil
	}
}

func testAccCheckVpnGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vpn" {
			continue
		}

		// Try to find the VPN
		instance, err := client.DescribeVpnGateway(rs.Primary.ID)

		if err != nil {
			//if IsExceptedError(err, VpnNotFound) || IsExceptedError(err, InstanceNotFound) {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Describe Vpn error %#v", err)
		}

		if instance.VpnGatewayId != "" {
			return fmt.Errorf("VPN %s still exist", instance.VpnGatewayId)
		}
	}

	return nil
}

const testAccVpnConfig = `
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
        instance_charge_type = "postpaid"
		description = "test_create_description"
}
`

const testAccVpnConfigUpdate = `
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
        instance_charge_type = "postpaid"
		description = "test_update_description"
}
`
