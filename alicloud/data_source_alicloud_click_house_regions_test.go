package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudClickHouseRegionsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	currentConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudClickHouseRegionsDataSourceName(rand, map[string]string{
			"current": "true",
		}),
		fakeConfig: "",
	}
	regionIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudClickHouseRegionsDataSourceName(rand, map[string]string{
			"region_id": `"cn-qingdao"`,
		}),
		fakeConfig: "",
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
	var alicloudClickHouseAccountBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_click_house_regions.default",
		existMapFunc: existAlicloudClickHouseRegionDataSourceNameMapFunc,
		fakeMapFunc:  fakeClickHouseRegionsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.ClickHouseSupportRegions)
	}
	alicloudClickHouseAccountBusesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, currentConf, regionIdConf)
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
