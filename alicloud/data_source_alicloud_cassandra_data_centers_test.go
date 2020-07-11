package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

var existCassandraDcMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"centers.#":                  CHECKSET,
		"centers.0.data_center_id":   CHECKSET,
		"centers.0.data_center_name": CHECKSET,
		"centers.0.status":           CHECKSET,
		"centers.0.zone_id":          CHECKSET,
		"ids.#":                      "1",
		"ids.0":                      CHECKSET,
		"names.#":                    "1",
	}
}

var fakeCassandraDcMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"centers.#": "0",
		"ids.#":     "0",
		"names.#":   "0",
	}
}

var checkCassandraDcInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cassandra_data_centers.default",
	existMapFunc: existCassandraDcMapFunc,
	fakeMapFunc:  fakeCassandraDcMapFunc,
}

func TestAccAlicloudCassandraDataCentersDataSourceNewDataCenter(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testAccCasDataSource_%d", rand)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCassandraDataCenterDataSourceConfigNewDataCenter(name, map[string]string{
			"name_regex": `"${alicloud_cassandra_data_center.dc_2.data_center_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCassandraDataCenterDataSourceConfigNewDataCenter(name, map[string]string{
			"name_regex": `"${alicloud_cassandra_data_center.dc_2.data_center_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCassandraDataCenterDataSourceConfigNewDataCenter(name, map[string]string{
			"ids": `["${alicloud_cassandra_data_center.dc_2.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCassandraDataCenterDataSourceConfigNewDataCenter(name, map[string]string{
			"ids": `["${alicloud_cassandra_data_center.dc_2.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCassandraDataCenterDataSourceConfigNewDataCenter(name, map[string]string{
			"name_regex": `"${alicloud_cassandra_data_center.dc_2.data_center_name}"`,
			"ids":        `["${alicloud_cassandra_data_center.dc_2.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCassandraDataCenterDataSourceConfigNewDataCenter(name, map[string]string{
			"name_regex": `"${alicloud_cassandra_data_center.dc_2.data_center_name}"`,
			"ids":        `["${alicloud_cassandra_data_center.dc_2.id}_fake"]`,
		}),
	}

	preCheck := func() {
	}
	checkCassandraDcInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
}

// new a cluster and a dataCenter config
func testAccCheckAlicloudCassandraDataCenterDataSourceConfigNewDataCenter(name string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_cassandra_zones" "default" {
}

data "alicloud_vpcs" "default" {
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)-1].id
}

resource "alicloud_vswitch" "this" {
  count = "${length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1}"
  name = "%[1]s"
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  availability_zone = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)-1].id
  cidr_block = "${cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)}"
}

data "alicloud_vswitches" "default_2" {
  vpc_id = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)-2].id
}

resource "alicloud_vswitch" "this_2" {
  count = "${length(data.alicloud_vswitches.default_2.ids) > 0 ? 0 : 1}"
  name = "tf_testAccCassandra_vpc_2"
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  availability_zone = data.alicloud_cassandra_zones.default.zones[length(data.alicloud_cassandra_zones.default.ids)-2].id
  cidr_block = "${cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 10)}"
}

resource "alicloud_cassandra_cluster" "default" {
  cluster_name = "%[1]s"
  data_center_name = "%[1]s"
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
  ip_white = "127.0.0.1"
}

resource "alicloud_cassandra_data_center" "dc_2" {
  cluster_id = "${alicloud_cassandra_cluster.default.id}"
  data_center_name = "%[1]s"
  auto_renew = false
  instance_type = "cassandra.c.large"
  node_count = "2"
  pay_type = "PayAsYouGo"
  vswitch_id = "${length(data.alicloud_vswitches.default_2.ids) > 0 ? data.alicloud_vswitches.default_2.ids[0] : alicloud_vswitch.this_2[0].id}"
  disk_size = "160"
  disk_type = "cloud_ssd"
}
`, name, strings.Join(pairs, "\n  "))
	return config
}
