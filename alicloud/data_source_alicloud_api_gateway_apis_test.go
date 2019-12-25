package alicloud

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudApigatewayApisDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_api_gateway_apis.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccApisDataSource_%d", rand),
		dataSourceApigatewayApisConfigDependence)

	groupAndIdsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{alicloud_api_gateway_api.default.api_id},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_api_gateway_api.default.api_id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": alicloud_api_gateway_api.default.name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_api_gateway_api.default.name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"group_id":   alicloud_api_gateway_group.default.id,
			"ids":        []string{alicloud_api_gateway_api.default.api_id},
			"name_regex": alicloud_api_gateway_api.default.name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"group_id":   "${alicloud_api_gateway_group.default.id}_fake",
			"ids":        []string{alicloud_api_gateway_api.default.api_id},
			"name_regex": alicloud_api_gateway_api.default.name,
		}),
	}

	var existApigatewayApisMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":              "1",
			"ids.0":              CHECKSET,
			"names.#":            "1",
			"names.0":            fmt.Sprintf("tf_testAccApisDataSource_%d", rand),
			"apis.#":             "1",
			"apis.0.name":        fmt.Sprintf("tf_testAccApisDataSource_%d", rand),
			"apis.0.group_name":  fmt.Sprintf("tf_testAccApisDataSource_%d", rand),
			"apis.0.description": "tf_testAcc_api description",
		}
	}

	var fakeApigatewayApisMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"apis.#":  "0",
		}
	}

	var apigatewayApisCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existApigatewayApisMapFunc,
		fakeMapFunc:  fakeApigatewayApisMapFunc,
	}

	apigatewayApisCheckInfo.dataSourceTestCheck(t, rand, groupAndIdsConf, nameRegexConf, allConf)
}

func dataSourceApigatewayApisConfigDependence(name string) string {
	return fmt.Sprintf(`

	variable "name" {
	  default = "%s"
	}

	variable "apigateway_group_description_test" {
	  default = "tf_testAcc_api group description"
	}

	resource "alicloud_api_gateway_group" "default" {
	  name = var.name
	  description = var.apigateway_group_description_test
	}

	resource "alicloud_api_gateway_api" "default" {
	  name = var.name
	  group_id = alicloud_api_gateway_group.default.id
	  description = "tf_testAcc_api description"
	  auth_type = "APP"
	  request_config {
	      protocol = "HTTP"
	      method = "GET"
	      path = "/test/path"
	      mode = "MAPPING"
	    }
	  service_type = "HTTP"
	  http_service_config {
	      address = "http://apigateway-backend.default.com:8080"
	      method = "GET"
	      path = "/web/cloudapi"
	      timeout = 20
	      aone_name = "cloudapi-openapi"
	    }

	  request_parameters {
	      name = "testparam"
	      type = "STRING"
	      required = "OPTIONAL"
	      in = "QUERY"
	      in_service = "QUERY"
	      name_service = "testparams"
	    }
	}
	`, name)
}
