package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudALBAclsDataSource(t *testing.T) {
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
			"ids.#":                            "1",
			"names.#":                          "1",
			"names.0":                          fmt.Sprintf("tf-testAccAlbAcl%d", rand),
			"acls.#":                           "1",
			"acls.0.acl_name":                  fmt.Sprintf("tf-testAccAlbAcl%d", rand),
			"acls.0.status":                    "Available",
			"acls.0.id":                        CHECKSET,
			"acls.0.acl_entries.#":             "1",
			"acls.0.acl_entries.0.description": "description",
			"acls.0.acl_entries.0.entry":       "10.0.0.0/24",
			"acls.0.acl_entries.0.status":      CHECKSET,
			"acls.0.acl_id":                    CHECKSET,
			"acls.0.address_ip_version":        CHECKSET,
			"acls.0.resource_group_id":         CHECKSET,
		}
	}
	var fakeDataAlicloudAlbAclsSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"acls.#":  "0",
		}
	}
	var alicloudAlbAclCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alb_acls.default",
		existMapFunc: existDataAlicloudAlbAclsSourceNameMapFunc,
		fakeMapFunc:  fakeDataAlicloudAlbAclsSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
	}
	alicloudAlbAclCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, aclIdsConf, idsConf, nameRegexConf, aclNamesConf, statusConf, ResourceGroupIdConf, allConf)
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
  acl_name          = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  acl_entries {
    description = "description"
    entry       = "10.0.0.0/24"
  }
}

data "alicloud_alb_acls" "default" {
  enable_details = true
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
