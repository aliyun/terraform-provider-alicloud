package cbwp_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCommonBandwidthPackagesDataSourceBasic(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"name_regex": `"${alicloud_common_bandwidth_package.default.bandwidth_package_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"name_regex": `"${alicloud_common_bandwidth_package.default.bandwidth_package_name}_fake"`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids":    `[ "${alicloud_common_bandwidth_package.default.id}" ]`,
			"status": `"Available"`,
		}),
		FakeConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids":    `[ "${alicloud_common_bandwidth_package.default.id}_fake" ]`,
			"status": `"Pending"`,
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids": `[ "${alicloud_common_bandwidth_package.default.id}" ]`,
		}),
		FakeConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids": `[ "${alicloud_common_bandwidth_package.default.id}_fake" ]`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids":        `[ "${alicloud_common_bandwidth_package.default.id}" ]`,
			"name_regex": `"${alicloud_common_bandwidth_package.default.bandwidth_package_name}"`,
			"status":     `"Available"`,
		}),
		FakeConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids":        `[ "${alicloud_common_bandwidth_package.default.id}_fake" ]`,
			"name_regex": `"${alicloud_common_bandwidth_package.default.bandwidth_package_name}_fake"`,
			"status":     `"Pending"`,
		}),
	}
	commonBandwidthPackagesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	}, nameRegexConf, statusConf, idsConf, allConf)
}

func testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccCommonBandwidthPackageDataSource%d"
}

data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}

resource "alicloud_common_bandwidth_package" "default" {
  bandwidth = "2"
  internet_charge_type = "PayByBandwidth"
  name = "${var.name}"
  description = "${var.name}_description"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}

data "alicloud_common_bandwidth_packages" "default"  {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existsCommonBandwidthPackagesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                       "1",
		"names.#":                                     "1",
		"packages.#":                                  "1",
		"packages.0.id":                               acctest.CHECKSET,
		"packages.0.isp":                              acctest.CHECKSET,
		"packages.0.status":                           acctest.CHECKSET,
		"packages.0.business_status":                  acctest.CHECKSET,
		"packages.0.bandwidth":                        "2",
		"packages.0.name":                             fmt.Sprintf("tf-testAccCommonBandwidthPackageDataSource%d", rand),
		"packages.0.bandwidth_package_name":           fmt.Sprintf("tf-testAccCommonBandwidthPackageDataSource%d", rand),
		"packages.0.bandwidth_package_id":             acctest.CHECKSET,
		"packages.0.description":                      fmt.Sprintf("tf-testAccCommonBandwidthPackageDataSource%d_description", rand),
		"packages.0.service_managed":                  acctest.CHECKSET,
		"packages.0.resource_group_id":                acctest.CHECKSET,
		"packages.0.reservation_order_type":           "",
		"packages.0.reservation_internet_charge_type": "",
		"packages.0.reservation_bandwidth":            "",
		"packages.0.reservation_active_time":          "",
		"packages.0.ratio":                            acctest.CHECKSET,
		"packages.0.public_ip_addresses.#":            acctest.CHECKSET,
		"packages.0.payment_type":                     acctest.CHECKSET,
		"packages.0.internet_charge_type":             acctest.CHECKSET,
		"packages.0.has_reservation_data":             acctest.CHECKSET,
		"packages.0.expired_time":                     "",
		"packages.0.deletion_protection":              acctest.CHECKSET,
	}
}

var fakeCommonBandwidthPackagesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":      "0",
		"names.#":    "0",
		"packages.#": "0",
	}
}

var commonBandwidthPackagesCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_common_bandwidth_packages.default",
	ExistMapFunc: existsCommonBandwidthPackagesMapFunc,
	FakeMapFunc:  fakeCommonBandwidthPackagesMapFunc,
}
