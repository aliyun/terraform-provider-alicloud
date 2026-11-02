package cen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCenTransitRouterCidrsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenTransitRouterCidrsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_cidr.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCenTransitRouterCidrsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_cidr.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenTransitRouterCidrsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_cidr.default.transit_router_cidr_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCenTransitRouterCidrsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_cidr.default.transit_router_cidr_name}_fake"`,
		}),
	}
	transitRouterCidrIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenTransitRouterCidrsDataSourceName(rand, map[string]string{
			"transit_router_cidr_id": `"${alicloud_cen_transit_router_cidr.default.transit_router_cidr_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudCenTransitRouterCidrsDataSourceName(rand, map[string]string{
			"transit_router_cidr_id": `"${alicloud_cen_transit_router_cidr.default.transit_router_cidr_id}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenTransitRouterCidrsDataSourceName(rand, map[string]string{
			"ids":                    `["${alicloud_cen_transit_router_cidr.default.id}"]`,
			"name_regex":             `"${alicloud_cen_transit_router_cidr.default.transit_router_cidr_name}"`,
			"transit_router_cidr_id": `"${alicloud_cen_transit_router_cidr.default.transit_router_cidr_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudCenTransitRouterCidrsDataSourceName(rand, map[string]string{
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
			"cidrs.0.id":                       acctest.CHECKSET,
			"cidrs.0.transit_router_id":        acctest.CHECKSET,
			"cidrs.0.transit_router_cidr_id":   acctest.CHECKSET,
			"cidrs.0.cidr":                     "192.168.0.0/16",
			"cidrs.0.transit_router_cidr_name": acctest.CHECKSET,
			"cidrs.0.description":              acctest.CHECKSET,
			"cidrs.0.publish_cidr_route":       "false",
			"cidrs.0.family":                   acctest.CHECKSET,
		}
	}
	var fakeAlicloudCenTransitRouterCidrsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"cidrs.#": "0",
		}
	}
	var alicloudCenTransitRouterCidrsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cen_transit_router_cidrs.default",
		ExistMapFunc: existAlicloudCenTransitRouterCidrsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCenTransitRouterCidrsDataSourceNameMapFunc,
	}
	alicloudCenTransitRouterCidrsCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, transitRouterCidrIdConf, allConf)
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
