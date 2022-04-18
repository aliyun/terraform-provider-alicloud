package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudHBROtsSnapshotsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 99999)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrOtsSnapshotsSourceConfig(rand, map[string]string{}),
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.HbrSupportRegions)
	}

	HbrOtsSnapshotsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
}

func testAccCheckAlicloudHbrOtsSnapshotsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_hbr_ots_snapshots" "default" {
%s
}
`, strings.Join(pairs, "\n   "))
	return config
}

var existHbrOtsSnapshotsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"snapshots.#":               "1",
		"snapshots.0.status":        "COMPLETE",
		"snapshots.0.source_type":   "OTS_TABLE",
		"snapshots.0.backup_type":   "COMPLETE",
		"snapshots.0.id":            CHECKSET,
		"snapshots.0.vault_id":      CHECKSET,
		"snapshots.0.start_time":    CHECKSET,
		"snapshots.0.retention":     CHECKSET,
		"snapshots.0.snapshot_hash": CHECKSET,
	}
}

var fakeHbrOtsSnapshotsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"snapshots.#": "0",
	}
}

var HbrOtsSnapshotsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_hbr_ots_backup_snapshots.default",
	existMapFunc: existHbrOtsSnapshotsMapFunc,
	fakeMapFunc:  fakeHbrOtsSnapshotsMapFunc,
}
