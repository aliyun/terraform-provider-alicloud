package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudRouteEntry_Basic(t *testing.T) {
	var rt vpc.RouteTable
	var rn vpc.RouteEntry

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_route_entry.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteEntryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRouteEntryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableEntryExists(
						"alicloud_route_entry.foo", &rt, &rn),
					resource.TestCheckResourceAttrSet(
						"alicloud_route_entry.foo", "nexthop_id"),
				),
			},
		},
	})

}

func TestAccAlicloudRouteEntry_RouteInterface(t *testing.T) {
	var rt vpc.RouteTable
	var rn vpc.RouteEntry

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_route_entry.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteEntryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRouteEntryInterfaceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteTableEntryExists(
						"alicloud_route_entry.foo", &rt, &rn),
					resource.TestCheckResourceAttrSet(
						"alicloud_route_entry.foo", "nexthop_id"),
				),
			},
		},
	})

}

func testAccCheckRouteTableExists(rtId string, t *vpc.RouteTable) error {
	client := testAccProvider.Meta().(*AliyunClient)
	//query route table
	rt, terr := client.QueryRouteTableById(rtId)

	if terr != nil {
		return terr
	}

	if rt.RouteTableId != rtId {
		return fmt.Errorf("Route Table not found")
	}

	t = &rt
	return nil
}

func testAccCheckRouteEntryExists(routeTableId, cidrBlock, nextHopType, nextHopId string, e *vpc.RouteEntry) error {
	client := testAccProvider.Meta().(*AliyunClient)
	//query route table entry
	re, rerr := client.QueryRouteEntry(routeTableId, cidrBlock, nextHopType, nextHopId)

	if rerr != nil {
		return rerr
	}

	e = &re
	return nil
}

func testAccCheckRouteTableEntryExists(n string, t *vpc.RouteTable, e *vpc.RouteEntry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Route Entry ID is set")
		}

		parts := strings.Split(rs.Primary.ID, ":")

		//query route table
		err := testAccCheckRouteTableExists(parts[0], t)

		if err != nil {
			return err
		}
		//query route table entry
		err = testAccCheckRouteEntryExists(parts[0], parts[2], parts[3], parts[4], e)
		return err
	}
}

func testAccCheckRouteEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "alicloud_route_entry" || rs.Type != "alicloud_route_entry" {
			continue
		}

		parts := strings.Split(rs.Primary.ID, ":")
		entry, err := client.QueryRouteEntry(parts[0], parts[2], parts[3], parts[4])
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if entry.RouteTableId != "" {
			return fmt.Errorf("Route entry still exist")
		}
	}

	testAccCheckRouterInterfaceDestroy(s)

	return nil
}

const testAccRouteEntryConfig = `
data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}

variable "name" {
	default = "testAccRouteEntryConfig"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_route_entry" "foo" {
	route_table_id = "${alicloud_vpc.foo.route_table_id}"
	destination_cidrblock = "172.11.1.1/32"
	nexthop_type = "Instance"
	nexthop_id = "${alicloud_instance.foo.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_security_group_rule" "ingress" {
	type = "ingress"
	ip_protocol = "tcp"
	nic_type = "intranet"
	policy = "accept"
	port_range = "22/22"
	priority = 1
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]

	vswitch_id = "${alicloud_vswitch.foo.id}"
	allocate_public_ip = true

	# series III
	instance_charge_type = "PostPaid"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5

	system_disk_category = "cloud_efficiency"
	image_id = "${data.alicloud_images.default.images.0.id}"
	instance_name = "${var.name}"
}`

const testAccRouteEntryInterfaceConfig = `
data "alicloud_zones" "default" {
  "available_resource_creation"= "VSwitch"
}
variable "name" {
	default = "testAccRouteEntryInterfaceConfig"
}
resource "alicloud_vpc" "foo" {
  name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
  vpc_id = "${alicloud_vpc.foo.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_route_entry" "foo" {
  route_table_id = "${alicloud_vpc.foo.route_table_id}"
  destination_cidrblock = "172.11.1.1/32"
  nexthop_type = "RouterInterface"
  nexthop_id = "${alicloud_router_interface.interface.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
  name = "${var.name}"
  description = "foo"
  vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_security_group_rule" "ingress" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "accept"
  port_range = "22/22"
  priority = 1
  security_group_id = "${alicloud_security_group.tf_test_foo.id}"
  cidr_ip = "0.0.0.0/0"
}

resource "alicloud_router_interface" "interface" {
  opposite_region = "cn-beijing"
  router_type = "VRouter"
  router_id = "${alicloud_vpc.foo.router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "${var.name}"
  description = "test1"
}`
