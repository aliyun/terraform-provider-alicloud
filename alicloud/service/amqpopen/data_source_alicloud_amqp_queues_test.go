package amqpopen_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudAmqpQueuesDataSource(t *testing.T) {

	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_amqp_queues.default"
	name := fmt.Sprintf("tf-testacc-amqpQueue%v", rand)

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceAmqpQueuesConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_amqp_queue.default.instance_id}",
			"virtual_host_name": "${alicloud_amqp_queue.default.virtual_host_name}",
			"name_regex":        "${alicloud_amqp_queue.default.queue_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_amqp_queue.default.instance_id}",
			"virtual_host_name": "${alicloud_amqp_queue.default.virtual_host_name}",
			"name_regex":        "fake_tf-testacc*",
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_amqp_queue.default.instance_id}",
			"virtual_host_name": "${alicloud_amqp_queue.default.virtual_host_name}",
			"ids":               []string{"${alicloud_amqp_queue.default.queue_name}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_amqp_queue.default.instance_id}",
			"virtual_host_name": "${alicloud_amqp_queue.default.virtual_host_name}",
			"ids":               []string{"${alicloud_amqp_queue.default.queue_name}_fake"},
		}),
	}

	var existAmqpQueuesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"names.#":                    "1",
			"queues.#":                   "1",
			"queues.0.queue_name":        name,
			"queues.0.instance_id":       acctest.CHECKSET,
			"queues.0.virtual_host_name": name,
			"queues.0.id":                name,
			"queues.0.attributes.%":      acctest.CHECKSET,
			"queues.0.auto_delete_state": "true",
			"queues.0.create_time":       acctest.CHECKSET,
			"queues.0.exclusive_state":   "false",
			"queues.0.last_consume_time": acctest.CHECKSET,
		}
	}

	var fakeAmqpQueuesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"queues.#": "0",
			"names.#":  "0",
			"ids.#":    "0",
		}
	}

	var AmqpQueuesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existAmqpQueuesMapFunc,
		FakeMapFunc:  fakeAmqpQueuesMapFunc,
	}
	AmqpQueuesCheckInfo.DataSourceTestCheck(t, rand, nameRegexConf, idsConf)
}

func dataSourceAmqpQueuesConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}
		data "alicloud_amqp_instances" "default" {
  			status = "SERVING"
		}
		resource "alicloud_amqp_virtual_host" "default" {
		  instance_id       = data.alicloud_amqp_instances.default.ids.0
		  virtual_host_name = "${var.name}"
		}

		resource "alicloud_amqp_queue" "default" {
		  instance_id = alicloud_amqp_virtual_host.default.instance_id
          virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
		  queue_name = var.name
		  auto_delete_state = true
		}
		`, name)
}
