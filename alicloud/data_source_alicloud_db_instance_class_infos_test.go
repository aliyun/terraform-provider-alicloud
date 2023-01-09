package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRdsDBInstanceClassInfosDatasource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_db_instance_class_infos.default"
	name := "tf-testAccAlicloudRdsDBInstanceClassInfos"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlicloudDBInstanceClassInfosDataSourceConfig)

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
			"ids.#":                   CHECKSET,
			"infos.#":                 CHECKSET,
			"infos.0.class_code":      CHECKSET,
			"infos.0.class_group":     CHECKSET,
			"infos.0.cpu":             CHECKSET,
			"infos.0.max_connections": CHECKSET,
			"infos.0.max_iombps":      CHECKSET,
			"infos.0.max_iops":        CHECKSET,
			"infos.0.memory_class":    CHECKSET,
			"infos.0.reference_price": CHECKSET,
		}
	}

	var existDBInstanceReadOnlyMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        CHECKSET,
			"infos.#":                      CHECKSET,
			"infos.0.class_code":           CHECKSET,
			"infos.0.class_group":          CHECKSET,
			"infos.0.cpu":                  CHECKSET,
			"infos.0.max_connections":      CHECKSET,
			"infos.0.max_iombps":           CHECKSET,
			"infos.0.max_iops":             CHECKSET,
			"infos.0.memory_class":         CHECKSET,
			"infos.0.instruction_set_arch": CHECKSET,
		}
	}

	var fakeDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"infos.#": "0",
			"ids.#":   "0",
		}
	}

	var DBInstanceMainCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_db_instance_class_infos.default",
		existMapFunc: existDBInstanceMainMapFunc,
		fakeMapFunc:  fakeDBInstanceMapFunc,
	}

	var DBInstanceReadOnlyCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_db_instance_class_infos.default",
		existMapFunc: existDBInstanceReadOnlyMapFunc,
		fakeMapFunc:  fakeDBInstanceMapFunc,
	}

	DBInstanceMainCheckInfo.dataSourceTestCheck(t, rand, BardsInstanceClassesConf, RdsInstanceClassesConf)
	DBInstanceReadOnlyCheckInfo.dataSourceTestCheck(t, rand, RordsInstanceClassesConf, RdsRordspreInstanceClassesConf)
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
