package rds_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudRdsDBEnginesDatasource(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_db_instance_engines.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, "", testAccCheckAlicloudDBEnginesDataSourceConfig)

	ZoneIDConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${data.alicloud_db_zones.default.zones.0.id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"zone_id": "fake_zoneid",
		}),
	}

	ChargeTypeConfPostpaid := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type": "PostPaid",
			"zone_id":              "${data.alicloud_db_zones.default.zones.0.id}",
		}),
	}
	ChargeTypeConfPrepaid := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type": "PrePaid",
			"zone_id":              "${data.alicloud_db_zones.default.zones.0.id}",
		}),
	}
	EngineConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"engine": "MySQL",
		}),
	}
	EngineVersionConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"engine":         "MySQL",
			"engine_version": "8.0",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"engine":         "MySQL",
			"engine_version": "3.0",
		}),
	}
	CategoryConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"engine_version": "8.0",
			"category":       "HighAvailability",
		}),
	}
	StorageTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"db_instance_storage_type": "local_ssd",
		}),
	}
	multiZoneConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"multi_zone": "true",
		}),
	}
	falseMultiZoneConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"multi_zone": "false",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type":     "PostPaid",
			"engine":                   "MySQL",
			"engine_version":           "8.0",
			"zone_id":                  "${data.alicloud_db_zones.default.zones.0.id}",
			"category":                 "HighAvailability",
			"db_instance_storage_type": "local_ssd",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"ids.#":                             acctest.CHECKSET,
			"instance_engines.#":                acctest.CHECKSET,
			"instance_engines.0.engine":         acctest.CHECKSET,
			"instance_engines.0.zone_ids.0.id":  acctest.CHECKSET,
			"instance_engines.0.engine_version": acctest.CHECKSET,
			"instance_engines.0.category":       acctest.CHECKSET,
		}
	}

	var fakeDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_engines.#": "0",
			"ids.#":              "0",
		}
	}

	var DBInstanceCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_db_instance_engines.default",
		ExistMapFunc: existDBInstanceMapFunc,
		FakeMapFunc:  fakeDBInstanceMapFunc,
	}
	DBInstanceCheckInfo.DataSourceTestCheck(t, rand, ZoneIDConf, ChargeTypeConfPostpaid, ChargeTypeConfPrepaid, EngineConf, EngineVersionConf, CategoryConf, StorageTypeConf, multiZoneConf, falseMultiZoneConf, allConf)
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
