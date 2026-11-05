// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package amqpopen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudAmqpOpenSourcePermissionDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAmqpOpenSourcePermissionSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_permission.default.id}"]`,
			"instance_id": `"${data.alicloud_amqp_instances.default.ids.0}"`,
			"user_name":   `"${var.user_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudAmqpOpenSourcePermissionSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_permission.default.id}_fake"]`,
			"instance_id": `"${data.alicloud_amqp_instances.default.ids.0}"`,
			"user_name":   `"${var.user_name}"`,
		}),
	}

	InstanceIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAmqpOpenSourcePermissionSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_permission.default.id}"]`,
			"instance_id": `"${data.alicloud_amqp_instances.default.ids.0}"`,
			"user_name":   `"${var.user_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudAmqpOpenSourcePermissionSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_permission.default.id}_fake"]`,
			"instance_id": `"${data.alicloud_amqp_instances.default.ids.0}"`,
			"user_name":   `"${var.user_name}"`,
		}),
	}
	UserNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAmqpOpenSourcePermissionSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_permission.default.id}"]`,
			"instance_id": `"${data.alicloud_amqp_instances.default.ids.0}"`,
			"user_name":   `"${var.user_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudAmqpOpenSourcePermissionSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_permission.default.id}_fake"]`,
			"instance_id": `"${data.alicloud_amqp_instances.default.ids.0}"`,
			"user_name":   `"${var.user_name}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAmqpOpenSourcePermissionSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_permission.default.id}"]`,
			"instance_id": `"${data.alicloud_amqp_instances.default.ids.0}"`,

			"user_name": `"${var.user_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudAmqpOpenSourcePermissionSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_permission.default.id}_fake"]`,
			"instance_id": `"${data.alicloud_amqp_instances.default.ids.0}"`,

			"user_name": `"${var.user_name}_fake"`,
		}),
	}

	AmqpOpenSourcePermissionCheckInfo.DataSourceTestCheck(t, rand, idsConf, InstanceIdConf, UserNameConf, allConf)
}

var existAmqpOpenSourcePermissionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"permissions.#":             "1",
		"permissions.0.write":       acctest.CHECKSET,
		"permissions.0.read":        acctest.CHECKSET,
		"permissions.0.vhost":       acctest.CHECKSET,
		"permissions.0.user_name":   acctest.CHECKSET,
		"permissions.0.instance_id": acctest.CHECKSET,
		"permissions.0.configure":   acctest.CHECKSET,
	}
}

var fakeAmqpOpenSourcePermissionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"permissions.#": "0",
	}
}

var AmqpOpenSourcePermissionCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_amqp_open_source_permissions.default",
	ExistMapFunc: existAmqpOpenSourcePermissionMapFunc,
	FakeMapFunc:  fakeAmqpOpenSourcePermissionMapFunc,
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

data "alicloud_amqp_instances" "default" {
  status = "SERVING"
}



resource "alicloud_amqp_open_source_permission" "default" {
  write       = ".*"
  read        = ".*"
  vhost       = var.vhost
  user_name   = var.user_name
  instance_id = data.alicloud_amqp_instances.default.ids.0
  configure   = ".*"
}

data "alicloud_amqp_open_source_permissions" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
