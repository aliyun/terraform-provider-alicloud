package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudNatGatewaysDataSourceBasic(t *testing.T) {
	rand := acctest.RandInt()
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNatGatewaysDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_nat_gateway.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlicloudNatGatewaysDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alicloud_nat_gateway.default.id}_fake" ]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNatGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_nat_gateway.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudNatGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_nat_gateway.default.name}_fake"`,
		}),
	}
	statusRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNatGatewaysDataSourceConfig(rand, map[string]string{
			"ids":    `[ "${alicloud_nat_gateway.default.id}" ]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudNatGatewaysDataSourceConfig(rand, map[string]string{
			"ids":    `[ "${alicloud_nat_gateway.default.id}" ]`,
			"status": `"Creating"`,
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNatGatewaysDataSourceConfig(rand, map[string]string{
			"vpc_id": `"${alicloud_vpc.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudNatGatewaysDataSourceConfig(rand, map[string]string{
			"vpc_id": `"${alicloud_vpc.default.id}_fake"`,
		}),
	}

	pagingConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNatGatewaysDataSourceConfig(rand, map[string]string{
			"vpc_id":      `"${alicloud_vpc.default.id}"`,
			"page_number": `1`,
		}),
		fakeConfig: testAccCheckAlicloudNatGatewaysDataSourceConfig(rand, map[string]string{
			"vpc_id":      `"${alicloud_vpc.default.id}"`,
			"page_number": `2`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNatGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex":  `"${alicloud_nat_gateway.default.nat_gateway_name}"`,
			"vpc_id":      `"${alicloud_vpc.default.id}"`,
			"ids":         `[ "${alicloud_nat_gateway.default.id}" ]`,
			"status":      `"Available"`,
			"page_number": `1`,
		}),
		fakeConfig: testAccCheckAlicloudNatGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex":  `"${alicloud_nat_gateway.default.nat_gateway_name}"`,
			"ids":         `[ "${alicloud_nat_gateway.default.id}" ]`,
			"vpc_id":      `"${alicloud_vpc.default.id}_fake"`,
			"status":      `"Creating"`,
			"page_number": `2`,
		}),
	}

	natGatewaysCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, statusRegexConf, IdsConf, vpcIdConf, pagingConf, allConf)
}

func testAccCheckAlicloudNatGatewaysDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccNatGatewaysDatasource%d"
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
    description = "${var.name}_decription"
	nat_type = "Enhanced"
	vswitch_id = alicloud_vswitch.default.id
    network_type = "internet"
}

data "alicloud_nat_gateways" "default" {
	%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existNatGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#":                      "1",
		"ids.#":                           "1",
		"names.#":                         "1",
		"total_count":                     CHECKSET,
		"gateways.0.id":                   CHECKSET,
		"gateways.0.business_status":      CHECKSET,
		"gateways.0.deletion_protection":  CHECKSET,
		"gateways.0.ecs_metric_enabled":   CHECKSET,
		"gateways.0.expired_time":         "",
		"gateways.0.internet_charge_type": CHECKSET,
		"gateways.0.nat_gateway_id":       CHECKSET,
		"gateways.0.nat_type":             CHECKSET,
		"gateways.0.payment_type":         CHECKSET,
		"gateways.0.resource_group_id":    CHECKSET,
		"gateways.0.vswitch_id":           CHECKSET,
		"gateways.0.specification":        "",
		"gateways.0.status":               "Available",
		"gateways.0.forward_table_ids.#":  CHECKSET,
		"gateways.0.snat_table_ids.#":     CHECKSET,
		"gateways.0.vpc_id":               CHECKSET,
		"gateways.0.ip_lists":             NOSET,
		"gateways.0.nat_gateway_name":     fmt.Sprintf("tf-testAccNatGatewaysDatasource%d", rand),
		"gateways.0.description":          fmt.Sprintf("tf-testAccNatGatewaysDatasource%d_decription", rand),
		"gateways.0.network_type":         "internet",
	}
}

var fakeNatGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#": "0",
		"ids.#":      "0",
		"names.#":    "0",
	}
}

var natGatewaysCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_nat_gateways.default",
	existMapFunc: existNatGatewaysMapFunc,
	fakeMapFunc:  fakeNatGatewaysMapFunc,
}
