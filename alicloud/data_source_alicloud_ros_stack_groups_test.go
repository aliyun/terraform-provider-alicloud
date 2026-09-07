package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudROSStackGroupsDataSource(t *testing.T) {
	resourceId := "data.alicloud_ros_stack_groups.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccRosStackGroups%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceRosStackGroupsDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_stack_group.default.stack_group_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_stack_group.default.stack_group_name}-fake",
			"enable_details": "true",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack_group.default.id}"},
			"status":         "ACTIVE",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack_group.default.id}"},
			"status":         "DELETED",
			"enable_details": "true",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack_group.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack_group.default.id}-fake"},
			"enable_details": "true",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_stack_group.default.stack_group_name}",
			"status":         "ACTIVE",
			"ids":            []string{"${alicloud_ros_stack_group.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_stack_group.default.stack_group_name}-fake",
			"status":         "DELETED",
			"ids":            []string{"${alicloud_ros_stack_group.default.id}-fake"},
			"enable_details": "true",
		}),
	}
	var existRosStackGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"ids.0":                             CHECKSET,
			"names.#":                           "1",
			"names.0":                           name,
			"groups.#":                          "1",
			"groups.0.id":                       CHECKSET,
			"groups.0.administration_role_name": CHECKSET,
			"groups.0.description":              "test for stack groups",
			"groups.0.parameters.#":             "0",
			"groups.0.stack_group_id":           CHECKSET,
			"groups.0.stack_group_name":         name,
			"groups.0.status":                   "ACTIVE",
			"groups.0.template_body":            CHECKSET,
		}
	}

	var fakeRosStackGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"groups.#": "0",
		}
	}

	var rosStackGroupsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existRosStackGroupsMapFunc,
		fakeMapFunc:  fakeRosStackGroupsMapFunc,
	}

	rosStackGroupsInfo.dataSourceTestCheck(t, 0, nameRegexConf, statusConf, idsConf, allConf)
}

func dataSourceRosStackGroupsDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_ros_stack_group" "default" {
	  stack_group_name=   "%s"
	  template_body=   "{\"ROSTemplateFormatVersion\":\"2015-09-01\"}"
      description = "test for stack groups"
	}
	`, name)
}
