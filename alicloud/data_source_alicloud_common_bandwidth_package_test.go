package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCommonBandwidthPackagesDataSourceBasic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"name_regex": `"${alicloud_common_bandwidth_package.default.bandwidth_package_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"name_regex": `"${alicloud_common_bandwidth_package.default.bandwidth_package_name}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids":    `[ "${alicloud_common_bandwidth_package.default.id}" ]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids":    `[ "${alicloud_common_bandwidth_package.default.id}_fake" ]`,
			"status": `"Pending"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids": `[ "${alicloud_common_bandwidth_package.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids": `[ "${alicloud_common_bandwidth_package.default.id}_fake" ]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids":        `[ "${alicloud_common_bandwidth_package.default.id}" ]`,
			"name_regex": `"${alicloud_common_bandwidth_package.default.bandwidth_package_name}"`,
			"status":     `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids":        `[ "${alicloud_common_bandwidth_package.default.id}_fake" ]`,
			"name_regex": `"${alicloud_common_bandwidth_package.default.bandwidth_package_name}_fake"`,
			"status":     `"Pending"`,
		}),
	}
	commonBandwidthPackagesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, statusConf, idsConf, allConf)
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
		"packages.0.id":                               CHECKSET,
		"packages.0.isp":                              CHECKSET,
		"packages.0.status":                           CHECKSET,
		"packages.0.business_status":                  CHECKSET,
		"packages.0.bandwidth":                        "2",
		"packages.0.name":                             fmt.Sprintf("tf-testAccCommonBandwidthPackageDataSource%d", rand),
		"packages.0.bandwidth_package_name":           fmt.Sprintf("tf-testAccCommonBandwidthPackageDataSource%d", rand),
		"packages.0.bandwidth_package_id":             CHECKSET,
		"packages.0.description":                      fmt.Sprintf("tf-testAccCommonBandwidthPackageDataSource%d_description", rand),
		"packages.0.service_managed":                  CHECKSET,
		"packages.0.resource_group_id":                CHECKSET,
		"packages.0.reservation_order_type":           "",
		"packages.0.reservation_internet_charge_type": "",
		"packages.0.reservation_bandwidth":            "",
		"packages.0.reservation_active_time":          "",
		"packages.0.ratio":                            CHECKSET,
		"packages.0.public_ip_addresses.#":            CHECKSET,
		"packages.0.payment_type":                     CHECKSET,
		"packages.0.internet_charge_type":             CHECKSET,
		"packages.0.has_reservation_data":             CHECKSET,
		"packages.0.expired_time":                     "",
		"packages.0.deletion_protection":              CHECKSET,
	}
}

var fakeCommonBandwidthPackagesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":      "0",
		"names.#":    "0",
		"packages.#": "0",
	}
}

var commonBandwidthPackagesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_common_bandwidth_packages.default",
	existMapFunc: existsCommonBandwidthPackagesMapFunc,
	fakeMapFunc:  fakeCommonBandwidthPackagesMapFunc,
}
