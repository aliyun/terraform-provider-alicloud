package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudApiGatewayModelsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_api_gateway_model.default.id}"]`,
			"group_id": `"${alicloud_api_gateway_model.default.group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_api_gateway_model.default.id}_fake"]`,
			"group_id": `"${alicloud_api_gateway_model.default.group_id}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_api_gateway_model.default.model_name}"`,
			"group_id":   `"${alicloud_api_gateway_model.default.group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_api_gateway_model.default.model_name}_fake"`,
			"group_id":   `"${alicloud_api_gateway_model.default.group_id}_fake"`,
		}),
	}
	groupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"group_id": `"${alicloud_api_gateway_model.default.group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"group_id": `"${alicloud_api_gateway_model.default.group_id}_fake"`,
		}),
	}
	modelNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"group_id":   `"${alicloud_api_gateway_model.default.group_id}"`,
			"model_name": `"${alicloud_api_gateway_model.default.model_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"group_id":   `"${alicloud_api_gateway_model.default.group_id}_fake"`,
			"model_name": `"${alicloud_api_gateway_model.default.model_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_api_gateway_model.default.id}"]`,
			"name_regex": `"${alicloud_api_gateway_model.default.model_name}"`,
			"group_id":   `"${alicloud_api_gateway_model.default.group_id}"`,
			"model_name": `"${alicloud_api_gateway_model.default.model_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudApiGatewayModelsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_api_gateway_model.default.id}_fake"]`,
			"name_regex": `"${alicloud_api_gateway_model.default.model_name}_fake"`,
			"group_id":   `"${alicloud_api_gateway_model.default.group_id}_fake"`,
			"model_name": `"${alicloud_api_gateway_model.default.model_name}_fake"`,
		}),
	}
	var existAlicloudApiGatewayModelsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "1",
			"names.#":                "1",
			"models.#":               "1",
			"models.0.id":            CHECKSET,
			"models.0.group_id":      CHECKSET,
			"models.0.model_name":    CHECKSET,
			"models.0.schema":        "{\"type\":\"object\",\"properties\":{\"id\":{\"format\":\"int64\",\"maximum\":100,\"exclusiveMaximum\":true,\"type\":\"integer\"},\"name\":{\"maxLength\":10,\"type\":\"string\"}}}",
			"models.0.description":   CHECKSET,
			"models.0.model_id":      CHECKSET,
			"models.0.model_ref":     CHECKSET,
			"models.0.modified_time": CHECKSET,
			"models.0.create_time":   CHECKSET,
		}
	}
	var fakeAlicloudApiGatewayModelsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"models.#": "0",
		}
	}
	var alicloudApiGatewayModelsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_api_gateway_models.default",
		existMapFunc: existAlicloudApiGatewayModelsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudApiGatewayModelsDataSourceNameMapFunc,
	}
	alicloudApiGatewayModelsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, groupIdConf, modelNameConf, allConf)
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
