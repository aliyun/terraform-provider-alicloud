package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCenPublishedRouteEntriesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCenPublishedRouteEntriesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_route_entries.entry"),
					resource.TestCheckResourceAttr("data.alicloud_cen_route_entries.entry", "entries.0.cidr_block", "11.0.0.0/16"),
					resource.TestCheckResourceAttr("data.alicloud_cen_route_entries.entry", "entries.0.next_hop_type", "Instance"),
					resource.TestCheckResourceAttr("data.alicloud_cen_route_entries.entry", "entries.0.route_type", "Custom"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_route_entries.entry", "entries.0.route_table_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_route_entries.entry", "entries.0.next_hop_id"),
					resource.TestCheckResourceAttr("data.alicloud_cen_route_entries.entry", "entries.0.publish_status", "Published"),
					resource.TestCheckResourceAttr("data.alicloud_cen_route_entries.entry", "entries.0.operational_mode", "true"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_route_entries.entry", "entries.0.conflicts"),
				),
			},
		},
	})
}

func TestAccAlicloudCenPublishedRouteEntriesDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCenPublishedRouteEntriesDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_route_entries.entry"),
					resource.TestCheckResourceAttr("data.alicloud_cen_route_entries.entry", "entries.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_route_entries.entry", "entries.0.cidr_block"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_route_entries.entry", "entries.0.next_hop_type"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_route_entries.entry", "entries.0.route_type"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_route_entries.entry", "entries.0.route_table_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_route_entries.entry", "entries.0.next_hop_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_route_entries.entry", "entries.0.publish_status"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_route_entries.entry", "entries.0.operational_mode"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_route_entries.entry", "entries.0.conflicts"),
				),
			},
		},
	})
}

const testAccCheckCenPublishedRouteEntriesDataSourceConfig = `
provider "alicloud" {
    alias = "bj"
    region = "cn-beijing"
}

variable "name" {
	default = "tf-testAccCenRoutes"
}

data "alicloud_zones" "default" {
    provider = "alicloud.bj"
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}

data "alicloud_instance_types" "default" {
    provider = "alicloud.bj"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

data "alicloud_images" "default" {
    provider = "alicloud.bj"
    name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}

resource "alicloud_vpc" "vpc" {
    provider = "alicloud.bj"
  	name = "${var.name}"
  	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
    provider = "alicloud.bj"
 	vpc_id = "${alicloud_vpc.vpc.id}"
 	cidr_block = "172.16.0.0/21"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 	name = "${var.name}"
}

resource "alicloud_security_group" "default" {
    provider = "alicloud.bj"
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_instance" "default" {
    provider = "alicloud.bj"
	vswitch_id = "${alicloud_vswitch.default.id}"
	image_id = "${data.alicloud_images.default.images.0.id}"

	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	system_disk_category = "cloud_efficiency"

	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5
	security_groups = ["${alicloud_security_group.default.id}"]
	instance_name = "${var.name}"
}

resource "alicloud_cen_instance" "cen" {
	name = "${var.name}"
	description = "terraform01"
}

resource "alicloud_cen_instance_attachment" "attach" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc.id}"
    child_instance_region_id = "cn-beijing"
}

resource "alicloud_route_entry" "route" {
    provider = "alicloud.bj"
    route_table_id = "${alicloud_vpc.vpc.route_table_id}"
    destination_cidrblock = "11.0.0.0/16"
    nexthop_type = "Instance"
    nexthop_id = "${alicloud_instance.default.id}"
}

resource "alicloud_cen_route_entry" "foo" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    route_table_id = "${alicloud_vpc.vpc.route_table_id}"
    cidr_block = "${alicloud_route_entry.route.destination_cidrblock}"
    depends_on = ["alicloud_cen_instance_attachment.attach"]
}

data "alicloud_cen_route_entries" "entry" {
	instance_id = "${alicloud_cen_route_entry.foo.instance_id}"
	route_table_id = "${alicloud_cen_route_entry.foo.route_table_id}"
 	cidr_block = "11.0.0.0/16"
}
`

const testAccCheckCenPublishedRouteEntriesDataSourceEmpty = `
provider "alicloud" {
    alias = "bj"
    region = "cn-beijing"
}

variable "name" {
	default = "tf-testAccCenRoutesEmpty"
}

data "alicloud_zones" "default" {
    provider = "alicloud.bj"
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}

data "alicloud_instance_types" "default" {
    provider = "alicloud.bj"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

data "alicloud_images" "default" {
    provider = "alicloud.bj"
    name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}

resource "alicloud_vpc" "vpc" {
    provider = "alicloud.bj"
  	name = "${var.name}"
  	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
    provider = "alicloud.bj"
 	vpc_id = "${alicloud_vpc.vpc.id}"
 	cidr_block = "172.16.0.0/21"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 	name = "${var.name}"
}

resource "alicloud_security_group" "default" {
    provider = "alicloud.bj"
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_instance" "default" {
    provider = "alicloud.bj"
	vswitch_id = "${alicloud_vswitch.default.id}"
	image_id = "${data.alicloud_images.default.images.0.id}"

	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	system_disk_category = "cloud_efficiency"

	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5
	security_groups = ["${alicloud_security_group.default.id}"]
	instance_name = "${var.name}"
}

resource "alicloud_cen_instance" "cen" {
	name = "${var.name}"
	description = "terraform01"
}

resource "alicloud_cen_instance_attachment" "attach" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc.id}"
    child_instance_region_id = "cn-beijing"
}

resource "alicloud_route_entry" "route" {
    provider = "alicloud.bj"
    route_table_id = "${alicloud_vpc.vpc.route_table_id}"
    destination_cidrblock = "11.0.0.0/16"
    nexthop_type = "Instance"
    nexthop_id = "${alicloud_instance.default.id}"
}

resource "alicloud_cen_route_entry" "foo" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    route_table_id = "${alicloud_vpc.vpc.route_table_id}"
    cidr_block = "${alicloud_route_entry.route.destination_cidrblock}"
    depends_on = ["alicloud_cen_instance_attachment.attach"]
}

data "alicloud_cen_route_entries" "entry" {
	instance_id = "${alicloud_cen_route_entry.foo.instance_id}"
	route_table_id = "${alicloud_cen_route_entry.foo.route_table_id}"
 	cidr_block = "11.1.0.0/16"
}
`
