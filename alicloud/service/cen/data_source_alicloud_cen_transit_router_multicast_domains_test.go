package cen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCenTransitRouterMulticastDomainsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_multicast_domain.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cen_transit_router_multicast_domain.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_name}_fake"`,
		}),
	}
	transitRouterMulticastDomainIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainsDataSourceName(rand, map[string]string{
			"transit_router_multicast_domain_id": `"${alicloud_cen_transit_router_multicast_domain.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainsDataSourceName(rand, map[string]string{
			"transit_router_multicast_domain_id": `"${alicloud_cen_transit_router_multicast_domain.default.id}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainsDataSourceName(rand, map[string]string{
			"transit_router_multicast_domain_id": `"${alicloud_cen_transit_router_multicast_domain.default.id}"`,
			"status":                             `"Active"`,
		}),
		FakeConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainsDataSourceName(rand, map[string]string{
			"transit_router_multicast_domain_id": `"${alicloud_cen_transit_router_multicast_domain.default.id}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainsDataSourceName(rand, map[string]string{
			"ids":                                `["${alicloud_cen_transit_router_multicast_domain.default.id}"]`,
			"name_regex":                         `"${alicloud_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_name}"`,
			"transit_router_multicast_domain_id": `"${alicloud_cen_transit_router_multicast_domain.default.id}"`,
			"status":                             `"Active"`,
		}),
		FakeConfig: testAccCheckAlicloudCenTransitRouterMulticastDomainsDataSourceName(rand, map[string]string{
			"ids":                                `["${alicloud_cen_transit_router_multicast_domain.default.id}_fake"]`,
			"name_regex":                         `"${alicloud_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_name}_fake"`,
			"transit_router_multicast_domain_id": `"${alicloud_cen_transit_router_multicast_domain.default.id}_fake"`,
		}),
	}
	var existAlicloudCenTransitRouterMulticastDomainsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"domains.#":                   "1",
			"domains.0.id":                acctest.CHECKSET,
			"domains.0.transit_router_id": acctest.CHECKSET,
			"domains.0.transit_router_multicast_domain_id":          acctest.CHECKSET,
			"domains.0.transit_router_multicast_domain_name":        acctest.CHECKSET,
			"domains.0.transit_router_multicast_domain_description": acctest.CHECKSET,
			"domains.0.status": "Active",
		}
	}
	var fakeAlicloudCenTransitRouterMulticastDomainsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"domains.#": "0",
		}
	}
	var alicloudCenTransitRouterMulticastDomainsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cen_transit_router_multicast_domains.default",
		ExistMapFunc: existAlicloudCenTransitRouterMulticastDomainsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCenTransitRouterMulticastDomainsDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudCenTransitRouterMulticastDomainsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, transitRouterMulticastDomainIdConf, statusConf, allConf)
}

func testAccCheckAlicloudCenTransitRouterMulticastDomainsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccCenTransitRouterMulticastDomain-%d"
	}

	resource "alicloud_cen_instance" "default" {
  		cen_instance_name = var.name
	}

	resource "alicloud_cen_transit_router" "default" {
  		cen_id            = alicloud_cen_instance.default.id
  		support_multicast = true
	}

	resource "alicloud_cen_transit_router_multicast_domain" "default" {
  		transit_router_id                           = alicloud_cen_transit_router.default.transit_router_id
  		transit_router_multicast_domain_name        = var.name
  		transit_router_multicast_domain_description = var.name
	}

	data "alicloud_cen_transit_router_multicast_domains" "default" {
  		transit_router_id = alicloud_cen_transit_router_multicast_domain.default.transit_router_id
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
