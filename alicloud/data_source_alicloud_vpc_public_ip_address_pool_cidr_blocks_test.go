package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVpcPublicIpAddressPoolCidrBlocksDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.VPCPublicIpAddressPoolCidrBlockSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPublicIpAddressPoolCidrBlocksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_public_ip_address_pool_cidr_block.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPublicIpAddressPoolCidrBlocksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_public_ip_address_pool_cidr_block.default.id}_fake"]`,
		}),
	}
	cidrBlockConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPublicIpAddressPoolCidrBlocksDataSourceName(rand, map[string]string{
			"cidr_block": `"${alicloud_vpc_public_ip_address_pool_cidr_block.default.cidr_block}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPublicIpAddressPoolCidrBlocksDataSourceName(rand, map[string]string{
			"cidr_block": `"${alicloud_vpc_public_ip_address_pool_cidr_block.default.cidr_block}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPublicIpAddressPoolCidrBlocksDataSourceName(rand, map[string]string{
			"status": `"Created"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPublicIpAddressPoolCidrBlocksDataSourceName(rand, map[string]string{
			"status": `"Deleting"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPublicIpAddressPoolCidrBlocksDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_vpc_public_ip_address_pool_cidr_block.default.id}"]`,
			"cidr_block": `"${alicloud_vpc_public_ip_address_pool_cidr_block.default.cidr_block}"`,
			"status":     `"Created"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPublicIpAddressPoolCidrBlocksDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_vpc_public_ip_address_pool_cidr_block.default.id}_fake"]`,
			"cidr_block": `"${alicloud_vpc_public_ip_address_pool_cidr_block.default.cidr_block}_fake"`,
			"status":     `"Deleting"`,
		}),
	}
	var existAlicloudVpcPublicIpAddressPoolCidrBlocksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                              "1",
			"blocks.#":                           "1",
			"blocks.0.id":                        CHECKSET,
			"blocks.0.public_ip_address_pool_id": CHECKSET,
			"blocks.0.cidr_block":                "47.118.126.0/25",
			"blocks.0.status":                    "Created",
			"blocks.0.used_ip_num":               CHECKSET,
			"blocks.0.total_ip_num":              CHECKSET,
			"blocks.0.create_time":               CHECKSET,
		}
	}
	var fakeAlicloudVpcPublicIpAddressPoolCidrBlocksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"blocks.#": "0",
		}
	}
	var alicloudVpcPublicIpAddressPoolCidrBlocksCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_vpc_public_ip_address_pool_cidr_blocks.default",
		existMapFunc: existAlicloudVpcPublicIpAddressPoolCidrBlocksDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudVpcPublicIpAddressPoolCidrBlocksDataSourceNameMapFunc,
	}
	alicloudVpcPublicIpAddressPoolCidrBlocksCheckInfo.dataSourceTestCheck(t, rand, idsConf, cidrBlockConf, statusConf, allConf)
}

func testAccCheckAlicloudVpcPublicIpAddressPoolCidrBlocksDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccVpcPublicIpAddressPoolCidrBlockBisic-%d"
	}

	resource "alicloud_vpc_public_ip_address_pool" "default" {
  		public_ip_address_pool_name = var.name
  		isp                         = "BGP"
  		description                 = var.name
	}

	resource "alicloud_vpc_public_ip_address_pool_cidr_block" "default" {
  		public_ip_address_pool_id = alicloud_vpc_public_ip_address_pool.default.id

		# Only users who have the required permissions can use the IP address pool feature. To apply for the required permissions, please submit a ticket.
  		cidr_block                         = "47.118.126.0/25"
	}

	data "alicloud_vpc_public_ip_address_pool_cidr_blocks" "default" {
		public_ip_address_pool_id = alicloud_vpc_public_ip_address_pool_cidr_block.default.public_ip_address_pool_id
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
