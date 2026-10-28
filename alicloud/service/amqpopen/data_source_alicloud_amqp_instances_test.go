package amqpopen_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudAmqpInstancesDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_amqp_instances.default"
	name := fmt.Sprintf("tf-testAcc-AmqpInstance%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceAmqpInstancesConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_amqp_instance.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_amqp_instance.default.id}_fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_amqp_instance.default.instance_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_amqp_instance.default.instance_name}_fake",
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_amqp_instance.default.id}"},
			"status": "SERVING",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_amqp_instance.default.id}"},
			"status": "RELEASED",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_amqp_instance.default.id}"},
			"name_regex": "${alicloud_amqp_instance.default.instance_name}",
			"status":     "SERVING",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"instances.0.id":                    acctest.CHECKSET,
			"instances.0.instance_id":           acctest.CHECKSET,
			"instances.0.instance_type":         acctest.CHECKSET,
			"instances.0.instance_name":         acctest.CHECKSET,
			"instances.0.private_end_point":     acctest.CHECKSET,
			"instances.0.support_eip":           acctest.CHECKSET,
			"instances.0.serverless_switch":     acctest.CHECKSET,
			"instances.0.payment_type":          "",
			"instances.0.renewal_status":        "",
			"instances.0.renewal_duration":      "0",
			"instances.0.renewal_duration_unit": "",
			"instances.0.status":                acctest.CHECKSET,
			"instances.0.expire_time":           acctest.CHECKSET,
			"instances.0.create_time":           acctest.CHECKSET,
		}
	}

	var fakeAliCloudAmqpInstancesDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"instances.#": "0",
		}
	}

	var aliCloudAmqpInstancesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_amqp_instances.default",
		ExistMapFunc: existAliCloudAmqpInstancesDefaultDataSourceMapFunc,
		FakeMapFunc:  fakeAliCloudAmqpInstancesDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		acctest.PreCheck(t)
		acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
	}

	aliCloudAmqpInstancesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}

func TestAccAliCloudAmqpInstancesDataSource_basic1(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_amqp_instances.default"
	name := fmt.Sprintf("tf-testAcc-AmqpInstance%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceAmqpInstancesConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_amqp_instance.default.id}"},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_amqp_instance.default.id}"},
			"enable_details": "false",
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_amqp_instance.default.instance_name}",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_amqp_instance.default.instance_name}",
			"enable_details": "false",
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_amqp_instance.default.id}"},
			"status":         "SERVING",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_amqp_instance.default.id}"},
			"status":         "SERVING",
			"enable_details": "false",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_amqp_instance.default.id}"},
			"name_regex":     "${alicloud_amqp_instance.default.instance_name}",
			"status":         "SERVING",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"instances.0.id":                    acctest.CHECKSET,
			"instances.0.instance_id":           acctest.CHECKSET,
			"instances.0.instance_type":         acctest.CHECKSET,
			"instances.0.instance_name":         acctest.CHECKSET,
			"instances.0.private_end_point":     acctest.CHECKSET,
			"instances.0.support_eip":           acctest.CHECKSET,
			"instances.0.serverless_switch":     acctest.CHECKSET,
			"instances.0.payment_type":          acctest.CHECKSET,
			"instances.0.renewal_status":        acctest.CHECKSET,
			"instances.0.renewal_duration":      acctest.CHECKSET,
			"instances.0.renewal_duration_unit": acctest.CHECKSET,
			"instances.0.status":                acctest.CHECKSET,
			"instances.0.expire_time":           acctest.CHECKSET,
			"instances.0.create_time":           acctest.CHECKSET,
		}
	}

	var fakeAliCloudAmqpInstancesDefaultDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"instances.#":                       "1",
			"instances.0.id":                    acctest.CHECKSET,
			"instances.0.instance_id":           acctest.CHECKSET,
			"instances.0.instance_type":         acctest.CHECKSET,
			"instances.0.instance_name":         acctest.CHECKSET,
			"instances.0.private_end_point":     acctest.CHECKSET,
			"instances.0.support_eip":           acctest.CHECKSET,
			"instances.0.serverless_switch":     acctest.CHECKSET,
			"instances.0.payment_type":          "",
			"instances.0.renewal_status":        "",
			"instances.0.renewal_duration":      "0",
			"instances.0.renewal_duration_unit": "",
			"instances.0.status":                acctest.CHECKSET,
			"instances.0.expire_time":           acctest.CHECKSET,
			"instances.0.create_time":           acctest.CHECKSET,
		}
	}

	var aliCloudAmqpInstancesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_amqp_instances.default",
		ExistMapFunc: existAliCloudAmqpInstancesDefaultDataSourceMapFunc,
		FakeMapFunc:  fakeAliCloudAmqpInstancesDefaultDataSourceMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		acctest.PreCheck(t)
		acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
	}

	aliCloudAmqpInstancesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}

func dataSourceAmqpInstancesConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id     = data.alicloud_vpcs.default.ids.0
  		name_regex = "^default-NODELETING-vsw"
	}

	data "alicloud_security_groups" "default" {
  		vpc_id     = data.alicloud_vpcs.default.ids.0
  		name_regex = "default-NODELETING"
	}

	resource "alicloud_amqp_instance" "default" {
  		instance_name          = var.name
  		payment_type           = "PayAsYouGo"
  		serverless_charge_type = "onDemand"
  		vpc_id                 = data.alicloud_vswitches.default.vpc_id
  		vswitch_ids            = [data.alicloud_vswitches.default.ids.0, data.alicloud_vswitches.default.ids.1]
  		security_group_id      = data.alicloud_security_groups.default.ids.0
	}
`, name)
}
