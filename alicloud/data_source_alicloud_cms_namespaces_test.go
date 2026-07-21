package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCmsNamespacesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsNamespacesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_namespace.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCmsNamespacesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_namespace.default.id}_fake"]`,
		}),
	}
	keywordConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCmsNamespacesDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_cms_namespace.default.id}"]`,
			"keyword": `"${alicloud_cms_namespace.default.namespace}"`,
		}),
		fakeConfig: testAccCheckAlicloudCmsNamespacesDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_cms_namespace.default.id}"]`,
			"keyword": `"${alicloud_cms_namespace.default.namespace}_fake"`,
		}),
	}
	var existAlicloudCmsNamespacesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"namespaces.#":               "1",
			"namespaces.0.description":   fmt.Sprintf("tf-testaccnamespace-%d", rand),
			"namespaces.0.namespace":     fmt.Sprintf("tf-testaccnamespace-%d", rand),
			"namespaces.0.specification": "cms.s1.large",
			"namespaces.0.namespace_id":  CHECKSET,
			"namespaces.0.id":            CHECKSET,
			"namespaces.0.create_time":   CHECKSET,
			"namespaces.0.modify_time":   CHECKSET,
		}
	}
	var fakeAlicloudCmsNamespacesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudCmsNamespacesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cms_namespaces.default",
		existMapFunc: existAlicloudCmsNamespacesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCmsNamespacesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudCmsNamespacesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, keywordConf)
}
func testAccCheckAlicloudCmsNamespacesDataSourceName(rand int, attrMap map[string]string) string {
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

data "alicloud_cms_namespaces" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
