// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlicloudWafv3AddressBookDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafv3AddressBookSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_wafv3_address_book.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudWafv3AddressBookSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_wafv3_address_book.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafv3AddressBookSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_wafv3_address_book.default.address_book_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudWafv3AddressBookSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_wafv3_address_book.default.address_book_name}_fake"`,
		}),
	}

	enableDetailsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafv3AddressBookSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_wafv3_address_book.default.id}"]`,
			"enable_details": `"true"`,
		}),
		existChangMap: map[string]string{
			"books.0.address_list.#": "3",
		},
		fakeConfig: testAccCheckAlicloudWafv3AddressBookSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_wafv3_address_book.default.id}_fake"]`,
			"enable_details": `"true"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafv3AddressBookSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_wafv3_address_book.default.id}"]`,
			"name_regex":     `"${alicloud_wafv3_address_book.default.address_book_name}"`,
			"enable_details": `"true"`,
		}),
		existChangMap: map[string]string{
			"books.0.address_list.#": "3",
		},
		fakeConfig: testAccCheckAlicloudWafv3AddressBookSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_wafv3_address_book.default.id}_fake"]`,
			"name_regex":     `"${alicloud_wafv3_address_book.default.address_book_name}_fake"`,
			"enable_details": `"true"`,
		}),
	}

	Wafv3AddressBookCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, enableDetailsConf, allConf)
}

var existWafv3AddressBookMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"books.#":                   "1",
		"books.0.description":       CHECKSET,
		"books.0.address_book_name": CHECKSET,
		"books.0.address_book_id":   CHECKSET,
		"books.0.address_book_type": CHECKSET,
		"books.0.id":                CHECKSET,
		"ids.#":                     "1",
		"names.#":                   "1",
	}
}

var fakeWafv3AddressBookMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"books.#": "0",
		"ids.#":   "0",
		"names.#": "0",
	}
}

var Wafv3AddressBookCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_wafv3_address_books.default",
	existMapFunc: existWafv3AddressBookMapFunc,
	fakeMapFunc:  fakeWafv3AddressBookMapFunc,
}

func testAccCheckAlicloudWafv3AddressBookSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccWafv3AddressBook%d"
}

data "alicloud_wafv3_instances" "default" {
}

resource "alicloud_wafv3_address_book" "default" {
  description       = "test"
  instance_id       = data.alicloud_wafv3_instances.default.ids.0
  address_book_name = var.name
  address_list      = ["100.100.100.100/32", "101.101.101.101/32", "102.102.102.102/32"]
  address_book_type = "ip"
}

data "alicloud_wafv3_address_books" "default" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  %s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
