package databasegateway_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudDatabaseGatewayGatewaysDataSource(t *testing.T) {
	rand := tfacctest.RandInt()

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_database_gateway_gateway.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_database_gateway_gateway.default.id}_fakeid"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_database_gateway_gateway.default.gateway_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_database_gateway_gateway.default.gateway_name}_fake"`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_database_gateway_gateway.default.id}"]`,
			"status": `"NEW"`,
		}),
		FakeConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_database_gateway_gateway.default.id}"]`,
			"status": `"STOPPED"`,
		}),
	}

	searchKeyConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"search_key": `"${alicloud_database_gateway_gateway.default.gateway_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"search_key": `"${alicloud_database_gateway_gateway.default.gateway_name}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_database_gateway_gateway.default.id}"]`,
			"name_regex": `"${alicloud_database_gateway_gateway.default.gateway_name}"`,
			"search_key": `"${alicloud_database_gateway_gateway.default.gateway_name}"`,
			"status":     `"NEW"`,
		}),
		FakeConfig: testAccCheckAlicloudDatabaseGatewayGatewayDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_database_gateway_gateway.default.id}_fake"]`,
			"name_regex": `"${alicloud_database_gateway_gateway.default.gateway_name}_fake"`,
			"search_key": `"${alicloud_database_gateway_gateway.default.gateway_name}_fake"`,
			"status":     `"STOPPED"`,
		}),
	}

	existDataAlicloudDatabaseGatewayGatewaysSourceNameMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"gateways.#":              "1",
			"gateways.0.gateway_name": fmt.Sprintf("tf-testAccDatabaseGatewayGateway%d", rand),
		}
	}
	fakeDataAlicloudDatabaseGatewayGatewaysSourceNameMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"gateways.#": "0",
		}
	}
	alicloudDatabaseGatewayGatewayCheckInfo := acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_database_gateway_gateways.default",
		ExistMapFunc: existDataAlicloudDatabaseGatewayGatewaysSourceNameMapFunc,
		FakeMapFunc:  fakeDataAlicloudDatabaseGatewayGatewaysSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.DatabaseGatewaySupportRegions)
	}
	alicloudDatabaseGatewayGatewayCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, searchKeyConf, allConf)
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
