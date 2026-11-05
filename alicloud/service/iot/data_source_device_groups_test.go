package iot_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudIotDeviceGroupsDataSource(t *testing.T) {
	os.Setenv("ALICLOUD_REGION", "cn-shanghai")
	rand := tfacctest.RandIntRange(1, 100)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudIotDeviceGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_iot_device_group.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudIotDeviceGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_iot_device_group.default.id}_fake"]`,
		}),
	}

	groupNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudIotDeviceGroupsDataSourceName(rand, map[string]string{
			"group_name": `"${alicloud_iot_device_group.default.group_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudIotDeviceGroupsDataSourceName(rand, map[string]string{
			"group_name": `"${alicloud_iot_device_group.default.group_name}_fake"`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudIotDeviceGroupsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_iot_device_group.default.group_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudIotDeviceGroupsDataSourceName(rand, map[string]string{
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
	var alicloudIotDeviceGroupsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_iot_device_groups.default",
		ExistMapFunc: existAlicloudIotDeviceGroupsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudIotDeviceGroupsDataSourceNameMapFunc,
	}
	alicloudIotDeviceGroupsCheckInfo.DataSourceTestCheck(t, rand, idsConf, groupNameConf, nameRegexConf)
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
