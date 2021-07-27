package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAmqpVirtualHostsDataSource(t *testing.T) {

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
		
		data "alicloud_amqp_instances" "default" {
			status = "SERVING"
		}
		resource "alicloud_amqp_virtual_host" "default" {
		  instance_id       = data.alicloud_amqp_instances.default.ids.0
		  virtual_host_name = "${var.name}"
		}
		`, name)
}
