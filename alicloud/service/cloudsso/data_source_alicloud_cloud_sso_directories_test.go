package cloudsso_test

import (
	"fmt"
	"strings"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCloudSSODirectoriesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudSsoDirectoryDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloud_sso_directory.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCloudSsoDirectoryDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloud_sso_directory.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudSsoDirectoryDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_sso_directory.default.directory_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCloudSsoDirectoryDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_sso_directory.default.directory_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudSsoDirectoryDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloud_sso_directory.default.id}"]`,
			"name_regex": `"${alicloud_cloud_sso_directory.default.directory_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCloudSsoDirectoryDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloud_sso_directory.default.id}_fake"]`,
			"name_regex": `"${alicloud_cloud_sso_directory.default.directory_name}_fake"`,
		}),
	}
	var existAlicloudCloudSsoDirectoriesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        "1",
			"names.#":                      "1",
			"directories.#":                "1",
			"directories.0.directory_name": fmt.Sprintf("tf-testacccloudssodirectory%d", rand),
		}
	}
	var fakeAlicloudCloudSsoDirectoriesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCloudSsoDirectoryCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cloud_sso_directories.default",
		ExistMapFunc: existAlicloudCloudSsoDirectoriesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCloudSsoDirectoriesDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckEnterpriseAccountEnabled(t)
		acctest.PreCheckWithRegions(t, true, connectivity.CloudSsoSupportRegions)
	}
	alicloudCloudSsoDirectoryCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudCloudSsoDirectoryDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testacccloudssodirectory%d"
}

resource "alicloud_cloud_sso_directory" "default" {	
	directory_name = var.name
}

data "alicloud_cloud_sso_directories" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
