package alicloud

import (
	"testing"
)

func TestAccAlicloudFcZonesDataSource(t *testing.T) {

	testAccConfig := dataSourceTestAccConfigFunc("data.alicloud_fc_zones.default", "",
		func(name string) string {
			return ""
		})

	multiConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"multi": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"multi": "false",
		}),
	}
	instanceChargeTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type": "PrePaid",
		}),
	}

	fcZonesMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"zones.#":                  CHECKSET,
			"zones.0.id":               CHECKSET,
			"zones.0.multi_zone_ids.#": CHECKSET,
			"ids.#":                    CHECKSET,
		}
	}

	fcZonesCheckInfo := dataSourceAttr{
		resourceId:   "data.alicloud_fc_zones.default",
		existMapFunc: fcZonesMapFunc,
		fakeMapFunc:  fcZonesMapFunc,
	}

	fcZonesCheckInfo.dataSourceTestCheck(t, 0, multiConf, instanceChargeTypeConf)
}
