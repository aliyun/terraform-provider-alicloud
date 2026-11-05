package cloudmonitorservice_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCmsNamespacesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsNamespacesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_namespace.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCmsNamespacesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cms_namespace.default.id}_fake"]`,
		}),
	}
	keywordConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCmsNamespacesDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_cms_namespace.default.id}"]`,
			"keyword": `"${alicloud_cms_namespace.default.namespace}"`,
		}),
		FakeConfig: testAccCheckAlicloudCmsNamespacesDataSourceName(rand, map[string]string{
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
			"namespaces.0.namespace_id":  acctest.CHECKSET,
			"namespaces.0.id":            acctest.CHECKSET,
			"namespaces.0.create_time":   acctest.CHECKSET,
			"namespaces.0.modify_time":   acctest.CHECKSET,
		}
	}
	var fakeAlicloudCmsNamespacesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudCmsNamespacesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cms_namespaces.default",
		ExistMapFunc: existAlicloudCmsNamespacesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCmsNamespacesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudCmsNamespacesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, keywordConf)
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
