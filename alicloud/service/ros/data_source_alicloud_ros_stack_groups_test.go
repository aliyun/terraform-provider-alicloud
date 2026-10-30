package ros_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudROSStackGroupsDataSource(t *testing.T) {
	resourceId := "data.alicloud_ros_stack_groups.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccRosStackGroups%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceRosStackGroupsDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_stack_group.default.stack_group_name}",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_stack_group.default.stack_group_name}-fake",
			"enable_details": "true",
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack_group.default.id}"},
			"status":         "ACTIVE",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack_group.default.id}"},
			"status":         "DELETED",
			"enable_details": "true",
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack_group.default.id}"},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_stack_group.default.id}-fake"},
			"enable_details": "true",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_stack_group.default.stack_group_name}",
			"status":         "ACTIVE",
			"ids":            []string{"${alicloud_ros_stack_group.default.id}"},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_ros_stack_group.default.stack_group_name}-fake",
			"status":         "DELETED",
			"ids":            []string{"${alicloud_ros_stack_group.default.id}-fake"},
			"enable_details": "true",
		}),
	}
	var existRosStackGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"ids.0":                             acctest.CHECKSET,
			"names.#":                           "1",
			"names.0":                           name,
			"groups.#":                          "1",
			"groups.0.id":                       acctest.CHECKSET,
			"groups.0.administration_role_name": acctest.CHECKSET,
			"groups.0.description":              "test for stack groups",
			"groups.0.parameters.#":             "0",
			"groups.0.stack_group_id":           acctest.CHECKSET,
			"groups.0.stack_group_name":         name,
			"groups.0.status":                   "ACTIVE",
			"groups.0.template_body":            acctest.CHECKSET,
		}
	}

	var fakeRosStackGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"groups.#": "0",
		}
	}

	var rosStackGroupsInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existRosStackGroupsMapFunc,
		FakeMapFunc:  fakeRosStackGroupsMapFunc,
	}

	rosStackGroupsInfo.DataSourceTestCheck(t, 0, nameRegexConf, statusConf, idsConf, allConf)
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
