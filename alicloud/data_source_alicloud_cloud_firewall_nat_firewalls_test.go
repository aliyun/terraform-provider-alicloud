package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudFirewallNatFirewallDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallNatFirewallSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_nat_firewall.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallNatFirewallSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_nat_firewall.default.id}_fake"]`,
		}),
	}

	LangConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallNatFirewallSourceConfig(rand, map[string]string{
			"ids":  `["${alicloud_cloud_firewall_nat_firewall.default.id}"]`,
			"lang": `"zh"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallNatFirewallSourceConfig(rand, map[string]string{
			"ids":  `["${alicloud_cloud_firewall_nat_firewall.default.id}_fake"]`,
			"lang": `"zh_fake"`,
		}),
	}
	ProxyNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallNatFirewallSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_cloud_firewall_nat_firewall.default.id}"]`,
			"proxy_name": `"YQC-防火墙测试"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallNatFirewallSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_cloud_firewall_nat_firewall.default.id}_fake"]`,
			"proxy_name": `"YQC-防火墙测试_fake"`,
		}),
	}
	VpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallNatFirewallSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cloud_firewall_nat_firewall.default.id}"]`,
			"vpc_id": `"${alicloud_vpc.defaultikZ0gD.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallNatFirewallSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cloud_firewall_nat_firewall.default.id}_fake"]`,
			"vpc_id": `"${alicloud_vpc.defaultikZ0gD.id}_fake"`,
		}),
	}
	NatGatewayIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallNatFirewallSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_cloud_firewall_nat_firewall.default.id}"]`,
			"nat_gateway_id": `"${alicloud_nat_gateway.default2iRZpC.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallNatFirewallSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_cloud_firewall_nat_firewall.default.id}_fake"]`,
			"nat_gateway_id": `"${alicloud_nat_gateway.default2iRZpC.id}_fake"`,
		}),
	}
	StatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallNatFirewallSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cloud_firewall_nat_firewall.default.id}"]`,
			"status": `"closed"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallNatFirewallSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cloud_firewall_nat_firewall.default.id}_fake"]`,
			"status": `"closed_fake"`,
		}),
	}
	RegionNoConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallNatFirewallSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cloud_firewall_nat_firewall.default.id}"]`,
			"region_no": `"cn-hangzhou"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallNatFirewallSourceConfig(rand, map[string]string{
			"ids":       `["${alicloud_cloud_firewall_nat_firewall.default.id}_fake"]`,
			"region_no": `"cn-shenzhen"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallNatFirewallSourceConfig(rand, map[string]string{
			"ids":  `["${alicloud_cloud_firewall_nat_firewall.default.id}"]`,
			"lang": `"zh"`,

			"proxy_name": `"YQC-防火墙测试"`,

			"vpc_id": `"${alicloud_vpc.defaultikZ0gD.id}"`,

			"nat_gateway_id": `"${alicloud_nat_gateway.default2iRZpC.id}"`,

			"status": `"closed"`,

			"region_no": `"cn-hangzhou"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallNatFirewallSourceConfig(rand, map[string]string{
			"ids":  `["${alicloud_cloud_firewall_nat_firewall.default.id}_fake"]`,
			"lang": `"zh_fake"`,

			"proxy_name": `"YQC-防火墙测试_fake"`,

			"vpc_id": `"${alicloud_vpc.defaultikZ0gD.id}_fake"`,

			"nat_gateway_id": `"${alicloud_nat_gateway.default2iRZpC.id}_fake"`,

			"status": `"closed_fake"`,

			"region_no": `"cn-shenzhen"`,
		}),
	}

	CloudFirewallNatFirewallCheckInfo.dataSourceTestCheck(t, rand, idsConf, LangConf, ProxyNameConf, VpcIdConf, NatGatewayIdConf, StatusConf, RegionNoConf, allConf)
}

var existCloudFirewallNatFirewallMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"firewalls.#":                        "1",
		"firewalls.0.member_uid":             CHECKSET,
		"firewalls.0.nat_gateway_name":       CHECKSET,
		"firewalls.0.proxy_id":               CHECKSET,
		"firewalls.0.strict_mode":            CHECKSET,
		"firewalls.0.vpc_id":                 CHECKSET,
		"firewalls.0.proxy_name":             CHECKSET,
		"firewalls.0.nat_route_entry_list.#": CHECKSET,
		"firewalls.0.nat_gateway_id":         CHECKSET,
		"firewalls.0.ali_uid":                CHECKSET,
	}
}

var fakeCloudFirewallNatFirewallMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"firewalls.#": "0",
	}
}

var CloudFirewallNatFirewallCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cloud_firewall_nat_firewalls.default",
	existMapFunc: existCloudFirewallNatFirewallMapFunc,
	fakeMapFunc:  fakeCloudFirewallNatFirewallMapFunc,
}

func testAccCheckAlicloudCloudFirewallNatFirewallSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCloudFirewallNatFirewall%d"
}
data "alicloud_regions" "current" {
  current = true
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultikZ0gD" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultp4O7qi" {
  vpc_id       = alicloud_vpc.defaultikZ0gD.id
  cidr_block   = "172.16.6.0/24"
  vswitch_name = var.name
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_nat_gateway" "default2iRZpC" {
  eip_bind_mode    = "MULTI_BINDED"
  vpc_id           = alicloud_vpc.defaultikZ0gD.id
  nat_gateway_name = var.name
  payment_type     = "PayAsYouGo"
  vswitch_id       = alicloud_vswitch.defaultp4O7qi.id
  nat_type         = "Enhanced"
  network_type = "internet"
}

resource "alicloud_eip_address" "defaultyiRwgs" {
  address_name = var.name
}

resource "alicloud_eip_association" "defaults2MTuO" {
  instance_id   = alicloud_nat_gateway.default2iRZpC.id
  allocation_id = alicloud_eip_address.defaultyiRwgs.id
  mode          = "NAT"
  instance_type = "Nat"
}

resource "alicloud_snat_entry" "defaultAKE43g" {
  snat_ip           = alicloud_eip_address.defaultyiRwgs.ip_address
  snat_table_id     = alicloud_nat_gateway.default2iRZpC.snat_table_ids
  source_vswitch_id = alicloud_vswitch.defaultp4O7qi.id
}




resource "alicloud_cloud_firewall_nat_firewall" "default" {
  vswitch_cidr = "172.16.5.0/24"
  region_no = "${data.alicloud_regions.current.ids.0}"
  nat_gateway_id = "${alicloud_nat_gateway.default2iRZpC.id}"
  vswitch_id = "${alicloud_snat_entry.defaultAKE43g.source_vswitch_id}"
  vswitch_auto = "true"
  strict_mode = "0"
  vpc_id = "${alicloud_vpc.defaultikZ0gD.id}"
  nat_route_entry_list {
    nexthop_id = "${alicloud_nat_gateway.default2iRZpC.id}"
    destination_cidr = "0.0.0.0/0"
    nexthop_type = "NatGateway"
    route_table_id = "${alicloud_vpc.defaultikZ0gD.route_table_id}"
  }
  
  proxy_name = "YQC-防火墙测试"
  status = "closed"
  firewall_switch = "close"
}

data "alicloud_cloud_firewall_nat_firewalls" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
