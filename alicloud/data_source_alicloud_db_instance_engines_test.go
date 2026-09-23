package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRdsDBEnginesDatasource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_db_instance_engines.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", testAccCheckAlicloudDBEnginesDataSourceConfig)

	ZoneIDConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${data.alicloud_db_zones.default.zones.0.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id": "fake_zoneid",
		}),
	}

	ChargeTypeConfPostpaid := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type": "PostPaid",
			"zone_id":              "${data.alicloud_db_zones.default.zones.0.id}",
		}),
	}
	ChargeTypeConfPrepaid := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type": "PrePaid",
			"zone_id":              "${data.alicloud_db_zones.default.zones.0.id}",
		}),
	}
	EngineConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine": "MySQL",
		}),
	}
	EngineVersionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine":         "MySQL",
			"engine_version": "8.0",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"engine":         "MySQL",
			"engine_version": "3.0",
		}),
	}
	CategoryConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine_version": "8.0",
			"category":       "HighAvailability",
		}),
	}
	StorageTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_storage_type": "local_ssd",
		}),
	}
	multiZoneConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"multi_zone": "true",
		}),
	}
	falseMultiZoneConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"multi_zone": "false",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type":     "PostPaid",
			"engine":                   "MySQL",
			"engine_version":           "8.0",
			"zone_id":                  "${data.alicloud_db_zones.default.zones.0.id}",
			"category":                 "HighAvailability",
			"db_instance_storage_type": "local_ssd",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":                  "${data.alicloud_db_zones.default.zones.0.id}",
			"instance_charge_type":     "PostPaid",
			"engine":                   "MySQL",
			"engine_version":           "3.0",
			"category":                 "HighAvailability",
			"db_instance_storage_type": "local_ssd",
		}),
	}

	var existDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             CHECKSET,
			"instance_engines.#":                CHECKSET,
			"instance_engines.0.engine":         CHECKSET,
			"instance_engines.0.zone_ids.0.id":  CHECKSET,
			"instance_engines.0.engine_version": CHECKSET,
			"instance_engines.0.category":       CHECKSET,
		}
	}

	var fakeDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_engines.#": "0",
			"ids.#":              "0",
		}
	}

	var DBInstanceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_db_instance_engines.default",
		existMapFunc: existDBInstanceMapFunc,
		fakeMapFunc:  fakeDBInstanceMapFunc,
	}
	DBInstanceCheckInfo.dataSourceTestCheck(t, rand, ZoneIDConf, ChargeTypeConfPostpaid, ChargeTypeConfPrepaid, EngineConf, EngineVersionConf, CategoryConf, StorageTypeConf, multiZoneConf, falseMultiZoneConf, allConf)
}

func testAccCheckAlicloudDBEnginesDataSourceConfig(name string) string {
	return fmt.Sprintf(`
data "alicloud_db_zones" "default" {
  instance_charge_type= "PostPaid"
  engine = "MySQL"
  db_instance_storage_type = "local_ssd"
}
`)
}
