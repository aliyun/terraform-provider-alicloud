package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAdbClustersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_adb_cluster.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudAdbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"^test1234"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_adb_cluster.default.description}"`,
			"status":            `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudAdbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_adb_cluster.default.description}"`,
			"status":            `"Creating"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_adb_cluster.default.description}"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
		}),
		fakeConfig: testAccCheckAlicloudAdbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_adb_cluster.default.description}"`,
			"tags": `{ 
						"key1" = "value1_fake" 
						"key2" = "value2_fake" 
					}`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_adb_cluster.default.description}"`,
			"status":            `"Running"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
		}),
		fakeConfig: testAccCheckAlicloudAdbClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"^test1234"`,
			"status":            `"Creating"`,
			"tags": `{ 
						"key1" = "value1_fake" 
						"key2" = "value2_fake" 
					}`,
		}),
	}

	var existAdbClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"descriptions.#":             "1",
			"clusters.#":                 "1",
			"clusters.0.id":              CHECKSET,
			"clusters.0.description":     CHECKSET,
			"clusters.0.charge_type":     "PostPaid",
			"clusters.0.region_id":       CHECKSET,
			"clusters.0.expired":         "false",
			"clusters.0.lock_mode":       "Unlock",
			"clusters.0.create_time":     CHECKSET,
			"clusters.0.db_node_count":   "2",
			"clusters.0.db_node_class":   "C8",
			"clusters.0.db_node_storage": "200",
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
		resourceId:   "data.alicloud_adb_clusters.default",
		existMapFunc: existAdbClusterMapFunc,
		fakeMapFunc:  fakeAdbClusterMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithNoDefaultVpc(t)
		testAccPreCheckWithNoDefaultVswitch(t)
	}

	AdbClusterCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameConf, statusConf, tagsConf, allConf)
}

func testAccCheckAlicloudAdbClusterDataSourceConfig(rand int, attrMap map[string]string) string {
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

	resource "alicloud_adb_cluster" "default" {
		db_cluster_version      = "3.0"
        db_cluster_category     = "Cluster"
        db_node_class           = "C8"
        db_node_count           = 2
        db_node_storage         = 200
		pay_type                = "PostPaid"
		vswitch_id              = "${data.alicloud_vswitches.default.ids.0}"
		description             = "${var.name}"
		tags = {
			"key1" = "value1"
			"key2" = "value2"
		}
	}
	data "alicloud_adb_clusters" "default" {
	  %s
	}
`, AdbCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
