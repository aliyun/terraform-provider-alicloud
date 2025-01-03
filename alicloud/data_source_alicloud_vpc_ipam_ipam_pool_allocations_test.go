package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVpcIpamIpamPoolAllocationDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpamIpamPoolAllocationSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_vpc_ipam_ipam_pool_allocation.default.id}"]`,
			"ipam_pool_id": `"${alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpamIpamPoolAllocationSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_vpc_ipam_ipam_pool_allocation.default.id}_fake"]`,
			"ipam_pool_id": `"${alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id}"`,
		}),
	}

	IpamPoolAllocationNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpamIpamPoolAllocationSourceConfig(rand, map[string]string{
			"ids":                       `["${alicloud_vpc_ipam_ipam_pool_allocation.default.id}"]`,
			"ipam_pool_allocation_name": `"${var.name}"`,
			"ipam_pool_id":              `"${alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpamIpamPoolAllocationSourceConfig(rand, map[string]string{
			"ids":                       `["${alicloud_vpc_ipam_ipam_pool_allocation.default.id}_fake"]`,
			"ipam_pool_allocation_name": `"${var.name}_fake"`,
			"ipam_pool_id":              `"${alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id}"`,
		}),
	}
	IpamPoolIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpamIpamPoolAllocationSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_vpc_ipam_ipam_pool_allocation.default.id}"]`,
			"ipam_pool_id": `"${alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpamIpamPoolAllocationSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_vpc_ipam_ipam_pool_allocation.default.id}_fake"]`,
			"ipam_pool_id": `"${alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id}_fake"`,
		}),
	}
	CidrConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpamIpamPoolAllocationSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_vpc_ipam_ipam_pool_allocation.default.id}"]`,
			"cidr":         `"10.0.0.0/20"`,
			"ipam_pool_id": `"${alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpamIpamPoolAllocationSourceConfig(rand, map[string]string{
			"ids":          `["${alicloud_vpc_ipam_ipam_pool_allocation.default.id}_fake"]`,
			"cidr":         `"10.0.0.0/20_fake"`,
			"ipam_pool_id": `"${alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcIpamIpamPoolAllocationSourceConfig(rand, map[string]string{
			"ids":                       `["${alicloud_vpc_ipam_ipam_pool_allocation.default.id}"]`,
			"ipam_pool_allocation_name": `"${var.name}"`,

			"ipam_pool_id": `"${alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id}"`,

			"cidr": `"10.0.0.0/20"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcIpamIpamPoolAllocationSourceConfig(rand, map[string]string{
			"ids":                       `["${alicloud_vpc_ipam_ipam_pool_allocation.default.id}_fake"]`,
			"ipam_pool_allocation_name": `"${var.name}_fake"`,

			"ipam_pool_id": `"${alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id}_fake"`,

			"cidr": `"10.0.0.0/20_fake"`,
		}),
	}

	VpcIpamIpamPoolAllocationCheckInfo.dataSourceTestCheck(t, rand, idsConf, IpamPoolAllocationNameConf, IpamPoolIdConf, CidrConf, allConf)
}

var existVpcIpamIpamPoolAllocationMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"allocations.#":        "1",
		"allocations.0.status": CHECKSET,
		"allocations.0.ipam_pool_allocation_description": CHECKSET,
		"allocations.0.source_cidr":                      CHECKSET,
		"allocations.0.ipam_pool_id":                     CHECKSET,
		"allocations.0.create_time":                      CHECKSET,
		"allocations.0.resource_type":                    CHECKSET,
		"allocations.0.resource_owner_id":                CHECKSET,
		"allocations.0.ipam_pool_allocation_name":        CHECKSET,
		"allocations.0.total_count":                      CHECKSET,
		"allocations.0.cidr":                             CHECKSET,
		"allocations.0.ipam_pool_allocation_id":          CHECKSET,
		"allocations.0.region_id":                        CHECKSET,
	}
}

var fakeVpcIpamIpamPoolAllocationMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"allocations.#": "0",
	}
}

var VpcIpamIpamPoolAllocationCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_vpc_ipam_ipam_pool_allocations.default",
	existMapFunc: existVpcIpamIpamPoolAllocationMapFunc,
	fakeMapFunc:  fakeVpcIpamIpamPoolAllocationMapFunc,
}

func testAccCheckAlicloudVpcIpamIpamPoolAllocationSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccVpcIpamIpamPoolAllocation%d"
}
resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
}

resource "alicloud_vpc_ipam_ipam_pool" "defaultIpamPool" {
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id = "cn-hangzhou"
}

resource "alicloud_vpc_ipam_ipam_pool_cidr" "defaultIpamPoolCidr" {
  cidr         = "10.0.0.0/8"
  ipam_pool_id = alicloud_vpc_ipam_ipam_pool.defaultIpamPool.id
}



resource "alicloud_vpc_ipam_ipam_pool_allocation" "default" {
  ipam_pool_allocation_description = "init alloc desc"
  ipam_pool_allocation_name        = var.name
  cidr                             = "10.0.0.0/20"
  ipam_pool_id                     = alicloud_vpc_ipam_ipam_pool_cidr.defaultIpamPoolCidr.ipam_pool_id
}

data "alicloud_vpc_ipam_ipam_pool_allocations" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
