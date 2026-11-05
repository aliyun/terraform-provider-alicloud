package mnsopen_test

import (
	"fmt"
	"strings"
	"testing"

	tfacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudMessageServiceQueuesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMessageServiceQueuesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_message_service_queue.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudMessageServiceQueuesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_message_service_queue.default.id}_fake"]`,
		}),
	}
	nameRegexConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMessageServiceQueuesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_message_service_queue.default.queue_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudMessageServiceQueuesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_message_service_queue.default.queue_name}_fake"`,
		}),
	}
	queueNameConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMessageServiceQueuesDataSourceName(rand, map[string]string{
			"queue_name": `"${alicloud_message_service_queue.default.queue_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudMessageServiceQueuesDataSourceName(rand, map[string]string{
			"queue_name": `"${alicloud_message_service_queue.default.queue_name}_fake"`,
		}),
	}
	allConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMessageServiceQueuesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_message_service_queue.default.id}"]`,
			"name_regex": `"${alicloud_message_service_queue.default.queue_name}"`,
			"queue_name": `"${alicloud_message_service_queue.default.queue_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudMessageServiceQueuesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_message_service_queue.default.id}_fake"]`,
			"name_regex": `"${alicloud_message_service_queue.default.queue_name}_fake"`,
			"queue_name": `"${alicloud_message_service_queue.default.queue_name}_fake"`,
		}),
	}
	var existAlicloudMessageServiceQueuesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"queues.#":                          "1",
			"queues.0.id":                       tfacctest.CHECKSET,
			"queues.0.queue_name":               tfacctest.CHECKSET,
			"queues.0.delay_seconds":            "60478",
			"queues.0.maximum_message_size":     "12357",
			"queues.0.message_retention_period": "256000",
			"queues.0.visibility_timeout":       "30",
			"queues.0.polling_wait_seconds":     "3",
			"queues.0.logging_enabled":          "true",
			"queues.0.active_messages":          tfacctest.CHECKSET,
			"queues.0.inactive_messages":        tfacctest.CHECKSET,
			"queues.0.delay_messages":           tfacctest.CHECKSET,
			"queues.0.queue_url":                tfacctest.CHECKSET,
			"queues.0.queue_internal_url":       tfacctest.CHECKSET,
			"queues.0.last_modify_time":         tfacctest.CHECKSET,
			"queues.0.create_time":              tfacctest.CHECKSET,
		}
	}
	var fakeAlicloudMessageServiceQueuesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"queues.#": "0",
		}
	}
	var alicloudMessageServiceQueuesCheckInfo = tfacctest.DataSourceAttr{
		ResourceId:   "data.alicloud_message_service_queues.default",
		ExistMapFunc: existAlicloudMessageServiceQueuesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudMessageServiceQueuesDataSourceNameMapFunc,
	}
	alicloudMessageServiceQueuesCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, queueNameConf, allConf)
}

func testAccCheckAlicloudMessageServiceQueuesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccMNSQueue-%d"
	}

	resource "alicloud_message_service_queue" "default" {
		queue_name               = var.name
		delay_seconds            = 60478
		maximum_message_size     = 12357
		message_retention_period = 256000
		visibility_timeout       = 30
		polling_wait_seconds     = 3
		logging_enabled          = true
	}
	
	data "alicloud_message_service_queues" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
