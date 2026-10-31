package sae_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudSAENamespaceDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 100)
	acctest.CheckoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSaeNamespaceDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sae_namespace.default.namespace_id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudSaeNamespaceDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sae_namespace.default.namespace_id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSaeNamespaceDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_sae_namespace.default.namespace_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudSaeNamespaceDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_sae_namespace.default.namespace_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSaeNamespaceDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_sae_namespace.default.namespace_id}"]`,
			"name_regex": `"${alicloud_sae_namespace.default.namespace_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudSaeNamespaceDataSourceName(rand, map[string]string{
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
	var alicloudSaeNamespaceCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_sae_namespaces.default",
		ExistMapFunc: existAlicloudSaeNamespaceDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudSaeNamespaceDataSourceNameMapFunc,
	}

	alicloudSaeNamespaceCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
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
`, rand, acctest.DefaultRegionToTest, rand, strings.Join(pairs, " \n "))
	return config
}
