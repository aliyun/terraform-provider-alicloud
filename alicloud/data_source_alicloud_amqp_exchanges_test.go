package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAmqpExchangesDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_amqp_exchanges.default"
	name := fmt.Sprintf("tf-testacc-amqpexchange%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAmqpExchangesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_amqp_exchange.default.instance_id}",
			"virtual_host_name": "${alicloud_amqp_exchange.default.virtual_host_name}",
			"name_regex":        "${alicloud_amqp_exchange.default.exchange_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_amqp_exchange.default.instance_id}",
			"virtual_host_name": "${alicloud_amqp_exchange.default.virtual_host_name}",
			"name_regex":        "fake_tf-testacc*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_amqp_exchange.default.instance_id}",
			"virtual_host_name": "${alicloud_amqp_exchange.default.virtual_host_name}",
			"ids":               []string{"${alicloud_amqp_exchange.default.exchange_name}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_amqp_exchange.default.instance_id}",
			"virtual_host_name": "${alicloud_amqp_exchange.default.virtual_host_name}",
			"ids":               []string{"${alicloud_amqp_exchange.default.exchange_name}_fake"},
		}),
	}

	var existAmqpExchangesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"exchanges.#":                   "1",
			"exchanges.0.exchange_name":     name,
			"exchanges.0.instance_id":       CHECKSET,
			"exchanges.0.virtual_host_name": name,
			"exchanges.0.id":                name,
			"exchanges.0.attributes.%":      CHECKSET,
			"exchanges.0.auto_delete_state": "true",
			"exchanges.0.exchange_type":     "FANOUT",
		}
	}

	var fakeAmqpExchangesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"exchanges.#": "0",
			"names.#":     "0",
			"ids.#":       "0",
		}
	}

	var AmqpExchangesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAmqpExchangesMapFunc,
		fakeMapFunc:  fakeAmqpExchangesMapFunc,
	}

	AmqpExchangesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf)
}

func dataSourceAmqpExchangesConfigDependence(name string) string {
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
		  exchange_name = var.name
		  auto_delete_state = true
          exchange_type = "FANOUT"
          internal = false
		}
		`, name)
}
