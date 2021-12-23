package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

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
	transitRouterRouteTableStatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":                               `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"transit_router_route_table_status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":                               `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"transit_router_route_table_status": `"Creating"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"name_regex": `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"name_regex": `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"status": `"Creating"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":                               `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"name_regex":                        `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_name}"`,
			"status":                            `"Active"`,
			"transit_router_route_table_ids":    `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"transit_router_route_table_names":  `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_name}"]`,
			"transit_router_route_table_status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":                               `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}_fake"]`,
			"name_regex":                        `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_name}_fake"`,
			"status":                            `"Creating"`,
			"transit_router_route_table_ids":    `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}_fake"]`,
			"transit_router_route_table_names":  `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_name}_fake"]`,
			"transit_router_route_table_status": `"Creating"`,
		}),
	}
	var existAlicloudCenTransitRouterRouteTablesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"names.#":                                "1",
			"tables.#":                               "1",
			"tables.0.transit_router_route_table_id": CHECKSET,
			"tables.0.transit_router_route_table_description": `desp`,
			"tables.0.transit_router_route_table_name":        CHECKSET,
			"tables.0.transit_router_route_table_type":        CHECKSET,
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
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.CenTRSupportRegions)
	}
	alicloudCenTransitRouterRouteTablesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, transitRouterRouteTableStatusConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudCenTransitRouterRouteTablesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccDataTransitRouterRouteTable-%d"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
cen_id= alicloud_cen_instance.default.id
transit_router_name = var.name
}

resource "alicloud_cen_transit_router_route_table" "default" {
transit_router_id = alicloud_cen_transit_router.default.transit_router_id
transit_router_route_table_description = "desp"
transit_router_route_table_name = var.name
depends_on = [alicloud_cen_transit_router.default]
}

data "alicloud_cen_transit_router_route_tables" "default" {	
transit_router_id = alicloud_cen_transit_router.default.transit_router_id
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
