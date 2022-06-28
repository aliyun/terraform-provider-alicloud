package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPCDhcpOptionsSetsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcDhcpOptionsSetsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_dhcp_options_set.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcDhcpOptionsSetsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_dhcp_options_set.default.id}_fake"]`,
		}),
	}
	dhcpOptionsSetNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcDhcpOptionsSetsDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_vpc_dhcp_options_set.default.id}"]`,
			"dhcp_options_set_name": `"${alicloud_vpc_dhcp_options_set.default.dhcp_options_set_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcDhcpOptionsSetsDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_vpc_dhcp_options_set.default.id}"]`,
			"dhcp_options_set_name": `"${alicloud_vpc_dhcp_options_set.default.dhcp_options_set_name}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcDhcpOptionsSetsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_vpc_dhcp_options_set.default.id}"]`,
			"name_regex": `"${alicloud_vpc_dhcp_options_set.default.dhcp_options_set_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcDhcpOptionsSetsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_vpc_dhcp_options_set.default.id}"]`,
			"name_regex": `"${alicloud_vpc_dhcp_options_set.default.dhcp_options_set_name}_fake"`,
		}),
	}
	domainNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcDhcpOptionsSetsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_vpc_dhcp_options_set.default.id}"]`,
			"domain_name": `"${alicloud_vpc_dhcp_options_set.default.domain_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcDhcpOptionsSetsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_vpc_dhcp_options_set.default.id}"]`,
			"domain_name": `"${alicloud_vpc_dhcp_options_set.default.domain_name}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcDhcpOptionsSetsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_dhcp_options_set.default.id}"]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcDhcpOptionsSetsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_vpc_dhcp_options_set.default.id}"]`,
			"status": `"Pending"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcDhcpOptionsSetsDataSourceName(rand, map[string]string{
			"dhcp_options_set_name": `"${alicloud_vpc_dhcp_options_set.default.dhcp_options_set_name}"`,
			"ids":                   `["${alicloud_vpc_dhcp_options_set.default.id}"]`,
			"name_regex":            `"${alicloud_vpc_dhcp_options_set.default.dhcp_options_set_name}"`,
			"domain_name":           `"${alicloud_vpc_dhcp_options_set.default.domain_name}"`,
			"status":                `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcDhcpOptionsSetsDataSourceName(rand, map[string]string{
			"dhcp_options_set_name": `"${alicloud_vpc_dhcp_options_set.default.dhcp_options_set_name}_fake"`,
			"ids":                   `["${alicloud_vpc_dhcp_options_set.default.id}"]`,
			"name_regex":            `"${alicloud_vpc_dhcp_options_set.default.dhcp_options_set_name}_fake"`,
			"domain_name":           `"${alicloud_vpc_dhcp_options_set.default.domain_name}_fake"`,
			"status":                `"Pending"`,
		}),
	}
	var existAlicloudVpcDhcpOptionsSetsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"sets.#":                              "1",
			"sets.0.dhcp_options_set_description": fmt.Sprintf("tf-testAccVpcDhcpOptionsSets-%d", rand),
			"sets.0.dhcp_options_set_name":        fmt.Sprintf("tf-testAccVpcDhcpOptionsSets-%d", rand),
		}
	}
	var fakeAlicloudVpcDhcpOptionsSetsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudVpcDhcpOptionsSetsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_vpc_dhcp_options_sets.default",
		existMapFunc: existAlicloudVpcDhcpOptionsSetsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudVpcDhcpOptionsSetsDataSourceNameMapFunc,
	}
	alicloudVpcDhcpOptionsSetsCheckInfo.dataSourceTestCheck(t, rand, idsConf, dhcpOptionsSetNameConf, nameRegexConf, domainNameConf, statusConf, allConf)
}
func testAccCheckAlicloudVpcDhcpOptionsSetsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccVpcDhcpOptionsSets-%d"
}

resource "alicloud_vpc_dhcp_options_set" "default" {
	dhcp_options_set_name = var.name
	dhcp_options_set_description = var.name
	domain_name = "example.com"
	domain_name_servers = "100.100.2.136"
}

data "alicloud_vpc_dhcp_options_sets" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
