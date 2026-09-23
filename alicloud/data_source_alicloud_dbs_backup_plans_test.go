package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDbsBackupPlansDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDbsBackupPlansDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dbs_backup_plan.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDbsBackupPlansDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dbs_backup_plan.default.id}_fake"]`,
		}),
	}
	backupPlanNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDbsBackupPlansDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_dbs_backup_plan.default.id}"]`,
			"backup_plan_name": `"${alicloud_dbs_backup_plan.default.backup_plan_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDbsBackupPlansDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_dbs_backup_plan.default.id}"]`,
			"backup_plan_name": `"${alicloud_dbs_backup_plan.default.backup_plan_name}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDbsBackupPlansDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_dbs_backup_plan.default.backup_plan_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDbsBackupPlansDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_dbs_backup_plan.default.backup_plan_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDbsBackupPlansDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_dbs_backup_plan.default.id}"]`,
			"status": `"${alicloud_dbs_backup_plan.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudDbsBackupPlansDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_dbs_backup_plan.default.id}"]`,
			"status": `"stop"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDbsBackupPlansDataSourceName(rand, map[string]string{
			"backup_plan_name": `"${alicloud_dbs_backup_plan.default.backup_plan_name}"`,
			"ids":              `["${alicloud_dbs_backup_plan.default.id}"]`,
			"name_regex":       `"${alicloud_dbs_backup_plan.default.backup_plan_name}"`,
			"status":           `"${alicloud_dbs_backup_plan.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudDbsBackupPlansDataSourceName(rand, map[string]string{
			"backup_plan_name": `"${alicloud_dbs_backup_plan.default.backup_plan_name}_fake"`,
			"ids":              `["${alicloud_dbs_backup_plan.default.id}_fake"]`,
			"name_regex":       `"${alicloud_dbs_backup_plan.default.backup_plan_name}_fake"`,
			"status":           `"stop"`,
		}),
	}
	var existAlicloudDbsBackupPlansDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                        "1",
			"names.#":                                      "1",
			"plans.#":                                      "1",
			"plans.0.backup_plan_name":                     CHECKSET,
			"plans.0.payment_type":                         "PayAsYouGo",
			"plans.0.status":                               "running",
			"plans.0.backup_gateway_id":                    "",
			"plans.0.backup_method":                        "logical",
			"plans.0.backup_objects":                       "[{\"DBName\":\"tftestdatabase\"}]",
			"plans.0.resource_group_id":                    "",
			"plans.0.backup_period":                        "Monday",
			"plans.0.id":                                   CHECKSET,
			"plans.0.backup_plan_id":                       CHECKSET,
			"plans.0.backup_retention_period":              "740",
			"plans.0.backup_start_time":                    "14:22",
			"plans.0.backup_storage_type":                  "system",
			"plans.0.instance_class":                       "xlarge",
			"plans.0.database_type":                        "MySQL",
			"plans.0.source_endpoint_instance_type":        "RDS",
			"plans.0.source_endpoint_region":               "cn-hangzhou",
			"plans.0.source_endpoint_instance_id":          CHECKSET,
			"plans.0.source_endpoint_user_name":            "tftestnormal000",
			"plans.0.cross_aliyun_id":                      "",
			"plans.0.cross_role_name":                      "",
			"plans.0.duplication_archive_period":           CHECKSET,
			"plans.0.duplication_infrequent_access_period": CHECKSET,
			"plans.0.enable_backup_log":                    CHECKSET,
			"plans.0.oss_bucket_name":                      CHECKSET,
			"plans.0.source_endpoint_database_name":        "",
			"plans.0.source_endpoint_sid":                  "",
		}
	}
	var fakeAlicloudDbsBackupPlansDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudDbsBackupPlansCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_dbs_backup_plans.default",
		existMapFunc: existAlicloudDbsBackupPlansDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudDbsBackupPlansDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudDbsBackupPlansCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, backupPlanNameConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudDbsBackupPlansDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccBackupPlan-%d"
}
variable "database_region" {
  default = "%s"
}
variable "storage_region" {
  default = "%s"
}
variable "source_endpoint_region" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "local_ssd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "local_ssd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids.0
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_db_instance" "default" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "local_ssd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = local.vswitch_id
	instance_name = var.name
	security_group_ids = alicloud_security_group.default.*.id
}
resource "alicloud_db_database" "default" {
  instance_id = alicloud_db_instance.default.id
  name        = "tftestdatabase"
}
resource "alicloud_rds_account" "default" {
  db_instance_id = alicloud_db_instance.default.id
  account_name        = "tftestnormal000"
  account_password    = "Test12345"
}
resource "alicloud_db_account_privilege" "default" {
  instance_id  = alicloud_db_instance.default.id
  account_name = alicloud_rds_account.default.account_name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.default.name]
}
resource "alicloud_dbs_backup_plan" "default" {
	backup_plan_name = var.name
	payment_type =                  "PayAsYouGo"
	instance_class =                "xlarge"
	backup_method =                 "logical"
	database_type =                 "MySQL"
	database_region =               var.database_region
	storage_region =                var.storage_region
	instance_type =                 "RDS"
	source_endpoint_instance_type = "RDS"
	resource_group_id =             data.alicloud_resource_manager_resource_groups.default.ids.0
	source_endpoint_region =        var.source_endpoint_region
	source_endpoint_instance_id =   alicloud_db_instance.default.id
	source_endpoint_user_name =     alicloud_db_account_privilege.default.account_name
	source_endpoint_password =      alicloud_rds_account.default.account_password
	backup_objects =                "[{\"DBName\":\"${alicloud_db_database.default.name}\"}]"
	backup_period =                 "Monday"
	backup_start_time =             "14:22"
	backup_storage_type =           "system"
	backup_retention_period =       740
}

data "alicloud_dbs_backup_plans" "default" {	
	enable_details = true
	%s	
}
`, rand, defaultRegionToTest, defaultRegionToTest, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
