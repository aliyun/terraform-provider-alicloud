package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSASGroupsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	groupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSasGroupSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_security_center_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSasGroupSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_security_center_group.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSasGroupSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_security_center_group.default.group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSasGroupSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_security_center_group.default.group_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSasGroupSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_security_center_group.default.id}"]`,
			"name_regex": `"${alicloud_security_center_group.default.group_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSasGroupSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_security_center_group.default.id}"]`,
			"name_regex": `"${alicloud_security_center_group.default.group_name}_fake"`,
		}),
	}

	SasGroupCheckInfo.dataSourceTestCheck(t, rand, groupIdConf, nameRegexConf, allConf)
}

func testAccCheckAlicloudSasGroupSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSasGroupsDataSource%d"
}

resource "alicloud_security_center_group" "default" {
  group_name = var.name
}

data "alicloud_security_center_groups" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}

var existSasGroupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#":            "1",
		"groups.0.group_id":   CHECKSET,
		"groups.0.group_name": fmt.Sprintf("tf-testAccSasGroupsDataSource%d", rand),
	}
}

var fakeSasGroupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#": "0",
	}
}

var SasGroupCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_security_center_groups.default",
	existMapFunc: existSasGroupMapFunc,
	fakeMapFunc:  fakeSasGroupMapFunc,
}
