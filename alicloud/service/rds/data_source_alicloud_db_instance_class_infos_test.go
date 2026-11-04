package rds_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudRdsDBInstanceClassInfosDatasource_main(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_db_instance_class_infos.default"
	name := "tf-testAccAlicloudRdsDBInstanceClassInfos"
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlicloudDBInstanceClassInfosDataSourceConfig)

	BardsInstanceClassesConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"commodity_code": "bards",
			"order_type":     "BUY",
		}),
	}
	RdsInstanceClassesConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"commodity_code": "rds",
			"order_type":     "CONVERT",
		}),
	}

	var existDBInstanceMainMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                acctest.CHECKSET,
			"infos.#":              acctest.CHECKSET,
			"infos.0.class_code":   acctest.CHECKSET,
			"infos.0.class_group":  acctest.CHECKSET,
			"infos.0.cpu":          acctest.CHECKSET,
			"infos.0.memory_class": acctest.CHECKSET,
		}
	}

	var fakeDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"infos.#": "0",
			"ids.#":   "0",
		}
	}

	var DBInstanceMainCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_db_instance_class_infos.default",
		ExistMapFunc: existDBInstanceMainMapFunc,
		FakeMapFunc:  fakeDBInstanceMapFunc,
	}

	DBInstanceMainCheckInfo.DataSourceTestCheck(t, rand, BardsInstanceClassesConf, RdsInstanceClassesConf)
}

func TestAccAlicloudRdsDBInstanceClassInfosDatasource_readOnly(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_db_instance_class_infos.default"
	name := "tf-testAccAlicloudRdsDBInstanceClassInfos"
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlicloudDBInstanceClassInfosDataSourceConfig)

	RordsInstanceClassesConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"commodity_code": "rords",
			"order_type":     "UPGRADE",
			"db_instance_id": "${alicloud_db_instance.default.id}",
		}),
	}
	RdsRordspreInstanceClassesConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"commodity_code": "rds_rordspre_public_cn",
			"order_type":     "RENEW",
			"db_instance_id": "${alicloud_db_instance.default.id}",
		}),
		ExistChangMap: map[string]string{
			"infos.0.instruction_set_arch": acctest.REMOVEKEY,
		},
	}

	var existDBInstanceReadOnlyMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        acctest.CHECKSET,
			"infos.#":                      acctest.CHECKSET,
			"infos.0.class_code":           acctest.CHECKSET,
			"infos.0.class_group":          acctest.CHECKSET,
			"infos.0.cpu":                  acctest.CHECKSET,
			"infos.0.max_connections":      acctest.CHECKSET,
			"infos.0.max_iombps":           acctest.CHECKSET,
			"infos.0.max_iops":             acctest.CHECKSET,
			"infos.0.memory_class":         acctest.CHECKSET,
			"infos.0.instruction_set_arch": acctest.CHECKSET,
		}
	}

	var fakeDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"infos.#": "0",
			"ids.#":   "0",
		}
	}

	var DBInstanceReadOnlyCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_db_instance_class_infos.default",
		ExistMapFunc: existDBInstanceReadOnlyMapFunc,
		FakeMapFunc:  fakeDBInstanceMapFunc,
	}

	DBInstanceReadOnlyCheckInfo.DataSourceTestCheck(t, rand, RordsInstanceClassesConf, RdsRordspreInstanceClassesConf)
}

func testAccCheckAlicloudDBInstanceClassInfosDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine                   = "SQLServer"
	engine_version           = "2017_ent"
	instance_charge_type     = "PostPaid"
 	db_instance_storage_type = "cloud_essd"
	category                 = "AlwaysOn"
}


data "alicloud_vswitches" "default" {
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
	name   = var.name
	vpc_id = data.alicloud_vswitches.default.vswitches.0.vpc_id
}

resource "alicloud_db_instance" "default" {
    engine                   = "SQLServer"
	engine_version           = "2017_ent"
 	db_instance_storage_type = "cloud_essd"
	category                 = "AlwaysOn"
	instance_type            = "mssql.x4.medium.e2"
	instance_storage         = 20
	vswitch_id               = data.alicloud_vswitches.default.vswitches.0.vswitch_id
	instance_name            = var.name
}
`, name)
}
