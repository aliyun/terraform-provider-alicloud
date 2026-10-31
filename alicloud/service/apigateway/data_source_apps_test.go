package apigateway_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func SkipTestAccAlicloudApigatewayAppsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_api_gateway_apps.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccApp_%d", rand),
		dataSourceApiGatewayAppsConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_api_gateway_app.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_api_gateway_app.default.name}_fake",
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_api_gateway_app.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_api_gateway_app.default.id}_fake"},
		}),
	}
	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "acceptance test",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]interface{}{
				"Created": "TF-fake",
				"For":     "acceptance test",
			},
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_api_gateway_app.default.name}",
			"ids":        []string{"${alicloud_api_gateway_app.default.id}"},
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "acceptance test",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_api_gateway_app.default.name}",
			"ids":        []string{"${alicloud_api_gateway_app.default.id}_fake"},
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "acceptance test",
			},
		}),
	}

	var existApiGatewayAppsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":              "1",
			"names.#":            "1",
			"names.0":            fmt.Sprintf("tf_testAccApp_%d", rand),
			"apps.#":             "1",
			"apps.0.name":        fmt.Sprintf("tf_testAccApp_%d", rand),
			"apps.0.description": "tf_testAcc api gateway description",
			"apps.0.app_code":    acctest.CHECKSET,
		}
	}

	var fakeApiGatewayAppsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"apps.#":  "0",
		}
	}

	var apiGatewayAppsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existApiGatewayAppsMapFunc,
		FakeMapFunc:  fakeApiGatewayAppsMapFunc,
	}

	apiGatewayAppsCheckInfo.DataSourceTestCheck(t, rand, nameRegexConf, idsConf, tagsConf, allConf)
}

func dataSourceApiGatewayAppsConfigDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

variable "description" {
  default = "tf_testAcc api gateway description"
}

resource "alicloud_api_gateway_app" "default" {
  name = "${var.name}"
  description = "${var.description}"
  tags 		= {
		Created = "TF"
		For 	= "acceptance test"
  }
}

`, name)
}
