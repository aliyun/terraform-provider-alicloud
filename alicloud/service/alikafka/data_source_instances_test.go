package alikafka_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudAlikafkaInstancesDataSource(t *testing.T) {

	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_alikafka_instances.default"
	name := fmt.Sprintf("tf-testacc-alikafkainstance%v", rand)

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceAlikafkaInstancesConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"name_regex":     "${alicloud_alikafka_instance.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_tf-testacc*",
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_alikafka_instance.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_alikafka_instance.default.id}_fake"},
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
			"ids":            []string{"${alicloud_alikafka_instance.default.id}"},
			"name_regex":     "${alicloud_alikafka_instance.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_alikafka_instance.default.id}_fake"},
			"name_regex": "${alicloud_alikafka_instance.default.name}_fake",
		}),
	}

	var existAlikafkaInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                     "1",
			"instances.#":                               "1",
			"instances.0.name":                          fmt.Sprintf("tf-testacc-alikafkainstance%v", rand),
			"instances.0.partition_num":                 "50",
			"instances.0.topic_quota":                   acctest.CHECKSET,
			"instances.0.disk_type":                     "1",
			"instances.0.disk_size":                     "500",
			"instances.0.deploy_type":                   "5",
			"instances.0.io_max":                        "20",
			"instances.0.paid_type":                     "PostPaid",
			"instances.0.spec_type":                     "normal",
			"instances.0.service_version":               acctest.CHECKSET,
			"instances.0.config":                        acctest.CHECKSET,
			"instances.0.security_group":                acctest.CHECKSET,
			"instances.0.vpc_id":                        acctest.CHECKSET,
			"instances.0.zone_id":                       acctest.CHECKSET,
			"instances.0.end_point":                     acctest.CHECKSET,
			"instances.0.expired_time":                  acctest.CHECKSET,
			"instances.0.msg_retain":                    acctest.CHECKSET,
			"instances.0.ssl_end_point":                 "",
			"instances.0.tags.%":                        "2",
			"instances.0.tags.Created":                  "TF",
			"instances.0.tags.For":                      "Test",
			"instances.0.upgrade_service_detail_info.#": "1",
			"instances.0.upgrade_service_detail_info.0.current2_open_source_version": acctest.CHECKSET,
			"instances.0.domain_endpoint":                                            acctest.CHECKSET,
			"instances.0.ssl_domain_endpoint":                                        "",
			"instances.0.sasl_domain_endpoint":                                       "",
			"instances.0.allowed_list.0.deploy_type":                                 acctest.CHECKSET,
			"instances.0.allowed_list.0.vpc_list.0.allowed_ip_list.#":                acctest.CHECKSET,
			"instances.0.allowed_list.0.vpc_list.0.port_range":                       acctest.CHECKSET,
		}
	}

	var fakeAlikafkaInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var alikafkaInstancesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existAlikafkaInstancesMapFunc,
		FakeMapFunc:  fakeAlikafkaInstancesMapFunc,
	}
	alikafkaInstancesCheckInfo.DataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)

}

func dataSourceAlikafkaInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
 default = "%v"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = "${data.alicloud_vswitches.default.vswitches.0.vpc_id}"
}

resource "alicloud_alikafka_instance" "default" {
  name = "${var.name}"
  partition_num = "50"
  disk_type = "1"
  disk_size = "500"
  deploy_type = "5"
  io_max = "20"
  vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
  security_group = "${alicloud_security_group.default.id}"
  tags = {
       Created = "TF"
      For =     "Test"
  }
}
`, name)
}
