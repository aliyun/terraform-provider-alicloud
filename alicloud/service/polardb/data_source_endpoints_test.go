package polardb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudPolarDBClusterEndPointsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()

	dbClusterIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPolarClusterEndPointsDataSourceConfig(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_cluster.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudPolarClusterEndPointsDataSourceConfig(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_cluster.default.id}"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudPolarClusterEndPointsDataSourceConfig(rand, map[string]string{
			"db_cluster_id": `"${alicloud_polardb_cluster.default.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudPolarClusterEndPointsDataSourceConfig(rand, map[string]string{
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

	var PolarClusterCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_polardb_endpoints.default",
		ExistMapFunc: existPolarClusterMapFunc,
		FakeMapFunc:  fakePolarClusterMapFunc,
	}

	PolarClusterCheckInfo.DataSourceTestCheck(t, rand, dbClusterIdConf, allConf)
}

func testAccCheckAlicloudPolarClusterEndPointsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
	%s
	variable "name" {
	  default = "tf-testAccPolarClusterConfig_%d"
	}

	resource "alicloud_polardb_cluster" "default" {
	  db_type           = "MySQL"
	  db_version        = "8.0"
      pay_type          = "PostPaid"
      db_node_class     = "polar.mysql.x4.large"
	  vswitch_id        = local.vswitch_id
	  description       = "${var.name}"
	  lifecycle {
	    ignore_changes = [encrypt_new_tables]
	  }
	}
	
	data "alicloud_polardb_endpoints" "default" {
	  %s
	}
`, PolarDBCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
