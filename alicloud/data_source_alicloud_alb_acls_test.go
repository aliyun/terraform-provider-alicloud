package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudAlbAclsDataSource(t *testing.T) {
	rand := acctest.RandInt()

	aclIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"acl_ids": `["${alicloud_alb_acl.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"acl_ids": `["${alicloud_alb_acl.default.id}_fake"]`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_acl.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_acl.default.id}_fakeid"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_acl.default.acl_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_acl.default.acl_name}_fake"`,
		}),
	}

	aclNamesConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"acl_name": `"${alicloud_alb_acl.default.acl_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"acl_name": `"${alicloud_alb_acl.default.acl_name}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_acl.default.id}"]`,
			"status": `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_acl.default.id}"]`,
			"status": `"Configuring"`,
		}),
	}

	ResourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"name_regex":        `"${alicloud_alb_acl.default.acl_name}"`,
			"resource_group_id": `"${alicloud_alb_acl.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"name_regex":        `"${alicloud_alb_acl.default.acl_name}"`,
			"resource_group_id": `"${alicloud_alb_acl.default.resource_group_id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_acl.default.id}"]`,
			"acl_ids":           `["${alicloud_alb_acl.default.id}"]`,
			"acl_name":          `"${alicloud_alb_acl.default.acl_name}"`,
			"name_regex":        `"${alicloud_alb_acl.default.acl_name}"`,
			"resource_group_id": `"${alicloud_alb_acl.default.resource_group_id}"`,
			"status":            `"Available"`,
		}),
		fakeConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_acl.default.id}_fake"]`,
			"acl_ids":           `["${alicloud_alb_acl.default.id}_fake"]`,
			"acl_name":          `"${alicloud_alb_acl.default.acl_name}"`,
			"status":            `"Configuring"`,
			"name_regex":        `"${alicloud_alb_acl.default.acl_name}_fake"`,
			"resource_group_id": `"${alicloud_alb_acl.default.resource_group_id}_fake"`,
		}),
	}

	var existDataAlicloudAlbAclsSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           "1",
			"acls.#":          "1",
			"acls.0.acl_name": fmt.Sprintf("tf-testAccAlbAcl%d", rand),
			"acls.0.status":   "Available",
		}
	}
	var fakeDataAlicloudAlbAclsSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":  "0",
			"acls.#": "0",
		}
	}
	var alicloudAlbAclCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alb_acls.default",
		existMapFunc: existDataAlicloudAlbAclsSourceNameMapFunc,
		fakeMapFunc:  fakeDataAlicloudAlbAclsSourceNameMapFunc,
	}
	alicloudAlbAclCheckInfo.dataSourceTestCheck(t, rand, aclIdsConf, idsConf, nameRegexConf, aclNamesConf, statusConf, ResourceGroupIdConf, allConf)
}
func testAccCheckAlicloudAlbAclDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccAlbAcl%d"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_alb_acl" "default" {
	acl_name = var.name
}

data "alicloud_alb_acls" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
