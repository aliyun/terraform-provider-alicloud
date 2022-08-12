package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVpcIpv4GatewaysDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpv4GatewaysDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_ipv4_gateway.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpv4GatewaysDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_ipv4_gateway.default.id}_fake"]`,
		}),
	}
	ipv4GatewayNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpv4GatewaysDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_vpc_ipv4_gateway.default.id}"]`,
			"ipv4_gateway_name": `"${alicloud_vpc_ipv4_gateway.default.ipv4_gateway_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpv4GatewaysDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_vpc_ipv4_gateway.default.id}"]`,
			"ipv4_gateway_name": `"${alicloud_vpc_ipv4_gateway.default.ipv4_gateway_name}_fake"`,
		}),
	}
	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpv4GatewaysDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_ipv4_gateway.default.id}"]`,
			"vpc_id": `"${alicloud_vpc_ipv4_gateway.default.vpc_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpv4GatewaysDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_ipv4_gateway.default.id}"]`,
			"vpc_id": `"${alicloud_vpc_ipv4_gateway.default.vpc_id}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpv4GatewaysDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vpc_ipv4_gateway.default.ipv4_gateway_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpv4GatewaysDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vpc_ipv4_gateway.default.ipv4_gateway_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpv4GatewaysDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_ipv4_gateway.default.id}"]`,
			"status": `"${alicloud_vpc_ipv4_gateway.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpv4GatewaysDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_ipv4_gateway.default.id}"]`,
			"status": `"Deleting"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpv4GatewaysDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_vpc_ipv4_gateway.default.id}"]`,
			"ipv4_gateway_name": `"${alicloud_vpc_ipv4_gateway.default.ipv4_gateway_name}"`,
			"name_regex":        `"${alicloud_vpc_ipv4_gateway.default.ipv4_gateway_name}"`,
			"status":            `"${alicloud_vpc_ipv4_gateway.default.status}"`,
			"vpc_id":            `"${alicloud_vpc_ipv4_gateway.default.vpc_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpv4GatewaysDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_vpc_ipv4_gateway.default.id}_fake"]`,
			"ipv4_gateway_name": `"${alicloud_vpc_ipv4_gateway.default.ipv4_gateway_name}_fake"`,
			"name_regex":        `"${alicloud_vpc_ipv4_gateway.default.ipv4_gateway_name}_fake"`,
			"status":            `"Deleting"`,
			"vpc_id":            `"${alicloud_vpc_ipv4_gateway.default.vpc_id}_fake"`,
		}),
	}
	var existAlicloudVpcIpv4GatewaysDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"names.#":                                "1",
			"gateways.#":                             "1",
			"gateways.0.ipv4_gateway_description":    fmt.Sprintf("tf-testAccIpv4Gateway-%d", rand),
			"gateways.0.ipv4_gateway_name":           fmt.Sprintf("tf-testAccIpv4Gateway-%d", rand),
			"gateways.0.vpc_id":                      CHECKSET,
			"gateways.0.create_time":                 CHECKSET,
			"gateways.0.enabled":                     CHECKSET,
			"gateways.0.id":                          CHECKSET,
			"gateways.0.ipv4_gateway_id":             CHECKSET,
			"gateways.0.ipv4_gateway_route_table_id": "",
			"gateways.0.status":                      CHECKSET,
		}
	}
	var fakeAlicloudVpcIpv4GatewaysDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudVpcIpv4GatewaysCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_vpc_ipv4_gateways.default",
		existMapFunc: existAlicloudVpcIpv4GatewaysDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudVpcIpv4GatewaysDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudVpcIpv4GatewaysCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, ipv4GatewayNameConf, vpcIdConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudVpcIpv4GatewaysDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccIpv4Gateway-%d"
}
data "alicloud_vpcs" "default" {
	name_regex = "default-NoDeleting"
}
resource "alicloud_vpc_ipv4_gateway" "default" {
	ipv4_gateway_description = var.name
	ipv4_gateway_name = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_vpc_ipv4_gateways" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
