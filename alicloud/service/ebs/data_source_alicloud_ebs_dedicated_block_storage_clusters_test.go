package ebs_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func SkipTestAccAlicloudEbsDedicatedBlockStorageClusterDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEbsDedicatedBlockStorageClusterSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ebs_dedicated_block_storage_cluster.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudEbsDedicatedBlockStorageClusterSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ebs_dedicated_block_storage_cluster.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEbsDedicatedBlockStorageClusterSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ebs_dedicated_block_storage_cluster.default.dedicated_block_storage_cluster_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudEbsDedicatedBlockStorageClusterSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ebs_dedicated_block_storage_cluster.default.dedicated_block_storage_cluster_name}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEbsDedicatedBlockStorageClusterSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ebs_dedicated_block_storage_cluster.default.id}"]`,
			"name_regex": `"${alicloud_ebs_dedicated_block_storage_cluster.default.dedicated_block_storage_cluster_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudEbsDedicatedBlockStorageClusterSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ebs_dedicated_block_storage_cluster.default.id}_fake"]`,
			"name_regex": `"${alicloud_ebs_dedicated_block_storage_cluster.default.dedicated_block_storage_cluster_name}_fake"`,
		}),
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.EbsDedicatedBlockStorageClusterRegions)
	}

	EbsDedicatedBlockStorageClusterCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

var existEbsDedicatedBlockStorageClusterMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"clusters.#":                                      "1",
		"clusters.0.id":                                   acctest.CHECKSET,
		"clusters.0.available_capacity":                   acctest.CHECKSET,
		"clusters.0.category":                             acctest.CHECKSET,
		"clusters.0.create_time":                          acctest.CHECKSET,
		"clusters.0.dedicated_block_storage_cluster_id":   acctest.CHECKSET,
		"clusters.0.dedicated_block_storage_cluster_name": acctest.CHECKSET,
		"clusters.0.delivery_capacity":                    acctest.CHECKSET,
		"clusters.0.description":                          fmt.Sprintf("tf-testAccCluster%d", rand),
		"clusters.0.expired_time":                         acctest.CHECKSET,
		"clusters.0.performance_level":                    acctest.CHECKSET,
		"clusters.0.resource_group_id":                    acctest.CHECKSET,
		"clusters.0.status":                               acctest.CHECKSET,
		"clusters.0.supported_category":                   acctest.CHECKSET,
		"clusters.0.total_capacity":                       "61440",
		"clusters.0.type":                                 "Premium",
		"clusters.0.used_capacity":                        acctest.CHECKSET,
		"clusters.0.zone_id":                              acctest.CHECKSET,
	}
}

var fakeEbsDedicatedBlockStorageClusterMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"clusters.#": "0",
	}
}

var EbsDedicatedBlockStorageClusterCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_ebs_dedicated_block_storage_clusters.default",
	ExistMapFunc: existEbsDedicatedBlockStorageClusterMapFunc,
	FakeMapFunc:  fakeEbsDedicatedBlockStorageClusterMapFunc,
}

func testAccCheckAlicloudEbsDedicatedBlockStorageClusterSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCluster%d"
}

variable "region" {
  default = "%s"
}

data "alicloud_ebs_regions" "default"{
  region_id = var.region
}

resource "alicloud_ebs_dedicated_block_storage_cluster" "default" {
  type                                 = "Premium"
  zone_id                              = data.alicloud_ebs_regions.default.regions[0].zones[0].zone_id
  dedicated_block_storage_cluster_name = var.name
  total_capacity                       = 61440
  description                          = var.name
}

data "alicloud_ebs_dedicated_block_storage_clusters" "default" {
%s
}
`, rand, acctest.DefaultRegionToTest, strings.Join(pairs, "\n   "))
	return config
}
