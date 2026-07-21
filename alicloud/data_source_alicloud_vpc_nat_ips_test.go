package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPCNatIpsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 1000)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcNatIpDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip.default.nat_gateway_id}"`,
			"ids":            `["${alicloud_vpc_nat_ip.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcNatIpDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip.default.nat_gateway_id}"`,
			"ids":            `["${alicloud_vpc_nat_ip.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcNatIpDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip.default.nat_gateway_id}"`,
			"name_regex":     `"${alicloud_vpc_nat_ip.default.nat_ip_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcNatIpDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip.default.nat_gateway_id}"`,
			"name_regex":     `"${alicloud_vpc_nat_ip.default.nat_ip_name}_fake"`,
		}),
	}

	natIpCidrConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcNatIpDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip.default.nat_gateway_id}"`,
			"nat_ip_cidr":    `"${alicloud_vpc_nat_ip.default.nat_ip_cidr}"`,
			"ids":            `["${alicloud_vpc_nat_ip.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcNatIpDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip.default.nat_gateway_id}"`,
			"nat_ip_cidr":    `"8.8.8.8/24"`,
			"ids":            `["${alicloud_vpc_nat_ip.default.id}_fake"]`,
		}),
	}

	natIpNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcNatIpDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip.default.nat_gateway_id}"`,
			"ids":            `["${alicloud_vpc_nat_ip.default.id}"]`,
			"nat_ip_name":    `["${alicloud_vpc_nat_ip.default.nat_ip_name}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcNatIpDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip.default.nat_gateway_id}"`,
			"ids":            `["${alicloud_vpc_nat_ip.default.id}_fake"]`,
			"nat_ip_name":    `["${alicloud_vpc_nat_ip.default.nat_ip_name}_fake"]`,
		}),
	}
	natIpIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcNatIpDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip.default.nat_gateway_id}"`,
			"ids":            `["${alicloud_vpc_nat_ip.default.id}"]`,
			"nat_ip_ids":     `["${alicloud_vpc_nat_ip.default.nat_ip_id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcNatIpDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip.default.nat_gateway_id}"`,
			"ids":            `["${alicloud_vpc_nat_ip.default.id}_fake"]`,
			"nat_ip_ids":     `["${alicloud_vpc_nat_ip.default.nat_ip_id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcNatIpDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip.default.nat_gateway_id}"`,
			"ids":            `["${alicloud_vpc_nat_ip.default.id}"]`,
			"status":         `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcNatIpDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip.default.nat_gateway_id}"`,
			"ids":            `["${alicloud_vpc_nat_ip.default.id}"]`,
			"status":         `"Deleting"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcNatIpDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip.default.nat_gateway_id}"`,
			"ids":            `["${alicloud_vpc_nat_ip.default.id}"]`,
			"nat_ip_name":    `["${alicloud_vpc_nat_ip.default.nat_ip_name}"]`,
			"name_regex":     `"${alicloud_vpc_nat_ip.default.nat_ip_name}"`,
			"nat_ip_ids":     `["${alicloud_vpc_nat_ip.default.nat_ip_id}"]`,
			"nat_ip_cidr":    `"${alicloud_vpc_nat_ip.default.nat_ip_cidr}"`,
			"status":         `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcNatIpDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip.default.nat_gateway_id}"`,
			"ids":            `["${alicloud_vpc_nat_ip.default.id}_fake"]`,
			"nat_ip_name":    `["${alicloud_vpc_nat_ip.default.nat_ip_name}_fake"]`,
			"name_regex":     `"${alicloud_vpc_nat_ip.default.nat_ip_name}_fake"`,
			"nat_ip_ids":     `["${alicloud_vpc_nat_ip.default.nat_ip_id}_fake"]`,
			"nat_ip_cidr":    `"8.8.8.8/24"`,
			"status":         `"Deleting"`,
		}),
	}

	var existDataAlicloudVpcNatIpsSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":             "1",
			"names.#":           "1",
			"ips.#":             "1",
			"ips.0.nat_ip_name": fmt.Sprintf("tf-testaccvpcnatip%d", rand),
			"ips.0.status":      "Available",
		}
	}
	var fakeDataAlicloudVpcNatIpsSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"ips.#":   "0",
		}
	}
	var alicloudVpcNatIpCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_vpc_nat_ips.default",
		existMapFunc: existDataAlicloudVpcNatIpsSourceNameMapFunc,
		fakeMapFunc:  fakeDataAlicloudVpcNatIpsSourceNameMapFunc,
	}

	alicloudVpcNatIpCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, natIpCidrConf, natIpNameConf, natIpIdsConf, statusConf, allConf)
}
func testAccCheckAlicloudVpcNatIpDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testaccvpcnatip%d"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alicloud_vpc" "default" {
	vpc_name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
	vpc_id = alicloud_vpc.default.id
	cidr_block = "172.16.0.0/21"
	zone_id = data.alicloud_zones.default.zones.0.id
	vswitch_name = var.name
}

resource "alicloud_nat_gateway" "default" {
	vpc_id = alicloud_vpc.default.id
	nat_gateway_name = var.name
    description = "${var.name}_description"
	nat_type = "Enhanced"
	vswitch_id = alicloud_vswitch.default.id
	network_type = "intranet"
}

resource "alicloud_vpc_nat_ip_cidr" "default" {
	nat_ip_cidr = "192.168.0.0/16"
	nat_gateway_id = alicloud_nat_gateway.default.id
	nat_ip_cidr_description = var.name
	nat_ip_cidr_name = var.name
}

resource "alicloud_vpc_nat_ip" "default" {
	nat_ip = "192.168.0.37"
	nat_gateway_id = alicloud_nat_gateway.default.id
	nat_ip_description = var.name
	nat_ip_name = var.name
	nat_ip_cidr = alicloud_vpc_nat_ip_cidr.default.nat_ip_cidr
}

data "alicloud_vpc_nat_ips" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
