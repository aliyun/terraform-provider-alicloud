// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package wafopenapi_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudWafv3AddressBookDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudWafv3AddressBookSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_wafv3_address_book.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudWafv3AddressBookSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_wafv3_address_book.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudWafv3AddressBookSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_wafv3_address_book.default.address_book_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudWafv3AddressBookSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_wafv3_address_book.default.address_book_name}_fake"`,
		}),
	}

	enableDetailsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudWafv3AddressBookSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_wafv3_address_book.default.id}"]`,
			"enable_details": `"true"`,
		}),
		ExistChangMap: map[string]string{
			"books.0.address_list.#": "3",
		},
		FakeConfig: testAccCheckAlicloudWafv3AddressBookSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_wafv3_address_book.default.id}_fake"]`,
			"enable_details": `"true"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudWafv3AddressBookSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_wafv3_address_book.default.id}"]`,
			"name_regex":     `"${alicloud_wafv3_address_book.default.address_book_name}"`,
			"enable_details": `"true"`,
		}),
		ExistChangMap: map[string]string{
			"books.0.address_list.#": "3",
		},
		FakeConfig: testAccCheckAlicloudWafv3AddressBookSourceConfig(rand, map[string]string{
			"ids":            `["${alicloud_wafv3_address_book.default.id}_fake"]`,
			"name_regex":     `"${alicloud_wafv3_address_book.default.address_book_name}_fake"`,
			"enable_details": `"true"`,
		}),
	}

	Wafv3AddressBookCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, enableDetailsConf, allConf)
}

var existWafv3AddressBookMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"books.#":                   "1",
		"books.0.description":       acctest.CHECKSET,
		"books.0.address_book_name": acctest.CHECKSET,
		"books.0.address_book_id":   acctest.CHECKSET,
		"books.0.address_book_type": acctest.CHECKSET,
		"books.0.id":                acctest.CHECKSET,
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

var Wafv3AddressBookCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_wafv3_address_books.default",
	ExistMapFunc: existWafv3AddressBookMapFunc,
	FakeMapFunc:  fakeWafv3AddressBookMapFunc,
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
