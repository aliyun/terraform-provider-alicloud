package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudFirewallAddressBooksDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallAddressBooksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_address_book.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallAddressBooksDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_address_book.default.id}_fake"]`,
		}),
	}

	groupTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallAddressBooksDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloud_firewall_address_book.default.id}"]`,
			"group_type": `"${alicloud_cloud_firewall_address_book.default.group_type}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallAddressBooksDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloud_firewall_address_book.default.id}"]`,
			"group_type": `"domain"`,
		}),
	}

	var existAlicloudCloudFirewallAddressBooksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"books.#":                  "1",
			"books.0.auto_add_tag_ecs": `0`,
			"books.0.description":      `tf-testAcc-jOgZg`,
			"books.0.group_name":       fmt.Sprintf("tf-testAccAddressBook-%d", rand),
			"books.0.group_type":       `ip`,
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
	alicloudCloudFirewallAddressBooksCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, groupTypeConf)
}
func testAccCheckAlicloudCloudFirewallAddressBooksDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccAddressBook-%d"
}

resource "alicloud_cloud_firewall_address_book" "default" {
description = "tf-testAcc-jOgZg"
group_name = "${var.name}"
group_type =       "ip"
address_list =     ["10.21.0.0/16", "10.168.0.0/16"]
auto_add_tag_ecs = 0
}

data "alicloud_cloud_firewall_address_books" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

func TestAccAlicloudCloudFirewallAddressBooksDataSource2(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallAddressBooksDataSourceName2(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_address_book.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallAddressBooksDataSourceName2(rand, map[string]string{
			"ids": `["${alicloud_cloud_firewall_address_book.default.id}_fake"]`,
		}),
	}

	containPortConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallAddressBooksDataSourceName2(rand, map[string]string{
			"ids":          `["${alicloud_cloud_firewall_address_book.default.id}"]`,
			"contain_port": "80",
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallAddressBooksDataSourceName2(rand, map[string]string{
			"ids":          `["${alicloud_cloud_firewall_address_book.default.id}"]`,
			"contain_port": `"8080"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallAddressBooksDataSourceName2(rand, map[string]string{
			"contain_port": "80",
			"group_type":   `"${alicloud_cloud_firewall_address_book.default.group_type}"`,
			"ids":          `["${alicloud_cloud_firewall_address_book.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudFirewallAddressBooksDataSourceName2(rand, map[string]string{
			"contain_port": `"8080"`,
			"group_type":   `"domain"`,
			"ids":          `["${alicloud_cloud_firewall_address_book.default.id}_fake"]`,
		}),
	}
	var existAlicloudCloudFirewallAddressBooksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"books.#":                  "1",
			"books.0.auto_add_tag_ecs": `0`,
			"books.0.description":      `tf-testAcc-jOgZg`,
			"books.0.group_name":       fmt.Sprintf("tf-testAccAddressBook-%d", rand),
			"books.0.group_type":       `port`,
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
	alicloudCloudFirewallAddressBooksCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, containPortConf, allConf)
}
func testAccCheckAlicloudCloudFirewallAddressBooksDataSourceName2(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "tf-testAccAddressBook-%d"
}

resource "alicloud_cloud_firewall_address_book" "default" {
description = "tf-testAcc-jOgZg"
group_name = "${var.name}"
group_type =       "port"
address_list =     ["22", "80"]
auto_add_tag_ecs = 0
}

data "alicloud_cloud_firewall_address_books" "default" {
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
