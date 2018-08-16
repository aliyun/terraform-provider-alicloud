package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudRouterInterfacesDataSource_oneRouter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRouterInterfacesDataSourceOneRouterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_router_interfaces.router_interfaces"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.status", "Idle"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.name", "testAccCheckAlicloudRouterInterfacesDataSourceOneRouterConfig"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.description", "testAccCheckAlicloudRouterInterfacesDataSourceOneRouterConfig_descr"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.role", "InitiatingSide"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.specification", "Large.2"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.router_id"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.router_type", "VRouter"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.vpc_id"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.access_point_id", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.opposite_region_id"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.opposite_interface_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.opposite_router_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.opposite_router_type", "VRouter"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.opposite_interface_owner_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.health_check_source_ip", ""),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.router_interfaces", "interfaces.0.health_check_target_ip", ""),
				),
			},
		},
	})
}

func TestAccAlicloudRouterInterfacesDataSource_twoRouters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRouterInterfacesDataSourceTwoRoutersConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_router_interfaces.foo_router_interfaces"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.0.status", "Active"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.0.name", "testAccCheckAlicloudRouterInterfacesDataSourceTwoRoutersConfig_initiating"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.0.description", "testAccCheckAlicloudRouterInterfacesDataSourceTwoRoutersConfig_initiating_descr"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.0.role", "InitiatingSide"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.0.specification", "Large.2"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.0.router_id"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.0.router_type", "VRouter"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.0.vpc_id"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.0.access_point_id", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.0.opposite_region_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.0.opposite_interface_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.0.opposite_router_id"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.0.opposite_router_type", "VRouter"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.0.opposite_interface_owner_id"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.0.health_check_source_ip", ""),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.foo_router_interfaces", "interfaces.0.health_check_target_ip", ""),

					testAccCheckAlicloudDataSourceID("data.alicloud_router_interfaces.bar_router_interfaces"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.0.status", "Active"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.0.name", "testAccCheckAlicloudRouterInterfacesDataSourceTwoRoutersConfig_accepting"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.0.description", "testAccCheckAlicloudRouterInterfacesDataSourceTwoRoutersConfig_accepting_descr"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.0.role", "AcceptingSide"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.0.specification", "Negative"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.0.router_id"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.0.router_type", "VRouter"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.0.vpc_id"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.0.access_point_id", ""),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.0.opposite_region_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.0.opposite_interface_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.0.opposite_router_id"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.0.opposite_router_type", "VRouter"),
					resource.TestCheckResourceAttrSet("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.0.opposite_interface_owner_id"),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.0.health_check_source_ip", ""),
					resource.TestCheckResourceAttr("data.alicloud_router_interfaces.bar_router_interfaces", "interfaces.0.health_check_target_ip", ""),
				),
			},
		},
	})
}

const testAccCheckAlicloudRouterInterfacesDataSourceOneRouterConfig = `
variable "name" {
	default = "testAccCheckAlicloudRouterInterfacesDataSourceOneRouterConfig"
}

resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

data "alicloud_regions" "current_regions" {
  current = true
}

resource "alicloud_router_interface" "interface" {
  opposite_region = "${data.alicloud_regions.current_regions.regions.0.id}"
  router_type = "VRouter"
  router_id = "${alicloud_vpc.foo.router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "${var.name}"
  description = "${var.name}_descr"
}

data "alicloud_router_interfaces" "router_interfaces" {
  router_id = "${alicloud_vpc.foo.router_id}"
  specification = "${alicloud_router_interface.interface.specification}"
}
`

const testAccCheckAlicloudRouterInterfacesDataSourceTwoRoutersConfig = `
variable "name" {
	default = "testAccCheckAlicloudRouterInterfacesDataSourceTwoRoutersConfig"
}

resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vpc" "bar" {
  name = "${var.name}"
  cidr_block = "192.168.0.0/16"
}

data "alicloud_regions" "current_regions" {
  current = true
}

resource "alicloud_router_interface" "initiating" {
  opposite_region = "${data.alicloud_regions.current_regions.regions.0.id}"
  router_type = "VRouter"
  router_id = "${alicloud_vpc.foo.router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "${var.name}_initiating"
  description = "${var.name}_initiating_descr"
}

resource "alicloud_router_interface" "accepting" {
  opposite_region = "${data.alicloud_regions.current_regions.regions.0.id}"
  router_type = "VRouter"
  router_id = "${alicloud_vpc.bar.router_id}"
  role = "AcceptingSide"
  specification = "Large.2"
  name = "${var.name}_accepting"
  description = "${var.name}_accepting_descr"
}

resource "alicloud_router_interface_connection" "bar" {
  interface_id = "${alicloud_router_interface.accepting.id}"
  opposite_interface_id = "${alicloud_router_interface.initiating.id}"
}

resource "alicloud_router_interface_connection" "foo" {
  interface_id = "${alicloud_router_interface.initiating.id}"
  opposite_interface_id = "${alicloud_router_interface.accepting.id}"
  depends_on = [
    "alicloud_router_interface_connection.bar"
  ]
}

data "alicloud_router_interfaces" "foo_router_interfaces" {
  opposite_interface_id = "${alicloud_router_interface_connection.foo.opposite_interface_id}"
  specification = "${alicloud_router_interface.initiating.specification}"
}

data "alicloud_router_interfaces" "bar_router_interfaces" {
  opposite_interface_id = "${alicloud_router_interface_connection.foo.interface_id}"
}
`
