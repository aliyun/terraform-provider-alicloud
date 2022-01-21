package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAmqpQueuesDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_amqp_queues.default"
	name := fmt.Sprintf("tf-testacc-amqpQueue%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAmqpQueuesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_amqp_queue.default.instance_id}",
			"virtual_host_name": "${alicloud_amqp_queue.default.virtual_host_name}",
			"name_regex":        "${alicloud_amqp_queue.default.queue_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_amqp_queue.default.instance_id}",
			"virtual_host_name": "${alicloud_amqp_queue.default.virtual_host_name}",
			"name_regex":        "fake_tf-testacc*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_amqp_queue.default.instance_id}",
			"virtual_host_name": "${alicloud_amqp_queue.default.virtual_host_name}",
			"ids":               []string{"${alicloud_amqp_queue.default.queue_name}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
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
			"queues.0.instance_id":       CHECKSET,
			"queues.0.virtual_host_name": name,
			"queues.0.id":                name,
			"queues.0.attributes.%":      CHECKSET,
			"queues.0.auto_delete_state": "true",
			"queues.0.create_time":       CHECKSET,
			"queues.0.exclusive_state":   "false",
			"queues.0.last_consume_time": CHECKSET,
		}
	}

	var fakeAmqpQueuesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"queues.#": "0",
			"names.#":  "0",
			"ids.#":    "0",
		}
	}

	var AmqpQueuesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAmqpQueuesMapFunc,
		fakeMapFunc:  fakeAmqpQueuesMapFunc,
	}
	AmqpQueuesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf)
}

func dataSourceAmqpQueuesConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}
		resource "alicloud_amqp_instance" "default" {
  			instance_type  = "professional"
  			max_tps        = 1000
  			queue_capacity = 50
  			support_eip    = true
  			max_eip_tps    = 128
  			payment_type   = "Subscription"
  			period         = 1
		}
		resource "alicloud_amqp_virtual_host" "default" {
		  instance_id       = alicloud_amqp_instance.default.id
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
