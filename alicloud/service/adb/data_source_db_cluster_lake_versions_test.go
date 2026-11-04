// Package adb_test.
package adb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudAdbDbClusterLakeVersionsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.ADBDBClusterLakeVersionSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAdbDbClusterLakeVersionsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_adb_db_cluster_lake_version.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudAdbDbClusterLakeVersionsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_adb_db_cluster_lake_version.default.id}_fake"]`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAdbDbClusterLakeVersionsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_adb_db_cluster_lake_version.default.id}"]`,
			"status": `"${alicloud_adb_db_cluster_lake_version.default.status}"`,
		}),
		FakeConfig: testAccCheckAlicloudAdbDbClusterLakeVersionsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_adb_db_cluster_lake_version.default.id}_fake"]`,
			"status": `"Deleting"`,
		}),
	}
	resourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAdbDbClusterLakeVersionsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_adb_db_cluster_lake_version.default.id}"]`,
			"resource_group_id": `"${alicloud_adb_db_cluster_lake_version.default.resource_group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudAdbDbClusterLakeVersionsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_adb_db_cluster_lake_version.default.id}_fake"]`,
			"resource_group_id": `"${alicloud_adb_db_cluster_lake_version.default.resource_group_id}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAdbDbClusterLakeVersionsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_adb_db_cluster_lake_version.default.id}"]`,
			"status":            `"${alicloud_adb_db_cluster_lake_version.default.status}"`,
			"resource_group_id": `"${alicloud_adb_db_cluster_lake_version.default.resource_group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudAdbDbClusterLakeVersionsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_adb_db_cluster_lake_version.default.id}_fake"]`,
			"status":            `"Deleting"`,
			"resource_group_id": `"${alicloud_adb_db_cluster_lake_version.default.resource_group_id}_fake"`,
		}),
	}
	var existAlicloudAdbDbClusterLakeVersionsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"versions.#":                    "1",
			"versions.0.id":                 acctest.CHECKSET,
			"versions.0.db_cluster_id":      acctest.CHECKSET,
			"versions.0.commodity_code":     acctest.CHECKSET,
			"versions.0.connection_string":  acctest.CHECKSET,
			"versions.0.expire_time":        "",
			"versions.0.compute_resource":   "16ACU",
			"versions.0.db_cluster_version": "5.0",
			"versions.0.payment_type":       "PayAsYouGo",
			"versions.0.storage_resource":   "0ACU",
			"versions.0.status":             acctest.CHECKSET,
			"versions.0.vswitch_id":         acctest.CHECKSET,
			"versions.0.vpc_id":             acctest.CHECKSET,
			"versions.0.zone_id":            acctest.CHECKSET,
			"versions.0.resource_group_id":  acctest.CHECKSET,
			"versions.0.engine_version":     acctest.CHECKSET,
			"versions.0.engine":             acctest.CHECKSET,
			"versions.0.create_time":        acctest.CHECKSET,
			"versions.0.expired":            "",
			"versions.0.lock_mode":          acctest.CHECKSET,
			"versions.0.lock_reason":        "",
			"versions.0.port":               acctest.CHECKSET,
		}
	}
	var fakeAlicloudAdbDbClusterLakeVersionsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudAdbDbClusterLakeVersionsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_adb_db_cluster_lake_versions.default",
		ExistMapFunc: existAlicloudAdbDbClusterLakeVersionsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudAdbDbClusterLakeVersionsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudAdbDbClusterLakeVersionsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, resourceGroupIdConf, allConf)
}
func testAccCheckAlicloudAdbDbClusterLakeVersionsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccDBClusterLakeVersion-%d"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_adb_zones" "default" {
}

data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_adb_zones.default.ids.0
}

resource "alicloud_adb_db_cluster_lake_version" "default" {
	compute_resource = "16ACU"
	db_cluster_version = "5.0"
	payment_type = "PayAsYouGo"
	storage_resource = "0ACU"
	vswitch_id = data.alicloud_vswitches.default.ids.0
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_adb_zones.default.ids.0
}

data "alicloud_adb_db_cluster_lake_versions" "default" {
	enable_details = true
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
