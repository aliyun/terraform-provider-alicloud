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
			"ids.#":                                     "1",
			"clusters.#":                                "1",
			"clusters.0.db_cluster_description":         fmt.Sprintf("tf-testAccClickhouseDbCluster-%d", rand),
			"clusters.0.payment_type":                   "PayAsYouGo",
			"clusters.0.category":                       "Basic",
			"clusters.0.db_cluster_access_white_list.#": "1",
			"clusters.0.db_cluster_access_white_list.0.db_cluster_ip_array_name": "test",
			"clusters.0.db_cluster_access_white_list.0.security_ip_list":         "192.168.0.1",
			"clusters.0.ali_uid":        CHECKSET,
			"clusters.0.bid":            CHECKSET,
			"clusters.0.commodity_code": CHECKSET,
			// TODO There is an api bug that did not return the connection string value. Reopen it after the bug is fixed.
			//"clusters.0.connection_string":        CHECKSET,
			"clusters.0.create_time":              CHECKSET,
			"clusters.0.id":                       CHECKSET,
			"clusters.0.db_cluster_id":            CHECKSET,
			"clusters.0.db_cluster_network_type":  "vpc",
			"clusters.0.db_cluster_type":          "Common",
			"clusters.0.db_node_class":            "S8",
			"clusters.0.db_node_count":            "1",
			"clusters.0.db_node_storage":          "100",
			"clusters.0.encryption_key":           "",
			"clusters.0.encryption_type":          "",
			"clusters.0.engine":                   "clickhouse",
			"clusters.0.engine_version":           CHECKSET,
			"clusters.0.expire_time":              "",
			"clusters.0.is_expired":               "",
			"clusters.0.lock_mode":                "Unlock",
			"clusters.0.lock_reason":              "",
			"clusters.0.maintain_time":            CHECKSET,
			"clusters.0.port":                     CHECKSET,
			"clusters.0.public_connection_string": "",
			"clusters.0.public_port":              "",
			"clusters.0.scale_out_status.#":       "0",
			"clusters.0.storage_type":             "CloudESSD",
			"clusters.0.support_backup":           "1",
			"clusters.0.support_https_port":       "true",
			"clusters.0.support_mysql_port":       "true",
			"clusters.0.vswitch_id":               CHECKSET,
			"clusters.0.vpc_cloud_instance_id":    CHECKSET,
			"clusters.0.vpc_id":                   CHECKSET,
			"clusters.0.zone_id":                  CHECKSET,
			"clusters.0.control_version":          CHECKSET,
			"clusters.0.status":                   "Running",
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
data "alicloud_click_house_regions" "default" {	
  current = true
}
data "alicloud_vpcs" "default"	{
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  zone_id = data.alicloud_click_house_regions.default.regions.0.zone_ids.0.zone_id
}
resource "alicloud_click_house_db_cluster" "default" {
  db_cluster_version      = "20.3.10.75"
  status                  = "Running"
  category                = "Basic"
  db_cluster_class        = "S8"
  db_cluster_network_type = "vpc"
  db_cluster_description  = var.name
  db_node_group_count     = "1"
  payment_type            = "PayAsYouGo"
  db_node_storage         = "100"
  storage_type            = "cloud_essd"
  vswitch_id              = data.alicloud_vswitches.default.vswitches.0.id
  db_cluster_access_white_list {
    db_cluster_ip_array_name      = "test"
    security_ip_list              = "192.168.0.1"
  }
}

data "alicloud_click_house_db_clusters" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
