package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCenTransitRouterRouterAvailableResourcesDataSource(t *testing.T) {
	rand := acctest.RandInt()

	noConf := dataSourceTestAccConfig{
		fakeConfig: testAccCheckAliCloudCenTransitRouterAvailableResourcesDataSourceName(rand, map[string]string{}),
	}

	supportMulticastConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCenTransitRouterAvailableResourcesDataSourceName(rand, map[string]string{
			"support_multicast": `"true"`,
		}),
		fakeConfig: testAccCheckAliCloudCenTransitRouterAvailableResourcesDataSourceName(rand, map[string]string{
			"support_multicast": `"false"`,
		}),
	}

	var existAliCloudCenTransitRouterAvailableResourcesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"resources.#":                   "1",
			"resources.0.support_multicast": "true",
			"resources.0.master_zones.#":    CHECKSET,
			"resources.0.slave_zones.#":     CHECKSET,
			"resources.0.available_zones.#": CHECKSET,
		}
	}

	var fakeAliCloudCenTransitRouterAvailableResourcesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"resources.#":                   "1",
			"resources.0.support_multicast": "false",
		}
	}

	var alicloudCenTransitRouterRouteEntriesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_transit_router_available_resources.default",
		existMapFunc: existAliCloudCenTransitRouterAvailableResourcesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudCenTransitRouterAvailableResourcesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	alicloudCenTransitRouterRouteEntriesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, noConf, supportMulticastConf)
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
