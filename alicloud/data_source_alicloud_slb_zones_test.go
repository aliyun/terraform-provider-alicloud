package alicloud

import (
	"testing"
)

func TestAccAlicloudSlbZonesDataSource(t *testing.T) {

	testAccConfig := dataSourceTestAccConfigFunc("data.alicloud_slb_zones.default", "",
		func(name string) string {
			return ""
		})

	instanceChargeTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type": "PrePaid",
		}),
	}

	slbZonesMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    CHECKSET,
			"ids.0":                    CHECKSET,
			"zones.#":                  CHECKSET,
			"zones.0.zone_id":          CHECKSET,
			"zones.0.slave_zone_ids.#": CHECKSET,
		}
	}

	slbZonesCheckInfo := dataSourceAttr{
		resourceId:   "data.alicloud_slb_zones.default",
		existMapFunc: slbZonesMapFunc,
		fakeMapFunc:  slbZonesMapFunc,
	}

	slbZonesCheckInfo.dataSourceTestCheck(t, 0, instanceChargeTypeConf)
}
