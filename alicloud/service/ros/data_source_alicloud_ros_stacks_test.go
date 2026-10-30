package ros_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudROSStacksDataSource(t *testing.T) {
	resourceId := "data.alicloud_ros_stacks.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccRosStacks%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceRosStacksDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_stack.default.stack_name}",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_stack.default.stack_name}-fake",
			"enable_details": "true",
		}),
	}
	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack.default.id}"},
			"enable_details": "true",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "ROS",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack.default.id}"},
			"enable_details": "true",
			"tags": map[string]string{
				"Created": "TF-fake",
				"For":     "update test fake",
			},
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack.default.id}"},
			"status":         "CREATE_COMPLETE",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack.default.id}"},
			"status":         "ROLLBACK_IN_PROGRESS",
			"enable_details": "true",
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack.default.id}"},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack.default.id}-fake"},
			"enable_details": "true",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_stack.default.stack_name}",
			"status":         "CREATE_COMPLETE",
			"ids":            []string{"${alicloud_ros_stack.default.id}"},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_stack.default.stack_name}-fake",
			"status":         "CREATE_IN_PROGRESS",
			"ids":            []string{"${alicloud_ros_stack.default.id}-fake"},
			"enable_details": "true",
		}),
	}
	var existRosStacksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"ids.0":                         acctest.CHECKSET,
			"names.#":                       "1",
			"names.0":                       name,
			"stacks.#":                      "1",
			"stacks.0.id":                   acctest.CHECKSET,
			"stacks.0.deletion_protection":  "Disabled",
			"stacks.0.description":          acctest.CHECKSET,
			"stacks.0.disable_rollback":     "false",
			"stacks.0.drift_detection_time": "",
			"stacks.0.parent_stack_id":      "",
			"stacks.0.ram_role_name":        "",
			"stacks.0.root_stack_id":        "",
			"stacks.0.stack_drift_status":   "",
			"stacks.0.stack_name":           name,
			"stacks.0.tags.%":               "2",
			"stacks.0.parameters.#":         "7",
			"stacks.0.stack_policy_body":    acctest.CHECKSET,
			"stacks.0.timeout_in_minutes":   "60",
			"stacks.0.template_description": acctest.CHECKSET,
			"stacks.0.status":               "CREATE_COMPLETE",
			"stacks.0.status_reason":        acctest.CHECKSET,
		}
	}

	var fakeRosStacksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"stacks.#": "0",
		}
	}

	var rosStacksInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existRosStacksMapFunc,
		FakeMapFunc:  fakeRosStacksMapFunc,
	}

	rosStacksInfo.DataSourceTestCheck(t, 0, nameRegexConf, tagsConf, statusConf, idsConf, allConf)
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
