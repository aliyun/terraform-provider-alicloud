package dfs_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudDfsAccessGroupsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDfsAccsessGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dfs_access_group.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudDfsAccsessGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dfs_access_group.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDfsAccsessGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_dfs_access_group.default.access_group_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudDfsAccsessGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_dfs_access_group.default.access_group_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDfsAccsessGroupsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_dfs_access_group.default.id}"]`,
			"name_regex": `"${alicloud_dfs_access_group.default.access_group_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudDfsAccsessGroupsDataSourceName(rand, map[string]string{
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
	var AlicloudDfsAccsessGroupsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_dfs_access_groups.default",
		ExistMapFunc: existAlicloudDfsAccsessGroupsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudDfsAccsessGroupsDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.DfsSupportRegions)
	}
	AlicloudDfsAccsessGroupsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
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
