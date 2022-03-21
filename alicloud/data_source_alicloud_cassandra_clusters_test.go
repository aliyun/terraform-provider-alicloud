package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

var existCassandraMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"clusters.#":              CHECKSET,
		"clusters.0.id":           CHECKSET,
		"clusters.0.cluster_name": CHECKSET,
		"clusters.0.status":       CHECKSET,
		"clusters.0.tags.%":       "2",
		"clusters.0.tags.Created": "TF",
		"clusters.0.tags.For":     "acceptance test",
		"ids.#":                   "1",
		"ids.0":                   CHECKSET,
		"names.#":                 "1",
	}
}

var fakeCassandraMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"clusters.#": "0",
		"ids.#":      "0",
		"names.#":    "0",
	}
}

var checkCassandraInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cassandra_clusters.default",
	existMapFunc: existCassandraMapFunc,
	fakeMapFunc:  fakeCassandraMapFunc,
}

func SkipTestAccAlicloudCassandraClustersDataSourceNewCluster(t *testing.T) {
	// Cassandra has been offline
	t.Skip("Cassandra has been offline")
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCassandraDataSourceConfigNewCluster(rand, map[string]string{
			"name_regex": `"${alicloud_cassandra_cluster.default.cluster_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCassandraDataSourceConfigNewCluster(rand, map[string]string{
			"name_regex": `"${alicloud_cassandra_cluster.default.cluster_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCassandraDataSourceConfigNewCluster(rand, map[string]string{
			"ids": `["${alicloud_cassandra_cluster.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCassandraDataSourceConfigNewCluster(rand, map[string]string{
			"ids": `["${alicloud_cassandra_cluster.default.id}_fake"]`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCassandraDataSourceConfigNewCluster(rand, map[string]string{
			"name_regex": `"${alicloud_cassandra_cluster.default.cluster_name}"`,
			"tags":       `{Created = "TF"}`,
		}),
		fakeConfig: testAccCheckAlicloudCassandraDataSourceConfigNewCluster(rand, map[string]string{
			"name_regex": `"${alicloud_cassandra_cluster.default.cluster_name}"`,
			"tags":       `{Created = "TF1"}`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCassandraDataSourceConfigNewCluster(rand, map[string]string{
			"name_regex": `"${alicloud_cassandra_cluster.default.cluster_name}"`,
			"ids":        `["${alicloud_cassandra_cluster.default.id}"]`,
			"tags":       `{Created = "TF"}`,
		}),
		fakeConfig: testAccCheckAlicloudCassandraDataSourceConfigNewCluster(rand, map[string]string{
			"name_regex": `"${alicloud_cassandra_cluster.default.cluster_name}"`,
			"ids":        `["${alicloud_cassandra_cluster.default.id}_fake"]`,
			"tags":       `{Created = "TF1"}`,
		}),
	}

	checkCassandraInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, tagsConf, allConf)
}

// new a cluster config
func testAccCheckAlicloudCassandraDataSourceConfigNewCluster(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
		variable "name" {
		  default = "tf-testAccCassandraCluster_datasource_%d"
		}
		data "alicloud_cassandra_zones" "default" {
		}
		
		data "alicloud_vpcs" "default" {
			name_regex = "default-NODELETING"
		}
		
		data "alicloud_vswitches" "default" {
		  vpc_id = data.alicloud_vpcs.default.ids[0]
		  zone_id = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)-1].id
		}
		
		resource "alicloud_vswitch" "this" {
		  count = "${length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1}"
		  vswitch_name = "${var.name}"
		  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
		  availability_zone = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)-1].id
		  cidr_block = "${cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)}"
		}
		
		resource "alicloud_cassandra_cluster" "default" {
		  cluster_name = "${var.name}"
		  data_center_name = "${var.name}"
		  auto_renew = false
		  instance_type = "cassandra.c.large"
		  major_version = "3.11"
		  node_count = "2"
		  pay_type = "PayAsYouGo"
		  vswitch_id = "${length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : alicloud_vswitch.this[0].id}"
		  disk_size = "160"
		  disk_type = "cloud_ssd"
		  maintain_start_time = "18:00Z"
		  maintain_end_time = "20:00Z"
		  ip_white = "127.0.0.2"
		  enable_public = "true"
		  zone_id = "${data.alicloud_vswitches.default.zone_id}"
		  tags = {
			Created = "TF"
			For     = "acceptance test"
		  }
		}
		
		data "alicloud_cassandra_clusters" "default" {
		  %s
		}
		`, rand, strings.Join(pairs, "\n  "))
	return config
}
