package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudApiGatewayTrafficControlsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApiGatewayTrafficControlsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_api_gateway_traffic_control.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudApiGatewayTrafficControlsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_api_gateway_traffic_control.default.id}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudApiGatewayTrafficControlsDataSourceName(rand, map[string]string{
			"ids":                `["${alicloud_api_gateway_traffic_control.default.id}"]`,
			"traffic_control_id": `"${alicloud_api_gateway_traffic_control.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudApiGatewayTrafficControlsDataSourceName(rand, map[string]string{
			"ids":                `["${alicloud_api_gateway_traffic_control.default.id}_fake"]`,
			"traffic_control_id": `"${alicloud_api_gateway_traffic_control.default.id}_fake"`,
		}),
	}
	var existAlicloudApiGatewayTrafficControlsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"controls.#":                      "1",
			"controls.0.id":                   CHECKSET,
			"controls.0.traffic_control_id":   CHECKSET,
			"controls.0.traffic_control_name": CHECKSET,
			"controls.0.traffic_control_unit": CHECKSET,
			"controls.0.api_default":          CHECKSET,
			"controls.0.region_id":            CHECKSET,
		}
	}
	var fakeAlicloudApiGatewayTrafficControlsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"controls.#": "0",
		}
	}
	var alicloudApiGatewayTrafficControlsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_api_gateway_traffic_controls.default",
		existMapFunc: existAlicloudApiGatewayTrafficControlsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudApiGatewayTrafficControlsDataSourceNameMapFunc,
	}
	alicloudApiGatewayTrafficControlsCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

func testAccCheckAlicloudApiGatewayTrafficControlsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
		default = "tf_testacc_trafficcontrol_%d"
	}

	resource "alicloud_api_gateway_traffic_control" "default" {
		traffic_control_name = var.name
		traffic_control_unit = "MINUTE"
		api_default          = 100
	}

	data "alicloud_api_gateway_traffic_controls" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
