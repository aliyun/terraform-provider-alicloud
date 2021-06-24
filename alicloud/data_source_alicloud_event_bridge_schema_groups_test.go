package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEventBridgeSchemaGroupsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeSchemaGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_event_bridge_schema_group.default.group_id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeSchemaGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_event_bridge_schema_group.default.group_id}_fake"]`,
		}),
	}
	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeSchemaGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_event_bridge_schema_group.default.group_id}"]`,
			"description_regex": `"${alicloud_event_bridge_schema_group.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeSchemaGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_event_bridge_schema_group.default.group_id}"]`,
			"description_regex": `"${alicloud_event_bridge_schema_group.default.description}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEventBridgeSchemaGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_event_bridge_schema_group.default.group_id}"]`,
			"description_regex": `"${alicloud_event_bridge_schema_group.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudEventBridgeSchemaGroupsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_event_bridge_schema_group.default.group_id}_fake"]`,
			"description_regex": `"${alicloud_event_bridge_schema_group.default.description}_fake"`,
		}),
	}
	var existAlicloudEventBridgeSchemaGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "1",
			"groups.#":             "1",
			"groups.0.group_id":    fmt.Sprintf(`tf-testAccSchemaGroup-%d`, rand),
			"groups.0.description": fmt.Sprintf(`tf-testAccSchemaGroup-%d`, rand),
			"groups.0.format":      `OPEN_API_3_0`,
		}
	}
	var fakeAlicloudEventBridgeSchemaGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"descriptions.#": "0",
		}
	}
	var alicloudEventBridgeSchemaGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_event_bridge_schema_groups.default",
		existMapFunc: existAlicloudEventBridgeSchemaGroupsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEventBridgeSchemaGroupsDataSourceNameMapFunc,
	}
	alicloudEventBridgeSchemaGroupsCheckInfo.dataSourceTestCheck(t, rand, idsConf, descriptionRegexConf, allConf)
}
func testAccCheckAlicloudEventBridgeSchemaGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccSchemaGroup-%d"
}

resource "alicloud_event_bridge_schema_group" "default" {
	group_id = var.name
	description = var.name
	format = "OPEN_API_3_0"
}

data "alicloud_event_bridge_schema_groups" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
