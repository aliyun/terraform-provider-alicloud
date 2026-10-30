// Package hbr_test. This file is generated automatically. Please do not modify it manually, thank you!
package hbr_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudHbrUdmSnapshotDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrUdmSnapshotSourceConfig(rand, map[string]string{
			"source_type": `"UDM_ECS"`,
			"start_time":  `"1642057551"`,
			"end_time":    `"1750927687"`,
			"instance_id": `"i-bp1bo1iozgsxrzfbktkp"`, // should be a valid instance id and not hard coded
		}),
		FakeConfig: testAccCheckAlicloudHbrUdmSnapshotSourceConfig(rand, map[string]string{
			"source_type": `"UDM_ECS"`,
			"start_time":  `"1642057551"`,
			"end_time":    `"1642057751"`,
			"instance_id": `"i-gw862uxmqlzyoudv57sb"`,
		}),
	}

	HbrUdmSnapshotCheckInfo.DataSourceTestCheck(t, rand, allConf)
}

var existHbrUdmSnapshotMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"snapshots.#":                 "1",
		"snapshots.0.instance_id":     acctest.CHECKSET,
		"snapshots.0.source_type":     acctest.CHECKSET,
		"snapshots.0.udm_snapshot_id": acctest.CHECKSET,
		"snapshots.0.job_id":          acctest.CHECKSET,
	}
}

var fakeHbrUdmSnapshotMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"snapshots.#": "0",
	}
}

var HbrUdmSnapshotCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_hbr_udm_snapshots.default",
	ExistMapFunc: existHbrUdmSnapshotMapFunc,
	FakeMapFunc:  fakeHbrUdmSnapshotMapFunc,
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
