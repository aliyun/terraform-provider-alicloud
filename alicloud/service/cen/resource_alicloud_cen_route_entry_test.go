package cen_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cen"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCenRouteEntry_basic(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.CenSupportRegions)
	resourceId := "alicloud_cen_route_entry.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCenRouteEntryMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cen.NewCenService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCenRouteEntry")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sCenRouteEntry-%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCenRouteEntryBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":    "${alicloud_cen_instance_attachment.default.instance_id}",
					"route_table_id": "${alicloud_route_entry.default.route_table_id}",
					"cidr_block":     "${alicloud_route_entry.default.destination_cidrblock}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    acctest.CHECKSET,
						"route_table_id": acctest.CHECKSET,
						"cidr_block":     acctest.CHECKSET,
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

var AliCloudCenRouteEntryMap0 = map[string]string{}

func AliCloudCenRouteEntryBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_regions" "default" {
 		current = true
	}

	data "alicloud_zones" "default" {
 		available_disk_category     = "cloud_efficiency"
 		available_resource_creation = "VSwitch"
	}

	data "alicloud_images" "default" {
 		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
 		most_recent = true
 		owners      = "system"
	}

	data "alicloud_instance_types" "default" {
 		availability_zone = data.alicloud_zones.default.zones.0.id
 		image_id          = data.alicloud_images.default.images.0.id
	}

	resource "alicloud_vpc" "default" {
 		vpc_name   = var.name
 		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
 		vswitch_name = var.name
 		vpc_id       = alicloud_vpc.default.id
 		cidr_block   = "192.168.192.0/24"
 		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
 		name   = var.name
 		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_instance" "default" {
 		image_id                   = data.alicloud_images.default.images.0.id
 		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
 		security_groups            = alicloud_security_group.default.*.id
 		internet_charge_type       = "PayByTraffic"
 		internet_max_bandwidth_out = "10"
 		availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
 		instance_charge_type       = "PostPaid"
 		system_disk_category       = "cloud_efficiency"
 		vswitch_id                 = alicloud_vswitch.default.id
 		instance_name              = var.name
	}

	resource "alicloud_cen_instance" "default" {
 		cen_instance_name = var.name
 		description       = var.name
	}

	resource "alicloud_cen_instance_attachment" "default" {
 		instance_id              = alicloud_cen_instance.default.id
 		child_instance_id        = alicloud_vswitch.default.vpc_id
 		child_instance_type      = "VPC"
 		child_instance_region_id = data.alicloud_regions.default.regions.0.id
	}

	resource "alicloud_route_entry" "default" {
 		route_table_id        = alicloud_vpc.default.route_table_id
 		destination_cidrblock = "11.0.0.0/16"
 		nexthop_type          = "Instance"
 		nexthop_id            = alicloud_instance.default.id
	}
`, name)
}
