// Package quickbi_test.
package quickbi_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudQuickBIUsersDataSource(t *testing.T) {
	t.Skip()
	rand := tfacctest.RandInt()

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudQuickBIUserDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_quick_bi_user.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudQuickBIUserDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_quick_bi_user.default.id}_fakeid"]`,
		}),
	}

	keywordConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudQuickBIUserDataSourceName(rand, map[string]string{
			"keyword": `"${alicloud_quick_bi_user.default.nick_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudQuickBIUserDataSourceName(rand, map[string]string{
			"keyword": `"${alicloud_quick_bi_user.default.nick_name}_fake"`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudQuickBIUserDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_quick_bi_user.default.id}"]`,
			"keyword": `"${alicloud_quick_bi_user.default.nick_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudQuickBIUserDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_quick_bi_user.default.id}_fake"]`,
			"keyword": `"${alicloud_quick_bi_user.default.nick_name}"`,
		}),
	}

	var existDataAlicloudQuickBIUsersSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"users.#":                 "1",
			"users.0.nick_name":       fmt.Sprintf("tf-testAccQuickBIUser%d", rand),
			"users.0.admin_user":      "false",
			"users.0.auth_admin_user": "false",
			"users.0.user_type":       "Developer",
		}
	}
	var fakeDataAlicloudQuickBIUsersSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"users.#": "0",
		}
	}
	var alicloudQuickBIUserCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_quick_bi_users.default",
		ExistMapFunc: existDataAlicloudQuickBIUsersSourceNameMapFunc,
		FakeMapFunc:  fakeDataAlicloudQuickBIUsersSourceNameMapFunc,
	}
	preCheck := func() {
		// QuickBI is not available in the default test region (cn-beijing); the
		// service serves cn-hangzhou (and other supported regions) per the official
		// endpoint metadata. Force a supported region so the acceptance tests
		// exercise a real endpoint instead of NXDOMAIN. Production code is unchanged.
		t.Setenv("ALICLOUD_REGION", "cn-hangzhou")
		acctest.PreCheck(t)
	}
	alicloudQuickBIUserCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, keywordConf, allConf)
}
func testAccCheckAlicloudQuickBIUserDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tf-testAccQuickBIUser%d"
}
data "alicloud_ram_users" "default" {
  name_regex  = "terraform*"
}
resource "alicloud_quick_bi_user" "default" {
  nick_name       = var.name
  account_name    = join(":",["%s",data.alicloud_ram_users.default.users.0.name])
  admin_user      = "false"
  auth_admin_user = "false"
  user_type       = "Developer"
}

data "alicloud_quick_bi_users" "default" {
	%s
}
`, rand, "openapiautomation@test.aliyunid.com", strings.Join(pairs, " \n "))
	return config
}
