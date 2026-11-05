package alikafka_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudAlikafkaSaslUsersDataSource(t *testing.T) {
	rand := tfacctest.RandInt()

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudAlikafkaSaslUsersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alikafka_sasl_user.default.id}"]`,
		}),
		FakeConfig: testAccCheckAliCloudAlikafkaSaslUsersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alikafka_sasl_user.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudAlikafkaSaslUsersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alikafka_sasl_user.default.username}"`,
		}),
		FakeConfig: testAccCheckAliCloudAlikafkaSaslUsersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alikafka_sasl_user.default.username}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudAlikafkaSaslUsersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_alikafka_sasl_user.default.id}"]`,
			"name_regex": `"${alicloud_alikafka_sasl_user.default.username}"`,
		}),
		FakeConfig: testAccCheckAliCloudAlikafkaSaslUsersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_alikafka_sasl_user.default.id}_fake"]`,
			"name_regex": `"${alicloud_alikafka_sasl_user.default.username}_fake"`,
		}),
	}

	var existAliCloudAlikafkaSaslUsersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            "1",
			"names.#":          "1",
			"users.#":          "1",
			"users.0.id":       acctest.CHECKSET,
			"users.0.username": acctest.CHECKSET,
			"users.0.password": acctest.CHECKSET,
			"users.0.type":     acctest.CHECKSET,
		}
	}

	var fakeAliCloudAlikafkaSaslUsersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"users.#": "0",
		}
	}

	var alicloudAlikafkaSaslUsersCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_alikafka_sasl_users.default",
		ExistMapFunc: existAliCloudAlikafkaSaslUsersDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudAlikafkaSaslUsersDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	alicloudAlikafkaSaslUsersCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func testAccCheckAliCloudAlikafkaSaslUsersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testacc-alikafkasasluser-%d"
	}

	data "alicloud_zones" "default" {
  		available_resource_creation = "VSwitch"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "10.4.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "10.4.0.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_alikafka_instance" "default" {
  		name            = var.name
  		partition_num   = 50
  		disk_type       = "1"
  		disk_size       = "500"
  		deploy_type     = "5"
  		io_max          = "20"
  		spec_type       = "professional"
  		service_version = "2.2.0"
  		config          = "{\"enable.acl\":\"true\"}"
  		vswitch_id      = alicloud_vswitch.default.id
  		security_group  = alicloud_security_group.default.id
	}

	resource "alicloud_alikafka_sasl_user" "default" {
  		instance_id = alicloud_alikafka_instance.default.id
  		username    = var.name
  		password    = "YourPassword1234!"
	}

	data "alicloud_alikafka_sasl_users" "default" {
  		instance_id = alicloud_alikafka_sasl_user.default.instance_id
 		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
