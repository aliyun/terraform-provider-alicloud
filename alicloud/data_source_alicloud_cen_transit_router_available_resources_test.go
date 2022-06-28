package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTransitRouterRouterAvailableResourcesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 2999)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTransitRouterAvailableResourcesDataSourceName(rand, map[string]string{}),
	}
	var existAlicloudCenTransitRouterAvailableResourcesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"resources.#": CHECKSET,
		}
	}
	var fakeAlicloudCenTransitRouterAvailableResourcesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"resources.#": "0",
		}
	}
	var alicloudCenTransitRouterRouteEntriesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_transit_router_available_resources.default",
		existMapFunc: existAlicloudCenTransitRouterAvailableResourcesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCenTransitRouterAvailableResourcesDataSourceNameMapFunc,
	}
	alicloudCenTransitRouterRouteEntriesCheckInfo.dataSourceTestCheck(t, rand, allConf)
}
func testAccCheckAlicloudCenTransitRouterAvailableResourcesDataSourceName(rand int, attrMap map[string]string) string {
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
