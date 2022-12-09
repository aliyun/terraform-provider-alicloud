package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudPolarDBGlobalDatabaseNetworkDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_polardb_global_database_network.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_polardb_global_database_network.default.id}_fake"]`,
		}),
	}
	gdnIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"gdn_id": `"${alicloud_polardb_global_database_network.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"gdn_id": `"${alicloud_polardb_global_database_network.default.id}_fake"`,
		}),
	}
	dbClusterIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_global_database_network.default.db_cluster_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_global_database_network.default.db_cluster_id}_fake"`,
		}),
	}
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"description": `"${alicloud_polardb_global_database_network.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"description": `"${alicloud_polardb_global_database_network.default.description}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_polardb_global_database_network.default.id}"]`,
			"status": `"active"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_polardb_global_database_network.default.id}_fake"]`,
			"status": `"locked"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_polardb_global_database_network.default.id}"]`,
			"gdn_id":        `"${alicloud_polardb_global_database_network.default.id}"`,
			"db_cluster_id": `"${alicloud_polardb_global_database_network.default.db_cluster_id}"`,
			"description":   `"${alicloud_polardb_global_database_network.default.description}"`,
			"status":        `"active"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_polardb_global_database_network.default.id}_fake"]`,
			"gdn_id":        `"${alicloud_polardb_global_database_network.default.id}_fake"`,
			"db_cluster_id": `"${alicloud_polardb_global_database_network.default.db_cluster_id}_fake"`,
			"description":   `"${alicloud_polardb_global_database_network.default.description}_fake"`,
			"status":        `"locked"`,
		}),
	}
	var existAlicloudPolarDBGlobalDatabaseNetworkDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"networks.#":                             "1",
			"networks.0.id":                          CHECKSET,
			"networks.0.gdn_id":                      CHECKSET,
			"networks.0.description":                 CHECKSET,
			"networks.0.db_type":                     CHECKSET,
			"networks.0.db_version":                  CHECKSET,
			"networks.0.create_time":                 CHECKSET,
			"networks.0.status":                      "active",
			"networks.0.db_clusters.#":               "1",
			"networks.0.db_clusters.0.db_cluster_id": CHECKSET,
			"networks.0.db_clusters.0.role":          CHECKSET,
			"networks.0.db_clusters.0.region_id":     CHECKSET,
		}
	}
	var fakeAlicloudPolarDBGlobalDatabaseNetworkDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"networks.#": "0",
		}
	}
	var alicloudPolarDBGlobalDatabaseNetworkCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_polardb_global_database_networks.default",
		existMapFunc: existAlicloudPolarDBGlobalDatabaseNetworkDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudPolarDBGlobalDatabaseNetworkDataSourceNameMapFunc,
	}
	alicloudPolarDBGlobalDatabaseNetworkCheckInfo.dataSourceTestCheck(t, rand, idsConf, gdnIdConf, dbClusterIdConf, descriptionConf, statusConf, allConf)
}

func testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
	variable "name" {
		default = "tf-testAcc-%d"
	}
	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}
	
	data "alicloud_vswitches" "default" {
		vpc_id = data.alicloud_vpcs.default.ids.0
	}
	
	data "alicloud_polardb_node_classes" "default" {
		zone_id    = data.alicloud_vswitches.default.vswitches.0.zone_id
		pay_type   = "PostPaid"
		db_type    = "MySQL"
		db_version = "8.0"
	}
	resource "alicloud_polardb_cluster" "default" {
		db_type       = "MySQL"
		db_version    = "8.0"
		pay_type      = "PostPaid"
		db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
		vswitch_id    = data.alicloud_vswitches.default.ids.0
		description   = "${var.name}"
	}
	resource "alicloud_polardb_global_database_network" "default" {
		db_cluster_id = "${alicloud_polardb_cluster.default.id}"
		description   = var.name
	}
	data "alicloud_polardb_global_database_networks" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
