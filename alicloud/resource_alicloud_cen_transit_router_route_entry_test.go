package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCenTransitRouterRouteEntry_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cen_transit_router_route_entry.default"
	ra := resourceAttrInit(resourceId, AlicloudCenTransitRouterRouteEntryMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenTransitRouterRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scentransitrouterrouteentry%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCenTransitRouterRouteEntryBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_route_entry_destination_cidr_block": "192.168.1.0/24",
					"transit_router_route_entry_name":                   name,
					"transit_router_route_entry_description":            "test",
					"transit_router_route_entry_next_hop_type":          "Attachment",
					"transit_router_route_entry_next_hop_id":            "${alicloud_cen_transit_router_vbr_attachment.default.transit_router_attachment_id}",
					"transit_router_route_table_id":                     "${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_route_entry_destination_cidr_block": "192.168.1.0/24",
						"transit_router_route_entry_name":                   name,
						"transit_router_route_entry_description":            "test",
						"transit_router_route_entry_next_hop_type":          "Attachment",
						"transit_router_route_entry_next_hop_id":            CHECKSET,
						"transit_router_route_table_id":                     CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_route_entry_description": "desc1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_route_entry_description": "desc1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_route_entry_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_route_entry_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_route_entry_description": "desc",
					"transit_router_route_entry_name":        name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transit_router_route_entry_description": "desc",
						"transit_router_route_entry_name":        name,
					}),
				),
			},
		},
	})
}

var AlicloudCenTransitRouterRouteEntryMap = map[string]string{
	"dry_run":                                NOSET,
	"status":                                 CHECKSET,
	"transit_router_route_entry_description": CHECKSET,
	"transit_router_route_entry_destination_cidr_block": CHECKSET,
	"transit_router_route_entry_name":                   CHECKSET,
	"transit_router_route_entry_next_hop_id":            CHECKSET,
	"transit_router_route_entry_next_hop_type":          CHECKSET,
	"transit_router_route_table_id":                     CHECKSET,
}

func AlicloudCenTransitRouterRouteEntryBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_route_table" "default" {
transit_router_id = alicloud_cen_transit_router.default.transit_router_id
transit_router_route_table_name = var.name
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
resource "alicloud_cen_transit_router_vbr_attachment" "default" {
  cen_id = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
  vbr_id = alicloud_express_connect_virtual_border_router.default.id
  auto_publish_route_enabled = true
  transit_router_attachment_name = var.name
  transit_router_attachment_description = var.name
}
`, name)
}
