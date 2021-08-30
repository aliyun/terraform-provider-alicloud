package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRdsDBZonesDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_db_zones.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceDBZonesConfigDependence)

	multiConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"multi": "true",
		}),
	}

	chargeTypeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type": "PrePaid",
		}),
	}
	engineVersionConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine_version": "8.0",
		}),
	}

	allconfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"multi":                "true",
			"instance_charge_type": "PrePaid",
			"engine_version":       "8.0",
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

	DBZonesCheckInfo.dataSourceTestCheck(t, rand, multiConfig, chargeTypeConfig, engineVersionConfig, allconfig)
}

func dataSourceDBZonesConfigDependence(name string) string {
	return ""
}
