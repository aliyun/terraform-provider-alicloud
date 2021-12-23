package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRdsBackupsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRdsBackupsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_backup.default.id}"]`,
			"db_instance_id": `"${alicloud_rds_backup.default.db_instance_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudRdsBackupsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_backup.default.id}_fake"]`,
			"db_instance_id": `"${alicloud_rds_backup.default.db_instance_id}"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRdsBackupsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_backup.default.id}"]`,
			"db_instance_id": `"${alicloud_rds_backup.default.db_instance_id}"`,
			"backup_status":  `"Success"`,
		}),
		fakeConfig: testAccCheckAlicloudRdsBackupsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_backup.default.id}"]`,
			"db_instance_id": `"${alicloud_rds_backup.default.db_instance_id}"`,
			"backup_status":  `"Failed"`,
		}),
	}
	modeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRdsBackupsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_backup.default.id}"]`,
			"db_instance_id": `"${alicloud_rds_backup.default.db_instance_id}"`,
			"backup_mode":    `"Manual"`,
		}),
		fakeConfig: testAccCheckAlicloudRdsBackupsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_backup.default.id}"]`,
			"db_instance_id": `"${alicloud_rds_backup.default.db_instance_id}"`,
			"backup_mode":    `"Automated"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRdsBackupsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_backup.default.id}"]`,
			"db_instance_id": `"${alicloud_rds_backup.default.db_instance_id}"`,
			"backup_status":  `"Success"`,
			"backup_mode":    `"Manual"`,
		}),
		fakeConfig: testAccCheckAlicloudRdsBackupsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_rds_backup.default.id}"]`,
			"db_instance_id": `"${alicloud_rds_backup.default.db_instance_id}"`,
			"backup_status":  `"Failed"`,
			"backup_mode":    `"Automated"`,
		}),
	}
	var existAlicloudRdsBackupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"backups.#":                              "1",
			"backups.0.backup_download_url":          ``,
			"backups.0.backup_end_time":              CHECKSET,
			"backups.0.id":                           CHECKSET,
			"backups.0.backup_id":                    CHECKSET,
			"backups.0.backup_initiator":             "User",
			"backups.0.backup_intranet_download_url": ``,
			"backups.0.backup_method":                CHECKSET,
			"backups.0.backup_mode":                  "Manual",
			"backups.0.backup_size":                  CHECKSET,
			"backups.0.backup_start_time":            CHECKSET,
			"backups.0.backup_type":                  "FullBackup",
			"backups.0.consistent_time":              CHECKSET,
			"backups.0.copy_only_backup":             ``,
			"backups.0.db_instance_id":               CHECKSET,
			"backups.0.encryption":                   "{}",
			"backups.0.host_instance_id":             CHECKSET,
			"backups.0.is_avail":                     "1",
			"backups.0.meta_status":                  ``,
			"backups.0.backup_status":                "Success",
			"backups.0.storage_class":                "0",
			"backups.0.store_status":                 "Disabled",
		}
	}
	var fakeAlicloudRdsBackupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"backups.#": "0",
		}
	}
	var alicloudRdsBackupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_rds_backups.default",
		existMapFunc: existAlicloudRdsBackupsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudRdsBackupsDataSourceNameMapFunc,
	}
	alicloudRdsBackupsCheckInfo.dataSourceTestCheck(t, rand, idsConf, statusConf, modeConf, allConf)
}
func testAccCheckAlicloudRdsBackupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-rds-backup"
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

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_db_zones.default.ids.0
  cidr_block   = "172.16.0.0/24"
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "cloud_essd"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
  vswitch_id               = alicloud_vswitch.default.id
  instance_name            = var.name
}

resource "alicloud_rds_backup" "default" {
  db_instance_id    = alicloud_db_instance.default.id
  remove_from_state = "true"
}

data "alicloud_rds_backups" "default" {	
  %s
}`, strings.Join(pairs, "\n"))
}
