package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVpcPrefixListsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPrefixListsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_prefix_list.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPrefixListsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_vpc_prefix_list.default.id}_fake"]`,
		}),
	}
	prefixListNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPrefixListsDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_vpc_prefix_list.default.id}"]`,
			"prefix_list_name": `"${alicloud_vpc_prefix_list.default.prefix_list_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPrefixListsDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_vpc_prefix_list.default.id}"]`,
			"prefix_list_name": `"${alicloud_vpc_prefix_list.default.prefix_list_name}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPrefixListsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vpc_prefix_list.default.prefix_list_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPrefixListsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_vpc_prefix_list.default.prefix_list_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudVpcPrefixListsDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_vpc_prefix_list.default.id}"]`,
			"name_regex":       `"${alicloud_vpc_prefix_list.default.prefix_list_name}"`,
			"prefix_list_name": `"${alicloud_vpc_prefix_list.default.prefix_list_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudVpcPrefixListsDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_vpc_prefix_list.default.id}_fake"]`,
			"name_regex":       `"${alicloud_vpc_prefix_list.default.prefix_list_name}_fake"`,
			"prefix_list_name": `"${alicloud_vpc_prefix_list.default.prefix_list_name}_fake"`,
		}),
	}
	var existAlicloudVpcPrefixListsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"lists.#":                         "1",
			"lists.0.entrys.#":                "1",
			"lists.0.entrys.0.cidr":           "192.168.0.0/16",
			"lists.0.entrys.0.description":    "description",
			"lists.0.ip_version":              "IPV4",
			"lists.0.max_entries":             "50",
			"lists.0.prefix_list_name":        CHECKSET,
			"lists.0.prefix_list_description": "description",
			"lists.0.create_time":             CHECKSET,
			"lists.0.id":                      CHECKSET,
			"lists.0.prefix_list_id":          CHECKSET,
			"lists.0.share_type":              "",
		}
	}
	var fakeAlicloudVpcPrefixListsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudVpcPrefixListsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_vpc_prefix_lists.default",
		existMapFunc: existAlicloudVpcPrefixListsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudVpcPrefixListsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudVpcPrefixListsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, prefixListNameConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudVpcPrefixListsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccPrefixList-%d"
}

resource "alicloud_vpc_prefix_list" "default" {
	entrys {
		cidr =  "192.168.0.0/16"
		description = "description"
	}
	ip_version = "IPV4"
	max_entries = 50
	prefix_list_name = var.name
	prefix_list_description = "description"
}

data "alicloud_vpc_prefix_lists" "default" {	
	enable_details = true
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
