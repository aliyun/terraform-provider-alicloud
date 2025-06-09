package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCenTransitRouterRouteTablesDataSource_basic0(t *testing.T) {
	rand := acctest.RandInt()

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_name}_fake"`,
		}),
	}

	transitRouterRouteTableTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"transit_router_route_table_type": `"Custom"`,
		}),
		fakeConfig: testAccCheckAliCloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":                             `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}_fake"]`,
			"transit_router_route_table_type": `"Custom"`,
		}),
	}

	transitRouterRouteTableStatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":                               `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"transit_router_route_table_status": `"Active"`,
		}),
		fakeConfig: testAccCheckAliCloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"transit_router_route_table_status": `"Creating"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAliCloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"status": `"Creating"`,
		}),
	}

	transitRouterRouteTableIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"transit_router_route_table_ids": `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"transit_router_route_table_ids": `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}_fake"]`,
		}),
	}

	transitRouterRouteTableNamesConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"transit_router_route_table_names": `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_name}"]`,
		}),
		fakeConfig: testAccCheckAliCloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"transit_router_route_table_names": `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_name}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":                               `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"name_regex":                        `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_name}"`,
			"transit_router_route_table_type":   `"Custom"`,
			"transit_router_route_table_status": `"Active"`,
			"status":                            `"Active"`,
			"transit_router_route_table_ids":    `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}"]`,
			"transit_router_route_table_names":  `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_name}"]`,
		}),
		fakeConfig: testAccCheckAliCloudCenTransitRouterRouteTablesDataSourceName(rand, map[string]string{
			"ids":                               `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}_fake"]`,
			"name_regex":                        `"${alicloud_cen_transit_router_route_table.default.transit_router_route_table_name}_fake"`,
			"transit_router_route_table_type":   `"System"`,
			"transit_router_route_table_status": `"Creating"`,
			"status":                            `"Creating"`,
			"transit_router_route_table_ids":    `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_id}_fake"]`,
			"transit_router_route_table_names":  `["${alicloud_cen_transit_router_route_table.default.transit_router_route_table_name}_fake"]`,
		}),
	}

	var existAliCloudCenTransitRouterRouteTablesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"names.#":                                "1",
			"tables.#":                               "1",
			"tables.0.id":                            CHECKSET,
			"tables.0.transit_router_route_table_id": CHECKSET,
			"tables.0.transit_router_route_table_type":        CHECKSET,
			"tables.0.transit_router_route_table_name":        CHECKSET,
			"tables.0.transit_router_route_table_description": CHECKSET,
			"tables.0.status": CHECKSET,
		}
	}

	var fakeAliCloudCenTransitRouterRouteTablesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"tables.#": "0",
		}
	}

	var aliCloudCenTransitRouterRouteTablesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_transit_router_route_tables.default",
		existMapFunc: existAliCloudCenTransitRouterRouteTablesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudCenTransitRouterRouteTablesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.CenTRSupportRegions)
	}

	aliCloudCenTransitRouterRouteTablesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, transitRouterRouteTableTypeConf, transitRouterRouteTableStatusConf, statusConf, transitRouterRouteTableIdsConf, transitRouterRouteTableNamesConf, allConf)
}

func testAccCheckAliCloudCenTransitRouterRouteTablesDataSourceName(rand int, attrMap map[string]string) string {
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
  		protection_level  = "REDUCED"
	}

	resource "alicloud_cen_transit_router" "default" {
  		cen_id              = alicloud_cen_instance.default.id
  		transit_router_name = var.name
	}

	resource "alicloud_cen_transit_router_route_table" "default" {
  		transit_router_id                      = alicloud_cen_transit_router.default.transit_router_id
  		transit_router_route_table_name        = var.name
  		transit_router_route_table_description = var.name
	}

	data "alicloud_cen_transit_router_route_tables" "default" {
  		transit_router_id = alicloud_cen_transit_router.default.transit_router_id
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
