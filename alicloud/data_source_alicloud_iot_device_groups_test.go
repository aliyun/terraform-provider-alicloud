package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudIotDeviceGroupsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 100)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudIotDeviceGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_iot_device_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudIotDeviceGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_iot_device_group.default.id}_fake"]`,
		}),
	}

	groupNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudIotDeviceGroupsDataSourceName(rand, map[string]string{
			"group_name": `"${alicloud_iot_device_group.default.group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudIotDeviceGroupsDataSourceName(rand, map[string]string{
			"group_name": `"${alicloud_iot_device_group.default.group_name}_fake"`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudIotDeviceGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_iot_device_group.default.group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudIotDeviceGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_iot_device_group.default.group_name}_fake"`,
		}),
	}

	var existAlicloudIotDeviceGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "1",
			"groups.#":               "1",
			"groups.0.group_name":    fmt.Sprintf("tf_testAccDeviceGroups_%d", rand),
			"groups.0.group_desc":    fmt.Sprintf("tf_testAccDeviceGroups_%d", rand),
			"groups.0.device_active": "",
			"groups.0.device_count":  "",
			"groups.0.device_online": "",
		}
	}
	var fakeAlicloudIotDeviceGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudIotDeviceGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_iot_device_groups.default",
		existMapFunc: existAlicloudIotDeviceGroupsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudIotDeviceGroupsDataSourceNameMapFunc,
	}
	alicloudIotDeviceGroupsCheckInfo.dataSourceTestCheck(t, rand, idsConf, groupNameConf, nameRegexConf)
}
func testAccCheckAlicloudIotDeviceGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf_testAccDeviceGroups_%d"
}
resource "alicloud_iot_device_group" "default"{
  group_name = var.name
  group_desc = var.name
}

data "alicloud_iot_device_groups" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
