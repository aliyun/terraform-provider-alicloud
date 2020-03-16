package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudPolarDBClusterEndPointsDataSource(t *testing.T) {
	rand := acctest.RandInt()

	dbClusterIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterEndPointsDataSourceConfig(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_cluster.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterEndPointsDataSourceConfig(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_cluster.default.id}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudPolarClusterEndPointsDataSourceConfig(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_cluster.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudPolarClusterEndPointsDataSourceConfig(rand, map[string]string{
			"db_cluster_id":  `"${alicloud_polardb_cluster.default.id}"`,
			"db_endpoint_id": `"^test1234"`,
		}),
	}

	var existPolarClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"endpoints.0.db_endpoint_id": CHECKSET,
		}
	}

	var fakePolarClusterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"endpoints.#": CHECKSET,
		}
	}

	var PolarClusterCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_polardb_endpoints.default",
		existMapFunc: existPolarClusterMapFunc,
		fakeMapFunc:  fakePolarClusterMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithNoDefaultVpc(t)
	}

	PolarClusterCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, dbClusterIdConf, allConf)
}

func testAccCheckAlicloudPolarClusterEndPointsDataSourceConfig(rand int, attrMap map[string]string) string {
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
	
	data "alicloud_polardb_endpoints" "default" {
	  %s
	}
`, PolarDBCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
