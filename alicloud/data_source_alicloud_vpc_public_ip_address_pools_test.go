package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVpcPublicIpAddressPoolsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPublicIpAddressPoolsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_public_ip_address_pool.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPublicIpAddressPoolsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_public_ip_address_pool.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPublicIpAddressPoolsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vpc_public_ip_address_pool.default.public_ip_address_pool_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPublicIpAddressPoolsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vpc_public_ip_address_pool.default.public_ip_address_pool_name}_fake"`,
		}),
	}
	publicIpAddressPoolIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPublicIpAddressPoolsDataSourceName(rand, map[string]string{
			"public_ip_address_pool_ids": `["${alicloud_vpc_public_ip_address_pool.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPublicIpAddressPoolsDataSourceName(rand, map[string]string{
			"public_ip_address_pool_ids": `["${alicloud_vpc_public_ip_address_pool.default.id}_fake"]`,
		}),
	}
	publicIpAddressPoolNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPublicIpAddressPoolsDataSourceName(rand, map[string]string{
			"public_ip_address_pool_name": `"${alicloud_vpc_public_ip_address_pool.default.public_ip_address_pool_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPublicIpAddressPoolsDataSourceName(rand, map[string]string{
			"public_ip_address_pool_name": `"${alicloud_vpc_public_ip_address_pool.default.public_ip_address_pool_name}_fake"`,
		}),
	}
	ispConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPublicIpAddressPoolsDataSourceName(rand, map[string]string{
			"isp": `"BGP_PRO"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPublicIpAddressPoolsDataSourceName(rand, map[string]string{
			"isp": `"ChinaTelecom"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPublicIpAddressPoolsDataSourceName(rand, map[string]string{
			"status": `"Created"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPublicIpAddressPoolsDataSourceName(rand, map[string]string{
			"status": `"Deleting"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPublicIpAddressPoolsDataSourceName(rand, map[string]string{
			"ids":                         `["${alicloud_vpc_public_ip_address_pool.default.id}"]`,
			"name_regex":                  `"${alicloud_vpc_public_ip_address_pool.default.public_ip_address_pool_name}"`,
			"public_ip_address_pool_ids":  `["${alicloud_vpc_public_ip_address_pool.default.id}"]`,
			"public_ip_address_pool_name": `"${alicloud_vpc_public_ip_address_pool.default.public_ip_address_pool_name}"`,
			"isp":                         `"BGP_PRO"`,
			"status":                      `"Created"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPublicIpAddressPoolsDataSourceName(rand, map[string]string{
			"ids":                         `["${alicloud_vpc_public_ip_address_pool.default.id}_fake"]`,
			"name_regex":                  `"${alicloud_vpc_public_ip_address_pool.default.public_ip_address_pool_name}_fake"`,
			"public_ip_address_pool_ids":  `["${alicloud_vpc_public_ip_address_pool.default.id}_fake"]`,
			"public_ip_address_pool_name": `"${alicloud_vpc_public_ip_address_pool.default.public_ip_address_pool_name}_fake"`,
			"isp":                         `"ChinaTelecom"`,
			"status":                      `"Deleting"`,
		}),
	}
	var existAlicloudVpcPublicIpAddressPoolsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"names.#":                             "1",
			"pools.#":                             "1",
			"pools.0.id":                          CHECKSET,
			"pools.0.public_ip_address_pool_id":   CHECKSET,
			"pools.0.public_ip_address_pool_name": CHECKSET,
			"pools.0.isp":                         "BGP_PRO",
			"pools.0.description":                 CHECKSET,
			"pools.0.status":                      "Created",
			"pools.0.region_id":                   CHECKSET,
			"pools.0.user_type":                   CHECKSET,
			"pools.0.total_ip_num":                CHECKSET,
			"pools.0.used_ip_num":                 CHECKSET,
			"pools.0.create_time":                 CHECKSET,
			"pools.0.ip_address_remaining":        CHECKSET,
		}
	}
	var fakeAlicloudVpcPublicIpAddressPoolsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"pools.#": "0",
		}
	}
	var alicloudVpcPublicIpAddressPoolsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_vpc_public_ip_address_pools.default",
		existMapFunc: existAlicloudVpcPublicIpAddressPoolsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudVpcPublicIpAddressPoolsDataSourceNameMapFunc,
	}
	alicloudVpcPublicIpAddressPoolsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, publicIpAddressPoolIdsConf, publicIpAddressPoolNameConf, ispConf, statusConf, allConf)
}

func testAccCheckAlicloudVpcPublicIpAddressPoolsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccVpcPublicIpAddressPoolBisic-%d"
	}

	resource "alicloud_vpc_public_ip_address_pool" "default" {
  		public_ip_address_pool_name = var.name
  		isp                         = "BGP_PRO"
  		description                 = var.name
	}

	data "alicloud_vpc_public_ip_address_pools" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
