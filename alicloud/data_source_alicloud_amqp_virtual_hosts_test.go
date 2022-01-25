package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAmqpVirtualHostsDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	rand := acctest.RandInt()
	resourceId := "data.alicloud_amqp_virtual_hosts.default"
	name := fmt.Sprintf("tf-testacc-amqpvirtualhost%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAmqpVirtualHostsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_amqp_virtual_host.default.instance_id}",
			"name_regex":  "${alicloud_amqp_virtual_host.default.virtual_host_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_amqp_virtual_host.default.instance_id}",
			"name_regex":  "fake_tf-testacc*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_amqp_virtual_host.default.instance_id}",
			"ids":         []string{"${alicloud_amqp_virtual_host.default.virtual_host_name}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_amqp_virtual_host.default.instance_id}",
			"ids":         []string{"${alicloud_amqp_virtual_host.default.virtual_host_name}_fake"},
		}),
	}

	var existAmqpVirtualHostsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"names.#":                   "1",
			"hosts.#":                   "1",
			"hosts.0.virtual_host_name": name,
			"hosts.0.instance_id":       CHECKSET,
			"hosts.0.id":                name,
		}
	}

	var fakeAmqpVirtualHostsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"hosts.#": "0",
			"names.#": "0",
			"ids.#":   "0",
		}
	}

	var AmqpVirtualHostsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAmqpVirtualHostsMapFunc,
		fakeMapFunc:  fakeAmqpVirtualHostsMapFunc,
	}

	AmqpVirtualHostsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf)
}

func dataSourceAmqpVirtualHostsConfigDependence(name string) string {
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
		  instance_id       =  alicloud_amqp_instance.default.id
		  virtual_host_name = "${var.name}"
		}
		`, name)
}
