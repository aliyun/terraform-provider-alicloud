package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterRouteTableAssociationsDataSource(t *testing.T) {
	rand := acctest.RandInt()
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
			"status": `"${alicloud_cen_transit_router_route_table_association.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTableAssociationsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_table_association.default.transit_router_attachment_id}"]`,
			"status": `"${alicloud_cen_transit_router_route_table_association.default.status}"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterRouteTableAssociationsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_table_association.default.transit_router_attachment_id}"]`,
			"status": `"${alicloud_cen_transit_router_route_table_association.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterRouteTableAssociationsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router_route_table_association.default.transit_router_attachment_id}_fake"]`,
			"status": `"${alicloud_cen_transit_router_route_table_association.default.status}"`,
		}),
	}
	var existAlicloudCenTransitRouterRouteTableAssociationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "1",
			"associations.#": "1",
			"associations.0.transit_router_attachment_id":  CHECKSET,
			"associations.0.transit_router_route_table_id": CHECKSET,
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
  vbr_id = "vbr-j6cd9pm9y6d6e20atoi6w"
  auto_publish_route_enabled = true
  transit_router_attachment_name = "tf-test"
  transit_router_attachment_description = "tf-test"
}

resource "alicloud_cen_transit_router_route_table_association" "default" {
transit_router_attachment_id = "${alicloud_cen_transit_router_vbr_attachment.default.id}"
transit_router_route_table_id = "${alicloud_cen_transit_router_route_table.default.id}"
}

data "alicloud_cen_transit_router_route_table_associations" "default" {	
	transit_router_route_table_id = "${alicloud_cen_transit_router_route_table.default.id}" 
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
