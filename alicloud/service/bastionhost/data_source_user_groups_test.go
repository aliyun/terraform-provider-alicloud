package bastionhost_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudBastionhostUserGroupsDataSource(t *testing.T) {
	resourceId := "data.alicloud_bastionhost_user_groups.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccBastionhostUserGroupsTest%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceBastionhostUserGroupsDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_user_group.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_user_group.default.user_group_id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_user_group.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_user_group.default.id}-fake"},
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_user_group.default.instance_id}",
			"name_regex":  name,
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_user_group.default.instance_id}",
			"name_regex":  name + "fake",
		}),
	}
	userGroupNameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id":     "${alicloud_bastionhost_user_group.default.instance_id}",
			"user_group_name": name,
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":     "${alicloud_bastionhost_user_group.default.instance_id}",
			"user_group_name": name + "fake",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id":     "${alicloud_bastionhost_user_group.default.instance_id}",
			"name_regex":      name,
			"user_group_name": name,
			"ids":             []string{"${alicloud_bastionhost_user_group.default.user_group_id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":     "${alicloud_bastionhost_user_group.default.instance_id}",
			"user_group_name": name + "fake",
			"name_regex":      name + "fake",
			"ids":             []string{"${alicloud_bastionhost_user_group.default.id}-fake"},
		}),
	}
	var existBastionhostUserGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"ids.0":                    acctest.CHECKSET,
			"names.#":                  "1",
			"names.0":                  name,
			"groups.#":                 "1",
			"groups.0.id":              acctest.CHECKSET,
			"groups.0.comment":         "",
			"groups.0.instance_id":     acctest.CHECKSET,
			"groups.0.user_group_id":   acctest.CHECKSET,
			"groups.0.user_group_name": name,
		}
	}

	var fakeBastionhostUserGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"groups.#": "0",
		}
	}

	var BastionhostUserGroupsInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existBastionhostUserGroupsMapFunc,
		FakeMapFunc:  fakeBastionhostUserGroupsMapFunc,
	}

	BastionhostUserGroupsInfo.DataSourceTestCheck(t, 0, idsConf, nameRegexConf, userGroupNameRegexConf, allConf)
}

func dataSourceBastionhostUserGroupsDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_bastionhost_instances" "default" {}

resource "alicloud_bastionhost_user_group" "default" {
  instance_id     = data.alicloud_bastionhost_instances.default.ids.0
  user_group_name = var.name
}
`, name)
}
