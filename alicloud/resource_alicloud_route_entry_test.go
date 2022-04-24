package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudVPCRouteEntryInstance(t *testing.T) {
	var v *vpc.RouteEntry
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, testAccRouteEntryCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteEntryConfig_instance(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nexthop_type": "Instance",
						"name":         fmt.Sprintf("tf-testAccRouteEntryConfigName%d", rand),
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudVPCRouteEntryInterface(t *testing.T) {
	var v *vpc.RouteEntry
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, testAccRouteEntryCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteEntryConfig_interface(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nexthop_type": "RouterInterface",
						"name":         fmt.Sprintf("tf-testAccRouteEntryInterfaceConfig%d", rand),
					}),
				),
			},
		},
	})
}

func TestAccAlicloudVPCRouteEntryNatGateway(t *testing.T) {
	var v *vpc.RouteEntry
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, testAccRouteEntryCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteEntryConfig_natGateway(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nexthop_type": "NatGateway",
						"name":         fmt.Sprintf("tf-testAccRouteEntryNatGatewayConfig%d", rand),
					}),
				),
			},
		},
	})
}

func TestAccAlicloudVPCRouteEntryMulti(t *testing.T) {
	var v *vpc.RouteEntry
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "alicloud_route_entry.default.4"
	ra := resourceAttrInit(resourceId, testAccRouteEntryCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRouteEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteEntryConfigMulti(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nexthop_type":          "NetworkInterface",
						"destination_cidrblock": "172.16.4.0/24",
						"name":                  fmt.Sprintf("tf-testAccRouteEntryConcurrence%d", rand),
					}),
				),
			},
		},
	})
}

func testAccCheckRouteEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "alicloud_route_entry" || rs.Type != "alicloud_route_entry" {
			continue
		}
		entry, err := vpcService.DescribeRouteEntry(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		if entry.RouteTableId != "" {
			return WrapError(Error("Route entry still exist"))
		}
	}

	testAccCheckRouterInterfaceDestroy(s)

	return nil
}

func testAccRouteEntryConfig_instance(rand int) string {
	return fmt.Sprintf(
		`
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}
data "alicloud_images" "default" {
	name_regex = "^ubuntu"
	most_recent = true
	owners = "system"
}

variable "name" {
	default = "tf-testAccRouteEntryConfigName%d"
}
resource "alicloud_vpc" "default" {
	vpc_name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	vswitch_name = "${var.name}"
}

resource "alicloud_security_group" "default" {
	name = "${var.name}"
	description = "default"
	vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group_rule" "default" {
	type = "ingress"
	ip_protocol = "tcp"
	nic_type = "intranet"
	policy = "accept"
	port_range = "22/22"
	priority = 1
	security_group_id = "${alicloud_security_group.default.id}"
	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_instance" "default" {
	security_groups = ["${alicloud_security_group.default.id}"]

	vswitch_id = "${alicloud_vswitch.default.id}"
	allocate_public_ip = true

	instance_charge_type = "PostPaid"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5

	system_disk_category = "cloud_efficiency"
	image_id = "${data.alicloud_images.default.images.0.id}"
	instance_name = "${var.name}"
}

resource "alicloud_route_entry" "default" {
	route_table_id = "${alicloud_vpc.default.route_table_id}"
	destination_cidrblock = "172.11.1.1/32"
	nexthop_type = "Instance"
	nexthop_id = "${alicloud_instance.default.id}"
	name = "${var.name}"
}
`, rand)
}

func testAccRouteEntryConfig_interface(rand int) string {
	return fmt.Sprintf(
		`
data "alicloud_zones" "default" {
  available_resource_creation= "VSwitch"
}
variable "name" {
	default = "tf-testAccRouteEntryInterfaceConfig%d"
}
resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_security_group" "default" {
  name = "${var.name}"
  description = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group_rule" "default" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "accept"
  port_range = "22/22"
  priority = 1
  security_group_id = "${alicloud_security_group.default.id}"
  cidr_ip = "0.0.0.0/0"
}

data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_router_interface" "default" {
  opposite_region = "${data.alicloud_regions.default.regions.0.id}"
  router_type = "VRouter"
  router_id = "${alicloud_vpc.default.router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "${var.name}"
  description = "${var.name}"
}

resource "alicloud_route_entry" "default" {
  route_table_id = "${alicloud_vpc.default.route_table_id}"
  destination_cidrblock = "172.11.1.1/32"
  nexthop_type = "RouterInterface"
  nexthop_id = "${alicloud_router_interface.default.id}"
  name = "${var.name}"
}
`, rand)
}

func testAccRouteEntryConfig_natGateway(rand int) string {
	return fmt.Sprintf(
		`
data "alicloud_zones" "default" {
  available_resource_creation= "VSwitch"
}
variable "name" {
	default = "tf-testAccRouteEntryNatGatewayConfig%d"
}
resource "alicloud_vpc" "default" {
  vpc_name = "${var.name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	nat_type = "Enhanced"
	vswitch_id = alicloud_vswitch.default.id
	nat_gateway_name = "${var.name}"
}

resource "alicloud_route_entry" "default" {
  route_table_id = "${alicloud_vpc.default.route_table_id}"
  destination_cidrblock = "172.11.1.1/32"
  nexthop_type = "NatGateway"
  nexthop_id = "${alicloud_nat_gateway.default.id}"
  name = "${var.name}"
}`, rand)
}

func testAccRouteEntryConfigMulti(rand int) string {
	return fmt.Sprintf(`
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

variable "name" {
	default = "tf-testAccRouteEntryConcurrence%d"
}
resource "alicloud_vpc" "default" {
	vpc_name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
    vswitch_name = "${var.name}"
    cidr_block = "10.1.1.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_network_interface" "default" {
    name = "${var.name}"
    vswitch_id = "${alicloud_vswitch.default.id}"
    security_groups = [ "${alicloud_security_group.default.id}" ]
}

resource "alicloud_route_entry" "default" {
	count = 5
	route_table_id = "${alicloud_vpc.default.route_table_id}"
	destination_cidrblock = "172.16.${count.index}.0/24"
	nexthop_type = "NetworkInterface"
	nexthop_id = "${alicloud_network_interface.default.id}"
	name = "${var.name}"
}
`, rand)
}

var testAccRouteEntryCheckMap = map[string]string{
	"route_table_id":        CHECKSET,
	"nexthop_id":            CHECKSET,
	"destination_cidrblock": "172.11.1.1/32",
}
