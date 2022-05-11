package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCddcDedicatedHostGroupsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 200)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCddcDedicatedHostGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cddc_dedicated_host_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCddcDedicatedHostGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cddc_dedicated_host_group.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCddcDedicatedHostGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cddc_dedicated_host_group.default.id}"]`,
			"name_regex": `"${alicloud_cddc_dedicated_host_group.default.dedicated_host_group_desc}"`,
		}),
		fakeConfig: testAccCheckAlicloudCddcDedicatedHostGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cddc_dedicated_host_group.default.id}_fake"]`,
			"name_regex": `"${alicloud_cddc_dedicated_host_group.default.dedicated_host_group_desc}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCddcDedicatedHostGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cddc_dedicated_host_group.default.id}"]`,
			"engine":     `"MySQL"`,
			"name_regex": `"${alicloud_cddc_dedicated_host_group.default.dedicated_host_group_desc}"`,
		}),
		fakeConfig: testAccCheckAlicloudCddcDedicatedHostGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cddc_dedicated_host_group.default.id}_fake"]`,
			"engine":     `"MySQL"`,
			"name_regex": `"${alicloud_cddc_dedicated_host_group.default.dedicated_host_group_desc}_fake"`,
		}),
	}
	var existAlicloudCddcDedicatedHostGroupsNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                            "1",
			"names.0":                            fmt.Sprintf("tf-testAccName-%d", rand),
			"ids.#":                              "1",
			"groups.#":                           "1",
			"groups.0.engine":                    "MySQL",
			"groups.0.dedicated_host_group_desc": fmt.Sprintf("tf-testAccName-%d", rand),
			"groups.0.allocation_policy":         "Evenly",
			"groups.0.cpu_allocation_ratio":      "101",
			"groups.0.mem_allocation_ratio":      "50",
			"groups.0.disk_allocation_ratio":     "200",
			"groups.0.host_replace_policy":       "Manual",
			"groups.0.create_time":               CHECKSET,
		}
	}
	var fakeAlicloudCddcDedicatedHostGroupsNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"groups.#": "0",
		}
	}
	var alicloudSaeNamespaceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cddc_dedicated_host_groups.default",
		existMapFunc: existAlicloudCddcDedicatedHostGroupsNameMapFunc,
		fakeMapFunc:  fakeAlicloudCddcDedicatedHostGroupsNameMapFunc,
	}

	alicloudSaeNamespaceCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudCddcDedicatedHostGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccName-%d"
}
data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

resource "alicloud_cddc_dedicated_host_group" "default" {
	engine = "MySQL"
	vpc_id = data.alicloud_vpcs.default.ids.0
	cpu_allocation_ratio = 101
	mem_allocation_ratio = 50
	disk_allocation_ratio = 200
	allocation_policy = "Evenly"
	host_replace_policy = "Manual"
	dedicated_host_group_desc = var.name
}

data "alicloud_cddc_dedicated_host_groups" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}

func TestAccAlicloudCddcDedicatedHostGroupsSqlServerDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 200)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCddcDedicatedHostGroupsSqlServerDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cddc_dedicated_host_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCddcDedicatedHostGroupsSqlServerDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cddc_dedicated_host_group.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCddcDedicatedHostGroupsSqlServerDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cddc_dedicated_host_group.default.id}"]`,
			"name_regex": `"${alicloud_cddc_dedicated_host_group.default.dedicated_host_group_desc}"`,
		}),
		fakeConfig: testAccCheckAlicloudCddcDedicatedHostGroupsSqlServerDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cddc_dedicated_host_group.default.id}_fake"]`,
			"name_regex": `"${alicloud_cddc_dedicated_host_group.default.dedicated_host_group_desc}_fake"`,
		}),
	}
	// TODO: There is an api bug that the request parameter does not support SQLServer. This should reopen after the bug is fixed.
	//allConf := dataSourceTestAccConfig{
	//	existConfig: testAccCheckAlicloudCddcDedicatedHostGroupsSqlServerDataSourceName(rand, map[string]string{
	//		"ids":        `["${alicloud_cddc_dedicated_host_group.default.id}"]`,
	//		"engine":     `"SQLServer"`,
	//		"name_regex": `"${alicloud_cddc_dedicated_host_group.default.dedicated_host_group_desc}"`,
	//	}),
	//	fakeConfig: testAccCheckAlicloudCddcDedicatedHostGroupsSqlServerDataSourceName(rand, map[string]string{
	//		"ids":        `["${alicloud_cddc_dedicated_host_group.default.id}_fake"]`,
	//		"engine":     `"SQLServer"`,
	//		"name_regex": `"${alicloud_cddc_dedicated_host_group.default.dedicated_host_group_desc}_fake"`,
	//	}),
	//}
	var existAlicloudCddcDedicatedHostGroupsNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                            "1",
			"names.0":                            fmt.Sprintf("tf-testAccName-%d", rand),
			"ids.#":                              "1",
			"groups.#":                           "1",
			"groups.0.engine":                    "SQLServer",
			"groups.0.dedicated_host_group_desc": fmt.Sprintf("tf-testAccName-%d", rand),
			"groups.0.allocation_policy":         "Evenly",
			"groups.0.cpu_allocation_ratio":      "101",
			"groups.0.mem_allocation_ratio":      "50",
			"groups.0.disk_allocation_ratio":     "100",
			"groups.0.host_replace_policy":       "Manual",
			"groups.0.create_time":               CHECKSET,
		}
	}
	var fakeAlicloudCddcDedicatedHostGroupsNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"groups.#": "0",
		}
	}
	var alicloudSaeNamespaceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cddc_dedicated_host_groups.default",
		existMapFunc: existAlicloudCddcDedicatedHostGroupsNameMapFunc,
		fakeMapFunc:  fakeAlicloudCddcDedicatedHostGroupsNameMapFunc,
	}

	alicloudSaeNamespaceCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf)
}
func testAccCheckAlicloudCddcDedicatedHostGroupsSqlServerDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccName-%d"
}
data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

resource "alicloud_cddc_dedicated_host_group" "default" {
	engine = "SQLServer"
	vpc_id = data.alicloud_vpcs.default.ids.0
	open_permission = true
	cpu_allocation_ratio = 101
	mem_allocation_ratio = 50
	allocation_policy = "Evenly"
	host_replace_policy = "Manual"
	dedicated_host_group_desc = var.name
}

data "alicloud_cddc_dedicated_host_groups" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
