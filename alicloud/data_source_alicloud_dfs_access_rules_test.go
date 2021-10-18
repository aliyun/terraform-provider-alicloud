package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDfsAccessRulesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.DfsSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDfsAccessRulesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dfs_access_rule.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDfsAccessRulesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dfs_access_rule.default.id}_fake"]`,
		}),
	}

	var existAlicloudDfsAccessRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"rules.#":                 "1",
			"rules.0.access_group_id": CHECKSET,
			"rules.0.description":     fmt.Sprintf("tf-testAccAccessRule-%d", rand),
			"rules.0.network_segment": "192.0.2.0/24",
			"rules.0.priority":        "10",
			"rules.0.rw_access_type":  "RDWR",
		}
	}
	var fakeAlicloudDfsAccessRulesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudDfsAccessRulesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_dfs_access_rules.default",
		existMapFunc: existAlicloudDfsAccessRulesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudDfsAccessRulesDataSourceNameMapFunc,
	}
	alicloudDfsAccessRulesCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}
func testAccCheckAlicloudDfsAccessRulesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccAccessRule-%d"
}

resource "alicloud_dfs_access_group" "default" {
  network_type      = "VPC"
  access_group_name = var.name
  description       = var.name
}

resource "alicloud_dfs_access_rule" "default" {
  network_segment = "192.0.2.0/24"
  access_group_id = alicloud_dfs_access_group.default.id
  description     = var.name
  rw_access_type  = "RDWR"
  priority        = "10"
}

data "alicloud_dfs_access_rules" "default" {
  access_group_id = alicloud_dfs_access_group.default.id
  %s
}

`, rand, strings.Join(pairs, " \n "))
	return config
}
