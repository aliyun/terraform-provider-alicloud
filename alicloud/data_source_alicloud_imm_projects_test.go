package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudIMMProjectsDataSource(t *testing.T) {
	rand := acctest.RandInt()

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudIMMProjectDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_imm_project.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudIMMProjectDataSourceName(rand, map[string]string{
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
	var alicloudIMMProjectCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_imm_projects.default",
		existMapFunc: existDataAlicloudIMMProjectsSourceNameMapFunc,
		fakeMapFunc:  fakeDataAlicloudIMMProjectsSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.IMMSupportRegions)
	}
	alicloudIMMProjectCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
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
