package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAmqpInstancesDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_amqp_instances.default"
	name := fmt.Sprintf("tf-testacc-amqpInstances%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAmqpInstancesConfigDependence)

	// Currently, there is missing OpenAPI to set instance name when creating the instance.
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_amqp_instance.default.instance_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_tf-testacc*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_amqp_instance.default.instance_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_amqp_instance.default.instance_id}_fake"},
		}),
	}

	var existAmqpInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"instances.#":                       "1",
			"instances.0.instance_name":         name,
			"instances.0.instance_id":           CHECKSET,
			"instances.0.id":                    CHECKSET,
			"instances.0.instance_type":         "professional",
			"instances.0.auto_delete_state":     "true",
			"instances.0.expire_time":           CHECKSET,
			"instances.0.create_time":           CHECKSET,
			"instances.0.payment_type":          "Subscription",
			"instances.0.private_end_point":     CHECKSET,
			"instances.0.public_endpoint":       "",
			"instances.0.renewal_duration_unit": "",
			"instances.0.renewal_status":        "ManualRenewal",
			"instances.0.status":                "SERVING",
			"instances.0.support_eip":           "false",
		}
	}

	var fakeAmqpInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#": "0",
			"names.#":     "0",
			"ids.#":       "0",
		}
	}

	var AmqpInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAmqpInstancesMapFunc,
		fakeMapFunc:  fakeAmqpInstancesMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithTime(t, []int{1})
	}
	AmqpInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf)
}

func dataSourceAmqpInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}

		resource "alicloud_amqp_instance" "default" {
		  instance_name  = var.name
		  instance_type  = "professional"
		  payment_type   = "Subscription"
		  renewal_status = "ManualRenewal"
		  support_eip    = false
		  max_tps = 1000
		  queue_capacity = 50
		}
		`, name)
}
