package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterCidrsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterCidrsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_cidr.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterCidrsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_cidr.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterCidrsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_cidr.default.transit_router_cidr_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterCidrsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_cidr.default.transit_router_cidr_name}_fake"`,
		}),
	}
	transitRouterCidrIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterCidrsDataSourceName(rand, map[string]string{
			"transit_router_cidr_id": `"${alicloud_cen_transit_router_cidr.default.transit_router_cidr_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterCidrsDataSourceName(rand, map[string]string{
			"transit_router_cidr_id": `"${alicloud_cen_transit_router_cidr.default.transit_router_cidr_id}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterCidrsDataSourceName(rand, map[string]string{
			"ids":                    `["${alicloud_cen_transit_router_cidr.default.id}"]`,
			"name_regex":             `"${alicloud_cen_transit_router_cidr.default.transit_router_cidr_name}"`,
			"transit_router_cidr_id": `"${alicloud_cen_transit_router_cidr.default.transit_router_cidr_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterCidrsDataSourceName(rand, map[string]string{
			"ids":                    `["${alicloud_cen_transit_router_cidr.default.id}_fake"]`,
			"name_regex":             `"${alicloud_cen_transit_router_cidr.default.transit_router_cidr_name}_fake"`,
			"transit_router_cidr_id": `"${alicloud_cen_transit_router_cidr.default.transit_router_cidr_id}_fake"`,
		}),
	}
	var existAlicloudCenTransitRouterCidrsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"names.#":                          "1",
			"cidrs.#":                          "1",
			"cidrs.0.id":                       CHECKSET,
			"cidrs.0.transit_router_id":        CHECKSET,
			"cidrs.0.transit_router_cidr_id":   CHECKSET,
			"cidrs.0.cidr":                     "192.168.0.0/16",
			"cidrs.0.transit_router_cidr_name": CHECKSET,
			"cidrs.0.description":              CHECKSET,
			"cidrs.0.publish_cidr_route":       "false",
			"cidrs.0.family":                   CHECKSET,
		}
	}
	var fakeAlicloudCenTransitRouterCidrsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"cidrs.#": "0",
		}
	}
	var alicloudCenTransitRouterCidrsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_transit_router_cidrs.default",
		existMapFunc: existAlicloudCenTransitRouterCidrsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCenTransitRouterCidrsDataSourceNameMapFunc,
	}
	alicloudCenTransitRouterCidrsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, transitRouterCidrIdConf, allConf)
}

func testAccCheckAlicloudCenTransitRouterCidrsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAcc-%d"
	}

	resource "alicloud_cen_instance" "default" {
  		cen_instance_name = var.name
	}

	resource "alicloud_cen_transit_router" "default" {
		cen_id = alicloud_cen_instance.default.id
	}

	resource "alicloud_cen_transit_router_cidr" "default" {
		transit_router_id        = alicloud_cen_transit_router.default.transit_router_id
  		cidr                     = "192.168.0.0/16"
  		transit_router_cidr_name = var.name
  		description              = var.name
  		publish_cidr_route       = false
	}

	data "alicloud_cen_transit_router_cidrs" "default" {
		transit_router_id = alicloud_cen_transit_router_cidr.default.transit_router_id
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
