// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEnsGatewayQosDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ens_gateway_qos.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ens_gateway_qos.default.id}_fake"]`,
		}),
	}

	GatewayQosTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_ens_gateway_qos.default.id}"]`,
			"gateway_qos_type": `"Nat"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_ens_gateway_qos.default.id}_fake"]`,
			"gateway_qos_type": `"Nat_fake"`,
		}),
	}
	GatewayQosNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_ens_gateway_qos.default.id}"]`,
			"gateway_qos_name": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_ens_gateway_qos.default.id}_fake"]`,
			"gateway_qos_name": `"${var.name}_fake"`,
		}),
	}
	NetworkIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ens_gateway_qos.default.id}"]`,
			"network_id": `"${alicloud_ens_network.defaultC7YqlT.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ens_gateway_qos.default.id}_fake"]`,
			"network_id": `"${alicloud_ens_network.defaultC7YqlT.id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_ens_gateway_qos.default.id}"]`,
			"gateway_qos_type": `"Nat"`,

			"gateway_qos_name": `"${var.name}"`,

			"network_id": `"${alicloud_ens_network.defaultC7YqlT.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsGatewayQosSourceConfig(rand, map[string]string{
			"ids":              `["${alicloud_ens_gateway_qos.default.id}_fake"]`,
			"gateway_qos_type": `"Nat_fake"`,

			"gateway_qos_name": `"${var.name}_fake"`,

			"network_id": `"${alicloud_ens_network.defaultC7YqlT.id}_fake"`,
		}),
	}

	EnsGatewayQosCheckInfo.dataSourceTestCheck(t, rand, idsConf, GatewayQosTypeConf, GatewayQosNameConf, NetworkIdConf, allConf)
}

var existEnsGatewayQosMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"qos.#":                  "1",
		"qos.0.status":           CHECKSET,
		"qos.0.gateway_qos_name": CHECKSET,
		"qos.0.bandwidth_in":     CHECKSET,
		"qos.0.gateway_qos_type": CHECKSET,
		"qos.0.network_id":       CHECKSET,
		"qos.0.bandwidth_out":    CHECKSET,
		"qos.0.gateway_qos_id":   CHECKSET,
		"qos.0.creation_time":    CHECKSET,
		"qos.0.ens_region_id":    CHECKSET,
	}
}

var fakeEnsGatewayQosMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"qos.#": "0",
	}
}

var EnsGatewayQosCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ens_gateway_qos.default",
	existMapFunc: existEnsGatewayQosMapFunc,
	fakeMapFunc:  fakeEnsGatewayQosMapFunc,
}

func testAccCheckAlicloudEnsGatewayQosSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccEnsGatewayQos%d"
}
variable "ens_region_id" {
  default = "cn-hangzhou-63"
}

resource "alicloud_ens_network" "defaultC7YqlT" {
  network_name  = "镇元-网关限速测试使用"
  cidr_block    = "10.0.0.0/10"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "default5giQWR" {
  cidr_block    = "10.0.8.0/24"
  vswitch_name  = "镇元-网关限速测试"
  ens_region_id = var.ens_region_id
  network_id    = alicloud_ens_network.defaultC7YqlT.id
}

resource "alicloud_ens_instance" "defaultvuczsY" {
  amount      = "1"
  period_unit = "Month"
  auto_renew  = false
  system_disk {
    size = "20"
  }
  scheduling_strategy        = "Concentrate"
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "Subscription"
  instance_type              = "ens.sn1.small"
  password_inherit           = false
  password                   = "12345678abcABC"
  status                     = "Running"
  vswitch_id                 = alicloud_ens_vswitch.default5giQWR.id
  internet_charge_type       = "95BandwidthByMonth"
  instance_name              = "镇元-网关限速测试"
  internet_max_bandwidth_out = "0"
  unique_suffix              = false
  auto_use_coupon            = "true"
  public_ip_identification   = false
  instance_charge_strategy   = "PriceHighPriority"
  ens_region_id              = var.ens_region_id
}



resource "alicloud_ens_gateway_qos" "default" {
  gateway_qos_name = "test"
  bandwidth_in     = "10"
  gateway_qos_type = "Nat"
  network_id       = alicloud_ens_network.defaultC7YqlT.id
  bandwidth_out    = "20"
}

data "alicloud_ens_gateway_qos" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
