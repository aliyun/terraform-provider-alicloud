package gpdb_test

import (
	"testing"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
)

func TestAccAlicloudGPDBZonesDataSource_basic(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_gpdb_zones.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, "", dataSourceGpdbZonesConfigDependence)

	multiConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"multi": "true",
		}),
	}

	var existGpdbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   acctest.CHECKSET,
			"zones.#": acctest.CHECKSET,
		}
	}

	var fakeGpdbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var gpdbZonesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existGpdbZonesMapFunc,
		FakeMapFunc:  fakeGpdbZonesMapFunc,
	}

	gpdbZonesCheckInfo.DataSourceTestCheck(t, rand, multiConfig)
}

func dataSourceGpdbZonesConfigDependence(name string) string {
	return ""
}
