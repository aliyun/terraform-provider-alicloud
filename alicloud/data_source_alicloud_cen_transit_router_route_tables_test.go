package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterRouteTablesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}_fake"]`,
		}),
	}
	transitRouterRouteTableIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":                            `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"transit_router_route_table_ids": `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_ids}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":                            `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"transit_router_route_table_ids": `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_ids}_fake"`,
		}),
	}
	transitRouterRouteTableNamesConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":                              `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"transit_router_route_table_names": `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_names}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":                              `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"transit_router_route_table_names": `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_names}_fake"`,
		}),
	}
	transitRouterRouteTableStatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":                               `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"transit_router_route_table_status": `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_status}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":                               `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"transit_router_route_table_status": `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_status}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"status": `"${alicloud_cen_transit_router_route_table.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"status": `"${alicloud_cen_transit_router_route_table.default.status}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":                               `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"name_regex":                        `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_name}"`,
			"status":                            `"${alicloud_cen_transit_router_route_table.default.status}"`,
			"transit_router_route_table_ids":    `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_ids}"`,
			"transit_router_route_table_names":  `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_names}"`,
			"transit_router_route_table_status": `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_status}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":                               `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}_fake"]`,
			"name_regex":                        `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_name}_fake"`,
			"status":                            `"${alicloud_cen_transit_router_route_table.default.status}_fake"`,
			"transit_router_route_table_ids":    `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_ids}_fake"`,
			"transit_router_route_table_names":  `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_names}_fake"`,
			"transit_router_route_table_status": `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_status}_fake"`,
		}),
	}
	var existAlicloudCenTransitRouterRouteTablesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"names.#":                    "1",
			"tables.#":                   "1",
			"tables.0.transit_router_id": CHECKSET,
			"tables.0.transit_router_route_table_description": `desp`,
			"tables.0.transit_router_route_table_name":        CHECKSET,
		}
	}
	var fakeAlicloudCenTransitRouterRouteTablesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCenTransitRouterRouteTablesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_transit_router_route_tables.default",
		existMapFunc: existAlicloudCenTransitRouterRouteTablesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCenTransitRouterRouteTablesDataSourceNameMapFunc,
	}
	alicloudCenTransitRouterRouteTablesCheckInfo.dataSourceTestCheck(t, rand, idsConf, transitRouterRouteTableIdsConf, transitRouterRouteTableNamesConf, transitRouterRouteTableStatusConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccTransitRouterRouteTable-%d"
}

resource "alicloud_cen_transit_router" "default" {
cen_id= "${alicloud_cen_instance.default.id}"
transit_router_name = var.name
}

resource "alicloud_cen_transit_router_route_table" "default" {
transit_router_id = "${alicloud_cen_transitrouter.default.transit_router_id}"
transit_router_route_table_description = "desp"
transit_router_route_table_name = var.name
}

data "alicloud_cen_transit_router_route_tables" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
