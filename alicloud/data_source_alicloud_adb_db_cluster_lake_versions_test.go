package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAdbDbClusterLakeVersionsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.ADBDBClusterLakeVersionSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbDbClusterLakeVersionsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_adb_db_cluster_lake_version.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAdbDbClusterLakeVersionsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_adb_db_cluster_lake_version.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbDbClusterLakeVersionsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_adb_db_cluster_lake_version.default.id}"]`,
			"status": `"${alicloud_adb_db_cluster_lake_version.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudAdbDbClusterLakeVersionsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_adb_db_cluster_lake_version.default.id}_fake"]`,
			"status": `"Deleting"`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbDbClusterLakeVersionsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_adb_db_cluster_lake_version.default.id}"]`,
			"resource_group_id": `"${alicloud_adb_db_cluster_lake_version.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAdbDbClusterLakeVersionsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_adb_db_cluster_lake_version.default.id}_fake"]`,
			"resource_group_id": `"${alicloud_adb_db_cluster_lake_version.default.resource_group_id}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAdbDbClusterLakeVersionsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_adb_db_cluster_lake_version.default.id}"]`,
			"status":            `"${alicloud_adb_db_cluster_lake_version.default.status}"`,
			"resource_group_id": `"${alicloud_adb_db_cluster_lake_version.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAdbDbClusterLakeVersionsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_adb_db_cluster_lake_version.default.id}_fake"]`,
			"status":            `"Deleting"`,
			"resource_group_id": `"${alicloud_adb_db_cluster_lake_version.default.resource_group_id}_fake"`,
		}),
	}
	var existAlicloudAdbDbClusterLakeVersionsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"versions.#":                    "1",
			"versions.0.id":                 CHECKSET,
			"versions.0.db_cluster_id":      CHECKSET,
			"versions.0.commodity_code":     CHECKSET,
			"versions.0.connection_string":  CHECKSET,
			"versions.0.expire_time":        "",
			"versions.0.compute_resource":   "16ACU",
			"versions.0.db_cluster_version": "5.0",
			"versions.0.payment_type":       "PayAsYouGo",
			"versions.0.storage_resource":   "0ACU",
			"versions.0.status":             CHECKSET,
			"versions.0.vswitch_id":         CHECKSET,
			"versions.0.vpc_id":             CHECKSET,
			"versions.0.zone_id":            CHECKSET,
			"versions.0.resource_group_id":  CHECKSET,
			"versions.0.engine_version":     CHECKSET,
			"versions.0.engine":             CHECKSET,
			"versions.0.create_time":        CHECKSET,
			"versions.0.expired":            "",
			"versions.0.lock_mode":          CHECKSET,
			"versions.0.lock_reason":        "",
			"versions.0.port":               CHECKSET,
		}
	}
	var fakeAlicloudAdbDbClusterLakeVersionsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudAdbDbClusterLakeVersionsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_adb_db_cluster_lake_versions.default",
		existMapFunc: existAlicloudAdbDbClusterLakeVersionsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudAdbDbClusterLakeVersionsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudAdbDbClusterLakeVersionsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, resourceGroupIdConf, allConf)
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
