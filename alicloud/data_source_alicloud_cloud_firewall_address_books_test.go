package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudCloudFirewallAddressBooksDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudFirewallAddressBooksDataSource(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_address_book.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudCloudFirewallAddressBooksDataSource(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_address_book.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudFirewallAddressBooksDataSource(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_firewall_address_book.default.group_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudCloudFirewallAddressBooksDataSource(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_firewall_address_book.default.group_name}_fake"`,
		}),
	}
	groupTypePortConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudFirewallAddressBooksDataSource(rand, map[string]string{
			"ids":        `["${alicloud_cloud_firewall_address_book.default.id}"]`,
			"group_type": `"ip"`,
		}),
		fakeConfig: testAccCheckAliCloudCloudFirewallAddressBooksDataSource(rand, map[string]string{
			"group_type": `"tag"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudFirewallAddressBooksDataSource(rand, map[string]string{
			"ids":        `["${alicloud_cloud_firewall_address_book.default.id}"]`,
			"name_regex": `"${alicloud_cloud_firewall_address_book.default.group_name}"`,
			"group_type": `"ip"`,
		}),
		fakeConfig: testAccCheckAliCloudCloudFirewallAddressBooksDataSource(rand, map[string]string{
			"ids":        `["${alicloud_cloud_firewall_address_book.default.id}_fake"]`,
			"name_regex": `"${alicloud_cloud_firewall_address_book.default.group_name}_fake"`,
			"group_type": `"tag"`,
		}),
	}
	var existAliCloudCloudFirewallAddressBooksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"names.#":                  "1",
			"books.#":                  "1",
			"books.0.id":               CHECKSET,
			"books.0.group_uuid":       CHECKSET,
			"books.0.group_name":       CHECKSET,
			"books.0.group_type":       "ip",
			"books.0.description":      CHECKSET,
			"books.0.auto_add_tag_ecs": "0",
			"books.0.tag_relation":     "",
			"books.0.address_list.#":   "2",
			"books.0.ecs_tags.#":       "0",
		}
	}
	var fakeAliCloudCloudFirewallAddressBooksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"books.#": "0",
		}
	}
	var alicloudCloudFirewallAddressBooksCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cloud_firewall_address_books.default",
		existMapFunc: existAliCloudCloudFirewallAddressBooksDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudCloudFirewallAddressBooksDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudCloudFirewallAddressBooksCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, groupTypePortConf, allConf)
}

func testAccCheckAliCloudCloudFirewallAddressBooksDataSource(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccAddressBook-%d"
	}

	resource "alicloud_cloud_firewall_address_book" "default" {
  		group_name       = var.name
  		group_type       = "ip"
  		description      = "tf-testAccAddressBook"
  		auto_add_tag_ecs = 0
  		address_list     = ["10.21.0.0/16", "10.168.0.0/16"]
	}

	data "alicloud_cloud_firewall_address_books" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
}
