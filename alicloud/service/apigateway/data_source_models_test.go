package apigateway_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudApiGatewayModelsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_api_gateway_model.default.id}"]`,
			"group_id": `"${alicloud_api_gateway_model.default.group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_api_gateway_model.default.id}_fake"]`,
			"group_id": `"${alicloud_api_gateway_model.default.group_id}"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_api_gateway_model.default.model_name}"`,
			"group_id":   `"${alicloud_api_gateway_model.default.group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"name_regex": `"^${alicloud_api_gateway_model.default.model_name}_fake$"`,
			"group_id":   `"${alicloud_api_gateway_model.default.group_id}"`,
		}),
	}
	groupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"group_id": `"${alicloud_api_gateway_model.default.group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"group_id": `"${alicloud_api_gateway_group.empty.id}"`,
		}) + `
resource "alicloud_api_gateway_group" "empty" {
  name        = "${var.name}-empty"
  description = "${var.name}-empty"
}
`,
	}
	modelNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"group_id":   `"${alicloud_api_gateway_model.default.group_id}"`,
			"model_name": `"${alicloud_api_gateway_model.default.model_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"group_id":   `"${alicloud_api_gateway_model.default.group_id}"`,
			"model_name": `"${alicloud_api_gateway_model.default.model_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_api_gateway_model.default.id}"]`,
			"name_regex": `"${alicloud_api_gateway_model.default.model_name}"`,
			"group_id":   `"${alicloud_api_gateway_model.default.group_id}"`,
			"model_name": `"${alicloud_api_gateway_model.default.model_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_api_gateway_model.default.id}_fake"]`,
			"name_regex": `"^${alicloud_api_gateway_model.default.model_name}_fake$"`,
			"group_id":   `"${alicloud_api_gateway_model.default.group_id}"`,
			"model_name": `"${alicloud_api_gateway_model.default.model_name}_fake"`,
		}),
	}
	var existAlicloudApiGatewayModelsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "1",
			"names.#":                "1",
			"models.#":               "1",
			"models.0.id":            acctest.CHECKSET,
			"models.0.group_id":      acctest.CHECKSET,
			"models.0.model_name":    acctest.CHECKSET,
			"models.0.schema":        "{\"type\":\"object\",\"properties\":{\"id\":{\"format\":\"int64\",\"maximum\":100,\"exclusiveMaximum\":true,\"type\":\"integer\"},\"name\":{\"maxLength\":10,\"type\":\"string\"}}}",
			"models.0.description":   acctest.CHECKSET,
			"models.0.model_id":      acctest.CHECKSET,
			"models.0.model_ref":     acctest.CHECKSET,
			"models.0.modified_time": acctest.CHECKSET,
			"models.0.create_time":   acctest.CHECKSET,
		}
	}
	var fakeAlicloudApiGatewayModelsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"models.#": "0",
		}
	}
	var alicloudApiGatewayModelsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_api_gateway_models.default",
		ExistMapFunc: existAlicloudApiGatewayModelsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudApiGatewayModelsDataSourceNameMapFunc,
	}
	alicloudApiGatewayModelsCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, groupIdConf, modelNameConf, allConf)
}

func testAccCheckAlicloudApiGatewayModelsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testacc-%d"
	}

	resource "alicloud_api_gateway_group" "default" {
  		name        = var.name
  		description = var.name
	}

	resource "alicloud_api_gateway_model" "default" {
  		group_id    = alicloud_api_gateway_group.default.id
  		model_name  = var.name
  		schema      = "{\"type\":\"object\",\"properties\":{\"id\":{\"format\":\"int64\",\"maximum\":100,\"exclusiveMaximum\":true,\"type\":\"integer\"},\"name\":{\"maxLength\":10,\"type\":\"string\"}}}"
  		description = var.name
	}

	data "alicloud_api_gateway_models" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
