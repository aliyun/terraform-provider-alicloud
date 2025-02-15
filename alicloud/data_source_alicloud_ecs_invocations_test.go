package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudECSInvocationsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsInvocationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_invocation.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsInvocationsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecs_invocation.default.id}_fake"]`,
		}),
	}
	commandIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsInvocationsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecs_invocation.default.id}"]`,
			"command_id": `"${alicloud_ecs_invocation.default.command_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsInvocationsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecs_invocation.default.id}"]`,
			"command_id": `"${alicloud_ecs_invocation.default.command_id}_fake"`,
		}),
	}

	invokeStatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsInvocationsDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_ecs_invocation.default.id}"]`,
			"invoke_status": `"Finished"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsInvocationsDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_ecs_invocation.default.id}"]`,
			"invoke_status": `"Running"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsInvocationsDataSourceName(rand, map[string]string{
			"command_id":    `"${alicloud_ecs_invocation.default.command_id}"`,
			"ids":           `["${alicloud_ecs_invocation.default.id}"]`,
			"invoke_status": `"Finished"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsInvocationsDataSourceName(rand, map[string]string{
			"command_id":    `"${alicloud_ecs_invocation.default.command_id}_fake"`,
			"ids":           `["${alicloud_ecs_invocation.default.id}_fake"]`,
			"invoke_status": `"Running"`,
		}),
	}
	var existAlicloudEcsInvocationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                                   "1",
			"invocations.#":                                           "1",
			"invocations.0.command_id":                                CHECKSET,
			"invocations.0.command_type":                              CHECKSET,
			"invocations.0.create_time":                               CHECKSET,
			"invocations.0.parameters":                                CHECKSET,
			"invocations.0.repeat_mode":                               CHECKSET,
			"invocations.0.invocation_status":                         CHECKSET,
			"invocations.0.invoke_status":                             CHECKSET,
			"invocations.0.command_content":                           CHECKSET,
			"invocations.0.command_name":                              CHECKSET,
			"invocations.0.timed":                                     CHECKSET,
			"invocations.0.id":                                        CHECKSET,
			"invocations.0.invocation_id":                             CHECKSET,
			"invocations.0.frequency":                                 "",
			"invocations.0.invoke_instances.#":                        CHECKSET,
			"invocations.0.invoke_instances.0.creation_time":          CHECKSET,
			"invocations.0.invoke_instances.0.update_time":            CHECKSET,
			"invocations.0.invoke_instances.0.finish_time":            CHECKSET,
			"invocations.0.invoke_instances.0.instance_id":            CHECKSET,
			"invocations.0.invoke_instances.0.invocation_status":      CHECKSET,
			"invocations.0.invoke_instances.0.repeats":                CHECKSET,
			"invocations.0.invoke_instances.0.output":                 "",
			"invocations.0.invoke_instances.0.dropped":                CHECKSET,
			"invocations.0.invoke_instances.0.stop_time":              "",
			"invocations.0.invoke_instances.0.exit_code":              CHECKSET,
			"invocations.0.invoke_instances.0.start_time":             CHECKSET,
			"invocations.0.invoke_instances.0.error_info":             "",
			"invocations.0.invoke_instances.0.timed":                  CHECKSET,
			"invocations.0.invoke_instances.0.error_code":             "",
			"invocations.0.invoke_instances.0.instance_invoke_status": CHECKSET,
		}
	}
	var fakeAlicloudEcsInvocationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudEcsInvocationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecs_invocations.default",
		existMapFunc: existAlicloudEcsInvocationsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcsInvocationsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEcsInvocationsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, commandIdConf, invokeStatusConf, allConf)
}
func testAccCheckAlicloudEcsInvocationsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccInvocation-%d"
}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone = data.alicloud_zones.default.zones.0.id
  		image_id          = data.alicloud_images.default.images.0.id
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_efficiency"
  		vswitch_id                 = alicloud_vswitch.default.id
  		instance_name              = var.name
	}

	resource "alicloud_ecs_command" "default" {
  		name            = var.name
  		command_content = "bHMK"
  		description     = "For Terraform Test"
  		type            = "RunShellScript"
  		working_dir     = "/root"
	}

resource "alicloud_ecs_invocation" "default" {
	command_id = alicloud_ecs_command.default.id
	instance_id = [alicloud_instance.default.id]
}

data "alicloud_ecs_invocations" "default" {
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
