package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudSSODirectoriesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudSsoDirectoryDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloud_sso_directory.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudSsoDirectoryDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_cloud_sso_directory.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudSsoDirectoryDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_sso_directory.default.directory_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudSsoDirectoryDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_cloud_sso_directory.default.directory_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudSsoDirectoryDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_cloud_sso_directory.default.id}"]`,
			"name_regex": `"${alicloud_cloud_sso_directory.default.directory_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudSsoDirectoryDataSourceName(rand, map[string]string{
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
	var alicloudCloudSsoDirectoryCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cloud_sso_directories.default",
		existMapFunc: existAlicloudCloudSsoDirectoriesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCloudSsoDirectoriesDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckEnterpriseAccountEnabled(t)
		testAccPreCheckWithRegions(t, true, connectivity.CloudSsoSupportRegions)
	}
	alicloudCloudSsoDirectoryCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
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
