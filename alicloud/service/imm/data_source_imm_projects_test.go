package imm_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudIMMProjectsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudIMMProjectDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_imm_project.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudIMMProjectDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_imm_project.default.id}_fake"]`,
		}),
	}

	var existDataAlicloudIMMProjectsSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":              "1",
			"projects.#":         "1",
			"projects.0.project": fmt.Sprintf("tf-testAccIMMProject%d", rand),
		}
	}
	var fakeDataAlicloudIMMProjectsSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"projects.#": "0",
		}
	}
	var alicloudIMMProjectCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_imm_projects.default",
		ExistMapFunc: existDataAlicloudIMMProjectsSourceNameMapFunc,
		FakeMapFunc:  fakeDataAlicloudIMMProjectsSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.IMMSupportRegions)
	}
	alicloudIMMProjectCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}
func testAccCheckAlicloudIMMProjectDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccIMMProject%d"
}


resource "alicloud_imm_project" "default" {
	project = var.name
}

data "alicloud_imm_projects" "default" {
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
