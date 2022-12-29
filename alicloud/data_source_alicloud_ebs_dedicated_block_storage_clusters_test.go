package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func SkipTestAccAlicloudEbsDedicatedBlockStorageClusterDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEbsDedicatedBlockStorageClusterSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ebs_dedicated_block_storage_cluster.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEbsDedicatedBlockStorageClusterSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ebs_dedicated_block_storage_cluster.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEbsDedicatedBlockStorageClusterSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ebs_dedicated_block_storage_cluster.default.dedicated_block_storage_cluster_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEbsDedicatedBlockStorageClusterSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_ebs_dedicated_block_storage_cluster.default.dedicated_block_storage_cluster_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEbsDedicatedBlockStorageClusterSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ebs_dedicated_block_storage_cluster.default.id}"]`,
			"name_regex": `"${alicloud_ebs_dedicated_block_storage_cluster.default.dedicated_block_storage_cluster_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEbsDedicatedBlockStorageClusterSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_ebs_dedicated_block_storage_cluster.default.id}_fake"]`,
			"name_regex": `"${alicloud_ebs_dedicated_block_storage_cluster.default.dedicated_block_storage_cluster_name}_fake"`,
		}),
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.EbsDedicatedBlockStorageClusterRegions)
	}

	EbsDedicatedBlockStorageClusterCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

var existEbsDedicatedBlockStorageClusterMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"clusters.#":                                      "1",
		"clusters.0.id":                                   CHECKSET,
		"clusters.0.available_capacity":                   CHECKSET,
		"clusters.0.category":                             CHECKSET,
		"clusters.0.create_time":                          CHECKSET,
		"clusters.0.dedicated_block_storage_cluster_id":   CHECKSET,
		"clusters.0.dedicated_block_storage_cluster_name": CHECKSET,
		"clusters.0.delivery_capacity":                    CHECKSET,
		"clusters.0.description":                          fmt.Sprintf("tf-testAccCluster%d", rand),
		"clusters.0.expired_time":                         CHECKSET,
		"clusters.0.performance_level":                    CHECKSET,
		"clusters.0.resource_group_id":                    CHECKSET,
		"clusters.0.status":                               CHECKSET,
		"clusters.0.supported_category":                   CHECKSET,
		"clusters.0.total_capacity":                       "61440",
		"clusters.0.type":                                 "Premium",
		"clusters.0.used_capacity":                        CHECKSET,
		"clusters.0.zone_id":                              CHECKSET,
	}
}

var fakeEbsDedicatedBlockStorageClusterMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"clusters.#": "0",
	}
}

var EbsDedicatedBlockStorageClusterCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ebs_dedicated_block_storage_clusters.default",
	existMapFunc: existEbsDedicatedBlockStorageClusterMapFunc,
	fakeMapFunc:  fakeEbsDedicatedBlockStorageClusterMapFunc,
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
`, rand, defaultRegionToTest, strings.Join(pairs, "\n   "))
	return config
}
