package dbfs_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudDBFSSnapshotsDataSource(t *testing.T) {
	t.Skip("This resource has been deprecated since v1.279.0 and will be removed in the future. See: https://help.aliyun.com/en/dbfs/")
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.DBFSSystemSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDbfsSnapshotsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dbfs_snapshot.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudDbfsSnapshotsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dbfs_snapshot.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDbfsSnapshotsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_dbfs_snapshot.default.snapshot_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudDbfsSnapshotsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_dbfs_snapshot.default.snapshot_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDbfsSnapshotsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_dbfs_snapshot.default.id}"]`,
			"status": `"accomplished"`,
		}),
		FakeConfig: testAccCheckAlicloudDbfsSnapshotsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_dbfs_snapshot.default.id}"]`,
			"status": `"failed"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDbfsSnapshotsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_dbfs_snapshot.default.id}"]`,
			"name_regex": `"${alicloud_dbfs_snapshot.default.snapshot_name}"`,
			"status":     `"accomplished"`,
		}),
		FakeConfig: testAccCheckAlicloudDbfsSnapshotsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_dbfs_snapshot.default.id}_fake"]`,
			"name_regex": `"${alicloud_dbfs_snapshot.default.snapshot_name}_fake"`,
			"status":     `"failed"`,
		}),
	}
	var existAlicloudDbfsSnapshotsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"names.#":                        "1",
			"snapshots.#":                    "1",
			"snapshots.0.description":        fmt.Sprintf("tf-testAccSnapshot-%d", rand),
			"snapshots.0.instance_id":        acctest.CHECKSET,
			"snapshots.0.retention_days":     "30",
			"snapshots.0.status":             "accomplished",
			"snapshots.0.category":           acctest.CHECKSET,
			"snapshots.0.create_time":        acctest.CHECKSET,
			"snapshots.0.last_modified_time": acctest.CHECKSET,
			"snapshots.0.progress":           acctest.CHECKSET,
			"snapshots.0.remain_time":        acctest.CHECKSET,
			"snapshots.0.id":                 acctest.CHECKSET,
			"snapshots.0.snapshot_id":        acctest.CHECKSET,
			"snapshots.0.snapshot_name":      fmt.Sprintf("tf-testAccSnapshot-%d", rand),
			"snapshots.0.snapshot_type":      acctest.CHECKSET,
			"snapshots.0.source_fs_size":     acctest.CHECKSET,
		}
	}
	var fakeAlicloudDbfsSnapshotsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudDbfsSnapshotsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_dbfs_snapshots.default",
		ExistMapFunc: existAlicloudDbfsSnapshotsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudDbfsSnapshotsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudDbfsSnapshotsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudDbfsSnapshotsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccSnapshot-%d"
}
locals {
  zone_id = "cn-hangzhou-i"
}
data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = local.zone_id
}

resource "alicloud_security_group" "default" {
  name        = var.name
  description = "tf test"
  vpc_id      = data.alicloud_vpcs.default.ids.0
}

data "alicloud_images" "default" {
  owners      = "system"
  name_regex  = "^centos_8"
  most_recent = true
}

resource "alicloud_instance" "default" {
  image_id          = data.alicloud_images.default.images[0].id
  instance_name     = var.name
  instance_type     = "ecs.g7se.large"
  availability_zone = local.zone_id
  vswitch_id        = data.alicloud_vswitches.default.ids[0]
  system_disk_category = "cloud_essd"
  security_groups = [
    alicloud_security_group.default.id
  ]
}

resource "alicloud_dbfs_instance" "default" {
  category          = "standard"
  zone_id           = alicloud_instance.default.availability_zone
  performance_level = "PL1"
  instance_name     = var.name
  size              = 100
}

resource "alicloud_dbfs_instance_attachment" "default" {
  ecs_id      = alicloud_instance.default.id
  instance_id = alicloud_dbfs_instance.default.id
}

resource "alicloud_dbfs_snapshot" "default" {
	depends_on = [alicloud_dbfs_instance_attachment.default]
	description = var.name
	instance_id = alicloud_dbfs_instance.default.id
	retention_days = 30
	snapshot_name = var.name
}

data "alicloud_dbfs_snapshots" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
