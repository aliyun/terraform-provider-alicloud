package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEipanycastAnycastEipAddressDataSource(t *testing.T) {
	rand := acctest.RandInt()

	anycastEipAddressNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"anycast_eip_address_name": `"${alicloud_eipanycast_anycast_eip_address.default.anycast_eip_address_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"anycast_eip_address_name": `"${alicloud_eipanycast_anycast_eip_address.default.anycast_eip_address_name}_fake"`,
		}),
	}
	serviceLocationConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_eipanycast_anycast_eip_address.default.id}"]`,
			"service_location": `"${alicloud_eipanycast_anycast_eip_address.default.service_location}"`,
		}),
		fakeConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_eipanycast_anycast_eip_address.default.id}"]`,
			"service_location": `"${alicloud_eipanycast_anycast_eip_address.default.service_location}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_eipanycast_anycast_eip_address.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_eipanycast_anycast_eip_address.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"name_regex": `"^${alicloud_eipanycast_anycast_eip_address.default.anycast_eip_address_name}$"`,
		}),
		fakeConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_eipanycast_anycast_eip_address.default.anycast_eip_address_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_eipanycast_anycast_eip_address.default.id}"]`,
			"status": `"${alicloud_eipanycast_anycast_eip_address.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_eipanycast_anycast_eip_address.default.id}"]`,
			"status": `"Unassociating"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_eipanycast_anycast_eip_address.default.id}"]`,
			"name_regex": `"^${alicloud_eipanycast_anycast_eip_address.default.anycast_eip_address_name}$"`,
			"status":     `"${alicloud_eipanycast_anycast_eip_address.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_eipanycast_anycast_eip_address.default.id}_fake"]`,
			"name_regex": `"${alicloud_eipanycast_anycast_eip_address.default.anycast_eip_address_name}_fake"`,
			"status":     `"Unassociating"`,
		}),
	}
	var existAlicloudEipanycastAnycastEipAddressDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"names.#":                              "1",
			"addresses.#":                          "1",
			"addresses.0.ali_uid":                  CHECKSET,
			"addresses.0.anycast_eip_address_name": fmt.Sprintf("tf-testAccAlicloudEipanycastAnycastEipAddress%d", rand),
			"addresses.0.anycast_eip_bind_info_list.#": "0",
			"addresses.0.id":                   CHECKSET,
			"addresses.0.anycast_id":           CHECKSET,
			"addresses.0.bandwidth":            CHECKSET,
			"addresses.0.bid":                  CHECKSET,
			"addresses.0.business_status":      CHECKSET,
			"addresses.0.description":          "tf-testAccAlicloudEipanycastAnycastEipAddress",
			"addresses.0.internet_charge_type": "PayByTraffic",
			"addresses.0.ip_address":           CHECKSET,
			"addresses.0.payment_type":         "PayAsYouGo",
			"addresses.0.service_location":     "international",
			"addresses.0.status":               "Allocated",
		}
	}
	var fakeAlicloudEipanycastAnycastEipAddressDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"addresses.#": "0",
		}
	}
	var alicloudEipanycastAnycastEipAddressCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_eipanycast_anycast_eip_addresses.default",
		existMapFunc: existAlicloudEipanycastAnycastEipAddressDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEipanycastAnycastEipAddressDataSourceNameMapFunc,
	}

	alicloudEipanycastAnycastEipAddressCheckInfo.dataSourceTestCheck(t, rand, anycastEipAddressNameConf, serviceLocationConf, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
resource "alicloud_eipanycast_anycast_eip_address" "default" {
  service_location         = "international"
  anycast_eip_address_name = "tf-testAccAlicloudEipanycastAnycastEipAddress%d"
  payment_type             = "PayAsYouGo"
  internet_charge_type     = "PayByTraffic"
  description              = "tf-testAccAlicloudEipanycastAnycastEipAddress"
}

data "alicloud_eipanycast_anycast_eip_addresses" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
