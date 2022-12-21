package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudCenChildInstanceRouteEntryToAttachment_basic1977(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_child_instance_route_entry_to_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudCenChildInstanceRouteEntryToAttachmentMap1977)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenChildInstanceRouteEntryToAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.CenChildInstanceRouteEntryToAttachmentSupportRegions)
	name := fmt.Sprintf("tf-testacc%sCenChildInstanceRouteEntryToAttachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenChildInstanceRouteEntryToAttachmentBasicDependence1977)
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
					"transit_router_attachment_id":  "${alicloud_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}",
					"cen_id":                        "${alicloud_cen_instance.default.id}",
					"destination_cidr_block":        "10.0.0.0/24",
					"child_instance_route_table_id": "${alicloud_route_table.foo.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_attachment_id":  CHECKSET,
						"cen_id":                        CHECKSET,
						"destination_cidr_block":        CHECKSET,
						"child_instance_route_table_id": CHECKSET,
						"status":                        "Available",
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudCenChildInstanceRouteEntryToAttachmentMap1977 = map[string]string{}

func AlicloudCenChildInstanceRouteEntryToAttachmentBasicDependence1977(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_cen_transit_router_available_resources" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

data "alicloud_zones" "default" {}

resource "alicloud_vswitch" "default_master" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.1.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_vswitch" "default_slave" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.2.0/24"
  zone_id      = data.alicloud_zones.default.zones.1.id
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level  = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
  support_multicast = true
}

resource "alicloud_cen_transit_router_vpc_attachment" "default" {
  cen_id            = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
  vpc_id            = alicloud_vpc.default.id
  zone_mappings {
    zone_id    = alicloud_vswitch.default_master.zone_id
    vswitch_id = alicloud_vswitch.default_master.id
  }
  zone_mappings {
    zone_id    = alicloud_vswitch.default_slave.zone_id
    vswitch_id = alicloud_vswitch.default_slave.id
  }
  transit_router_attachment_name        = var.name
  transit_router_attachment_description = var.name
}

resource "alicloud_route_table" "foo" {
  vpc_id           = alicloud_vpc.default.id
  route_table_name = var.name
  description      = var.name
}

`, name)
}
