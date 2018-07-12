package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudNatGateway_specification(t *testing.T) {
	var nat vpc.NatGateway

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_nat_gateway.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckNatGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNatGatewayConfigSpec,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists(
						"alicloud_nat_gateway.foo", &nat),
					resource.TestCheckResourceAttr(
						"alicloud_nat_gateway.foo",
						"name",
						"testAccNatGatewayConfigSpec"),
					resource.TestCheckResourceAttr(
						"alicloud_nat_gateway.foo",
						"specification",
						"Small"),
				),
			},

			resource.TestStep{
				Config: testAccNatGatewayConfigSpecUpgrade,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayExists(
						"alicloud_nat_gateway.foo", &nat),
					resource.TestCheckResourceAttr(
						"alicloud_nat_gateway.foo",
						"specification",
						"Middle"),
				),
			},
		},
	})
}

func testAccCheckNatGatewayExists(n string, nat *vpc.NatGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Gateway ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribeNatGateway(rs.Primary.ID)

		if err != nil {
			return err
		}
		if instance.NatGatewayId != rs.Primary.ID {
			return fmt.Errorf("Nat gateway not found")
		}

		*nat = instance
		return nil
	}
}

func testAccCheckNatGatewayDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_nat_gateway" {
			continue
		}

		// Try to find the Nat gateway
		if _, err := client.DescribeNatGateway(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		return fmt.Errorf("Nat gateway %s still exist", rs.Primary.ID)
	}

	return nil
}

const testAccNatGatewayConfigSpec = `
variable "name" {
	default = "testAccNatGatewayConfigSpec"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Small"
	name = "${var.name}"
}
`

const testAccNatGatewayConfigSpecUpgrade = `
variable "name" {
	default = "testAccNatGatewayConfigSpec"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	specification = "Middle"
	name = "${var.name}"
}
`
