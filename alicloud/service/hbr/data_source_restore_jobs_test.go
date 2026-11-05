package hbr_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudHbrRestoreJobsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)

	ecsBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrRestoreJobSourceConfig(rand, map[string]string{
			"restore_type":       `"ECS_FILE"`,
			"vault_id":           `["${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}"]`,
			"target_instance_id": `["${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudHbrRestoreJobSourceConfig(rand, map[string]string{
			"restore_type":       `"ECS_FILE"`,
			"vault_id":           `["${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}_fake"]`,
			"target_instance_id": `["${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}_fake"]`,
		}),
	}

	ossBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrRestoreJobSourceConfig(rand, map[string]string{
			"restore_type":  `"OSS"`,
			"vault_id":      `["${data.alicloud_hbr_oss_backup_plans.default.plans.0.vault_id}"]`,
			"target_bucket": `["${data.alicloud_hbr_oss_backup_plans.default.plans.0.bucket}"]`,
		}),
		FakeConfig: testAccCheckAlicloudHbrRestoreJobSourceConfig(rand, map[string]string{
			"restore_type":  `"OSS"`,
			"vault_id":      `["${data.alicloud_hbr_oss_backup_plans.default.plans.0.vault_id}_fake"]`,
			"target_bucket": `["${data.alicloud_hbr_oss_backup_plans.default.plans.0.bucket}_fake"]`,
		}),
	}

	nasBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrRestoreJobSourceConfig(rand, map[string]string{
			"restore_type":          `"NAS"`,
			"vault_id":              `["${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}"]`,
			"target_file_system_id": `["${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudHbrRestoreJobSourceConfig(rand, map[string]string{
			"restore_type":          `"NAS"`,
			"vault_id":              `["${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}_fake"]`,
			"target_file_system_id": `["${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}_fake"]`,
		}),
	}

	statusBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrRestoreJobSourceConfig(rand, map[string]string{
			"restore_type":       `"ECS_FILE"`,
			"vault_id":           `["${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}"]`,
			"target_instance_id": `["${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}"]`,
			"status":             `"COMPLETE"`,
		}),
		FakeConfig: testAccCheckAlicloudHbrRestoreJobSourceConfig(rand, map[string]string{
			"restore_type":       `"ECS_FILE"`,
			"vault_id":           `["${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}_fake"]`,
			"target_instance_id": `["${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}_fake"]`,
			"status":             `"CANCELING"`,
		}),
	}

	restoreBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrRestoreJobSourceConfig(rand, map[string]string{
			"restore_type": `"NAS"`,
			"restore_id":   `[split(":", alicloud_hbr_restore_job.default.id)[0]]`,
		}),
		FakeConfig: testAccCheckAlicloudHbrRestoreJobSourceConfig(rand, map[string]string{
			"restore_type": `"NAS"`,
			"restore_id":   `["fakeId"]`,
		}),
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.HbrSupportRegions)
	}
	HbrRestoreJobCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, ecsBackupConf, ossBackupConf, nasBackupConf, statusBackupConf, restoreBackupConf)
}

func testAccCheckAlicloudHbrRestoreJobSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAcc%d"
}

data "alicloud_hbr_ecs_backup_plans" "default" {
}

data "alicloud_hbr_oss_backup_plans" "default" {
}

data "alicloud_hbr_nas_backup_plans" "default" {
}

data "alicloud_hbr_snapshots" "nas_snapshots" {
    source_type     = "NAS"
    vault_id        =  data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id
    file_system_id  =  data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id
    create_time     =  data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time
}

resource "alicloud_hbr_restore_job" "default" {
    snapshot_hash =         data.alicloud_hbr_snapshots.nas_snapshots.snapshots.0.snapshot_hash
    vault_id =              data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id
    source_type =          "NAS"
    restore_type =         "NAS"
    snapshot_id =           data.alicloud_hbr_snapshots.nas_snapshots.snapshots.0.snapshot_id
    target_file_system_id = data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id
    target_create_time =    data.alicloud_hbr_nas_backup_plans.default.plans.0.create_time
    target_path =           "/"
    options = <<EOF
    {"includes":[], "excludes":[]}
    EOF
}

data "alicloud_hbr_restore_jobs" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

var existHbrRestoreJobMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"jobs.#": acctest.CHECKSET,
	}
}

var fakeHbrRestoreJobMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"jobs.#": "0",
	}
}

var HbrRestoreJobCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_hbr_restore_jobs.default",
	ExistMapFunc: existHbrRestoreJobMapFunc,
	FakeMapFunc:  fakeHbrRestoreJobMapFunc,
}
