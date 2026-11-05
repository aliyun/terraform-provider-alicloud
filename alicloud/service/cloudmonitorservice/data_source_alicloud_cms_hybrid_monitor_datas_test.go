package cloudmonitorservice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCmsHybridMonitorDatasDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsHybridMonitorDatasDataSourceName(rand, map[string]string{
			"period": `"60"`,
		}),
		FakeConfig: "",
	}
	var existAlicloudCmsHybridMonitorDatasDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"datas.#": acctest.CHECKSET,
		}
	}
	var fakeAlicloudCmsHybridMonitorDatasDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"datas.#": "0",
		}
	}
	var alicloudCmsHybridMonitorDatasCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cms_hybrid_monitor_datas.default",
		ExistMapFunc: existAlicloudCmsHybridMonitorDatasDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCmsHybridMonitorDatasDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudCmsHybridMonitorDatasCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
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
