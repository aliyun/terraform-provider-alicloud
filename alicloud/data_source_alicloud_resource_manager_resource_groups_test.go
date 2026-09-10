package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudResourceManagerResourceGroupsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerResourceGroupsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_resource_group.example.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerResourceGroupsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_resource_group.example.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerResourceGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_resource_group.example.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerResourceGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_resource_group.example.id}_fake"]`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerResourceGroupsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_resource_group.example.name}"`,
			"status":     `"OK"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerResourceGroupsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_resource_group.example.name}"`,
			"status":     `"Creating"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerResourceGroupsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_resource_group.example.name}"`,
			"ids":        `["${alicloud_resource_manager_resource_group.example.id}"]`,
			"status":     `"OK"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerResourceGroupsSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_resource_manager_resource_group.example.name}_fake"`,
			"ids":        `["${alicloud_resource_manager_resource_group.example.id}"]`,
			"status":     `"OK"`,
		}),
	}

	var existResourceManagerResourceGroupsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#":                     "1",
			"names.#":                      "1",
			"ids.#":                        "1",
			"groups.0.id":                  CHECKSET,
			"groups.0.name":                fmt.Sprintf("tf-%d", rand),
			"groups.0.resource_group_name": fmt.Sprintf("tf-%d", rand),
			"groups.0.display_name":        fmt.Sprintf("terraform_test_%d", rand),
			"groups.0.account_id":          CHECKSET,
			"groups.0.status":              "OK",
			"groups.0.region_statuses.#":   CHECKSET,
		}
	}

	var fakeResourceManagerResourceGroupsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#": "0",
			"ids.#":    "0",
			"names.#":  "0",
		}
	}

	var resourceGroupsRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_resource_groups.example",
		existMapFunc: existResourceManagerResourceGroupsRecordsMapFunc,
		fakeMapFunc:  fakeResourceManagerResourceGroupsRecordsMapFunc,
	}

	resourceGroupsRecordsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, statusConf, allConf)

}

func testAccCheckAlicloudResourceManagerResourceGroupsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
resource "alicloud_resource_manager_resource_group" "example"{
	resource_group_name = "tf-%[1]d"
	display_name = "terraform_test_%[1]d"
}

data "alicloud_resource_manager_resource_groups" "example"{
%s
}

`, rand, strings.Join(pairs, "\n   "))
	return config
}
