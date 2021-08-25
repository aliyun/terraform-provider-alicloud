package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDfsAccsessGroupsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDfsAccsessGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dfs_access_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDfsAccsessGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dfs_access_group.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDfsAccsessGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_dfs_access_group.default.access_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDfsAccsessGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_dfs_access_group.default.access_group_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDfsAccsessGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_dfs_access_group.default.id}"]`,
			"name_regex": `"${alicloud_dfs_access_group.default.access_group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudDfsAccsessGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_dfs_access_group.default.id}_fake"]`,
			"name_regex": `"${alicloud_dfs_access_group.default.access_group_name}_fake"`,
		}),
	}
	var existAlicloudDfsAccsessGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"groups.#":                   "1",
			"groups.0.network_type":      "VPC",
			"groups.0.description":       fmt.Sprintf("tf-testAccDfsAccsessGroupAcc-%d", rand),
			"groups.0.access_group_name": fmt.Sprintf("tf-testAccDfsAccsessGroupAcc-%d", rand),
		}
	}
	var fakeAlicloudDfsAccsessGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"groups.#": "0",
		}
	}
	var AlicloudDfsAccsessGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_dfs_access_groups.default",
		existMapFunc: existAlicloudDfsAccsessGroupsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudDfsAccsessGroupsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.DfsSupportRegions)
	}
	AlicloudDfsAccsessGroupsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudDfsAccsessGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccDfsAccsessGroupAcc-%d"
}

resource "alicloud_dfs_access_group" "default" {
	network_type = "VPC"
	access_group_name = var.name
	description =  var.name
}

data "alicloud_dfs_access_groups" "default" {
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
