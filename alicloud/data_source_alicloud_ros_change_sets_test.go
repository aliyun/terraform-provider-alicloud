package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudROSChangeSetsDataSource(t *testing.T) {
	resourceId := "data.alicloud_ros_change_sets.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccRosChangeSets%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceRosChangeSetsDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"stack_id":       "${alicloud_ros_change_set.default.stack_id}",
			"name_regex":     "${alicloud_ros_change_set.default.change_set_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"stack_id":       "${alicloud_ros_change_set.default.stack_id}",
			"name_regex":     "${alicloud_ros_change_set.default.change_set_name}-fake",
			"enable_details": "true",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"stack_id":       "${alicloud_ros_change_set.default.stack_id}",
			"status":         "CREATE_COMPLETE",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"stack_id":       "${alicloud_ros_change_set.default.stack_id}",
			"status":         "CREATE_PENDING",
			"enable_details": "true",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"stack_id":       "${alicloud_ros_change_set.default.stack_id}",
			"ids":            []string{"${alicloud_ros_change_set.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"stack_id":       "${alicloud_ros_change_set.default.stack_id}",
			"ids":            []string{"${alicloud_ros_change_set.default.id}-fake"},
			"enable_details": "true",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"stack_id":       "${alicloud_ros_change_set.default.stack_id}",
			"name_regex":     "${alicloud_ros_change_set.default.change_set_name}",
			"status":         "CREATE_COMPLETE",
			"ids":            []string{"${alicloud_ros_change_set.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"stack_id":       "${alicloud_ros_change_set.default.stack_id}",
			"name_regex":     "${alicloud_ros_change_set.default.change_set_name}-fake",
			"status":         "CREATE_IN_PROGRESS",
			"ids":            []string{"${alicloud_ros_change_set.default.id}-fake"},
			"enable_details": "true",
		}),
	}
	var existRosChangeSetsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"ids.0":                     CHECKSET,
			"names.#":                   "1",
			"names.0":                   name,
			"sets.#":                    "1",
			"sets.0.id":                 CHECKSET,
			"sets.0.change_set_id":      CHECKSET,
			"sets.0.change_set_name":    name,
			"sets.0.change_set_type":    "CREATE",
			"sets.0.description":        CHECKSET,
			"sets.0.disable_rollback":   "false",
			"sets.0.execution_status":   "AVAILABLE",
			"sets.0.stack_id":           CHECKSET,
			"sets.0.stack_name":         CHECKSET,
			"sets.0.status":             "CREATE_COMPLETE",
			"sets.0.template_body":      CHECKSET,
			"sets.0.timeout_in_minutes": "60",
		}
	}

	var fakeRosChangeSetsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"sets.#":  "0",
		}
	}

	var rosChangeSetsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existRosChangeSetsMapFunc,
		fakeMapFunc:  fakeRosChangeSetsMapFunc,
	}

	rosChangeSetsInfo.dataSourceTestCheck(t, 0, nameRegexConf, statusConf, idsConf, allConf)
}

func dataSourceRosChangeSetsDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_ros_change_set" "default" {
	  change_set_name = "%[1]s"
	  stack_name=   "%[1]sstack"
	  change_set_type= "CREATE"
	  description=   "Test From Terraform"
	  template_body=   "{\"ROSTemplateFormatVersion\":\"2015-09-01\"}"
	}
	`, name)
}
