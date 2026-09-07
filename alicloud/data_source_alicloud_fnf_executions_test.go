package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudFnFExecutionsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.FnFSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFnFExecutionsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_fnf_execution.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudFnFExecutionsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_fnf_execution.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFnFExecutionsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_fnf_execution.default.execution_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudFnFExecutionsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_fnf_execution.default.execution_name}_fake"`,
		}),
	}
	statusRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFnFExecutionsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_fnf_execution.default.id}"]`,
			"status": `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudFnFExecutionsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_fnf_execution.default.id}_fake"]`,
			"status": `"TimedOut"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFnFExecutionsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_fnf_execution.default.id}"]`,
			"name_regex": `"${alicloud_fnf_execution.default.execution_name}"`,
			"status":     `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudFnFExecutionsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_fnf_execution.default.id}_fake"]`,
			"name_regex": `"${alicloud_fnf_execution.default.execution_name}_fake"`,
			"status":     `"TimedOut"`,
		}),
	}
	var existAlicloudFnFExecutionsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"executions.#":                "1",
			"executions.0.id":             CHECKSET,
			"executions.0.execution_name": fmt.Sprintf("tf-testAccExecution-%d", rand),
			"executions.0.flow_name":      fmt.Sprintf("tf-testAccExecution-%d", rand),
			"executions.0.input":          CHECKSET,
			"executions.0.output":         "",
			"executions.0.started_time":   CHECKSET,
			"executions.0.status":         "Running",
			"executions.0.stopped_time":   "",
		}
	}
	var fakeAlicloudFnFExecutionsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudFnFExecutionsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_fnf_executions.default",
		existMapFunc: existAlicloudFnFExecutionsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudFnFExecutionsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudFnFExecutionsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusRegexConf, allConf)
}
func testAccCheckAlicloudFnFExecutionsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccExecution-%d"
}

resource "alicloud_ram_role" "default" {
  name     = var.name
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
  definition  = <<EOF
  version: v1beta1
  type: flow
  steps:
    - type: wait
      name: custom_wait
      duration: $.wait
  EOF
  role_arn    = alicloud_ram_role.default.arn
  description = "Test for terraform fnf_flow."
  name        = var.name
  type        = "FDL"
}


resource "alicloud_fnf_execution" "default" {
  execution_name = var.name
  flow_name      = alicloud_fnf_flow.default.name
  input          = "{\"wait\": 600}"
}

data "alicloud_fnf_executions" "default" {	
	enable_details = true
    flow_name      = alicloud_fnf_flow.default.name
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
