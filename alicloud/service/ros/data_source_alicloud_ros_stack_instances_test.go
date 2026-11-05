package ros_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudROSStackInstancesDataSource(t *testing.T) {
	resourceId := "data.alicloud_ros_stack_instances.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-rosstackinstance-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceRosStackInstancesDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"stack_group_name": "${alicloud_ros_stack_instance.default.stack_group_name}",
			"ids":              []string{"${alicloud_ros_stack_instance.default.id}"},
			"enable_details":   "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"stack_group_name": "${alicloud_ros_stack_instance.default.stack_group_name}",
			"ids":              []string{"${alicloud_ros_stack_instance.default.id}-fake"},
			"enable_details":   "true",
		}),
	}
	accountIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"stack_group_name":          "${alicloud_ros_stack_instance.default.stack_group_name}",
			"stack_instance_account_id": "${alicloud_ros_stack_instance.default.stack_instance_account_id}",
			"ids":                       []string{"${alicloud_ros_stack_instance.default.id}"},
			"enable_details":            "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"stack_group_name":          "${alicloud_ros_stack_instance.default.stack_group_name}",
			"stack_instance_account_id": "1234567890",
			"ids":                       []string{"${alicloud_ros_stack_instance.default.id}"},
			"enable_details":            "true",
		}),
	}
	regionIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"stack_group_name":         "${alicloud_ros_stack_instance.default.stack_group_name}",
			"stack_instance_region_id": "${alicloud_ros_stack_instance.default.stack_instance_region_id}",
			"ids":                      []string{"${alicloud_ros_stack_instance.default.id}"},
			"enable_details":           "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"stack_group_name":         "${alicloud_ros_stack_instance.default.stack_group_name}",
			"stack_instance_region_id": "${alicloud_ros_stack_instance.default.stack_instance_region_id}-fake",
			"ids":                      []string{"${alicloud_ros_stack_instance.default.id}"},
			"enable_details":           "true",
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"stack_group_name": "${alicloud_ros_stack_instance.default.stack_group_name}",
			"ids":              []string{"${alicloud_ros_stack_instance.default.id}"},
			"status":           "CURRENT",
			"enable_details":   "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"stack_group_name": "${alicloud_ros_stack_instance.default.stack_group_name}",
			"ids":              []string{"${alicloud_ros_stack_instance.default.id}"},
			"status":           "OUTDATED",
			"enable_details":   "true",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"stack_group_name":          "${alicloud_ros_stack_instance.default.stack_group_name}",
			"stack_instance_account_id": "${alicloud_ros_stack_instance.default.stack_instance_account_id}",
			"stack_instance_region_id":  "${alicloud_ros_stack_instance.default.stack_instance_region_id}",
			"ids":                       []string{"${alicloud_ros_stack_instance.default.id}"},
			"status":                    "CURRENT",
			"enable_details":            "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"ids.0":                             acctest.CHECKSET,
			"instances.#":                       "1",
			"instances.0.status":                "CURRENT",
			"instances.0.parameter_overrides.#": "1",
			"instances.0.parameter_overrides.0.parameter_key":   "VpcName",
			"instances.0.parameter_overrides.0.parameter_value": "VpcName",
			"instances.0.stack_group_id":                        acctest.CHECKSET,
			"instances.0.stack_group_name":                      acctest.CHECKSET,
			"instances.0.stack_id":                              acctest.CHECKSET,
			"instances.0.stack_instance_account_id":             acctest.CHECKSET,
			"instances.0.id":                                    acctest.CHECKSET,
			"instances.0.stack_instance_region_id":              acctest.CHECKSET,
			"instances.0.status_reason":                         "",
		}
	}

	var fakeRosStackInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var RosStackInstanceCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existRosStackInstanceMapFunc,
		FakeMapFunc:  fakeRosStackInstanceMapFunc,
	}

	RosStackInstanceCheckInfo.DataSourceTestCheck(t, rand, idsConf, accountIdConf, regionIdConf, statusConf, allConf)
}

func dataSourceRosStackInstancesDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_account" "this" {}

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
  stack_instance_account_id = data.alicloud_account.this.id
  stack_group_name          = alicloud_ros_stack_group.default.stack_group_name
  stack_instance_region_id  = data.alicloud_ros_regions.default.regions.0.region_id
  parameter_overrides {
    parameter_value = "VpcName"
    parameter_key   = "VpcName"
  }
}`, name)
}
