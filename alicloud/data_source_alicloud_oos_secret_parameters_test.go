package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOOSSecretParametersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.OOSSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosSecretParametersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_secret_parameter.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudOosSecretParametersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_secret_parameter.default.id}_fake"]`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosSecretParametersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_oos_secret_parameter.default.id}"]`,
			"resource_group_id": `"${alicloud_oos_secret_parameter.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudOosSecretParametersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_oos_secret_parameter.default.id}"]`,
			"resource_group_id": `"fake"`,
		}),
	}
	secretParameterNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosSecretParametersDataSourceName(rand, map[string]string{
			"secret_parameter_name": `"${alicloud_oos_secret_parameter.default.secret_parameter_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudOosSecretParametersDataSourceName(rand, map[string]string{
			"secret_parameter_name": `"${alicloud_oos_secret_parameter.default.secret_parameter_name}_fake"`,
		}),
	}
	sortOrderConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosSecretParametersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_oos_secret_parameter.default.id}"]`,
			"sort_field": `"Name"`,
			"sort_order": `"Ascending"`,
		}),
		fakeConfig: "",
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosSecretParametersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_secret_parameter.default.id}"]`,
			"tags": `{
				"Created" = "TF"
				"For" = "Acceptance-test"
		}`,
		}),
		fakeConfig: testAccCheckAlicloudOosSecretParametersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_secret_parameter.default.id}"]`,
			"tags": `{
				"Created" = "TF-fake"
				"For" = "Acceptance-test"
			}`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosSecretParametersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_oos_secret_parameter.default.secret_parameter_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudOosSecretParametersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_oos_secret_parameter.default.secret_parameter_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosSecretParametersDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_oos_secret_parameter.default.id}"]`,
			"name_regex":            `"${alicloud_oos_secret_parameter.default.secret_parameter_name}"`,
			"resource_group_id":     `"${alicloud_oos_secret_parameter.default.resource_group_id}"`,
			"sort_field":            `"Name"`,
			"sort_order":            `"Ascending"`,
			"secret_parameter_name": `"${alicloud_oos_secret_parameter.default.secret_parameter_name}"`,
			"tags": `{
				"Created" = "TF"
				"For" = "Acceptance-test"
		}`,
		}),
		fakeConfig: testAccCheckAlicloudOosSecretParametersDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_oos_secret_parameter.default.id}"]`,
			"name_regex":            `"${alicloud_oos_secret_parameter.default.secret_parameter_name}_fake"`,
			"resource_group_id":     `"fake"`,
			"secret_parameter_name": `"${alicloud_oos_secret_parameter.default.secret_parameter_name}_fake"`,
			"tags": `{
				"Created" = "TF-fake"
				"For" = "Acceptance-test"
			}`,
		}),
	}
	var existAlicloudOosSecretParametersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"names.#":                            "1",
			"parameters.#":                       "1",
			"parameters.0.id":                    fmt.Sprintf("tf-testAccSecretParameter-%d", rand),
			"parameters.0.constraints":           CHECKSET,
			"parameters.0.create_time":           CHECKSET,
			"parameters.0.created_by":            CHECKSET,
			"parameters.0.description":           "oos_secret_parameter_description",
			"parameters.0.key_id":                "",
			"parameters.0.parameter_version":     CHECKSET,
			"parameters.0.resource_group_id":     CHECKSET,
			"parameters.0.secret_parameter_id":   CHECKSET,
			"parameters.0.secret_parameter_name": fmt.Sprintf("tf-testAccSecretParameter-%d", rand),
			"parameters.0.share_type":            CHECKSET,
			"parameters.0.type":                  "Secret",
			"parameters.0.tags.%":                "2",
			"parameters.0.tags.Created":          "TF",
			"parameters.0.tags.For":              "Acceptance-test",
			"parameters.0.updated_by":            CHECKSET,
			"parameters.0.updated_date":          CHECKSET,
		}
	}
	var fakeAlicloudOosSecretParametersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudOosSecretParametersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_oos_secret_parameters.default",
		existMapFunc: existAlicloudOosSecretParametersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudOosSecretParametersDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudOosSecretParametersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, resourceGroupIdConf, secretParameterNameConf, sortOrderConf, tagsConf, nameRegexConf, allConf)
}

func testAccCheckAlicloudOosSecretParametersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tf-testAccSecretParameter-%d"
}

resource "alicloud_oos_secret_parameter" "default" {
  constraints           = "{\"AllowedValues\":[\"tf-testacc-oos_secret_parameter\"], \"AllowedPattern\": \"tf-testacc-oos_secret_parameter\", \"MinLength\": 1, \"MaxLength\": 100}"
  secret_parameter_name = var.name
  value                 = "tf-testacc-oos_secret_parameter"
  type                  = "Secret"
  description           = "oos_secret_parameter_description"
  tags                  = {
    Created = "TF"
    For     = "Acceptance-test"
  }
}

data "alicloud_oos_secret_parameters" "default" {
  enable_details = true
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
