package rkvstore_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudKvstoreInstanceEngines(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_kvstore_instance_engines.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, "KVStore", kvstoreConfigHeader)

	EngineVersionConfRedis := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":        "${data.alicloud_zones.resources.zones.0.id}",
			"engine":         "Redis",
			"engine_version": "5.0",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":        "${data.alicloud_zones.resources.zones.0.id}",
			"engine":         "Redis",
			"engine_version": "4.9",
		}),
	}

	EngineVersionConfMemcache := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${data.alicloud_zones.resources.zones.0.id}",
			"engine":  "Memcache",
		}),
	}

	ChargeTypeConfPostpaid := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_zones.resources.zones.0.id}",
			"instance_charge_type": "PostPaid",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_zones.resources.zones.0.id}",
			"instance_charge_type": "PostPaid",
			"engine":               "Redis",
			"engine_version":       "5.0",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_zones.resources.zones.0.id}",
			"instance_charge_type": "PostPaid",
			"engine":               "Redis",
			"engine_version":       "5.6",
		}),
	}

	var existKVStoreInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_engines.#":                acctest.CHECKSET,
			"instance_engines.0.engine":         acctest.CHECKSET,
			"instance_engines.0.zone_id":        acctest.CHECKSET,
			"instance_engines.0.engine_version": acctest.CHECKSET,
		}
	}

	var fakeKVStoreInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_engines.#": "0",
		}
	}

	var KVStoreInstanceCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_kvstore_instance_engines.default",
		ExistMapFunc: existKVStoreInstanceMapFunc,
		FakeMapFunc:  fakeKVStoreInstanceMapFunc,
	}
	// At present, the datasource does not support memcache
	//KVStoreInstanceCheckInfo.DataSourceTestCheck(t, rand, EngineVersionConfRedis, EngineVersionConfMemcache, ChargeTypeConfPostpaid, allConf)
	KVStoreInstanceCheckInfo.DataSourceTestCheck(t, rand, EngineVersionConfRedis, ChargeTypeConfPostpaid, EngineVersionConfMemcache, allConf)
}
