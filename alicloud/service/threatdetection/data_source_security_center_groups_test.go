package threatdetection_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudSASGroupsDataSource(t *testing.T) {
	rand := tFacctest.RandIntRange(1000000, 9999999)

	groupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSasGroupSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_security_center_group.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudSasGroupSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_security_center_group.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSasGroupSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_security_center_group.default.group_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudSasGroupSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_security_center_group.default.group_name}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSasGroupSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_security_center_group.default.id}"]`,
			"name_regex": `"${alicloud_security_center_group.default.group_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudSasGroupSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_security_center_group.default.id}"]`,
			"name_regex": `"${alicloud_security_center_group.default.group_name}_fake"`,
		}),
	}

	SasGroupCheckInfo.DataSourceTestCheck(t, rand, groupIdConf, nameRegexConf, allConf)
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
		"groups.0.group_id":   acctest.CHECKSET,
		"groups.0.group_name": fmt.Sprintf("tf-testAccSasGroupsDataSource%d", rand),
	}
}

var fakeSasGroupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#": "0",
	}
}

var SasGroupCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_security_center_groups.default",
	ExistMapFunc: existSasGroupMapFunc,
	FakeMapFunc:  fakeSasGroupMapFunc,
}
