package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudHBRBackupJobsDataSource(t *testing.T) {
	checkoutAccount(t, true)
	defer checkoutAccount(t, false)
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_hbr_backup_jobs.default"
	name := fmt.Sprintf("tf-testAccHbrBackupJobTest%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceHbrBackupJobSourceConfig)

	planIdBackupConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"source_type": "ECS_FILE",
			"filter": []map[string]interface{}{
				{
					"key":      "PlanId",
					"operator": "IN",
					"values":   []string{"${data.alicloud_hbr_ecs_backup_plans.default.plans.0.id}"},
				},
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
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

	jobIdBackupConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"source_type": "ECS_FILE",
			"filter": []map[string]interface{}{
				{
					"key":      "JobId",
					"operator": "IN",
					"values":   []string{"job-000fdy7y0b7g99dp7isg"},
				},
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"source_type": "ECS_FILE",
			"filter": []map[string]interface{}{
				{
					"key":      "PlanId",
					"operator": "IN",
					"values":   []string{"job-000fdy7y0b7g99dp7isg_fake"},
				},
			},
		}),
	}

	ecsBackupConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
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
		fakeConfig: testAccConfig(map[string]interface{}{
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

	ossBackupConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
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
		fakeConfig: testAccConfig(map[string]interface{}{
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

	nasBackupConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
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
		fakeConfig: testAccConfig(map[string]interface{}{
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

	completeTimeBackupConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
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
					"operator": "GREATER_THAN_OR_EQUAL",
					"values":   []string{"2021-08-23T14:17:15CST"},
				},
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
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

	betweenTimeBackupConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
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
					"values":   []string{"2021-08-23T14:17:15CST", "2022-08-23T14:17:15CST"},
				},
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
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

	statusBackupConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"source_type": "ECS_FILE",
			"status":      "COMPLETE",
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
		fakeConfig: testAccConfig(map[string]interface{}{
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
		testAccPreCheckWithRegions(t, true, connectivity.HbrSupportRegions)
	}
	HbrBackupJobCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, planIdBackupConf, jobIdBackupConf, ecsBackupConf, ossBackupConf, nasBackupConf, statusBackupConf, completeTimeBackupConf, betweenTimeBackupConf)
}

func TestAccAlicloudHBRBackupJobsDataSource_ots(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_hbr_backup_jobs.default"
	name := fmt.Sprintf("tf-testAccHbrBackupJobTest%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceHbrBackupJobSourceConfig)
	otsBackupConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
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
		testAccPreCheckWithRegions(t, true, connectivity.HbrSupportRegions)
	}
	HbrBackupJobCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, otsBackupConf)
}

func dataSourceHbrBackupJobSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_hbr_ecs_backup_plans" "default" {
    name_regex = "plan-tf-used-dont-delete"
}

data "alicloud_hbr_oss_backup_plans" "default" {
	name_regex = "plan-tf-used-dont-delete"
}

data "alicloud_hbr_nas_backup_plans" "default" {
	name_regex = "plan-tf-used-dont-delete"
}`, name)
}

var existHbrBackupJobMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"jobs.#":    CHECKSET,
		"jobs.0.id": CHECKSET,
	}
}

var fakeHbrBackupJobMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"jobs.#": "0",
	}
}

var HbrBackupJobCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_hbr_backup_jobs.default",
	existMapFunc: existHbrBackupJobMapFunc,
	fakeMapFunc:  fakeHbrBackupJobMapFunc,
}
