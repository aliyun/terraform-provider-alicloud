package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudBastionhostHostGroupsDataSource(t *testing.T) {
	resourceId := "data.alicloud_bastionhost_host_groups.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccBastionhostHostGroupsTest%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceBastionhostHostGroupsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host_group.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_host_group.default.host_group_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host_group.default.instance_id}",
			"ids":         []string{"${alicloud_bastionhost_host_group.default.id}-fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host_group.default.instance_id}",
			"name_regex":  "${alicloud_bastionhost_host_group.default.host_group_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_bastionhost_host_group.default.instance_id}",
			"name_regex":  "${alicloud_bastionhost_host_group.default.host_group_name}" + "fake",
		}),
	}
	userNameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":     "${alicloud_bastionhost_host_group.default.instance_id}",
			"host_group_name": "${alicloud_bastionhost_host_group.default.host_group_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":     "${alicloud_bastionhost_host_group.default.instance_id}",
			"host_group_name": "${alicloud_bastionhost_host_group.default.host_group_name}" + "fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":     "${alicloud_bastionhost_host_group.default.instance_id}",
			"name_regex":      name,
			"host_group_name": name,
			"ids":             []string{"${alicloud_bastionhost_host_group.default.host_group_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":     "${alicloud_bastionhost_host_group.default.instance_id}",
			"host_group_name": name + "fake",
			"name_regex":      name + "fake",
			"ids":             []string{"${alicloud_bastionhost_host_group.default.id}-fake"},
		}),
	}
	var existBastionhostHostGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"ids.0":                    CHECKSET,
			"names.#":                  "1",
			"names.0":                  name,
			"groups.#":                 "1",
			"groups.0.id":              CHECKSET,
			"groups.0.comment":         "",
			"groups.0.instance_id":     CHECKSET,
			"groups.0.host_group_id":   CHECKSET,
			"groups.0.host_group_name": name,
		}
	}

	var fakeBastionhostHostGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"groups.#": "0",
		}
	}

	var BastionhostHostGroupsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existBastionhostHostGroupsMapFunc,
		fakeMapFunc:  fakeBastionhostHostGroupsMapFunc,
	}

	BastionhostHostGroupsInfo.dataSourceTestCheck(t, 0, idsConf, nameRegexConf, userNameRegexConf, allConf)
}

func dataSourceBastionhostHostGroupsDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_bastionhost_instances" "default" {}

resource "alicloud_bastionhost_host_group" "default" {
  instance_id     = data.alicloud_bastionhost_instances.default.ids.0
  host_group_name      = var.name
}
`, name)
}
