package vpcipam_test

import (
	"fmt"
	"strings"
	"testing"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudVpcIpamIpamPoolCidrDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	IpamPoolIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamPoolCidrSourceConfig(rand, map[string]string{
			"ipam_pool_id": `"${alicloud_vpc_ipam_ipam_pool_cidr.default.ipam_pool_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamPoolCidrSourceConfig(rand, map[string]string{
			"ipam_pool_id": `"${alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id}_fake"`,
		}),
	}
	CidrConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamPoolCidrSourceConfig(rand, map[string]string{
			"ipam_pool_id": `"${alicloud_vpc_ipam_ipam_pool_cidr.default.ipam_pool_id}"`,
			"cidr":         `"10.0.0.0/8"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamPoolCidrSourceConfig(rand, map[string]string{
			"ipam_pool_id": `"${alicloud_vpc_ipam_ipam_pool_cidr.default.ipam_pool_id}"`,
			"cidr":         `"10.0.0.0/8_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamPoolCidrSourceConfig(rand, map[string]string{
			"ipam_pool_id": `"${alicloud_vpc_ipam_ipam_pool_cidr.default.ipam_pool_id}"`,

			"cidr": `"10.0.0.0/8"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamPoolCidrSourceConfig(rand, map[string]string{
			"ipam_pool_id": `"${alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id}_fake"`,

			"cidr": `"10.0.0.0/8_fake"`,
		}),
	}

	VpcIpamIpamPoolCidrCheckInfo.DataSourceTestCheck(t, rand, IpamPoolIdConf, CidrConf, allConf)
}

var existVpcIpamIpamPoolCidrMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"cidrs.#":              "1",
		"cidrs.0.status":       acctest.CHECKSET,
		"cidrs.0.cidr":         acctest.CHECKSET,
		"cidrs.0.ipam_pool_id": acctest.CHECKSET,
	}
}

var fakeVpcIpamIpamPoolCidrMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"cidrs.#": "0",
	}
}

var VpcIpamIpamPoolCidrCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_vpc_ipam_ipam_pool_cidrs.default",
	ExistMapFunc: existVpcIpamIpamPoolCidrMapFunc,
	FakeMapFunc:  fakeVpcIpamIpamPoolCidrMapFunc,
}

func testAccCheckAlicloudVpcIpamIpamPoolCidrSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccVpcIpamIpamPoolCidr%d"
}
resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
}

resource "alicloud_vpc_ipam_ipam_pool" "defaultIpamPool" {
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id = alicloud_vpc_ipam_ipam.defaultIpam.region_id
  ip_version     = "IPv4"
}



resource "alicloud_vpc_ipam_ipam_pool_cidr" "default" {
  cidr         = "10.0.0.0/8"
  ipam_pool_id = alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id
}

data "alicloud_vpc_ipam_ipam_pool_cidrs" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
