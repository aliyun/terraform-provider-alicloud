package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGaAclsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaAclsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_acl.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudGaAclsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ga_acl.default.id}_fake"]`,
		}),
	}
	aclNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaAclsDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_ga_acl.default.id}"]`,
			"acl_name": `"${alicloud_ga_acl.default.acl_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudGaAclsDataSourceName(rand, map[string]string{
			"ids":      `["${alicloud_ga_acl.default.id}"]`,
			"acl_name": `"${alicloud_ga_acl.default.acl_name}_fake"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaAclsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ga_acl.default.acl_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudGaAclsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ga_acl.default.acl_name}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaAclsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ga_acl.default.id}"]`,
			"status": `"active"`,
		}),
		fakeConfig: testAccCheckAlicloudGaAclsDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ga_acl.default.id}"]`,
			"status": `"configuring"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGaAclsDataSourceName(rand, map[string]string{
			"acl_name":   `"${alicloud_ga_acl.default.acl_name}"`,
			"ids":        `["${alicloud_ga_acl.default.id}"]`,
			"name_regex": `"${alicloud_ga_acl.default.acl_name}"`,
			"status":     `"active"`,
		}),
		fakeConfig: testAccCheckAlicloudGaAclsDataSourceName(rand, map[string]string{
			"acl_name":   `"${alicloud_ga_acl.default.acl_name}_fake"`,
			"ids":        `["${alicloud_ga_acl.default.id}_fake"]`,
			"name_regex": `"${alicloud_ga_acl.default.acl_name}_fake"`,
			"status":     `"configuring"`,
		}),
	}
	var existAlicloudGaAclsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"names.#":                                "1",
			"acls.#":                                 "1",
			"acls.0.acl_entries.#":                   "1",
			"acls.0.acl_entries.0.entry":             "192.168.1.0/24",
			"acls.0.acl_entries.0.entry_description": "tf-test1",
			"acls.0.id":                              CHECKSET,
			"acls.0.acl_id":                          CHECKSET,
			"acls.0.acl_name":                        fmt.Sprintf("tf-testAccAcl-%d", rand),
			"acls.0.address_ip_version":              "IPv4",
			"acls.0.status":                          "active",
		}
	}
	var fakeAlicloudGaAclsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudGaAclsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ga_acls.default",
		existMapFunc: existAlicloudGaAclsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudGaAclsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudGaAclsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, aclNameConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudGaAclsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccAcl-%d"
}

resource "alicloud_ga_acl" "default" {
  acl_name           = var.name
  address_ip_version = "IPv4"
  acl_entries {
    entry             = "192.168.1.0/24"
    entry_description = "tf-test1"
  }
}

data "alicloud_ga_acls" "default" {	
	enable_details = true
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
