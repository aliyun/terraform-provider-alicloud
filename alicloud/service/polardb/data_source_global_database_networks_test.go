package polardb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudPolarDBGlobalDatabaseNetworkDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_polardb_global_database_network.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_polardb_global_database_network.default.id}_fake"]`,
		}),
	}
	gdnIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"gdn_id": `"${alicloud_polardb_global_database_network.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"gdn_id": `"${alicloud_polardb_global_database_network.default.id}_fake"`,
		}),
	}
	dbClusterIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_global_database_network.default.db_cluster_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_global_database_network.default.db_cluster_id}_fake"`,
		}),
	}
	descriptionConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"description": `"${alicloud_polardb_global_database_network.default.description}"`,
		}),
		FakeConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"description": `"${alicloud_polardb_global_database_network.default.description}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_polardb_global_database_network.default.id}"]`,
			"status": `"active"`,
		}),
		FakeConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_polardb_global_database_network.default.id}_fake"]`,
			"status": `"locked"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_polardb_global_database_network.default.id}"]`,
			"gdn_id":        `"${alicloud_polardb_global_database_network.default.id}"`,
			"db_cluster_id": `"${alicloud_polardb_global_database_network.default.db_cluster_id}"`,
			"description":   `"${alicloud_polardb_global_database_network.default.description}"`,
			"status":        `"active"`,
		}),
		FakeConfig: testAccCheckAlicloudPolarDBGlobalDatabaseNetworkDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_polardb_global_database_network.default.id}_fake"]`,
			"gdn_id":        `"${alicloud_polardb_global_database_network.default.id}_fake"`,
			"db_cluster_id": `"${alicloud_polardb_global_database_network.default.db_cluster_id}_fake"`,
			"description":   `"${alicloud_polardb_global_database_network.default.description}_fake"`,
			"status":        `"locked"`,
		}),
	}
	// DescribeGlobalDatabaseNetworks does not return ResourceGroupId, so the
	// value may be empty until the API contract is confirmed. The assertion
	// requires the key and, once populated, a resource group ID shape.
	var existAlicloudPolarDBGlobalDatabaseNetworkDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"networks.#":                             "1",
			"networks.0.id":                          CHECKSET,
			"networks.0.gdn_id":                      CHECKSET,
			"networks.0.description":                 CHECKSET,
			"networks.0.gdn_version":                 CHECKSET,
			"networks.0.resource_group_id":           REGEXMATCH + "^(|rg-[0-9a-z]+)$",
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
	var alicloudPolarDBGlobalDatabaseNetworkCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_polardb_global_database_networks.default",
		ExistMapFunc: existAlicloudPolarDBGlobalDatabaseNetworkDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudPolarDBGlobalDatabaseNetworkDataSourceNameMapFunc,
	}
	alicloudPolarDBGlobalDatabaseNetworkCheckInfo.DataSourceTestCheck(t, rand, idsConf, gdnIdConf, dbClusterIdConf, descriptionConf, statusConf, allConf)
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
	data "alicloud_resource_manager_resource_groups" "default" {}

	data "alicloud_polardb_zones" "default" {}

	resource "alicloud_vpc" "default" {
		vpc_name = var.name
	}

	resource "alicloud_vswitch" "default" {
		zone_id    = data.alicloud_polardb_zones.default.ids[length(data.alicloud_polardb_zones.default.ids)-1]
		vpc_id     = alicloud_vpc.default.id
		cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
	}

	resource "alicloud_polardb_cluster" "default" {
		db_type       = "MySQL"
		db_version    = "8.0"
		pay_type      = "PostPaid"
		db_node_class = "polar.mysql.x4.medium"
		vswitch_id    = alicloud_vswitch.default.id
		description   = "${var.name}"
		lifecycle {
			ignore_changes = [encrypt_new_tables]
		}
	}
	resource "alicloud_polardb_global_database_network" "default" {
		db_cluster_id     = "${alicloud_polardb_cluster.default.id}"
		description       = var.name
		resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.1
		gdn_version       = "1.0"
	}
	data "alicloud_polardb_global_database_networks" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
