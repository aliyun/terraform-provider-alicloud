package cr_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCRChartRepositoriesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCRRepoDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_cr_chart_namespace.default.instance_id}"`,
			"ids":         `["${alicloud_cr_chart_repository.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCRRepoDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_cr_chart_namespace.default.instance_id}"`,
			"ids":         `["${alicloud_cr_chart_repository.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCRRepoDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_cr_chart_namespace.default.instance_id}"`,
			"name_regex":  `"${alicloud_cr_chart_repository.default.repo_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCRRepoDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_cr_chart_namespace.default.instance_id}"`,
			"name_regex":  `"${alicloud_cr_chart_repository.default.repo_name}_fake"`,
		}),
	}
	var existAlicloudCrRepoDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"names.#":                  "1",
			"repositories.#":           acctest.CHECKSET,
			"repositories.0.repo_type": "PUBLIC",
			"repositories.0.repo_name": acctest.CHECKSET,
			"repositories.0.summary":   acctest.CHECKSET,
		}
	}
	var fakeAlicloudCrRepoDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"names.#":        "0",
			"repositories.#": "0",
		}
	}
	var alicloudCrRepoCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cr_chart_repositories.default",
		ExistMapFunc: existAlicloudCrRepoDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCrRepoDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudCrRepoCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsRegexConf, nameRegexConf)
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
