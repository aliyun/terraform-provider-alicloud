package hbr_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudHbrSnapshotsDataSource(t *testing.T) {
	defer acctest.CheckoutAccount(t, false)
	acctest.CheckoutAccount(t, true)
	acctest.CheckoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	rand := tfacctest.RandIntRange(1000000, 9999999)

	ecsBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type": `"ECS_FILE"`,
			"vault_id":    `"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}"`,
			"instance_id": `"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type": `"ECS_FILE"`,
			"vault_id":    `"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}_fake"`,
			"instance_id": `"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}_fake"`,
		}),
	}

	ossBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type": `"OSS"`,
			"vault_id":    `"${data.alicloud_hbr_oss_backup_plans.default.plans.0.vault_id}"`,
			"bucket":      `"${data.alicloud_hbr_oss_backup_plans.default.plans.0.bucket}"`,
		}),
		FakeConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type": `"OSS"`,
			"vault_id":    `"${data.alicloud_hbr_oss_backup_plans.default.plans.0.vault_id}_fake"`,
			"bucket":      `"${data.alicloud_hbr_oss_backup_plans.default.plans.0.bucket}_fake"`,
		}),
	}

	nasBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type":    `"NAS"`,
			"vault_id":       `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}"`,
			"file_system_id": `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}"`,
			"create_time":    `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time}"`,
		}),
		FakeConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type":    `"NAS"`,
			"vault_id":       `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}_fake"`,
			"file_system_id": `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}"`,
			"create_time":    `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time}"`,
		}),
	}

	statusBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type": `"ECS_FILE"`,
			"vault_id":    `"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}"`,
			"instance_id": `"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}"`,
			"status":      `"COMPLETE"`,
		}),
		FakeConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type": `"ECS_FILE"`,
			"vault_id":    `"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}"`,
			"instance_id": `"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}"`,
			"status":      `"PARTIAL_COMPLETE"`,
		}),
	}

	completeTimeBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type":           `"NAS"`,
			"vault_id":              `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}"`,
			"file_system_id":        `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}"`,
			"create_time":           `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time}"`,
			"complete_time":         `"2021-08-23T14:17:15CST"`,
			"complete_time_checker": `"GREATER_THAN_OR_EQUAL"`,
		}),
		FakeConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type":           `"NAS"`,
			"vault_id":              `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}"`,
			"file_system_id":        `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}"`,
			"create_time":           `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time}"`,
			"complete_time":         `"2021-08-23T14:17:15CST"`,
			"complete_time_checker": `"LESS_THAN"`,
		}),
	}

	betweenTimeBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type":           `"NAS"`,
			"vault_id":              `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}"`,
			"file_system_id":        `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}"`,
			"create_time":           `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time}"`,
			"complete_time":         `"2021-08-20T14:17:15CST,2025-08-26T14:17:15CST"`,
			"complete_time_checker": `"BETWEEN"`,
		}),
		FakeConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type":           `"NAS"`,
			"vault_id":              `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}"`,
			"file_system_id":        `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}"`,
			"create_time":           `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time}"`,
			"complete_time":         `"2021-07-20T14:17:15CST,2021-07-24T14:17:15CST"`,
			"complete_time_checker": `"BETWEEN"`,
		}),
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.HbrSupportRegions)
	}

	HbrSnapshotCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, ecsBackupConf, ossBackupConf, nasBackupConf, statusBackupConf, completeTimeBackupConf, betweenTimeBackupConf)
}

func testAccCheckAlicloudHbrSnapshotSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAcc%d"
}

data "alicloud_hbr_ecs_backup_plans" "default" {}

data "alicloud_hbr_oss_backup_plans" "default" {}

data "alicloud_hbr_nas_backup_plans" "default" {}

data "alicloud_hbr_snapshots" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

var existHbrSnapshotMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"snapshots.#": acctest.CHECKSET,
	}
}

var fakeHbrSnapshotMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"snapshots.#": "0",
	}
}

var HbrSnapshotCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_hbr_snapshots.default",
	ExistMapFunc: existHbrSnapshotMapFunc,
	FakeMapFunc:  fakeHbrSnapshotMapFunc,
}
