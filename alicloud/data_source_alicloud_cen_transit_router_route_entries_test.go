package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterRouteEntriesDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	rand := acctest.RandIntRange(1, 2999)
	transitRouterRouteEntryIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteEntriesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_route_entry.default.transit_router_route_entry_id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteEntriesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_route_entry.default.transit_router_route_entry_id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteEntriesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_entry.default.transit_router_route_entry_id}"]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteEntriesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_entry.default.transit_router_route_entry_id}"]`,
			"status": `"Creating"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteEntriesDataSourceName(rand, map[string]string{
			"ids":                              `["${alicloud_cen_transit_router_route_entry.default.transit_router_route_entry_id}"]`,
			"name_regex":                       `"${alicloud_cen_transit_router_route_entry.default.transit_router_route_entry_name}"`,
			"status":                           `"Active"`,
			"transit_router_route_entry_ids":   `["${alicloud_cen_transit_router_route_entry.default.transit_router_route_entry_id}"]`,
			"transit_router_route_entry_names": `["${alicloud_cen_transit_router_route_entry.default.transit_router_route_entry_name}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteEntriesDataSourceName(rand, map[string]string{
			"ids":                              `["${alicloud_cen_transit_router_route_entry.default.transit_router_route_entry_id}_fake"]`,
			"name_regex":                       `"${alicloud_cen_transit_router_route_entry.default.transit_router_route_entry_name}_fake"`,
			"status":                           `"Creating"`,
			"transit_router_route_entry_ids":   `["${alicloud_cen_transit_router_route_entry.default.transit_router_route_entry_id}_fake"]`,
			"transit_router_route_entry_names": `["${alicloud_cen_transit_router_route_entry.default.transit_router_route_entry_name}_fake"]`,
		}),
	}
	var existAlicloudCenTransitRouterRouteEntriesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "1",
			"names.#":   "1",
			"entries.#": "1",
			"entries.0.transit_router_route_entry_description":            `desc`,
			"entries.0.transit_router_route_entry_destination_cidr_block": `192.168.1.0/24`,
			"entries.0.transit_router_route_entry_name":                   CHECKSET,
			"entries.0.transit_router_route_entry_next_hop_id":            CHECKSET,
			"entries.0.transit_router_route_entry_next_hop_type":          `Attachment`,
			"entries.0.transit_router_route_entry_id":                     CHECKSET,
		}
	}
	var fakeAlicloudCenTransitRouterRouteEntriesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCenTransitRouterRouteEntriesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_transit_router_route_entries.default",
		existMapFunc: existAlicloudCenTransitRouterRouteEntriesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCenTransitRouterRouteEntriesDataSourceNameMapFunc,
	}
	alicloudCenTransitRouterRouteEntriesCheckInfo.dataSourceTestCheck(t, rand, transitRouterRouteEntryIdsConf, statusConf, allConf)
}
func testAccCheckAlicloudCenTransitRouterRouteEntriesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccDataTransitRouterRouteEntry-%d"
}
resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
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
  auto_publish_route_enabled = true
  cen_id = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
  vbr_id = alicloud_express_connect_virtual_border_router.default.id
  transit_router_attachment_description = "desp"
  transit_router_attachment_name = var.name
}

resource "alicloud_cen_transit_router_route_table" "default" {
	transit_router_id = alicloud_cen_transit_router_vbr_attachment.default.transit_router_id
	transit_router_route_table_name = var.name
}

resource "alicloud_cen_transit_router_route_entry" "default" {
	transit_router_route_entry_description = "desc"
	transit_router_route_entry_destination_cidr_block = "192.168.1.0/24"
	transit_router_route_entry_name = var.name
	transit_router_route_entry_next_hop_id = alicloud_cen_transit_router_vbr_attachment.default.transit_router_attachment_id
	transit_router_route_entry_next_hop_type = "Attachment"
	transit_router_route_table_id = alicloud_cen_transit_router_route_table.default.transit_router_route_table_id
}

data "alicloud_cen_transit_router_route_entries" "default" {
transit_router_route_table_id = alicloud_cen_transit_router_route_table.default.transit_router_route_table_id
	%s
}
`, rand, rand, strings.Join(pairs, " \n "))
	return config
}
