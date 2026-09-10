package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSchedulerxNamespacesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.SchedulerxSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSchedulerxNamespacesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_schedulerx_namespace.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSchedulerxNamespacesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_schedulerx_namespace.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSchedulerxNamespacesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_schedulerx_namespace.default.namespace_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSchedulerxNamespacesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_schedulerx_namespace.default.namespace_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSchedulerxNamespacesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_schedulerx_namespace.default.id}"]`,
			"name_regex": `"${alicloud_schedulerx_namespace.default.namespace_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSchedulerxNamespacesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_schedulerx_namespace.default.id}_fake"]`,
			"name_regex": `"${alicloud_schedulerx_namespace.default.namespace_name}_fake"`,
		}),
	}
	var existAlicloudSchedulerxNamespacesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"namespaces.#":                "1",
			"namespaces.0.description":    fmt.Sprintf("tf-testAccNamespace-%d", rand),
			"namespaces.0.namespace_name": fmt.Sprintf("tf-testAccNamespace-%d", rand),
			"namespaces.0.namespace_id":   CHECKSET,
			"namespaces.0.id":             CHECKSET,
		}
	}
	var fakeAlicloudSchedulerxNamespacesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudSchedulerxNamespacesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_schedulerx_namespaces.default",
		existMapFunc: existAlicloudSchedulerxNamespacesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudSchedulerxNamespacesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudSchedulerxNamespacesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudSchedulerxNamespacesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccNamespace-%d"
}

resource "alicloud_schedulerx_namespace" "default" {
	description = var.name
	namespace_name = var.name
}

data "alicloud_schedulerx_namespaces" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
