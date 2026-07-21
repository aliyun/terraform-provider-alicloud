package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterPrefixListAssociationsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterPrefixListAssociationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_prefix_list_association.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterPrefixListAssociationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_prefix_list_association.default.id}_fake"]`,
		}),
	}
	prefixListIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterPrefixListAssociationsDataSourceName(rand, map[string]string{
			"prefix_list_id": `"${alicloud_cen_transit_router_prefix_list_association.default.prefix_list_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterPrefixListAssociationsDataSourceName(rand, map[string]string{
			"prefix_list_id": `"${alicloud_cen_transit_router_prefix_list_association.default.prefix_list_id}_fake"`,
		}),
	}
	ownerUidConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterPrefixListAssociationsDataSourceName(rand, map[string]string{
			"owner_uid": `"${data.alicloud_account.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterPrefixListAssociationsDataSourceName(rand, map[string]string{
			"owner_uid": `"${data.alicloud_account.default.id}0"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterPrefixListAssociationsDataSourceName(rand, map[string]string{
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterPrefixListAssociationsDataSourceName(rand, map[string]string{
			"status": `"Updating"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterPrefixListAssociationsDataSourceName(rand, map[string]string{
			"prefix_list_id": `"${alicloud_cen_transit_router_prefix_list_association.default.prefix_list_id}"`,
			"owner_uid":      `"${data.alicloud_account.default.id}"`,
			"status":         `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTransitRouterPrefixListAssociationsDataSourceName(rand, map[string]string{
			"prefix_list_id": `"${alicloud_cen_transit_router_prefix_list_association.default.prefix_list_id}_fake"`,
			"owner_uid":      `"${data.alicloud_account.default.id}0"`,
			"status":         `"Updating"`,
		}),
	}
	var existAlicloudCenTransitRouterPrefixListAssociationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"associations.#":                         "1",
			"associations.0.id":                      CHECKSET,
			"associations.0.prefix_list_id":          CHECKSET,
			"associations.0.transit_router_id":       CHECKSET,
			"associations.0.transit_router_table_id": CHECKSET,
			"associations.0.next_hop":                "BlackHole",
			"associations.0.next_hop_type":           "BlackHole",
			"associations.0.next_hop_instance_id":    CHECKSET,
			"associations.0.owner_uid":               CHECKSET,
			"associations.0.status":                  "Active",
		}
	}
	var fakeAlicloudCenTransitRouterPrefixListAssociationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"associations.#": "0",
		}
	}
	var alicloudCenTransitRouterPrefixListAssociationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_transit_router_prefix_list_associations.default",
		existMapFunc: existAlicloudCenTransitRouterPrefixListAssociationsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCenTransitRouterPrefixListAssociationsDataSourceNameMapFunc,
	}
	alicloudCenTransitRouterPrefixListAssociationsCheckInfo.dataSourceTestCheck(t, rand, idsConf, prefixListIdConf, ownerUidConf, statusConf, allConf)
}

func testAccCheckAlicloudCenTransitRouterPrefixListAssociationsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf_testAcc_%d"
	}

	data "alicloud_account" "default" {
	}

	resource "alicloud_vpc_prefix_list" "default" {
  		entrys {
    		cidr = "192.168.0.0/16"
  		}
	}

	resource "alicloud_cen_instance" "default" {
  		cen_instance_name = var.name
	}

	resource "alicloud_cen_transit_router" "default" {
  		cen_id = alicloud_cen_instance.default.id
	}

	resource "alicloud_cen_transit_router_route_table" "default" {
  		transit_router_id = alicloud_cen_transit_router.default.transit_router_id
	}

	resource "alicloud_cen_transit_router_prefix_list_association" "default" {
  		prefix_list_id          = alicloud_vpc_prefix_list.default.id
  		transit_router_id       = alicloud_cen_transit_router.default.transit_router_id
  		transit_router_table_id = alicloud_cen_transit_router_route_table.default.transit_router_route_table_id
  		next_hop                = "BlackHole"
  		next_hop_type           = "BlackHole"
  		owner_uid               = data.alicloud_account.default.id
	}

	data "alicloud_cen_transit_router_prefix_list_associations" "default" {
  		transit_router_id       = alicloud_cen_transit_router_prefix_list_association.default.transit_router_id
  		transit_router_table_id = alicloud_cen_transit_router_prefix_list_association.default.transit_router_table_id
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
