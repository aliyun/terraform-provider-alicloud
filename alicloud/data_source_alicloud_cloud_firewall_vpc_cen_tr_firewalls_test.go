package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudFirewallVpcCenTrFirewallDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcCenTrFirewallSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_vpc_cen_tr_firewall.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcCenTrFirewallSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_vpc_cen_tr_firewall.default.id}_fake"]`,
		}),
	}

	CenIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcCenTrFirewallSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cloud_firewall_vpc_cen_tr_firewall.default.id}"]`,
			"cen_id": `"${alicloud_cen_instance.cen.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcCenTrFirewallSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cloud_firewall_vpc_cen_tr_firewall.default.id}_fake"]`,
			"cen_id": `"${alicloud_cen_instance.cen.id}_fake"`,
		}),
	}
	RouteModeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcCenTrFirewallSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_cloud_firewall_vpc_cen_tr_firewall.default.id}"]`,
			"route_mode": `"managed"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcCenTrFirewallSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_cloud_firewall_vpc_cen_tr_firewall.default.id}_fake"]`,
			"route_mode": `"managed_fake"`,
		}),
	}
	RegionNoConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcCenTrFirewallSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cloud_firewall_vpc_cen_tr_firewall.default.id}"]`,
			"region_no": `"${var.region}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcCenTrFirewallSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cloud_firewall_vpc_cen_tr_firewall.default.id}_fake"]`,
			"region_no": `"${var.region}_fake"`,
		}),
	}
	TransitRouterIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcCenTrFirewallSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_cloud_firewall_vpc_cen_tr_firewall.default.id}"]`,
			"transit_router_id": `"${alicloud_cen_transit_router.tr.transit_router_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcCenTrFirewallSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_cloud_firewall_vpc_cen_tr_firewall.default.id}_fake"]`,
			"transit_router_id": `"${alicloud_cen_transit_router.tr.transit_router_id}_fake"`,
		}),
	}
	FirewallNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcCenTrFirewallSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_cloud_firewall_vpc_cen_tr_firewall.default.id}"]`,
			"firewall_name": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcCenTrFirewallSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_cloud_firewall_vpc_cen_tr_firewall.default.id}_fake"]`,
			"firewall_name": `"${var.name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcCenTrFirewallSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cloud_firewall_vpc_cen_tr_firewall.default.id}"]`,
			"cen_id": `"${alicloud_cen_instance.cen.id}"`,

			"route_mode": `"managed"`,

			"region_no": `"${var.region}"`,

			"transit_router_id": `"${alicloud_cen_transit_router.tr.transit_router_id}"`,

			"firewall_name": `"${var.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcCenTrFirewallSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cloud_firewall_vpc_cen_tr_firewall.default.id}_fake"]`,
			"cen_id": `"${alicloud_cen_instance.cen.id}_fake"`,

			"route_mode": `"managed_fake"`,

			"region_no": `"${var.region}_fake"`,

			"transit_router_id": `"${alicloud_cen_transit_router.tr.transit_router_id}_fake"`,

			"firewall_name": `"${var.name}_fake"`,
		}),
	}

	CloudFirewallVpcCenTrFirewallCheckInfo.dataSourceTestCheck(t, rand, idsConf, CenIdConf, RouteModeConf, RegionNoConf, TransitRouterIdConf, FirewallNameConf, allConf)
}

var existCloudFirewallVpcCenTrFirewallMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"firewalls.#":                        "1",
		"firewalls.0.route_mode":             CHECKSET,
		"firewalls.0.region_no":              CHECKSET,
		"firewalls.0.firewall_id":            CHECKSET,
		"firewalls.0.firewall_name":          CHECKSET,
		"firewalls.0.transit_router_id":      CHECKSET,
		"firewalls.0.cen_id":                 CHECKSET,
		"firewalls.0.firewall_switch_status": CHECKSET,
	}
}

var fakeCloudFirewallVpcCenTrFirewallMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"firewalls.#": "0",
	}
}

var CloudFirewallVpcCenTrFirewallCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cloud_firewall_vpc_cen_tr_firewalls.default",
	existMapFunc: existCloudFirewallVpcCenTrFirewallMapFunc,
	fakeMapFunc:  fakeCloudFirewallVpcCenTrFirewallMapFunc,
}

func testAccCheckAlicloudCloudFirewallVpcCenTrFirewallSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCloudFirewallVpcCenTrFirewall%d"
}

variable "description" {
  default = "Created by Terraform"
}

variable "firewall_name" {
  default = "tf-test"
}

variable "tr_attachment_master_cidr" {
  default = "192.168.3.192/26"
}

variable "firewall_subnet_cidr" {
  default = "192.168.3.0/25"
}

variable "region" {
  default = "cn-hangzhou"
}

variable "tr_attachment_slave_cidr" {
  default = "192.168.3.128/26"
}

variable "firewall_vpc_cidr" {
  default = "192.168.3.0/24"
}

variable "zone1" {
  default = "cn-hangzhou-h"
}

variable "firewall_name_update" {
  default = "tf-test-1"
}

variable "zone2" {
  default = "cn-hangzhou-i"
}

data "alicloud_cen_transit_router_available_resources" "default" {
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_cen_instance" "cen" {
  description       = "terraform test"
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "tr" {
  transit_router_name        = var.name
  transit_router_description = "tr-created-by-terraform"
  cen_id                     = alicloud_cen_instance.cen.id
}

resource "alicloud_vpc" "vpc1" {
  description = "created by terraform"
  cidr_block  = "192.168.1.0/24"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "vpc1vsw1" {
  cidr_block   = "192.168.1.0/25"
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.vpc1.id
  zone_id      = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[1]
}

resource "alicloud_vswitch" "vpc1vsw2" {
  vpc_id       = alicloud_vpc.vpc1.id
  cidr_block   = "192.168.1.128/26"
  vswitch_name = var.name
  zone_id      = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[2]
}

resource "alicloud_route_table" "foo" {
  vpc_id           = alicloud_vpc.vpc1.id
  route_table_name = var.name
  description      = var.name
}

resource "alicloud_cen_transit_router_vpc_attachment" "tr-vpc1" {
  zone_mappings {
    vswitch_id = alicloud_vswitch.vpc1vsw1.id
    zone_id    = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[1]
  }
  zone_mappings {
    zone_id    = data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[2]
    vswitch_id = alicloud_vswitch.vpc1vsw2.id
  }
  vpc_id = alicloud_vpc.vpc1.id
  cen_id = alicloud_cen_instance.cen.id
  transit_router_id = alicloud_cen_transit_router.tr.transit_router_id
  depends_on = [alicloud_route_table.foo]
}


resource "alicloud_cloud_firewall_vpc_cen_tr_firewall" "default" {
  cen_id = "${alicloud_cen_transit_router_vpc_attachment.tr-vpc1.cen_id}"
  firewall_name = "${var.name}"
  firewall_subnet_cidr = "${var.firewall_subnet_cidr}"
  firewall_description = "VpcCenTrFirewall created by terraform"
  tr_attachment_master_cidr = "${var.tr_attachment_master_cidr}"
  tr_attachment_master_zone = "${data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[1]}"
  transit_router_id = "${alicloud_cen_transit_router.tr.transit_router_id}"
  tr_attachment_slave_cidr = "${var.tr_attachment_slave_cidr}"
  region_no = "${var.region}"
  route_mode = "managed"
  tr_attachment_slave_zone = "${data.alicloud_cen_transit_router_available_resources.default.resources[0].master_zones[2]}"
  firewall_vpc_cidr = "${var.firewall_vpc_cidr}"
}

data "alicloud_cloud_firewall_vpc_cen_tr_firewalls" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
