package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCRChartRepositoriesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCRRepoDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_cr_chart_namespace.default.instance_id}"`,
			"ids":         `["${alicloud_cr_chart_repository.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCRRepoDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_cr_chart_namespace.default.instance_id}"`,
			"ids":         `["${alicloud_cr_chart_repository.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCRRepoDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_cr_chart_namespace.default.instance_id}"`,
			"name_regex":  `"${alicloud_cr_chart_repository.default.repo_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCRRepoDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_cr_chart_namespace.default.instance_id}"`,
			"name_regex":  `"${alicloud_cr_chart_repository.default.repo_name}_fake"`,
		}),
	}
	var existAlicloudCrRepoDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"names.#":                  "1",
			"repositories.#":           CHECKSET,
			"repositories.0.repo_type": "PUBLIC",
			"repositories.0.repo_name": CHECKSET,
			"repositories.0.summary":   CHECKSET,
		}
	}
	var fakeAlicloudCrRepoDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"names.#":        "0",
			"repositories.#": "0",
		}
	}
	var alicloudCrRepoCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cr_chart_repositories.default",
		existMapFunc: existAlicloudCrRepoDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCrRepoDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudCrRepoCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsRegexConf, nameRegexConf)
}
func testAccCheckAlicloudCRRepoDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {	
	default = "tf_testacc_cr_repo%d"
}

data "alicloud_cr_ee_instances" "default" {}

resource "alicloud_cr_chart_namespace" "default" {
	instance_id        = data.alicloud_cr_ee_instances.default.ids.0
	namespace_name       = var.name
}

resource "alicloud_cr_chart_repository" "default" {
	instance_id        		  = alicloud_cr_chart_namespace.default.instance_id
	repo_namespace_name       = alicloud_cr_chart_namespace.default.namespace_name
	repo_name				  = var.name
	repo_type				  = "PUBLIC"
	summary					  = var.name
}


data "alicloud_cr_chart_repositories" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
