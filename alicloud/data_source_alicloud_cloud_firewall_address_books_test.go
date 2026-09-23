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

func TestAccAliCloudCloudFirewallAddressBooksDataSource_assetType(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudFirewallAddressBooksDataSourceAssetType(rand, map[string]string{
			"ids":        `["${alicloud_cloud_firewall_address_book.default.id}"]`,
			"group_type": `"asset"`,
		}),
		fakeConfig: testAccCheckAliCloudCloudFirewallAddressBooksDataSourceAssetType(rand, map[string]string{
			"ids":        `["${alicloud_cloud_firewall_address_book.default.id}_fake"]`,
			"group_type": `"asset"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudCloudFirewallAddressBooksDataSourceAssetType(rand, map[string]string{
			"ids":        `["${alicloud_cloud_firewall_address_book.default.id}"]`,
			"name_regex": `"${alicloud_cloud_firewall_address_book.default.group_name}"`,
			"group_type": `"asset"`,
		}),
		fakeConfig: testAccCheckAliCloudCloudFirewallAddressBooksDataSourceAssetType(rand, map[string]string{
			"ids":        `["${alicloud_cloud_firewall_address_book.default.id}_fake"]`,
			"name_regex": `"${alicloud_cloud_firewall_address_book.default.group_name}_fake"`,
			"group_type": `"asset"`,
		}),
	}
	var existAliCloudCloudFirewallAddressBooksDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 "1",
			"names.#":                               "1",
			"books.#":                               "1",
			"books.0.id":                            CHECKSET,
			"books.0.group_uuid":                    CHECKSET,
			"books.0.group_name":                    CHECKSET,
			"books.0.group_type":                    "asset",
			"books.0.description":                   CHECKSET,
			"books.0.address_list_count":            CHECKSET,
			"books.0.reference_count":               CHECKSET,
			"books.0.asset_region_resource_types.#": "1",
			"books.0.asset_region_resource_types.0.asset_region_id":                        "all",
			"books.0.asset_region_resource_types.0.resource_type.#":                        "1",
			"books.0.asset_region_resource_types.0.resource_type.0.ipv4.#":                 "1",
			"books.0.asset_region_resource_types.0.resource_type.0.ipv4.0.bastion_host_ip": "true",
			"books.0.asset_region_resource_types.0.resource_type.0.ipv4.0.havip":           "true",
			"books.0.asset_region_resource_types.0.resource_type.0.ipv4.0.eip":             "false",
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
	alicloudCloudFirewallAddressBooksCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, allConf)
}

func testAccCheckAliCloudCloudFirewallAddressBooksDataSourceAssetType(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testAccAddressBookAsset-%d"
	}

	resource "alicloud_cloud_firewall_address_book" "default" {
  		group_name  = var.name
  		group_type  = "asset"
  		description = "tf-testAccAddressBookAsset"
  		asset_region_resource_types {
    		asset_region_id = "all"
    		resource_type {
      		ipv4 {
        		eip                     = false
        		ecs_eip                 = false
        		ecs_public_ip           = false
        		slb_eip                 = false
        		slb_public_ip           = false
        		nlb_eip                 = false
        		alb_eip                 = false
        		nat_eip                 = false
        		nat_public_ip           = false
        		eni_eip                 = false
        		ga_eip                  = false
        		api_gateway_eip         = false
        		ai_gateway_eip          = false
        		bastion_host_ip         = true
        		bastion_host_ingress_ip = false
        		bastion_host_egress_ip  = false
        		havip                   = true
      		}
    		}
  		}
	}

	data "alicloud_cloud_firewall_address_books" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))
	return config
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
