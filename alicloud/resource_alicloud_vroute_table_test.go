package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudRouteTable_basic(t *testing.T) {
	var routeTable vpc.DescribeRouteTablesResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: "alicloud_route_table.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteTableDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRouteTableConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTablesExists("alicloud_route_table.foo", &routeTable),
					resource.TestCheckResourceAttrSet(
						"alicloud_route_table.foo", "vpc_id"),
					resource.TestCheckResourceAttrSet(
						"alicloud_route_table.foo", "route_table_name"),
					resource.TestCheckResourceAttrSet(
						"alicloud_route_table.foo", "description"),
				),
			},
		},
	})
}

func testAccCheckRouteTablesExists(n string, routeTable *vpc.DescribeRouteTablesResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Route Table ID is set")
		}
		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribeRouteTable(rs.Primary.ID)
		if err != nil {
			return err
		}
		if routeTable == nil || len((*routeTable).RouteTables.RouteTable) <= 0 {
			return err
		}
		(*routeTable).RouteTables.RouteTable[0].RouteTableId = instance.RouteTableId
		return nil
	}
}

func testAccCheckRouteTableDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_route_table" {
			continue
		}
		instance, err := client.DescribeRouteTable(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Describe Route Table error %#v", err)
		}
		if instance.RouteTableId != "" {
			return fmt.Errorf("Route Table %s still exist", instance.RouteTableId)
		}
	}
	return nil
}

const testAccRouteTableConfig = `

provider "alicloud" {
  region     = "cn-shanghai"
}

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

resource "alicloud_route_table" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  route_table_name = "test_route_table"
  description = "test_route_table"
}

`
