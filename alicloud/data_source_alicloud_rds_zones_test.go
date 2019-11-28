package alicloud

import (
	"testing"
)

func TestAccAlicloudRdsZonesDataSource(t *testing.T) {

	testAccConfig := dataSourceTestAccConfigFunc("data.alicloud_rds_zones.default", "",
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

	rdsZonesMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"ids.#": CHECKSET,
			"ids.0": CHECKSET,
		}
	}

	rdsZonesCheckInfo := dataSourceAttr{
		resourceId:   "data.alicloud_rds_zones.default",
		existMapFunc: rdsZonesMapFunc,
		fakeMapFunc:  rdsZonesMapFunc,
	}

	rdsZonesCheckInfo.dataSourceTestCheck(t, 0, multiConf, instanceChargeTypeConf)
}
