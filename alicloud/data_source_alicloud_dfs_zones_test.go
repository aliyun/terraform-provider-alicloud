package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDfsZonesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	checkoutSupportedRegions(t, true, connectivity.DfsSupportRegions)
	regionIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDfsZonesDataSourceName(rand, map[string]string{}),
		fakeConfig:  "",
	}

	var existAlicloudDfsZoneDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": CHECKSET,
		}
	}
	var fakeDfsZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": "0",
		}
	}
	var alicloudDfsZonesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_dfs_zones.default",
		existMapFunc: existAlicloudDfsZoneDataSourceNameMapFunc,
		fakeMapFunc:  fakeDfsZonesMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudDfsZonesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, regionIdConf)
}
func testAccCheckAlicloudDfsZonesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_dfs_zones" "default" {	
	%s
}
`, strings.Join(pairs, " \n "))
	return config
}
