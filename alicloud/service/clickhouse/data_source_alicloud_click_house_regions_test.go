package clickhouse_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudClickHouseRegionsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	currentConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudClickHouseRegionsDataSourceName(rand, map[string]string{
			"current": "true",
		}),
		FakeConfig: "",
	}
	regionIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudClickHouseRegionsDataSourceName(rand, map[string]string{
			"region_id": `"cn-qingdao"`,
		}),
		FakeConfig: "",
	}

	var existAlicloudClickHouseRegionDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"regions.#": "1",
		}
	}
	var fakeClickHouseRegionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"regions.#": "0",
		}
	}
	var alicloudClickHouseAccountBusesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_click_house_regions.default",
		ExistMapFunc: existAlicloudClickHouseRegionDataSourceNameMapFunc,
		FakeMapFunc:  fakeClickHouseRegionsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.ClickHouseSupportRegions)
	}
	alicloudClickHouseAccountBusesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, currentConf, regionIdConf)
}
func testAccCheckAlicloudClickHouseRegionsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_click_house_regions" "default" {	
	%s
}
`, strings.Join(pairs, " \n "))
	return config
}
