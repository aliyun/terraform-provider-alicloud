package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudHBROtsBackupPlansDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 99999)

	otsBackupIdsconf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrOtsBackupPlanSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_hbr_ots_backup_plan.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudHbrOtsBackupPlanSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_hbr_ots_backup_plan.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrOtsBackupPlanSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_hbr_ots_backup_plan.default.ots_backup_plan_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrOtsBackupPlanSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_hbr_ots_backup_plan.default.ots_backup_plan_name}_fake"`,
		}),
	}

	vaultIdconf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrOtsBackupPlanSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_hbr_ots_backup_plan.default.id}"]`,
			"vault_id": `"${alicloud_hbr_ots_backup_plan.default.vault_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrOtsBackupPlanSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_hbr_ots_backup_plan.default.id}"]`,
			"vault_id": `"${alicloud_hbr_ots_backup_plan.default.vault_id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrOtsBackupPlanSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_hbr_ots_backup_plan.default.id}"]`,
			"name_regex": `"${alicloud_hbr_ots_backup_plan.default.ots_backup_plan_name}"`,
			"vault_id":   `"${alicloud_hbr_ots_backup_plan.default.vault_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrOtsBackupPlanSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_hbr_ots_backup_plan.default.id}_fake"]`,
			"name_regex": `"${alicloud_hbr_ots_backup_plan.default.ots_backup_plan_name}_fake"`,
			"vault_id":   `"${alicloud_hbr_ots_backup_plan.default.vault_id}_fake"`,
		}),
	}

	HbrOtsBackupPlanCheckInfo.dataSourceTestCheck(t, rand, otsBackupIdsconf, nameRegexConf, vaultIdconf, allConf)
}

func testAccCheckAlicloudHbrOtsBackupPlanSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "testAcc%d"
}
resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
  vault_type = "OTS_BACKUP"
}

resource "alicloud_ots_instance" "foo" {
  name        = var.name
  description = var.name
  accessed_by = "Any"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

resource "alicloud_ots_table" "basic" {
  instance_name = alicloud_ots_instance.foo.name
  table_name    = var.name
  primary_key {
    name = "pk1"
    type = "Integer"
  }
  time_to_live                  = -1
  max_version                   = 1
  deviation_cell_version_in_sec = 1
}


resource "alicloud_hbr_ots_backup_plan" "default" {
  ots_backup_plan_name = var.name
  vault_id =             alicloud_hbr_vault.default.id
  backup_type =          "COMPLETE"
  schedule =             "I|1602673264|PT2H"
  retention =            "2"
  instance_name=      alicloud_ots_instance.foo.name
  ots_detail {
    table_names = [alicloud_ots_table.basic.table_name]
  }
}

data "alicloud_hbr_ots_backup_plans" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

var existHbrOtsBackupPlanMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"plans.#":                      "1",
		"plans.0.id":                   CHECKSET,
		"plans.0.ots_backup_plan_name": fmt.Sprintf("testAcc%d", rand),
		"plans.0.instance_id":          NOSET,
		"plans.0.schedule":             "I|1602673264|PT2H",
		"plans.0.backup_type":          "COMPLETE",
		"plans.0.vault_id":             CHECKSET,
		"plans.0.ots_detail.#":         "1",
	}
}

var fakeHbrOtsBackupPlanMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"plans.#": "0",
	}
}

var HbrOtsBackupPlanCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_hbr_ots_backup_plans.default",
	existMapFunc: existHbrOtsBackupPlanMapFunc,
	fakeMapFunc:  fakeHbrOtsBackupPlanMapFunc,
}
