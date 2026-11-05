package cen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// Skip this testcase because of the account cannot purchase non-internal products.
func SkipTestAccAlicloudCenBandwidthPackagesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 99999999)
	idConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_cen_bandwidth_package_attachment.default.instance_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_cen_bandwidth_package_attachment.default.instance_id}-fake"`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_bandwidth_package.default.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_bandwidth_package.default.name}-fake"`,
		}),
	}
	statusRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_bandwidth_package.default.name}"`,
			"status":     `"InUse"`,
		}),
		FakeConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_bandwidth_package.default.name}-fake"`,
			"status":     `"Idle"`,
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cen_bandwidth_package.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cen_bandwidth_package.default.id}-fake"]`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_cen_bandwidth_package_attachment.default.instance_id}"`,
			"name_regex":  `"${alicloud_cen_bandwidth_package.default.name}"`,
			"ids":         `["${alicloud_cen_bandwidth_package.default.id}"]`,
			"status":      `"InUse"`,
		}),
		FakeConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_cen_bandwidth_package_attachment.default.instance_id}-fake"`,
			"name_regex":  `"${alicloud_cen_bandwidth_package.default.name}"`,
			"ids":         `["${alicloud_cen_bandwidth_package.default.id}"]`,
			"status":      `"Idle"`,
		}),
	}
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
		acctest.PreCheckWithRegions(t, true, connectivity.CenNoSkipRegions)
	}
	cenBandwidthPackagesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idConf, nameRegexConf, statusRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	  default = "tf-testAcc%sCenBandwidthLimitsDataSource-%d"
	}
resource "alicloud_cen_instance" "default" {
	name = "${var.name}"
	description = "tf-testAccCenConfigDescription"
}

resource "alicloud_cen_bandwidth_package" "default" {
    bandwidth = 5
	cen_bandwidth_package_name = "${var.name}"
    geographic_region_a_id = "China"
    geographic_region_b_id = "China"
}

resource "alicloud_cen_bandwidth_package_attachment" "default" {
	instance_id = "${alicloud_cen_instance.default.id}"
	bandwidth_package_id = "${alicloud_cen_bandwidth_package.default.id}"
}

data "alicloud_cen_bandwidth_packages" "default" {
	%s

}
`, acctest.DefaultRegionToTest, rand, strings.Join(pairs, "\n  "))
	return config
}

var existCenBandwidthPackagesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                       "1",
		"packages.#":                                  "1",
		"packages.0.geographic_region_a_id":           "China",
		"packages.0.geographic_region_b_id":           "China",
		"packages.0.status":                           "InUse",
		"packages.0.bandwidth":                        "5",
		"packages.0.business_status":                  "Normal",
		"packages.0.bandwidth_package_charge_type":    "POSTPAY",
		"packages.0.description":                      "",
		"packages.0.instance_id":                      acctest.CHECKSET,
		"packages.0.name":                             fmt.Sprintf("tf-testAcc%sCenBandwidthLimitsDataSource-%d", acctest.DefaultRegionToTest, rand),
		"packages.0.id":                               acctest.CHECKSET,
		"packages.0.cen_bandwidth_package_id":         acctest.CHECKSET,
		"packages.0.cen_bandwidth_package_name":       fmt.Sprintf("tf-testAcc%sCenBandwidthLimitsDataSource-%d", acctest.DefaultRegionToTest, rand),
		"packages.0.cen_ids.#":                        "1",
		"packages.0.expired_time":                     acctest.CHECKSET,
		"packages.0.geographic_span_id":               acctest.CHECKSET,
		"packages.0.has_reservation_data":             acctest.CHECKSET,
		"packages.0.is_cross_border":                  acctest.CHECKSET,
		"packages.0.payment_type":                     "POSTPAY",
		"packages.0.reservation_active_time":          "",
		"packages.0.reservation_bandwidth":            "",
		"packages.0.reservation_internet_charge_type": "",
		"packages.0.reservation_order_type":           "",
	}
}

var fakeCenBandwidthPackagesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":      "0",
		"packages.#": "0",
	}
}

var cenBandwidthPackagesCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_cen_bandwidth_packages.default",
	ExistMapFunc: existCenBandwidthPackagesMapFunc,
	FakeMapFunc:  fakeCenBandwidthPackagesMapFunc,
}
