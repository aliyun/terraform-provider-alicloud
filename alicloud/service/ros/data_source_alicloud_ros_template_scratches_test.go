package ros_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudROSTemplateScratchesDataSource(t *testing.T) {
	resourceId := "data.alicloud_ros_template_scratches.default"
	rand := tfacctest.RandIntRange(100000, 999999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.ROSSupportRegions)
	name := fmt.Sprintf("tf-testacc-templatescratch-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceRosTemplateScratchesDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_template_scratch.default.id}"},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_template_scratch.default.id}-fake"},
			"enable_details": "true",
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_template_scratch.default.id}"},
			"status":         "GENERATE_COMPLETE",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_ros_template_scratch.default.id}"},
			"status":         "GENERATE_FAILED",
			"enable_details": "true",
		}),
	}
	templateScratchTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":                   []string{"${alicloud_ros_template_scratch.default.id}"},
			"template_scratch_type": "ResourceImport",
			"enable_details":        "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":                   []string{"${alicloud_ros_template_scratch.default.id}"},
			"template_scratch_type": "ArchitectureReplication",
			"enable_details":        "true",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":                   []string{"${alicloud_ros_template_scratch.default.id}"},
			"status":                "GENERATE_COMPLETE",
			"template_scratch_type": "ResourceImport",
			"enable_details":        "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":                   []string{"${alicloud_ros_template_scratch.default.id}"},
			"status":                "GENERATE_FAILED",
			"template_scratch_type": "ArchitectureReplication",
			"enable_details":        "true",
		}),
	}
	var existRosTemplateScratchMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"ids.0":                               acctest.CHECKSET,
			"scratches.#":                         "1",
			"scratches.0.create_time":             acctest.CHECKSET,
			"scratches.0.description":             fmt.Sprintf("tf-testacc-templatescratch-%d", rand),
			"scratches.0.logical_id_strategy":     acctest.CHECKSET,
			"scratches.0.preference_parameters.#": "1",
			"scratches.0.preference_parameters.0.parameter_key":          "DeletionPolicy",
			"scratches.0.preference_parameters.0.parameter_value":        "Retain",
			"scratches.0.source_tag.#":                                   "0",
			"scratches.0.source_resource_group.#":                        "1",
			"scratches.0.source_resource_group.0.resource_group_id":      acctest.CHECKSET,
			"scratches.0.source_resource_group.0.resource_type_filter.#": "1",
			"scratches.0.source_resource_group.0.resource_type_filter.0": "ALIYUN::ECS::VPC",
			"scratches.0.source_resources.#":                             "0",
			"scratches.0.stacks.#":                                       "0",
			"scratches.0.status":                                         "GENERATE_COMPLETE",
			"scratches.0.id":                                             acctest.CHECKSET,
			"scratches.0.template_scratch_id":                            acctest.CHECKSET,
			"scratches.0.template_scratch_type":                          "ResourceImport",
		}
	}

	var fakeRosTemplateScratchMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"scratches.#": "0",
		}
	}

	var RosTemplateScratchCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existRosTemplateScratchMapFunc,
		FakeMapFunc:  fakeRosTemplateScratchMapFunc,
	}

	RosTemplateScratchCheckInfo.DataSourceTestCheck(t, rand, idsConf, statusConf, templateScratchTypeConf, allConf)
}

func dataSourceRosTemplateScratchesDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_ros_template_scratch" "default" {
  description           = var.name
  template_scratch_type = "ResourceImport"
  preference_parameters {
    parameter_key   = "DeletionPolicy"
    parameter_value = "Retain"
  }
  source_resource_group {
    resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids.0
    resource_type_filter = ["ALIYUN::ECS::VPC"]
  }
}`, name)
}
