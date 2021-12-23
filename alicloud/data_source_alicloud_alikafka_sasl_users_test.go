package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAlikafkaSaslUsersDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_alikafka_sasl_users.default"
	name := fmt.Sprintf("tf-testacc-alikafkasasluser%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAlikafkaSaslUsersConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_alikafka_instance.default.id}",
			"name_regex":  "${alicloud_alikafka_sasl_user.default.username}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_alikafka_instance.default.id}",
			"name_regex":  "fake_tf-testacc*",
		}),
	}

	var existAlikafkaSaslUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":          "1",
			"users.#":          "1",
			"users.0.username": fmt.Sprintf("tf-testacc-alikafkasasluser%v", rand),
			"users.0.password": "password",
		}
	}

	var fakeAlikafkaSaslUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"users.#": "0",
			"names.#": "0",
		}
	}

	var alikafkaSaslUsersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAlikafkaSaslUsersMapFunc,
		fakeMapFunc:  fakeAlikafkaSaslUsersMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithAlikafkaAclEnable(t)
		testAccPreCheckWithRegions(t, true, connectivity.AlikafkaSupportedRegions)
	}
	alikafkaSaslUsersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf)
}

func dataSourceAlikafkaSaslUsersConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
 default = "%v"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_alikafka_instance" "default" {
  name = "${var.name}"
  topic_quota = "50"
  disk_type = "1"
  disk_size = "500"
  deploy_type = "5"
  io_max = "20"
  vswitch_id = "${data.alicloud_vswitches.default.ids.0}"
  security_group = alicloud_security_group.default.id
}

resource "alicloud_alikafka_sasl_user" "default" {
  instance_id = "${alicloud_alikafka_instance.default.id}"
  username = "${var.name}"
  password = "password"
}
`, name)
}
