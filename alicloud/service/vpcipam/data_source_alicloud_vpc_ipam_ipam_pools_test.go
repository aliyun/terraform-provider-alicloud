package vpcipam_test

import (
	"fmt"
	"strings"
	"testing"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudVpcIpamIpamPoolDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamPoolSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpc_ipam_ipam_pool.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamPoolSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_vpc_ipam_ipam_pool.default.id}_fake"]`,
		}),
	}

	PoolRegionIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamPoolSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_vpc_ipam_ipam_pool.default.id}"]`,
			"pool_region_id": `"${alicloud_vpc_ipam_ipam_pool.parentIpamPool.pool_region_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamPoolSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_vpc_ipam_ipam_pool.default.id}_fake"]`,
			"pool_region_id": `"${alicloud_vpc_ipam_ipam_pool.parentIpamPool.pool_region_id}_fake"`,
		}),
	}
	ResourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamPoolSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_vpc_ipam_ipam_pool.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamPoolSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_vpc_ipam_ipam_pool.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}
	SourceIpamPoolIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamPoolSourceConfig(rand, map[string]string{
			"ids":                 `["${alicloud_vpc_ipam_ipam_pool.default.id}"]`,
			"source_ipam_pool_id": `"${alicloud_vpc_ipam_ipam_pool.parentIpamPool.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamPoolSourceConfig(rand, map[string]string{
			"ids":                 `["${alicloud_vpc_ipam_ipam_pool.default.id}_fake"]`,
			"source_ipam_pool_id": `"${alicloud_vpc_ipam_ipam_pool.parentIpamPool.id}_fake"`,
		}),
	}
	IpamScopeIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamPoolSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_vpc_ipam_ipam_pool.default.id}"]`,
			"ipam_scope_id": `"${alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamPoolSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_vpc_ipam_ipam_pool.default.id}_fake"]`,
			"ipam_scope_id": `"${alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id}_fake"`,
		}),
	}
	IpamPoolNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamPoolSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_vpc_ipam_ipam_pool.default.id}"]`,
			"ipam_pool_name": `"${var.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamPoolSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_vpc_ipam_ipam_pool.default.id}_fake"]`,
			"ipam_pool_name": `"${var.name}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudVpcIpamIpamPoolSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_vpc_ipam_ipam_pool.default.id}"]`,
			"pool_region_id": `"${alicloud_vpc_ipam_ipam_pool.parentIpamPool.pool_region_id}"`,

			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,

			"source_ipam_pool_id": `"${alicloud_vpc_ipam_ipam_pool.parentIpamPool.id}"`,

			"ipam_scope_id": `"${alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id}"`,

			"ipam_pool_name": `"${var.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudVpcIpamIpamPoolSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_vpc_ipam_ipam_pool.default.id}_fake"]`,
			"pool_region_id": `"${alicloud_vpc_ipam_ipam_pool.parentIpamPool.pool_region_id}_fake"`,

			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,

			"source_ipam_pool_id": `"${alicloud_vpc_ipam_ipam_pool.parentIpamPool.id}_fake"`,

			"ipam_scope_id": `"${alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id}_fake"`,

			"ipam_pool_name": `"${var.name}_fake"`,
		}),
	}

	VpcIpamIpamPoolCheckInfo.DataSourceTestCheck(t, rand, idsConf, PoolRegionIdConf, ResourceGroupIdConf, SourceIpamPoolIdConf, IpamScopeIdConf, IpamPoolNameConf, allConf)
}

var existVpcIpamIpamPoolMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"pools.#":                              "1",
		"pools.0.status":                       acctest.CHECKSET,
		"pools.0.resource_group_id":            acctest.CHECKSET,
		"pools.0.ipam_pool_id":                 acctest.CHECKSET,
		"pools.0.ipam_pool_name":               acctest.CHECKSET,
		"pools.0.source_ipam_pool_id":          acctest.CHECKSET,
		"pools.0.ip_version":                   acctest.CHECKSET,
		"pools.0.create_time":                  acctest.CHECKSET,
		"pools.0.ipam_id":                      acctest.CHECKSET,
		"pools.0.allocation_default_cidr_mask": acctest.CHECKSET,
		"pools.0.allocation_min_cidr_mask":     acctest.CHECKSET,
		"pools.0.ipam_scope_id":                acctest.CHECKSET,
		"pools.0.ipam_pool_description":        acctest.CHECKSET,
		"pools.0.pool_region_id":               acctest.CHECKSET,
		"pools.0.pool_depth":                   acctest.CHECKSET,
		"pools.0.has_sub_pool":                 acctest.CHECKSET,
		"pools.0.auto_import":                  acctest.CHECKSET,
		"pools.0.tags.%":                       acctest.CHECKSET,
		"pools.0.allocation_max_cidr_mask":     acctest.CHECKSET,
	}
}

var fakeVpcIpamIpamPoolMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"pools.#": "0",
	}
}

var VpcIpamIpamPoolCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_vpc_ipam_ipam_pools.default",
	ExistMapFunc: existVpcIpamIpamPoolMapFunc,
	FakeMapFunc:  fakeVpcIpamIpamPoolMapFunc,
}

func testAccCheckAlicloudVpcIpamIpamPoolSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccVpcIpamIpamPool%d"
}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
}

resource "alicloud_vpc_ipam_ipam_pool" "parentIpamPool" {
  ipam_scope_id  = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id = "cn-hangzhou"
}



resource "alicloud_vpc_ipam_ipam_pool" "default" {
  ipam_scope_id       = alicloud_vpc_ipam_ipam.defaultIpam.private_default_scope_id
  pool_region_id      = alicloud_vpc_ipam_ipam_pool.parentIpamPool.pool_region_id
  ipam_pool_name      = var.name
  source_ipam_pool_id = alicloud_vpc_ipam_ipam_pool.parentIpamPool.id
  ip_version          = "IPv4"
  ipam_pool_description = var.name
}

data "alicloud_vpc_ipam_ipam_pools" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
