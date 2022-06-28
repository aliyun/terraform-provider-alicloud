package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMscSubContactsDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMscSubContactsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_msc_sub_contact.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMscSubContactsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_msc_sub_contact.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMscSubContactsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_msc_sub_contact.default.contact_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMscSubContactsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_msc_sub_contact.default.contact_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMscSubContactsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_msc_sub_contact.default.id}"]`,
			"name_regex": `"${alicloud_msc_sub_contact.default.contact_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMscSubContactsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_msc_sub_contact.default.id}_fake"]`,
			"name_regex": `"${alicloud_msc_sub_contact.default.contact_name}fake"`,
		}),
	}
	var existAlicloudEventBridgeEventBusesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"names.#":                 "1",
			"contacts.#":              "1",
			"contacts.0.contact_name": "testtfac",
			"contacts.0.email":        "123@163.com",
			"contacts.0.position":     "CEO",
			"contacts.0.mobile":       "12345257908",
		}
	}
	var fakeAlicloudEventBridgeEventBusesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"contacts.#": "0",
		}
	}
	var alicloudEventBridgeEventBusesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_msc_sub_contacts.default",
		existMapFunc: existAlicloudEventBridgeEventBusesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEventBridgeEventBusesDataSourceNameMapFunc,
	}
	alicloudEventBridgeEventBusesCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudMscSubContactsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "testtfac"
}
resource "alicloud_msc_sub_contact" "default" {
	contact_name = var.name
	position = "CEO"
    email =  "123@163.com"
    mobile = "12345257908"
}
data "alicloud_msc_sub_contacts" "default" {	
	%s
}
`, strings.Join(pairs, " \n "))
	return config
}
