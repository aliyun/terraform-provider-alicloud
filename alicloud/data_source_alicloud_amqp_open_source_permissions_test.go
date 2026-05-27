// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlicloudAmqpOpenSourcePermissionDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAmqpOpenSourcePermissionSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_permission.default.id}"]`,
			"instance_id": `"${alicloud_amqp_instance.CreateInstance.id}"`,
			"user_name":   `"${var.user_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAmqpOpenSourcePermissionSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_permission.default.id}_fake"]`,
			"instance_id": `"${alicloud_amqp_instance.CreateInstance.id}"`,
			"user_name":   `"${var.user_name}"`,
		}),
	}

	InstanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAmqpOpenSourcePermissionSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_permission.default.id}"]`,
			"instance_id": `"${alicloud_amqp_instance.CreateInstance.id}"`,
			"user_name":   `"${var.user_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAmqpOpenSourcePermissionSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_permission.default.id}_fake"]`,
			"instance_id": `"${alicloud_amqp_instance.CreateInstance.id}"`,
			"user_name":   `"${var.user_name}"`,
		}),
	}
	UserNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAmqpOpenSourcePermissionSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_permission.default.id}"]`,
			"instance_id": `"${alicloud_amqp_instance.CreateInstance.id}"`,
			"user_name":   `"${var.user_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAmqpOpenSourcePermissionSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_permission.default.id}_fake"]`,
			"instance_id": `"${alicloud_amqp_instance.CreateInstance.id}"`,
			"user_name":   `"${var.user_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAmqpOpenSourcePermissionSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_permission.default.id}"]`,
			"instance_id": `"${alicloud_amqp_instance.CreateInstance.id}"`,

			"user_name": `"${var.user_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAmqpOpenSourcePermissionSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_permission.default.id}_fake"]`,
			"instance_id": `"${alicloud_amqp_instance.CreateInstance.id}"`,

			"user_name": `"${var.user_name}_fake"`,
		}),
	}

	AmqpOpenSourcePermissionCheckInfo.dataSourceTestCheck(t, rand, idsConf, InstanceIdConf, UserNameConf, allConf)
}

var existAmqpOpenSourcePermissionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"permissions.#":             "1",
		"permissions.0.write":       CHECKSET,
		"permissions.0.read":        CHECKSET,
		"permissions.0.vhost":       CHECKSET,
		"permissions.0.user_name":   CHECKSET,
		"permissions.0.instance_id": CHECKSET,
		"permissions.0.configure":   CHECKSET,
	}
}

var fakeAmqpOpenSourcePermissionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"permissions.#": "0",
	}
}

var AmqpOpenSourcePermissionCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_amqp_open_source_permissions.default",
	existMapFunc: existAmqpOpenSourcePermissionMapFunc,
	fakeMapFunc:  fakeAmqpOpenSourcePermissionMapFunc,
}

func testAccCheckAlicloudAmqpOpenSourcePermissionSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccAmqpOpenSourcePermission%d"
}
variable "instance_name" {
  default = "测试开源鉴权实例"
}

variable "vhost" {
  default = "/"
}

variable "user_name" {
  default = "Suhao123_WithPer"
}

resource "alicloud_amqp_instance" "CreateInstance" {
  renewal_duration      = "1"
  max_tps               = "3000"
  period_cycle          = "Month"
  max_connections       = "2000"
  support_eip           = true
  auto_renew            = false
  renewal_status        = "AutoRenewal"
  period                = "12"
  instance_name         = var.instance_name
  support_tracing       = false
  payment_type          = "Subscription"
  renewal_duration_unit = "Month"
  instance_type         = "enterprise"
  queue_capacity        = "200"
  max_eip_tps           = "128"
  vpc_id                = alicloud_vpc.default.id
  vswitch_ids           = [alicloud_vswitch.default_b.id, alicloud_vswitch.default_g.id]
  security_group_id     = alicloud_security_group.default.id
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default_b" {
  vswitch_name = "${var.name}-b"
  cidr_block   = "172.16.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-hangzhou-b"
}

resource "alicloud_vswitch" "default_g" {
  vswitch_name = "${var.name}-g"
  cidr_block   = "172.16.1.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-hangzhou-g"
}

resource "alicloud_security_group" "default" {
  security_group_name = var.name
  vpc_id              = alicloud_vpc.default.id
}



resource "alicloud_amqp_open_source_permission" "default" {
  write       = ".*"
  read        = ".*"
  vhost       = var.vhost
  user_name   = var.user_name
  instance_id = alicloud_amqp_instance.CreateInstance.id
  configure   = ".*"
}

data "alicloud_amqp_open_source_permissions" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
