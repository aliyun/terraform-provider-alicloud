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

func TestAccAlicloudAmqpOpenSourceAccountDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAmqpOpenSourceAccountSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_account.default.id}"]`,
			"instance_id": `"${data.alicloud_amqp_instances.default.ids.0}"`,
		}),
		FakeConfig: testAccCheckAlicloudAmqpOpenSourceAccountSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_account.default.id}_fake"]`,
			"instance_id": `"${data.alicloud_amqp_instances.default.ids.0}"`,
		}),
	}

	outputFileConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAmqpOpenSourceAccountSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_account.default.id}"]`,
			"instance_id": `"${data.alicloud_amqp_instances.default.ids.0}"`,
			"output_file": `"./tf-testacc-amqp-open-source-accounts.txt"`,
		}),
		FakeConfig: testAccCheckAlicloudAmqpOpenSourceAccountSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_account.default.id}_fake"]`,
			"instance_id": `"${data.alicloud_amqp_instances.default.ids.0}"`,
			"output_file": `"./tf-testacc-amqp-open-source-accounts-fake.txt"`,
		}),
	}

	AmqpOpenSourceAccountCheckInfo.DataSourceTestCheck(t, rand, idsConf, outputFileConf)
}

var existAmqpOpenSourceAccountMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#":             "1",
		"accounts.0.user_name":   acctest.CHECKSET,
		"accounts.0.description": acctest.CHECKSET,
		"accounts.0.instance_id": acctest.CHECKSET,
		"accounts.0.password":    acctest.CHECKSET,
	}
}

var fakeAmqpOpenSourceAccountMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#": "0",
	}
}

var AmqpOpenSourceAccountCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_amqp_open_source_accounts.default",
	ExistMapFunc: existAmqpOpenSourceAccountMapFunc,
	FakeMapFunc:  fakeAmqpOpenSourceAccountMapFunc,
}

func testAccCheckAlicloudAmqpOpenSourceAccountSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccAmqpOpenSourceAccount%d"
}
variable "instance_name" {
  default = "测试开源鉴权实例"
}

variable "user_name" {
  default = "Suhao123_"
}

variable "user_name_update" {
  default = "Suhao456_"
}

data "alicloud_amqp_instances" "default" {
  status = "SERVING"
}


resource "alicloud_amqp_open_source_account" "default" {
  user_name   = var.user_name
  description = var.user_name
  password    = var.user_name
  instance_id = data.alicloud_amqp_instances.default.ids.0
}

data "alicloud_amqp_open_source_accounts" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
