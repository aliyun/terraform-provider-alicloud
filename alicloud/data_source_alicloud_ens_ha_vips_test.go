// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEnsHaVipDataSource(t *testing.T) {
	testAccPreCheck(t)
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsHaVipSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ens_ha_vip.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEnsHaVipSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ens_ha_vip.default.id}_fake"]`,
		}),
	}

	IpAddressConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsHaVipSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ens_ha_vip.default.id}"]`,
			"ip_address": `"10.0.9.5"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsHaVipSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ens_ha_vip.default.id}_fake"]`,
			"ip_address": `"10.0.9.5_fake"`,
		}),
	}
	HaVipNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsHaVipSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_ens_ha_vip.default.id}"]`,
			"ha_vip_name": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsHaVipSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_ens_ha_vip.default.id}_fake"]`,
			"ha_vip_name": `"${var.name}_fake"`,
		}),
	}
	VSwitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsHaVipSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ens_ha_vip.default.id}"]`,
			"vswitch_id": `"${alicloud_ens_vswitch.defaultcW3Eib.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsHaVipSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ens_ha_vip.default.id}_fake"]`,
			"vswitch_id": `"${alicloud_ens_vswitch.defaultcW3Eib.id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEnsHaVipSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ens_ha_vip.default.id}"]`,
			"ip_address": `"10.0.9.5"`,

			"ha_vip_name": `"${var.name}"`,

			"vswitch_id": `"${alicloud_ens_vswitch.defaultcW3Eib.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEnsHaVipSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ens_ha_vip.default.id}_fake"]`,
			"ip_address": `"10.0.9.5_fake"`,

			"ha_vip_name": `"${var.name}_fake"`,

			"vswitch_id": `"${alicloud_ens_vswitch.defaultcW3Eib.id}_fake"`,
		}),
	}

	EnsHaVipCheckInfo.dataSourceTestCheck(t, rand, idsConf, IpAddressConf, HaVipNameConf, VSwitchIdConf, allConf)
}

var existEnsHaVipMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"vips.#":               "1",
		"vips.0.status":        CHECKSET,
		"vips.0.description":   CHECKSET,
		"vips.0.create_time":   CHECKSET,
		"vips.0.vswitch_id":    CHECKSET,
		"vips.0.ha_vip_name":   CHECKSET,
		"vips.0.ha_vip_id":     CHECKSET,
		"vips.0.network_id":    CHECKSET,
		"vips.0.ip_address":    CHECKSET,
		"vips.0.ens_region_id": CHECKSET,
	}
}

var fakeEnsHaVipMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"vips.#": "0",
	}
}

var EnsHaVipCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ens_ha_vips.default",
	existMapFunc: existEnsHaVipMapFunc,
	fakeMapFunc:  fakeEnsHaVipMapFunc,
}

func testAccCheckAlicloudEnsHaVipSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccEnsHaVip%d"
}
variable "ens_region_id" {
  default = "cn-hangzhou-58"
}

resource "alicloud_ens_network" "default4wYgcV" {
  network_name  = "tf-testAccEnsHaVip"
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "defaultcW3Eib" {
  cidr_block    = "10.0.9.0/24"
  vswitch_name  = "tf-testAccEnsHaVip"
  ens_region_id = var.ens_region_id
  network_id    = alicloud_ens_network.default4wYgcV.id
}



resource "alicloud_ens_ha_vip" "default" {
  description = "desc1"
  vswitch_id  = alicloud_ens_vswitch.defaultcW3Eib.id
  amount      = 1
  ip_address  = "10.0.9.5"
  ha_vip_name = var.name
}

data "alicloud_ens_ha_vips" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
