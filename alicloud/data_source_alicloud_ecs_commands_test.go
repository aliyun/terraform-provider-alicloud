package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECSCommandsDataSource(t *testing.T) {
	resourceId := "data.alicloud_ecs_commands.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccEcsCommandsTest%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEcsCommandsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_command.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_command.default.id}-fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
			"ids":        []string{"${alicloud_ecs_command.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
			"ids":        []string{"${alicloud_ecs_command.default.id}-fake"},
		}),
	}
	var existEcsCommandsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"ids.0":                       CHECKSET,
			"names.#":                     "1",
			"names.0":                     name,
			"commands.#":                  "1",
			"commands.0.id":               CHECKSET,
			"commands.0.command_content":  "bHMK",
			"commands.0.command_id":       CHECKSET,
			"commands.0.description":      "For Terraform Test",
			"commands.0.enable_parameter": "false",
			"commands.0.name":             name,
			"commands.0.type":             "RunShellScript",
			"commands.0.working_dir":      "/root",
		}
	}

	var fakeEcsCommandsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"commands.#": "0",
		}
	}

	var EcsCommandsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEcsCommandsMapFunc,
		fakeMapFunc:  fakeEcsCommandsMapFunc,
	}

	EcsCommandsInfo.dataSourceTestCheck(t, 0, idsConf, nameRegexConf)
}

func dataSourceEcsCommandsDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_ecs_command" "default" {
		name              = "%s"
		command_content   = "bHMK"
		description       = "For Terraform Test"
		type              = "RunShellScript"
		working_dir       = "/root"
	}`, name)
}
