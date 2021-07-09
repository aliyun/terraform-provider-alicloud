package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

/**
This resource has buried point data.
VBR is buried point data.
*/
func SkipTestAccAlicloudCenTransitRouterRouteEntriesDataSource(t *testing.T) {
	rand := acctest.RandInt()
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

resource "alicloud_cen_transit_router_vbr_attachment" "default" {
  auto_publish_route_enabled = true
  cen_id = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
  vbr_id = "vbr-j6cd9pm9y6d6e20atoi6w"
  transit_router_attachment_description = "desp"
  transit_router_attachment_name = var.name
}

resource "alicloud_cen_transit_router_route_table" "default" {
transit_router_id = alicloud_cen_transit_router.default.transit_router_id
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
`, rand, strings.Join(pairs, " \n "))
	return config
}
