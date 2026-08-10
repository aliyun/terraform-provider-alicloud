// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEnsNatGatewayForwardEntryDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsNatGatewayForwardEntrySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ens_nat_gateway_forward_entry.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEnsNatGatewayForwardEntrySourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ens_nat_gateway_forward_entry.default.id}_fake"]`,
		}),
	}

	InternalIpConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsNatGatewayForwardEntrySourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_ens_nat_gateway_forward_entry.default.id}"]`,
			"internal_ip": `"${alicloud_ens_instance.defaulth6OQ3p.private_ip_address}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsNatGatewayForwardEntrySourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_ens_nat_gateway_forward_entry.default.id}_fake"]`,
			"internal_ip": `"${alicloud_ens_instance.defaulth6OQ3p.private_ip_address}_fake"`,
		}),
	}
	IpProtocolConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsNatGatewayForwardEntrySourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_ens_nat_gateway_forward_entry.default.id}"]`,
			"ip_protocol": `"TCP"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsNatGatewayForwardEntrySourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_ens_nat_gateway_forward_entry.default.id}_fake"]`,
			"ip_protocol": `"TCP_fake"`,
		}),
	}
	NatGatewayIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsNatGatewayForwardEntrySourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ens_nat_gateway_forward_entry.default.id}"]`,
			"nat_gateway_id": `"${alicloud_ens_nat_gateway.defaultlZ7YKl.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsNatGatewayForwardEntrySourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_ens_nat_gateway_forward_entry.default.id}_fake"]`,
			"nat_gateway_id": `"${alicloud_ens_nat_gateway.defaultlZ7YKl.id}_fake"`,
		}),
	}
	ExternalIpConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsNatGatewayForwardEntrySourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_ens_nat_gateway_forward_entry.default.id}"]`,
			"external_ip": `"${alicloud_ens_eip.defaultLQgQB6.ip_address}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsNatGatewayForwardEntrySourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_ens_nat_gateway_forward_entry.default.id}_fake"]`,
			"external_ip": `"${alicloud_ens_eip.defaultLQgQB6.ip_address}_fake"`,
		}),
	}
	ForwardEntryNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsNatGatewayForwardEntrySourceConfig(rand, map[string]string{
			"ids":                `["${alicloud_ens_nat_gateway_forward_entry.default.id}"]`,
			"forward_entry_name": `"测试用例-dnat"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsNatGatewayForwardEntrySourceConfig(rand, map[string]string{
			"ids":                `["${alicloud_ens_nat_gateway_forward_entry.default.id}_fake"]`,
			"forward_entry_name": `"测试用例-dnat_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsNatGatewayForwardEntrySourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_ens_nat_gateway_forward_entry.default.id}"]`,
			"internal_ip": `"${alicloud_ens_instance.defaulth6OQ3p.private_ip_address}"`,

			"ip_protocol": `"TCP"`,

			"nat_gateway_id": `"${alicloud_ens_nat_gateway.defaultlZ7YKl.id}"`,

			"external_ip": `"${alicloud_ens_eip.defaultLQgQB6.ip_address}"`,

			"forward_entry_name": `"测试用例-dnat"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsNatGatewayForwardEntrySourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_ens_nat_gateway_forward_entry.default.id}_fake"]`,
			"internal_ip": `"${alicloud_ens_instance.defaulth6OQ3p.private_ip_address}_fake"`,

			"ip_protocol": `"TCP_fake"`,

			"nat_gateway_id": `"${alicloud_ens_nat_gateway.defaultlZ7YKl.id}_fake"`,

			"external_ip": `"${alicloud_ens_eip.defaultLQgQB6.ip_address}_fake"`,

			"forward_entry_name": `"测试用例-dnat_fake"`,
		}),
	}

	EnsNatGatewayForwardEntryCheckInfo.dataSourceTestCheck(t, rand, idsConf, InternalIpConf, IpProtocolConf, NatGatewayIdConf, ExternalIpConf, ForwardEntryNameConf, allConf)
}

var existEnsNatGatewayForwardEntryMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"entries.#":                    "1",
		"entries.0.status":             CHECKSET,
		"entries.0.external_port":      CHECKSET,
		"entries.0.external_ip":        CHECKSET,
		"entries.0.forward_entry_id":   CHECKSET,
		"entries.0.ip_protocol":        CHECKSET,
		"entries.0.internal_port":      CHECKSET,
		"entries.0.health_check_port":  CHECKSET,
		"entries.0.forward_entry_name": CHECKSET,
		"entries.0.nat_gateway_id":     CHECKSET,
		"entries.0.internal_ip":        CHECKSET,
	}
}

var fakeEnsNatGatewayForwardEntryMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"entries.#": "0",
	}
}

var EnsNatGatewayForwardEntryCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ens_nat_gateway_forward_entries.default",
	existMapFunc: existEnsNatGatewayForwardEntryMapFunc,
	fakeMapFunc:  fakeEnsNatGatewayForwardEntryMapFunc,
}

func testAccCheckAlicloudEnsNatGatewayForwardEntrySourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccEnsNatGatewayForwardEntry%d"
}
variable "ens_region_id" {
  default = "cn-hangzhou-44"
}

resource "alicloud_ens_network" "default6T9qR2" {
  network_name  = "测试用例_Dnat"
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "default5BAAN2" {
  cidr_block    = "10.0.6.0/24"
  vswitch_name  = "测试用例-dnat"
  ens_region_id = alicloud_ens_network.default6T9qR2.ens_region_id
  network_id    = alicloud_ens_network.default6T9qR2.id
}

resource "alicloud_ens_nat_gateway" "defaultlZ7YKl" {
  vswitch_id    = alicloud_ens_vswitch.default5BAAN2.id
  ens_region_id = alicloud_ens_vswitch.default5BAAN2.ens_region_id
  network_id    = alicloud_ens_vswitch.default5BAAN2.network_id
  instance_type = "enat.default"
  nat_name      = "测试用例-dnat"
}

resource "alicloud_ens_eip" "defaultLQgQB6" {
  bandwidth            = "5"
  payment_type         = "PayAsYouGo"
  ens_region_id        = var.ens_region_id
  eip_name             = "测试用例-dnat"
  internet_charge_type = "95BandwidthByMonth"
}

resource "alicloud_ens_eip_instance_attachment" "defaultc19VZl" {
  instance_id   = alicloud_ens_nat_gateway.defaultlZ7YKl.id
  allocation_id = alicloud_ens_eip.defaultLQgQB6.id
  instance_type = "Nat"
}

resource "alicloud_ens_instance" "defaulth6OQ3p" {
  auto_renew = false
  system_disk {
    size     = "20"
    category = "cloud_efficiency"
  }
  scheduling_strategy        = "Concentrate"
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "Subscription"
  instance_type              = "ens.sn1.stiny"
  password_inherit           = false
  password                   = "12345678abcABC"
  status                     = "Running"
  amount                     = "1"
  vswitch_id                 = alicloud_ens_vswitch.default5BAAN2.id
  internet_charge_type       = "95BandwidthByMonth"
  instance_name              = "测试用例-dnat"
  internet_max_bandwidth_out = "0"
  unique_suffix              = false
  auto_use_coupon            = "true"
  public_ip_identification   = false
  instance_charge_strategy   = "PriceHighPriority"
  ens_region_id              = var.ens_region_id
  period_unit                = "Month"
}

resource "alicloud_ens_eip" "eip2" {
  bandwidth            = "5"
  payment_type         = "PayAsYouGo"
  ens_region_id        = var.ens_region_id
  eip_name             = "测试用例-dnat2"
  internet_charge_type = "95BandwidthByMonth"
}

resource "alicloud_ens_eip_instance_attachment" "default4Ph8bE" {
  instance_id   = alicloud_ens_nat_gateway.defaultlZ7YKl.id
  allocation_id = alicloud_ens_eip.eip2.id
  instance_type = "Nat"
}

resource "alicloud_ens_instance" "instance2" {
  auto_renew = false
  system_disk {
    size     = "20"
    category = "cloud_efficiency"
  }
  scheduling_strategy        = "Concentrate"
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "Subscription"
  instance_type              = "ens.sn1.stiny"
  password_inherit           = false
  password                   = "12345678abcABC"
  status                     = "Running"
  amount                     = "1"
  vswitch_id                 = alicloud_ens_vswitch.default5BAAN2.id
  internet_charge_type       = "95BandwidthByMonth"
  instance_name              = "测试用例-dnat2"
  internet_max_bandwidth_out = "0"
  unique_suffix              = false
  auto_use_coupon            = "true"
  public_ip_identification   = false
  instance_charge_strategy   = "PriceHighPriority"
  ens_region_id              = var.ens_region_id
  period_unit                = "Month"
}



resource "alicloud_ens_nat_gateway_forward_entry" "default" {
  external_port      = "100/200"
  external_ip        = alicloud_ens_eip.defaultLQgQB6.ip_address
  ip_protocol        = "TCP"
  internal_port      = "100/200"
  health_check_port  = "150"
  nat_gateway_id     = alicloud_ens_nat_gateway.defaultlZ7YKl.id
  forward_entry_name = "测试用例-dnat"
  internal_ip        = alicloud_ens_instance.defaulth6OQ3p.private_ip_address
}

data "alicloud_ens_nat_gateway_forward_entries" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
