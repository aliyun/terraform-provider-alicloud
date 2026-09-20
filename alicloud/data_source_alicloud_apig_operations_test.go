// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudApigOperationsDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigOperationsSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_operation.default.id}"]`,
			"http_api_id": `"${alicloud_apig_http_api.operations_httpapi.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudApigOperationsSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_operation.default.id}_fake"]`,
			"http_api_id": `"${alicloud_apig_http_api.operations_httpapi.id}"`,
		}),
	}

	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigOperationsSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_operation.default.id}"]`,
			"http_api_id": `"${alicloud_apig_http_api.operations_httpapi.id}"`,
			"name":        `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudApigOperationsSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_operation.default.id}_fake"]`,
			"http_api_id": `"${alicloud_apig_http_api.operations_httpapi.id}"`,
			"name":        `"${var.name}_fake"`,
		}),
	}

	methodConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigOperationsSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_operation.default.id}"]`,
			"http_api_id": `"${alicloud_apig_http_api.operations_httpapi.id}"`,
			"method":      `"GET"`,
		}),
		fakeConfig: testAccCheckAlicloudApigOperationsSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_operation.default.id}_fake"]`,
			"http_api_id": `"${alicloud_apig_http_api.operations_httpapi.id}"`,
			"method":      `"DELETE"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApigOperationsSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_operation.default.id}"]`,
			"http_api_id": `"${alicloud_apig_http_api.operations_httpapi.id}"`,
			"name":        `"${var.name}"`,
			"method":      `"GET"`,
		}),
		fakeConfig: testAccCheckAlicloudApigOperationsSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_apig_operation.default.id}_fake"]`,
			"http_api_id": `"${alicloud_apig_http_api.operations_httpapi.id}"`,
			"name":        `"${var.name}_fake"`,
			"method":      `"DELETE"`,
		}),
	}

	ApigOperationsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameConf, methodConf, allConf)
}

var existApigOperationsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"operations.#":                "1",
		"operations.0.id":             CHECKSET,
		"operations.0.http_api_id":    CHECKSET,
		"operations.0.operation_id":   CHECKSET,
		"operations.0.operation_name": CHECKSET,
		"operations.0.path":           CHECKSET,
		"operations.0.method":         CHECKSET,
		"operations.0.description":    CHECKSET,
		"operations.0.create_time":    CHECKSET,
		"operations.0.mock.#":         CHECKSET,
	}
}

var fakeApigOperationsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"operations.#": "0",
	}
}

var ApigOperationsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_apig_operations.default",
	existMapFunc: existApigOperationsMapFunc,
	fakeMapFunc:  fakeApigOperationsMapFunc,
}

func testAccCheckAlicloudApigOperationsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testaccapigop%d"
}

resource "alicloud_apig_http_api" "operations_httpapi" {
  http_api_name = "${var.name}api"
  protocols     = ["HTTP"]
  type          = "Rest"
  description   = "operations datasource test httpapi"
  base_path     = "/${var.name}"
}

resource "alicloud_apig_operation" "default" {
  http_api_id    = alicloud_apig_http_api.operations_httpapi.id
  operation_name = "${var.name}"
  path           = "/ds-op"
  method         = "GET"
  description    = "datasource test operation"
  mock {
    enable           = true
    response_code    = 200
    response_content = "hello"
  }
}

data "alicloud_apig_operations" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
