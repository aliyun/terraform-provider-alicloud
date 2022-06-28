package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudSSOUsersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_user.default.id}"]`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id": `"${local.directory_id}"`,
			"ids":          `["${alicloud_cloud_sso_user.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id":   `"${local.directory_id}"`,
			"name_regex":     `"${alicloud_cloud_sso_user.default.user_name}"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id": `"${local.directory_id}"`,
			"name_regex":   `"${alicloud_cloud_sso_user.default.user_name}_fake"`,
		}),
	}
	provisionTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_user.default.id}"]`,
			"provision_type": `"Manual"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_user.default.id}"]`,
			"provision_type": `"Synchronized"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_user.default.id}"]`,
			"status":         `"Enabled"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id": `"${local.directory_id}"`,
			"ids":          `["${alicloud_cloud_sso_user.default.id}"]`,
			"status":       `"Disabled"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_user.default.id}"]`,
			"name_regex":     `"${alicloud_cloud_sso_user.default.user_name}"`,
			"provision_type": `"Manual"`,
			"status":         `"Enabled"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_user.default.id}_fake"]`,
			"name_regex":     `"${alicloud_cloud_sso_user.default.user_name}_fake"`,
			"provision_type": `"Synchronized"`,
			"status":         `"Disabled"`,
		}),
	}
	var existAlicloudCloudSsoUsersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":             "1",
			"names.#":           "1",
			"users.#":           "1",
			"users.0.user_name": fmt.Sprintf("tf-testacccloudssouser%d", rand),
		}
	}
	var fakeAlicloudCloudSsoUsersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCloudSsoUserCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cloud_sso_users.default",
		existMapFunc: existAlicloudCloudSsoUsersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCloudSsoUsersDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckEnterpriseAccountEnabled(t)
	}
	alicloudCloudSsoUserCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, provisionTypeConf, statusConf, allConf)
}
func testAccCheckAlicloudCloudSsoUserDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testacccloudssouser%d"
}

data "alicloud_cloud_sso_directories" "default" {}

resource "alicloud_cloud_sso_directory" "default" {
  count             = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  directory_name    = var.name
}

locals{
  directory_id =  length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
}

resource "alicloud_cloud_sso_user" "default" {	
	user_name = var.name
	directory_id = local.directory_id
}

data "alicloud_cloud_sso_users" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
