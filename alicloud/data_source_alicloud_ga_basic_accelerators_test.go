package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGaBasicAcceleratorsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaBasicAcceleratorsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_basic_accelerator.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudGaBasicAcceleratorsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_basic_accelerator.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaBasicAcceleratorsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ga_basic_accelerator.default.basic_accelerator_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudGaBasicAcceleratorsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ga_basic_accelerator.default.basic_accelerator_name}_fake"`,
		}),
	}
	acceleratorIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaBasicAcceleratorsDataSourceName(rand, map[string]string{
			"accelerator_id": `"${alicloud_ga_basic_accelerator.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudGaBasicAcceleratorsDataSourceName(rand, map[string]string{
			"accelerator_id": `"${alicloud_ga_basic_accelerator.default.id}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaBasicAcceleratorsDataSourceName(rand, map[string]string{
			"accelerator_id": `"${alicloud_ga_basic_accelerator.default.id}"`,
			"status":         `"active"`,
		}),
		fakeConfig: testAccCheckAlicloudGaBasicAcceleratorsDataSourceName(rand, map[string]string{
			"accelerator_id": `"${alicloud_ga_basic_accelerator.default.id}_fake"`,
			"status":         `"deleting"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaBasicAcceleratorsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_ga_basic_accelerator.default.id}"]`,
			"name_regex":     `"${alicloud_ga_basic_accelerator.default.basic_accelerator_name}"`,
			"accelerator_id": `"${alicloud_ga_basic_accelerator.default.id}"`,
			"status":         `"active"`,
		}),
		fakeConfig: testAccCheckAlicloudGaBasicAcceleratorsDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_ga_basic_accelerator.default.id}_fake"]`,
			"name_regex":     `"${alicloud_ga_basic_accelerator.default.basic_accelerator_name}_fake"`,
			"accelerator_id": `"${alicloud_ga_basic_accelerator.default.id}_fake"`,
			"status":         `"deleting"`,
		}),
	}
	var existAlicloudGaBasicAcceleratorsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                           "1",
			"names.#":                                         "1",
			"accelerators.#":                                  "1",
			"accelerators.0.id":                               CHECKSET,
			"accelerators.0.basic_accelerator_id":             CHECKSET,
			"accelerators.0.basic_accelerator_name":           CHECKSET,
			"accelerators.0.basic_endpoint_group_id":          "",
			"accelerators.0.basic_ip_set_id":                  "",
			"accelerators.0.bandwidth_billing_type":           "BandwidthPackage",
			"accelerators.0.instance_charge_type":             "PREPAY",
			"accelerators.0.description":                      CHECKSET,
			"accelerators.0.region_id":                        CHECKSET,
			"accelerators.0.create_time":                      CHECKSET,
			"accelerators.0.expired_time":                     CHECKSET,
			"accelerators.0.status":                           "active",
			"accelerators.0.basic_bandwidth_package.#":        "1",
			"accelerators.0.cross_domain_bandwidth_package.#": "1",
		}
	}
	var fakeAlicloudGaBasicAcceleratorsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"names.#":        "0",
			"accelerators.#": "0",
		}
	}
	var alicloudGaBasicAcceleratorsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ga_basic_accelerators.default",
		existMapFunc: existAlicloudGaBasicAcceleratorsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudGaBasicAcceleratorsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithTime(t, []int{1})
	}
	alicloudGaBasicAcceleratorsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, acceleratorIdConf, statusConf, allConf)
}

func testAccCheckAlicloudGaBasicAcceleratorsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccGaBasicAccelerator-%d"
	}

	resource "alicloud_ga_basic_accelerator" "default" {
		duration               = 1
  		pricing_cycle          = "Month"
  		basic_accelerator_name = var.name
  		description            = var.name
  		bandwidth_billing_type = "BandwidthPackage"
  		auto_pay               = true
  		auto_use_coupon        = "true"
  		auto_renew             = false
  		auto_renew_duration    = 1
	}

	data "alicloud_ga_basic_accelerators" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
