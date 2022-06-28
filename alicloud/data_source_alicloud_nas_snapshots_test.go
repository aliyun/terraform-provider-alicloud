package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudNASSnapshotsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.NASSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasSnapshotsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nas_snapshot.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNasSnapshotsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_nas_snapshot.default.id}_fake"]`,
		}),
	}
	fileSystemIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasSnapshotsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_nas_snapshot.default.id}"]`,
			"file_system_id": `"${alicloud_nas_snapshot.default.file_system_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudNasSnapshotsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_nas_snapshot.default.id}"]`,
			"file_system_id": `"${alicloud_nas_snapshot.default.file_system_id}_fake"`,
		}),
	}
	snapshotNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasSnapshotsDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_nas_snapshot.default.id}"]`,
			"snapshot_name": `"${alicloud_nas_snapshot.default.snapshot_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudNasSnapshotsDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_nas_snapshot.default.id}"]`,
			"snapshot_name": `"${alicloud_nas_snapshot.default.snapshot_name}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasSnapshotsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_nas_snapshot.default.snapshot_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudNasSnapshotsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_nas_snapshot.default.snapshot_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasSnapshotsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nas_snapshot.default.id}"]`,
			"status": `"accomplished"`,
		}),
		fakeConfig: testAccCheckAlicloudNasSnapshotsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_nas_snapshot.default.id}"]`,
			"status": `"failed"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasSnapshotsDataSourceName(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_snapshot.default.file_system_id}"`,
			"ids":            `["${alicloud_nas_snapshot.default.id}"]`,
			"name_regex":     `"${alicloud_nas_snapshot.default.snapshot_name}"`,
			"snapshot_name":  `"${alicloud_nas_snapshot.default.snapshot_name}"`,
			"status":         `"accomplished"`,
		}),
		fakeConfig: testAccCheckAlicloudNasSnapshotsDataSourceName(rand, map[string]string{
			"file_system_id": `"${alicloud_nas_snapshot.default.file_system_id}_fake"`,
			"ids":            `["${alicloud_nas_snapshot.default.id}_fake"]`,
			"name_regex":     `"${alicloud_nas_snapshot.default.snapshot_name}_fake"`,
			"snapshot_name":  `"${alicloud_nas_snapshot.default.snapshot_name}_fake"`,
			"status":         `"failed"`,
		}),
	}
	var existAlicloudNasSnapshotsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"names.#":                                "1",
			"snapshots.#":                            "1",
			"snapshots.0.create_time":                CHECKSET,
			"snapshots.0.description":                fmt.Sprintf("tf-testAccSnapshot-%d", rand),
			"snapshots.0.encrypt_type":               CHECKSET,
			"snapshots.0.progress":                   CHECKSET,
			"snapshots.0.remain_time":                CHECKSET,
			"snapshots.0.retention_days":             CHECKSET,
			"snapshots.0.id":                         CHECKSET,
			"snapshots.0.snapshot_id":                CHECKSET,
			"snapshots.0.snapshot_name":              fmt.Sprintf("tf-testAccSnapshot-%d", rand),
			"snapshots.0.source_file_system_id":      CHECKSET,
			"snapshots.0.source_file_system_size":    CHECKSET,
			"snapshots.0.source_file_system_version": CHECKSET,
			"snapshots.0.status":                     CHECKSET,
		}
	}
	var fakeAlicloudNasSnapshotsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudNasSnapshotsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_nas_snapshots.default",
		existMapFunc: existAlicloudNasSnapshotsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudNasSnapshotsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudNasSnapshotsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, fileSystemIdConf, snapshotNameConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudNasSnapshotsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccSnapshot-%d"
}

data "alicloud_nas_zones" "default" {
  file_system_type = "extreme"
}

locals {
  count_size = length(data.alicloud_nas_zones.default.zones)
}

resource "alicloud_nas_file_system" "default" {
  file_system_type = "extreme"
  protocol_type    = "NFS"
  zone_id          = data.alicloud_nas_zones.default.zones[local.count_size - 1].zone_id
  storage_type     = "standard"
  description      = var.name
  capacity         = 100
}

resource "alicloud_nas_snapshot" "default" {
  file_system_id = alicloud_nas_file_system.default.id
  description    = var.name
  retention_days = 20
  snapshot_name  = var.name
}

data "alicloud_nas_snapshots" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
