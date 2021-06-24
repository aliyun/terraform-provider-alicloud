package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterRouteTablePropagationsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTablePropagationsDataSourceName(rand, map[string]string{
			"ids":                           `["${alicloud_cen_transit_router_route_table_propagation.default.transit_router_attachment_id}"]`,
			"transit_router_route_table_id": `"${alicloud_cen_transit_router_route_table.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTablePropagationsDataSourceName(rand, map[string]string{
			"ids":                           `["${alicloud_cen_transit_router_route_table_propagation.default.transit_router_attachment_id}_fake"]`,
			"transit_router_route_table_id": `"${alicloud_cen_transit_router_route_table.default.id}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTablePropagationsDataSourceName(rand, map[string]string{
			"ids":                           `["${alicloud_cen_transit_router_route_table_propagation.default.transit_router_attachment_id}"]`,
			"status":                        `"${alicloud_cen_transit_router_route_table_propagation.default.status}"`,
			"transit_router_route_table_id": `"${alicloud_cen_transit_router_route_table.default.id}_fake"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTablePropagationsDataSourceName(rand, map[string]string{
			"ids":                           `["${alicloud_cen_transit_router_route_table_propagation.default.transit_router_attachment_id}"]`,
			"status":                        `"${alicloud_cen_transit_router_route_table_propagation.default.status}_fake"`,
			"transit_router_route_table_id": `"${alicloud_cen_transit_router_route_table.default.id}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTablePropagationsDataSourceName(rand, map[string]string{
			"ids":                           `["${alicloud_cen_transit_router_route_table_propagation.default.transit_router_attachment_id}"]`,
			"status":                        `"${alicloud_cen_transit_router_route_table_propagation.default.status}"`,
			"transit_router_route_table_id": `"${alicloud_cen_transit_router_route_table.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTablePropagationsDataSourceName(rand, map[string]string{
			"ids":                           `["${alicloud_cen_transit_router_route_table_propagation.default.transit_router_attachment_id}_fake"]`,
			"status":                        `"${alicloud_cen_transit_router_route_table_propagation.default.status}"`,
			"transit_router_route_table_id": `"${alicloud_cen_transit_router_route_table.default.id}_fake"`,
		}),
	}
	var existAlicloudCenTransitRouterRouteTablePropagationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "1",
			"propagations.#": "1",
			"propagations.0.transit_router_attachment_id":  CHECKSET,
			"propagations.0.transit_router_route_table_id": CHECKSET,
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
	default = "tf-testAccTransitRouterRouteTablePropagation-%d"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = "${var.name}"
  protection_level = "REDUCED"
}

resource "alicloud_cen_transit_router" "default" {
cen_id= "${alicloud_cen_instance.default.id}"
region_id = "cn-hongkong"
}

resource "alicloud_cen_transit_router_route_table" "default" {
  transit_router_id = "${alicloud_cen_transit_router.default.id}"
}

resource "alicloud_cen_transit_router_vbr_attachment" "default" {
  cen_id = "${alicloud_cen_instance.default.id}"
  transit_router_id = "${alicloud_cen_transit_router.default.id}"
  vbr_id = "vbr-j6cxhs879lwxzosc4h0lv"
  auto_publish_route_enabled = true
  transit_router_attachment_name = "tf-test"
  transit_router_attachment_description = "tf-test"
}

resource "alicloud_cen_transit_router_route_table_propagation" "default" {
transit_router_attachment_id = "${alicloud_cen_transit_router_vbr_attachment.default.id}"
transit_router_route_table_id = "${alicloud_cen_transit_router_route_table.default.id}"
}

data "alicloud_cen_transit_router_route_table_propagations" "default" {
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
