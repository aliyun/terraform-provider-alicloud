package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAdbDbClustersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_adb_db_cluster.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"^test1234"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_adb_db_cluster.default.description}"`,
			"status":            `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_adb_db_cluster.default.description}"`,
			"status":            `"Creating"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_adb_db_cluster.default.description}"`,
			"tags": `{ 
						"Created" = "TF-update"
    					"For"     = "acceptance-test-update" 
					}`,
		}),
		fakeConfig: testAccCheckAlicloudAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_adb_db_cluster.default.description}"`,
			"tags": `{ 
						"Created" = "TF-update-fake"
    					"For"     = "acceptance-test-update-fake" 
					}`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbDbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_adb_db_cluster.default.description}"`,
			"status":            `"Running"`,
			"tags": `{ 
						"Created" = "TF-update"
    					"For"     = "acceptance-test-update" 
					}`,
		}),
		fakeConfig: testAccCheckAlicloudAdbDbClusterDataSourceConfig(rand, map[string]string{
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
			"clusters.0.id":                  CHECKSET,
			"clusters.0.description":         CHECKSET,
			"clusters.0.charge_type":         "PostPaid",
			"clusters.0.region_id":           CHECKSET,
			"clusters.0.expired":             "false",
			"clusters.0.lock_mode":           "Unlock",
			"clusters.0.create_time":         CHECKSET,
			"clusters.0.db_node_count":       "4",
			"clusters.0.db_node_class":       "C8",
			"clusters.0.db_node_storage":     "400",
			"clusters.0.db_cluster_category": "Cluster",
			"clusters.0.maintain_time":       "23:00Z-00:00Z",
			"clusters.0.security_ips.#":      "2",
		}
	}

	var fakeAdbClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"clusters.#":     CHECKSET,
			"ids.#":          CHECKSET,
			"descriptions.#": CHECKSET,
		}
	}

	var AdbClusterCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_adb_db_clusters.default",
		existMapFunc: existAdbClusterMapFunc,
		fakeMapFunc:  fakeAdbClusterMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithNoDefaultVpc(t)
		testAccPreCheckWithNoDefaultVswitch(t)
	}

	AdbClusterCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameConf, statusConf, tagsConf, allConf)
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

resource "alicloud_adb_db_cluster" "default" {
  db_cluster_category = "Cluster"
  db_cluster_class    = "C8"
  db_node_count       = "4"
  db_node_storage     = "400"
  mode                = "reserver"
  db_cluster_version  = "3.0"
  payment_type        = "PostPaid"
  vswitch_id          = data.alicloud_vswitches.default.ids[0]
  description         = "Test new adb again."
  maintain_time       = "23:00Z-00:00Z"
  tags = {
    Created = "TF-update"
    For     = "acceptance-test-update"
  }
  resource_group_id = "rg-aek2s7ylxx66kca"
  security_ips      = ["10.168.1.12", "10.168.1.11"]
}

data "alicloud_adb_db_clusters" "default" {	
	enable_details = true
	%s
}
`, AdbCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
