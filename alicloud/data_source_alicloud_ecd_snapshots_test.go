package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEcdSnapshotsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.EcdSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdSnapshotsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_snapshot.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcdSnapshotsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ecd_snapshot.default.id}_fake"]`,
		}),
	}
	desktopIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdSnapshotsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_snapshot.default.id}"]`,
			"desktop_id": `"${alicloud_ecd_snapshot.default.desktop_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdSnapshotsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ecd_snapshot.default.id}"]`,
			"desktop_id": `"${alicloud_ecd_snapshot.default.desktop_id}_fake"`,
		}),
	}
	snapshotIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdSnapshotsDataSourceName(rand, map[string]string{
			"snapshot_id": `"${alicloud_ecd_snapshot.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdSnapshotsDataSourceName(rand, map[string]string{
			"snapshot_id": `"${alicloud_ecd_snapshot.default.id}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdSnapshotsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_snapshot.default.snapshot_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdSnapshotsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ecd_snapshot.default.snapshot_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcdSnapshotsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_ecd_snapshot.default.id}"]`,
			"desktop_id":  `"${alicloud_ecd_snapshot.default.desktop_id}"`,
			"snapshot_id": `"${alicloud_ecd_snapshot.default.id}"`,
			"name_regex":  `"${alicloud_ecd_snapshot.default.snapshot_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcdSnapshotsDataSourceName(rand, map[string]string{
			"ids":         `["${alicloud_ecd_snapshot.default.id}_fake"]`,
			"desktop_id":  `"${alicloud_ecd_snapshot.default.desktop_id}_fake"`,
			"snapshot_id": `"${alicloud_ecd_snapshot.default.id}_fake"`,
			"name_regex":  `"${alicloud_ecd_snapshot.default.snapshot_name}_fake"`,
		}),
	}
	var existAlicloudEcdSnapshotsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"names.#":                      "1",
			"snapshots.#":                  "1",
			"snapshots.0.id":               CHECKSET,
			"snapshots.0.create_time":      CHECKSET,
			"snapshots.0.description":      fmt.Sprintf("tf-testaccecdsnapshot%d", rand),
			"snapshots.0.desktop_id":       CHECKSET,
			"snapshots.0.progress":         CHECKSET,
			"snapshots.0.remain_time":      CHECKSET,
			"snapshots.0.snapshot_id":      CHECKSET,
			"snapshots.0.snapshot_name":    fmt.Sprintf("tf-testaccecdsnapshot%d", rand),
			"snapshots.0.snapshot_type":    CHECKSET,
			"snapshots.0.source_disk_size": CHECKSET,
			"snapshots.0.source_disk_type": "SYSTEM",
			"snapshots.0.status":           CHECKSET,
		}
	}
	var fakeAlicloudEcdSnapshotsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEcdSnapshotsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecd_snapshots.default",
		existMapFunc: existAlicloudEcdSnapshotsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcdSnapshotsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEcdSnapshotsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, desktopIdConf, snapshotIdConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudEcdSnapshotsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testaccecdsnapshot%d"
}

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block             = "172.16.0.0/12"
  desktop_access_type    = "Internet"
  office_site_name       = var.name
  enable_internet_access = false
}

data "alicloud_ecd_bundles" "default" {
  bundle_type = "SYSTEM"
}

resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = var.name
  clipboard         = "readwrite"
  local_drive       = "read"
  authorize_access_policy_rules {
    description = "example_value"
    cidr_ip     = "1.2.3.4/24"
  }
  authorize_security_policy_rules {
    type        = "inflow"
    policy      = "accept"
    description = "example_value"
    port_range  = "80/80"
    ip_protocol = "TCP"
    priority    = "1"
    cidr_ip     = "0.0.0.0/0"
  }
}

resource "alicloud_ecd_desktop" "default" {
  office_site_id  = alicloud_ecd_simple_office_site.default.id
  policy_group_id = alicloud_ecd_policy_group.default.id
  bundle_id       = data.alicloud_ecd_bundles.default.bundles.0.id
  desktop_name    = var.name
}

resource "alicloud_ecd_snapshot" "default" {
	description = var.name
	desktop_id = alicloud_ecd_desktop.default.id
	snapshot_name = var.name
	source_disk_type = "SYSTEM"
}

data "alicloud_ecd_snapshots" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
