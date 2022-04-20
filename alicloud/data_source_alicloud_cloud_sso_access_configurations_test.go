package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudSSOAccessConfigurationsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	checkoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudSsoAccessConfigurationDataSourceName(rand, map[string]string{
			"enable_details": "true",
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_access_configuration.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCloudSsoAccessConfigurationDataSourceName(rand, map[string]string{
			"enable_details": "true",
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_access_configuration.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudSsoAccessConfigurationDataSourceName(rand, map[string]string{
			"enable_details": "true",
			"directory_id":   `"${local.directory_id}"`,
			"name_regex":     `"${alicloud_cloud_sso_access_configuration.default.access_configuration_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudSsoAccessConfigurationDataSourceName(rand, map[string]string{
			"enable_details": "true",
			"directory_id":   `"${local.directory_id}"`,
			"name_regex":     `"${alicloud_cloud_sso_access_configuration.default.access_configuration_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudSsoAccessConfigurationDataSourceName(rand, map[string]string{
			"enable_details": "true",
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_access_configuration.default.id}"]`,
			"name_regex":     `"${alicloud_cloud_sso_access_configuration.default.access_configuration_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCloudSsoAccessConfigurationDataSourceName(rand, map[string]string{
			"enable_details": "true",
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_access_configuration.default.id}_fake"]`,
			"name_regex":     `"${alicloud_cloud_sso_access_configuration.default.access_configuration_name}_fake"`,
		}),
	}
	var existAlicloudCloudSsoAccessConfigurationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            "1",
			"names.#":          "1",
			"configurations.#": "1",
			"configurations.0.access_configuration_name": fmt.Sprintf("tf-testaccconfiguration%d", rand),
		}
	}
	var fakeAlicloudCloudSsoAccessConfigurationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudCloudSsoAccessConfigurationCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cloud_sso_access_configurations.default",
		existMapFunc: existAlicloudCloudSsoAccessConfigurationsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCloudSsoAccessConfigurationsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckEnterpriseAccountEnabled(t)
	}
	alicloudCloudSsoAccessConfigurationCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudCloudSsoAccessConfigurationDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testaccconfiguration%d"
}

data "alicloud_cloud_sso_directories" "default" {}
resource "alicloud_cloud_sso_directory" "default" {
  count             = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  directory_name    = var.name
}
locals{
  directory_id =  length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
}

resource "alicloud_cloud_sso_access_configuration" "default" {	
	access_configuration_name = var.name
    directory_id = local.directory_id
    permission_policies {
		permission_policy_document = <<EOF
		{
        "Statement":[
        {
        "Action":"cs:Get*",
        "Effect":"Allow",
        "Resource":[
            "*"
        ]
        }
        ],
			"Version": "1"
		}
	  EOF
      permission_policy_type =    "Inline"
	  permission_policy_name =    var.name
	}
  	
}

data "alicloud_cloud_sso_access_configurations" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
