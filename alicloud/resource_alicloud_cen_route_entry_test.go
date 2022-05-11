package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCenRouteEntry_basic(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.CenNoSkipRegions)
	var routeEntry cbn.PublishedRouteEntry
	resourceId := "alicloud_cen_route_entry.default"
	ra := resourceAttrInit(resourceId, cenRouteEntryBasicMap)

	serviceFunc := func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &routeEntry, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf(`"tf-testAcc%sCenRouteEntryConfig-%d"`, defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceCenRouteEntryConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":    "${alicloud_cen_instance.default.id}",
					"route_table_id": "${alicloud_vpc.default.route_table_id}",
					"cidr_block":     "${alicloud_route_entry.default1.destination_cidrblock}",
					"depends_on":     []string{"alicloud_cen_instance_attachment.default"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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

func resourceCenRouteEntryConfigDependence(name string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
	    default = %s
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

	resource "alicloud_cen_instance" "default" {
	    name = "${var.name}"
	}

	resource "alicloud_cen_instance_attachment" "default" {
	    instance_id = "${alicloud_cen_instance.default.id}"
	    child_instance_id = "${alicloud_vpc.default.id}"
	    child_instance_type = "VPC"
	    child_instance_region_id = "%s"
	    depends_on = [
	        "alicloud_vswitch.default"]
	}

	resource "alicloud_route_entry" "default1" {
	    route_table_id = "${alicloud_vpc.default.route_table_id}"
	    destination_cidrblock = "11.0.0.0/16"
	    nexthop_type = "Instance"
	    nexthop_id = "${alicloud_instance.default.id}"
	}
	`, EcsInstanceCommonTestCase, name, defaultRegionToTest)
}

var cenRouteEntryBasicMap = map[string]string{
	"instance_id":    CHECKSET,
	"route_table_id": CHECKSET,
	"cidr_block":     "11.0.0.0/16",
}
