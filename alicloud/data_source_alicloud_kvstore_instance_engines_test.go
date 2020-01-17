package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudKVStoreInstanceEngines(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_kvstore_instance_engines.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "KVStore", kvstoreConfigHeader)

	EngineVersionConfRedis := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":        "${data.alicloud_zones.resources.zones.0.id}",
			"engine":         "Redis",
			"engine_version": "5.0",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":        "${data.alicloud_zones.resources.zones.0.id}",
			"engine":         "Redis",
			"engine_version": "4.9",
		}),
	}

	EngineVersionConfMemcache := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${data.alicloud_zones.resources.zones.0.id}",
			"engine":  "Memcache",
		}),
	}

	ChargeTypeConfPostpaid := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_zones.resources.zones.0.id}",
			"instance_charge_type": "PostPaid",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_zones.resources.zones.0.id}",
			"instance_charge_type": "PostPaid",
			"engine":               "Redis",
			"engine_version":       "5.0",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_zones.resources.zones.0.id}",
			"instance_charge_type": "PostPaid",
			"engine":               "Redis",
			"engine_version":       "5.6",
		}),
	}

	var existKVStoreInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_engines.#":                CHECKSET,
			"instance_engines.0.engine":         CHECKSET,
			"instance_engines.0.zone_id":        CHECKSET,
			"instance_engines.0.engine_version": CHECKSET,
		}
	}

	var fakeKVStoreInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_engines.#": "0",
		}
	}

	var KVStoreInstanceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_kvstore_instance_engines.default",
		existMapFunc: existKVStoreInstanceMapFunc,
		fakeMapFunc:  fakeKVStoreInstanceMapFunc,
	}
	// At present, the datasource does not support memcache
	//KVStoreInstanceCheckInfo.dataSourceTestCheck(t, rand, EngineVersionConfRedis, EngineVersionConfMemcache, ChargeTypeConfPostpaid, allConf)
	KVStoreInstanceCheckInfo.dataSourceTestCheck(t, rand, EngineVersionConfRedis, ChargeTypeConfPostpaid, EngineVersionConfMemcache, allConf)
}
