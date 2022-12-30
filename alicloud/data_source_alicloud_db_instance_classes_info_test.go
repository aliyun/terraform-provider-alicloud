package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"testing"
)

func TestAccAlicloudRdsDBInstanceClassesInfoDatasource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_db_instance_classes_info.default"
	name := "tf-testAccAlicloudRdsDBInstanceClassesInfo"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlicloudDBInstanceClassesInfoDataSourceConfig)

	BardsInstanceClassesConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"commodity_code": "bards",
			"order_type":     "BUY",
		}),
	}
	RdsInstanceClassesConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"commodity_code": "rds",
			"order_type":     "CONVERT",
		}),
	}
	RordsInstanceClassesConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"commodity_code": "rords",
			"order_type":     "UPGRADE",
			"db_instance_id": "${alicloud_db_instance.default.id}",
		}),
	}
	RdsRordspreInstanceClassesConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"commodity_code": "rds_rordspre_public_cn",
			"order_type":     "RENEW",
			"db_instance_id": "${alicloud_db_instance.default.id}",
		}),
	}

	var existDBInstanceMainMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                    CHECKSET,
			"instance_classes_infos.#":                 CHECKSET,
			"instance_classes_infos.0.class_code":      CHECKSET,
			"instance_classes_infos.0.class_group":     CHECKSET,
			"instance_classes_infos.0.cpu":             CHECKSET,
			"instance_classes_infos.0.max_connections": CHECKSET,
			"instance_classes_infos.0.max_iombps":      CHECKSET,
			"instance_classes_infos.0.max_iops":        CHECKSET,
			"instance_classes_infos.0.memory_class":    CHECKSET,
			"instance_classes_infos.0.reference_price": CHECKSET,
		}
	}

	var existDBInstanceReadOnlyMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                         CHECKSET,
			"instance_classes_infos.#":                      CHECKSET,
			"instance_classes_infos.0.class_code":           CHECKSET,
			"instance_classes_infos.0.class_group":          CHECKSET,
			"instance_classes_infos.0.cpu":                  CHECKSET,
			"instance_classes_infos.0.max_connections":      CHECKSET,
			"instance_classes_infos.0.max_iombps":           CHECKSET,
			"instance_classes_infos.0.max_iops":             CHECKSET,
			"instance_classes_infos.0.memory_class":         CHECKSET,
			"instance_classes_infos.0.instruction_set_arch": CHECKSET,
		}
	}

	var fakeDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_classes_infos.#": "0",
		}
	}

	var DBInstanceMainCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_db_instance_classes_info.default",
		existMapFunc: existDBInstanceMainMapFunc,
		fakeMapFunc:  fakeDBInstanceMapFunc,
	}

	var DBInstanceReadOnlyCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_db_instance_classes_info.default",
		existMapFunc: existDBInstanceReadOnlyMapFunc,
		fakeMapFunc:  fakeDBInstanceMapFunc,
	}

	DBInstanceMainCheckInfo.dataSourceTestCheck(t, rand, BardsInstanceClassesConf, RdsInstanceClassesConf)
	DBInstanceReadOnlyCheckInfo.dataSourceTestCheck(t, rand, RordsInstanceClassesConf, RdsRordspreInstanceClassesConf)
}

func testAccCheckAlicloudDBInstanceClassesInfoDataSourceConfig(name string) string {
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

data "alicloud_db_instance_classes" "master" {
    zone_id                  = data.alicloud_db_zones.default.zones.0.id
	engine                   = "SQLServer"
	engine_version           = "2017_ent"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type     = "PostPaid"
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
	instance_type            = data.alicloud_db_instance_classes.master.instance_classes.0.instance_class
	instance_storage         = data.alicloud_db_instance_classes.master.instance_classes.0.storage_range.min
	vswitch_id               = data.alicloud_vswitches.default.vswitches.0.vswitch_id
	instance_name            = var.name
}
`, name)
}
