package ga_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudGaBasicAcceleratorsDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"ids": `["${alicloud_ga_basic_accelerator.default.id}"]`,
		}),
		FakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"ids": `["${alicloud_ga_basic_accelerator.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"name_regex": `"${alicloud_ga_basic_accelerator.default.basic_accelerator_name}"`,
		}),
		FakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"name_regex": `"${alicloud_ga_basic_accelerator.default.basic_accelerator_name}_fake"`,
		}),
	}

	acceleratorIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"accelerator_id": `"${alicloud_ga_basic_accelerator.default.id}"`,
		}),
		FakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"accelerator_id": `"${alicloud_ga_basic_accelerator.default.id}_fake"`,
		}),
	}

	bandwidthBillingTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"accelerator_id":         `"${alicloud_ga_basic_accelerator.default.id}"`,
			"bandwidth_billing_type": `"BandwidthPackage"`,
		}),
		FakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"bandwidth_billing_type": `"CDT95"`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"accelerator_id": `"${alicloud_ga_basic_accelerator.default.id}"`,
			"status":         `"active"`,
		}),
		FakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"status": `"deleting"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
			"ids":                    `["${alicloud_ga_basic_accelerator.default.id}"]`,
			"name_regex":             `"${alicloud_ga_basic_accelerator.default.basic_accelerator_name}"`,
			"accelerator_id":         `"${alicloud_ga_basic_accelerator.default.id}"`,
			"bandwidth_billing_type": `"BandwidthPackage"`,
			"status":                 `"active"`,
		}),
		FakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName0(rand, map[string]string{
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
			"accelerators.0.id":                               acctest.CHECKSET,
			"accelerators.0.basic_accelerator_id":             acctest.CHECKSET,
			"accelerators.0.basic_accelerator_name":           acctest.CHECKSET,
			"accelerators.0.basic_endpoint_group_id":          "",
			"accelerators.0.basic_ip_set_id":                  "",
			"accelerators.0.bandwidth_billing_type":           "BandwidthPackage",
			"accelerators.0.instance_charge_type":             "PREPAY",
			"accelerators.0.description":                      acctest.CHECKSET,
			"accelerators.0.region_id":                        acctest.CHECKSET,
			"accelerators.0.create_time":                      acctest.CHECKSET,
			"accelerators.0.expired_time":                     acctest.CHECKSET,
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

	var alicloudGaBasicAcceleratorsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ga_basic_accelerators.default",
		ExistMapFunc: existAliCloudGaBasicAcceleratorsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudGaBasicAcceleratorsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckPrePaidResources(t)
		//acctest.PreCheckWithTime(t, []int{1})
	}

	alicloudGaBasicAcceleratorsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, acceleratorIdConf, bandwidthBillingTypeConf, statusConf, allConf)
}

func TestAccAliCloudGaBasicAcceleratorsDataSource_basic1(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"ids": `["${alicloud_ga_basic_accelerator.default.id}"]`,
		}),
		FakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"ids": `["${alicloud_ga_basic_accelerator.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"name_regex": `"${alicloud_ga_basic_accelerator.default.basic_accelerator_name}"`,
		}),
		FakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"name_regex": `"${alicloud_ga_basic_accelerator.default.basic_accelerator_name}_fake"`,
		}),
	}

	acceleratorIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"accelerator_id": `"${alicloud_ga_basic_accelerator.default.id}"`,
		}),
		FakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"accelerator_id": `"${alicloud_ga_basic_accelerator.default.id}_fake"`,
		}),
	}

	bandwidthBillingTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"accelerator_id":         `"${alicloud_ga_basic_accelerator.default.id}"`,
			"bandwidth_billing_type": `"CDT"`,
		}),
		FakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"bandwidth_billing_type": `"CDT95"`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"accelerator_id": `"${alicloud_ga_basic_accelerator.default.id}"`,
			"status":         `"active"`,
		}),
		FakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"status": `"deleting"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
			"ids":                    `["${alicloud_ga_basic_accelerator.default.id}"]`,
			"name_regex":             `"${alicloud_ga_basic_accelerator.default.basic_accelerator_name}"`,
			"accelerator_id":         `"${alicloud_ga_basic_accelerator.default.id}"`,
			"bandwidth_billing_type": `"CDT"`,
			"status":                 `"active"`,
		}),
		FakeConfig: testAccCheckAliCloudGaBasicAcceleratorsDataSourceName1(rand, map[string]string{
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
			"accelerators.0.id":                               acctest.CHECKSET,
			"accelerators.0.basic_accelerator_id":             acctest.CHECKSET,
			"accelerators.0.basic_accelerator_name":           acctest.CHECKSET,
			"accelerators.0.basic_endpoint_group_id":          "",
			"accelerators.0.basic_ip_set_id":                  "",
			"accelerators.0.bandwidth_billing_type":           "CDT",
			"accelerators.0.instance_charge_type":             "POSTPAY",
			"accelerators.0.description":                      acctest.CHECKSET,
			"accelerators.0.region_id":                        acctest.CHECKSET,
			"accelerators.0.create_time":                      acctest.CHECKSET,
			"accelerators.0.expired_time":                     acctest.CHECKSET,
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

	var alicloudGaBasicAcceleratorsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ga_basic_accelerators.default",
		ExistMapFunc: existAliCloudGaBasicAcceleratorsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudGaBasicAcceleratorsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		//acctest.PreCheckWithTime(t, []int{1})
	}

	alicloudGaBasicAcceleratorsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, acceleratorIdConf, bandwidthBillingTypeConf, statusConf, allConf)
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
