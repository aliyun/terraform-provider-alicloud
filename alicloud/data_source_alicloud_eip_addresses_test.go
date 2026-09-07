package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEIPAddressesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"ids": `[alicloud_eip_address.default.id]`,
		}),
		fakeConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_eip_address.default.address_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_eip_address.default.address_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"ids":    `[alicloud_eip_address.default.id]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"ids":    `[alicloud_eip_address.default.id]`,
			"status": `"InUse"`,
		}),
	}
	paymentTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"ids":          `[alicloud_eip_address.default.id]`,
			"payment_type": `"PayAsYouGo"`,
		}),
		fakeConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"ids":          `[alicloud_eip_address.default.id]`,
			"payment_type": `"Subscription"`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"ids":               `[alicloud_eip_address.default.id]`,
			"resource_group_id": `"${alicloud_eip_address.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"ids":               `["fake"]`,
			"resource_group_id": `"${alicloud_eip_address.default.resource_group_id}"`,
		}),
	}
	ispConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"ids": `[alicloud_eip_address.default.id]`,
			"isp": `"BGP"`,
		}),
		fakeConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"ids": `[alicloud_eip_address.default.id]`,
			"isp": `"BGP_PRO"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"ids": `[alicloud_eip_address.default.id]`,
			"tags": `{
							Created = "tfTest"
							For 	= "tfTest 123"
					  }`,
		}),
		fakeConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"ids": `["fake"]`,
			"tags": `{
							Created = "tfTest-fake"
							For 	= "tfTest 123"
					  }`,
		}),
	}
	ipAddressConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"ids":          `[alicloud_eip_address.default.id]`,
			"ip_addresses": `["${alicloud_eip_address.default.ip_address}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"ids":          `["fake"]`,
			"ip_addresses": `["${alicloud_eip_address.default.ip_address}_fake"]`,
		}),
	}
	eipAddressConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"ids":        `[alicloud_eip_address.default.id]`,
			"ip_address": `"${alicloud_eip_address.default.ip_address}"`,
		}),
		fakeConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"ids":        `["fake"]`,
			"ip_address": `"${alicloud_eip_address.default.ip_address}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"ids":               `[alicloud_eip_address.default.id]`,
			"name_regex":        `"${alicloud_eip_address.default.address_name}"`,
			"status":            `"Available"`,
			"payment_type":      `"PayAsYouGo"`,
			"resource_group_id": `"${alicloud_eip_address.default.resource_group_id}"`,
			"ip_addresses":      `["${alicloud_eip_address.default.ip_address}"]`,
			"ip_address":        `"${alicloud_eip_address.default.ip_address}"`,
			"tags": `{
							Created = "tfTest"
							For 	= "tfTest 123"
					  }`,
			"isp": `"BGP"`,
		}),
		fakeConfig: testAccCheckAlicloudEipAddressesDataSourceName(rand, map[string]string{
			"ids":               `["fake"]`,
			"name_regex":        `"${alicloud_eip_address.default.address_name}"`,
			"status":            `"InUse"`,
			"payment_type":      `"Subscription"`,
			"resource_group_id": `"${alicloud_eip_address.default.resource_group_id}"`,
			"ip_addresses":      `["${alicloud_eip_address.default.ip_address}_fake"]`,
			"ip_address":        `"${alicloud_eip_address.default.ip_address}_fake"`,
			"tags": `{
							Created = "tfTest-fake"
							For 	= "tfTest 123"
					  }`,
			"isp": `"BGP_PRO"`,
		}),
	}
	var existAlicloudEipAddressesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"names.#":                          "1",
			"addresses.#":                      "1",
			"addresses.0.id":                   CHECKSET,
			"addresses.0.address_name":         fmt.Sprintf("tf-testAccEipAddress-%d", rand),
			"addresses.0.bandwidth":            CHECKSET,
			"addresses.0.description":          fmt.Sprintf("tf-testAccEipAddress-%d", rand),
			"addresses.0.isp":                  "BGP",
			"addresses.0.payment_type":         "PayAsYouGo",
			"addresses.0.resource_group_id":    CHECKSET,
			"addresses.0.status":               "Available",
			"addresses.0.internet_charge_type": "PayByBandwidth",
			"addresses.0.tags.%":               "2",
			"addresses.0.tags.Created":         "tfTest",
			"addresses.0.tags.For":             "tfTest 123",
			"eips.0.id":                        CHECKSET,
			"eips.0.bandwidth":                 CHECKSET,
			"eips.0.status":                    "Available",
			"eips.0.internet_charge_type":      "PayByBandwidth",
		}
	}
	var fakeAlicloudEipAddressesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"addresses.#": "0",
			"eips.#":      "0",
		}
	}
	var alicloudEipAddressesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_eip_addresses.default",
		existMapFunc: existAlicloudEipAddressesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEipAddressesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEipAddressesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, paymentTypeConf, resourceGroupIdConf, ispConf, tagsConf, ipAddressConf, eipAddressConf, allConf)
}
func testAccCheckAlicloudEipAddressesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccEipAddress-%d"
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_eip_address" "default" {
	address_name  		= var.name
	description         = var.name
	isp          		= "BGP"
	internet_charge_type = "PayByBandwidth"
	payment_type        =  "PayAsYouGo"
	tags ={
			Created = "tfTest"
			For 	= "tfTest 123"
	}
}

data "alicloud_eip_addresses" "default" {
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
