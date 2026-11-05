package eipanycast_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEipanycastAnycastEipAddressDataSource(t *testing.T) {
	rand := tfacctest.RandInt()

	anycastEipAddressNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"anycast_eip_address_name": `"${alicloud_eipanycast_anycast_eip_address.default.anycast_eip_address_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"anycast_eip_address_name": `"${alicloud_eipanycast_anycast_eip_address.default.anycast_eip_address_name}_fake"`,
		}),
	}
	serviceLocationConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_eipanycast_anycast_eip_address.default.id}"]`,
			"service_location": `"${alicloud_eipanycast_anycast_eip_address.default.service_location}"`,
		}),
		FakeConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_eipanycast_anycast_eip_address.default.id}"]`,
			"service_location": `"${alicloud_eipanycast_anycast_eip_address.default.service_location}_fake"`,
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_eipanycast_anycast_eip_address.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_eipanycast_anycast_eip_address.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"name_regex": `"^${alicloud_eipanycast_anycast_eip_address.default.anycast_eip_address_name}$"`,
		}),
		FakeConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_eipanycast_anycast_eip_address.default.anycast_eip_address_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_eipanycast_anycast_eip_address.default.id}"]`,
			"status": `"${alicloud_eipanycast_anycast_eip_address.default.status}"`,
		}),
		FakeConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_eipanycast_anycast_eip_address.default.id}"]`,
			"status":            `"Unassociating"`,
			"ip_address":        `"1.1.1.1"`,
			"bind_instance_ids": `["fake-instance-id-1"]`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_eipanycast_anycast_eip_address.default.id}"]`,
			"name_regex":           `"^${alicloud_eipanycast_anycast_eip_address.default.anycast_eip_address_name}$"`,
			"status":               `"${alicloud_eipanycast_anycast_eip_address.default.status}"`,
			"internet_charge_type": `"${alicloud_eipanycast_anycast_eip_address.default.internet_charge_type}"`,
			"payment_type":         `"${alicloud_eipanycast_anycast_eip_address.default.payment_type}"`,
			"business_status":      `"Normal"`,
			"output_file":          `"/tmp/tf-test-eipanycast-ds-exist.json"`,
		}),
		FakeConfig: testAccCheckAlicloudEipanycastAnycastEipAddressDataSourceName(rand, map[string]string{
			"ids":                  `["${alicloud_eipanycast_anycast_eip_address.default.id}_fake"]`,
			"name_regex":           `"${alicloud_eipanycast_anycast_eip_address.default.anycast_eip_address_name}_fake"`,
			"status":               `"Unassociating"`,
			"internet_charge_type": `"PayByBandwidth"`,
			"payment_type":         `"PayAsYouGo"`,
			"business_status":      `"FinancialLocked"`,
			"ip_address":           `"2.2.2.2"`,
			"bind_instance_ids":    `["fake-instance-id-2"]`,
			"output_file":          `"/tmp/tf-test-eipanycast-ds-fake.json"`,
		}),
	}
	var existAlicloudEipanycastAnycastEipAddressDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"names.#":                              "1",
			"addresses.#":                          "1",
			"addresses.0.ali_uid":                  acctest.CHECKSET,
			"addresses.0.anycast_eip_address_name": fmt.Sprintf("tf-testAccAlicloudEipanycastAnycastEipAddress%d", rand),
			"addresses.0.anycast_eip_bind_info_list.#": "0",
			"addresses.0.id":                   acctest.CHECKSET,
			"addresses.0.anycast_id":           acctest.CHECKSET,
			"addresses.0.bandwidth":            acctest.CHECKSET,
			"addresses.0.bid":                  acctest.CHECKSET,
			"addresses.0.business_status":      acctest.CHECKSET,
			"addresses.0.description":          "tf-testAccAlicloudEipanycastAnycastEipAddress",
			"addresses.0.internet_charge_type": "PayByTraffic",
			"addresses.0.ip_address":           acctest.CHECKSET,
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
	var alicloudEipanycastAnycastEipAddressCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_eipanycast_anycast_eip_addresses.default",
		ExistMapFunc: existAlicloudEipanycastAnycastEipAddressDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudEipanycastAnycastEipAddressDataSourceNameMapFunc,
	}

	alicloudEipanycastAnycastEipAddressCheckInfo.DataSourceTestCheck(t, rand, anycastEipAddressNameConf, serviceLocationConf, idsConf, nameRegexConf, statusConf, allConf)
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
