package mongodb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudMongodbBackupsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	existConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongodbBackupsDataSourceConfig(rand, map[string]string{
			"db_instance_id": `"${alicloud_mongodb_instance.defaultHrZmxC.id}"`,
			"ids":            `["${alicloud_mongodb_backup.default.id}"]`,
			"backup_id":      `"${alicloud_mongodb_backup.default.backup_id}"`,
			"enable_details": `true`,
			"output_file":    `"./test_output_file"`,
		}),
		FakeConfig: testAccCheckAliCloudMongodbBackupsDataSourceConfig(rand, map[string]string{
			"db_instance_id": `"${alicloud_mongodb_instance.defaultHrZmxC.id}"`,
			"start_time":     `"2020-01-01T00:00Z"`,
			"end_time":       `"2020-01-02T00:00Z"`,
		}),
	}
	var existAliCloudMongodbBackupsDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"backups.#":                   "1",
			"backups.0.id":                acctest.CHECKSET,
			"backups.0.backup_id":         acctest.CHECKSET,
			"backups.0.backup_job_id":     acctest.CHECKSET,
			"backups.0.backup_method":     acctest.CHECKSET,
			"backups.0.backup_mode":       acctest.CHECKSET,
			"backups.0.backup_size":       acctest.CHECKSET,
			"backups.0.backup_type":       acctest.CHECKSET,
			"backups.0.backup_start_time": acctest.CHECKSET,
			"backups.0.backup_end_time":   acctest.CHECKSET,
			"backups.0.status":            acctest.CHECKSET,
		}
	}
	var fakeAliCloudMongodbBackupsDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"backups.#": "0",
		}
	}
	var alicloudMongodbBackupsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_mongodb_backups.default",
		ExistMapFunc: existAliCloudMongodbBackupsDataSourceMapFunc,
		FakeMapFunc:  fakeAliCloudMongodbBackupsDataSourceMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
		acctest.PreCheck(t)
	}
	alicloudMongodbBackupsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, existConf)
}

func testAccCheckAliCloudMongodbBackupsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tfaccmongodb%d"
	}

	variable "zone_id" {
	  default = "cn-shanghai-b"
	}

	resource "alicloud_vpc" "defaultie35CW" {
	  cidr_block = "10.0.0.0/8"
	  vpc_name   = var.name
	}

	resource "alicloud_vswitch" "defaultg0DCAR" {
	  vpc_id     = alicloud_vpc.defaultie35CW.id
	  zone_id    = var.zone_id
	  cidr_block = "10.0.0.0/24"
	}

	resource "alicloud_mongodb_instance" "defaultHrZmxC" {
	  engine_version      = "4.4"
	  storage_type        = "cloud_essd1"
	  vswitch_id          = alicloud_vswitch.defaultg0DCAR.id
	  db_instance_storage = "20"
	  vpc_id              = alicloud_vpc.defaultie35CW.id
	  db_instance_class   = "mdb.shard.4x.large.d"
	  storage_engine      = "WiredTiger"
	  network_type        = "VPC"
	  zone_id             = var.zone_id
	}

	resource "alicloud_mongodb_backup" "default" {
	  backup_method           = "Snapshot"
	  db_instance_id          = alicloud_mongodb_instance.defaultHrZmxC.id
	  backup_retention_period = 7
	}

	data "alicloud_mongodb_backups" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
