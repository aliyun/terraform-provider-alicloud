package rds_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudRdsDBInstanceClassesDatasource(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_db_instance_classes.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, "", testAccCheckAliCloudDBInstanceClassesDataSourceConfig)

	ZoneIDConf := acctest.DataSourceTestAccConfig{
		FakeConfig: testAccConfig(map[string]interface{}{
			"zone_id": "fake_zoneid",
		}),
	}
	EngineVersionConf := acctest.DataSourceTestAccConfig{
		FakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":        "${data.alicloud_db_zones.default.zones.0.id}",
			"engine":         "MySQL",
			"engine_version": "3.0",
		}),
	}

	ChargeTypeConfPrepaid := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_db_zones.default.zones.0.id}",
			"instance_charge_type": "PrePaid",
			"engine_version":       "8.0",
		}),
	}
	ChargeTypeConfPostpaid := acctest.DataSourceTestAccConfig{
		FakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_db_zones.default.zones.0.id}",
			"instance_charge_type": "PostPaid",
		}),
	}
	StorageTypeConf := acctest.DataSourceTestAccConfig{
		FakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":      "${data.alicloud_db_zones.default.zones.0.id}",
			"storage_type": "local_ssd",
		}),
	}

	CategoryConf := acctest.DataSourceTestAccConfig{
		FakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":  "${data.alicloud_db_zones.default.zones.0.id}",
			"category": "HighAvailability",
		}),
	}

	CommodityCodeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_db_zones.default.zones.0.id}",
			"instance_charge_type": "PrePaid",
			"engine_version":       "8.0",
			"commodity_code":       "bards",
		}),
	}

	GeneralEssdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":                  "${data.alicloud_db_zones.default.zones.0.id}",
			"engine":                   "MySQL",
			"engine_version":           "8.0",
			"category":                 "HighAvailability",
			"db_instance_storage_type": "general_essd",
			"instance_charge_type":     "PostPaid",
			"commodity_code":           "bards",
		}),
	}

	ServerlessConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":                  "${data.alicloud_db_zones.default.zones.0.id}",
			"engine":                   "MySQL",
			"engine_version":           "8.0",
			"category":                 "serverless_basic",
			"db_instance_storage_type": "cloud_essd",
			"instance_charge_type":     "Serverless",
			"commodity_code":           "rds_serverless_public_cn",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_db_zones.default.zones.0.id}",
			"instance_charge_type": "PostPaid",
			"storage_type":         "local_ssd",
			"category":             "HighAvailability",
			"engine":               "MySQL",
			"engine_version":       "8.0",
			"commodity_code":       "bards",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_db_zones.default.zones.0.id}",
			"instance_charge_type": "PostPaid",
			"engine":               "MySQL",
			"engine_version":       "5.0",
		}),
	}

	var existDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             acctest.CHECKSET,
			"instance_classes.#":                acctest.CHECKSET,
			"instance_classes.0.instance_class": acctest.CHECKSET,
		}
	}

	var fakeDBInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":              "0",
			"instance_classes.#": "0",
		}
	}

	var DBInstanceCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_db_instance_classes.default",
		ExistMapFunc: existDBInstanceMapFunc,
		FakeMapFunc:  fakeDBInstanceMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.RDSInstanceClassesSupportRegions)
	}

	DBInstanceCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, ZoneIDConf, EngineVersionConf, ChargeTypeConfPrepaid,
		ChargeTypeConfPostpaid, CategoryConf, StorageTypeConf, CommodityCodeConf, GeneralEssdConf, ServerlessConf, allConf)
}

func testAccCheckAliCloudDBInstanceClassesDataSourceConfig(name string) string {
	return fmt.Sprintf(`
data "alicloud_db_zones" "default" {
}
`)
}
