package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMessageServiceQueuesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMessageServiceQueuesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_message_service_queue.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMessageServiceQueuesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_message_service_queue.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMessageServiceQueuesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_message_service_queue.default.queue_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMessageServiceQueuesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_message_service_queue.default.queue_name}_fake"`,
		}),
	}
	queueNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMessageServiceQueuesDataSourceName(rand, map[string]string{
			"queue_name": `"${alicloud_message_service_queue.default.queue_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMessageServiceQueuesDataSourceName(rand, map[string]string{
			"queue_name": `"${alicloud_message_service_queue.default.queue_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMessageServiceQueuesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_message_service_queue.default.id}"]`,
			"name_regex": `"${alicloud_message_service_queue.default.queue_name}"`,
			"queue_name": `"${alicloud_message_service_queue.default.queue_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMessageServiceQueuesDataSourceName(rand, map[string]string{
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
			"queues.0.id":                       CHECKSET,
			"queues.0.queue_name":               CHECKSET,
			"queues.0.delay_seconds":            "60478",
			"queues.0.maximum_message_size":     "12357",
			"queues.0.message_retention_period": "256000",
			"queues.0.visibility_timeout":       "30",
			"queues.0.polling_wait_seconds":     "3",
			"queues.0.logging_enabled":          "true",
			"queues.0.active_messages":          CHECKSET,
			"queues.0.inactive_messages":        CHECKSET,
			"queues.0.delay_messages":           CHECKSET,
			"queues.0.queue_url":                CHECKSET,
			"queues.0.queue_internal_url":       CHECKSET,
			"queues.0.last_modify_time":         CHECKSET,
			"queues.0.create_time":              CHECKSET,
		}
	}
	var fakeAlicloudMessageServiceQueuesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"queues.#": "0",
		}
	}
	var alicloudMessageServiceQueuesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_message_service_queues.default",
		existMapFunc: existAlicloudMessageServiceQueuesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudMessageServiceQueuesDataSourceNameMapFunc,
	}
	alicloudMessageServiceQueuesCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, queueNameConf, allConf)
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
