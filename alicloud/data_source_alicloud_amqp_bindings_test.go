package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAmqpBindingsDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_amqp_bindings.default"
	name := fmt.Sprintf("tf-testacc-amqpbinding%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAmqpBindingsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_amqp_binding.default.instance_id}",
			"virtual_host_name": "${alicloud_amqp_binding.default.virtual_host_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_amqp_binding.default.instance_id}",
			"virtual_host_name": "${alicloud_amqp_binding.default.virtual_host_name}-fake",
		}),
	}

	var existAmqpBindingsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"bindings.#":                   "1",
			"bindings.0.destination_name":  name,
			"bindings.0.instance_id":       CHECKSET,
			"bindings.0.virtual_host_name": name,
			"bindings.0.source_exchange":   name,
			"bindings.0.binding_type":      "EXCHANGE",
			"bindings.0.binding_key":       name + "-2",
			"bindings.0.argument":          "x-match:all",
		}
	}

	var fakeAmqpBindingsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"bindings.#": "0",
			"ids.#":      "0",
		}
	}

	var AmqpBindingsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAmqpBindingsMapFunc,
		fakeMapFunc:  fakeAmqpBindingsMapFunc,
	}

	AmqpBindingsCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

func dataSourceAmqpBindingsConfigDependence(name string) string {
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

		resource "alicloud_amqp_exchange" "default" {
			instance_id = alicloud_amqp_virtual_host.default.instance_id
			virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
			auto_delete_state = true
			exchange_name = var.name
			exchange_type = "HEADERS"
			internal = false
		}
		resource "alicloud_amqp_exchange" "default2" {
			instance_id = alicloud_amqp_virtual_host.default.instance_id
			virtual_host_name = alicloud_amqp_virtual_host.default.virtual_host_name
			auto_delete_state = true
			exchange_name = "${var.name}-2"
			exchange_type = "HEADERS"
			internal = false
		}

		resource "alicloud_amqp_binding" "default" {
		  instance_id = alicloud_amqp_exchange.default.instance_id
          virtual_host_name = alicloud_amqp_exchange.default.virtual_host_name
		  argument = "x-match:all"
          binding_key = alicloud_amqp_exchange.default2.exchange_name
          binding_type = "EXCHANGE"
		  destination_name = var.name
		  source_exchange = alicloud_amqp_exchange.default.exchange_name
		}
		`, name)
}
