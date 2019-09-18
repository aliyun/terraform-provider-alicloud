package alicloud

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudAlikafkaConsumerGroupsDataSource(t *testing.T) {

	testAccPreCheckWithAlikafkaInstanceSetting(t)

	rand := acctest.RandInt()
	resourceId := "data.alicloud_alikafka_consumer_groups.default"
	name := fmt.Sprintf("tf-testacc-alikafkaconsumer%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAlikafkaConsumerGroupsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":       os.Getenv("ALICLOUD_INSTANCE_ID"),
			"consumer_id_regex": "${alicloud_alikafka_consumer_group.default.consumer_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":       os.Getenv("ALICLOUD_INSTANCE_ID"),
			"consumer_id_regex": "fake_tf-testacc*",
		}),
	}

	var existAlikafkaConsumerGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"consumer_ids.#":                         "1",
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
		variable "consumer_id" {
		 default = "%v"
		}
		
		resource "alicloud_alikafka_consumer_group" "default" {
		  instance_id = "%v"
		  consumer_id = "${var.consumer_id}"
		}
		`, name, os.Getenv("ALICLOUD_INSTANCE_ID"))
}
