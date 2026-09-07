package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudVPCNetworkAclsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNetworkAclsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_network_acl.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudNetworkAclsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_network_acl.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNetworkAclsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_network_acl.default.network_acl_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudNetworkAclsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_network_acl.default.network_acl_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNetworkAclsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_network_acl.default.id}"]`,
			"status": `"${alicloud_network_acl.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudNetworkAclsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_network_acl.default.id}"]`,
			"status": `"Modifying"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNetworkAclsDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_network_acl.default.id}"]`,
			"name_regex":       `"${alicloud_network_acl.default.network_acl_name}"`,
			"network_acl_name": `"${alicloud_network_acl.default.network_acl_name}"`,
			"status":           `"${alicloud_network_acl.default.status}"`,
		}),
		fakeConfig: testAccCheckAlicloudNetworkAclsDataSourceName(rand, map[string]string{
			"ids":              `["${alicloud_network_acl.default.id}"]`,
			"name_regex":       `"${alicloud_network_acl.default.network_acl_name}_fake"`,
			"network_acl_name": `"${alicloud_network_acl.default.network_acl_name}_fake"`,
			"status":           `"Modifying"`,
		}),
	}
	var existAlicloudNetworkAclsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"names.#":                      "1",
			"acls.#":                       "1",
			"acls.0.description":           fmt.Sprintf("tf-testAccNetworkAcl-%d", rand),
			"acls.0.egress_acl_entries.#":  "1",
			"acls.0.ingress_acl_entries.#": "1",
			"acls.0.network_acl_name":      fmt.Sprintf("tf-testAccNetworkAcl-%d", rand),
			"acls.0.vpc_id":                CHECKSET,
			"acls.0.status":                "Available",
		}
	}
	var fakeAlicloudNetworkAclsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudNetworkAclsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_network_acls.default",
		existMapFunc: existAlicloudNetworkAclsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudNetworkAclsDataSourceNameMapFunc,
	}
	alicloudNetworkAclsCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudNetworkAclsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccNetworkAcl-%d"
}

resource "alicloud_vpc" "default" {
	vpc_name = "${var.name}"
}

resource "alicloud_network_acl" "default" {
	description = "${var.name}"
	network_acl_name = "${var.name}"
	vpc_id = "${alicloud_vpc.default.id}"
}

data "alicloud_network_acls" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
