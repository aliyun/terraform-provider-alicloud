package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEventBridgeEventStreamingsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.EventBridgeSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeEventStreamingsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_event_bridge_event_streaming.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeEventStreamingsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_event_bridge_event_streaming.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeEventStreamingsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_event_bridge_event_streaming.default.event_streaming_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeEventStreamingsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_event_bridge_event_streaming.default.event_streaming_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeEventStreamingsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_event_bridge_event_streaming.default.id}"]`,
			"status": `"${alicloud_event_bridge_event_streaming.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeEventStreamingsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_event_bridge_event_streaming.default.id}"]`,
			"status": `"PAUSED"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeEventStreamingsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_event_bridge_event_streaming.default.id}"]`,
			"name_regex": `"${alicloud_event_bridge_event_streaming.default.event_streaming_name}"`,
			"status":     `"${alicloud_event_bridge_event_streaming.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeEventStreamingsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_event_bridge_event_streaming.default.id}_fake"]`,
			"name_regex": `"${alicloud_event_bridge_event_streaming.default.event_streaming_name}_fake"`,
			"status":     `"PAUSED"`,
		}),
	}
	var existAlicloudEventBridgeEventStreamingsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"streamings.#":                      "1",
			"streamings.0.description":          fmt.Sprintf("tf-testAccEventStreaming-%d", rand),
			"streamings.0.event_streaming_name": fmt.Sprintf("tf-testAccEventStreaming-%d", rand),
			"streamings.0.filter_pattern":       "{}",
			"streamings.0.source.#":             "1",
			"streamings.0.source.0.source_mns_parameters.#":                  "1",
			"streamings.0.source.0.source_mns_parameters.0.queue_name":       "test",
			"streamings.0.source.0.source_mns_parameters.0.is_base64_decode": "true",
			"streamings.0.sink.#":                                                "1",
			"streamings.0.sink.0.sink_mns_parameters.#":                          "1",
			"streamings.0.sink.0.sink_mns_parameters.0.queue_name.#":             "1",
			"streamings.0.sink.0.sink_mns_parameters.0.queue_name.0.value":       "test",
			"streamings.0.sink.0.sink_mns_parameters.0.queue_name.0.form":        "CONSTANT",
			"streamings.0.sink.0.sink_mns_parameters.0.body.#":                   "1",
			"streamings.0.sink.0.sink_mns_parameters.0.body.0.value":             "$.data",
			"streamings.0.sink.0.sink_mns_parameters.0.body.0.form":              "JSONPATH",
			"streamings.0.sink.0.sink_mns_parameters.0.is_base64_encode.#":       "1",
			"streamings.0.sink.0.sink_mns_parameters.0.is_base64_encode.0.value": "true",
			"streamings.0.sink.0.sink_mns_parameters.0.is_base64_encode.0.form":  "CONSTANT",
			"streamings.0.run_options.#":                                         "1",
			"streamings.0.run_options.0.errors_tolerance":                        "ALL",
			"streamings.0.run_options.0.retry_strategy.#":                        "1",
			"streamings.0.run_options.0.retry_strategy.0.push_retry_strategy":    "BACKOFF_RETRY",
		}
	}
	var fakeAlicloudEventBridgeEventStreamingsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEventBridgeEventStreamingsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_event_bridge_event_streamings.default",
		existMapFunc: existAlicloudEventBridgeEventStreamingsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEventBridgeEventStreamingsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEventBridgeEventStreamingsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudEventBridgeEventStreamingsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccEventStreaming-%d"
}

resource "alicloud_event_bridge_event_streaming" "default" {
	event_streaming_name =  var.name
	description          =  var.name
	source {
		source_mns_parameters {
			queue_name =       "test"
			is_base64_decode = "true"
		}
	}
	sink {
		sink_mns_parameters {
			queue_name {
				value = "test"
				form =  "CONSTANT"
			}
			body {
				value = "$.data"
				form =  "JSONPATH"
			}
			is_base64_encode {
				value = "true"
				form =  "CONSTANT"
			}
		}
	}
	run_options {
		errors_tolerance = "ALL"
		retry_strategy {
			push_retry_strategy = "BACKOFF_RETRY"
		}
	}
	filter_pattern = "{}"
}

data "alicloud_event_bridge_event_streamings" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
