package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudPolarDBZonesDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_polardb_zones.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourcePolarDBZonesConfigDependence)

	multiConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"multi": "true",
		}),
	}

	var existPolarDBZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   CHECKSET,
			"zones.#": CHECKSET,
		}
	}

	var fakePolarDBZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var polarDBZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPolarDBZonesMapFunc,
		fakeMapFunc:  fakePolarDBZonesMapFunc,
	}

	polarDBZonesCheckInfo.dataSourceTestCheck(t, rand, multiConfig)
}

func dataSourcePolarDBZonesConfigDependence(name string) string {
	return ""
}
