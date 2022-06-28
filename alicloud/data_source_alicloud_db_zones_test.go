package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRdsDBZonesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_db_zones.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceDBZonesConfigDependence)

	multiZoneTrueConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"multi_zone": "true",
		}),
	}

	multiZoneFalseConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"multi_zone": "false",
		}),
	}
	chargeTypePostPaidConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type": "PostPaid",
		}),
	}
	chargeTypePrePaidConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type": "PrePaid",
		}),
	}
	engineVersionConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine":         "MySQL",
			"engine_version": "8.0",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"engine":         "MySQL",
			"engine_version": "2.0",
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

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"multi_zone":               "true",
			"instance_charge_type":     "PostPaid",
			"engine":                   "MySQL",
			"engine_version":           "8.0",
			"category":                 "HighAvailability",
			"db_instance_storage_type": "local_ssd",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
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
			"ids.#":                    CHECKSET,
			"ids.0":                    CHECKSET,
			"zones.#":                  CHECKSET,
			"zones.0.id":               CHECKSET,
			"zones.0.multi_zone_ids.#": CHECKSET,
		}
	}

	var fakeDBZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var DBZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDBZonesMapFunc,
		fakeMapFunc:  fakeDBZonesMapFunc,
	}

	DBZonesCheckInfo.dataSourceTestCheck(t, rand, multiZoneTrueConfig, multiZoneFalseConfig, chargeTypePostPaidConfig, chargeTypePrePaidConfig, engineVersionConfig, CategoryConf, StorageTypeConf, allConfig)
}

func dataSourceDBZonesConfigDependence(name string) string {
	return ""
}
