package gdb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudGraphDatabaseDbInstancesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 1000)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGraphDatabaseDbInstanceDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_graph_database_db_instance.default.id}"]`,
			"enable_details": "true",
		}),
		FakeConfig: testAccCheckAlicloudGraphDatabaseDbInstanceDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_graph_database_db_instance.default.id}_fakeid"]`,
			"enable_details": "true",
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGraphDatabaseDbInstanceDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_graph_database_db_instance.default.id}"]`,
			"status": `"Running"`,
		}),
		FakeConfig: testAccCheckAlicloudGraphDatabaseDbInstanceDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_graph_database_db_instance.default.id}"]`,
			"status": `"Creating"`,
		}),
	}

	descriptionConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGraphDatabaseDbInstanceDataSourceName(rand, map[string]string{
			"ids":                     `["${alicloud_graph_database_db_instance.default.id}"]`,
			"db_instance_description": `"${alicloud_graph_database_db_instance.default.db_instance_description}"`,
		}),
		FakeConfig: testAccCheckAlicloudGraphDatabaseDbInstanceDataSourceName(rand, map[string]string{
			"ids":                     `["${alicloud_graph_database_db_instance.default.id}"]`,
			"db_instance_description": `"${alicloud_graph_database_db_instance.default.db_instance_description}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGraphDatabaseDbInstanceDataSourceName(rand, map[string]string{
			"ids":                     `["${alicloud_graph_database_db_instance.default.id}"]`,
			"status":                  `"Running"`,
			"db_instance_description": `"${alicloud_graph_database_db_instance.default.db_instance_description}"`,
		}),
		FakeConfig: testAccCheckAlicloudGraphDatabaseDbInstanceDataSourceName(rand, map[string]string{
			"ids":                     `["${alicloud_graph_database_db_instance.default.id}"]`,
			"status":                  `"Creating"`,
			"db_instance_description": `"${alicloud_graph_database_db_instance.default.db_instance_description}_fake"`,
		}),
	}

	var existDataAlicloudGraphDatabaseDbInstancesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                               "1",
			"instances.#":                         "1",
			"instances.0.db_instance_description": fmt.Sprintf("tf-testaccgraphdatabasedbinstance%d", rand),
		}
	}
	var fakeDataGraphDatabaseDbInstancesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}
	var alicloudGraphDatabaseDbInstanceCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_graph_database_db_instances.default",
		ExistMapFunc: existDataAlicloudGraphDatabaseDbInstancesSourceNameMapFunc,
		FakeMapFunc:  fakeDataGraphDatabaseDbInstancesSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.GraphDatabaseSupportRegions)
	}

	alicloudGraphDatabaseDbInstanceCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, descriptionConf, allConf)
}

func testAccCheckAlicloudGraphDatabaseDbInstanceDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testaccgraphdatabasedbinstance%d"
}

resource "alicloud_graph_database_db_instance" "default" {
  db_node_class = "gdb.r.xlarge"
  db_instance_network_type = "vpc"
  db_version = "1.0"
  db_instance_category = "HA"
  db_instance_storage_type = "cloud_ssd"
  db_node_storage = "50"
  payment_type = "PayAsYouGo"
  db_instance_description = var.name
}

data "alicloud_graph_database_db_instances" "default" {
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
