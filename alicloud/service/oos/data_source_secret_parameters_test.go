package oos_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudOosSecretParametersDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_oos_secret_parameters.default"
	name := fmt.Sprintf("tf-testAcc-OosSecretParameter%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceOosSecretParametersConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_oos_secret_parameter.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_oos_secret_parameter.default.id}_fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_oos_secret_parameter.default.secret_parameter_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_oos_secret_parameter.default.secret_parameter_name}_fake",
		}),
	}

	secretParameterNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"secret_parameter_name": "${alicloud_oos_secret_parameter.default.secret_parameter_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"secret_parameter_name": "${alicloud_oos_secret_parameter.default.secret_parameter_name}_fake",
		}),
	}

	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": "${alicloud_oos_secret_parameter.default.resource_group_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "SecretParameter",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "SecretParameter_Fake",
			},
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":                   []string{"${alicloud_oos_secret_parameter.default.id}"},
			"name_regex":            "${alicloud_oos_secret_parameter.default.secret_parameter_name}",
			"secret_parameter_name": "${alicloud_oos_secret_parameter.default.secret_parameter_name}",
			"resource_group_id":     "${alicloud_oos_secret_parameter.default.resource_group_id}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "SecretParameter",
			},
			"sort_field": "Name",
			"sort_order": "Ascending",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":                   []string{"${alicloud_oos_secret_parameter.default.id}_fake"},
			"name_regex":            "${alicloud_oos_secret_parameter.default.secret_parameter_name}_fake",
			"secret_parameter_name": "${alicloud_oos_secret_parameter.default.secret_parameter_name}_fake",
			"resource_group_id":     "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "SecretParameter_Fake",
			},
			"sort_field": "Name",
			"sort_order": "Ascending",
		}),
	}

	var existAliCloudOosSecretParametersDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"names.#":                            "1",
			"parameters.#":                       "1",
			"parameters.0.id":                    acctest.CHECKSET,
			"parameters.0.secret_parameter_id":   acctest.CHECKSET,
			"parameters.0.type":                  acctest.CHECKSET,
			"parameters.0.parameter_version":     acctest.CHECKSET,
			"parameters.0.share_type":            acctest.CHECKSET,
			"parameters.0.key_id":                acctest.CHECKSET,
			"parameters.0.resource_group_id":     acctest.CHECKSET,
			"parameters.0.secret_parameter_name": acctest.CHECKSET,
			"parameters.0.description":           acctest.CHECKSET,
			"parameters.0.tags.%":                "2",
			"parameters.0.tags.Created":          "TF",
			"parameters.0.tags.For":              "SecretParameter",
			"parameters.0.constraints":           "",
			"parameters.0.value":                 "",
			"parameters.0.created_by":            acctest.CHECKSET,
			"parameters.0.create_time":           acctest.CHECKSET,
			"parameters.0.updated_by":            acctest.CHECKSET,
			"parameters.0.updated_date":          acctest.CHECKSET,
		}
	}

	var fakeAliCloudOosSecretParametersDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"names.#":      "0",
			"parameters.#": "0",
		}
	}

	var aliCloudOosSecretParametersCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_oos_secret_parameters.default",
		ExistMapFunc: existAliCloudOosSecretParametersDefaultDataSourceMapFunc,
		FakeMapFunc:  fakeAliCloudOosSecretParametersDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	aliCloudOosSecretParametersCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, secretParameterNameConf, resourceGroupIdConf, tagsConf, allConf)
}

func TestAccAliCloudOosSecretParametersDataSource_basic1(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_oos_secret_parameters.default"
	name := fmt.Sprintf("tf-testAcc-OosSecretParameter%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceOosSecretParametersConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":             []string{"${alicloud_oos_secret_parameter.default.id}"},
			"enable_details":  "true",
			"with_decryption": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_oos_secret_parameter.default.id}"},
			"enable_details": "true",
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":      "${alicloud_oos_secret_parameter.default.secret_parameter_name}",
			"enable_details":  "true",
			"with_decryption": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_oos_secret_parameter.default.secret_parameter_name}",
			"enable_details": "true",
		}),
	}

	secretParameterNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"secret_parameter_name": "${alicloud_oos_secret_parameter.default.secret_parameter_name}",
			"enable_details":        "true",
			"with_decryption":       "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"secret_parameter_name": "${alicloud_oos_secret_parameter.default.secret_parameter_name}",
			"enable_details":        "true",
		}),
	}

	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": "${alicloud_oos_secret_parameter.default.resource_group_id}",
			"enable_details":    "true",
			"with_decryption":   "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": "${alicloud_oos_secret_parameter.default.resource_group_id}",
			"enable_details":    "true",
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "SecretParameter",
			},
			"enable_details":  "true",
			"with_decryption": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "SecretParameter",
			},
			"enable_details": "true",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":                   []string{"${alicloud_oos_secret_parameter.default.id}"},
			"name_regex":            "${alicloud_oos_secret_parameter.default.secret_parameter_name}",
			"secret_parameter_name": "${alicloud_oos_secret_parameter.default.secret_parameter_name}",
			"resource_group_id":     "${alicloud_oos_secret_parameter.default.resource_group_id}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "SecretParameter",
			},
			"sort_field":      "Name",
			"sort_order":      "Ascending",
			"enable_details":  "true",
			"with_decryption": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":                   []string{"${alicloud_oos_secret_parameter.default.id}"},
			"name_regex":            "${alicloud_oos_secret_parameter.default.secret_parameter_name}",
			"secret_parameter_name": "${alicloud_oos_secret_parameter.default.secret_parameter_name}",
			"resource_group_id":     "${alicloud_oos_secret_parameter.default.resource_group_id}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "SecretParameter",
			},
			"sort_field":     "Name",
			"sort_order":     "Ascending",
			"enable_details": "true",
		}),
	}

	var existAliCloudOosSecretParametersDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"names.#":                            "1",
			"parameters.#":                       "1",
			"parameters.0.id":                    acctest.CHECKSET,
			"parameters.0.secret_parameter_id":   acctest.CHECKSET,
			"parameters.0.type":                  acctest.CHECKSET,
			"parameters.0.parameter_version":     acctest.CHECKSET,
			"parameters.0.share_type":            acctest.CHECKSET,
			"parameters.0.key_id":                acctest.CHECKSET,
			"parameters.0.resource_group_id":     acctest.CHECKSET,
			"parameters.0.secret_parameter_name": acctest.CHECKSET,
			"parameters.0.description":           acctest.CHECKSET,
			"parameters.0.tags.%":                "2",
			"parameters.0.tags.Created":          "TF",
			"parameters.0.tags.For":              "SecretParameter",
			"parameters.0.constraints":           acctest.CHECKSET,
			"parameters.0.value":                 acctest.CHECKSET,
			"parameters.0.created_by":            acctest.CHECKSET,
			"parameters.0.create_time":           acctest.CHECKSET,
			"parameters.0.updated_by":            acctest.CHECKSET,
			"parameters.0.updated_date":          acctest.CHECKSET,
		}
	}

	var fakeAliCloudOosSecretParametersDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"names.#":                            "1",
			"parameters.#":                       "1",
			"parameters.0.id":                    acctest.CHECKSET,
			"parameters.0.secret_parameter_id":   acctest.CHECKSET,
			"parameters.0.type":                  acctest.CHECKSET,
			"parameters.0.parameter_version":     acctest.CHECKSET,
			"parameters.0.share_type":            acctest.CHECKSET,
			"parameters.0.key_id":                acctest.CHECKSET,
			"parameters.0.resource_group_id":     acctest.CHECKSET,
			"parameters.0.secret_parameter_name": acctest.CHECKSET,
			"parameters.0.description":           acctest.CHECKSET,
			"parameters.0.tags.%":                "2",
			"parameters.0.tags.Created":          "TF",
			"parameters.0.tags.For":              "SecretParameter",
			"parameters.0.constraints":           acctest.CHECKSET,
			"parameters.0.value":                 "",
			"parameters.0.created_by":            acctest.CHECKSET,
			"parameters.0.create_time":           acctest.CHECKSET,
			"parameters.0.updated_by":            acctest.CHECKSET,
			"parameters.0.updated_date":          acctest.CHECKSET,
		}
	}

	var aliCloudOosSecretParametersCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_oos_secret_parameters.default",
		ExistMapFunc: existAliCloudOosSecretParametersDefaultDataSourceMapFunc,
		FakeMapFunc:  fakeAliCloudOosSecretParametersDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	aliCloudOosSecretParametersCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, secretParameterNameConf, resourceGroupIdConf, tagsConf, allConf)
}

func dataSourceOosSecretParametersConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
	}

	resource "alicloud_kms_key" "default" {
  		status                 = "Enabled"
  		pending_window_in_days = 7
	}

	resource "alicloud_oos_secret_parameter" "default" {
  		secret_parameter_name = var.name
  		value                 = "tf-testacc-oos_secret_parameter"
  		type                  = "Secret"
		key_id                = alicloud_kms_key.default.id
  		description           = var.name
		constraints           = "{\"AllowedValues\":[\"tf-testacc-oos_secret_parameter\"], \"AllowedPattern\": \"tf-testacc-oos_secret_parameter\", \"MinLength\": 1, \"MaxLength\": 100}"
  		tags = {
    		Created = "TF"
    		For     = "SecretParameter"
  		}
	}
`, name)
}
