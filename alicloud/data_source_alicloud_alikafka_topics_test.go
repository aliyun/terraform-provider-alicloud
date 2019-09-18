package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudAlikafkaTopicsDataSource(t *testing.T) {

	testAccPreCheckWithAlikafkaInstanceSetting(t)

	rand := acctest.RandInt()
	resourceId := "data.alicloud_alikafka_topics.default"
	name := fmt.Sprintf("tf-testacc-alikafkatopic%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAlikafkaTopicsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": os.Getenv("ALICLOUD_INSTANCE_ID"),
			"name_regex":  "${alicloud_alikafka_topic.default.topic}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": os.Getenv("ALICLOUD_INSTANCE_ID"),
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
		variable "topic" {
		 default = "%v"
		}
		
		resource "alicloud_alikafka_topic" "default" {
		  instance_id = "%v"
		  topic = "${var.topic}"
		  local_topic = "false"
		  compact_topic = "false"
		  partition_num = "12"
		  remark = "alicloud_alikafka_topic_remark"
		}
		`, name, os.Getenv("ALICLOUD_INSTANCE_ID"))
}
