package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudHBROssBackupPlansDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	ossBackupIdsconf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrOssBackupPlanSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_hbr_oss_backup_plan.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudHbrOssBackupPlanSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_hbr_oss_backup_plan.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrOssBackupPlanSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_hbr_oss_backup_plan.default.oss_backup_plan_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrOssBackupPlanSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_hbr_oss_backup_plan.default.oss_backup_plan_name}_fake"`,
		}),
	}

	ossBucketconf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrOssBackupPlanSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_hbr_oss_backup_plan.default.id}"]`,
			"bucket": `"${alicloud_oss_bucket.default.bucket}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrOssBackupPlanSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_hbr_oss_backup_plan.default.id}"]`,
			"bucket": `"${alicloud_oss_bucket.default.bucket}_fake"`,
		}),
	}

	vaultIdconf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrOssBackupPlanSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_hbr_oss_backup_plan.default.id}"]`,
			"vault_id": `"${alicloud_hbr_oss_backup_plan.default.vault_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrOssBackupPlanSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_hbr_oss_backup_plan.default.id}"]`,
			"vault_id": `"${alicloud_hbr_oss_backup_plan.default.vault_id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudHbrOssBackupPlanSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_hbr_oss_backup_plan.default.id}"]`,
			"name_regex": `"${alicloud_hbr_oss_backup_plan.default.oss_backup_plan_name}"`,
			"bucket":     `"${alicloud_oss_bucket.default.bucket}"`,
			"vault_id":   `"${alicloud_hbr_oss_backup_plan.default.vault_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudHbrOssBackupPlanSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_hbr_oss_backup_plan.default.id}_fake"]`,
			"name_regex": `"${alicloud_hbr_oss_backup_plan.default.oss_backup_plan_name}_fake"`,
			"bucket":     `"${alicloud_oss_bucket.default.bucket}_fake"`,
			"vault_id":   `"${alicloud_hbr_oss_backup_plan.default.vault_id}_fake"`,
		}),
	}

	HbrOssBackupPlanCheckInfo.dataSourceTestCheck(t, rand, ossBackupIdsconf, nameRegexConf, ossBucketconf, vaultIdconf, allConf)
}

func testAccCheckAlicloudHbrOssBackupPlanSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-test%d"
}
resource "alicloud_hbr_vault" "default" {
  vault_name = "${var.name}"
}
resource "alicloud_oss_bucket" "default" {
  bucket = "${var.name}"
}
resource "alicloud_hbr_oss_backup_plan" "default" {
  oss_backup_plan_name = var.name
  bucket   =             alicloud_oss_bucket.default.bucket
  vault_id =             alicloud_hbr_vault.default.id
  schedule =             "I|1602673264|PT2H"
  backup_type =          "COMPLETE"
  retention =            "2"
}
data "alicloud_hbr_oss_backup_plans" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

var existHbrOssBackupPlanMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"plans.#":                      "1",
		"plans.0.id":                   CHECKSET,
		"plans.0.oss_backup_plan_name": fmt.Sprintf("tf-test%d", rand),
		"plans.0.instance_id":          NOSET,
		"plans.0.schedule":             "I|1602673264|PT2H",
		"plans.0.backup_type":          "COMPLETE",
		"plans.0.vault_id":             CHECKSET,
		"plans.0.bucket":               CHECKSET,
	}
}

var fakeHbrOssBackupPlanMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"plans.#": "0",
	}
}

var HbrOssBackupPlanCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_hbr_oss_backup_plans.default",
	existMapFunc: existHbrOssBackupPlanMapFunc,
	fakeMapFunc:  fakeHbrOssBackupPlanMapFunc,
}
