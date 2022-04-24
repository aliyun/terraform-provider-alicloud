package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

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

		response, err := vpcService.DescribeRouterInterfaceConnection(rs.Primary.ID, client.RegionId)
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

		ri, err := vpcService.DescribeRouterInterfaceConnection(rs.Primary.ID, client.RegionId)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		if ri.Status == string(Active) {
			return WrapError(Error("Interface connection %s still exists.", rs.Primary.ID))
		}
	}
	return nil
}

func TestAccAlicloudVPCRouterInterfaceConnectionBasic(t *testing.T) {
	resourceId := "alicloud_router_interface_connection.foo"
	ra := resourceAttrInit(resourceId, testAccRouterInterfaceConnectionCheckMap)
	rand := acctest.RandInt()
	testAccCheck := ra.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouterInterfaceConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouterInterfaceConnectionConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceConnectionExists(resourceId),
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccRouterInterfaceConnectionConfigBasic(rand int) string {
	return fmt.Sprintf(
		`
provider "alicloud" {
  region = "${var.region}"
}
variable "region" {
  default = "cn-hangzhou"
}
variable "name" {
  default = "tf-testAccAlicloudRIConnection_basic%d"
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
}
`, rand)
}

var testAccRouterInterfaceConnectionCheckMap = map[string]string{
	"interface_id":                CHECKSET,
	"opposite_interface_id":       CHECKSET,
	"opposite_router_type":        "VRouter",
	"opposite_router_id":          CHECKSET,
	"opposite_interface_owner_id": CHECKSET,
}
