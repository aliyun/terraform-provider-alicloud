package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterRouteEntriesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	transitRouterRouteEntryIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteEntriesDataSourceName(rand, map[string]string{
			"transit_router_route_table_id":  `"${alicloud_cen_transit_router_route_table.default.id}"`,
			"transit_router_route_entry_ids": `["${alicloud_cen_transit_router_route_entry.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteEntriesDataSourceName(rand, map[string]string{
			"transit_router_route_table_id":  `"${alicloud_cen_transit_router_route_table.default.id}_fake"`,
			"transit_router_route_entry_ids": `["${alicloud_cen_transit_router_route_entry.default.id}_fake"]`,
		}),
	}
	transitRouterRouteEntryNamesConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteEntriesDataSourceName(rand, map[string]string{
			"transit_router_route_table_id":    `["${alicloud_cen_transit_router_route_table.default.id}"]`,
			"transit_router_route_entry_names": `["${alicloud_cen_transit_router_route_entry.default.name}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteEntriesDataSourceName(rand, map[string]string{
			"transit_router_route_table_id":    `["${alicloud_cen_transit_router_route_table.default.id}_fake"]`,
			"transit_router_route_entry_names": `["${alicloud_cen_transit_router_route_entry.default.name}_fake"]`,
		}),
	}
	transitRouterRouteEntryStatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteEntriesDataSourceName(rand, map[string]string{
			"transit_router_route_table_id":     `"${alicloud_cen_transit_router_route_table.default.id}"`,
			"transit_router_route_entry_status": `"${alicloud_cen_transit_router_route_entry.default.transit_router_route_entry_status}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteEntriesDataSourceName(rand, map[string]string{
			"transit_router_route_table_id":     `"${alicloud_cen_transit_router_route_table.default.id}_fake"`,
			"transit_router_route_entry_status": `"${alicloud_cen_transit_router_route_entry.default.transit_router_route_entry_status}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteEntriesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_entry.default.id}"]`,
			"status": `"${alicloud_cen_transit_router_route_entry.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteEntriesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_entry.default.id}"]`,
			"status": `"${alicloud_cen_transit_router_route_entry.default.status}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteEntriesDataSourceName(rand, map[string]string{
			"transit_router_route_table_id":     `"${alicloud_cen_transit_router_route_table.default.id}"`,
			"ids":                               `["${alicloud_cen_transit_router_route_entry.default.id}"]`,
			"name_regex":                        `"${alicloud_cen_transit_router_route_entry.default.name}"`,
			"status":                            `"${alicloud_cen_transit_router_route_entry.default.status}"`,
			"transit_router_route_entry_ids":    `"${alicloud_cen_transit_router_route_entry.default.transit_router_route_entry_ids}"`,
			"transit_router_route_entry_names":  `"${alicloud_cen_transit_router_route_entry.default.transit_router_route_entry_names}"`,
			"transit_router_route_entry_status": `"${alicloud_cen_transit_router_route_entry.default.transit_router_route_entry_status}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteEntriesDataSourceName(rand, map[string]string{
			"ids":                               `["${alicloud_cen_transit_router_route_entry.default.id}_fake"]`,
			"name_regex":                        `"${alicloud_cen_transit_router_route_entry.default.name}_fake"`,
			"status":                            `"${alicloud_cen_transit_router_route_entry.default.status}"`,
			"transit_router_route_table_id":     `"${alicloud_cen_transit_router_route_table.default.id}"`,
			"transit_router_route_entry_ids":    `"${alicloud_cen_transit_router_route_entry.default.ids}_fake"`,
			"transit_router_route_entry_names":  `"${alicloud_cen_transit_router_route_entry.default.names}_fake"`,
			"transit_router_route_entry_status": `"${alicloud_cen_transit_router_route_entry.default.status}"`,
		}),
	}
	var existAlicloudCenTransitRouterRouteEntriesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"transit_router_route_entries.#":                                                   "1",
			"transit_router_route_entries.0.transit_router_route_entry_description":            `desc`,
			"transit_router_route_entries.0.transit_router_route_entry_destination_cidr_block": `192.168.0.1/24`,
			"transit_router_route_entries.0.transit_router_route_entry_name":                   CHECKSET,
			"transit_router_route_entries.0.transit_router_route_entry_next_hop_id":            CHECKSET,
			"transit_router_route_entries.0.transit_router_route_entry_next_hop_type":          `Attachment`,
			"transit_router_route_entries.0.transit_router_route_table_id":                     CHECKSET,
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
	alicloudCenTransitRouterRouteEntriesCheckInfo.dataSourceTestCheck(t, rand, transitRouterRouteEntryIdsConf, transitRouterRouteEntryNamesConf, transitRouterRouteEntryStatusConf, statusConf, allConf)
}
func testAccCheckAlicloudCenTransitRouterRouteEntriesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccTransitRouterRouteEntry-%d"
}
resource "alicloud_cen_instance" "default" {
  cen_instance_name = "${var.name}"
  protection_level = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = "${alicloud_cen_instance.default.id}"
}

resource "alicloud_cen_transit_router_route_table" "default" {
transit_router_id = "${alicloud_cen_transit_router.default.id}"
transit_router_route_table_name = "${var.name}"
}
resource "alicloud_cen_transit_router_vbr_attachment" "default" {
  cen_id = "${alicloud_cen_instance.default.id}"
  transit_router_id = "${alicloud_cen_transit_router.default.id}"
  vbr_id = "vbr-j6cd9pm9y6d6e20atoi6w"
  auto_publish_route_enabled = true
  transit_router_attachment_name = "name"
  transit_router_attachment_description = "name"
}
resource "alicloud_cen_transit_router_route_entry" "default" {
transit_router_route_entry_description = "desc"
transit_router_route_entry_destination_cidr_block = "192.168.0.1/24"
transit_router_route_entry_name = "${var.name}"
transit_router_route_entry_next_hop_id = "vbr-j6cd9pm9y6d6e20atoi6w"
transit_router_route_entry_next_hop_type = "Attachment"
transit_router_route_table_id = "${alicloud_cen_transit_router_route_table.default.id}"
}

data "alicloud_cen_transit_router_route_entries" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
