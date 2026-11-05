package cen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCenTransitRouterRouterAvailableResourcesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()

	noConf := acctest.DataSourceTestAccConfig{
		FakeConfig: testAccCheckAliCloudCenTransitRouterAvailableResourcesDataSourceName(rand, map[string]string{}),
	}

	supportMulticastConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudCenTransitRouterAvailableResourcesDataSourceName(rand, map[string]string{
			"support_multicast": `"true"`,
		}),
		FakeConfig: testAccCheckAliCloudCenTransitRouterAvailableResourcesDataSourceName(rand, map[string]string{
			"support_multicast": `"false"`,
		}),
	}

	var existAliCloudCenTransitRouterAvailableResourcesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"resources.#":                   "1",
			"resources.0.support_multicast": "true",
			"resources.0.master_zones.#":    acctest.CHECKSET,
			"resources.0.slave_zones.#":     acctest.CHECKSET,
			"resources.0.available_zones.#": acctest.CHECKSET,
		}
	}

	var fakeAliCloudCenTransitRouterAvailableResourcesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"resources.#":                   "1",
			"resources.0.support_multicast": "false",
		}
	}

	var alicloudCenTransitRouterRouteEntriesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cen_transit_router_available_resources.default",
		ExistMapFunc: existAliCloudCenTransitRouterAvailableResourcesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudCenTransitRouterAvailableResourcesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	alicloudCenTransitRouterRouteEntriesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, noConf, supportMulticastConf)
}

func testAccCheckAliCloudCenTransitRouterAvailableResourcesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	data "alicloud_cen_transit_router_available_resources" "default" {
		%s
	}
`, strings.Join(pairs, " \n "))
	return config
}
