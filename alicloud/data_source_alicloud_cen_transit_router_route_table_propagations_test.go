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
func SkipTestAccAlicloudCenTransitRouterRouteTablePropagationsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTablePropagationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_route_table_propagation.default.transit_router_attachment_id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTablePropagationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_route_table_propagation.default.transit_router_attachment_id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTablePropagationsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_table_propagation.default.transit_router_attachment_id}"]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTablePropagationsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_table_propagation.default.transit_router_attachment_id}"]`,
			"status": `"Creating"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTablePropagationsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_table_propagation.default.transit_router_attachment_id}"]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTablePropagationsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_table_propagation.default.transit_router_attachment_id}_fake"]`,
			"status": `"Creating"`,
		}),
	}
	var existAlicloudCenTransitRouterRouteTablePropagationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "1",
			"propagations.#": "1",
			"propagations.0.transit_router_attachment_id": CHECKSET,
			"propagations.0.resource_id":                  CHECKSET,
			"propagations.0.resource_type":                "VBR",
		}
	}
	var fakeAlicloudCenTransitRouterRouteTablePropagationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCenTransitRouterRouteTablePropagationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_transit_router_route_table_propagations.default",
		existMapFunc: existAlicloudCenTransitRouterRouteTablePropagationsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCenTransitRouterRouteTablePropagationsDataSourceNameMapFunc,
	}
	alicloudCenTransitRouterRouteTablePropagationsCheckInfo.dataSourceTestCheck(t, rand, idsConf, statusConf, allConf)
}
func testAccCheckAlicloudCenTransitRouterRouteTablePropagationsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccDataTransitRouterRouteTablePropagation-%d"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  protection_level = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
cen_id= alicloud_cen_instance.default.id
}

resource "alicloud_cen_transit_router_vbr_attachment" "default" {
  cen_id = alicloud_cen_instance.default.id
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
  vbr_id = "vbr-j6cxhs879lwxzosc4h0lv"
  auto_publish_route_enabled = true
  transit_router_attachment_name = "tf-test"
  transit_router_attachment_description = "tf-test"
}
resource "alicloud_cen_transit_router_route_table" "default" {
  transit_router_id = alicloud_cen_transit_router.default.transit_router_id
  transit_router_route_table_name = "testRouteTable"
}

resource "alicloud_cen_transit_router_route_table_propagation" "default" {
transit_router_attachment_id = alicloud_cen_transit_router_vbr_attachment.default.transit_router_attachment_id
transit_router_route_table_id = alicloud_cen_transit_router_route_table.default.transit_router_route_table_id
}

data "alicloud_cen_transit_router_route_table_propagations" "default" {
transit_router_route_table_id = alicloud_cen_transit_router_route_table.default.transit_router_route_table_id
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
