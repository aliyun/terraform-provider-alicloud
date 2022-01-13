package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAAlikafkaConsumerGroupsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.AlikafkaSupportedRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAAlikafkaConsumerGroupsDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_alikafka_instance.default.id}"`,
			"ids":         `["${alicloud_alikafka_consumer_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAAlikafkaConsumerGroupsDataSourceName(rand, map[string]string{
			"instance_id": `"${alicloud_alikafka_instance.default.id}"`,
			"ids":         `["${alicloud_alikafka_consumer_group.default.id}_fake"]`,
		}),
	}

	consumerIdRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAAlikafkaConsumerGroupsDataSourceName(rand, map[string]string{
			"instance_id":       `"${alicloud_alikafka_instance.default.id}"`,
			"consumer_id_regex": `"${alicloud_alikafka_consumer_group.default.consumer_id}"`,
		}),
		fakeConfig: testAccCheckAAlikafkaConsumerGroupsDataSourceName(rand, map[string]string{
			"instance_id":       `"${alicloud_alikafka_instance.default.id}"`,
			"consumer_id_regex": `"${alicloud_alikafka_consumer_group.default.consumer_id}_fake"`,
		}),
	}

	var existAAlikafkaConsumerGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 "1",
			"names.#":               "1",
			"groups.#":              "1",
			"groups.0.id":           CHECKSET,
			"groups.0.consumer_id":  fmt.Sprintf("tf-testacc-alikafkaconsumer%v", rand),
			"groups.0.remark":       "",
			"groups.0.instance_id":  CHECKSET,
			"groups.0.tags.%":       "2",
			"groups.0.tags.Created": "TF",
			"groups.0.tags.For":     "Test",
		}
	}
	var fakeAAlikafkaConsumerGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"groups.#": "0",
			"names.#":  "0",
		}
	}
	var AAlikafkaConsumerGroupsBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alikafka_consumer_groups.default",
		existMapFunc: existAAlikafkaConsumerGroupsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAAlikafkaConsumerGroupsDataSourceNameMapFunc,
	}

	AAlikafkaConsumerGroupsBusesCheckInfo.dataSourceTestCheck(t, rand, idsConf, consumerIdRegexConf)
}
func testAccCheckAAlikafkaConsumerGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tf-testacc-alikafkaconsumer%d"
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

resource "alicloud_alikafka_consumer_group" "default" {
  instance_id = "${alicloud_alikafka_instance.default.id}"
  consumer_id = "${var.name}"
  tags = {
    	Created = "TF"
		For =     "Test"
  }
}

data "alicloud_alikafka_consumer_groups" "default" {	
	%s
}

`, rand, strings.Join(pairs, " \n "))
	return config
}
