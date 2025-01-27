package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudGaBasicAcceleratorsDataSource_basic0(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"ids": `["${alicloud_ga_basic_accelerator.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"ids": `["${alicloud_ga_basic_accelerator.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"name_regex": `"${alicloud_ga_basic_accelerator.default.basic_accelerator_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"name_regex": `"${alicloud_ga_basic_accelerator.default.basic_accelerator_name}_fake"`,
		}),
	}

	acceleratorIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"accelerator_id": `"${alicloud_ga_basic_accelerator.default.id}"`,
		}),
		fakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"accelerator_id": `"${alicloud_ga_basic_accelerator.default.id}_fake"`,
		}),
	}

	bandwidthBillingTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"accelerator_id":         `"${alicloud_ga_basic_accelerator.default.id}"`,
			"bandwidth_billing_type": `"BandwidthPackage"`,
		}),
		fakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"bandwidth_billing_type": `"CDT95"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"accelerator_id": `"${alicloud_ga_basic_accelerator.default.id}"`,
			"status":         `"active"`,
		}),
		fakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"status": `"deleting"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"ids":                    `["${alicloud_ga_basic_accelerator.default.id}"]`,
			"name_regex":             `"${alicloud_ga_basic_accelerator.default.basic_accelerator_name}"`,
			"accelerator_id":         `"${alicloud_ga_basic_accelerator.default.id}"`,
			"bandwidth_billing_type": `"BandwidthPackage"`,
			"status":                 `"active"`,
		}),
		fakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"ids":                    `["${alicloud_ga_basic_accelerator.default.id}_fake"]`,
			"name_regex":             `"${alicloud_ga_basic_accelerator.default.basic_accelerator_name}_fake"`,
			"accelerator_id":         `"${alicloud_ga_basic_accelerator.default.id}_fake"`,
			"bandwidth_billing_type": `"CDT95"`,
			"status":                 `"deleting"`,
		}),
	}

	var existAliCloudGaBasicAcceleratorsDataSourceNameMapFunc = func(rand int) map[string]string {
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

	var fakeAliCloudGaBasicAcceleratorsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"names.#":        "0",
			"accelerators.#": "0",
		}
	}

	var alicloudGaBasicAcceleratorsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ga_basic_accelerators.default",
		existMapFunc: existAliCloudGaBasicAcceleratorsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudGaBasicAcceleratorsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		//testAccPreCheckWithTime(t, []int{1})
	}

	alicloudGaBasicAcceleratorsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, acceleratorIdConf, bandwidthBillingTypeConf, statusConf, allConf)
}

func TestAccAliCloudGaBasicAcceleratorsDataSource_basic1(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"ids": `["${alicloud_ga_basic_accelerator.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"ids": `["${alicloud_ga_basic_accelerator.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"name_regex": `"${alicloud_ga_basic_accelerator.default.basic_accelerator_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"name_regex": `"${alicloud_ga_basic_accelerator.default.basic_accelerator_name}_fake"`,
		}),
	}

	acceleratorIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"accelerator_id": `"${alicloud_ga_basic_accelerator.default.id}"`,
		}),
		fakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"accelerator_id": `"${alicloud_ga_basic_accelerator.default.id}_fake"`,
		}),
	}

	bandwidthBillingTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"accelerator_id":         `"${alicloud_ga_basic_accelerator.default.id}"`,
			"bandwidth_billing_type": `"CDT"`,
		}),
		fakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"bandwidth_billing_type": `"CDT95"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"accelerator_id": `"${alicloud_ga_basic_accelerator.default.id}"`,
			"status":         `"active"`,
		}),
		fakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"status": `"deleting"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"ids":                    `["${alicloud_ga_basic_accelerator.default.id}"]`,
			"name_regex":             `"${alicloud_ga_basic_accelerator.default.basic_accelerator_name}"`,
			"accelerator_id":         `"${alicloud_ga_basic_accelerator.default.id}"`,
			"bandwidth_billing_type": `"CDT"`,
			"status":                 `"active"`,
		}),
		fakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"ids":                    `["${alicloud_ga_basic_accelerator.default.id}_fake"]`,
			"name_regex":             `"${alicloud_ga_basic_accelerator.default.basic_accelerator_name}_fake"`,
			"accelerator_id":         `"${alicloud_ga_basic_accelerator.default.id}_fake"`,
			"bandwidth_billing_type": `"CDT95"`,
			"status":                 `"deleting"`,
		}),
	}

	var existAliCloudGaBasicAcceleratorsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                           "1",
			"names.#":                                         "1",
			"accelerators.#":                                  "1",
			"accelerators.0.id":                               CHECKSET,
			"accelerators.0.basic_accelerator_id":             CHECKSET,
			"accelerators.0.basic_accelerator_name":           CHECKSET,
			"accelerators.0.basic_endpoint_group_id":          "",
			"accelerators.0.basic_ip_set_id":                  "",
			"accelerators.0.bandwidth_billing_type":           "CDT",
			"accelerators.0.instance_charge_type":             "POSTPAY",
			"accelerators.0.description":                      CHECKSET,
			"accelerators.0.region_id":                        CHECKSET,
			"accelerators.0.create_time":                      CHECKSET,
			"accelerators.0.expired_time":                     CHECKSET,
			"accelerators.0.status":                           "active",
			"accelerators.0.basic_bandwidth_package.#":        "1",
			"accelerators.0.cross_domain_bandwidth_package.#": "1",
		}
	}

	var fakeAliCloudGaBasicAcceleratorsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"names.#":        "0",
			"accelerators.#": "0",
		}
	}

	var alicloudGaBasicAcceleratorsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ga_basic_accelerators.default",
		existMapFunc: existAliCloudGaBasicAcceleratorsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudGaBasicAcceleratorsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
		//testAccPreCheckWithTime(t, []int{1})
	}

	alicloudGaBasicAcceleratorsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, acceleratorIdConf, bandwidthBillingTypeConf, statusConf, allConf)
}

func testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand int, attrMap map[string]string) string {
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

func testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccGaBasicAccelerator-%d"
	}

	resource "alicloud_ga_basic_accelerator" "default" {
  		basic_accelerator_name = var.name
  		description            = var.name
		bandwidth_billing_type = "CDT"
  		payment_type           = "PayAsYouGo"
	}

	data "alicloud_ga_basic_accelerators" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
