package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDfsMountPointsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.DfsSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDfsMountPointsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dfs_mount_point.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDfsMountPointsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dfs_mount_point.default.id}_fake"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDfsMountPointsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_dfs_mount_point.default.id}"]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudDfsMountPointsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_dfs_mount_point.default.id}"]`,
			"status": `"Inactive"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDfsMountPointsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_dfs_mount_point.default.id}"]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudDfsMountPointsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_dfs_mount_point.default.id}_fake"]`,
			"status": `"Inactive"`,
		}),
	}
	var existAlicloudDfsMountPointsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"points.#":                "1",
			"points.0.description":    fmt.Sprintf("tf-testAccMountPoint-%d", rand),
			"points.0.file_system_id": CHECKSET,
			"points.0.network_type":   "VPC",
			"points.0.vswitch_id":     CHECKSET,
		}
	}
	var fakeAlicloudDfsMountPointsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudDfsMountPointsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_dfs_mount_points.default",
		existMapFunc: existAlicloudDfsMountPointsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudDfsMountPointsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudDfsMountPointsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, allConf)
}
func testAccCheckAlicloudDfsMountPointsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccMountPoint-%d"
}

data "alicloud_dfs_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_dfs_zones.default.zones.0.zone_id
}

resource "alicloud_vswitch" "default" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id      = data.alicloud_dfs_zones.default.zones.0.zone_id
  vswitch_name = var.name
}

resource "alicloud_dfs_file_system" "default" {
  storage_type     = data.alicloud_dfs_zones.default.zones.0.options.0.storage_type
  zone_id          = data.alicloud_dfs_zones.default.zones.0.zone_id
  protocol_type    = "HDFS"
  description      = var.name
  file_system_name = var.name
  throughput_mode  = "Standard"
  space_capacity   = "1024"
}

resource "alicloud_dfs_access_group" "default" {
  network_type      = "VPC"
  access_group_name = var.name
  description       = var.name
}

resource "alicloud_dfs_mount_point" "default" {
  description     = var.name
  vpc_id          = data.alicloud_vpcs.default.ids.0
  file_system_id  = alicloud_dfs_file_system.default.id
  access_group_id = alicloud_dfs_access_group.default.id
  network_type    = "VPC"
  vswitch_id      = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.default.*.id, [""])[0]
}

data "alicloud_dfs_mount_points" "default" {	
	file_system_id = alicloud_dfs_file_system.default.id 
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
