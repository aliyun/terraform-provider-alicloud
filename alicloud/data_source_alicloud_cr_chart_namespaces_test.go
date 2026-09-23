package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCRNameSpaceDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCRNameSpaceDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_cr_chart_namespace.default.instance_id}"`,
			"ids":         `["${alicloud_cr_chart_namespace.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCRNameSpaceDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_cr_chart_namespace.default.instance_id}"`,
			"ids":         `["${alicloud_cr_chart_namespace.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCRNameSpaceDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_cr_chart_namespace.default.instance_id}"`,
			"name_regex":  `"${alicloud_cr_chart_namespace.default.namespace_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCRNameSpaceDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_cr_chart_namespace.default.instance_id}"`,
			"name_regex":  `"${alicloud_cr_chart_namespace.default.namespace_name}_fake"`,
		}),
	}
	var existAlicloudCrNameSpaceDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"namespaces.#":                   CHECKSET,
			"namespaces.0.auto_create_repo":  "false",
			"namespaces.0.default_repo_type": "PRIVATE",
			"namespaces.0.namespace_name":    CHECKSET,
		}
	}
	var fakeAlicloudCrNameSpaceDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"names.#":      "0",
			"namespaces.#": "0",
		}
	}
	var alicloudCrNameSpaceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cr_chart_namespaces.default",
		existMapFunc: existAlicloudCrNameSpaceDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCrNameSpaceDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudCrNameSpaceCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsRegexConf, nameRegexConf)
}
func testAccCheckAlicloudCRNameSpaceDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {	
	default = "tftest%d"
}

data "alicloud_cr_ee_instances" "default" {}

resource "alicloud_cr_chart_namespace" "default" {
	instance_id        = data.alicloud_cr_ee_instances.default.ids.0
	namespace_name       = var.name
}

data "alicloud_cr_chart_namespaces" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
