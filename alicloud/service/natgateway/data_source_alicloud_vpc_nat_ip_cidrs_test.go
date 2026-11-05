package natgateway_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudVPCNatIpCidrsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 1000)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcNatIpCidrDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip_cidr.default.nat_gateway_id}"`,
			"ids":            `["${alicloud_vpc_nat_ip_cidr.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudVpcNatIpCidrDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip_cidr.default.nat_gateway_id}"`,
			"ids":            `["${alicloud_vpc_nat_ip_cidr.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcNatIpCidrDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip_cidr.default.nat_gateway_id}"`,
			"name_regex":     `"${alicloud_vpc_nat_ip_cidr.default.nat_ip_cidr_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcNatIpCidrDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip_cidr.default.nat_gateway_id}"`,
			"name_regex":     `"${alicloud_vpc_nat_ip_cidr.default.nat_ip_cidr_name}_fake"`,
		}),
	}

	natIpCidrConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcNatIpCidrDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip_cidr.default.nat_gateway_id}"`,
			"nat_ip_cidrs":   `["${alicloud_vpc_nat_ip_cidr.default.nat_ip_cidr}"]`,
		}),
		FakeConfig: testAccCheckAlicloudVpcNatIpCidrDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip_cidr.default.nat_gateway_id}"`,
			"nat_ip_cidrs":   `["8.8.8.8/24"]`,
		}),
	}

	natIpCidrNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcNatIpCidrDataSourceName(rand, map[string]string{
			"nat_gateway_id":   `"${alicloud_vpc_nat_ip_cidr.default.nat_gateway_id}"`,
			"nat_ip_cidr_name": `["${alicloud_vpc_nat_ip_cidr.default.nat_ip_cidr_name}"]`,
		}),
		FakeConfig: testAccCheckAlicloudVpcNatIpCidrDataSourceName(rand, map[string]string{
			"nat_gateway_id":   `"${alicloud_vpc_nat_ip_cidr.default.nat_gateway_id}"`,
			"nat_ip_cidr_name": `["${alicloud_vpc_nat_ip_cidr.default.nat_ip_cidr_name}_fake"]`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcNatIpCidrDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip_cidr.default.nat_gateway_id}"`,
			"ids":            `["${alicloud_vpc_nat_ip_cidr.default.id}"]`,
			"status":         `"Available"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcNatIpCidrDataSourceName(rand, map[string]string{
			"nat_gateway_id": `"${alicloud_vpc_nat_ip_cidr.default.nat_gateway_id}"`,
			"ids":            `["${alicloud_vpc_nat_ip_cidr.default.id}_fake"]`,
			"status":         `"Available"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcNatIpCidrDataSourceName(rand, map[string]string{
			"nat_gateway_id":   `"${alicloud_vpc_nat_ip_cidr.default.nat_gateway_id}"`,
			"ids":              `["${alicloud_vpc_nat_ip_cidr.default.id}"]`,
			"name_regex":       `"${alicloud_vpc_nat_ip_cidr.default.nat_ip_cidr_name}"`,
			"nat_ip_cidrs":     `["${alicloud_vpc_nat_ip_cidr.default.nat_ip_cidr}"]`,
			"nat_ip_cidr_name": `["${alicloud_vpc_nat_ip_cidr.default.nat_ip_cidr_name}"]`,
			"status":           `"Available"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcNatIpCidrDataSourceName(rand, map[string]string{
			"nat_gateway_id":   `"${alicloud_vpc_nat_ip_cidr.default.nat_gateway_id}"`,
			"ids":              `["${alicloud_vpc_nat_ip_cidr.default.id}_fake"]`,
			"name_regex":       `"${alicloud_vpc_nat_ip_cidr.default.nat_ip_cidr_name}_fake"`,
			"nat_ip_cidrs":     `["${alicloud_vpc_nat_ip_cidr.default.nat_ip_cidr}"]`,
			"nat_ip_cidr_name": `["${alicloud_vpc_nat_ip_cidr.default.nat_ip_cidr_name}"]`,
			"status":           `"Available"`,
		}),
	}

	var existDataAlicloudVpcNatIpCidrsSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"names.#":                  "1",
			"cidrs.#":                  "1",
			"cidrs.0.nat_ip_cidr_name": fmt.Sprintf("tf-testaccvpcnatipcidr%d", rand),
			"cidrs.0.status":           "Available",
		}
	}
	var fakeDataAlicloudVpcNatIpCidrsSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"cidrs.#": "0",
		}
	}
	var alicloudVpcNatIpCidrCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_vpc_nat_ip_cidrs.default",
		ExistMapFunc: existDataAlicloudVpcNatIpCidrsSourceNameMapFunc,
		FakeMapFunc:  fakeDataAlicloudVpcNatIpCidrsSourceNameMapFunc,
	}

	alicloudVpcNatIpCidrCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, natGatewayDataSourcePreCheck(t), idsConf, nameRegexConf, natIpCidrConf, natIpCidrNameConf, statusConf, allConf)
}
func testAccCheckAlicloudVpcNatIpCidrDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testaccvpcnatipcidr%d"
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
	internet_charge_type = "PayByLcu"
	nat_gateway_name = var.name
    description = "${var.name}_description"
	nat_type = "Enhanced"
	vswitch_id = alicloud_vswitch.default.id
	network_type = "intranet"
}

resource "alicloud_vpc_nat_ip_cidr" "default" {
	nat_ip_cidr = "192.168.0.0/16"
	nat_gateway_id =  alicloud_nat_gateway.default.id
	nat_ip_cidr_description = var.name
	nat_ip_cidr_name = var.name
}

data "alicloud_vpc_nat_ip_cidrs" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
