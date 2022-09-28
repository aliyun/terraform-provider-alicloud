package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudFirewallAddressBooksDataSource2(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallAddressBooksDataSource(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_address_book.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallAddressBooksDataSource(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_address_book.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallAddressBooksDataSource(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_firewall_address_book.default.group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallAddressBooksDataSource(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_firewall_address_book.default.group_name}_fake"`,
		}),
	}
	groupTypePortConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallAddressBooksDataSource(rand, map[string]string{
			"ids":        `["${alicloud_cloud_firewall_address_book.default.id}"]`,
			"group_type": `"ip"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallAddressBooksDataSource(rand, map[string]string{
			"ids":        `["${alicloud_cloud_firewall_address_book.default.id}"]`,
			"group_type": `"tag"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallAddressBooksDataSource(rand, map[string]string{
			"group_type": `"ip"`,
			"name_regex": `"${alicloud_cloud_firewall_address_book.default.group_name}"`,
			"ids":        `["${alicloud_cloud_firewall_address_book.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallAddressBooksDataSource(rand, map[string]string{
			"group_type": `"tag"`,
			"name_regex": `"${alicloud_cloud_firewall_address_book.default.group_name}_fake"`,
			"ids":        `["${alicloud_cloud_firewall_address_book.default.id}_fake"]`,
		}),
	}
	var existAlicloudCloudFirewallAddressBooksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"books.#":                  "1",
			"books.0.address_list.#":   `2`,
			"books.0.auto_add_tag_ecs": `0`,
			"books.0.description":      `tf-testAccAddressBook`,
			"books.0.group_name":       fmt.Sprintf("tf-testAccAddressBook-%d", rand),
			"books.0.group_type":       `ip`,
			"books.0.id":               CHECKSET,
			"books.0.group_uuid":       CHECKSET,
			"books.0.tag_relation":     "",
			"books.0.ecs_tags.#":       `0`,
		}
	}
	var fakeAlicloudCloudFirewallAddressBooksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCloudFirewallAddressBooksCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cloud_firewall_address_books.default",
		existMapFunc: existAlicloudCloudFirewallAddressBooksDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCloudFirewallAddressBooksDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudCloudFirewallAddressBooksCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, groupTypePortConf, allConf)
}
func testAccCheckAlicloudCloudFirewallAddressBooksDataSource(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccAddressBook-%d"
}

resource "alicloud_cloud_firewall_address_book" "default" {
  description      = "tf-testAccAddressBook"
  group_name       = "${var.name}"
  group_type       = "ip"
  address_list     = ["10.21.0.0/16", "10.168.0.0/16"]
  auto_add_tag_ecs = 0
}

data "alicloud_cloud_firewall_address_books" "default" {
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
