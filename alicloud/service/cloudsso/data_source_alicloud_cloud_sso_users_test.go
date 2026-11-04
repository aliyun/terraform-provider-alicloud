package cloudsso_test

import (
	"fmt"
	"strings"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCloudSSOUsersDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_user.default.id}"]`,
			"enable_details": "true",
		}),
		FakeConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id": `"${local.directory_id}"`,
			"ids":          `["${alicloud_cloud_sso_user.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id":   `"${local.directory_id}"`,
			"name_regex":     `"${alicloud_cloud_sso_user.default.user_name}"`,
			"enable_details": "true",
		}),
		FakeConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id": `"${local.directory_id}"`,
			"name_regex":   `"${alicloud_cloud_sso_user.default.user_name}_fake"`,
		}),
	}
	provisionTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_user.default.id}"]`,
			"provision_type": `"Manual"`,
			"enable_details": "true",
		}),
		FakeConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_user.default.id}"]`,
			"provision_type": `"Synchronized"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_user.default.id}"]`,
			"status":         `"Enabled"`,
			"enable_details": "true",
		}),
		FakeConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id": `"${local.directory_id}"`,
			"ids":          `["${alicloud_cloud_sso_user.default.id}"]`,
			"status":       `"Disabled"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_user.default.id}"]`,
			"name_regex":     `"${alicloud_cloud_sso_user.default.user_name}"`,
			"provision_type": `"Manual"`,
			"status":         `"Enabled"`,
			"enable_details": "true",
		}),
		FakeConfig: testAccCheckAlicloudCloudSsoUserDataSourceName(rand, map[string]string{
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
	var alicloudCloudSsoUserCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cloud_sso_users.default",
		ExistMapFunc: existAlicloudCloudSsoUsersDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCloudSsoUsersDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckEnterpriseAccountEnabled(t)
	}
	alicloudCloudSsoUserCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, provisionTypeConf, statusConf, allConf)
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


resource "alicloud_cloud_sso_directory" "default" {
  directory_name    = var.name
}

locals{
  directory_id = alicloud_cloud_sso_directory.default.id
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
