package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMessageServiceTopicsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMessageServiceTopicsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_message_service_topic.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMessageServiceTopicsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_message_service_topic.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMessageServiceTopicsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_message_service_topic.default.topic_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMessageServiceTopicsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_message_service_topic.default.topic_name}_fake"`,
		}),
	}
	topicNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMessageServiceTopicsDataSourceName(rand, map[string]string{
			"topic_name": `"${alicloud_message_service_topic.default.topic_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMessageServiceTopicsDataSourceName(rand, map[string]string{
			"topic_name": `"${alicloud_message_service_topic.default.topic_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMessageServiceTopicsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_message_service_topic.default.id}"]`,
			"name_regex": `"${alicloud_message_service_topic.default.topic_name}"`,
			"topic_name": `"${alicloud_message_service_topic.default.topic_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMessageServiceTopicsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_message_service_topic.default.id}_fake"]`,
			"name_regex": `"${alicloud_message_service_topic.default.topic_name}_fake"`,
			"topic_name": `"${alicloud_message_service_topic.default.topic_name}_fake"`,
		}),
	}
	var existAlicloudMessageServiceTopicsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"topics.#":                          "1",
			"topics.0.id":                       CHECKSET,
			"topics.0.topic_name":               CHECKSET,
			"topics.0.message_count":            CHECKSET,
			"topics.0.max_message_size":         "12357",
			"topics.0.message_retention_period": CHECKSET,
			"topics.0.logging_enabled":          "true",
			"topics.0.topic_url":                CHECKSET,
			"topics.0.topic_inner_url":          CHECKSET,
			"topics.0.last_modify_time":         CHECKSET,
			"topics.0.create_time":              CHECKSET,
		}
	}
	var fakeAlicloudMessageServiceTopicsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"topics.#": "0",
		}
	}
	var alicloudMessageServiceTopicsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_message_service_topics.default",
		existMapFunc: existAlicloudMessageServiceTopicsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudMessageServiceTopicsDataSourceNameMapFunc,
	}
	alicloudMessageServiceTopicsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, topicNameConf, allConf)
}

func testAccCheckAlicloudMessageServiceTopicsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccMNSTopic-%d"
	}

	resource "alicloud_message_service_topic" "default" {
  		topic_name       = var.name
  		max_message_size = 12357
  		logging_enabled   = true
	}
	
	data "alicloud_message_service_topics" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
