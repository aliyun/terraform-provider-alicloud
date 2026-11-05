package alb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudALBAclsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()

	aclIdsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"acl_ids": `["${alicloud_alb_acl.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"acl_ids": `["${alicloud_alb_acl.default.id}_fake"]`,
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_acl.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_alb_acl.default.id}_fakeid"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_acl.default.acl_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_alb_acl.default.acl_name}_fake"`,
		}),
	}

	aclNamesConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"acl_name": `"${alicloud_alb_acl.default.acl_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"acl_name": `"${alicloud_alb_acl.default.acl_name}_fake"`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_acl.default.id}"]`,
			"status": `"Available"`,
		}),
		FakeConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_alb_acl.default.id}"]`,
			"status": `"Configuring"`,
		}),
	}

	ResourceGroupIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"name_regex":        `"${alicloud_alb_acl.default.acl_name}"`,
			"resource_group_id": `"${alicloud_alb_acl.default.resource_group_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"name_regex":        `"${alicloud_alb_acl.default.acl_name}"`,
			"resource_group_id": `"${alicloud_alb_acl.default.resource_group_id}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_alb_acl.default.id}"]`,
			"acl_ids":           `["${alicloud_alb_acl.default.id}"]`,
			"acl_name":          `"${alicloud_alb_acl.default.acl_name}"`,
			"name_regex":        `"${alicloud_alb_acl.default.acl_name}"`,
			"resource_group_id": `"${alicloud_alb_acl.default.resource_group_id}"`,
			"status":            `"Available"`,
		}),
		FakeConfig: testAccCheckAlicloudAlbAclDataSourceName(rand, map[string]string{
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
			"acls.0.id":                        acctest.CHECKSET,
			"acls.0.acl_entries.#":             "1",
			"acls.0.acl_entries.0.description": "description",
			"acls.0.acl_entries.0.entry":       "10.0.0.0/24",
			"acls.0.acl_entries.0.status":      acctest.CHECKSET,
			"acls.0.acl_id":                    acctest.CHECKSET,
			"acls.0.address_ip_version":        acctest.CHECKSET,
			"acls.0.resource_group_id":         acctest.CHECKSET,
		}
	}
	var fakeDataAlicloudAlbAclsSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"acls.#":  "0",
		}
	}
	var alicloudAlbAclCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_alb_acls.default",
		ExistMapFunc: existDataAlicloudAlbAclsSourceNameMapFunc,
		FakeMapFunc:  fakeDataAlicloudAlbAclsSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
	}
	alicloudAlbAclCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, aclIdsConf, idsConf, nameRegexConf, aclNamesConf, statusConf, ResourceGroupIdConf, allConf)
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
