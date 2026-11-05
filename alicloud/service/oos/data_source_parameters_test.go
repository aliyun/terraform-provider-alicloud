package oos_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudOosParametersDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_oos_parameters.default"
	name := fmt.Sprintf("tf-testAcc-OosParameter%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceOosParametersConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_oos_parameter.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_oos_parameter.default.id}_fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_oos_parameter.default.parameter_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_oos_parameter.default.parameter_name}_fake",
		}),
	}

	parameterNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"parameter_name": "${alicloud_oos_parameter.default.parameter_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"parameter_name": "${alicloud_oos_parameter.default.parameter_name}_fake",
		}),
	}

	typeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"type": "${alicloud_oos_parameter.default.type}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"type": "StringList",
		}),
	}

	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": "${alicloud_oos_parameter.default.resource_group_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Parameter",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Parameter_Fake",
			},
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_oos_parameter.default.id}"},
			"name_regex":        "${alicloud_oos_parameter.default.parameter_name}",
			"parameter_name":    "${alicloud_oos_parameter.default.parameter_name}",
			"type":              "${alicloud_oos_parameter.default.type}",
			"resource_group_id": "${alicloud_oos_parameter.default.resource_group_id}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Parameter",
			},
			"sort_field": "Name",
			"sort_order": "Ascending",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_oos_parameter.default.id}_fake"},
			"name_regex":        "${alicloud_oos_parameter.default.parameter_name}_fake",
			"parameter_name":    "${alicloud_oos_parameter.default.parameter_name}_fake",
			"type":              "StringList",
			"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Parameter_Fake",
			},
			"sort_field": "Name",
			"sort_order": "Ascending",
		}),
	}

	var existAliCloudOosParametersDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"parameters.#":                   "1",
			"parameters.0.id":                acctest.CHECKSET,
			"parameters.0.parameter_name":    acctest.CHECKSET,
			"parameters.0.parameter_id":      acctest.CHECKSET,
			"parameters.0.type":              acctest.CHECKSET,
			"parameters.0.parameter_version": acctest.CHECKSET,
			"parameters.0.share_type":        acctest.CHECKSET,
			"parameters.0.resource_group_id": acctest.CHECKSET,
			"parameters.0.description":       acctest.CHECKSET,
			"parameters.0.tags.%":            "2",
			"parameters.0.tags.Created":      "TF",
			"parameters.0.tags.For":          "Parameter",
			"parameters.0.constraints":       "",
			"parameters.0.value":             "",
			"parameters.0.created_by":        acctest.CHECKSET,
			"parameters.0.create_time":       acctest.CHECKSET,
			"parameters.0.updated_by":        acctest.CHECKSET,
			"parameters.0.updated_date":      acctest.CHECKSET,
		}
	}

	var fakeAliCloudOosParametersDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"names.#":      "0",
			"parameters.#": "0",
		}
	}

	var aliCloudOosParametersCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_oos_parameters.default",
		ExistMapFunc: existAliCloudOosParametersDefaultDataSourceMapFunc,
		FakeMapFunc:  fakeAliCloudOosParametersDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	aliCloudOosParametersCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, parameterNameConf, typeConf, resourceGroupIdConf, tagsConf, allConf)
}

func TestAccAliCloudOosParametersDataSource_basic1(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_oos_parameters.default"
	name := fmt.Sprintf("tf-testAcc-OosParameter%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceOosParametersConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_oos_parameter.default.id}"},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_oos_parameter.default.id}"},
			"enable_details": "false",
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_oos_parameter.default.parameter_name}",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_oos_parameter.default.parameter_name}",
			"enable_details": "false",
		}),
	}

	parameterNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"parameter_name": "${alicloud_oos_parameter.default.parameter_name}",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"parameter_name": "${alicloud_oos_parameter.default.parameter_name}",
			"enable_details": "false",
		}),
	}

	typeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"type":           "${alicloud_oos_parameter.default.type}",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"type":           "${alicloud_oos_parameter.default.type}",
			"enable_details": "false",
		}),
	}

	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": "${alicloud_oos_parameter.default.resource_group_id}",
			"enable_details":    "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": "${alicloud_oos_parameter.default.resource_group_id}",
			"enable_details":    "false",
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Parameter",
			},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Parameter",
			},
			"enable_details": "false",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_oos_parameter.default.id}"},
			"name_regex":        "${alicloud_oos_parameter.default.parameter_name}",
			"parameter_name":    "${alicloud_oos_parameter.default.parameter_name}",
			"type":              "${alicloud_oos_parameter.default.type}",
			"resource_group_id": "${alicloud_oos_parameter.default.resource_group_id}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Parameter",
			},
			"sort_field":     "Name",
			"sort_order":     "Ascending",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_oos_parameter.default.id}"},
			"name_regex":        "${alicloud_oos_parameter.default.parameter_name}",
			"parameter_name":    "${alicloud_oos_parameter.default.parameter_name}",
			"type":              "${alicloud_oos_parameter.default.type}",
			"resource_group_id": "${alicloud_oos_parameter.default.resource_group_id}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Parameter",
			},
			"sort_field":     "Name",
			"sort_order":     "Ascending",
			"enable_details": "false",
		}),
	}

	var existAliCloudOosParametersDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"parameters.#":                   "1",
			"parameters.0.id":                acctest.CHECKSET,
			"parameters.0.parameter_name":    acctest.CHECKSET,
			"parameters.0.parameter_id":      acctest.CHECKSET,
			"parameters.0.type":              acctest.CHECKSET,
			"parameters.0.parameter_version": acctest.CHECKSET,
			"parameters.0.share_type":        acctest.CHECKSET,
			"parameters.0.resource_group_id": acctest.CHECKSET,
			"parameters.0.description":       acctest.CHECKSET,
			"parameters.0.tags.%":            "2",
			"parameters.0.tags.Created":      "TF",
			"parameters.0.tags.For":          "Parameter",
			"parameters.0.constraints":       acctest.CHECKSET,
			"parameters.0.value":             acctest.CHECKSET,
			"parameters.0.created_by":        acctest.CHECKSET,
			"parameters.0.create_time":       acctest.CHECKSET,
			"parameters.0.updated_by":        acctest.CHECKSET,
			"parameters.0.updated_date":      acctest.CHECKSET,
		}
	}

	var fakeAliCloudOosParametersDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"parameters.#":                   "1",
			"parameters.0.id":                acctest.CHECKSET,
			"parameters.0.parameter_name":    acctest.CHECKSET,
			"parameters.0.parameter_id":      acctest.CHECKSET,
			"parameters.0.type":              acctest.CHECKSET,
			"parameters.0.parameter_version": acctest.CHECKSET,
			"parameters.0.share_type":        acctest.CHECKSET,
			"parameters.0.resource_group_id": acctest.CHECKSET,
			"parameters.0.description":       acctest.CHECKSET,
			"parameters.0.tags.%":            "2",
			"parameters.0.tags.Created":      "TF",
			"parameters.0.tags.For":          "Parameter",
			"parameters.0.constraints":       "",
			"parameters.0.value":             "",
			"parameters.0.created_by":        acctest.CHECKSET,
			"parameters.0.create_time":       acctest.CHECKSET,
			"parameters.0.updated_by":        acctest.CHECKSET,
			"parameters.0.updated_date":      acctest.CHECKSET,
		}
	}

	var aliCloudOosParametersCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_oos_parameters.default",
		ExistMapFunc: existAliCloudOosParametersDefaultDataSourceMapFunc,
		FakeMapFunc:  fakeAliCloudOosParametersDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	aliCloudOosParametersCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, parameterNameConf, typeConf, resourceGroupIdConf, tagsConf, allConf)
}

func dataSourceOosParametersConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	resource "alicloud_oos_parameter" "default" {
  		parameter_name = var.name
  		value          = "tf-testacc-oos_parameter"
  		type           = "String"
  		description    = var.name
  		constraints    = "{\"AllowedValues\":[\"tf-testacc-oos_parameter\"], \"AllowedPattern\": \"tf-testacc-oos_parameter\", \"MinLength\": 1, \"MaxLength\": 100}"
  		tags = {
    		Created = "TF"
    		For     = "Parameter"
  		}
	}
`, name)
}
