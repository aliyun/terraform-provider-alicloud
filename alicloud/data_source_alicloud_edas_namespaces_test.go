package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEdasNamespacesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.EdasSupportedRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEdasNamespacesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_edas_namespace.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEdasNamespacesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_edas_namespace.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEdasNamespacesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_edas_namespace.default.namespace_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEdasNamespacesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_edas_namespace.default.namespace_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEdasNamespacesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_edas_namespace.default.id}"]`,
			"name_regex": `"${alicloud_edas_namespace.default.namespace_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEdasNamespacesDataSourceName(rand, map[string]string{
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
			"namespaces.0.namespace_logical_id": fmt.Sprintf("%s:tftest%d", defaultRegionToTest, rand),
			"namespaces.0.namespace_name":       fmt.Sprintf("tf-testAccNamespace-%d", rand),
			"namespaces.0.user_id":              CHECKSET,
			"namespaces.0.belong_region":        CHECKSET,
		}
	}
	var fakeAlicloudEdasNamespacesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEdasNamespacesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_edas_namespaces.default",
		existMapFunc: existAlicloudEdasNamespacesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEdasNamespacesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEdasNamespacesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
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
`, rand, defaultRegionToTest, rand, strings.Join(pairs, " \n "))
	return config
}
