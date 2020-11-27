package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudFnfFlowsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFnfFlowsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_fnf_flow.default.name}"]`,
		}),
		fakeConfig: testAccCheckAlicloudFnfFlowsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_fnf_flow.default.name}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFnfFlowsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_fnf_flow.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudFnfFlowsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_fnf_flow.default.name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFnfFlowsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_fnf_flow.default.name}"]`,
			"name_regex": `"${alicloud_fnf_flow.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudFnfFlowsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_fnf_flow.default.name}_fake"]`,
			"name_regex": `"${alicloud_fnf_flow.default.name}_fake"`,
		}),
	}
	var existAlicloudFnfFlowsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"names.#":             "1",
			"flows.#":             "1",
			"flows.0.definition":  strings.Replace(`version: v1beta1\ntype: flow\nsteps:\n  - type: pass\n    name: helloworld`, `\n`, "\n", -1),
			"flows.0.description": `tf-testaccFnFFlow983041`,
			"flows.0.name":        CHECKSET,
			"flows.0.role_arn":    CHECKSET,
			"flows.0.type":        `FDL`,
		}
	}
	var fakeAlicloudFnfFlowsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudFnfFlowsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_fnf_flows.default",
		existMapFunc: existAlicloudFnfFlowsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudFnfFlowsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.FnfSupportRegions)
	}
	alicloudFnfFlowsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudFnfFlowsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccFlow-%d"
}
resource "alicloud_ram_role" "default" {
  name = var.name
  document = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "fnf.aliyuncs.com"
          ]
        }
      }
    ],
    "Version": "1"
  }
  EOF
}

resource "alicloud_fnf_flow" "default" {
definition = "version: v1beta1\ntype: flow\nsteps:\n  - type: pass\n    name: helloworld"
description = "tf-testaccFnFFlow983041"
name = var.name
role_arn = "${alicloud_ram_role.default.arn}"
type = "FDL"
}

data "alicloud_fnf_flows" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
