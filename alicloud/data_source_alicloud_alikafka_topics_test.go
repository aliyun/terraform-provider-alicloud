package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlikafkaTopicsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.AlikafkaSupportedRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlikafkaTopicsDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_alikafka_instance.default.id}"`,
			"ids":         `["${alicloud_alikafka_topic.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlikafkaTopicsDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_alikafka_instance.default.id}"`,
			"ids":         `["${alicloud_alikafka_topic.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlikafkaTopicsDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_alikafka_instance.default.id}"`,
			"name_regex":  `"${alicloud_alikafka_topic.default.topic}"`,
		}),
		fakeConfig: testAccCheckAlikafkaTopicsDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_alikafka_instance.default.id}"`,
			"name_regex":  `"${alicloud_alikafka_topic.default.topic}_fake"`,
		}),
	}

	topicConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlikafkaTopicsDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_alikafka_instance.default.id}"`,
			"topic":       `"${alicloud_alikafka_topic.default.topic}"`,
		}),
		fakeConfig: testAccCheckAlikafkaTopicsDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_alikafka_instance.default.id}"`,
			"topic":       `"${alicloud_alikafka_topic.default.topic}_fake"`,
		}),
	}

	pagingConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlikafkaTopicsDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_alikafka_instance.default.id}"`,
			"topic":       `"${alicloud_alikafka_topic.default.topic}"`,
			"page_number": `1`,
		}),
	}

	var existAlikafkaTopicsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  CHECKSET,
			"names.#":                CHECKSET,
			"total_count":            CHECKSET,
			"topics.#":               CHECKSET,
			"topics.0.id":            CHECKSET,
			"topics.0.topic":         fmt.Sprintf("tf-testacc-alikafkatopic%v", rand),
			"topics.0.local_topic":   "false",
			"topics.0.compact_topic": "false",
			"topics.0.partition_num": "12",
			"topics.0.remark":        "remark",
			"topics.0.status":        "0",
			"topics.0.status_name":   CHECKSET,
			"topics.0.instance_id":   CHECKSET,
			"topics.0.tags.%":        "2",
			"topics.0.tags.Created":  "TF",
			"topics.0.tags.For":      "Test",
		}
	}
	var fakeAlikafkaTopicsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"topics.#": "0",
			"names.#":  "0",
		}
	}
	var AlikafkaTopicsBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alikafka_topics.default",
		existMapFunc: existAlikafkaTopicsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlikafkaTopicsDataSourceNameMapFunc,
	}

	AlikafkaTopicsBusesCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, topicConf, pagingConf)
}
func testAccCheckAlikafkaTopicsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tf-testacc-alikafkatopic%d"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_alikafka_instance" "default" {
  name = "${var.name}"
  topic_quota = "50"
  disk_type = "1"
  disk_size = "500"
  deploy_type = "5"
  io_max = "20"
  vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
  security_group = alicloud_security_group.default.id
}

resource "alicloud_alikafka_topic" "default" {
  instance_id = "${alicloud_alikafka_instance.default.id}"
  topic = "${var.name}"
  local_topic = "false"
  compact_topic = "false"
  partition_num = "12"
  remark = "remark"
  tags = {
       Created = "TF"
      For =     "Test"
  }
}

data "alicloud_alikafka_topics" "default" {    
   %s
}

`, rand, strings.Join(pairs, " \n "))
	return config
}
