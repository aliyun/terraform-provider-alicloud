package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudAlikafkaConsumerGroupsDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_alikafka_consumer_groups.default"
	name := fmt.Sprintf("tf-testacc-alikafkaconsumer%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAlikafkaConsumerGroupsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_alikafka_instance.default.id}",
			"consumer_id_regex": "${alicloud_alikafka_consumer_group.default.consumer_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_alikafka_instance.default.id}",
			"consumer_id_regex": "fake_tf-testacc*",
		}),
	}

	var existAlikafkaConsumerGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"consumer_ids.#": "1",
			"consumer_ids.0": fmt.Sprintf("tf-testacc-alikafkaconsumer%v", rand),
		}
	}

	var fakeAlikafkaConsumerGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"consumer_ids.#": "0",
		}
	}

	var alikafkaConsumerGroupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAlikafkaConsumerGroupsMapFunc,
		fakeMapFunc:  fakeAlikafkaConsumerGroupsMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.AlikafkaSupportedRegions)
		testAccPreCheckWithNoDefaultVswitch(t)
	}
	alikafkaConsumerGroupsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf)
}

func dataSourceAlikafkaConsumerGroupsConfigDependence(name string) string {
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

		resource "alicloud_alikafka_consumer_group" "default" {
		  instance_id = "${alicloud_alikafka_instance.default.id}"
		  consumer_id = "${var.name}"
		}
		`, name)
}
