package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudNASAllAccessRuleDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100000, 999999)
	ipConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAllAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"source_cidr_ip":    `"${alicloud_nas_all_access_rule.default.source_cidr_ip}"`,
		}),
		fakeConfig: testAccCheckAlicloudAllAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"source_cidr_ip":    `"${alicloud_nas_all_access_rule.default.source_cidr_ip}_fake"`,
		}),
	}
	RWAccessConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAllAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"rw_access":         `"${alicloud_nas_all_access_rule.default.rw_access_type}"`,
		}),
		fakeConfig: testAccCheckAlicloudAllAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"rw_access":         `"${alicloud_nas_all_access_rule.default.rw_access_type}_fake"`,
		}),
	}
	UserAccessConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAllAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"user_access":       `"${alicloud_nas_all_access_rule.default.user_access_type}"`,
		}),
		fakeConfig: testAccCheckAlicloudAllAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"user_access":       `"${alicloud_nas_all_access_rule.default.user_access_type}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAllAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"ids":               `["${alicloud_nas_all_access_rule.default.access_rule_id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAllAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"ids":               `["${alicloud_nas_all_access_rule.default.access_rule_id}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAllAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"user_access":       `"${alicloud_nas_all_access_rule.default.user_access_type}"`,
			"rw_access":         `"${alicloud_nas_all_access_rule.default.rw_access_type}"`,
			"ids":               `["${alicloud_nas_all_access_rule.default.access_rule_id}"]`,
			"source_cidr_ip":    `"${alicloud_nas_all_access_rule.default.source_cidr_ip}"`,
		}),
		fakeConfig: testAccCheckAlicloudAllAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${alicloud_nas_access_group.default.access_group_name}"`,
			"user_access":       `"${alicloud_nas_all_access_rule.default.user_access_type}_fake"`,
			"rw_access":         `"${alicloud_nas_all_access_rule.default.rw_access_type}_fake"`,
			"ids":               `["${alicloud_nas_all_access_rule.default.access_rule_id}_fake"]`,
			"source_cidr_ip":    `"${alicloud_nas_all_access_rule.default.source_cidr_ip}_fake"`,
		}),
	}
	allAccessRuleCheckInfo.dataSourceTestCheck(t, rand, ipConf, RWAccessConf, UserAccessConf, idsConf, allConf)
}

func TestAccAlicloudNASAllAccessRuleDataSourceIpv6SourceCidrIp(t *testing.T) {
	rand := acctest.RandIntRange(100000, 999999)
	Ipv6SourceCidrIpConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAllAccessRuleIpv6SourceCidrIpDataSourceConfig(rand, map[string]string{
			"access_group_name":   `"${alicloud_nas_access_group.default.access_group_name}"`,
			"user_access":         `"${alicloud_nas_all_access_rule.default.user_access_type}"`,
			"rw_access":           `"${alicloud_nas_all_access_rule.default.rw_access_type}"`,
			"ids":                 `["${alicloud_nas_all_access_rule.default.access_rule_id}"]`,
			"ipv6_source_cidr_ip": `"${alicloud_nas_all_access_rule.default.ipv6_source_cidr_ip}"`,
		}),
		fakeConfig: testAccCheckAlicloudAllAccessRuleIpv6SourceCidrIpDataSourceConfig(rand, map[string]string{
			"access_group_name":   `"${alicloud_nas_access_group.default.access_group_name}"`,
			"user_access":         `"${alicloud_nas_all_access_rule.default.user_access_type}_fake"`,
			"rw_access":           `"${alicloud_nas_all_access_rule.default.rw_access_type}_fake"`,
			"ids":                 `["${alicloud_nas_all_access_rule.default.access_rule_id}_fake"]`,
			"ipv6_source_cidr_ip": `"${alicloud_nas_all_access_rule.default.ipv6_source_cidr_ip}_fake"`,
		}),
	}
	allAccessRuleCheckIpv6SourceCidrIpInfo.dataSourceTestCheck(t, rand, Ipv6SourceCidrIpConf)
}

func testAccCheckAlicloudAllAccessRuleDataSourceConfig(rand int, attrMap map[string]string) string {
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
	        file_system_type = "standard"
}
resource "alicloud_nas_all_access_rule" "default" {
        	access_group_name = "${alicloud_nas_access_group.default.access_group_name}"
	        source_cidr_ip = "168.1.1.0/16"
        	rw_access_type = "RDWR"
	        user_access_type = "no_squash"
	        priority = 2
	        file_system_type = "${alicloud_nas_access_group.default.file_system_type}"
}
data "alicloud_nas_all_access_rules" "default" {
            file_system_type = alicloud_nas_access_group.default.file_system_type
		%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

func testAccCheckAlicloudAllAccessRuleIpv6SourceCidrIpDataSourceConfig(rand int, attrMap map[string]string) string {
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
	        file_system_type = "standard"
}
resource "alicloud_nas_all_access_rule" "default" {
        	access_group_name = "${alicloud_nas_access_group.default.access_group_name}"
	        ipv6_source_cidr_ip = "::1/128"
        	rw_access_type = "RDWR"
	        user_access_type = "no_squash"
	        priority = 2
	        file_system_type = "${alicloud_nas_access_group.default.file_system_type}"
}
data "alicloud_nas_all_access_rules" "default" {
            file_system_type = alicloud_nas_access_group.default.file_system_type
		%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existAllAccessRuleMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"rules.#":                     "1",
		"rules.0.source_cidr_ip":      "168.1.1.0/16",
		"rules.0.priority":            "2",
		"rules.0.access_rule_id":      CHECKSET,
		"rules.0.user_access":         "no_squash",
		"rules.0.rw_access":           "RDWR",
		"rules.0.ipv6_source_cidr_ip": "",
		"ids.#":                       "1",
		"ids.0":                       CHECKSET,
	}
}

var existAllAccessRuleIpv6SourceCidrIpMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"rules.#":                     "1",
		"rules.0.source_cidr_ip":      "",
		"rules.0.priority":            "2",
		"rules.0.access_rule_id":      CHECKSET,
		"rules.0.user_access":         "no_squash",
		"rules.0.rw_access":           "RDWR",
		"rules.0.ipv6_source_cidr_ip": "::1/128",
		"ids.#":                       "1",
		"ids.0":                       CHECKSET,
	}
}

var fakeAllAccessRuleMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"rules.#": "0",
		"ids.#":   "0",
	}
}

var allAccessRuleCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_nas_all_access_rules.default",
	existMapFunc: existAllAccessRuleMapCheck,
	fakeMapFunc:  fakeAllAccessRuleMapCheck,
}

var allAccessRuleCheckIpv6SourceCidrIpInfo = dataSourceAttr{
	resourceId:   "data.alicloud_nas_all_access_rules.default",
	existMapFunc: existAllAccessRuleIpv6SourceCidrIpMapCheck,
	fakeMapFunc:  fakeAllAccessRuleMapCheck,
}
