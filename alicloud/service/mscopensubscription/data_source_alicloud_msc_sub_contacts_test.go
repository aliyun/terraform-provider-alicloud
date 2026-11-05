package mscopensubscription_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudMscSubContactsDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMscSubContactsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_msc_sub_contact.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudMscSubContactsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_msc_sub_contact.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMscSubContactsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_msc_sub_contact.default.contact_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudMscSubContactsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_msc_sub_contact.default.contact_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMscSubContactsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_msc_sub_contact.default.id}"]`,
			"name_regex": `"${alicloud_msc_sub_contact.default.contact_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudMscSubContactsDataSourceName(rand, map[string]string{
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
	var alicloudEventBridgeEventBusesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_msc_sub_contacts.default",
		ExistMapFunc: existAlicloudEventBridgeEventBusesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudEventBridgeEventBusesDataSourceNameMapFunc,
	}
	alicloudEventBridgeEventBusesCheckInfo.DataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
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
