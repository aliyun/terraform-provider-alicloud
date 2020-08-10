package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudAlikafkaTopicsDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_alikafka_topics.default"
	name := fmt.Sprintf("tf-testacc-alikafkatopic%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAlikafkaTopicsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_alikafka_instance.default.id}",
			"name_regex":  "${alicloud_alikafka_topic.default.topic}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_alikafka_instance.default.id}",
			"name_regex":  "fake_tf-testacc*",
		}),
	}

	var existAlikafkaTopicsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                "1",
			"topics.#":               "1",
			"topics.0.topic":         fmt.Sprintf("tf-testacc-alikafkatopic%v", rand),
			"topics.0.local_topic":   "false",
			"topics.0.compact_topic": "false",
			"topics.0.partition_num": "12",
			"topics.0.remark":        "alicloud_alikafka_topic_remark",
			"topics.0.status":        "0",
		}
	}

	var fakeAlikafkaTopicsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"topics.#": "0",
			"names.#":  "0",
		}
	}

	var alikafkaTopicsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAlikafkaTopicsMapFunc,
		fakeMapFunc:  fakeAlikafkaTopicsMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.AlikafkaSupportedRegions)
	}
	alikafkaTopicsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf)
}

func dataSourceAlikafkaTopicsConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}


        data "alicloud_vswitches" "default" {
		  is_default = "true"
		}

		resource "alicloud_alikafka_instance" "default" {
          name = "${var.name}"
		  topic_quota = "50"
		  disk_type = "1"
		  disk_size = "500"
		  deploy_type = "5"
		  io_max = "20"
          vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
		}
		
		resource "alicloud_alikafka_topic" "default" {
		  instance_id = "${alicloud_alikafka_instance.default.id}"
		  topic = "${var.name}"
		  local_topic = "false"
		  compact_topic = "false"
		  partition_num = "12"
		  remark = "alicloud_alikafka_topic_remark"
		}
		`, name)
}
