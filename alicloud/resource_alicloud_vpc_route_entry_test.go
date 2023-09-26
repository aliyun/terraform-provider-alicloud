package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudVpcRouteEntry_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_route_entry.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcrouteentry%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcRouteEntryBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "172.11.1.1/32",
					"nexthop_type":          "Instance",
					"nexthop_id":            "${alicloud_instance.default.id}",
					"name":                  "${var.name}",
					"description":           "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": "172.11.1.1/32",
						"nexthop_type":          "Instance",
						"nexthop_id":            CHECKSET,
						"name":                  CHECKSET,
						"description":           CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_cidrblock": "2.2.2.2/32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_cidrblock": "2.2.2.2/32",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"nexthop_id", "nexthop_type"},
			},
		},
	})
}

var AlicloudVpcRouteEntryMap0 = map[string]string{
	"status":         CHECKSET,
	"route_entry_id": CHECKSET,
}

func AlicloudVpcRouteEntryBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count    = 1
	memory_size       = 2
}

data "alicloud_images" "default" {
	name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
	most_recent = true
	owners      = "system"
}

resource "alicloud_vpc" "default" {
	vpc_name   = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
	vpc_id            = "${alicloud_vpc.default.id}"
	cidr_block        = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	vswitch_name      = "${var.name}"
}

resource "alicloud_security_group" "default" {
	name        = "${var.name}"
	description = "default"
	vpc_id      = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group_rule" "default" {
	type              = "ingress"
	ip_protocol       = "tcp"
	nic_type          = "intranet"
	policy            = "accept"
	port_range        = "22/22"
	priority          = 1
	security_group_id = "${alicloud_security_group.default.id}"
	cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_instance" "default" {
	security_groups            = ["${alicloud_security_group.default.id}"]
	vswitch_id                 = "${alicloud_vswitch.default.id}"
	allocate_public_ip         = true
	instance_charge_type       = "PostPaid"
	instance_type              = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type       = "PayByTraffic"
	internet_max_bandwidth_out = 5
	system_disk_category       = "cloud_efficiency"
	image_id                   = "${data.alicloud_images.default.images.0.id}"
	instance_name              = "${var.name}"
}

`, name)
}

func TestAccAliCloudVpcRouteEntry_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_route_entry.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcRouteEntryMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcrouteentry%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcRouteEntryBasicDependence1)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "172.11.1.1/32",
					"nexthop_type":          "RouterInterface",
					"nexthop_id":            "${alicloud_router_interface.default.id}",
					"route_entry_name":      "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": "172.11.1.1/32",
						"nexthop_type":          "RouterInterface",
						"nexthop_id":            CHECKSET,
						"route_entry_name":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"route_entry_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_entry_name": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_cidrblock": "2.2.2.2/32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_cidrblock": "2.2.2.2/32",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"nexthop_id", "nexthop_type"},
			},
		},
	})
}

var AlicloudVpcRouteEntryMap1 = map[string]string{
	"status":         CHECKSET,
	"route_entry_id": CHECKSET,
}

func AlicloudVpcRouteEntryBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

data "alicloud_regions" "default" {
	current = true
}

resource "alicloud_vpc" "default" {
	vpc_name   = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
	vpc_id            = "${alicloud_vpc.default.id}"
	cidr_block        = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	vswitch_name      = "${var.name}"
}

resource "alicloud_security_group" "default" {
	name        = "${var.name}"
	description = "${var.name}"
	vpc_id      = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group_rule" "default" {
	type              = "ingress"
	ip_protocol       = "tcp"
	nic_type          = "intranet"
	policy            = "accept"
	port_range        = "22/22"
	priority          = 1
	security_group_id = "${alicloud_security_group.default.id}"
	cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_router_interface" "default" {
	opposite_region = "${data.alicloud_regions.default.regions.0.id}"
	router_type     = "VRouter"
	router_id       = "${alicloud_vpc.default.router_id}"
	role            = "InitiatingSide"
	specification   = "Large.2"
	name            = "${var.name}"
	description     = "${var.name}"
}

`, name)
}

func TestAccAliCloudVpcRouteEntry_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_route_entry.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcRouteEntryMap2)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcrouteentry%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcRouteEntryBasicDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "172.11.1.1/32",
					"nexthop_type":          "NatGateway",
					"nexthop_id":            "${alicloud_nat_gateway.default.id}",
					"name":                  "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": "172.11.1.1/32",
						"nexthop_type":          "NatGateway",
						"nexthop_id":            CHECKSET,
						"name":                  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_cidrblock": "2.2.2.2/32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_cidrblock": "2.2.2.2/32",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"nexthop_id", "nexthop_type"},
			},
		},
	})
}

var AlicloudVpcRouteEntryMap2 = map[string]string{
	"status":         CHECKSET,
	"route_entry_id": CHECKSET,
}

func AlicloudVpcRouteEntryBasicDependence2(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_enhanced_nat_available_zones" "default" {}

resource "alicloud_vpc" "default" {
	vpc_name   = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
	vpc_id            = "${alicloud_vpc.default.id}"
	cidr_block        = "10.1.1.0/24"
	availability_zone = "${data.alicloud_enhanced_nat_available_zones.default.zones.0.zone_id}"
	vswitch_name      = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
	vpc_id           = "${alicloud_vpc.default.id}"
	nat_type         = "Enhanced"
	vswitch_id       = alicloud_vswitch.default.id
	nat_gateway_name = "${var.name}"
}

`, name)
}

func TestAccAliCloudVpcRouteEntry_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_route_entry.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcRouteEntryMap3)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcrouteentry%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcRouteEntryBasicDependence3)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "172.11.1.1/32",
					"nexthop_type":          "NetworkInterface",
					"nexthop_id":            "${alicloud_ecs_network_interface.default.id}",
					"name":                  "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": "172.11.1.1/32",
						"nexthop_type":          "NetworkInterface",
						"nexthop_id":            CHECKSET,
						"name":                  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_cidrblock": "2.2.2.2/32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_cidrblock": "2.2.2.2/32",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"nexthop_id", "nexthop_type"},
			},
		},
	})
}

var AlicloudVpcRouteEntryMap3 = map[string]string{
	"status":         CHECKSET,
	"route_entry_id": CHECKSET,
}

func AlicloudVpcRouteEntryBasicDependence3(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
	vpc_name   = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
	vswitch_name      = "${var.name}"
	cidr_block        = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	vpc_id            = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group" "default" {
	name   = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_ecs_network_interface" "default" {
	name            = "${var.name}"
	vswitch_id      = "${alicloud_vswitch.default.id}"
	security_groups = ["${alicloud_security_group.default.id}"]
}

resource "alicloud_route_entry" "defaultMultiple" {
	count                 = 5
	route_table_id        = "${alicloud_vpc.default.route_table_id}"
	destination_cidrblock = "172.16.${count.index}.0/24"
	nexthop_type          = "NetworkInterface"
	nexthop_id            = "${alicloud_ecs_network_interface.default.id}"
	name                  = "${var.name}"
}
`, name)
}

func TestAccAliCloudVpcRouteEntry_basic4(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_route_entry.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcRouteEntryMap4)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcrouteentry%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcRouteEntryBasicDependence4)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "2001:ffff:ffff:ffff::/64",
					"nexthop_type":          "Instance",
					"nexthop_id":            "${alicloud_instance.default.id}",
					"name":                  "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": "2001:ffff:ffff:ffff::/64",
						"nexthop_type":          "Instance",
						"nexthop_id":            CHECKSET,
						"name":                  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"destination_cidrblock": "2.2.2.2/32",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_cidrblock": "2.2.2.2/32",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"nexthop_id", "nexthop_type"},
			},
		},
	})
}

var AlicloudVpcRouteEntryMap4 = map[string]string{
	"status":         CHECKSET,
	"route_entry_id": CHECKSET,
}

func AlicloudVpcRouteEntryBasicDependence4(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
}

data "alicloud_instance_types" "default" {
	availability_zone = data.alicloud_zones.default.zones.0.id
	system_disk_category = "cloud_efficiency"
	cpu_core_count = 4
	minimum_eni_ipv6_address_quantity = 1
}

data "alicloud_images" "default" {
	name_regex  = "^ubuntu_18.*64"
	most_recent = true
	owners      = "system"
}

resource "alicloud_vpc" "default" {
	vpc_name    = "${var.name}"
	cidr_block  = "10.1.0.0/21"
	enable_ipv6 = "true"
}

resource "alicloud_vswitch" "default" {
	vpc_id               = "${alicloud_vpc.default.id}"
	cidr_block           = "10.1.1.0/24"
	availability_zone    = "${data.alicloud_zones.default.zones.0.id}"
	vswitch_name         = "${var.name}"
	ipv6_cidr_block_mask = 64
}

resource "alicloud_security_group" "default" {
	name        = "${var.name}"
	description = "default"
	vpc_id      = "${alicloud_vpc.default.id}"
}

resource "alicloud_instance" "default" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	ipv6_address_count = 1
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	system_disk_category = "cloud_efficiency"
	image_id = "${data.alicloud_images.default.images.0.id}"
	instance_name = var.name
	vswitch_id = "${alicloud_vswitch.default.id}"
	internet_max_bandwidth_out = 10
	security_groups = "${alicloud_security_group.default.*.id}"
}

`, name)
}
