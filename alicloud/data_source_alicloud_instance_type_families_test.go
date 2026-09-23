package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECSInstanceTypeFamiliesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_instance_type_families.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccInstanceTypeFamiliesDataSource_%d", rand),
		dataSourceInstanceTypeFamiliesConfigDependence)

	zoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${data.alicloud_zones.default.zones.0.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${data.alicloud_zones.default.zones.0.id}-F",
		}),
	}

	generationConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"generation": "ecs-3",
		}),
	}

	instanceChargeTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type": "PrePaid",
		}),
	}

	spotStrategyConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"spot_strategy": "SpotAsPriceGo",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_zones.default.zones.0.id}",
			"generation":           "ecs-3",
			"instance_charge_type": "PrePaid",
			"spot_strategy":        "SpotAsPriceGo",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id":              "${data.alicloud_zones.default.zones.0.id}-F",
			"generation":           "ecs-3",
			"instance_charge_type": "PrePaid",
			"spot_strategy":        "SpotAsPriceGo",
		}),
	}

	var existInstanceTypeFamiliesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 CHECKSET,
			"ids.0":                 REGEXMATCH + "^ecs.*",
			"families.#":            CHECKSET,
			"families.0.id":         REGEXMATCH + "^ecs.*",
			"families.0.generation": REGEXMATCH + "^ecs-.*",
			"families.0.zone_ids.#": CHECKSET,
			"families.0.zone_ids.0": CHECKSET,
		}
	}

	var fakeInstanceTypeFamiliesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"families.#": "0",
		}
	}

	var instanceTypeFamiliesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existInstanceTypeFamiliesMapFunc,
		fakeMapFunc:  fakeInstanceTypeFamiliesMapFunc,
	}

	instanceTypeFamiliesCheckInfo.dataSourceTestCheck(t, rand, zoneIdConf, generationConf, instanceChargeTypeConf, spotStrategyConf, allConf)
}

func dataSourceInstanceTypeFamiliesConfigDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_zones" "default" {
	  available_resource_creation = "Instance"
	}
`)
}
