package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudPolarDBClusterDatabasesDataSource(t *testing.T) {
	rand := acctest.RandInt()

	idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterDatabasesDataSourceConfig(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_database.database.db_cluster_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterDatabasesDataSourceConfig(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_database.database.db_cluster_id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterDatabasesDataSourceConfig(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_database.database.db_cluster_id}"`,
			"name_regex":    `"${alicloud_polardb_database.database.db_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterDatabasesDataSourceConfig(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_database.database.db_cluster_id}"`,
			"name_regex":    `"^test1234"`,
		}),
	}

	var existPolarClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                        CHECKSET,
			"databases.0.character_set_name": CHECKSET,
			"databases.0.db_name":            CHECKSET,
			"databases.0.db_description":     CHECKSET,
			"databases.0.db_status":          CHECKSET,
			"databases.0.engine":             CHECKSET,
			"databases.0.accounts.#":         CHECKSET,
		}
	}

	var fakePolarClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"databases.#": CHECKSET,
		}
	}

	var PolarClusterCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_polardb_databases.default",
		existMapFunc: existPolarClusterMapFunc,
		fakeMapFunc:  fakePolarClusterMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithNoDefaultVpc(t)
	}

	PolarClusterCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idConf, allConf)
}

func testAccCheckAlicloudPolarClusterDatabasesDataSourceConfig(rand int, attrMap map[string]string) string {
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
	  default = "tf-testAccPolarClusterConfig_%d"
	}
	
	resource "alicloud_polardb_cluster" "default" {
	  db_type           = "MySQL"
	  db_version        = "8.0"
      pay_type          = "PostPaid"
	  db_node_class     = "polar.mysql.x4.large"
	  vswitch_id        = "${data.alicloud_vswitches.default.ids.0}"
	  description       = "${var.name}"
	}

	resource "alicloud_polardb_database" "database" {
	  db_cluster_id     = "${alicloud_polardb_cluster.default.id}"
	  db_name           = "tftestdatabase"
      db_description    = "${var.name}"
	}

	data "alicloud_polardb_databases" "default" {
	  %s
	}
`, PolarDBCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
