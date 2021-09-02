package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudClickHouseDbClusterDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudClickHouseDbClusterDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_click_house_db_cluster.default.id}"]`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudClickHouseDbClusterDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_click_house_db_cluster.default.id}_fake"]`,
			"enable_details": "true",
		}),
	}
	dbClusterDescConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudClickHouseDbClusterDataSourceName(rand, map[string]string{
			"db_cluster_description": `"${alicloud_click_house_db_cluster.default.db_cluster_description}"`,
			"enable_details":         "true",
		}),
		fakeConfig: testAccCheckAlicloudClickHouseDbClusterDataSourceName(rand, map[string]string{
			"db_cluster_description": `"${alicloud_click_house_db_cluster.default.db_cluster_description}_fake"`,
			"enable_details":         "true",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudClickHouseDbClusterDataSourceName(rand, map[string]string{
			"status":         `"${alicloud_click_house_db_cluster.default.status}"`,
			"ids":            `["${alicloud_click_house_db_cluster.default.id}"]`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudClickHouseDbClusterDataSourceName(rand, map[string]string{
			"status":         `"Preparing"`,
			"ids":            `["${alicloud_click_house_db_cluster.default.id}"]`,
			"enable_details": "true",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudClickHouseDbClusterDataSourceName(rand, map[string]string{
			"ids":                    `["${alicloud_click_house_db_cluster.default.id}"]`,
			"db_cluster_description": `"${alicloud_click_house_db_cluster.default.db_cluster_description}"`,
			"status":                 `"${alicloud_click_house_db_cluster.default.status}"`,
			"enable_details":         "true",
		}),
		fakeConfig: testAccCheckAlicloudClickHouseDbClusterDataSourceName(rand, map[string]string{
			"ids":                    `["${alicloud_click_house_db_cluster.default.id}_fake"]`,
			"db_cluster_description": `"${alicloud_click_house_db_cluster.default.db_cluster_description}_fake"`,
			"status":                 `"Preparing"`,
			"enable_details":         "true",
		}),
	}
	var existAlicloudClickHouseDbClusterDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"clusters.#":                        "1",
			"clusters.0.db_cluster_description": fmt.Sprintf("tf-testAccClickhouseDbCluster-%d", rand),
			"clusters.0.payment_type":           "PayAsYouGo",
			"clusters.0.category":               "Basic",
		}
	}
	var fakeAlicloudClickHouseDbClusterDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudClickHouseDbClusterCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_click_house_db_clusters.default",
		existMapFunc: existAlicloudClickHouseDbClusterDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudClickHouseDbClusterDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.ClickHouseSupportRegions)
	}

	alicloudClickHouseDbClusterCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, dbClusterDescConf, statusConf, allConf)
}
func testAccCheckAlicloudClickHouseDbClusterDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {	
	default = "tf-testAccClickhouseDbCluster-%d"
}
data "alicloud_vpcs" "default"	{
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
}
resource "alicloud_click_house_db_cluster" "default" {
  db_cluster_version    =  "20.3.10.75"
  status = "Running"
  category=                "Basic"
  db_cluster_class =        "S8"
  db_cluster_network_type= "vpc"
  db_cluster_description = var.name
  db_node_group_count=     "1"
  payment_type=            "PayAsYouGo"
  db_node_storage=         "500"
  storage_type=            "cloud_essd"
  vswitch_id=              "${data.alicloud_vswitches.default.vswitches.0.id}"
}

data "alicloud_click_house_db_clusters" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
