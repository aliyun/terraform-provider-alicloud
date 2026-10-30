package hbr_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudHbrOtsSnapshotsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 99999)
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrOtsSnapshotsSourceConfig(rand, map[string]string{}),
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.HbrSupportRegions)
	}

	HbrOtsSnapshotsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
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
		"snapshots.0.id":            acctest.CHECKSET,
		"snapshots.0.vault_id":      acctest.CHECKSET,
		"snapshots.0.start_time":    acctest.CHECKSET,
		"snapshots.0.retention":     acctest.CHECKSET,
		"snapshots.0.snapshot_hash": acctest.CHECKSET,
	}
}

var fakeHbrOtsSnapshotsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"snapshots.#": "0",
	}
}

var HbrOtsSnapshotsCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_hbr_ots_snapshots.default",
	ExistMapFunc: existHbrOtsSnapshotsMapFunc,
	FakeMapFunc:  fakeHbrOtsSnapshotsMapFunc,
}
