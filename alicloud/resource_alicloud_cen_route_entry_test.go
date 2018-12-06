package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCenRouteEntry_basic(t *testing.T) {
	var routeEntry cbn.PublishedRouteEntry

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cen_route_entry.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCenRouteEntryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenRouteEntryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenRouteEntryExists("alicloud_cen_route_entry.foo", &routeEntry),
					resource.TestCheckResourceAttr("alicloud_cen_route_entry.foo", "cidr_block", "11.0.0.0/16"),
					resource.TestCheckResourceAttrSet("alicloud_cen_route_entry.foo", "instance_id"),
					resource.TestCheckResourceAttrSet("alicloud_cen_route_entry.foo", "route_table_id"),
				),
			},
		},
	})
}

func testAccCheckCenRouteEntryExists(n string, routeEntry *cbn.PublishedRouteEntry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cen Route Entry Publishment ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		cenService := CenService{client}

		routeEntryItem, err := cenService.DescribePublishedRouteEntriesById(rs.Primary.ID)
		if err != nil {
			return err
		}

		if routeEntryItem.PublishStatus != string(PUBLISHED) {
			return fmt.Errorf("CEN route entry %s status error", rs.Primary.ID)
		}

		*routeEntry = routeEntryItem
		return nil
	}
}

func testAccCheckCenRouteEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	cenService := CenService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cen_route_entry" {
			continue
		}

		routeEntryItem, err := cenService.DescribePublishedRouteEntriesById(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if routeEntryItem.PublishStatus == string(NOPUBLISHED) {
			continue
		} else {
			return fmt.Errorf("CEN route entry %s status error", rs.Primary.ID)
		}
	}

	return nil
}

const testAccCenRouteEntryConfig = `
provider "alicloud" {
    alias = "bj"
    region = "cn-beijing"
}

variable "name" {
    default = "tf-testAccCenRouteEntryConfig"
}

data "alicloud_zones" "default" {
    provider = "alicloud.bj"
    available_disk_category = "cloud_efficiency"
    available_resource_creation = "VSwitch"
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
    depends_on = [
        "alicloud_cen_instance_attachment.attach"]
}
`
