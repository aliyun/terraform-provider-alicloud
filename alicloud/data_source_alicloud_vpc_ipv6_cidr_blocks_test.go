// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudVPCIpv6CidrBlockDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheckExist, testAccCheckEmpty := VpcIpv6CidrBlockCheckInfo.checkDataSourceAttr(rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou", "cn-beijing"})
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudVpcIpv6CidrBlockSourceConfig(rand, map[string]interface{}{
					"ids":         "split(\",\", alicloud_vpc_ipv6_cidr_block.default.id)",
					"vpc_id":      "alicloud_vpc.defaultVpc.id",
					"output_file": "\"./test_output_file\"",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExist(nil),
				),
			},
			{
				Config: testAccCheckAlicloudVpcIpv6CidrBlockSourceConfig(rand, map[string]interface{}{
					"ids":    "split(\",\", format(\"%s_fake\", alicloud_vpc_ipv6_cidr_block.default.id))",
					"vpc_id": "alicloud_vpc.defaultVpc.id",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEmpty(nil),
				),
			},
		},
	})
}

var existVpcIpv6CidrBlockMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"blocks.#":                 "1",
		"blocks.0.ipv6_cidr_block": CHECKSET,
	}
}

var fakeVpcIpv6CidrBlockMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"blocks.#": "0",
	}
}

var VpcIpv6CidrBlockCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_vpc_ipv6_cidr_blocks.default",
	existMapFunc: existVpcIpv6CidrBlockMapFunc,
	fakeMapFunc:  fakeVpcIpv6CidrBlockMapFunc,
}

func testAccCheckAlicloudVpcIpv6CidrBlockSourceConfig(rand int, attrMap map[string]interface{}) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, fmt.Sprintf("%s = %v", k, v))
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccVpcIpv6CidrBlock%d"
}
resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou", "cn-beijing"]
}

resource "alicloud_vpc_ipam_ipam_pool" "defaultIpv6Pool" {
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id = alicloud_vpc_ipam_ipam.defaultIpam.region_id
  ip_version     = "IPv6"
}

resource "alicloud_vpc_ipam_ipam_pool_cidr" "defaultIpv6PoolCidr" {
  ipam_pool_id = alicloud_vpc_ipam_ipam_pool.defaultIpv6Pool.id
  cidr         = "fd03:d00:a000::/48"
}

resource "alicloud_vpc" "defaultVpc" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = "test-ipv6-cidr-block"

  lifecycle {
    ignore_changes = [enable_ipv6, ipv6_cidr_block, ipv6_cidr_blocks]
  }
}



resource "alicloud_vpc_ipv6_cidr_block" "default" {
  ipv6_ipam_pool_id = alicloud_vpc_ipam_ipam_pool_cidr.defaultIpv6PoolCidr.ipam_pool_id
  vpc_id            = alicloud_vpc.defaultVpc.id
  ipv6_cidr_block   = "fd03:d00:a000::/60"
}

data "alicloud_vpc_ipv6_cidr_blocks" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
