package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEventBridgeEventBusesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeEventBusesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_event_bridge_event_bus.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeEventBusesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_event_bridge_event_bus.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeEventBusesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_event_bridge_event_bus.default.event_bus_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeEventBusesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_event_bridge_event_bus.default.event_bus_name}_fake"`,
		}),
	}
	namePrefixConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeEventBusesDataSourceName(rand, map[string]string{
			"name_prefix": `"${alicloud_event_bridge_event_bus.default.event_bus_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeEventBusesDataSourceName(rand, map[string]string{
			"name_prefix": `"${alicloud_event_bridge_event_bus.default.event_bus_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeEventBusesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_event_bridge_event_bus.default.id}"]`,
			"name_regex":  `"${alicloud_event_bridge_event_bus.default.event_bus_name}"`,
			"name_prefix": `"tf-testAcc"`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeEventBusesDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_event_bridge_event_bus.default.id}_fake"]`,
			"name_regex":  `"${alicloud_event_bridge_event_bus.default.event_bus_name}_fake"`,
			"name_prefix": `"tf-testAcc_fake"`,
		}),
	}
	var existAlicloudEventBridgeEventBusesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "1",
			"names.#":                "1",
			"buses.#":                "1",
			"buses.0.description":    fmt.Sprintf("tf-testAccEventBus-%d", rand),
			"buses.0.event_bus_name": fmt.Sprintf("tf-testAccEventBus-%d", rand),
		}
	}
	var fakeAlicloudEventBridgeEventBusesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEventBridgeEventBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_event_bridge_event_buses.default",
		existMapFunc: existAlicloudEventBridgeEventBusesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEventBridgeEventBusesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.EventBridgeSupportRegions)
	}
	alicloudEventBridgeEventBusesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, namePrefixConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudEventBridgeEventBusesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccEventBus-%d"
}

resource "alicloud_event_bridge_event_bus" "default" {
	description = var.name
	event_bus_name = var.name
}

data "alicloud_event_bridge_event_buses" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
