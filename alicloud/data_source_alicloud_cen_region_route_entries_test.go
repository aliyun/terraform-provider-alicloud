package alicloud

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCenRegionRouteEntriesDataSource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCenRegionRouteEntriesDataSourceBasic(EcsInstanceCommonTestCase, defaultRegionToTest),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_region_route_entries.entry"),
					resource.TestCheckResourceAttr("data.alicloud_cen_region_route_entries.entry", "entries.0.cidr_block", "11.0.0.0/16"),
					resource.TestCheckResourceAttr("data.alicloud_cen_region_route_entries.entry", "entries.0.type", "CBN"),
					resource.TestCheckResourceAttr("data.alicloud_cen_region_route_entries.entry", "entries.0.next_hop_type", "VPC"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_region_route_entries.entry", "entries.0.next_hop_id"),
					resource.TestCheckResourceAttr("data.alicloud_cen_region_route_entries.entry", "entries.0.next_hop_region_id", defaultRegionToTest),
					resource.TestCheckResourceAttr("data.alicloud_cen_region_route_entries.entry", "entries.#", "2"),
				),
			},
		},
	})
}

func testAccCheckCenRegionRouteEntriesDataSourceBasic(common, region string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccRegionRouteEntriesDataSourceBasic"
	}
	
	resource "alicloud_instance" "default" {
		vswitch_id = "${alicloud_vswitch.default.id}"
		image_id = "${data.alicloud_images.default.images.0.id}"
	
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		system_disk_category = "cloud_efficiency"
	
		internet_charge_type = "PayByTraffic"
		internet_max_bandwidth_out = 5
		security_groups = ["${alicloud_security_group.default.id}"]
		instance_name = "${var.name}-region-route-entry"
	}
	
	resource "alicloud_cen_instance" "cen" {
		name = "${var.name}-cen"
		description = "terraform01"
	}
	
	resource "alicloud_cen_instance_attachment" "attach" {
	    instance_id = "${alicloud_cen_instance.cen.id}"
	    child_instance_id = "${alicloud_vpc.default.id}"
	    child_instance_region_id = "%s"
	}
	
	resource "alicloud_route_entry" "route" {
	    route_table_id = "${alicloud_vpc.default.route_table_id}"
	    destination_cidrblock = "11.0.0.0/16"
	    nexthop_type = "Instance"
	    nexthop_id = "${alicloud_instance.default.id}"
	}
	
	resource "alicloud_cen_route_entry" "foo" {
	    instance_id = "${alicloud_cen_instance.cen.id}"
	    route_table_id = "${alicloud_vpc.default.route_table_id}"
	    cidr_block = "${alicloud_route_entry.route.destination_cidrblock}"
	    depends_on = [
			"alicloud_cen_instance_attachment.attach"]
	}
	
	data "alicloud_cen_region_route_entries" "entry" {
		instance_id = "${alicloud_cen_route_entry.foo.instance_id}"
		region_id = "%s"
	}
	`, common, region, region)
}
