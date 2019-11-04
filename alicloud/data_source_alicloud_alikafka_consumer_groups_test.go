package alicloud

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

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
	}
	alikafkaConsumerGroupsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf)
}

func dataSourceAlikafkaConsumerGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}

		data "alicloud_zones" "default" {
			available_resource_creation= "VSwitch"
		}
		resource "alicloud_vpc" "default" {
		  cidr_block = "172.16.0.0/12"
		  name       = "${var.name}"
		}
		
		resource "alicloud_vswitch" "default" {
		  vpc_id = "${alicloud_vpc.default.id}"
		  cidr_block = "172.16.0.0/24"
		  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		  name       = "${var.name}"
		}

		resource "alicloud_alikafka_instance" "default" {
          name = "tf-testacc-alikafkainstance"
		  topic_quota = "50"
		  disk_type = "1"
		  disk_size = "500"
		  deploy_type = "5"
		  io_max = "20"
          vswitch_id = "${alicloud_vswitch.default.id}"
		}

		resource "alicloud_alikafka_consumer_group" "default" {
		  instance_id = "${alicloud_alikafka_instance.default.id}"
		  consumer_id = "${var.name}"
		}
		`, name)
}
