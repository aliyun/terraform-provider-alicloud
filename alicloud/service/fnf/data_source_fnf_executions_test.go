package fnf_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlicloudFnFExecutionsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	executionName := fmt.Sprintf("tf-testAccExecution-%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fnfExecutionsTestAccConfig(rand, map[string]string{
					"ids": `["${alicloud_fnf_execution.default.id}"]`,
				}),
				Check: fnfExecutionsExistChecks(executionName),
			},
			{
				Config: fnfExecutionsTestAccConfig(rand, map[string]string{
					"ids": `["${alicloud_fnf_execution.default.id}_fake"]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_fnf_executions.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_fnf_executions.default", "names.#", "0"),
				),
			},
			{
				Config: fnfExecutionsTestAccConfig(rand, map[string]string{
					"name_regex": `"${alicloud_fnf_execution.default.execution_name}"`,
				}),
				Check: fnfExecutionsExistChecks(executionName),
			},
			{
				Config: fnfExecutionsTestAccConfig(rand, map[string]string{
					"name_regex": `"${alicloud_fnf_execution.default.execution_name}_fake"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_fnf_executions.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_fnf_executions.default", "names.#", "0"),
				),
			},
			{
				Config: fnfExecutionsTestAccConfig(rand, map[string]string{
					"ids":    `["${alicloud_fnf_execution.default.id}"]`,
					"status": `"Running"`,
				}),
				Check: fnfExecutionsExistChecks(executionName),
			},
			{
				Config: fnfExecutionsTestAccConfig(rand, map[string]string{
					"ids":    `["${alicloud_fnf_execution.default.id}_fake"]`,
					"status": `"TimedOut"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_fnf_executions.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_fnf_executions.default", "names.#", "0"),
				),
			},
			{
				Config: fnfExecutionsTestAccConfig(rand, map[string]string{
					"ids":        `["${alicloud_fnf_execution.default.id}"]`,
					"name_regex": `"${alicloud_fnf_execution.default.execution_name}"`,
					"status":     `"Running"`,
				}),
				Check: fnfExecutionsExistChecks(executionName),
			},
			{
				Config: fnfExecutionsTestAccConfig(rand, map[string]string{
					"ids":        `["${alicloud_fnf_execution.default.id}_fake"]`,
					"name_regex": `"${alicloud_fnf_execution.default.execution_name}_fake"`,
					"status":     `"TimedOut"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_fnf_executions.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_fnf_executions.default", "names.#", "0"),
				),
			},
		},
	})
}

func fnfExecutionsExistChecks(executionName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("data.alicloud_fnf_executions.default", "ids.#", "1"),
		resource.TestCheckResourceAttr("data.alicloud_fnf_executions.default", "names.#", "1"),
		resource.TestCheckResourceAttr("data.alicloud_fnf_executions.default", "executions.#", "1"),
		resource.TestCheckResourceAttrSet("data.alicloud_fnf_executions.default", "executions.0.id"),
		resource.TestCheckResourceAttr("data.alicloud_fnf_executions.default", "executions.0.execution_name", executionName),
		resource.TestCheckResourceAttr("data.alicloud_fnf_executions.default", "executions.0.flow_name", executionName),
		resource.TestCheckResourceAttrSet("data.alicloud_fnf_executions.default", "executions.0.input"),
		resource.TestCheckResourceAttr("data.alicloud_fnf_executions.default", "executions.0.output", ""),
		resource.TestCheckResourceAttrSet("data.alicloud_fnf_executions.default", "executions.0.started_time"),
		resource.TestCheckResourceAttr("data.alicloud_fnf_executions.default", "executions.0.status", "Running"),
		resource.TestCheckResourceAttr("data.alicloud_fnf_executions.default", "executions.0.stopped_time", ""),
	)
}

func fnfExecutionsTestAccConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`

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
}
