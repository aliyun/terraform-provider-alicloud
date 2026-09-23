package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCmsHybridMonitorDatasDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsHybridMonitorDatasDataSourceName(rand, map[string]string{
			"period": `"60"`,
		}),
		fakeConfig: "",
	}
	var existAlicloudCmsHybridMonitorDatasDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"datas.#": CHECKSET,
		}
	}
	var fakeAlicloudCmsHybridMonitorDatasDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"datas.#": "0",
		}
	}
	var alicloudCmsHybridMonitorDatasCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_hybrid_monitor_datas.default",
		existMapFunc: existAlicloudCmsHybridMonitorDatasDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCmsHybridMonitorDatasDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudCmsHybridMonitorDatasCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
}
func testAccCheckAlicloudCmsHybridMonitorDatasDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testaccnamespace-%d"
}
resource "alicloud_cms_namespace" "default" {
	description = var.name
	namespace = var.name
	specification = "cms.s1.large"
}

data "alicloud_cms_hybrid_monitor_datas" "default" {
	namespace = alicloud_cms_namespace.default.id
	prom_sql = "AliyunEcs_cpu_total"
	start = "1657505665"
	end = "1657520065"
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
