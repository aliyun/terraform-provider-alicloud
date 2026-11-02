package privatezone_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudPvtzResolverZonesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	regionIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPvtzResolverZonesDataSourceName(rand, map[string]string{}),
		FakeConfig:  "",
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPvtzResolverZonesDataSourceName(rand, map[string]string{
			"status": `"NORMAL"`,
		}),
		FakeConfig: "",
	}

	var existAlicloudPvtzResolverZoneDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": acctest.CHECKSET,
		}
	}
	var fakePvtzResolverZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": "0",
		}
	}
	var alicloudPvtzResolverZonesAccountBusesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_pvtz_resolver_zones.default",
		ExistMapFunc: existAlicloudPvtzResolverZoneDataSourceNameMapFunc,
		FakeMapFunc:  fakePvtzResolverZonesMapFunc,
	}

	alicloudPvtzResolverZonesAccountBusesCheckInfo.DataSourceTestCheck(t, rand, regionIdConf, statusConf)
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
