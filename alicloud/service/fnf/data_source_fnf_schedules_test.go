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

func TestAccAlicloudFnfSchedulesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.FnfSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fnfSchedulesTestAccConfig(rand, map[string]string{
					"ids": `["${alicloud_fnf_schedule.default.schedule_name}"]`,
				}),
				Check: fnfSchedulesExistChecks(),
			},
			{
				Config: fnfSchedulesTestAccConfig(rand, map[string]string{
					"ids": `["${alicloud_fnf_schedule.default.schedule_name}_fake"]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_fnf_schedules.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_fnf_schedules.default", "names.#", "0"),
				),
			},
			{
				Config: fnfSchedulesTestAccConfig(rand, map[string]string{
					"name_regex": `"${alicloud_fnf_schedule.default.schedule_name}"`,
				}),
				Check: fnfSchedulesExistChecks(),
			},
			{
				Config: fnfSchedulesTestAccConfig(rand, map[string]string{
					"name_regex": `"${alicloud_fnf_schedule.default.schedule_name}_fake"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_fnf_schedules.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_fnf_schedules.default", "names.#", "0"),
				),
			},
			{
				Config: fnfSchedulesTestAccConfig(rand, map[string]string{
					"ids":        `["${alicloud_fnf_schedule.default.schedule_name}"]`,
					"name_regex": `"${alicloud_fnf_schedule.default.schedule_name}"`,
				}),
				Check: fnfSchedulesExistChecks(),
			},
			{
				Config: fnfSchedulesTestAccConfig(rand, map[string]string{
					"ids":        `["${alicloud_fnf_schedule.default.schedule_name}_fake"]`,
					"name_regex": `"${alicloud_fnf_schedule.default.schedule_name}_fake"`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_fnf_schedules.default", "ids.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_fnf_schedules.default", "names.#", "0"),
				),
			},
		},
	})
}

func fnfSchedulesExistChecks() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("data.alicloud_fnf_schedules.default", "ids.#", "1"),
		resource.TestCheckResourceAttr("data.alicloud_fnf_schedules.default", "names.#", "1"),
		resource.TestCheckResourceAttr("data.alicloud_fnf_schedules.default", "schedules.#", "1"),
		resource.TestCheckResourceAttr("data.alicloud_fnf_schedules.default", "schedules.0.cron_expression", "30 9 * * * *"),
		resource.TestCheckResourceAttr("data.alicloud_fnf_schedules.default", "schedules.0.description", "tf-testaccFnFSchedule983041"),
		resource.TestCheckResourceAttr("data.alicloud_fnf_schedules.default", "schedules.0.enable", "true"),
		resource.TestCheckResourceAttr("data.alicloud_fnf_schedules.default", "schedules.0.payload", `{"tf-test": "test success"}`),
		resource.TestCheckResourceAttrSet("data.alicloud_fnf_schedules.default", "schedules.0.schedule_name"),
	)
}

func fnfSchedulesTestAccConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`

variable "name" {
	default = "tf-testAccSchedule-%d"
}

resource "alicloud_fnf_flow" "default" {
definition= "version: v1beta1\ntype: flow\nsteps:\n  - type: pass\n    name: helloworld"
description= "tf-testaccFnFFlow983041"
name = var.name
type= "FDL"
}

resource "alicloud_fnf_schedule" "default" {
cron_expression = "30 9 * * * *"
description = "tf-testaccFnFSchedule983041"
enable = "true"
flow_name = "${alicloud_fnf_flow.default.name}"
payload = "{\"tf-test\": \"test success\"}"
schedule_name = var.name
}

data "alicloud_fnf_schedules" "default" {
	flow_name = "${alicloud_fnf_flow.default.name}"
	%s
}
`, rand, strings.Join(pairs, " \n "))
}
