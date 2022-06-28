package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECSPrefixListsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc")
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsPrefixListDataSourceName(name, map[string]string{
			"enable_details": "true",
			"ids":            `["${alicloud_ecs_prefix_list.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsPrefixListDataSourceName(name, map[string]string{
			"ids": `["${alicloud_ecs_prefix_list.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsPrefixListDataSourceName(name, map[string]string{
			"enable_details": "true",
			"name_regex":     `"${alicloud_ecs_prefix_list.default.prefix_list_name}"`,
			"ids":            `["${alicloud_ecs_prefix_list.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsPrefixListDataSourceName(name, map[string]string{
			"name_regex": `"${alicloud_ecs_prefix_list.default.prefix_list_name}_fake"`,
			"ids":        `["${alicloud_ecs_prefix_list.default.id}"]`,
		}),
	}
	addressFamilyRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsPrefixListDataSourceName(name, map[string]string{
			"enable_details": "true",
			"address_family": `"${alicloud_ecs_prefix_list.default.address_family}"`,
			"ids":            `["${alicloud_ecs_prefix_list.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsPrefixListDataSourceName(name, map[string]string{
			"address_family": `"IPv6"`,
			"ids":            `["${alicloud_ecs_prefix_list.default.id}"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsPrefixListDataSourceName(name, map[string]string{
			"enable_details": "true",
			"ids":            `["${alicloud_ecs_prefix_list.default.id}"]`,
			"name_regex":     `"${alicloud_ecs_prefix_list.default.prefix_list_name}"`,
			"address_family": `"${alicloud_ecs_prefix_list.default.address_family}"`,
		}),
		fakeConfig: testAccCheckAlicloudEcsPrefixListDataSourceName(name, map[string]string{
			"ids":            `["${alicloud_ecs_prefix_list.default.id}_fake"]`,
			"name_regex":     `"${alicloud_ecs_prefix_list.default.prefix_list_name}_fake"`,
			"address_family": `"IPv6"`,
		}),
	}
	var existAlicloudEcsSnapshotsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"names.#":                   "1",
			"lists.#":                   "1",
			"lists.0.address_family":    `IPv4`,
			"lists.0.description":       name,
			"lists.0.association_count": CHECKSET,
			"lists.0.create_time":       CHECKSET,
			"lists.0.entry.#":           "1",
			"lists.0.max_entries":       "2",
			"lists.0.prefix_list_id":    CHECKSET,
			"lists.0.prefix_list_name":  CHECKSET,
		}
	}
	var fakeAlicloudEcsSnapshotsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEcsSnapshotsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ecs_prefix_lists.default",
		existMapFunc: existAlicloudEcsSnapshotsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEcsSnapshotsDataSourceNameMapFunc,
	}
	alicloudEcsSnapshotsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, addressFamilyRegexConf, allConf)
}
func testAccCheckAlicloudEcsPrefixListDataSourceName(name string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
	default = "%s"
}
resource "alicloud_ecs_prefix_list" "default"{
	address_family = "IPv4"
	max_entries = 2
	prefix_list_name = var.name
	description = var.name
	entry {
		cidr = "192.168.0.0/24"
		description = "description"
	}
}

data "alicloud_ecs_prefix_lists" "default" {	
	%s
}
`, name, strings.Join(pairs, " \n "))
	return config
}
