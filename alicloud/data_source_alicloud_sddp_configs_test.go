package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSddpConfigDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 100)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpConfigDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sddp_config.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSddpConfigDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sddp_config.default.id}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpConfigDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_sddp_config.default.id}"]`,
			"lang": `"zh"`,
		}),
		fakeConfig: testAccCheckAlicloudSddpConfigDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sddp_config.default.id}_fake"]`,
		}),
	}

	var existAlicloudSddpConfigDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "1",
			"configs.#":             "1",
			"configs.0.code":        "access_failed_cnt",
			"configs.0.description": CHECKSET,
			"configs.0.value":       "50",
		}
	}
	var fakeAlicloudSddpConfigDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"configs.#": "0",
		}
	}
	var alicloudSaeNamespaceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_sddp_configs.default",
		existMapFunc: existAlicloudSddpConfigDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudSddpConfigDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.SddpSupportRegions)
	}
	alicloudSaeNamespaceCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, allConf)

}
func testAccCheckAlicloudSddpConfigDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
resource "alicloud_sddp_config" "default" {
  code = "access_failed_cnt"
  description = "tf-testacc"
  value = 50
}

data "alicloud_sddp_configs" "default" {	
	%s
}
`, strings.Join(pairs, " \n "))
	return config
}
