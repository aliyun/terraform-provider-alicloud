package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOnsGroupsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_ons_groups.default"
	name := fmt.Sprintf("GID-tf-testacconsgroup%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceOnsGroupsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":    "${alicloud_ons_instance.default.id}",
			"group_id_regex": "${alicloud_ons_group.default.group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":    "${alicloud_ons_instance.default.id}",
			"group_id_regex": "${alicloud_ons_group.default.group_id}_fake",
		}),
	}

	var existOnsGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"groups.#":                    "1",
			"groups.0.id":                 fmt.Sprintf("GID-tf-testacconsgroup%v", rand),
			"groups.0.independent_naming": "true",
			"groups.0.remark":             "alicloud_ons_group_remark",
		}
	}

	var fakeOnsGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"groups.#": "0",
		}
	}

	var onsGroupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOnsGroupsMapFunc,
		fakeMapFunc:  fakeOnsGroupsMapFunc,
	}

	onsGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf)
}

func dataSourceOnsGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "group_id" {
 default = "%v"
}

resource "alicloud_ons_instance" "default" {
name = "tf-testaccOnsInstanceGroupbasic"
}

resource "alicloud_ons_group" "default" {
  instance_id = "${alicloud_ons_instance.default.id}"
  group_id = "${var.group_id}"
  remark = "alicloud_ons_group_remark"
  read_enable = "true"
}
`, name)
}
