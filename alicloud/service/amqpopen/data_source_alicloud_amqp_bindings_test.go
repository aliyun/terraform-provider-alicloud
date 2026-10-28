package amqpopen_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudAmqpBindingsDataSource(t *testing.T) {

	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_amqp_bindings.default"
	name := fmt.Sprintf("tf-testacc-amqpbinding%v", rand)

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceAmqpBindingsConfigDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_amqp_binding.default.instance_id}",
			"virtual_host_name": "${alicloud_amqp_binding.default.virtual_host_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_amqp_binding.default.instance_id}",
			"virtual_host_name": "${alicloud_amqp_virtual_host.empty.virtual_host_name}",
		}),
	}

	var existAmqpBindingsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"bindings.#":                   "1",
			"bindings.0.destination_name":  name,
			"bindings.0.instance_id":       acctest.CHECKSET,
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

	var AmqpBindingsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existAmqpBindingsMapFunc,
		FakeMapFunc:  fakeAmqpBindingsMapFunc,
	}

	AmqpBindingsCheckInfo.DataSourceTestCheck(t, rand, idsConf)
}

func dataSourceAmqpBindingsConfigDependence(name string) string {
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

		resource "alicloud_amqp_virtual_host" "empty" {
		  instance_id       = data.alicloud_amqp_instances.default.ids.0
		  virtual_host_name = "${var.name}-empty"
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
