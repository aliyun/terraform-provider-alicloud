package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudRouterInterfaceConnection_basic(t *testing.T) {
	var vpcInstance vpc.DescribeVpcAttributeResponse
	var ri, oppoRI vpc.RouterInterfaceType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: "alicloud_router_interface_connection.foo",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouterInterfaceConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouterInterfaceConnectionConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("alicloud_vpc.foo", &vpcInstance),
					testAccCheckRouterInterfaceExists("alicloud_router_interface.initiate", &ri),
					testAccCheckRouterInterfaceExists("alicloud_router_interface.opposite", &oppoRI),
					testAccCheckRouterInterfaceConnectionExists("alicloud_router_interface_connection.foo"),
					testAccCheckRouterInterfaceConnectionExists("alicloud_router_interface_connection.bar"),
					resource.TestCheckResourceAttr(
						"alicloud_router_interface.initiate", "instance_charge_type", "PostPaid"),
				),
			},
		},
	})

}

func testAccCheckRouterInterfaceConnectionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No interface ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpcService := VpcService{client}

		response, err := vpcService.DescribeRouterInterface(client.RegionId, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error finding interface %s: %#v", rs.Primary.ID, err)
		}
		if response.Status != string(Active) {
			return fmt.Errorf("Error connection router interface id %s is not Active.", response.RouterInterfaceId)
		}

		return nil
	}
}

func testAccCheckRouterInterfaceConnectionDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_router_interface_connection" {
			continue
		}

		// Try to find the interface
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpcService := VpcService{client}

		ri, err := vpcService.DescribeRouterInterface(client.RegionId, rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if ri.Status == string(Active) {
			return fmt.Errorf("Interface connection %s still exists.", rs.Primary.ID)
		}
	}
	return nil
}

const testAccRouterInterfaceConnectionConfig = `
provider "alicloud" {
  region = "${var.region}"
}

variable "region" {
  default = "cn-hangzhou"
}
variable "name" {
  default = "tf-testAccAlicloudRIConnection_basic"
}
resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vpc" "bar" {
  provider = "alicloud"
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_router_interface" "initiate" {
  opposite_region = "${var.region}"
  router_type = "VRouter"
  router_id = "${alicloud_vpc.foo.router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "${var.name}"
	description = "${var.name}"
	instance_charge_type = "PostPaid"
}

resource "alicloud_router_interface" "opposite" {
  provider = "alicloud"
  opposite_region = "${var.region}"
  router_type = "VRouter"
  router_id = "${alicloud_vpc.bar.router_id}"
  role = "AcceptingSide"
  specification = "Large.1"
  name = "${var.name}-opposite"
  description = "${var.name}-opposite"
}

resource "alicloud_router_interface_connection" "foo" {
  interface_id = "${alicloud_router_interface.initiate.id}"
  opposite_interface_id = "${alicloud_router_interface.opposite.id}"
  depends_on = ["alicloud_router_interface_connection.bar"]
}

resource "alicloud_router_interface_connection" "bar" {
  provider = "alicloud"
  interface_id = "${alicloud_router_interface.opposite.id}"
  opposite_interface_id = "${alicloud_router_interface.initiate.id}"
}`
