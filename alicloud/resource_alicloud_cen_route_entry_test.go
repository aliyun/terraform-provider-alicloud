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
			{
				Config: testAccCenRouteEntryConfig(EcsInstanceCommonTestCase, defaultRegionToTest),
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

func testAccCenRouteEntryConfig(common, region string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
	    default = "tf-testAccCenRouteEntryConfig"
	}

	resource "alicloud_instance" "default" {
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
	    child_instance_id = "${alicloud_vpc.default.id}"
	    child_instance_region_id = "%s"
	    depends_on = [
	        "alicloud_vswitch.default"]
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
	`, common, region)
}
