package edas_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEdasNamespacesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.EdasSupportedRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEdasNamespacesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_edas_namespace.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudEdasNamespacesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_edas_namespace.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEdasNamespacesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_edas_namespace.default.namespace_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudEdasNamespacesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_edas_namespace.default.namespace_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEdasNamespacesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_edas_namespace.default.id}"]`,
			"name_regex": `"${alicloud_edas_namespace.default.namespace_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudEdasNamespacesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_edas_namespace.default.id}_fake"]`,
			"name_regex": `"${alicloud_edas_namespace.default.namespace_name}_fake"`,
		}),
	}
	var existAlicloudEdasNamespacesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"namespaces.#":                      "1",
			"namespaces.0.debug_enable":         "false",
			"namespaces.0.description":          fmt.Sprintf("tf-testAccNamespace-%d", rand),
			"namespaces.0.namespace_logical_id": fmt.Sprintf("%s:tftest%d", acctest.DefaultRegionToTest, rand),
			"namespaces.0.namespace_name":       fmt.Sprintf("tf-testAccNamespace-%d", rand),
			"namespaces.0.user_id":              acctest.CHECKSET,
			"namespaces.0.belong_region":        acctest.CHECKSET,
		}
	}
	var fakeAlicloudEdasNamespacesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEdasNamespacesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_edas_namespaces.default",
		ExistMapFunc: existAlicloudEdasNamespacesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudEdasNamespacesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudEdasNamespacesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudEdasNamespacesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccNamespace-%d"
}

variable "logical_id" {
  default = "%s:tftest%d"
}

resource "alicloud_edas_namespace" "default" {
	debug_enable = false
	description = var.name
	namespace_logical_id = var.logical_id
	namespace_name = var.name
}

data "alicloud_edas_namespaces" "default" {
	%s
}
`, rand, acctest.DefaultRegionToTest, rand, strings.Join(pairs, " \n "))
	return config
}
