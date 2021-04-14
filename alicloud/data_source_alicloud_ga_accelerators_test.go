package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGaAcceleratorsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_ga_accelerators.default"
	name := fmt.Sprintf("tf-testAccelerators_datasource-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceGaAcceleratorsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ga_accelerator.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ga_accelerator.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ga_accelerator.default.accelerator_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ga_accelerator.default.accelerator_name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ga_accelerator.default.id}"},
			"name_regex": "${alicloud_ga_accelerator.default.accelerator_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_ga_accelerator.default.id}_fake"},
			"name_regex": "${alicloud_ga_accelerator.default.accelerator_name}_fake",
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                           "1",
			"names.#":                                         "1",
			"accelerators.#":                                  CHECKSET,
			"accelerators.0.id":                               CHECKSET,
			"accelerators.0.accelerator_name":                 fmt.Sprintf("tf-testAccelerators_datasource-%d", rand),
			"accelerators.0.cen_id":                           "",
			"accelerators.0.ddos_id":                          "",
			"accelerators.0.description":                      "",
			"accelerators.0.dns_name":                         CHECKSET,
			"accelerators.0.expired_time":                     CHECKSET,
			"accelerators.0.payment_type":                     "PREPAY",
			"accelerators.0.second_dns_name":                  "",
			"accelerators.0.spec":                             "1",
			"accelerators.0.status":                           "active",
			"accelerators.0.cross_domain_bandwidth_package.#": "0",
			"accelerators.0.basic_bandwidth_package.#":        "0",
		}
	}

	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"accelerators.#": "0",
			"ids.#":          "0",
		}
	}

	var CheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}
	preCheck := func() {}

	CheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceGaAcceleratorsConfigDependence(name string) string {
	return fmt.Sprintf(`resource "alicloud_ga_accelerator" "default" {
  spec ="1"
  auto_use_coupon=true
  duration=1
  accelerator_name="%s"
}`, name)
}
