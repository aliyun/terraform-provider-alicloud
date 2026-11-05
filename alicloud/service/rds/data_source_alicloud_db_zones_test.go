package rds_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudRdsDBZonesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_db_zones.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, "", dataSourceDBZonesConfigDependence)

	multiZoneTrueConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"multi_zone": "true",
		}),
	}

	multiZoneFalseConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"multi_zone": "false",
		}),
	}
	chargeTypePostPaidConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type": "PostPaid",
		}),
	}
	chargeTypePrePaidConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type": "PrePaid",
		}),
	}
	engineVersionConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"engine":         "MySQL",
			"engine_version": "8.0",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"engine":         "MySQL",
			"engine_version": "2.0",
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

	ServerlessConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"engine":                   "MySQL",
			"engine_version":           "8.0",
			"instance_charge_type":     "Serverless",
			"category":                 "serverless_basic",
			"db_instance_storage_type": "cloud_essd",
		}),
	}

	allConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"multi_zone":               "true",
			"instance_charge_type":     "PostPaid",
			"engine":                   "MySQL",
			"engine_version":           "8.0",
			"category":                 "HighAvailability",
			"db_instance_storage_type": "local_ssd",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"multi_zone":               "true",
			"instance_charge_type":     "PostPaid",
			"engine":                   "MySQL",
			"engine_version":           "2.0",
			"category":                 "HighAvailability",
			"db_instance_storage_type": "local_ssd",
		}),
	}

	var existDBZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    acctest.CHECKSET,
			"ids.0":                    acctest.CHECKSET,
			"zones.#":                  acctest.CHECKSET,
			"zones.0.id":               acctest.CHECKSET,
			"zones.0.multi_zone_ids.#": acctest.CHECKSET,
		}
	}

	var fakeDBZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var DBZonesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existDBZonesMapFunc,
		FakeMapFunc:  fakeDBZonesMapFunc,
	}

	DBZonesCheckInfo.DataSourceTestCheck(t, rand, multiZoneTrueConfig, multiZoneFalseConfig, chargeTypePostPaidConfig,
		chargeTypePrePaidConfig, engineVersionConfig, CategoryConf, StorageTypeConf, ServerlessConf, allConfig)
}

func dataSourceDBZonesConfigDependence(name string) string {
	return ""
}
