package dts_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const dtsMigrationJobsDataSourceId = "data.alicloud_dts_migration_jobs.default"

func TestAccAlicloudDTSMigrationJobsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			// idConf
			{
				Config: testAccCheckAlicloudDtsMigrationJobsDataSourceName(rand, map[string]string{
					"ids": `["${alicloud_dts_migration_job.default.id}"]`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsMigrationJobsExistChecks()...),
			},
			{
				Config: testAccCheckAlicloudDtsMigrationJobsDataSourceName(rand, map[string]string{
					"ids": `["${alicloud_dts_migration_job.default.id}_fake"]`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsMigrationJobsFakeChecks()...),
			},
			// nameRegexConf
			{
				Config: testAccCheckAlicloudDtsMigrationJobsDataSourceName(rand, map[string]string{
					"name_regex": `"${alicloud_dts_migration_job.default.dts_job_name}"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsMigrationJobsExistChecks()...),
			},
			{
				Config: testAccCheckAlicloudDtsMigrationJobsDataSourceName(rand, map[string]string{
					"name_regex": `"${alicloud_dts_migration_job.default.dts_job_name}_fake"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsMigrationJobsFakeChecks()...),
			},
			// allConf
			{
				Config: testAccCheckAlicloudDtsMigrationJobsDataSourceName(rand, map[string]string{
					"ids":        `["${alicloud_dts_migration_job.default.id}"]`,
					"name_regex": `"${alicloud_dts_migration_job.default.dts_job_name}"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsMigrationJobsExistChecks()...),
			},
			{
				Config: testAccCheckAlicloudDtsMigrationJobsDataSourceName(rand, map[string]string{
					"ids":        `["${alicloud_dts_migration_job.default.id}_fake"]`,
					"name_regex": `"${alicloud_dts_migration_job.default.dts_job_name}_fake"`,
				}),
				Check: resource.ComposeTestCheckFunc(dtsMigrationJobsFakeChecks()...),
			},
		},
	})
}

func dtsMigrationJobsExistChecks() []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(dtsMigrationJobsDataSourceId, "ids.#", "1"),
		resource.TestCheckResourceAttr(dtsMigrationJobsDataSourceId, "names.#", "1"),
		resource.TestCheckResourceAttr(dtsMigrationJobsDataSourceId, "jobs.#", "1"),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.data_initialization"),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.data_synchronization"),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.db_list"),
		resource.TestCheckResourceAttr(dtsMigrationJobsDataSourceId, "jobs.0.destination_endpoint_data_base_name", ""),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.destination_endpoint_engine_name"),
		resource.TestCheckResourceAttr(dtsMigrationJobsDataSourceId, "jobs.0.destination_endpoint_ip", ""),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.destination_endpoint_instance_id"),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.destination_endpoint_instance_type"),
		resource.TestCheckResourceAttr(dtsMigrationJobsDataSourceId, "jobs.0.destination_endpoint_oracle_sid", ""),
		resource.TestCheckResourceAttr(dtsMigrationJobsDataSourceId, "jobs.0.destination_endpoint_port", ""),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.destination_endpoint_region"),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.destination_endpoint_user_name"),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.dts_instance_id"),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.id"),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.dts_job_id"),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.dts_job_name"),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.payment_type"),
		resource.TestCheckResourceAttr(dtsMigrationJobsDataSourceId, "jobs.0.source_endpoint_database_name", ""),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.source_endpoint_engine_name"),
		resource.TestCheckResourceAttr(dtsMigrationJobsDataSourceId, "jobs.0.source_endpoint_ip", ""),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.source_endpoint_instance_id"),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.source_endpoint_instance_type"),
		resource.TestCheckResourceAttr(dtsMigrationJobsDataSourceId, "jobs.0.source_endpoint_oracle_sid", ""),
		resource.TestCheckResourceAttr(dtsMigrationJobsDataSourceId, "jobs.0.source_endpoint_owner_id", ""),
		resource.TestCheckResourceAttr(dtsMigrationJobsDataSourceId, "jobs.0.source_endpoint_port", ""),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.source_endpoint_region"),
		resource.TestCheckResourceAttr(dtsMigrationJobsDataSourceId, "jobs.0.source_endpoint_role", ""),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.source_endpoint_user_name"),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.status"),
		resource.TestCheckResourceAttrSet(dtsMigrationJobsDataSourceId, "jobs.0.structure_initialization"),
	}
}

func dtsMigrationJobsFakeChecks() []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(dtsMigrationJobsDataSourceId, "ids.#", "0"),
		resource.TestCheckResourceAttr(dtsMigrationJobsDataSourceId, "jobs.#", "0"),
		resource.TestCheckResourceAttr(dtsMigrationJobsDataSourceId, "names.#", "0"),
	}
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
