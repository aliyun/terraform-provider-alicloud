package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGaIpSetsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_ga_ip_sets.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceGaIpSetsConfigDependence)
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
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

	CheckInfo.dataSourceTestCheck(t, rand, idsConf, statusConf, allConf)
}

func dataSourceGaIpSetsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default  = "%s"
}
data "alicloud_ga_accelerators" "default" {
  status = "active"
}

resource "alicloud_ga_bandwidth_package" "default" {
   	bandwidth              =  100
  	type                   = "Basic"
  	bandwidth_type         = "Basic"
	payment_type           = "PayAsYouGo"
  	billing_type           = "PayBy95"
	ratio       = 30
	bandwidth_package_name = var.name
    auto_pay               = true
    auto_use_coupon        = true
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
	// Please run resource ga_accelerator test case to ensure this account has at least one accelerator before run this case.
	accelerator_id = data.alicloud_ga_accelerators.default.ids.0
	bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_ip_set" "default" {
  accelerate_region_id = "%s"
  bandwidth            = "5"
  accelerator_id       = "${alicloud_ga_bandwidth_package_attachment.default.accelerator_id}"
}`, name, defaultRegionToTest)
}
