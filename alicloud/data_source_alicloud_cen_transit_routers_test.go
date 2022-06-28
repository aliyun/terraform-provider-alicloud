package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRoutersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
			"ids":    `["${alicloud_cen_transit_router.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"cen_id": `"${alicloud_cen_instance.default.id}"`,
			"ids":    `["${alicloud_cen_transit_router.default.id}_fake"]`,
		}),
	}
	transitRouteridConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"cen_id":            `"${alicloud_cen_instance.default.id}"`,
			"transit_router_id": `"${alicloud_cen_transit_router.default.transit_router_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"cen_id":            `"${alicloud_cen_instance.default.id}"`,
			"transit_router_id": `"${alicloud_cen_transit_router.default.transit_router_id}_fake"`,
		}),
	}
	transitRouteridsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"cen_id":             `"${alicloud_cen_instance.default.id}"`,
			"transit_router_ids": `["${alicloud_cen_transit_router.default.transit_router_id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"cen_id":             `"${alicloud_cen_instance.default.id}"`,
			"transit_router_ids": `["${alicloud_cen_transit_router.default.transit_router_id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"cen_id":     `"${alicloud_cen_instance.default.id}"`,
			"name_regex": `"${alicloud_cen_transit_router.default.transit_router_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"cen_id":     `"${alicloud_cen_instance.default.id}"`,
			"name_regex": `"${alicloud_cen_transit_router.default.transit_router_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"cen_id":            `"${alicloud_cen_instance.default.id}"`,
			"transit_router_id": `"${alicloud_cen_transit_router.default.transit_router_id}"`,
			"status":            `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"cen_id":            `"${alicloud_cen_instance.default.id}"`,
			"transit_router_id": `"${alicloud_cen_transit_router.default.transit_router_id}"`,
			"status":            `"Creating"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"cen_id":             `"${alicloud_cen_instance.default.id}"`,
			"ids":                `["${alicloud_cen_transit_router.default.id}"]`,
			"transit_router_id":  `"${alicloud_cen_transit_router.default.transit_router_id}"`,
			"transit_router_ids": `["${alicloud_cen_transit_router.default.transit_router_id}"]`,
			"name_regex":         `"${alicloud_cen_transit_router.default.transit_router_name}"`,
			"status":             `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"cen_id":             `"${alicloud_cen_instance.default.id}"`,
			"transit_router_id":  `"${alicloud_cen_transit_router.default.transit_router_id}"`,
			"ids":                `["${alicloud_cen_transit_router.default.id}_fake"]`,
			"transit_router_ids": `["${alicloud_cen_transit_router.default.transit_router_id}_fake"]`,
			"name_regex":         `"${alicloud_cen_transit_router.default.transit_router_name}_fake"`,
			"status":             `"Creating"`,
		}),
	}
	var existAlicloudCenTransitRoutersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"transit_routers.#":                   "1",
			"transit_routers.0.cen_id":            CHECKSET,
			"transit_routers.0.id":                CHECKSET,
			"transit_routers.0.transit_router_id": CHECKSET,
			"transit_routers.0.transit_router_description": "desd",
			"transit_routers.0.ali_uid":                    CHECKSET,
			"transit_routers.0.status":                     CHECKSET,
			"transit_routers.0.type":                       CHECKSET,
			"transit_routers.0.xgw_vip":                    "",
			"transit_routers.0.transit_router_name":        CHECKSET,
		}
	}
	var fakeAlicloudCenTransitRoutersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCenTransitRoutersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_transit_routers.default",
		existMapFunc: existAlicloudCenTransitRoutersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCenTransitRoutersDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.CenTRSupportRegions)
	}
	alicloudCenTransitRoutersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, transitRouteridConf, transitRouteridsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudCenTransitRoutersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

	variable "name" {	
		default = "tf-testAccDataTransitRouter-%d"
	}

	resource "alicloud_cen_instance" "default" {
		cen_instance_name = var.name
	}

	resource "alicloud_cen_transit_router" "default" {
		cen_id = alicloud_cen_instance.default.id
		transit_router_description = "desd"
		transit_router_name = var.name
	}

	data "alicloud_cen_transit_routers" "default" {
		%s
	}
	`, rand, strings.Join(pairs, " \n "))
	return config
}
