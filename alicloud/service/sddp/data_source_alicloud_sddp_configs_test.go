package sddp_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudSddpConfigDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 100)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSddpConfigDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sddp_config.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudSddpConfigDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sddp_config.default.id}_fake"]`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSddpConfigDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_sddp_config.default.id}"]`,
			"lang": `"zh"`,
		}),
		FakeConfig: testAccCheckAlicloudSddpConfigDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sddp_config.default.id}_fake"]`,
		}),
	}

	var existAlicloudSddpConfigDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "1",
			"configs.#":             "1",
			"configs.0.code":        "access_failed_cnt",
			"configs.0.description": acctest.CHECKSET,
			"configs.0.value":       "50",
		}
	}
	var fakeAlicloudSddpConfigDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"configs.#": "0",
		}
	}
	var alicloudSaeNamespaceCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_sddp_configs.default",
		ExistMapFunc: existAlicloudSddpConfigDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudSddpConfigDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-zhangjiakou"})
	}
	alicloudSaeNamespaceCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, allConf)

}
func testAccCheckAlicloudSddpConfigDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
resource "alicloud_sddp_config" "default" {
  code        = "access_failed_cnt"
  description = "非授权资源多次访问尝试：目前的阈值定义为50次"
  value       = 50
}

data "alicloud_sddp_configs" "default" {
	%s
}
`, strings.Join(pairs, " \n "))
}
