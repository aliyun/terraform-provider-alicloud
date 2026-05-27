// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAmqpOpenSourceAccountDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAmqpOpenSourceAccountSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_account.default.id}"]`,
			"instance_id": `"${alicloud_amqp_instance.CreateInstance.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAmqpOpenSourceAccountSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_account.default.id}_fake"]`,
			"instance_id": `"${alicloud_amqp_instance.CreateInstance.id}"`,
		}),
	}

	outputFileConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAmqpOpenSourceAccountSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_account.default.id}"]`,
			"instance_id": `"${alicloud_amqp_instance.CreateInstance.id}"`,
			"output_file": `"./tf-testacc-amqp-open-source-accounts.txt"`,
		}),
		fakeConfig: testAccCheckAlicloudAmqpOpenSourceAccountSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_amqp_open_source_account.default.id}_fake"]`,
			"instance_id": `"${alicloud_amqp_instance.CreateInstance.id}"`,
			"output_file": `"./tf-testacc-amqp-open-source-accounts-fake.txt"`,
		}),
	}

	AmqpOpenSourceAccountCheckInfo.dataSourceTestCheck(t, rand, idsConf, outputFileConf)
}

var existAmqpOpenSourceAccountMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#":             "1",
		"accounts.0.user_name":   CHECKSET,
		"accounts.0.description": CHECKSET,
		"accounts.0.instance_id": CHECKSET,
		"accounts.0.password":    CHECKSET,
	}
}

var fakeAmqpOpenSourceAccountMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"accounts.#": "0",
	}
}

var AmqpOpenSourceAccountCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_amqp_open_source_accounts.default",
	existMapFunc: existAmqpOpenSourceAccountMapFunc,
	fakeMapFunc:  fakeAmqpOpenSourceAccountMapFunc,
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


resource "alicloud_amqp_open_source_account" "default" {
  user_name   = var.user_name
  description = var.user_name
  password    = var.user_name
  instance_id = alicloud_amqp_instance.CreateInstance.id
}

data "alicloud_amqp_open_source_accounts" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
