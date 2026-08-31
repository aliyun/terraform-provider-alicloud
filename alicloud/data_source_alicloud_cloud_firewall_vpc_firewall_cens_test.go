package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCloudFirewallVpcFirewallCenDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}_fake"]`,
		}),
	}

	CenIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}"]`,
			"cen_id": `"${alicloud_cen_instance_attachment.attach1.instance_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}_fake"]`,
			"cen_id": `"${alicloud_cen_instance_attachment.attach1.instance_id}_fake"`,
		}),
	}
	StatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}"]`,
			"status": `"opened"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}_fake"]`,
			"status": `"closed"`,
		}),
	}
	VpcFirewallNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}"]`,
			"vpc_firewall_name": `"${alicloud_cloud_firewall_vpc_firewall_cen.default.vpc_firewall_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}_fake"]`,
			"vpc_firewall_name": `"${alicloud_cloud_firewall_vpc_firewall_cen.default.vpc_firewall_name}fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}"]`,
			"cen_id":            `"${alicloud_cen_instance_attachment.attach1.instance_id}"`,
			"status":            `"opened"`,
			"vpc_firewall_name": `"${alicloud_cloud_firewall_vpc_firewall_cen.default.vpc_firewall_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_cloud_firewall_vpc_firewall_cen.default.id}_fake"]`,
			"cen_id":            `"${alicloud_cen_instance_attachment.attach1.instance_id}_fake"`,
			"status":            `"closed"`,
			"vpc_firewall_name": `"${alicloud_cloud_firewall_vpc_firewall_cen.default.vpc_firewall_name}_fake"`,
		}),
	}

	CloudFirewallVpcFirewallCenCheckInfo.dataSourceTestCheck(t, rand, idsConf, CenIdConf, StatusConf, VpcFirewallNameConf, allConf)
}

var existCloudFirewallVpcFirewallCenMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"cens.#":    "1",
		"cens.0.id": CHECKSET,
	}
}

var fakeCloudFirewallVpcFirewallCenMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"cens.#": "0",
	}
}

var CloudFirewallVpcFirewallCenCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cloud_firewall_vpc_firewall_cens.default",
	existMapFunc: existCloudFirewallVpcFirewallCenMapFunc,
	fakeMapFunc:  fakeCloudFirewallVpcFirewallCenMapFunc,
}

func testAccCheckAlicloudCloudFirewallVpcFirewallCenSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
    default = "tf-testacc-%d"
}
data "alicloud_regions" "current" {
  current = true
}
data "alicloud_zones" "zone" {
  available_instance_type = "ecs.sn1ne.large"
  available_resource_creation = "VSwitch"
}

data "alicloud_account" "current" {
}

resource "alicloud_vpc" "foo" {
  vpc_name   = "${var.name}-foo"
  cidr_block = "192.168.0.0/16"
}
resource "alicloud_vpc" "bar" {
  vpc_name   = "${var.name}-bar"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  vpc_id       = alicloud_vpc.foo.id
  cidr_block   = "192.168.10.0/24"
  zone_id      = data.alicloud_zones.zone.zones.0.id
  vswitch_name = "${var.name}-foo"
}

resource "alicloud_vswitch" "bar" {
  vpc_id       = alicloud_vpc.bar.id
  cidr_block   = "172.16.10.0/24"
  zone_id      = data.alicloud_zones.zone.zones.0.id
  vswitch_name = "${var.name}-bar"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  description       = var.name
}

resource "alicloud_cen_instance_attachment" "attach1" {
  instance_id = alicloud_cen_instance.default.id
  child_instance_id = alicloud_vpc.foo.id
  child_instance_type = "VPC"
  child_instance_region_id = data.alicloud_regions.current.ids.0
  child_instance_owner_id = data.alicloud_account.current.id
}
resource "alicloud_cen_instance_attachment" "attach2" {
  instance_id = alicloud_cen_instance.default.id
  child_instance_id = alicloud_vpc.bar.id
  child_instance_type = "VPC"
  child_instance_region_id = data.alicloud_regions.current.ids.0
  child_instance_owner_id = data.alicloud_account.current.id
}

resource "alicloud_cloud_firewall_vpc_firewall_cen" "default" {
  cen_id = alicloud_cen_instance_attachment.attach1.instance_id
  local_vpc {
    network_instance_id = alicloud_cen_instance_attachment.attach1.child_instance_id
  }
  status            = "open"
  member_uid        = data.alicloud_account.current.id
  vpc_region        = data.alicloud_regions.current.ids.0
  vpc_firewall_name = var.name
  lang = "zh"
}

data "alicloud_cloud_firewall_vpc_firewall_cens" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
