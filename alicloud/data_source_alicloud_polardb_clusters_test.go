package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudPolarDBClustersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_polardb_cluster.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"^test1234"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_polardb_cluster.default.description}"`,
			"status":            `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_polardb_cluster.default.description}"`,
			"status":            `"run"`,
		}),
	}
	dbtypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_polardb_cluster.default.description}"`,
			"db_type":           `"${alicloud_polardb_cluster.default.db_type}"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_polardb_cluster.default.description}"`,
			"db_type":           `"Oracle"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_polardb_cluster.default.description}"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_polardb_cluster.default.description}"`,
			"tags": `{ 
						"key1" = "value1_fake" 
						"key2" = "value2_fake" 
					}`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_polardb_cluster.default.description}"`,
			"status":            `"Running"`,
			"db_type":           `"${alicloud_polardb_cluster.default.db_type}"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_polardb_cluster.default.description}"`,
			"status":            `"run"`,
			"db_type":           `"Oracle"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
		}),
	}

	var existPolarClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"descriptions.#":            "1",
			"clusters.#":                "1",
			"clusters.0.id":             CHECKSET,
			"clusters.0.description":    CHECKSET,
			"clusters.0.charge_type":    "PostPaid",
			"clusters.0.network_type":   "VPC",
			"clusters.0.region_id":      CHECKSET,
			"clusters.0.zone_id":        CHECKSET,
			"clusters.0.expired":        "false",
			"clusters.0.status":         "Running",
			"clusters.0.engine":         "POLARDB",
			"clusters.0.db_type":        "MySQL",
			"clusters.0.db_version":     "8.0",
			"clusters.0.lock_mode":      "Unlock",
			"clusters.0.delete_lock":    "0",
			"clusters.0.create_time":    CHECKSET,
			"clusters.0.vpc_id":         CHECKSET,
			"clusters.0.db_node_number": "2",
			"clusters.0.db_node_class":  CHECKSET,
			"clusters.0.storage_used":   CHECKSET,
			"clusters.0.db_nodes.#":     "2",
		}
	}

	var fakePolarClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"clusters.#":     CHECKSET,
			"ids.#":          CHECKSET,
			"descriptions.#": CHECKSET,
		}
	}

	var PolarClusterCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_polardb_clusters.default",
		existMapFunc: existPolarClusterMapFunc,
		fakeMapFunc:  fakePolarClusterMapFunc,
	}

	PolarClusterCheckInfo.dataSourceTestCheck(t, rand, nameConf, statusConf, dbtypeConf, tagsConf, allConf)
}

func testAccCheckAlicloudPolarClusterDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccPolarDBConfig_%d"
	}
	data "alicloud_polardb_node_classes" "this" {
	  db_type    = "MySQL"
	  db_version = "8.0"
      pay_type   = "PostPaid"
	  zone_id    = local.zone_id
	}

	resource "alicloud_polardb_cluster" "default" {
		db_type = "MySQL"
		db_version = "8.0"
		pay_type = "PostPaid"
	    db_node_class     = data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class
		vswitch_id = local.vswitch_id
		description = "${var.name}"
		tags = {
			"key1" = "value1"
			"key2" = "value2"
		}
	}
	data "alicloud_polardb_clusters" "default" {
	  %s
	}
`, PolarDBCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
