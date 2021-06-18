package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"testing"

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
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":            "${alicloud_cen_instance.default.id}",
					"transit_router_id": "${alicloud_cen_transit_router.default.id}",
					"vpc_id":            "${alicloud_vpc.default.id}",
					"zone_mappings": []map[string]interface{}{
						{
							"vswitch_id": "${alicloud_vswitch.default_master.id}",
							"zone_id":    "cn-hangzhou-h",
						},
						{
							"vswitch_id": "${alicloud_vswitch.default_slave.id}",
							"zone_id":    "cn-hangzhou-i",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":                     CHECKSET,
						"transit_router_id":          CHECKSET,
						"vpc_id":                     CHECKSET,
						"zone_mappings.1.vswitch_id": CHECKSET,
						"zone_mappings.1.zone_id":    "cn-hangzhou-h",
						"zone_mappings.2.vswitch_id": CHECKSET,
						"zone_mappings.2.zone_id":    "cn-hangzhou-i",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_create_vpc_route", "cen_id", "charge_type", "dry_run", "route_table_association_enabled", "route_table_propagation_enabled", "transit_router_id"},
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
					"transit_router_attachment_name": "name1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_attachment_name": "name1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_attachment_description": "descp",
					"transit_router_attachment_name":        "name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_attachment_description": "descp",
						"transit_router_attachment_name":        "name",
					}),
				),
			},
		},
	})
}

var AlicloudCenTransitRouterVpcAttachmentMap = map[string]string{
	"auto_create_vpc_route":                 NOSET,
	"cen_id":                                CHECKSET,
	"charge_type":                           NOSET,
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
	"zone_mappings.#":                       "1",
}

func AlicloudCenTransitRouterVpcAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {	
	default = "%s"
}

resource "alicloud_vpc" "default" {
  vpc_name = "tf-test"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default_master" {
  vswitch_name = "tf-test"
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "192.168.1.0/24"
  zone_id = "cn-hangzhou-h"
}

resource "alicloud_vswitch" "default_slave" {
  vswitch_name = "tf-test"
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "192.168.2.0/24"
  zone_id = "cn-hangzhou-i"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = "${var.name}"
  protection_level = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
cen_id= "${alicloud_cen_instance.default.id}"
region_id = "cn-hangzhou"
}
`, name)
}
