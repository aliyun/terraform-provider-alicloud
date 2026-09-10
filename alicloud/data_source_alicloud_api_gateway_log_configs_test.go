package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudApiGatewayLogConfigsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApiGatewayLogConfigsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_api_gateway_log_config.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudApiGatewayLogConfigsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_api_gateway_log_config.default.id}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApiGatewayLogConfigsDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_api_gateway_log_config.default.id}"]`,
			"log_type": `"${alicloud_api_gateway_log_config.default.log_type}"`,
		}),
		fakeConfig: testAccCheckAlicloudApiGatewayLogConfigsDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_api_gateway_log_config.default.id}_fake"]`,
			"log_type": `"${alicloud_api_gateway_log_config.default.log_type}"`,
		}),
	}
	var existAlicloudApiGatewayLogConfigsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"configs.#":               "1",
			"configs.0.id":            CHECKSET,
			"configs.0.log_type":      CHECKSET,
			"configs.0.sls_project":   CHECKSET,
			"configs.0.sls_log_store": CHECKSET,
			"configs.0.region_id":     CHECKSET,
		}
	}
	var fakeAlicloudApiGatewayLogConfigsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"configs.#": "0",
		}
	}
	var alicloudApiGatewayLogConfigsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_api_gateway_log_configs.default",
		existMapFunc: existAlicloudApiGatewayLogConfigsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudApiGatewayLogConfigsDataSourceNameMapFunc,
	}
	alicloudApiGatewayLogConfigsCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
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
