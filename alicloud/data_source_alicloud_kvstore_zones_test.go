package alicloud

import (
	"testing"
)

func TestAccAlicloudKVStoreZonesDataSource(t *testing.T) {

	testAccConfig := dataSourceTestAccConfigFunc("data.alicloud_kvstore_zones.default", "",
		func(name string) string {
			return ""
		})

	multiConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"multi": "false",
		}),
	}
	instanceChargeTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type": "PrePaid",
		}),
	}

	kvStoreZonesMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"ids.#": CHECKSET,
			"ids.0": CHECKSET,
		}
	}

	kvStoreZonesCheckInfo := dataSourceAttr{
		resourceId:   "data.alicloud_kvstore_zones.default",
		existMapFunc: kvStoreZonesMapFunc,
		fakeMapFunc:  kvStoreZonesMapFunc,
	}

	kvStoreZonesCheckInfo.dataSourceTestCheck(t, 0, multiConf, instanceChargeTypeConf)
}
