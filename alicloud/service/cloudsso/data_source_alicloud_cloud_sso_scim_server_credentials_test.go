package cloudsso_test

import (
	"fmt"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCloudSSOScimServerDataSource(t *testing.T) {

	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_cloud_sso_scim_server_credentials.default"
	name := fmt.Sprintf("tf-testacc-cloudssocredentials%v", rand)

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceCloudSsoScimServerConfigDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"directory_id": "${local.directory_id}",
			"ids":          []string{"${alicloud_cloud_sso_scim_server_credential.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"directory_id": "${local.directory_id}",
			"ids":          []string{"${alicloud_cloud_sso_scim_server_credential.default.id}_fake"},
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alicloud_cloud_sso_scim_server_credential.default.id}"},
			"directory_id": "${local.directory_id}",
			"status":       "Enabled",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"directory_id": "${local.directory_id}",
			"ids":          []string{"${alicloud_cloud_sso_scim_server_credential.default.id}_fake"},
			"status":       "Disabled",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alicloud_cloud_sso_scim_server_credential.default.id}"},
			"directory_id": "${local.directory_id}",
			"status":       "Enabled",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"directory_id": "${local.directory_id}",
			"ids":          []string{"${alicloud_cloud_sso_scim_server_credential.default.id}_fake"},
			"status":       "Disabled",
		}),
	}
	var existCloudSsoScimServerMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "1",
			"credentials.#": "1",
		}
	}

	var fakeCloudSsoScimServerMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"credentials.#": "0",
			"ids.#":         "0",
		}
	}

	var CloudSsoScimServerCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existCloudSsoScimServerMapFunc,
		FakeMapFunc:  fakeCloudSsoScimServerMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckEnterpriseAccountEnabled(t)
		acctest.PreCheckWithRegions(t, true, connectivity.CloudSsoSupportRegions)
	}
	CloudSsoScimServerCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, allConf)
}

func dataSourceCloudSsoScimServerConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}
		resource "alicloud_cloud_sso_directory" "default" {
		  directory_name    = var.name
		}
		locals{
		  directory_id = alicloud_cloud_sso_directory.default.id
		}
		resource "alicloud_cloud_sso_scim_server_credential" "default" {
		  directory_id = local.directory_id
		}
		`, name)
}
