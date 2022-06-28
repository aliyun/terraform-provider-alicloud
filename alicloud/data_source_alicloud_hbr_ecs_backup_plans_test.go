package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudHBREcsBackupPlansDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	ecsBackupIdsconf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrEcsBackupPlanSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_hbr_ecs_backup_plan.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudHbrEcsBackupPlanSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_hbr_ecs_backup_plan.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrEcsBackupPlanSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_hbr_ecs_backup_plan.default.ecs_backup_plan_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrEcsBackupPlanSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_hbr_ecs_backup_plan.default.ecs_backup_plan_name}_fake"`,
		}),
	}

	ecsInstanceIdconf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrEcsBackupPlanSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_hbr_ecs_backup_plan.default.id}"]`,
			"instance_id": `"${alicloud_hbr_ecs_backup_plan.default.instance_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrEcsBackupPlanSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_hbr_ecs_backup_plan.default.id}"]`,
			"instance_id": `"${alicloud_hbr_ecs_backup_plan.default.instance_id}_fake"`,
		}),
	}

	vaultIdconf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrEcsBackupPlanSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_hbr_ecs_backup_plan.default.id}"]`,
			"vault_id": `"${alicloud_hbr_ecs_backup_plan.default.vault_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrEcsBackupPlanSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_hbr_ecs_backup_plan.default.id}"]`,
			"vault_id": `"${alicloud_hbr_ecs_backup_plan.default.vault_id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrEcsBackupPlanSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_hbr_ecs_backup_plan.default.id}"]`,
			"name_regex":  `"${alicloud_hbr_ecs_backup_plan.default.ecs_backup_plan_name}"`,
			"instance_id": `"${alicloud_hbr_ecs_backup_plan.default.instance_id}"`,
			"vault_id":    `"${alicloud_hbr_ecs_backup_plan.default.vault_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrEcsBackupPlanSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_hbr_ecs_backup_plan.default.id}_fake"]`,
			"name_regex":  `"${alicloud_hbr_ecs_backup_plan.default.ecs_backup_plan_name}_fake"`,
			"instance_id": `"${alicloud_hbr_ecs_backup_plan.default.instance_id}_fake"`,
			"vault_id":    `"${alicloud_hbr_ecs_backup_plan.default.vault_id}_fake"`,
		}),
	}

	HbrEcsBackupPlanCheckInfo.dataSourceTestCheck(t, rand, ecsBackupIdsconf, nameRegexConf, ecsInstanceIdconf, vaultIdconf, allConf)
}

func testAccCheckAlicloudHbrEcsBackupPlanSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAcc%d"
}

resource "alicloud_hbr_vault" "default" {
  vault_name = "${var.name}"
}

data "alicloud_instances" "default" {
  name_regex = "no-deleteing-hbr-ecs-backup-plan"
  status     = "Running"
}

resource "alicloud_hbr_ecs_backup_plan" "default" {
  ecs_backup_plan_name = var.name
  instance_id =          "${data.alicloud_instances.default.instances.0.id}"
  vault_id =             "${alicloud_hbr_vault.default.id}"
  schedule =             "I|1602673264|PT2H"
  backup_type =          "COMPLETE"
  retention =            "2"
}

data "alicloud_hbr_ecs_backup_plans" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

var existHbrEcsBackupPlanMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"plans.#":                      "1",
		"plans.0.id":                   CHECKSET,
		"plans.0.ecs_backup_plan_name": fmt.Sprintf("tf-testAcc%d", rand),
		"plans.0.instance_id":          CHECKSET,
		"plans.0.schedule":             "I|1602673264|PT2H",
		"plans.0.backup_type":          "COMPLETE",
		"plans.0.vault_id":             CHECKSET,
	}
}

var fakeHbrEcsBackupPlanMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"plans.#": "0",
	}
}

var HbrEcsBackupPlanCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_hbr_ecs_backup_plans.default",
	existMapFunc: existHbrEcsBackupPlanMapFunc,
	fakeMapFunc:  fakeHbrEcsBackupPlanMapFunc,
}
