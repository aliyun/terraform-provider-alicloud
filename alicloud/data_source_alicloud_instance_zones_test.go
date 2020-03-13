package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudInstanceZonesDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_instance_zones.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceInstanceZonesConfigDependence)

	instanceTypeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_instance_type": "ecs.n4.large",
		}),
	}

	diskCategoryConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_disk_category": "cloud_efficiency",
		}),
	}

	chargeTypeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type": "PrePaid",
		}),
	}

	networkTypeConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"network_type": "Vpc",
		}),
	}

	spotStrategyConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"spot_strategy": "SpotAsPriceGo",
		}),
	}

	allConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_instance_type": "ecs.n4.large",
			"available_disk_category": "cloud_efficiency",
			"instance_charge_type":    "PostPaid",
			"network_type":            "Vpc",
			"spot_strategy":           "SpotAsPriceGo",
		}),
	}

	var existInstanceZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 CHECKSET,
			"ids.0":                                 CHECKSET,
			"zones.#":                               CHECKSET,
			"zones.0.id":                            CHECKSET,
			"zones.0.available_disk_categories.#":   CHECKSET,
			"zones.0.available_instance_types.#":    CHECKSET,
			"zones.0.available_resource_creation.#": CHECKSET,
		}
	}

	var fakeInstanceZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var instanceZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existInstanceZonesMapFunc,
		fakeMapFunc:  fakeInstanceZonesMapFunc,
	}

	instanceZonesCheckInfo.dataSourceTestCheck(t, rand, instanceTypeConfig, diskCategoryConfig, chargeTypeConfig, networkTypeConfig, spotStrategyConfig, allConfig)
}

func dataSourceInstanceZonesConfigDependence(name string) string {
	return ""
}
