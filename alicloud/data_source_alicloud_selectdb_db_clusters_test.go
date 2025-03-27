package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSelectDBDbClustersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	dbClusterIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSelectDBDbClustersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_selectdb_db_cluster.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSelectDBDbClustersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_selectdb_db_cluster.default.id}_fake"]`,
		}),
	}
	var existAlicloudSelectDBDbClustersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"clusters.#":                     "1",
			"clusters.0.payment_type":        "PayAsYouGo",
			"clusters.0.db_cluster_id":       CHECKSET,
			"clusters.0.db_instance_id":      CHECKSET,
			"clusters.0.db_cluster_class":    "selectdb.2xlarge",
			"clusters.0.engine":              "selectdb",
			"clusters.0.create_time":         CHECKSET,
			"clusters.0.cpu":                 "8",
			"clusters.0.memory":              "32",
			"clusters.0.cache_size":          CHECKSET,
			"clusters.0.region_id":           CHECKSET,
			"clusters.0.vpc_id":              CHECKSET,
			"clusters.0.zone_id":             CHECKSET,
			"clusters.0.params.#":            "0",
			"clusters.0.param_change_logs.#": "0",
		}
	}
	var fakeAlicloudSelectDBDbClustersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudSelectDBDbClustersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_selectdb_db_clusters.default",
		existMapFunc: existAlicloudSelectDBDbClustersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudSelectDBDbClustersDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.SelectDBSupportRegions)
	}

	alicloudSelectDBDbClustersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, dbClusterIdsConf)
}

func testAccCheckAlicloudSelectDBDbClustersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf_testAccSelectDBDbCluster_%d"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.ids.0
}
resource "alicloud_selectdb_db_instance" "default" {
  db_instance_class       = "selectdb.2xlarge"
  db_instance_description = var.name
  cache_size              = "400"
  engine_minor_version    = "3.0.12"
  payment_type            = "PayAsYouGo"
  vpc_id                  = "${data.alicloud_vpcs.default.ids.0}"
  zone_id                 = "${data.alicloud_zones.default.ids.0}"
  vswitch_id              = "${data.alicloud_vswitches.default.ids.0}"
}
resource "alicloud_selectdb_db_cluster" "default" {
  db_instance_id         = "${alicloud_selectdb_db_instance.default.id}"
  db_cluster_description = var.name
  db_cluster_class       = "selectdb.2xlarge"
  cache_size             = "400"
  payment_type           = "PayAsYouGo"
}

data "alicloud_selectdb_db_clusters" "default" {
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
