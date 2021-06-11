package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRoutersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router.default.transit_router_id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router.default.transit_router_id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router.default.transit_router_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router.default.transit_router_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router.default.transit_router_id}"]`,
			"status": `"${alicloud_cen_transit_router.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_cen_transit_router.default.transit_router_id}"]`,
			"status": `"${alicloud_cen_transit_router.default.status}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cen_transit_router.default.transit_router_id}"]`,
			"name_regex": `"${alicloud_cen_transit_router.default.transit_router_name}"`,
			"status":     `"${alicloud_cen_transit_router.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRoutersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cen_transit_router.default.transit_router_id}_fake"]`,
			"name_regex": `"${alicloud_cen_transit_router.default.transit_router_name}_fake"`,
			"status":     `"${alicloud_cen_transit_router.default.status}_fake"`,
		}),
	}
	var existAlicloudCenTransitRoutersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"names.#":                              "1",
			"routers.#":                            "1",
			"routers.0.cen_id":                     CHECKSET,
			"routers.0.transit_router_description": `desd`,
			"routers.0.transit_router_name":        CHECKSET,
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
	alicloudCenTransitRoutersCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudCenTransitRoutersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccTransitRouter-%d"
}

resource "alicloud_cen_instance" "default" {
cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "default" {
cen_id = "${alicloud_cen_ceninstance.default.cen_id}"
transit_router_description = "desd"
transit_router_name = var.name
}

data "alicloud_cen_transit_routers" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
