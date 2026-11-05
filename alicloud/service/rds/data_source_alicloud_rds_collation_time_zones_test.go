package rds_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudRdsCollationTimeZonesDatasource(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_rds_collation_time_zones.default"
	name := "tf-testAccAlicloudRdsCollationTimeZones"
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, testAccAlicloudRdsCollationTimeZonesConfig)

	rdsCollationTimeZonesConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{}),
	}

	var existRdsCollationTimeZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  acctest.CHECKSET,
			"ids.0":                  acctest.CHECKSET,
			"collation_time_zones.#": acctest.CHECKSET,
			"collation_time_zones.0.standard_time_offset": acctest.CHECKSET,
			"collation_time_zones.0.description":          acctest.CHECKSET,
			"collation_time_zones.0.time_zone":            acctest.CHECKSET,
		}
	}

	var fakeRdsCollationTimeZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "0",
			"collation_time_zones.#": "0",
		}
	}

	var RdsCollationTimeZonesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existRdsCollationTimeZonesMapFunc,
		FakeMapFunc:  fakeRdsCollationTimeZonesMapFunc,
	}

	RdsCollationTimeZonesCheckInfo.DataSourceTestCheck(t, rand, rdsCollationTimeZonesConfig)
}

func testAccAlicloudRdsCollationTimeZonesConfig(name string) string {
	return ""
}
