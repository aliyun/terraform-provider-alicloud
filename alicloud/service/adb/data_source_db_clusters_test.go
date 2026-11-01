// Package adb_test.
package adb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudADBDbClustersDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	nameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_adb_db_cluster.default.description}"`,
		}),
		FakeConfig: testAccCheckAlicloudAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"^test1234"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_adb_db_cluster.default.description}"`,
			"status":            `"Running"`,
		}),
		FakeConfig: testAccCheckAlicloudAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_adb_db_cluster.default.description}"`,
			"status":            `"Creating"`,
		}),
	}
	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_adb_db_cluster.default.description}"`,
			"tags": `{ 
						"Created" = "TF-update"
    					"For"     = "acceptance-test-update" 
					}`,
		}),
		FakeConfig: testAccCheckAlicloudAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_adb_db_cluster.default.description}"`,
			"tags": `{ 
						"Created" = "TF-update-fake"
    					"For"     = "acceptance-test-update-fake" 
					}`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_adb_db_cluster.default.description}"`,
			"status":            `"Running"`,
			"tags": `{ 
						"Created" = "TF-update"
    					"For"     = "acceptance-test-update" 
					}`,
		}),
		FakeConfig: testAccCheckAlicloudAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"^test1234"`,
			"status":            `"Creating"`,
			"tags": `{ 
						"Created" = "TF-update-fake"
    					"For"     = "acceptance-test-update-fake" 
					}`,
		}),
	}

	var existAdbClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"descriptions.#":                 "1",
			"clusters.#":                     "1",
			"total_count":                    acctest.CHECKSET,
			"clusters.0.id":                  acctest.CHECKSET,
			"clusters.0.description":         acctest.CHECKSET,
			"clusters.0.payment_type":        "PayAsYouGo",
			"clusters.0.charge_type":         "PostPaid",
			"clusters.0.region_id":           acctest.CHECKSET,
			"clusters.0.expired":             acctest.CHECKSET,
			"clusters.0.lock_mode":           "Unlock",
			"clusters.0.create_time":         acctest.CHECKSET,
			"clusters.0.db_cluster_version":  "3.0",
			"clusters.0.db_node_class":       "E8",
			"clusters.0.db_node_count":       "1",
			"clusters.0.db_node_storage":     acctest.CHECKSET,
			"clusters.0.compute_resource":    "8Core32GB",
			"clusters.0.elastic_io_resource": "0",
			"clusters.0.zone_id":             acctest.CHECKSET,
			"clusters.0.db_cluster_category": "MixedStorage",
			"clusters.0.maintain_time":       "23:00Z-00:00Z",
			"clusters.0.security_ips.#":      "2",
			"clusters.0.mode":                "flexible",
			"clusters.0.kernel_version":      acctest.CHECKSET,
		}
	}

	var fakeAdbClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"clusters.#":     acctest.CHECKSET,
			"ids.#":          acctest.CHECKSET,
			"descriptions.#": acctest.CHECKSET,
		}
	}

	var AdbClusterCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_adb_db_clusters.default",
		ExistMapFunc: existAdbClusterMapFunc,
		FakeMapFunc:  fakeAdbClusterMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
	}
	AdbClusterCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameConf, statusConf, tagsConf, allConf)
}

func testAccCheckAlicloudAdbDbClusterDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
	%s
variable "creation" {	
	default = "ADB"
}

variable "name" {
	default = "tf-testAccADBConfig_%d"
}
data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}

resource "alicloud_adb_db_cluster" "default" {
  db_cluster_category = "MixedStorage"
  mode = "flexible"
  compute_resource = "8Core32GB"
  payment_type        = "PayAsYouGo"
  vswitch_id          = local.vswitch_id
  description         = var.name
  maintain_time       = "23:00Z-00:00Z"
  tags = {
    Created = "TF-update"
    For     = "acceptance-test-update"
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  security_ips      = ["10.168.1.12", "10.168.1.11"]
}

data "alicloud_adb_db_clusters" "default" {	
	enable_details = true
	%s
}
`, AdbCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
