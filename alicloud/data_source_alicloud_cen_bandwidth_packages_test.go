package alicloud

import (
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"fmt"
)

func TestAccAlicloudCenBandwidthPackagesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 99999999)
	idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_cen_bandwidth_package_attachment.default.instance_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_cen_bandwidth_package_attachment.default.instance_id}-fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_bandwidth_package.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_bandwidth_package.default.name}-fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cen_bandwidth_package.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cen_bandwidth_package.default.id}-fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_cen_bandwidth_package_attachment.default.instance_id}"`,
			"name_regex":  `"${alicloud_cen_bandwidth_package.default.name}"`,
			"ids":         `["${alicloud_cen_bandwidth_package.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"instance_id": `"${alicloud_cen_bandwidth_package_attachment.default.instance_id}-fake"`,
			"name_regex":  `"${alicloud_cen_bandwidth_package.default.name}"`,
			"ids":         `["${alicloud_cen_bandwidth_package.default.id}"]`,
		}),
	}
	preCheck := func() {
		testAccPreCheckWithAccountSiteType(t, DomesticSite)
		testAccPreCheckWithRegions(t, true, connectivity.CenNoSkipRegions)
	}
	cenBandwidthPackagesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idConf, nameRegexConf, idsConf, allConf)
}

func TestAccAlicloudCenBandwidthPackagesDataSource_multi(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 99999999)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig_multi(rand, map[string]string{
			"ids": `"${alicloud_cen_bandwidth_package.default.*.id}"`,
		}),
	}
	preCheck := func() {
		testAccPreCheckWithAccountSiteType(t, DomesticSite)
	}
	cenBandwidthPackagesCheckInfo_multi.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
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
	name = "${var.name}"
    geographic_region_ids = [
		"China",
		"Asia-Pacific"]
}

resource "alicloud_cen_bandwidth_package_attachment" "default" {
	instance_id = "${alicloud_cen_instance.default.id}"
	bandwidth_package_id = "${alicloud_cen_bandwidth_package.default.id}"
}

data "alicloud_cen_bandwidth_packages" "default" {
	%s

}
`, defaultRegionToTest, rand, strings.Join(pairs, "\n  "))
	return config
}

func testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig_multi(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
resource "alicloud_cen_bandwidth_package" "default" {
	name = "tf-testAcc%sCenBandwidthLimitsDataSource-%d"
    bandwidth = 5
    geographic_region_ids = [
		"China",
		"China"]
	count = 6
}

data "alicloud_cen_bandwidth_packages" "default" {
	%s
}
`, defaultRegionToTest, rand, strings.Join(pairs, "\n  "))
	return config
}

var existCenBandwidthPackagesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                    "1",
		"packages.#":                               "1",
		"packages.0.geographic_region_a_id":        "Asia-Pacific",
		"packages.0.geographic_region_b_id":        "China",
		"packages.0.status":                        "InUse",
		"packages.0.bandwidth":                     "5",
		"packages.0.business_status":               "Normal",
		"packages.0.bandwidth_package_charge_type": "POSTPAY",
		"packages.0.description":                   "",
		"packages.0.name":                          CHECKSET,
		"packages.0.creation_time":                 CHECKSET,
		"packages.0.id":                            CHECKSET,
		"packages.0.instance_id":                   CHECKSET,
	}
}
var existCenBandwidthPackagesMapFunc_multi = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                    "6",
		"packages.#":                               "6",
		"packages.0.geographic_region_a_id":        "China",
		"packages.0.geographic_region_b_id":        "China",
		"packages.0.status":                        CHECKSET,
		"packages.0.bandwidth":                     "5",
		"packages.0.business_status":               "Normal",
		"packages.0.bandwidth_package_charge_type": "POSTPAY",
		"packages.0.description":                   "",
		"packages.0.name":                          CHECKSET,
		"packages.0.creation_time":                 CHECKSET,
		"packages.0.id":                            CHECKSET,
	}
}

var fakeCenBandwidthPackagesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":      "0",
		"packages.#": "0",
	}
}

var cenBandwidthPackagesCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cen_bandwidth_packages.default",
	existMapFunc: existCenBandwidthPackagesMapFunc,
	fakeMapFunc:  fakeCenBandwidthPackagesMapFunc,
}
var cenBandwidthPackagesCheckInfo_multi = dataSourceAttr{
	resourceId:   "data.alicloud_cen_bandwidth_packages.default",
	existMapFunc: existCenBandwidthPackagesMapFunc_multi,
	fakeMapFunc:  fakeCenBandwidthPackagesMapFunc,
}
