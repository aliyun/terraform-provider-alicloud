package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudApiGatewayDatasetsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApiGatewayDatasetsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_api_gateway_dataset.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudApiGatewayDatasetsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_api_gateway_dataset.default.id}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApiGatewayDatasetsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_api_gateway_dataset.default.id}"]`,
			"name_regex": `"${alicloud_api_gateway_dataset.default.dataset_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudApiGatewayDatasetsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_api_gateway_dataset.default.id}_fake"]`,
			"name_regex": `"${alicloud_api_gateway_dataset.default.dataset_name}_fake"`,
		}),
	}
	var existAlicloudApiGatewayDatasetsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"datasets.#":              "1",
			"datasets.0.id":           CHECKSET,
			"datasets.0.dataset_id":   CHECKSET,
			"datasets.0.dataset_name": CHECKSET,
			"datasets.0.dataset_type": CHECKSET,
		}
	}
	var fakeAlicloudApiGatewayDatasetsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"datasets.#": "0",
		}
	}
	var alicloudApiGatewayDatasetsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_api_gateway_datasets.default",
		existMapFunc: existAlicloudApiGatewayDatasetsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudApiGatewayDatasetsDataSourceNameMapFunc,
	}
	alicloudApiGatewayDatasetsCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

func testAccCheckAlicloudApiGatewayDatasetsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  	default = "tf_testacc_ds_%d"
}

resource "alicloud_api_gateway_dataset" "default" {
  	dataset_name = var.name
  	dataset_type = "JWT_BLOCKING"
  	description  = "tf test acc dataset for datasource"
}

data "alicloud_api_gateway_datasets" "default" {
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
