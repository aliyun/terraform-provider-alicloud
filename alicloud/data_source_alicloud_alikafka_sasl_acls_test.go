package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAlikafkaSaslAclsDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_alikafka_sasl_acls.default"
	name := fmt.Sprintf("tf-testacc-alikafkasaslacl%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAlikafkaSaslAclsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_alikafka_instance.default.id}",
			"username":          "${alicloud_alikafka_sasl_acl.default.username}",
			"acl_resource_type": "Topic",
			"acl_resource_name": "${alicloud_alikafka_sasl_acl.default.acl_resource_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":       "${alicloud_alikafka_instance.default.id}",
			"username":          "fake_tf-testacc*",
			"acl_resource_type": "Topic",
			"acl_resource_name": "fake_tf-testacc*",
		}),
	}

	var existAlikafkaSaslAclsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"acls.#":                           "1",
			"acls.0.id":                        CHECKSET,
			"acls.0.username":                  fmt.Sprintf("tf-testacc-alikafkasaslacl%v", rand),
			"acls.0.acl_resource_type":         "TOPIC",
			"acls.0.acl_resource_name":         fmt.Sprintf("tf-testacc-alikafkasaslacl%v", rand),
			"acls.0.acl_resource_pattern_type": "LITERAL",
			"acls.0.host":                      "*",
			"acls.0.acl_operation_type":        "WRITE",
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
		testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		testAccPreCheck(t)
	}
	alikafkaSaslAclsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf)
}

func dataSourceAlikafkaSaslAclsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%v"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vswitches" "default" {
  zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vswitches.default.vswitches.0.vpc_id
}

resource "alicloud_alikafka_instance" "default" {
  name            = var.name
  deploy_type     = "4"
  instance_type   = "alikafka_serverless"
  vswitch_id      = data.alicloud_vswitches.default.vswitches.0.id
  spec_type       = "normal"
  service_version = "3.3.1"
  security_group  = alicloud_security_group.default.id
  config          = "{\"enable.acl\":\"true\"}"
  serverless_config {
    reserved_publish_capacity   = 60
    reserved_subscribe_capacity = 60
  }
}

resource "alicloud_alikafka_topic" "default" {
  instance_id = alicloud_alikafka_instance.default.id
  topic       = var.name
  remark      = "topic-remark"
}

resource "alicloud_alikafka_sasl_user" "default" {
  instance_id = alicloud_alikafka_instance.default.id
  username    = var.name
  type        = "scram"
  password    = "YourPassword1234!"
}

resource "alicloud_alikafka_sasl_acl" "default" {
  instance_id               = alicloud_alikafka_instance.default.id
  username                  = alicloud_alikafka_sasl_user.default.username
  acl_resource_type         = "Topic"
  acl_resource_name         = alicloud_alikafka_topic.default.topic
  acl_resource_pattern_type = "LITERAL"
  acl_operation_type        = "Write"
}
`, name)
}
