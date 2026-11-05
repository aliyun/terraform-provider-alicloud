package cdn_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCdnBlockedRegionsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	languageConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCdnBlockedRegionsDataSourceName(rand, map[string]string{
			"language": `"zh"`,
		}),
		fakeConfig: "",
	}

	var existAlicloudClickHouseRegionDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"regions.#": acctest.CHECKSET,
		}
	}
	var fakeClickHouseRegionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"regions.#": "0",
		}
	}
	var alicloudCdnBlockedRegionsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cdn_blocked_regions.default",
		existMapFunc: existAlicloudClickHouseRegionDataSourceNameMapFunc,
		fakeMapFunc:  fakeClickHouseRegionsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudCdnBlockedRegionsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, languageConf)
}
func testAccCheckAlicloudCdnBlockedRegionsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_cdn_blocked_regions" "default" {	
	%s
}
`, strings.Join(pairs, " \n "))
	return config
}
