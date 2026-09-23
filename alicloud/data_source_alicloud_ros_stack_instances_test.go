package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudROSStackInstancesDataSource(t *testing.T) {
	resourceId := "data.alicloud_ros_stack_instances.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-rosstackinstance-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceRosStackInstancesDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"stack_group_name": "${alicloud_ros_stack_instance.default.stack_group_name}",
			"ids":              []string{"${alicloud_ros_stack_instance.default.id}"},
			"enable_details":   "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"stack_group_name": "${alicloud_ros_stack_instance.default.stack_group_name}",
			"ids":              []string{"${alicloud_ros_stack_instance.default.id}-fake"},
			"enable_details":   "true",
		}),
	}
	accountIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"stack_group_name":          "${alicloud_ros_stack_instance.default.stack_group_name}",
			"stack_instance_account_id": "${alicloud_ros_stack_instance.default.stack_instance_account_id}",
			"ids":                       []string{"${alicloud_ros_stack_instance.default.id}"},
			"enable_details":            "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"stack_group_name":          "${alicloud_ros_stack_instance.default.stack_group_name}",
			"stack_instance_account_id": "1234567890",
			"ids":                       []string{"${alicloud_ros_stack_instance.default.id}"},
			"enable_details":            "true",
		}),
	}
	regionIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"stack_group_name":         "${alicloud_ros_stack_instance.default.stack_group_name}",
			"stack_instance_region_id": "${alicloud_ros_stack_instance.default.stack_instance_region_id}",
			"ids":                      []string{"${alicloud_ros_stack_instance.default.id}"},
			"enable_details":           "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"stack_group_name":         "${alicloud_ros_stack_instance.default.stack_group_name}",
			"stack_instance_region_id": "${alicloud_ros_stack_instance.default.stack_instance_region_id}-fake",
			"ids":                      []string{"${alicloud_ros_stack_instance.default.id}"},
			"enable_details":           "true",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"stack_group_name": "${alicloud_ros_stack_instance.default.stack_group_name}",
			"ids":              []string{"${alicloud_ros_stack_instance.default.id}"},
			"status":           "CURRENT",
			"enable_details":   "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"stack_group_name": "${alicloud_ros_stack_instance.default.stack_group_name}",
			"ids":              []string{"${alicloud_ros_stack_instance.default.id}"},
			"status":           "OUTDATED",
			"enable_details":   "true",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"stack_group_name":          "${alicloud_ros_stack_instance.default.stack_group_name}",
			"stack_instance_account_id": "${alicloud_ros_stack_instance.default.stack_instance_account_id}",
			"stack_instance_region_id":  "${alicloud_ros_stack_instance.default.stack_instance_region_id}",
			"ids":                       []string{"${alicloud_ros_stack_instance.default.id}"},
			"status":                    "CURRENT",
			"enable_details":            "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"stack_group_name":          "${alicloud_ros_stack_instance.default.stack_group_name}",
			"stack_instance_account_id": "1234567890",
			"stack_instance_region_id":  "${alicloud_ros_stack_instance.default.stack_instance_region_id}-fake",
			"ids":                       []string{"${alicloud_ros_stack_instance.default.id}"},
			"status":                    "OUTDATED",
			"enable_details":            "true",
		}),
	}
	var existRosStackInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"ids.0":                             CHECKSET,
			"instances.#":                       "1",
			"instances.0.status":                "CURRENT",
			"instances.0.parameter_overrides.#": "1",
			"instances.0.parameter_overrides.0.parameter_key":   "VpcName",
			"instances.0.parameter_overrides.0.parameter_value": "VpcName",
			"instances.0.stack_group_id":                        CHECKSET,
			"instances.0.stack_group_name":                      CHECKSET,
			"instances.0.stack_id":                              CHECKSET,
			"instances.0.stack_instance_account_id":             CHECKSET,
			"instances.0.id":                                    CHECKSET,
			"instances.0.stack_instance_region_id":              CHECKSET,
			"instances.0.status_reason":                         "",
		}
	}

	var fakeRosStackInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var RosStackInstanceCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existRosStackInstanceMapFunc,
		fakeMapFunc:  fakeRosStackInstanceMapFunc,
	}

	RosStackInstanceCheckInfo.dataSourceTestCheck(t, rand, idsConf, accountIdConf, regionIdConf, statusConf, allConf)
}

func dataSourceRosStackInstancesDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_ros_stack_group" "default" {
  stack_group_name=   var.name
  template_body=   "{\"ROSTemplateFormatVersion\":\"2015-09-01\", \"Parameters\": {\"VpcName\": {\"Type\": \"String\"},\"InstanceType\": {\"Type\": \"String\"}}}"
  description = "test for stack groups"
  parameters {
	parameter_key =   "VpcName"
	parameter_value = "VpcName"
  }
  parameters {
	parameter_key =   "InstanceType"
	parameter_value = "InstanceType"
  }
}

data "alicloud_ros_regions" "default" {}

resource "alicloud_ros_stack_instance" "default" {
  stack_instance_account_id = "%s"
  stack_group_name          = alicloud_ros_stack_group.default.stack_group_name
  stack_instance_region_id  = data.alicloud_ros_regions.default.regions.0.region_id
  parameter_overrides {
    parameter_value = "VpcName"
    parameter_key   = "VpcName"
  }
}`, name, os.Getenv("ALICLOUD_ACCOUNT_ID"))
}
