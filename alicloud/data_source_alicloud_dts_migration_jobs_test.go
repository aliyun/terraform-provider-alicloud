package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlicloudDTSMigrationJobsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDtsMigrationJobsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dts_migration_job.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDtsMigrationJobsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dts_migration_job.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDtsMigrationJobsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_dts_migration_job.default.dts_job_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDtsMigrationJobsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_dts_migration_job.default.dts_job_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDtsMigrationJobsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_dts_migration_job.default.id}"]`,
			"name_regex": `"${alicloud_dts_migration_job.default.dts_job_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDtsMigrationJobsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_dts_migration_job.default.id}_fake"]`,
			"name_regex": `"${alicloud_dts_migration_job.default.dts_job_name}_fake"`,
		}),
	}

	var existAlicloudDtsMigrationJobsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"jobs.#":                      "1",
			"jobs.0.data_initialization":  CHECKSET,
			"jobs.0.data_synchronization": CHECKSET,
			"jobs.0.db_list":              CHECKSET,
			"jobs.0.destination_endpoint_data_base_name": "",
			"jobs.0.destination_endpoint_engine_name":    CHECKSET,
			"jobs.0.destination_endpoint_ip":             "",
			"jobs.0.destination_endpoint_instance_id":    CHECKSET,
			"jobs.0.destination_endpoint_instance_type":  CHECKSET,
			"jobs.0.destination_endpoint_oracle_sid":     "",
			"jobs.0.destination_endpoint_port":           "",
			"jobs.0.destination_endpoint_region":         CHECKSET,
			"jobs.0.destination_endpoint_user_name":      CHECKSET,
			"jobs.0.dts_instance_id":                     CHECKSET,
			"jobs.0.id":                                  CHECKSET,
			"jobs.0.dts_job_id":                          CHECKSET,
			"jobs.0.dts_job_name":                        CHECKSET,
			"jobs.0.payment_type":                        CHECKSET,
			"jobs.0.source_endpoint_database_name":       "",
			"jobs.0.source_endpoint_engine_name":         CHECKSET,
			"jobs.0.source_endpoint_ip":                  "",
			"jobs.0.source_endpoint_instance_id":         CHECKSET,
			"jobs.0.source_endpoint_instance_type":       CHECKSET,
			"jobs.0.source_endpoint_oracle_sid":          "",
			"jobs.0.source_endpoint_owner_id":            "",
			"jobs.0.source_endpoint_port":                "",
			"jobs.0.source_endpoint_region":              CHECKSET,
			"jobs.0.source_endpoint_role":                "",
			"jobs.0.source_endpoint_user_name":           CHECKSET,
			"jobs.0.status":                              CHECKSET,
			"jobs.0.structure_initialization":            CHECKSET,
		}
	}
	var fakeAlicloudDtsMigrationJobsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"jobs.#":  "0",
			"names.#": "0",
		}
	}
	var alicloudDtsMigrationJobsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_dts_migration_jobs.default",
		existMapFunc: existAlicloudDtsMigrationJobsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudDtsMigrationJobsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	}
	alicloudDtsMigrationJobsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudDtsMigrationJobsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {	
	default = "tf-testAccMigrationJob-%d"
}

variable "password" {
  default = "Test12345"
}

variable "database_name" {
  default = "tftestdatabase"
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_db_instance" "default" {
  count            = 2
  engine           = "MySQL"
  engine_version   = "8.0"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
  vswitch_id       = data.alicloud_vswitches.default.ids.0
  instance_name    = join("", [var.name, count.index])
}

resource "alicloud_rds_account" "default" {
  count            = 2
  db_instance_id   = alicloud_db_instance.default[count.index].id
  account_name     = join("", [var.database_name, count.index])
  account_password = var.password
}

resource "alicloud_db_database" "default" {
  count       = 2
  instance_id = alicloud_db_instance.default[count.index].id
  name        = var.database_name
}

resource "alicloud_db_account_privilege" "default" {
  count        = 2
  instance_id  = alicloud_db_instance.default[count.index].id
  account_name = alicloud_rds_account.default[count.index].name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.default[count.index].name]
}

resource "alicloud_dts_migration_instance" "default" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "MySQL"
  source_endpoint_region           = data.alicloud_regions.default.regions.0.id
  destination_endpoint_engine_name = "MySQL"
  destination_endpoint_region      = data.alicloud_regions.default.regions.0.id
  instance_class                   = "small"
  sync_architecture                = "oneway"
}

resource "alicloud_dts_migration_job" "default" {
  dts_instance_id                    = alicloud_dts_migration_instance.default.id
  dts_job_name                       = var.name
  source_endpoint_instance_type      = "RDS"
  source_endpoint_instance_id        = alicloud_db_instance.default.0.id
  source_endpoint_engine_name        = "MySQL"
  source_endpoint_region             = data.alicloud_regions.default.regions.0.id
  source_endpoint_user_name          = alicloud_rds_account.default.0.name
  source_endpoint_password           = var.password
  destination_endpoint_instance_type = "RDS"
  destination_endpoint_instance_id   = alicloud_db_instance.default.1.id
  destination_endpoint_engine_name   = "MySQL"
  destination_endpoint_region        = data.alicloud_regions.default.regions.0.id
  destination_endpoint_user_name     = alicloud_rds_account.default.1.name
  destination_endpoint_password      = var.password
  db_list                            = "{\"tftestdatabase\":{\"name\":\"tftestdatabase\",\"all\":true}}"
  structure_initialization           = true
  data_initialization                = true
  data_synchronization               = true
  status                             = "Migrating"
  depends_on                         = [alicloud_db_account_privilege.default]
}

data "alicloud_dts_migration_jobs" "default" {	
	enable_details = true
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
