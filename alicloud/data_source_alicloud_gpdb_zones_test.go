package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGPDBZonesDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_gpdb_zones.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceGpdbZonesConfigDependence)

	multiConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"multi": "true",
		}),
	}

	var existGpdbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   CHECKSET,
			"zones.#": CHECKSET,
		}
	}

	var fakeGpdbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var gpdbZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existGpdbZonesMapFunc,
		fakeMapFunc:  fakeGpdbZonesMapFunc,
	}

	gpdbZonesCheckInfo.dataSourceTestCheck(t, rand, multiConfig)
}

func dataSourceGpdbZonesConfigDependence(name string) string {
	return ""
}
