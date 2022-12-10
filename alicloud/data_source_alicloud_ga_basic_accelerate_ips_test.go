package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGaBasicAccelerateIpsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaBasicAccelerateIpsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_basic_accelerate_ip.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudGaBasicAccelerateIpsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_basic_accelerate_ip.default.id}_fake"]`,
		}),
	}
	accelerateIpIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaBasicAccelerateIpsDataSourceName(rand, map[string]string{
			"accelerate_ip_id": `"${alicloud_ga_basic_accelerate_ip.default.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudGaBasicAccelerateIpsDataSourceName(rand, map[string]string{
			"accelerate_ip_id": `"${alicloud_ga_basic_accelerate_ip.default.id}_fake"`,
		}),
	}
	accelerateIpAddressConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaBasicAccelerateIpsDataSourceName(rand, map[string]string{
			"accelerate_ip_address": `"${alicloud_ga_basic_accelerate_ip.default.accelerate_ip_address}"`,
		}),
		fakeConfig: testAccCheckAlicloudGaBasicAccelerateIpsDataSourceName(rand, map[string]string{
			"accelerate_ip_address": `"${alicloud_ga_basic_accelerate_ip.default.accelerate_ip_address}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaBasicAccelerateIpsDataSourceName(rand, map[string]string{
			"accelerate_ip_id": `"${alicloud_ga_basic_accelerate_ip.default.id}"`,
			"status":           `"active"`,
		}),
		fakeConfig: testAccCheckAlicloudGaBasicAccelerateIpsDataSourceName(rand, map[string]string{
			"status": `"deleting"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaBasicAccelerateIpsDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_ga_basic_accelerate_ip.default.id}"]`,
			"accelerate_ip_id":      `"${alicloud_ga_basic_accelerate_ip.default.id}"`,
			"accelerate_ip_address": `"${alicloud_ga_basic_accelerate_ip.default.accelerate_ip_address}"`,
			"status":                `"active"`,
		}),
		fakeConfig: testAccCheckAlicloudGaBasicAccelerateIpsDataSourceName(rand, map[string]string{
			"ids":                   `["${alicloud_ga_basic_accelerate_ip.default.id}_fake"]`,
			"accelerate_ip_id":      `"${alicloud_ga_basic_accelerate_ip.default.id}_fake"`,
			"accelerate_ip_address": `"${alicloud_ga_basic_accelerate_ip.default.accelerate_ip_address}_fake"`,
			"status":                `"deleting"`,
		}),
	}
	var existAlicloudGaBasicAccelerateIpsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"ips.#":                       "1",
			"ips.0.id":                    CHECKSET,
			"ips.0.accelerate_ip_id":      CHECKSET,
			"ips.0.accelerate_ip_address": CHECKSET,
			"ips.0.accelerator_id":        CHECKSET,
			"ips.0.ip_set_id":             CHECKSET,
			"ips.0.status":                "active",
		}
	}
	var fakeAlicloudGaBasicAccelerateIpsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
			"ips.#": "0",
		}
	}
	var alicloudGaBasicAccelerateIpsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ga_basic_accelerate_ips.default",
		existMapFunc: existAlicloudGaBasicAccelerateIpsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudGaBasicAccelerateIpsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithTime(t, []int{1})
	}
	alicloudGaBasicAccelerateIpsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, accelerateIpIdConf, accelerateIpAddressConf, statusConf, allConf)
}

func testAccCheckAlicloudGaBasicAccelerateIpsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccGaBasicAccelerateIpEndpointRelation-%d"
	}

	resource "alicloud_ga_basic_accelerator" "default" {
  		duration               = 1
  		pricing_cycle          = "Month"
  		basic_accelerator_name = var.name
  		description            = var.name
  		bandwidth_billing_type = "CDT"
  		auto_pay               = true
  		auto_use_coupon        = "true"
  		auto_renew             = false
  		auto_renew_duration    = 1
	}

	resource "alicloud_ga_basic_ip_set" "default" {
  		accelerator_id       = alicloud_ga_basic_accelerator.default.id
  		accelerate_region_id = "cn-hangzhou"
  		isp_type             = "BGP"
  		bandwidth            = "5"
	}

	resource "alicloud_ga_basic_accelerate_ip" "default" {
  		accelerator_id = alicloud_ga_basic_ip_set.default.accelerator_id
  		ip_set_id      = alicloud_ga_basic_ip_set.default.id
	}

	data "alicloud_ga_basic_accelerate_ips" "default" {
  		ip_set_id = alicloud_ga_basic_accelerate_ip.default.ip_set_id
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
