package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDatabaseGatewayGatewaysDataSource(t *testing.T) {
	rand := acctest.RandInt()

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_database_gateway_gateway.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_database_gateway_gateway.default.id}_fakeid"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_database_gateway_gateway.default.gateway_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_database_gateway_gateway.default.gateway_name}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_database_gateway_gateway.default.id}"]`,
			"status": `"NEW"`,
		}),
		fakeConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_database_gateway_gateway.default.id}"]`,
			"status": `"STOPPED"`,
		}),
	}

	searchKeyConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"search_key": `"${alicloud_database_gateway_gateway.default.gateway_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"search_key": `"${alicloud_database_gateway_gateway.default.gateway_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_database_gateway_gateway.default.id}"]`,
			"name_regex": `"${alicloud_database_gateway_gateway.default.gateway_name}"`,
			"search_key": `"${alicloud_database_gateway_gateway.default.gateway_name}"`,
			"status":     `"NEW"`,
		}),
		fakeConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_database_gateway_gateway.default.id}_fake"]`,
			"name_regex": `"${alicloud_database_gateway_gateway.default.gateway_name}_fake"`,
			"search_key": `"${alicloud_database_gateway_gateway.default.gateway_name}_fake"`,
			"status":     `"STOPPED"`,
		}),
	}

	var existDataAlicloudDatabaseGatewayGatewaysSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"gateways.#":              "1",
			"gateways.0.gateway_name": fmt.Sprintf("tf-testAccDatabaseGatewayGateway%d", rand),
		}
	}
	var fakeDataAlicloudDatabaseGatewayGatewaysSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"gateways.#": "0",
		}
	}
	var alicloudDatabaseGatewayGatewayCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_database_gateway_gateways.default",
		existMapFunc: existDataAlicloudDatabaseGatewayGatewaysSourceNameMapFunc,
		fakeMapFunc:  fakeDataAlicloudDatabaseGatewayGatewaysSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.DatabaseGatewaySupportRegions)
	}
	alicloudDatabaseGatewayGatewayCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, searchKeyConf, allConf)
}
func testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccDatabaseGatewayGateway%d"
}
resource "alicloud_database_gateway_gateway" "default" {
	gateway_name = var.name
}
data "alicloud_database_gateway_gateways" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
