package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterRouteTableAssociationsDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	rand := acctest.RandIntRange(1, 2999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTableAssociationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_route_table_association.default.transit_router_attachment_id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTableAssociationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_route_table_association.default.transit_router_attachment_id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTableAssociationsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_table_association.default.transit_router_attachment_id}"]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTableAssociationsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_table_association.default.transit_router_attachment_id}"]`,
			"status": `"Creating"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTableAssociationsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_table_association.default.transit_router_attachment_id}"]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTableAssociationsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_table_association.default.transit_router_attachment_id}_fake"]`,
			"status": `"Creating"`,
		}),
	}
	var existAlicloudCenTransitRouterRouteTableAssociationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "1",
			"associations.#": "1",
			"associations.0.transit_router_attachment_id": CHECKSET,
			"associations.0.resource_id":                  CHECKSET,
			"associations.0.resource_type":                CHECKSET,
		}
	}
	var fakeAlicloudCenTransitRouterRouteTableAssociationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCenTransitRouterRouteTableAssociationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_transit_router_route_table_associations.default",
		existMapFunc: existAlicloudCenTransitRouterRouteTableAssociationsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCenTransitRouterRouteTableAssociationsDataSourceNameMapFunc,
	}
	alicloudCenTransitRouterRouteTableAssociationsCheckInfo.dataSourceTestCheck(t, rand, idsConf, statusConf, allConf)
}
func testAccCheckAlicloudCenTransitRouterRouteTableAssociationsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
  default = "tf-testAccTransitRouterRouteTableAssociation-%d"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level = "REDUCED"
  description = "testRouteTableAssociation"
}
resource "alicloud_cen_transit_router" "default" {
  cen_id= alicloud_cen_instance.default.id
  transit_router_description = "testRouteTableAssociation"
  transit_router_name = var.name
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
  vlan_id                    = %d
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
  transit_router_attachment_description = "tf-test"
}
resource "alicloud_cen_transit_router_route_table" "default" {
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
}

resource "alicloud_cen_transit_router_route_table_association" "default" {
  transit_router_attachment_id = alicloud_cen_transit_router_vbr_attachment.default.transit_router_attachment_id
  transit_router_route_table_id = alicloud_cen_transit_router_route_table.default.transit_router_route_table_id
}

data "alicloud_cen_transit_router_route_table_associations" "default" {	
	transit_router_route_table_id = alicloud_cen_transit_router_route_table.default.transit_router_route_table_id 
	%s
}
`, rand, rand, strings.Join(pairs, " \n "))
	return config
}
