package alicloud

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudAlikafkaSaslAclsDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_alikafka_sasl_acls.default"
	name := fmt.Sprintf("tf-testacc-alikafkasaslacl%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAlikafkaSaslAclsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":       alicloud_alikafka_instance.default.id,
			"username":          alicloud_alikafka_sasl_acl.default.username,
			"acl_resource_type": "Topic",
			"acl_resource_name": alicloud_alikafka_sasl_acl.default.acl_resource_name,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":       alicloud_alikafka_instance.default.id,
			"username":          "fake_tf-testacc*",
			"acl_resource_type": "Topic",
			"acl_resource_name": "fake_tf-testacc*",
		}),
	}

	var existAlikafkaSaslAclsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"acls.#":                           "1",
			"acls.0.username":                  fmt.Sprintf("tf-testacc-alikafkasaslacl%v", rand),
			"acls.0.acl_resource_type":         "Topic",
			"acls.0.acl_resource_name":         fmt.Sprintf("tf-testacc-alikafkasaslacl%v", rand),
			"acls.0.acl_resource_pattern_type": "LITERAL",
			"acls.0.host":                      "*",
			"acls.0.acl_operation_type":        "Write",
		}
	}

	var fakeAlikafkaSaslAclsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"acls.#": "0",
		}
	}

	var alikafkaSaslAclsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAlikafkaSaslAclsMapFunc,
		fakeMapFunc:  fakeAlikafkaSaslAclsMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithAlikafkaAclEnable(t)
		testAccPreCheckWithRegions(t, true, connectivity.AlikafkaSupportedRegions)
	}
	alikafkaSaslAclsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf)
}

func dataSourceAlikafkaSaslAclsConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}

		data "alicloud_zones" "default" {
			available_resource_creation= "VSwitch"
		}
		resource "alicloud_vpc" "default" {
		  cidr_block = "172.16.0.0/12"
		  name       = var.name
		}
		
		resource "alicloud_vswitch" "default" {
		  vpc_id = alicloud_vpc.default.id
		  cidr_block = "172.16.0.0/24"
		  availability_zone = data.alicloud_zones.default.zones.0.id
		  name       = var.name
		}

		resource "alicloud_alikafka_instance" "default" {
          name = "tf-testacc-alikafkainstance"
		  topic_quota = "50"
		  disk_type = "1"
		  disk_size = "500"
		  deploy_type = "5"
		  io_max = "20"
          vswitch_id = alicloud_vswitch.default.id
		}

		resource "alicloud_alikafka_topic" "default" {
          instance_id = alicloud_alikafka_instance.default.id
          topic = var.name
		  remark = "topic-remark"
		}
		
		resource "alicloud_alikafka_sasl_user" "default" {
		  instance_id = alicloud_alikafka_instance.default.id
		  username = var.name
		  password = "password"
		}

		resource "alicloud_alikafka_sasl_acl" "default" {
		  instance_id = alicloud_alikafka_instance.default.id
		  username = alicloud_alikafka_sasl_user.default.username
		  acl_resource_type = "Topic"
		  acl_resource_name = alicloud_alikafka_topic.default.topic
		  acl_resource_pattern_type = "LITERAL"
		  acl_operation_type = "Write"
		}
		`, name)
}
