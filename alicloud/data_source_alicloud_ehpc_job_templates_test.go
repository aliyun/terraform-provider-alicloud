package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEhpcJobTemplatesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEhpcJobTemplatesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ehpc_job_template.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEhpcJobTemplatesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ehpc_job_template.default.id}_fake"]`,
		}),
	}

	var existAlicloudEhpcJobTemplatesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"templates.#":                   "1",
			"templates.0.command_line":      "./LammpsTest/lammps.pbs",
			"templates.0.job_template_name": fmt.Sprintf("tf-testAccTemplates-%d", rand),
		}
	}
	var fakeAlicloudEhpcJobTemplatesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudEhpcJobTemplatesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ehpc_job_templates.default",
		existMapFunc: existAlicloudEhpcJobTemplatesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEhpcJobTemplatesDataSourceNameMapFunc,
	}
	alicloudEhpcJobTemplatesCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}
func testAccCheckAlicloudEhpcJobTemplatesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccTemplates-%d"
}

resource "alicloud_ehpc_job_template" "default"{
  job_template_name =  var.name
  command_line=       "./LammpsTest/lammps.pbs"
}

data "alicloud_ehpc_job_templates" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
