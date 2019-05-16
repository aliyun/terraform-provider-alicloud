package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"strings"
	"testing"
)

func TestAccAlicloudSlbAclsDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_slb_acl.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbAclsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_slb_acl.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbAclsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_slb_acl.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSlbAclsDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_slb_acl.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlbAclsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_slb_acl.default.id}"]`,
			"name_regex": `"${alicloud_slb_acl.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlbAclsDataSourceConfig(rand, map[string]string{
			"ids":        `["${alicloud_slb_acl.default.id}_fake"]`,
			"name_regex": `"${alicloud_slb_acl.default.name}"`,
		}),
	}

	var existDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"acls.#":                     "1",
			"ids.#":                      "1",
			"names.#":                    "1",
			"acls.0.id":                  CHECKSET,
			"acls.0.name":                fmt.Sprintf("tf-testAccSlbAclDataSourceBisic-%d", rand),
			"acls.0.ip_version":          "ipv4",
			"acls.0.entry_list.#":        "2",
			"acls.0.related_listeners.#": "0",
		}
	}

	var fakeDnsRecordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"acls.#":  "0",
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var slbaclsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_slb_acls.default",
		existMapFunc: existDnsRecordsMapFunc,
		fakeMapFunc:  fakeDnsRecordsMapFunc,
	}

	slbaclsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlicloudSlbAclsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSlbAclDataSourceBisic-%d"
}
variable "ip_version" {
	default = "ipv4"
}

resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list = [
    {
      entry="10.10.10.0/24"
      comment="first"
    },
    {
      entry="168.10.10.0/24"
      comment="second"
    }
  ]
}

data "alicloud_slb_acls" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
