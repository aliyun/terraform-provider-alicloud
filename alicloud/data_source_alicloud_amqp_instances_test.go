package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudAmqpInstancesDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_amqp_instances.default"
	name := fmt.Sprintf("tf-testAcc-AmqpInstance%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAmqpInstancesConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_amqp_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_amqp_instance.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_amqp_instance.default.instance_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_amqp_instance.default.instance_name}_fake",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_amqp_instance.default.id}"},
			"status": "SERVING",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_amqp_instance.default.id}"},
			"status": "RELEASED",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_amqp_instance.default.id}"},
			"name_regex": "${alicloud_amqp_instance.default.instance_name}",
			"status":     "SERVING",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_amqp_instance.default.id}_fake"},
			"name_regex": "${alicloud_amqp_instance.default.instance_name}_fake",
			"status":     "RELEASED",
		}),
	}

	var existAliCloudAmqpInstancesDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"instances.#":                       "1",
			"instances.0.id":                    CHECKSET,
			"instances.0.instance_id":           CHECKSET,
			"instances.0.instance_type":         CHECKSET,
			"instances.0.instance_name":         CHECKSET,
			"instances.0.public_endpoint":       CHECKSET,
			"instances.0.private_end_point":     CHECKSET,
			"instances.0.support_eip":           CHECKSET,
			"instances.0.payment_type":          "",
			"instances.0.renewal_status":        "",
			"instances.0.renewal_duration":      "0",
			"instances.0.renewal_duration_unit": "",
			"instances.0.status":                CHECKSET,
			"instances.0.expire_time":           CHECKSET,
			"instances.0.create_time":           CHECKSET,
		}
	}

	var fakeAliCloudAmqpInstancesDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"instances.#": "0",
		}
	}

	var aliCloudAmqpInstancesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_amqp_instances.default",
		existMapFunc: existAliCloudAmqpInstancesDefaultDataSourceMapFunc,
		fakeMapFunc:  fakeAliCloudAmqpInstancesDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithAccountSiteType(t, DomesticSite)
	}

	aliCloudAmqpInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}

func TestAccAliCloudAmqpInstancesDataSource_basic1(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_amqp_instances.default"
	name := fmt.Sprintf("tf-testAcc-AmqpInstance%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAmqpInstancesConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_amqp_instance.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_amqp_instance.default.id}"},
			"enable_details": "false",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_amqp_instance.default.instance_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_amqp_instance.default.instance_name}",
			"enable_details": "false",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_amqp_instance.default.id}"},
			"status":         "SERVING",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_amqp_instance.default.id}"},
			"status":         "SERVING",
			"enable_details": "false",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_amqp_instance.default.id}"},
			"name_regex":     "${alicloud_amqp_instance.default.instance_name}",
			"status":         "SERVING",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_amqp_instance.default.id}"},
			"name_regex":     "${alicloud_amqp_instance.default.instance_name}",
			"status":         "SERVING",
			"enable_details": "false",
		}),
	}

	var existAliCloudAmqpInstancesDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"instances.#":                       "1",
			"instances.0.id":                    CHECKSET,
			"instances.0.instance_id":           CHECKSET,
			"instances.0.instance_type":         CHECKSET,
			"instances.0.instance_name":         CHECKSET,
			"instances.0.public_endpoint":       CHECKSET,
			"instances.0.private_end_point":     CHECKSET,
			"instances.0.support_eip":           CHECKSET,
			"instances.0.payment_type":          CHECKSET,
			"instances.0.renewal_status":        CHECKSET,
			"instances.0.renewal_duration":      CHECKSET,
			"instances.0.renewal_duration_unit": CHECKSET,
			"instances.0.status":                CHECKSET,
			"instances.0.expire_time":           CHECKSET,
			"instances.0.create_time":           CHECKSET,
		}
	}

	var fakeAliCloudAmqpInstancesDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"instances.#":                       "1",
			"instances.0.id":                    CHECKSET,
			"instances.0.instance_id":           CHECKSET,
			"instances.0.instance_type":         CHECKSET,
			"instances.0.instance_name":         CHECKSET,
			"instances.0.public_endpoint":       CHECKSET,
			"instances.0.private_end_point":     CHECKSET,
			"instances.0.support_eip":           CHECKSET,
			"instances.0.payment_type":          "",
			"instances.0.renewal_status":        "",
			"instances.0.renewal_duration":      "0",
			"instances.0.renewal_duration_unit": "",
			"instances.0.status":                CHECKSET,
			"instances.0.expire_time":           CHECKSET,
			"instances.0.create_time":           CHECKSET,
		}
	}

	var aliCloudAmqpInstancesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_amqp_instances.default",
		existMapFunc: existAliCloudAmqpInstancesDefaultDataSourceMapFunc,
		fakeMapFunc:  fakeAliCloudAmqpInstancesDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithAccountSiteType(t, DomesticSite)
	}

	aliCloudAmqpInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}

func dataSourceAmqpInstancesConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_amqp_instance" "default" {
  		instance_name         = var.name
  		instance_type         = "enterprise"
  		max_tps               = 3000
  		max_connections       = 2000
  		queue_capacity        = 200
  		payment_type          = "Subscription"
  		renewal_status        = "AutoRenewal"
  		renewal_duration      = 1
  		renewal_duration_unit = "Year"
  		support_eip           = true
	}
`, name)
}
