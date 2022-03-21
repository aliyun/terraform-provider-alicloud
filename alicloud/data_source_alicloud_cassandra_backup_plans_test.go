package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func SkipTestAccAlicloudCassandraBackupPlansDataSource(t *testing.T) {
	// Cassandra has been offline
	t.Skip("Cassandra has been offline")
	rand := acctest.RandInt()

	var existAlicloudCassandraBackupPlanDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "1",
			"plans.data_center_id":   CHECKSET,
			"plans.cluster_id":       CHECKSET,
			"plans.backup_time":      "00:10Z",
			"plans.active":           "false",
			"plans.create_time":      CHECKSET,
			"plans.backup_period":    CHECKSET,
			"plans.retention_period": CHECKSET,
		}
	}
	var fakeAlicloudCassandraBackupPlanDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}

	var alicloudCassandraBackupPlanCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cassandra_backup_plans.default_2",
		existMapFunc: existAlicloudCassandraBackupPlanDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCassandraBackupPlanDataSourceNameMapFunc,
	}

	alicloudCassandraBackupPlanCheckInfo.dataSourceTestCheck(t, rand)
}

func testAccAlicloudCassandraBackupPlansDataSource(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
		variable "name" {
		  default = "tf-testAccCassandrBackupPlans_datasource_%d"
		}

		data "alicloud_cassandra_zones" "default" {
		}
		
		data "alicloud_vpcs" "default" {
		  name_regex = "default-NODELETING"
		}
		
		data "alicloud_vswitches" "default_1" {
		  vpc_id = data.alicloud_vpcs.default.ids[0]
		  zone_id = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)+(-1)].id
		}
		
		resource "alicloud_vswitch" "this_1" {
		  count = length(data.alicloud_vswitches.default_1.ids) > 0 ? 0 : 1
		  vswitch_name = var.name
		  vpc_id = data.alicloud_vpcs.default.ids.0
		  zone_id = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)+(-1)].id
		  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
		}
		resource "alicloud_cassandra_cluster" "default" {
		  cluster_name = var.name
		  data_center_name = var.name
		  auto_renew = "false"
		  instance_type = "cassandra.c.large"
		  major_version = "3.11"
		  node_count = "2"
		  pay_type = "PayAsYouGo"
		  vswitch_id = length(data.alicloud_vswitches.default_1.ids) > 0 ? data.alicloud_vswitches.default_1.ids[0] : alicloud_vswitch.this_1[0].id
		  disk_size = "160"
		  disk_type = "cloud_ssd"
		  maintain_start_time = "18:00Z"
		  maintain_end_time = "20:00Z"
		  ip_white = "127.0.0.1"
		}

		resource "alicloud_cassandra_backup_plan" "default" {
		  cluster_id = alicloud_cassandra_cluster.default.id
          data_center_id = alicloud_cassandra_cluster.default.zone_id
          backup_time = "00:10Z"
          active = false
		}

		data "alicloud_cassandra_backup_plans" "default_2" {
			cluster_id = alicloud_cassandra_backup_plan.default.cluster_id
			%s
		}
		`, rand, strings.Join(pairs, " \n "))
	return config
}
