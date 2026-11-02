package nas_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudNASAccessRuleDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100000, 999999)
	ipConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"source_cidr_ip":    `"${alicloud_nas_access_rule.default.source_cidr_ip}"`,
		}),
		FakeConfig: testAccCheckAlicloudAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"source_cidr_ip":    `"${alicloud_nas_access_rule.default.source_cidr_ip}_fake"`,
		}),
	}
	RWAccessConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"rw_access":         `"${alicloud_nas_access_rule.default.rw_access_type}"`,
		}),
		FakeConfig: testAccCheckAlicloudAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"rw_access":         `"${alicloud_nas_access_rule.default.rw_access_type}_fake"`,
		}),
	}
	UserAccessConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"user_access":       `"${alicloud_nas_access_rule.default.user_access_type}"`,
		}),
		FakeConfig: testAccCheckAlicloudAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"user_access":       `"${alicloud_nas_access_rule.default.user_access_type}_fake"`,
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"ids":               `["${alicloud_nas_access_rule.default.access_rule_id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"ids":               `["${alicloud_nas_access_rule.default.access_rule_id}_fake"]`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"user_access":       `"${alicloud_nas_access_rule.default.user_access_type}"`,
			"rw_access":         `"${alicloud_nas_access_rule.default.rw_access_type}"`,
			"ids":               `["${alicloud_nas_access_rule.default.access_rule_id}"]`,
			"source_cidr_ip":    `"${alicloud_nas_access_rule.default.source_cidr_ip}"`,
		}),
		FakeConfig: testAccCheckAlicloudAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"user_access":       `"${alicloud_nas_access_rule.default.user_access_type}_fake"`,
			"rw_access":         `"${alicloud_nas_access_rule.default.rw_access_type}_fake"`,
			"ids":               `["${alicloud_nas_access_rule.default.access_rule_id}"]`,
			"source_cidr_ip":    `"${alicloud_nas_access_rule.default.source_cidr_ip}_fake"`,
		}),
	}
	accessRuleCheckInfo.DataSourceTestCheck(t, rand, ipConf, RWAccessConf, UserAccessConf, idsConf, allConf)
}

func testAccCheckAlicloudAccessRuleDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
        	default = "tf-testAccAccessGroupsdatasource-%d"
}
resource "alicloud_nas_access_group" "default" {
        	access_group_name = "${var.name}"
	        access_group_type = "Vpc"
	        description = "tf-testAccAccessGroupsdatasource"
}
resource "alicloud_nas_access_rule" "default" {
        	access_group_name = "${alicloud_nas_access_group.default.access_group_name}"
	        source_cidr_ip = "168.1.1.0/16"
        	rw_access_type = "RDWR"
	        user_access_type = "no_squash"
	        priority = 2
}
data "alicloud_nas_access_rules" "default" {
		%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existAccessRuleMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"rules.#":                "1",
		"rules.0.source_cidr_ip": "168.1.1.0/16",
		"rules.0.priority":       "2",
		"rules.0.access_rule_id": acctest.CHECKSET,
		"rules.0.user_access":    "no_squash",
		"rules.0.rw_access":      "RDWR",
		"ids.#":                  "1",
		"ids.0":                  "1",
	}
}

var fakeAccessRuleMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"rules.#": "0",
		"ids.#":   "0",
	}
}

var accessRuleCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_nas_access_rules.default",
	ExistMapFunc: existAccessRuleMapCheck,
	FakeMapFunc:  fakeAccessRuleMapCheck,
}
