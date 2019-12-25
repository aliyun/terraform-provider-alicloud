package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudPolarDBClustersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `alicloud_polardb_cluster.default.description`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"^test1234"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `alicloud_polardb_cluster.default.description`,
			"status":            `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `alicloud_polardb_cluster.default.description`,
			"status":            `"run"`,
		}),
	}
	dbtypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `alicloud_polardb_cluster.default.description`,
			"db_type":           `alicloud_polardb_cluster.default.db_type`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `alicloud_polardb_cluster.default.description`,
			"db_type":           `"Oracle"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `alicloud_polardb_cluster.default.description`,
			"status":            `"Running"`,
			"db_type":           `alicloud_polardb_cluster.default.db_type`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `alicloud_polardb_cluster.default.description`,
			"status":            `"run"`,
			"db_type":           `"Oracle"`,
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
			"clusters.0.db_node_class":  "polar.mysql.x4.large",
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

	PolarClusterCheckInfo.dataSourceTestCheck(t, rand, nameConf, statusConf, dbtypeConf, allConf)
}

func testAccCheckAlicloudPolarClusterDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
	%s
	variable "creation" {
		default = "PolarDB"
	}

	variable "name" {
		default = "pc-testAccDBInstanceConfig_%d"
	}

	resource "alicloud_polardb_cluster" "default" {
		db_type = "MySQL"
		db_version = "8.0"
		pay_type = "PostPaid"
		db_node_class = "polar.mysql.x4.large"
		vswitch_id = alicloud_vswitch.default.id
		description = var.name
	}
	data "alicloud_polardb_clusters" "default" {
	  %s
	}
`, PolarDBCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
