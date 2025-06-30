// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudHbrUdmSnapshotDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
	rand := acctest.RandIntRange(1000000, 9999999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrUdmSnapshotSourceConfig(rand, map[string]string{
			"source_type": `"UDM_ECS"`,
			"start_time":  `"1642057551"`,
			"end_time":    `"1750927687"`,
			"instance_id": `"i-gw862uxmqlzyoudv57sb"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrUdmSnapshotSourceConfig(rand, map[string]string{
			"source_type": `"UDM_ECS"`,
			"start_time":  `"1642057551"`,
			"end_time":    `"1642057751"`,
			"instance_id": `"i-gw862uxmqlzyoudv57sb"`,
		}),
	}

	HbrUdmSnapshotCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

var existHbrUdmSnapshotMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"snapshots.#":                 "1",
		"snapshots.0.instance_id":     CHECKSET,
		"snapshots.0.source_type":     CHECKSET,
		"snapshots.0.udm_snapshot_id": CHECKSET,
		"snapshots.0.job_id":          CHECKSET,
	}
}

var fakeHbrUdmSnapshotMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"snapshots.#": "0",
	}
}

var HbrUdmSnapshotCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_hbr_udm_snapshots.default",
	existMapFunc: existHbrUdmSnapshotMapFunc,
	fakeMapFunc:  fakeHbrUdmSnapshotMapFunc,
}

func testAccCheckAlicloudHbrUdmSnapshotSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccHbrUdmSnapshot%d"
}

data "alicloud_hbr_udm_snapshots" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
