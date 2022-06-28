package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOOSParametersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.OOSSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosParametersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_parameter.default.parameter_name}"]`,
		}),
		fakeConfig: testAccCheckAlicloudOosParametersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_parameter.default.parameter_name}_fake"]`,
		}),
	}
	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosParametersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_oos_parameter.default.parameter_name}"]`,
			"parameter_name": `"${alicloud_oos_parameter.default.parameter_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudOosParametersDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_oos_parameter.default.parameter_name}"]`,
			"parameter_name": `"${alicloud_oos_parameter.default.parameter_name}_fake"`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosParametersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_oos_parameter.default.parameter_name}"]`,
			"resource_group_id": `"${alicloud_oos_parameter.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudOosParametersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_oos_parameter.default.parameter_name}"]`,
			"resource_group_id": `"fake"`,
		}),
	}
	sortOrderConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosParametersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_oos_parameter.default.parameter_name}"]`,
			"sort_field": `"Name"`,
			"sort_order": `"Ascending"`,
		}),
		fakeConfig: "",
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosParametersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_parameter.default.parameter_name}"]`,
			"tags": `{
				"Created" = "TF"
				"For" = "OosParameter"
		    }`,
		}),
		fakeConfig: testAccCheckAlicloudOosParametersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_oos_parameter.default.parameter_name}"]`,
			"tags": `{
				"Created" = "OosParameter"
				"For" = "TF"
			}`,
		}),
	}
	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosParametersDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_oos_parameter.default.parameter_name}"]`,
			"type": `"String"`,
		}),
		fakeConfig: testAccCheckAlicloudOosParametersDataSourceName(rand, map[string]string{
			"ids":  `["${alicloud_oos_parameter.default.parameter_name}"]`,
			"type": `"StringList"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosParametersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_oos_parameter.default.parameter_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudOosParametersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_oos_parameter.default.parameter_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudOosParametersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_oos_parameter.default.parameter_name}"]`,
			"parameter_name":    `"${alicloud_oos_parameter.default.parameter_name}"`,
			"name_regex":        `"${alicloud_oos_parameter.default.parameter_name}"`,
			"resource_group_id": `"${alicloud_oos_parameter.default.resource_group_id}"`,
			"sort_field":        `"Name"`,
			"sort_order":        `"Ascending"`,
			"tags": `{
				"Created" = "TF"
				"For" = "OosParameter"
		    }`,
			"type": `"String"`,
		}),
		fakeConfig: testAccCheckAlicloudOosParametersDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_oos_parameter.default.parameter_name}_fake"]`,
			"parameter_name":    `"${alicloud_oos_parameter.default.parameter_name}_fake"`,
			"name_regex":        `"${alicloud_oos_parameter.default.parameter_name}_fake"`,
			"resource_group_id": `"fake"`,
			"tags": `{
				"Created" = "OosParameter"
				"For" = "TF"
			}`,
			"type": `"StringList"`,
		}),
	}
	var existAlicloudOosParametersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"parameters.#":                   "1",
			"parameters.0.parameter_name":    fmt.Sprintf("tf-testAccParameter-%d", rand),
			"parameters.0.type":              "String",
			"parameters.0.constraints":       CHECKSET,
			"parameters.0.create_time":       CHECKSET,
			"parameters.0.created_by":        CHECKSET,
			"parameters.0.description":       "oos_parameter_description",
			"parameters.0.parameter_id":      CHECKSET,
			"parameters.0.id":                CHECKSET,
			"parameters.0.parameter_version": CHECKSET,
			"parameters.0.resource_group_id": CHECKSET,
			"parameters.0.share_type":        CHECKSET,
			"parameters.0.tags.%":            "2",
			"parameters.0.tags.Created":      "TF",
			"parameters.0.tags.For":          "OosParameter",
			"parameters.0.updated_by":        CHECKSET,
			"parameters.0.updated_date":      CHECKSET,
			"parameters.0.value":             "tf-testacc-oos_parameter",
		}
	}
	var fakeAlicloudOosParametersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudOosParametersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_oos_parameters.default",
		existMapFunc: existAlicloudOosParametersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudOosParametersDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudOosParametersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameConf, resourceGroupIdConf, sortOrderConf, tagsConf, typeConf, nameRegexConf, allConf)
}

func testAccCheckAlicloudOosParametersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tf-testAccParameter-%d"
}

resource "alicloud_oos_parameter" "default" {
  parameter_name = var.name
  type           = "String"
  value          = "tf-testacc-oos_parameter"
  tags = {
    Created = "TF"
    For     = "OosParameter"
  }
  description    = "oos_parameter_description"
  constraints    = "{\"AllowedValues\":[\"tf-testacc-oos_parameter\"], \"AllowedPattern\": \"tf-testacc-oos_parameter\", \"MinLength\": 1, \"MaxLength\": 100}"
}

data "alicloud_oos_parameters" "default" {
  enable_details = true
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
