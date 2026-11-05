package hbr_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudHbrHanaBackupPlansDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.HBRSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrHanaBackupPlansDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_hbr_hana_backup_plan.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudHbrHanaBackupPlansDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_hbr_hana_backup_plan.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrHanaBackupPlansDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_hbr_hana_backup_plan.default.plan_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudHbrHanaBackupPlansDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_hbr_hana_backup_plan.default.plan_name}_fake"`,
		}),
	}
	vaultIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrHanaBackupPlansDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_hbr_hana_backup_plan.default.id}"]`,
			"vault_id": `"${alicloud_hbr_hana_backup_plan.default.vault_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudHbrHanaBackupPlansDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_hbr_hana_instance.default.id}"]`,
			"vault_id": `"${alicloud_hbr_hana_instance.default.vault_id}_fake"`,
		}),
	}
	databaseNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrHanaBackupPlansDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_hbr_hana_backup_plan.default.id}"]`,
			"database_name": `"${alicloud_hbr_hana_backup_plan.default.database_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudHbrHanaBackupPlansDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_hbr_hana_backup_plan.default.id}"]`,
			"database_name": `"${alicloud_hbr_hana_backup_plan.default.database_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudHbrHanaBackupPlansDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_hbr_hana_backup_plan.default.id}"]`,
			"name_regex":    `"${alicloud_hbr_hana_backup_plan.default.plan_name}"`,
			"vault_id":      `"${alicloud_hbr_hana_backup_plan.default.vault_id}"`,
			"database_name": `"${alicloud_hbr_hana_backup_plan.default.database_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudHbrHanaBackupPlansDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_hbr_hana_backup_plan.default.id}_fake"]`,
			"name_regex":    `"${alicloud_hbr_hana_backup_plan.default.plan_name}_fake"`,
			"vault_id":      `"${alicloud_hbr_hana_instance.default.vault_id}_fake"`,
			"database_name": `"${alicloud_hbr_hana_backup_plan.default.database_name}_fake"`,
		}),
	}
	var existAlicloudHbrHanaBackupPlansDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":               "1",
			"ids.#":                 "1",
			"plans.#":               "1",
			"plans.0.backup_prefix": "DIFF_DATA_BACKUP",
			"plans.0.backup_type":   "COMPLETE",
			"plans.0.cluster_id":    acctest.CHECKSET,
			"plans.0.database_name": "SYSTEMDB",
			"plans.0.plan_id":       acctest.CHECKSET,
			"plans.0.id":            acctest.CHECKSET,
			"plans.0.plan_name":     fmt.Sprintf("tf-testacchanabackupplan-%d", rand),
			"plans.0.schedule":      "I|1602673264|P1D",
			"plans.0.vault_id":      acctest.CHECKSET,
			"plans.0.status":        acctest.CHECKSET,
		}
	}
	var fakeAlicloudHbrHanaBackupPlansDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudHbrHanaBackupPlansCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_hbr_hana_backup_plans.default",
		ExistMapFunc: existAlicloudHbrHanaBackupPlansDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudHbrHanaBackupPlansDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudHbrHanaBackupPlansCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, vaultIdConf, databaseNameConf, allConf)
}
func testAccCheckAlicloudHbrHanaBackupPlansDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testacchanabackupplan-%d"
}

resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
}
data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_hbr_hana_instance" "default" {
	alert_setting = "INHERITED"
	hana_name = var.name
	host = "1.1.1.1"
	instance_number = "1"
	password = "YouPassword123"
	resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
	sid = "HXE"
	use_ssl = false
	user_name = "admin"
	validate_certificate = false
	vault_id = alicloud_hbr_vault.default.id
}

resource "alicloud_hbr_hana_backup_plan" "default" {
	backup_prefix = "DIFF_DATA_BACKUP"
	backup_type = "COMPLETE"
	cluster_id = alicloud_hbr_hana_instance.default.hana_instance_id
	database_name = "SYSTEMDB"
	plan_name = var.name
	resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
	schedule = "I|1602673264|P1D"
	vault_id = alicloud_hbr_hana_instance.default.vault_id
}

data "alicloud_hbr_hana_backup_plans" "default" {
	cluster_id = alicloud_hbr_hana_instance.default.hana_instance_id
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
