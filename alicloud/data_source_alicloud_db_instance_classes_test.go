package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRdsDBInstanceClassesDatasource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_db_instance_classes.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", testAccCheckAlicloudDBInstanceClassesDataSourceConfig)

	ZoneIDConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${data.alicloud_db_zones.default.zones.0.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id": "fake_zoneid",
		}),
	}
	EngineVersionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":        "${data.alicloud_db_zones.default.zones.0.id}",
			"engine":         "MySQL",
			"engine_version": "8.0",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":        "${data.alicloud_db_zones.default.zones.0.id}",
			"engine":         "MySQL",
			"engine_version": "3.0",
		}),
	}

	ChargeTypeConfPrepaid := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_db_zones.default.zones.0.id}",
			"instance_charge_type": "PrePaid",
			"engine_version":       "8.0",
		}),
	}
	ChargeTypeConfPostpaid := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_db_zones.default.zones.0.id}",
			"instance_charge_type": "PostPaid",
		}),
	}
	StorageTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":      "${data.alicloud_db_zones.default.zones.0.id}",
			"storage_type": "local_ssd",
		}),
	}

	CategoryConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":  "${data.alicloud_db_zones.default.zones.0.id}",
			"category": "HighAvailability",
		}),
	}

	CommodityCodeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":        "${data.alicloud_db_zones.default.zones.0.id}",
			"commodity_code": "bards",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_db_zones.default.zones.0.id}",
			"instance_charge_type": "PostPaid",
			"storage_type":         "local_ssd",
			"category":             "HighAvailability",
			"engine":               "MySQL",
			"engine_version":       "8.0",
			"commodity_code":       "bards",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_db_zones.default.zones.0.id}",
			"instance_charge_type": "PostPaid",
			"engine":               "MySQL",
			"engine_version":       "5.0",
		}),
	}

	var existDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_classes.#":                    CHECKSET,
			"instance_classes.0.instance_class":     CHECKSET,
			"instance_classes.0.storage_range.min":  CHECKSET,
			"instance_classes.0.storage_range.max":  CHECKSET,
			"instance_classes.0.storage_range.step": CHECKSET,
			"instance_classes.0.zone_ids.0.id":      CHECKSET,
		}
	}

	var fakeDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_classes.#": "0",
		}
	}

	var DBInstanceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_db_instance_classes.default",
		existMapFunc: existDBInstanceMapFunc,
		fakeMapFunc:  fakeDBInstanceMapFunc,
	}

	DBInstanceCheckInfo.dataSourceTestCheck(t, rand, ZoneIDConf, EngineVersionConf, ChargeTypeConfPrepaid,
		ChargeTypeConfPostpaid, CategoryConf, StorageTypeConf, CommodityCodeConf, allConf)
}

func testAccCheckAlicloudDBInstanceClassesDataSourceConfig(name string) string {
	return fmt.Sprintf(`
data "alicloud_db_zones" "default" {
  instance_charge_type= "PostPaid"
  engine = "MySQL"
  db_instance_storage_type = "local_ssd"
}
data "alicloud_db_zones" "true" {
  instance_charge_type= "PostPaid"
  engine = "MySQL"
  db_instance_storage_type = "local_ssd"
  multi = true
}
`)
}
