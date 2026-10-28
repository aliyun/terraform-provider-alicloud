package polardb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudPolarDBClustersDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	nameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_polardb_cluster.default.description}"`,
		}),
		FakeConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"^test1234"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_polardb_cluster.default.description}"`,
			"status":            `"Running"`,
		}),
		FakeConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_polardb_cluster.default.description}"`,
			"status":            `"run"`,
		}),
	}
	dbtypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_polardb_cluster.default.description}"`,
			"db_type":           `"${alicloud_polardb_cluster.default.db_type}"`,
		}),
		FakeConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_polardb_cluster.default.description}"`,
			"db_type":           `"Oracle"`,
		}),
	}
	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_polardb_cluster.default.description}"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
		}),
		FakeConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_polardb_cluster.default.description}"`,
			"tags": `{ 
						"key1" = "value1_fake" 
						"key2" = "value2_fake" 
					}`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alicloud_polardb_cluster.default.description}"`,
			"status":            `"Running"`,
			"db_type":           `"${alicloud_polardb_cluster.default.db_type}"`,
			"tags": `{ 
						"key1" = "value1" 
						"key2" = "value2" 
					}`,
		}),
		FakeConfig: testAccCheckAlicloudPolarClusterDataSourceConfig(rand, map[string]string{
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
			"ids.#":                        "1",
			"descriptions.#":               "1",
			"clusters.#":                   "1",
			"clusters.0.id":                CHECKSET,
			"clusters.0.description":       CHECKSET,
			"clusters.0.charge_type":       "PostPaid",
			"clusters.0.network_type":      "VPC",
			"clusters.0.region_id":         CHECKSET,
			"clusters.0.zone_id":           CHECKSET,
			"clusters.0.expired":           "false",
			"clusters.0.status":            "Running",
			"clusters.0.engine":            "POLARDB",
			"clusters.0.db_type":           "MySQL",
			"clusters.0.db_version":        "8.0",
			"clusters.0.lock_mode":         "Unlock",
			"clusters.0.delete_lock":       "0",
			"clusters.0.create_time":       CHECKSET,
			"clusters.0.vpc_id":            CHECKSET,
			"clusters.0.db_node_number":    "2",
			"clusters.0.db_node_class":     CHECKSET,
			"clusters.0.storage_used":      CHECKSET,
			"clusters.0.db_nodes.#":        "2",
			"clusters.0.connection_string": CHECKSET,
			"clusters.0.port":              "3306",
		}
	}

	var fakePolarClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"clusters.#":     CHECKSET,
			"ids.#":          CHECKSET,
			"descriptions.#": CHECKSET,
		}
	}

	var PolarClusterCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_polardb_clusters.default",
		ExistMapFunc: existPolarClusterMapFunc,
		FakeMapFunc:  fakePolarClusterMapFunc,
	}

	PolarClusterCheckInfo.DataSourceTestCheck(t, rand, nameConf, statusConf, dbtypeConf, tagsConf, allConf)
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
      category   = "Normal"
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
		security_ips = ["10.168.1.12", "100.69.7.112"]
	}
	data "alicloud_polardb_clusters" "default" {
	  %s
	}
`, PolarDBCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
