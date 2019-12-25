package alicloud

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudAlikafkaSaslUsersDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_alikafka_sasl_users.default"
	name := fmt.Sprintf("tf-testacc-alikafkasasluser%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAlikafkaSaslUsersConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": alicloud_alikafka_instance.default.id,
			"name_regex":  alicloud_alikafka_sasl_user.default.username,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": alicloud_alikafka_instance.default.id,
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
		
		resource "alicloud_alikafka_sasl_user" "default" {
		  instance_id = alicloud_alikafka_instance.default.id
		  username = var.name
		  password = "password"
		}
		`, name)
}
