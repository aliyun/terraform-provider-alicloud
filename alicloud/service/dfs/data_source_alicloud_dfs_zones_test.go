package dfs_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudDfsZonesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.DfsSupportRegions)
	regionIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDfsZonesDataSourceName(rand, map[string]string{}),
		FakeConfig:  "",
	}

	var existAlicloudDfsZoneDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": acctest.CHECKSET,
		}
	}
	var fakeDfsZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": "0",
		}
	}
	var alicloudDfsZonesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_dfs_zones.default",
		ExistMapFunc: existAlicloudDfsZoneDataSourceNameMapFunc,
		FakeMapFunc:  fakeDfsZonesMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudDfsZonesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, regionIdConf)
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
