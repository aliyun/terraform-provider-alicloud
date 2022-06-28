package alicloud

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRdsDBInstancePerformanceDataSource(t *testing.T) {
	rand := acctest.RandInt()

	testConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRdsDBInstancePerformanceDataSourceName(rand, map[string]string{
			"db_instance_id": `"${alicloud_db_instance.default.id}"`,
			"start_time":     `"${var.startTime}"`,
			"end_time":       `"${var.endTime}"`,
			"key":            `"MySQL_ThreadStatus"`,
		}),
		fakeConfig: testAccCheckAlicloudRdsDBInstancePerformanceDataSourceName(rand, map[string]string{
			"db_instance_id": `"${alicloud_db_instance.default.id}"`,
			"start_time":     `"${var.startTime}"`,
			"end_time":       `"${var.errEndTime}"`,
			"key":            `"MySQL_ThreadStatus"`,
		}),
	}
	var existAlicloudRdsPerformanceDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"performance_keys.#":              "1",
			"performance_keys.0.key":          "MySQL_ThreadStatus",
			"performance_keys.0.unit":         CHECKSET,
			"performance_keys.0.value_format": CHECKSET,
			"performance_keys.0.values.#":     CHECKSET,
		}
	}
	var fakeAlicloudRdsPerformanceDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"performance_keys.#":              "1",
			"performance_keys.0.key":          "MySQL_ThreadStatus",
			"performance_keys.0.unit":         CHECKSET,
			"performance_keys.0.value_format": CHECKSET,
			"performance_keys.0.values.#":     "0",
		}
	}
	var alicloudRdsPerformanceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_rds_db_instance_performance.default",
		existMapFunc: existAlicloudRdsPerformanceDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudRdsPerformanceDataSourceNameMapFunc,
	}
	alicloudRdsPerformanceCheckInfo.dataSourceTestCheck(t, rand, testConf)
}

func testAccCheckAlicloudRdsDBInstancePerformanceDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	startTime := time.Now().AddDate(0, 0, -2).Format("2006-01-02T15:04Z")
	endTime := time.Now().AddDate(0, 0, 2).Format("2006-01-02T15:04Z")
	errEndTime := time.Now().AddDate(0, 0, -1).Format("2006-01-02T15:04Z")
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc-rds-db-instance-performance"
}

variable "startTime" {
 default = "%v"
}

variable "endTime" {
 default = "%v"
}

variable "errEndTime" {
 default = "%v"
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
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.ids.0
}

resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_db_zones.default.ids.0
  vswitch_name = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "cloud_essd"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
  vswitch_id               = local.vswitch_id
  instance_name            = var.name
}

data "alicloud_rds_db_instance_performance" "default" {
  %s
}`, startTime, endTime, errEndTime, strings.Join(pairs, "\n"))
}
