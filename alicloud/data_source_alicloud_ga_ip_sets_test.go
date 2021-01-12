package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGaIpSetsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_ga_ip_sets.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceGaIpSetsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"ids":            []string{"${alicloud_ga_ip_set.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"ids":            []string{"${alicloud_ga_ip_set.default.id}_fake"},
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"ids":            []string{"${alicloud_ga_ip_set.default.id}"},
			"status":         "active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"ids":            []string{"${alicloud_ga_ip_set.default.id}_fake"},
			"status":         "init",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"ids":            []string{"${alicloud_ga_ip_set.default.id}"},
			"status":         "active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"accelerator_id": "${data.alicloud_ga_accelerators.default.ids.0}",
			"ids":            []string{"${alicloud_ga_ip_set.default.id}_fake"},
			"status":         "init",
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"sets.#":                      CHECKSET,
			"sets.0.accelerate_region_id": defaultRegionToTest,
			"sets.0.bandwidth":            "5",
			"sets.0.ip_address_list.#":    CHECKSET,
			"sets.0.id":                   CHECKSET,
			"sets.0.ip_set_id":            CHECKSET,
			"sets.0.ip_version":           "IPv4",
			"sets.0.status":               "active",
		}
	}

	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"sets.#": "0",
			"ids.#":  "0",
		}
	}

	var CheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithNoDefaultVpc(t)
	}

	CheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, allConf)
}

func dataSourceGaIpSetsConfigDependence(name string) string {
	return fmt.Sprintf(`
data "alicloud_ga_accelerators" "default"{
}
data "alicloud_ga_bandwidth_packages" "default"{
}
resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = "${data.alicloud_ga_accelerators.default.ids.0}"
  bandwidth_package_id = "${data.alicloud_ga_bandwidth_packages.default.ids.0}"
}
resource "alicloud_ga_ip_set" "default" {
  depends_on           = [alicloud_ga_bandwidth_package_attachment.default]
  accelerate_region_id = "%s"
  bandwidth            = "5"
  accelerator_id       = "${data.alicloud_ga_accelerators.default.ids.0}"
}`, defaultRegionToTest)
}
