package ehpc_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudEhpcJobTemplatesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEhpcJobTemplatesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ehpc_job_template.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudEhpcJobTemplatesDataSourceName(rand, map[string]string{
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
	var alicloudEhpcJobTemplatesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ehpc_job_templates.default",
		ExistMapFunc: existAlicloudEhpcJobTemplatesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudEhpcJobTemplatesDataSourceNameMapFunc,
	}
	alicloudEhpcJobTemplatesCheckInfo.DataSourceTestCheck(t, rand, idsConf)
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
