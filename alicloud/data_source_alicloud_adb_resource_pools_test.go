package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAdbResourcePoolsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.ADBSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbResourcePoolsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_adb_resource_pool.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAdbResourcePoolsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_adb_resource_pool.default.id}_fake"]`,
		}),
	}
	resourcePoolnameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbResourcePoolsDataSourceName(rand, map[string]string{
			"resource_pool_name": `"${alicloud_adb_resource_pool.default.resource_pool_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAdbResourcePoolsDataSourceName(rand, map[string]string{
			"resource_pool_name": `"${alicloud_adb_resource_pool.default.resource_pool_name}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbResourcePoolsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_adb_resource_pool.default.resource_pool_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAdbResourcePoolsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_adb_resource_pool.default.resource_pool_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbResourcePoolsDataSourceName(rand, map[string]string{
			"ids":                `["${alicloud_adb_resource_pool.default.id}"]`,
			"resource_pool_name": `"${alicloud_adb_resource_pool.default.resource_pool_name}"`,
			"name_regex":         `"${alicloud_adb_resource_pool.default.resource_pool_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAdbResourcePoolsDataSourceName(rand, map[string]string{
			"ids":                `["${alicloud_adb_resource_pool.default.id}_fake"]`,
			"resource_pool_name": `"${alicloud_adb_resource_pool.default.resource_pool_name}_fake"`,
			"name_regex":         `"${alicloud_adb_resource_pool.default.resource_pool_name}_fake"`,
		}),
	}
	var existAlicloudAdbResourcePoolsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"pools.#":                    "1",
			"names.#":                    "1",
			"pools.0.id":                 CHECKSET,
			"pools.0.create_time":        CHECKSET,
			"pools.0.db_cluster_id":      CHECKSET,
			"pools.0.node_num":           "2",
			"pools.0.resource_pool_name": fmt.Sprintf("TF-TESTACCRESOURCEPOOL%d", rand),
			"pools.0.query_type":         "batch",
		}
	}
	var fakeAlicloudAdbResourcePoolsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudAdbResourcePoolsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_adb_resource_pools.default",
		existMapFunc: existAlicloudAdbResourcePoolsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudAdbResourcePoolsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudAdbResourcePoolsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, resourcePoolnameConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudAdbResourcePoolsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "TF-TESTACCRESOURCEPOOL%d"
}

data "alicloud_adb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_adb_zones.default.zones[0].id
}

resource "alicloud_adb_db_cluster" "default" {
  db_cluster_category = "MixedStorage"
  mode                = "flexible"
  compute_resource    = "32Core128GB"
  payment_type        = "PayAsYouGo"
  vswitch_id          = data.alicloud_vswitches.default.ids[0]
  description         = var.name
  maintain_time       = "23:00Z-00:00Z"
  tags = {
    Created = "TF"
    For     = "acceptance-test-update"
  }
}

resource "alicloud_adb_resource_pool" "default" {
  db_cluster_id      = alicloud_adb_db_cluster.default.id
  resource_pool_name = var.name
  query_type         = "batch"
  node_num           = 2
}


data "alicloud_adb_resource_pools" "default" {	
	db_cluster_id = alicloud_adb_db_cluster.default.id
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
