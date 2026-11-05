package hbr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudHbrBackupJobsDataSource(t *testing.T) {
	acctest.CheckoutAccount(t, true)
	defer acctest.CheckoutAccount(t, false)
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_hbr_backup_jobs.default"
	name := fmt.Sprintf("tf-testAccHbrBackupJobTest%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceHbrBackupJobSourceConfig)

	planIdBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"source_type": "ECS_FILE",
			"filter": []map[string]interface{}{
				{
					"key":      "PlanId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.id}"},
				},
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"source_type": "ECS_FILE",
			"filter": []map[string]interface{}{
				{
					"key":      "PlanId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.id}_fake"},
				},
			},
		}),
	}

	jobIdBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"source_type": "ECS_FILE",
			"filter": []map[string]interface{}{
				{
					"key":      "JobId",
					"operator": "IN",
					"values":   []string{"job-000e18wmkcpdqq0n5hna"},
				},
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"source_type": "ECS_FILE",
			"filter": []map[string]interface{}{
				{
					"key":      "JobId",
					"operator": "IN",
					"values":   []string{"job-000e18wmkcpdqq0n5hna_fake"},
				},
			},
		}),
	}

	ecsBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"source_type": "ECS_FILE",
			"filter": []map[string]interface{}{
				{
					"key":      "VaultId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}"},
				},
				{
					"key":      "InstanceId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}"},
				},
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"source_type": "ECS_FILE",
			"filter": []map[string]interface{}{
				{
					"key":      "VaultId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}_fake"},
				},
				{
					"key":      "InstanceId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}_fake"},
				},
			},
		}),
	}

	ossBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"source_type": "OSS",
			"filter": []map[string]interface{}{
				{
					"key":      "VaultId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_oss_backup_plans.default.plans.0.vault_id}"},
				},
				{
					"key":      "Bucket",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_oss_backup_plans.default.plans.0.bucket}"},
				},
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"source_type": "OSS",
			"filter": []map[string]interface{}{
				{
					"key":      "VaultId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_oss_backup_plans.default.plans.0.vault_id}"},
				},
				{
					"key":      "Bucket",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_oss_backup_plans.default.plans.0.bucket}_fake"},
				},
			},
		}),
	}

	nasBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"source_type": "NAS",
			"filter": []map[string]interface{}{
				{
					"key":      "VaultId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}"},
				},
				{
					"key":      "FileSystemId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}"},
				},
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"source_type": "NAS",
			"filter": []map[string]interface{}{
				{
					"key":      "VaultId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}"},
				},
				{
					"key":      "FileSystemId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}_fake"},
				},
			},
		}),
	}

	completeTimeBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"source_type": "NAS",
			"filter": []map[string]interface{}{
				{
					"key":      "VaultId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}"},
				},
				{
					"key":      "FileSystemId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}"},
				},
				{
					"key":      "CompleteTime",
					"operator": "GREATER_THAN",
					"values":   []string{"2021-08-23T14:17:15CST"},
				},
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"source_type": "NAS",
			"filter": []map[string]interface{}{
				{
					"key":      "VaultId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}"},
				},
				{
					"key":      "FileSystemId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}"},
				},
				{
					"key":      "CompleteTime",
					"operator": "LESS_THAN",
					"values":   []string{"2020-08-23T14:17:15CST"},
				},
			},
		}),
	}

	betweenTimeBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"source_type": "NAS",
			"filter": []map[string]interface{}{
				{
					"key":      "VaultId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}"},
				},
				{
					"key":      "FileSystemId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}"},
				},
				{
					"key":      "CompleteTime",
					"operator": "BETWEEN",
					"values":   []string{"2021-08-23T14:17:15CST", "2022-10-20T14:17:15CST"},
				},
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"source_type": "NAS",
			"filter": []map[string]interface{}{
				{
					"key":      "VaultId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_nas_backup_plans.default.plans.0.vault_id}"},
				},
				{
					"key":      "FileSystemId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_nas_backup_plans.default.plans.0.file_system_id}"},
				},
				{
					"key":      "CompleteTime",
					"operator": "BETWEEN",
					"values":   []string{"2019-08-23T14:17:15CST", "2020-08-23T14:17:15CST"},
				},
			},
		}),
	}

	statusBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"source_type": "ECS_FILE",
			"status":      "FAILED",
			"filter": []map[string]interface{}{
				{
					"key":      "VaultId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}"},
				},
				{
					"key":      "InstanceId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}"},
				},
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"source_type": "ECS_FILE",
			"status":      "PARTIAL_COMPLETE",
			"filter": []map[string]interface{}{
				{
					"key":      "VaultId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.vault_id}"},
				},
				{
					"key":      "InstanceId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.instance_id}_fake"},
				},
			},
		}),
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.HbrSupportRegions)
	}
	HbrBackupJobCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, planIdBackupConf, jobIdBackupConf, ecsBackupConf, ossBackupConf, nasBackupConf, statusBackupConf, completeTimeBackupConf, betweenTimeBackupConf)
}

func TestAccAliCloudHbrBackupJobsDataSource_ots(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_hbr_backup_jobs.default"
	name := fmt.Sprintf("tf-testAccHbrBackupJobTest%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceHbrBackupJobSourceConfig)
	otsBackupConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"source_type": "OTS",
			"filter": []map[string]interface{}{
				{
					"key":      "Status",
					"operator": "EQUAL",
					"values":   []string{"COMPLETE"},
				},
			},
		}),
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.HbrSupportRegions)
	}
	HbrBackupJobCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, otsBackupConf)
}

func dataSourceHbrBackupJobSourceConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	data "alicloud_hbr_ecs_backup_plans" "default" {
	}

	data "alicloud_hbr_oss_backup_plans" "default" {
	}

	data "alicloud_hbr_nas_backup_plans" "default" {

	}
`, name)
}

var existHbrBackupJobMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"jobs.#": acctest.CHECKSET,
		//"jobs.0.cross_account_type":    "SELF_ACCOUNT",
		//"jobs.0.cross_account_user_id": "0",
	}
}

var fakeHbrBackupJobMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"jobs.#": "0",
	}
}

var HbrBackupJobCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_hbr_backup_jobs.default",
	ExistMapFunc: existHbrBackupJobMapFunc,
	FakeMapFunc:  fakeHbrBackupJobMapFunc,
}
