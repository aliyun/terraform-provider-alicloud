package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSAENamespaceDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 100)
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSaeNamespaceDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sae_namespace.default.namespace_id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSaeNamespaceDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sae_namespace.default.namespace_id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSaeNamespaceDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_sae_namespace.default.namespace_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSaeNamespaceDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_sae_namespace.default.namespace_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSaeNamespaceDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_sae_namespace.default.namespace_id}"]`,
			"name_regex": `"${alicloud_sae_namespace.default.namespace_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSaeNamespaceDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_sae_namespace.default.namespace_id}_fake"]`,
			"name_regex": `"${alicloud_sae_namespace.default.namespace_name}_fake"`,
		}),
	}
	var existAlicloudSaeNamespaceDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"names.#":                            "1",
			"namespaces.#":                       "1",
			"namespaces.0.namespace_description": fmt.Sprintf("tf-testAccsaenamespace-%d", rand),
			"namespaces.0.namespace_name":        fmt.Sprintf("tf-testAccsaenamespace-%d", rand),
		}
	}
	var fakeAlicloudSaeNamespaceDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudSaeNamespaceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_sae_namespaces.default",
		existMapFunc: existAlicloudSaeNamespaceDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudSaeNamespaceDataSourceNameMapFunc,
	}

	alicloudSaeNamespaceCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudSaeNamespaceDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccsaenamespace-%d"
}
variable "namespace_id" {	
	default = "%s:tftestacc%d"
}

resource "alicloud_sae_namespace" "default" {
	namespace_id = var.namespace_id
	namespace_name = var.name
	namespace_description = var.name
}

data "alicloud_sae_namespaces" "default" {	
	%s
}
`, rand, defaultRegionToTest, rand, strings.Join(pairs, " \n "))
	return config
}
