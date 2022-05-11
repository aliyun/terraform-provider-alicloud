package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudROSStacksDataSource(t *testing.T) {
	resourceId := "data.alicloud_ros_stacks.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccRosStacks%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceRosStacksDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_stack.default.stack_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_stack.default.stack_name}-fake",
			"enable_details": "true",
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack.default.id}"},
			"enable_details": "true",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "ROS",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack.default.id}"},
			"enable_details": "true",
			"tags": map[string]string{
				"Created": "TF-fake",
				"For":     "update test fake",
			},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack.default.id}"},
			"status":         "CREATE_COMPLETE",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack.default.id}"},
			"status":         "ROLLBACK_IN_PROGRESS",
			"enable_details": "true",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack.default.id}-fake"},
			"enable_details": "true",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_stack.default.stack_name}",
			"status":         "CREATE_COMPLETE",
			"ids":            []string{"${alicloud_ros_stack.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_stack.default.stack_name}-fake",
			"status":         "CREATE_IN_PROGRESS",
			"ids":            []string{"${alicloud_ros_stack.default.id}-fake"},
			"enable_details": "true",
		}),
	}
	var existRosStacksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"ids.0":                         CHECKSET,
			"names.#":                       "1",
			"names.0":                       name,
			"stacks.#":                      "1",
			"stacks.0.id":                   CHECKSET,
			"stacks.0.deletion_protection":  "Disabled",
			"stacks.0.description":          CHECKSET,
			"stacks.0.disable_rollback":     "false",
			"stacks.0.drift_detection_time": "",
			"stacks.0.parent_stack_id":      "",
			"stacks.0.ram_role_name":        "",
			"stacks.0.root_stack_id":        "",
			"stacks.0.stack_drift_status":   "",
			"stacks.0.stack_name":           name,
			"stacks.0.tags.%":               "2",
			"stacks.0.parameters.#":         "6",
			"stacks.0.stack_policy_body":    CHECKSET,
			"stacks.0.timeout_in_minutes":   "60",
			"stacks.0.template_description": CHECKSET,
			"stacks.0.status":               "CREATE_COMPLETE",
			"stacks.0.status_reason":        CHECKSET,
		}
	}

	var fakeRosStacksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"stacks.#": "0",
		}
	}

	var rosStacksInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existRosStacksMapFunc,
		fakeMapFunc:  fakeRosStacksMapFunc,
	}

	rosStacksInfo.dataSourceTestCheck(t, 0, nameRegexConf, tagsConf, statusConf, idsConf, allConf)
}

func dataSourceRosStacksDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_ros_stack" "default" {
	  stack_name=   "%s"
	  template_body=   "{\"ROSTemplateFormatVersion\":\"2015-09-01\"}"
      stack_policy_body = "{\"Statement\": [{\"Action\": \"Update:Delete\", \"Resource\": \"*\", \"Effect\": \"Allow\", \"Principal\": \"*\"}]}"
	   tags = {
   			"Created": "TF",
			"For":     "ROS",
 		}
	}
	`, name)
}
