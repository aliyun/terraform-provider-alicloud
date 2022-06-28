package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudHBRNasBackupPlansDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	nasBackupIdsconf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrNasBackupPlanSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_hbr_nas_backup_plan.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudHbrNasBackupPlanSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_hbr_nas_backup_plan.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrNasBackupPlanSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_hbr_nas_backup_plan.default.nas_backup_plan_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrNasBackupPlanSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_hbr_nas_backup_plan.default.nas_backup_plan_name}_fake"`,
		}),
	}

	fileSystemIdconf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrNasBackupPlanSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_hbr_nas_backup_plan.default.id}"]`,
			"file_system_id": `"${alicloud_hbr_nas_backup_plan.default.file_system_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrNasBackupPlanSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_hbr_nas_backup_plan.default.id}"]`,
			"file_system_id": `"${alicloud_hbr_nas_backup_plan.default.file_system_id}_fake"`,
		}),
	}

	vaultIdconf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrNasBackupPlanSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_hbr_nas_backup_plan.default.id}"]`,
			"vault_id": `"${alicloud_hbr_nas_backup_plan.default.vault_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrNasBackupPlanSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_hbr_nas_backup_plan.default.id}"]`,
			"vault_id": `"${alicloud_hbr_nas_backup_plan.default.vault_id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrNasBackupPlanSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_hbr_nas_backup_plan.default.id}"]`,
			"name_regex":     `"${alicloud_hbr_nas_backup_plan.default.nas_backup_plan_name}"`,
			"file_system_id": `"${alicloud_hbr_nas_backup_plan.default.file_system_id}"`,
			"vault_id":       `"${alicloud_hbr_nas_backup_plan.default.vault_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrNasBackupPlanSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_hbr_nas_backup_plan.default.id}_fake"]`,
			"name_regex":     `"${alicloud_hbr_nas_backup_plan.default.nas_backup_plan_name}_fake"`,
			"file_system_id": `"${alicloud_hbr_nas_backup_plan.default.file_system_id}_fake"`,
			"vault_id":       `"${alicloud_hbr_nas_backup_plan.default.vault_id}_fake"`,
		}),
	}

	HbrNasBackupPlanCheckInfo.dataSourceTestCheck(t, rand, nasBackupIdsconf, nameRegexConf, fileSystemIdconf, vaultIdconf, allConf)
}

func testAccCheckAlicloudHbrNasBackupPlanSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAcc%d"
}

resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
}

resource "alicloud_nas_file_system" "default" {
  protocol_type = "NFS"
  storage_type  = "Performance"
  description   = var.name
  encrypt_type  = "1"
}

data "alicloud_nas_file_systems" "default" {
  protocol_type       = "NFS"
  description_regex   = alicloud_nas_file_system.default.description
}

resource "alicloud_hbr_nas_backup_plan" "default" {
  depends_on =           ["alicloud_nas_file_system.default"]
  nas_backup_plan_name = var.name
  file_system_id =      "${data.alicloud_nas_file_systems.default.systems.0.id}"
  schedule =            "I|1602673264|PT2H"
  backup_type =         "COMPLETE"
  vault_id =            "${alicloud_hbr_vault.default.id}"
  create_time =         "${data.alicloud_nas_file_systems.default.systems.0.create_time}"
  retention =			"2"
  path =                ["/"]
}

data "alicloud_hbr_nas_backup_plans" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

var existHbrNasBackupPlanMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"plans.#":                      "1",
		"plans.0.id":                   CHECKSET,
		"plans.0.nas_backup_plan_name": fmt.Sprintf("tf-testAcc%d", rand),
		"plans.0.file_system_id":       CHECKSET,
		"plans.0.schedule":             "I|1602673264|PT2H",
		"plans.0.backup_type":          "COMPLETE",
		"plans.0.vault_id":             CHECKSET,
		"plans.0.create_time":          CHECKSET,
	}
}

var fakeHbrNasBackupPlanMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"plans.#": "0",
	}
}

var HbrNasBackupPlanCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_hbr_nas_backup_plans.default",
	existMapFunc: existHbrNasBackupPlanMapFunc,
	fakeMapFunc:  fakeHbrNasBackupPlanMapFunc,
}
