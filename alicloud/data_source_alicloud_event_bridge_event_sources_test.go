package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEventBridgeEventSourcesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeEventSourcesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_event_bridge_event_source.default.event_source_name}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeEventSourcesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_event_bridge_event_source.default.event_source_name}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeEventSourcesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_event_bridge_event_source.default.event_source_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeEventSourcesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_event_bridge_event_source.default.event_source_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeEventSourcesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_event_bridge_event_source.default.event_source_name}"]`,
			"name_regex": `"${alicloud_event_bridge_event_source.default.event_source_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeEventSourcesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_event_bridge_event_source.default.event_source_name}_fake"]`,
			"name_regex": `"${alicloud_event_bridge_event_source.default.event_source_name}_fake"`,
		}),
	}
	var existAlicloudEventBridgeEventSourcesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"names.#":                            "1",
			"sources.#":                          "1",
			"sources.0.description":              fmt.Sprintf("tf-testAccEventSources-%d", rand),
			"sources.0.event_source_name":        fmt.Sprintf("tf-testAccEventSources-%d", rand),
			"sources.0.id":                       fmt.Sprintf("tf-testAccEventSources-%d", rand),
			"sources.0.external_source_config.%": "1",
			"sources.0.external_source_config.QueueName": fmt.Sprintf("tf-testAccEventSources-%d", rand),
			"sources.0.external_source_type":             "MNS",
			"sources.0.linked_external_source":           `true`,
			"sources.0.type":                             CHECKSET,
		}
	}
	var fakeAlicloudEventBridgeEventSourcesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"sources.#": "0",
		}
	}
	var alicloudEventBridgeEventSourcesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_event_bridge_event_sources.default",
		existMapFunc: existAlicloudEventBridgeEventSourcesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEventBridgeEventSourcesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.EventBridgeSupportRegions)
	}
	alicloudEventBridgeEventSourcesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudEventBridgeEventSourcesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccEventSources-%d"
}

resource "alicloud_event_bridge_event_bus" "default" {
	description = var.name
	event_bus_name = var.name
}

resource "alicloud_mns_queue" "default" {
  name = var.name
}

resource "alicloud_event_bridge_event_source" "default" {
  event_bus_name = alicloud_event_bridge_event_bus.default.id
  event_source_name = var.name
  description = var.name
  linked_external_source = true
  external_source_type = "MNS"
  external_source_config = {
    QueueName = alicloud_mns_queue.default.name
  }
}

data "alicloud_event_bridge_event_sources" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
