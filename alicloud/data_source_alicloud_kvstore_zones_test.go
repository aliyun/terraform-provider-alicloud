package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudKVStoreZonesDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_kvstore_zones.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceKVStoreZonesDependence)
	// // Redis has no multi zone
	//multiConfig := dataSourceTestAccConfig{
	//	existConfig: testAccConfig(map[string]interface{}{
	//		"multi": "true",
	//	}),
	//}

	chargeTypeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type": "PrePaid",
		}),
	}

	allconfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			//"multi":                "true",
			"instance_charge_type": "PrePaid",
		}),
	}

	var existKVStoreDBZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    CHECKSET,
			"ids.0":                    CHECKSET,
			"zones.#":                  CHECKSET,
			"zones.0.id":               CHECKSET,
			"zones.0.multi_zone_ids.#": NOSET,
		}
	}

	var fakeKVStoreDBZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var kvStoreDBZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKVStoreDBZonesMapFunc,
		fakeMapFunc:  fakeKVStoreDBZonesMapFunc,
	}

	kvStoreDBZonesCheckInfo.dataSourceTestCheck(t, rand, chargeTypeConfig, allconfig)
}

func dataSourceKVStoreZonesDependence(name string) string {
	return ""
}
