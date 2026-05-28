package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

const (
	sddpDataLimitsDataSourceTestRegion = connectivity.APSouthEast1 // Singapore
)

func TestAccAliCloudSddpDataLimitsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpDataLimitsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sddp_data_limit.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSddpDataLimitsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sddp_data_limit.default.id}_fake"]`,
		}),
	}
	resourceTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpDataLimitsDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_sddp_data_limit.default.id}"]`,
			"resource_type": `"RDS"`,
		}),
		fakeConfig: testAccCheckAlicloudSddpDataLimitsDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_sddp_data_limit.default.id}"]`,
			"resource_type": `"OSS"`,
		}),
	}
	parentIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpDataLimitsDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_sddp_data_limit.default.id}"]`,
			"parent_id": `"${alicloud_sddp_data_limit.default.parent_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSddpDataLimitsDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_sddp_data_limit.default.id}_fake"]`,
			"parent_id": `"${alicloud_sddp_data_limit.default.parent_id}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpDataLimitsDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_sddp_data_limit.default.id}"]`,
			"resource_type": `"RDS"`,
			"parent_id":     `"${alicloud_sddp_data_limit.default.parent_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudSddpDataLimitsDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_sddp_data_limit.default.id}_fake"]`,
			"parent_id":     `"${alicloud_sddp_data_limit.default.parent_id}_fake"`,
			"resource_type": `"OSS"`,
		}),
	}

	var alicloudSddpDataLimitsCheckInfo = dataSourceAttr{
		resourceId: "data.alicloud_sddp_data_limits.default",
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                  "1",
				"limits.#":               "1",
				"limits.0.audit_status":  CHECKSET,
				"limits.0.check_status":  CHECKSET,
				"limits.0.id":            CHECKSET,
				"limits.0.data_limit_id": CHECKSET,
				"limits.0.engine_type":   CHECKSET,
				"limits.0.local_name":    CHECKSET,
				"limits.0.log_store_day": CHECKSET,
				"limits.0.parent_id":     CHECKSET,
				"limits.0.port":          CHECKSET,
				"limits.0.resource_type": "RDS",
				"limits.0.user_name":     CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#": "0",
			}
		},
	}

	preCheck := func() {
		testAccPreCheck(t)
		checkoutSupportedRegions(t, true, []connectivity.Region{sddpDataLimitsDataSourceTestRegion})
	}
	alicloudSddpDataLimitsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, resourceTypeConf, parentIdConf, allConf)
}
func testAccCheckAlicloudSddpDataLimitsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testaccsddp-%d"
}

variable "region" {
  default = "%s"
}

variable "password" {
  default = "Test12345"
}

variable "database_name" {
  default = "tfaccdatabase"
}

data "alicloud_db_zones" "default" {
  engine = "MySQL"
  engine_version = "8.0"
  instance_charge_type = "PostPaid"
}

data "alicloud_db_instance_classes" "default" {
  zone_id = data.alicloud_db_zones.default.zones.1.id
  engine = "MySQL"
  engine_version = "8.0"
  category = "HighAvailability"
  instance_charge_type = "PostPaid"
  storage_type = "cloud_essd"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 2)
  zone_id      = data.alicloud_db_zones.default.zones.1.id
  vswitch_name = var.name
}

resource "alicloud_db_instance" "default" {
  engine           = "MySQL"
  engine_version   = "8.0"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
  instance_storage = "20"
  vswitch_id       = alicloud_vswitch.default.id
  instance_name    = var.name
}

resource "alicloud_rds_account" "default" {
  db_instance_id   = alicloud_db_instance.default.id
  account_name     = var.database_name
  account_password = var.password
}

resource "alicloud_db_database" "default" {
  instance_id = alicloud_db_instance.default.id
  name        = var.database_name
}

resource "alicloud_db_account_privilege" "default" {
  instance_id  = alicloud_db_instance.default.id
  account_name = alicloud_rds_account.default.name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.default.name]
}

resource "alicloud_sddp_data_limit" "default" {
  audit_status      = 0
  engine_type       = "MySQL"
  parent_id         = join(".", [alicloud_db_instance.default.id, var.database_name])
  resource_type     = "RDS"
  user_name         = var.database_name
  password          = var.password
  port              = 3306
  service_region_id = var.region
  depends_on        = [alicloud_db_account_privilege.default]
}

data "alicloud_sddp_data_limits" "default" {
	%s
}`, rand, sddpDataLimitsDataSourceTestRegion, strings.Join(pairs, " \n "))
	return config
}
