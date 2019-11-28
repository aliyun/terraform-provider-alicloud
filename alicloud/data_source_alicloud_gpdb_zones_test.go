package alicloud

import (
	"testing"
)

func TestAccAlicloudGpdbZonesDataSource(t *testing.T) {

	testAccConfig := dataSourceTestAccConfigFunc("data.alicloud_gpdb_zones.default", "",
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

	mongoDBZonesMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"ids.#": CHECKSET,
			"ids.0": CHECKSET,
		}
	}

	mongoDBZonesCheckInfo := dataSourceAttr{
		resourceId:   "data.alicloud_gpdb_zones.default",
		existMapFunc: mongoDBZonesMapFunc,
		fakeMapFunc:  mongoDBZonesMapFunc,
	}

	mongoDBZonesCheckInfo.dataSourceTestCheck(t, 0, multiConf, instanceChargeTypeConf)
}
