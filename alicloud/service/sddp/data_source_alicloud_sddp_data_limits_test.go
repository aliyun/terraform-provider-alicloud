package sddp_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

const (
	sddpDataLimitsDataSourceTestRegion = connectivity.APSouthEast1 // Singapore
)

func TestAccAliCloudSddpDataLimitsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100, 999)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSddpDataLimitsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sddp_data_limit.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudSddpDataLimitsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sddp_data_limit.default.id}_fake"]`,
		}),
	}
	resourceTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSddpDataLimitsDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_sddp_data_limit.default.id}"]`,
			"resource_type": `"RDS"`,
		}),
		FakeConfig: testAccCheckAlicloudSddpDataLimitsDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_sddp_data_limit.default.id}"]`,
			"resource_type": `"OSS"`,
		}),
	}
	parentIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSddpDataLimitsDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_sddp_data_limit.default.id}"]`,
			"parent_id": `"${alicloud_sddp_data_limit.default.parent_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudSddpDataLimitsDataSourceName(rand, map[string]string{
			"ids":       `["${alicloud_sddp_data_limit.default.id}_fake"]`,
			"parent_id": `"${alicloud_sddp_data_limit.default.parent_id}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSddpDataLimitsDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_sddp_data_limit.default.id}"]`,
			"resource_type": `"RDS"`,
			"parent_id":     `"${alicloud_sddp_data_limit.default.parent_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudSddpDataLimitsDataSourceName(rand, map[string]string{
			"ids":           `["${alicloud_sddp_data_limit.default.id}_fake"]`,
			"parent_id":     `"${alicloud_sddp_data_limit.default.parent_id}_fake"`,
			"resource_type": `"OSS"`,
		}),
	}

	var alicloudSddpDataLimitsCheckInfo = acctest.DataSourceAttr{
		ResourceId: "data.alicloud_sddp_data_limits.default",
		ExistMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                  "1",
				"limits.#":               "1",
				"limits.0.audit_status":  acctest.CHECKSET,
				"limits.0.check_status":  acctest.CHECKSET,
				"limits.0.id":            acctest.CHECKSET,
				"limits.0.data_limit_id": acctest.CHECKSET,
				"limits.0.engine_type":   acctest.CHECKSET,
				"limits.0.local_name":    acctest.CHECKSET,
				"limits.0.log_store_day": acctest.CHECKSET,
				"limits.0.parent_id":     acctest.CHECKSET,
				"limits.0.port":          acctest.CHECKSET,
				"limits.0.resource_type": "RDS",
				"limits.0.user_name":     acctest.CHECKSET,
			}
		},
		FakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#": "0",
			}
		},
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.CheckoutSupportedRegions(t, true, []connectivity.Region{sddpDataLimitsDataSourceTestRegion})
	}
	alicloudSddpDataLimitsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, resourceTypeConf, parentIdConf, allConf)
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

data "alicloud_db_instance_classes" "default" {
  zone_id = "ap-southeast-1c"
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
  zone_id      = "ap-southeast-1c"
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
