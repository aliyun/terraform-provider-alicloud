package apigateway_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudApiGatewayLogConfigsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudApiGatewayLogConfigsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_api_gateway_log_config.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudApiGatewayLogConfigsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_api_gateway_log_config.default.id}_fake"]`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudApiGatewayLogConfigsDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_api_gateway_log_config.default.id}"]`,
			"log_type": `"${alicloud_api_gateway_log_config.default.log_type}"`,
		}),
		FakeConfig: testAccCheckAlicloudApiGatewayLogConfigsDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_api_gateway_log_config.default.id}_fake"]`,
			"log_type": `"${alicloud_api_gateway_log_config.default.log_type}"`,
		}),
	}
	var existAlicloudApiGatewayLogConfigsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"configs.#":               "1",
			"configs.0.id":            acctest.CHECKSET,
			"configs.0.log_type":      acctest.CHECKSET,
			"configs.0.sls_project":   acctest.CHECKSET,
			"configs.0.sls_log_store": acctest.CHECKSET,
			"configs.0.region_id":     acctest.CHECKSET,
		}
	}
	var fakeAlicloudApiGatewayLogConfigsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"configs.#": "0",
		}
	}
	var alicloudApiGatewayLogConfigsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_api_gateway_log_configs.default",
		ExistMapFunc: existAlicloudApiGatewayLogConfigsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudApiGatewayLogConfigsDataSourceNameMapFunc,
	}
	alicloudApiGatewayLogConfigsCheckInfo.DataSourceTestCheck(t, rand, idsConf, allConf)
}

func testAccCheckAlicloudApiGatewayLogConfigsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf_testacc_%d"
	}

	resource "alicloud_api_gateway_log_config" "default" {
  		sls_project   = var.name
  		sls_log_store = var.name
  		log_type      = "PROVIDER"
	}

	data "alicloud_api_gateway_log_configs" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
