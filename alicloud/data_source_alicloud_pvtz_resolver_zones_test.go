package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudPvtzResolverZonesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	regionIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPvtzResolverZonesDataSourceName(rand, map[string]string{}),
		fakeConfig:  "",
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPvtzResolverZonesDataSourceName(rand, map[string]string{
			"status": `"NORMAL"`,
		}),
		fakeConfig: "",
	}

	var existAlicloudPvtzResolverZoneDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": CHECKSET,
		}
	}
	var fakePvtzResolverZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": "0",
		}
	}
	var alicloudPvtzResolverZonesAccountBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_pvtz_resolver_zones.default",
		existMapFunc: existAlicloudPvtzResolverZoneDataSourceNameMapFunc,
		fakeMapFunc:  fakePvtzResolverZonesMapFunc,
	}

	alicloudPvtzResolverZonesAccountBusesCheckInfo.dataSourceTestCheck(t, rand, regionIdConf, statusConf)
}

func testAccCheckAlicloudPvtzResolverZonesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_pvtz_resolver_zones" "default" {  
   %s
}
`, strings.Join(pairs, " \n "))
	return config
}
