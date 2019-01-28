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
				Config: testAccCheckCenRegionRouteEntriesDataSourceBasic(defaultRegionToTest),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_region_route_entries.entry"),
					resource.TestCheckResourceAttr("data.alicloud_cen_region_route_entries.entry", "entries.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_cen_region_route_entries.entry", "entries.0.cidr_block", "100.64.0.0/10"),
					resource.TestCheckResourceAttr("data.alicloud_cen_region_route_entries.entry", "entries.0.type", "System"),
					resource.TestCheckResourceAttr("data.alicloud_cen_region_route_entries.entry", "entries.0.next_hop_type", "local_service"),
					resource.TestCheckResourceAttr("data.alicloud_cen_region_route_entries.entry", "entries.0.next_hop_id", ""),
					resource.TestCheckResourceAttr("data.alicloud_cen_region_route_entries.entry", "entries.0.next_hop_region_id", ""),
				),
			},
		},
	})
}

func TestAccAlicloudCenRegionRouteEntriesDataSource_empty(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCenRegionRouteEntriesDataSourceEmpty(defaultRegionToTest),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_region_route_entries.entry"),
					resource.TestCheckResourceAttr("data.alicloud_cen_region_route_entries.entry", "entries.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_region_route_entries.entry", "entries.0.cidr_block"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_region_route_entries.entry", "entries.0.type"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_region_route_entries.entry", "entries.0.next_hop_type"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_region_route_entries.entry", "entries.0.next_hop_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_region_route_entries.entry", "entries.0.next_hop_region_id"),
				),
			},
		},
	})
}

func testAccCheckCenRegionRouteEntriesDataSourceBasic(region string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccRegionRouteEntriesDataSourceBasic"
	}

	resource "alicloud_vpc" "default" {
  		name = "${var.name}-vpc"
  		cidr_block = "172.16.0.0/12"
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
	
	data "alicloud_cen_region_route_entries" "entry" {
		instance_id = "${alicloud_cen_instance_attachment.attach.instance_id}"
		region_id = "%s"
    	output_file = "data_resource_tmp.txt"
	}
	`, region, region)
}

func testAccCheckCenRegionRouteEntriesDataSourceEmpty(region string) string {
	return fmt.Sprintf(`
	data "alicloud_cen_region_route_entries" "entry" {
		instance_id = "cen-cidwnnenc"
		region_id = "%s"
	}
	`, region)
}
