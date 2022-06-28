package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSddpDataLimitsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	checkoutSupportedRegions(t, true, connectivity.SddpSupportRegions)
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

	var existAlicloudSddpDataLimitsDataSourceNameMapFunc = func(rand int) map[string]string {
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
	}
	var fakeAlicloudSddpDataLimitsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudSddpDataLimitsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_sddp_data_limits.default",
		existMapFunc: existAlicloudSddpDataLimitsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudSddpDataLimitsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
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

data "alicloud_db_zones" "default" {}

data "alicloud_db_instance_classes" "default" {
  engine         = "MySQL"
  engine_version = "5.6"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_db_zones.default.zones[0].id
}

resource "alicloud_db_instance" "default" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
  instance_storage = "10"
  vswitch_id       = data.alicloud_vswitches.default.ids[0]
  instance_name    = var.name
}

locals {
  parent_id = join(".", [alicloud_db_instance.default.id, var.database_name])
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
  parent_id         = local.parent_id
  resource_type     = "RDS"
  user_name         = var.database_name
  password          = var.password
  port              = 3306
  service_region_id = var.region
  depends_on        = [alicloud_db_account_privilege.default]
}

data "alicloud_sddp_data_limits" "default" {	
	%s
}
`, rand, defaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
