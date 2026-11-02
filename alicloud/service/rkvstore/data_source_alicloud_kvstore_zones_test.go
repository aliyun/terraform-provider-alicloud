package rkvstore_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudKvstoreZonesDataSource_basic(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_kvstore_zones.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, "", dataSourceKVStoreZonesDependence)
	// // Redis has no multi zone
	//multiConfig := acctest.DataSourceTestAccConfig{
	//	ExistConfig: testAccConfig(map[string]interface{}{
	//		"multi": "true",
	//	}),
	//}

	chargeTypeConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type": "PrePaid",
		}),
	}

	allconfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			//"multi":                "true",
			"instance_charge_type": "PrePaid",
		}),
	}

	var existKVStoreDBZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    acctest.CHECKSET,
			"ids.0":                    acctest.CHECKSET,
			"zones.#":                  acctest.CHECKSET,
			"zones.0.id":               acctest.CHECKSET,
			"zones.0.multi_zone_ids.#": acctest.NOSET,
		}
	}

	var fakeKVStoreDBZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var kvStoreDBZonesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existKVStoreDBZonesMapFunc,
		FakeMapFunc:  fakeKVStoreDBZonesMapFunc,
	}

	kvStoreDBZonesCheckInfo.DataSourceTestCheck(t, rand, chargeTypeConfig, allconfig)
}

func dataSourceKVStoreZonesDependence(name string) string {
	return ""
}
