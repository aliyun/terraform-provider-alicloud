package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudRouterInterface_basic(t *testing.T) {
	var vpc vpc.DescribeVpcAttributeResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_router_interface.interface",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouterInterfaceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRouterInterfaceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists(
						"alicloud_vpc.foo", &vpc),
					testAccCheckRouterInterfaceExists(
						"alicloud_router_interface.interface"),
				),
			},
		},
	})

}

func testAccCheckRouterInterfaceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No interface ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)

		response, err := client.DescribeRouterInterface(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error finding interface %s: %#v", rs.Primary.ID, err)
		}
		if response.RouterInterfaceId != rs.Primary.ID {
			return fmt.Errorf("Error finding interface %s", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckRouterInterfaceDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_router_interface" {
			continue
		}

		// Try to find the interface
		client := testAccProvider.Meta().(*AliyunClient)

		ri, err := client.DescribeRouterInterface(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if ri.RouterInterfaceId == rs.Primary.ID {
			return fmt.Errorf("Interface %s still exists.", rs.Primary.ID)
		}
	}
	return nil
}

const testAccRouterInterfaceConfig = `
resource "alicloud_vpc" "foo" {
  name = "tf_test_foo12345"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_router_interface" "interface" {
  opposite_region = "cn-beijing"
  router_type = "VRouter"
  router_id = "${alicloud_vpc.foo.router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "test1"
  description = "test1"
}`
