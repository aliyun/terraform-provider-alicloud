package schedulerx_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudSchedulerxNamespacesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.SchedulerxSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSchedulerxNamespacesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_schedulerx_namespace.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudSchedulerxNamespacesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_schedulerx_namespace.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSchedulerxNamespacesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_schedulerx_namespace.default.namespace_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudSchedulerxNamespacesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_schedulerx_namespace.default.namespace_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSchedulerxNamespacesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_schedulerx_namespace.default.id}"]`,
			"name_regex": `"${alicloud_schedulerx_namespace.default.namespace_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudSchedulerxNamespacesDataSourceName(rand, map[string]string{
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
			"namespaces.0.namespace_id":   acctest.CHECKSET,
			"namespaces.0.id":             acctest.CHECKSET,
		}
	}
	var fakeAlicloudSchedulerxNamespacesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudSchedulerxNamespacesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_schedulerx_namespaces.default",
		ExistMapFunc: existAlicloudSchedulerxNamespacesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudSchedulerxNamespacesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudSchedulerxNamespacesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
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
