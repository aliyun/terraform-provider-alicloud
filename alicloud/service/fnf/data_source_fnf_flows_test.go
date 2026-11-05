package fnf_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const fnfFlowsDefinition = `version: v1beta1\ntype: flow\nsteps:\n  - type: pass\n    name: helloworld`

func TestAccAlicloudFnfFlowsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.FnfSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fnfFlowsTestAccConfig(rand, map[string]string{
					"ids": `["${alicloud_fnf_flow.default.name}"]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_fnf_flows.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_fnf_flows.default", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_fnf_flows.default", "flows.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_fnf_flows.default", "flows.0.definition", strings.Replace(fnfFlowsDefinition, `\n`, "\n", -1)),
					resource.TestCheckResourceAttr("data.alicloud_fnf_flows.default", "flows.0.description", "tf-testaccFnFFlow983041"),
					resource.TestCheckResourceAttrSet("data.alicloud_fnf_flows.default", "flows.0.name"),
					resource.TestCheckResourceAttrSet("data.alicloud_fnf_flows.default", "flows.0.role_arn"),
					resource.TestCheckResourceAttr("data.alicloud_fnf_flows.default", "flows.0.type", "FDL"),
				),
			},
			{
				Config: fnfFlowsTestAccConfig(rand, map[string]string{
					"ids": `["${alicloud_fnf_flow.default.name}_fake"]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_fnf_flows.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_fnf_flows.default", "names.#", "0"),
				),
			},
			{
				Config: fnfFlowsTestAccConfig(rand, map[string]string{
					"name_regex": `"${alicloud_fnf_flow.default.name}"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_fnf_flows.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_fnf_flows.default", "names.#", "1"),
				),
			},
			{
				Config: fnfFlowsTestAccConfig(rand, map[string]string{
					"name_regex": `"${alicloud_fnf_flow.default.name}_fake"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_fnf_flows.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_fnf_flows.default", "names.#", "0"),
				),
			},
			{
				Config: fnfFlowsTestAccConfig(rand, map[string]string{
					"ids":        `["${alicloud_fnf_flow.default.name}"]`,
					"name_regex": `"${alicloud_fnf_flow.default.name}"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_fnf_flows.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_fnf_flows.default", "names.#", "1"),
				),
			},
			{
				Config: fnfFlowsTestAccConfig(rand, map[string]string{
					"ids":        `["${alicloud_fnf_flow.default.name}_fake"]`,
					"name_regex": `"${alicloud_fnf_flow.default.name}_fake"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_fnf_flows.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_fnf_flows.default", "names.#", "0"),
				),
			},
		},
	})
}

func fnfFlowsTestAccConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`

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
}
