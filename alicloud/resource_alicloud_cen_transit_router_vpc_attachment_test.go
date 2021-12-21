package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCenTransitRouterVpcAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_vpc_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterVpcAttachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterVpcAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenTransitRouterVpcAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenTransitRouterVpcAttachmentBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.CenTransitRouterVpcAttachmentSupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":                                "${alicloud_cen_instance.default.id}",
					"transit_router_id":                     "${alicloud_cen_transit_router.default.transit_router_id}",
					"transit_router_attachment_name":        name,
					"transit_router_attachment_description": "tf-test",
					"vpc_id":                                "${data.alicloud_vpcs.default.ids.0}",
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${local.vswitch_id_master}",
							"zone_id":    "${alicloud_vswitch.default_master.zone_id}",
						},
						{
							"vswitch_id": "${local.vswitch_id}",
							"zone_id":    "${alicloud_vswitch.default_slave.zone_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":                                CHECKSET,
						"transit_router_id":                     CHECKSET,
						"transit_router_attachment_name":        name,
						"transit_router_attachment_description": "tf-test",
						"vpc_id":                                CHECKSET,
						"zone_mappings.0.vswitch_id":            CHECKSET,
						"zone_mappings.0.zone_id":               CHECKSET,
						"zone_mappings.1.vswitch_id":            CHECKSET,
						"zone_mappings.1.zone_id":               CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_create_vpc_route", "dry_run", "route_table_association_enabled", "route_table_propagation_enabled", "transit_router_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_type": "VPC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_type": "VPC",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_attachment_description": "desp1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_attachment_description": "desp1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_attachment_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_attachment_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_attachment_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_attachment_name": name,
					}),
				),
			},
		},
	})
}

var AlicloudCenTransitRouterVpcAttachmentMap = map[string]string{
	"auto_create_vpc_route":                 NOSET,
	"cen_id":                                CHECKSET,
	"dry_run":                               NOSET,
	"resource_type":                         "VPC",
	"route_table_association_enabled":       NOSET,
	"route_table_propagation_enabled":       NOSET,
	"status":                                CHECKSET,
	"transit_router_attachment_description": CHECKSET,
	"transit_router_attachment_name":        CHECKSET,
	"transit_router_id":                     CHECKSET,
	"vpc_id":                                CHECKSET,
	"vpc_owner_id":                          CHECKSET,
	"zone_mappings.#":                       "2",
}

func AlicloudCenTransitRouterVpcAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {	
	default = "%s"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_zones.default.zones.0.id
}

data "alicloud_vswitches" "default_master" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_zones.default.zones.1.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}

resource "alicloud_vswitch" "vswitch_master" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.1.id
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  vswitch_id_master = length(data.alicloud_vswitches.default_master.ids) > 0 ? data.alicloud_vswitches.default_master.ids[0] : concat(alicloud_vswitch.vswitch_master.*.id, [""])[0]
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
cen_id= alicloud_cen_instance.default.id
}
`, name)
}
