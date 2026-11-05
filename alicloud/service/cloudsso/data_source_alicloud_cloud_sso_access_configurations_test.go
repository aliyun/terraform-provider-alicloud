package cloudsso_test

import (
	"fmt"
	"strings"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCloudSSOAccessConfigurationsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudSsoAccessConfigurationDataSourceName(rand, map[string]string{
			"enable_details": "true",
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_access_configuration.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCloudSsoAccessConfigurationDataSourceName(rand, map[string]string{
			"enable_details": "true",
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_access_configuration.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudSsoAccessConfigurationDataSourceName(rand, map[string]string{
			"enable_details": "true",
			"directory_id":   `"${local.directory_id}"`,
			"name_regex":     `"${alicloud_cloud_sso_access_configuration.default.access_configuration_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCloudSsoAccessConfigurationDataSourceName(rand, map[string]string{
			"enable_details": "true",
			"directory_id":   `"${local.directory_id}"`,
			"name_regex":     `"${alicloud_cloud_sso_access_configuration.default.access_configuration_name}_fake"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudSsoAccessConfigurationDataSourceName(rand, map[string]string{
			"enable_details": "true",
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_access_configuration.default.id}"]`,
			"name_regex":     `"${alicloud_cloud_sso_access_configuration.default.access_configuration_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudCloudSsoAccessConfigurationDataSourceName(rand, map[string]string{
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
	var alicloudCloudSsoAccessConfigurationCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cloud_sso_access_configurations.default",
		ExistMapFunc: existAlicloudCloudSsoAccessConfigurationsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCloudSsoAccessConfigurationsDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckEnterpriseAccountEnabled(t)
	}
	alicloudCloudSsoAccessConfigurationCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
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

resource "alicloud_cloud_sso_directory" "default" {
  directory_name    = var.name
}
locals{
  directory_id = alicloud_cloud_sso_directory.default.id
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
