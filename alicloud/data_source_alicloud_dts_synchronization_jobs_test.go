package alicloud

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDTSSynchronizationJobsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	synchronizationJobidconf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDtsSynchronizationJobSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_dts_synchronization_job.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDtsSynchronizationJobSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_dts_synchronization_job.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDtsSynchronizationJobSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_dts_synchronization_job.default.dts_job_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDtsSynchronizationJobSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_dts_synchronization_job.default.dts_job_name}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDtsSynchronizationJobSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_dts_synchronization_job.default.id}"]`,
			"status": `"Synchronizing"`,
		}),
		fakeConfig: testAccCheckAlicloudDtsSynchronizationJobSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_dts_synchronization_job.default.id}"]`,
			"status": `"Failed"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDtsSynchronizationJobSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_dts_synchronization_job.default.id}"]`,
			"name_regex": `"${alicloud_dts_synchronization_job.default.dts_job_name}"`,
			"status":     `"Synchronizing"`,
		}),
		fakeConfig: testAccCheckAlicloudDtsSynchronizationJobSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_dts_synchronization_job.default.id}"]`,
			"name_regex": `"${alicloud_dts_synchronization_job.default.dts_job_name}"`,
			"status":     `"Failed"`,
		}),
	}

	DtsSynchronizationJobCheckInfo.dataSourceTestCheck(t, rand, synchronizationJobidconf, nameRegexConf, statusConf, allConf)
}

func testAccCheckAlicloudDtsSynchronizationJobSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccDtsSynchronizationJobs%d"
}
variable "region_id" {
	default = "%s"
}

data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
 name_regex = "^default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_db_instance" "source" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = data.alicloud_vswitches.default.ids.0
	instance_name = var.name
	tags = {
		"key1" = "value1"
		"key2" = "value2"
	}
}
resource "alicloud_db_instance" "dest" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
	vswitch_id = data.alicloud_vswitches.default.ids.0
	instance_name = var.name
	tags = {
		"key1" = "value1"
		"key2" = "value2"
	}
}

resource "alicloud_dts_synchronization_instance" "default" {
  payment_type                        = "PayAsYouGo"
  source_endpoint_engine_name         = "MySQL"
  source_endpoint_region              = var.region_id
  destination_endpoint_engine_name    = "MySQL"
  destination_endpoint_region         = var.region_id
  instance_class                      = "small"
  sync_architecture                   = "oneway"
}


resource "alicloud_db_database" "db" {
  count       = 2
  instance_id = alicloud_db_instance.dest.id
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
}

resource "alicloud_db_account" "account" {
  db_instance_id      = alicloud_db_instance.dest.id
  account_name        = "tftestdts"
  account_password    = "Test12345"
  account_description = "from terraform"
}

resource "alicloud_db_account_privilege" "privilege" {
  instance_id  = alicloud_db_account.account.instance_id
  account_name = alicloud_db_account.account.name
  privilege    = "ReadWrite"
  db_names     = alicloud_db_database.db.*.name
}

resource "alicloud_db_database" "db_r" {
  count       = 2
  instance_id = alicloud_db_instance.source.id
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
}

resource "alicloud_db_account" "account_r" {
  db_instance_id      = alicloud_db_instance.source.id
  account_name        = "tftestdts"
  account_password    = "Test12345"
  account_description = "from terraform"
}

resource "alicloud_db_account_privilege" "privilege_r" {
  instance_id  = alicloud_db_account.account_r.instance_id
  account_name = alicloud_db_account.account_r.name
  privilege    = "ReadWrite"
  db_names     = alicloud_db_database.db_r.*.name
}

resource "alicloud_dts_synchronization_job" "default" {
  dts_instance_id                     = alicloud_dts_synchronization_instance.default.id
  dts_job_name                        = "tf-testAccCase1"
  source_endpoint_instance_type       = "RDS"
  source_endpoint_instance_id         = alicloud_db_instance.source.id
  source_endpoint_engine_name         = "MySQL"
  source_endpoint_region              = var.region_id
  source_endpoint_database_name       = "tfaccountpri_0"
  source_endpoint_user_name           = "tftestdts"
  source_endpoint_password            = "Test12345"
  destination_endpoint_instance_type  = "RDS"
  destination_endpoint_instance_id    = alicloud_db_instance.dest.id
  destination_endpoint_engine_name    = "MySQL"
  destination_endpoint_region         = var.region_id
  destination_endpoint_database_name  = "tfaccountpri_0"
  destination_endpoint_user_name      = "tftestdts"
  destination_endpoint_password       = "Test12345"
  db_list                             = "{\"tfaccountpri_0\":{\"name\":\"tfaccountpri_0\",\"all\":true,\"state\":\"normal\"}}"
  structure_initialization            = "true"
  data_initialization                 = "true"
  data_synchronization                = "true"
}

data "alicloud_dts_synchronization_jobs" "default" {
%s
}
`, rand, os.Getenv("ALICLOUD_REGION"), strings.Join(pairs, "\n   "))
	return config
}

var existDtsSynchronizationJobMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"jobs.#":            "1",
		"jobs.0.dts_job_id": CHECKSET,
	}
}

var fakeDtsSynchronizationJobMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"jobs.#": "0",
	}
}

var DtsSynchronizationJobCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_dts_synchronization_jobs.default",
	existMapFunc: existDtsSynchronizationJobMapFunc,
	fakeMapFunc:  fakeDtsSynchronizationJobMapFunc,
}
