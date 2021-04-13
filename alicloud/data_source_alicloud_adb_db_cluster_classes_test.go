package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAdbDbClusterClassesDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_adb_db_cluster_classes.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceAdbDbClusterClassesConfigDependence)

	paymentTypeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"payment_type": "PayAsYouGo",
		}),
	}

	var existAdbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"available_zone_list.#": CHECKSET,
		}
	}

	var fakeAdbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"available_zone_list.#": "0",
		}
	}

	var adbZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAdbZonesMapFunc,
		fakeMapFunc:  fakeAdbZonesMapFunc,
	}

	adbZonesCheckInfo.dataSourceTestCheck(t, rand, paymentTypeConfig)
}

func dataSourceAdbDbClusterClassesConfigDependence(name string) string {
	return ""
}
