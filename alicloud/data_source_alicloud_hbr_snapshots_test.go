package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudHBRSnapshotsDataSource(t *testing.T) {
	defer checkoutAccount(t, false)
	checkoutAccount(t, true)
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	rand := acctest.RandIntRange(1000000, 9999999)

	ecsBackupConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type": `"ECS_FILE"`,
			"vault_id":    `"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}"`,
			"instance_id": `"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type": `"ECS_FILE"`,
			"vault_id":    `"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}_fake"`,
			"instance_id": `"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}_fake"`,
		}),
	}

	ossBackupConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type": `"OSS"`,
			"vault_id":    `"${data.alicloud_hbr_oss_backup_plans.default.plans.0.vault_id}"`,
			"bucket":      `"${data.alicloud_hbr_oss_backup_plans.default.plans.0.bucket}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type": `"OSS"`,
			"vault_id":    `"${data.alicloud_hbr_oss_backup_plans.default.plans.0.vault_id}_fake"`,
			"bucket":      `"${data.alicloud_hbr_oss_backup_plans.default.plans.0.bucket}_fake"`,
		}),
	}

	nasBackupConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type":    `"NAS"`,
			"vault_id":       `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}"`,
			"file_system_id": `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}"`,
			"create_time":    `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type":    `"NAS"`,
			"vault_id":       `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}_fake"`,
			"file_system_id": `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}"`,
			"create_time":    `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time}"`,
		}),
	}

	statusBackupConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type": `"ECS_FILE"`,
			"vault_id":    `"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}"`,
			"instance_id": `"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}"`,
			"status":      `"COMPLETE"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type": `"ECS_FILE"`,
			"vault_id":    `"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}"`,
			"instance_id": `"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}"`,
			"status":      `"PARTIAL_COMPLETE"`,
		}),
	}

	completeTimeBackupConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type":           `"NAS"`,
			"vault_id":              `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}"`,
			"file_system_id":        `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}"`,
			"create_time":           `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time}"`,
			"complete_time":         `"2021-08-23T14:17:15CST"`,
			"complete_time_checker": `"GREATER_THAN_OR_EQUAL"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type":           `"NAS"`,
			"vault_id":              `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}"`,
			"file_system_id":        `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}"`,
			"create_time":           `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time}"`,
			"complete_time":         `"2021-08-23T14:17:15CST"`,
			"complete_time_checker": `"LESS_THAN"`,
		}),
	}

	betweenTimeBackupConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type":           `"NAS"`,
			"vault_id":              `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}"`,
			"file_system_id":        `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}"`,
			"create_time":           `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time}"`,
			"complete_time":         `"2021-08-20T14:17:15CST,2025-08-26T14:17:15CST"`,
			"complete_time_checker": `"BETWEEN"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrSnapshotSourceConfig(rand, map[string]string{
			"source_type":           `"NAS"`,
			"vault_id":              `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}"`,
			"file_system_id":        `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}"`,
			"create_time":           `"${data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time}"`,
			"complete_time":         `"2021-07-20T14:17:15CST,2021-07-24T14:17:15CST"`,
			"complete_time_checker": `"BETWEEN"`,
		}),
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.HbrSupportRegions)
	}

	HbrSnapshotCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, ecsBackupConf, ossBackupConf, nasBackupConf, statusBackupConf, completeTimeBackupConf, betweenTimeBackupConf)
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

data "alicloud_hbr_ecs_backup_plans" "default" {
    name_regex = "plan-tf-used-dont-delete"
}

data "alicloud_hbr_oss_backup_plans" "default" {
	name_regex = "plan-tf-used-dont-delete"
}

data "alicloud_hbr_nas_backup_plans" "default" {
	name_regex = "plan-tf-used-dont-delete"
}

data "alicloud_hbr_snapshots" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

var existHbrSnapshotMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"snapshots.#":    CHECKSET,
		"snapshots.0.id": CHECKSET,
	}
}

var fakeHbrSnapshotMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"snapshots.#": "0",
	}
}

var HbrSnapshotCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_hbr_snapshots.default",
	existMapFunc: existHbrSnapshotMapFunc,
	fakeMapFunc:  fakeHbrSnapshotMapFunc,
}
