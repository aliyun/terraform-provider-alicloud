package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccAlicloudCenTransitRouterVbrAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_vbr_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterVbrAttachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterVbrAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenTransitRouterVbrAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenTransitRouterVbrAttachmentBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VbrSupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":                                "${alicloud_cen_instance.default.id}",
					"transit_router_id":                     "${alicloud_cen_transit_router.default.transit_router_id}",
					"transit_router_attachment_name":        "test",
					"transit_router_attachment_description": "test",
					"vbr_id":                                "${alicloud_express_connect_virtual_border_router.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":                                CHECKSET,
						"transit_router_id":                     CHECKSET,
						"transit_router_attachment_name":        "test",
						"transit_router_attachment_description": "test",
						"vbr_id":                                CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run", "resource_type", "route_table_association_enabled", "route_table_propagation_enabled"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_publish_route_enabled": `false`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_publish_route_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_type": "VBR",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_type": "VBR",
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
					"auto_publish_route_enabled":            `true`,
					"transit_router_attachment_description": "desp",
					"transit_router_attachment_name":        name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_publish_route_enabled":            "true",
						"transit_router_attachment_description": "desp",
						"transit_router_attachment_name":        name,
					}),
				),
			},
		},
	})
}

var AlicloudCenTransitRouterVbrAttachmentMap = map[string]string{
	"auto_publish_route_enabled":            CHECKSET,
	"cen_id":                                CHECKSET,
	"dry_run":                               NOSET,
	"resource_type":                         "VBR",
	"route_table_association_enabled":       NOSET,
	"route_table_propagation_enabled":       NOSET,
	"status":                                CHECKSET,
	"transit_router_attachment_description": CHECKSET,
	"transit_router_attachment_name":        CHECKSET,
	"transit_router_id":                     CHECKSET,
	"vbr_id":                                CHECKSET,
	"vbr_owner_id":                          CHECKSET,
}

func AlicloudCenTransitRouterVbrAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {	
	default = "%s"
}
resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level = "REDUCED"
}
resource "alicloud_cen_transit_router" "default" {
cen_id= alicloud_cen_instance.default.id
}
data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^preserved-NODELETING"
}

resource "alicloud_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = 1
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}
`, name)
}
