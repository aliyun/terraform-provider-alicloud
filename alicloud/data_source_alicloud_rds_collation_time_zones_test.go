package alicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"testing"
)

func TestAccAlicloudRdsCollationTimeZonesDatasource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_rds_collation_time_zones.default"
	name := "tf-testAccAlicloudRdsCollationTimeZones"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccAlicloudRdsCollationTimeZonesConfig)

	rdsCollationTimeZonesConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{}),
	}

	var existRdsCollationTimeZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  CHECKSET,
			"ids.0":                  CHECKSET,
			"collation_time_zones.#": CHECKSET,
			"collation_time_zones.0.standard_time_offset": CHECKSET,
			"collation_time_zones.0.description":          CHECKSET,
			"collation_time_zones.0.time_zone":            CHECKSET,
		}
	}

	var fakeRdsCollationTimeZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "0",
			"collation_time_zones.#": "0",
		}
	}

	var RdsCollationTimeZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existRdsCollationTimeZonesMapFunc,
		fakeMapFunc:  fakeRdsCollationTimeZonesMapFunc,
	}

	RdsCollationTimeZonesCheckInfo.dataSourceTestCheck(t, rand, rdsCollationTimeZonesConfig)
}

func testAccAlicloudRdsCollationTimeZonesConfig(name string) string {
	return ""
}
