package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudSSOScimServerDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_cloud_sso_scim_server_credentials.default"
	name := fmt.Sprintf("tf-testacc-cloudssoCredentials%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceCloudSsoScimServerConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"directory_id": "${local.directory_id}",
			"ids":          []string{"${alicloud_cloud_sso_scim_server_credential.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"directory_id": "${local.directory_id}",
			"ids":          []string{"${alicloud_cloud_sso_scim_server_credential.default.id}_fake"},
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alicloud_cloud_sso_scim_server_credential.default.id}"},
			"directory_id": "${local.directory_id}",
			"status":       "Enabled",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"directory_id": "${local.directory_id}",
			"ids":          []string{"${alicloud_cloud_sso_scim_server_credential.default.id}_fake"},
			"status":       "Disabled",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alicloud_cloud_sso_scim_server_credential.default.id}"},
			"directory_id": "${local.directory_id}",
			"status":       "Enabled",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
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

	var CloudSsoScimServerCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCloudSsoScimServerMapFunc,
		fakeMapFunc:  fakeCloudSsoScimServerMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckEnterpriseAccountEnabled(t)
		testAccPreCheckWithRegions(t, true, connectivity.CloudSsoSupportRegions)
	}
	CloudSsoScimServerCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, allConf)
}

func dataSourceCloudSsoScimServerConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}
		data "alicloud_cloud_sso_directories" "default" {}
		resource "alicloud_cloud_sso_directory" "default" {
		  count             = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
		  directory_name    = var.name
		}
		locals{
		  directory_id =  length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
		}
		resource "alicloud_cloud_sso_scim_server_credential" "default" {
		  directory_id = local.directory_id
		}
		`, name)
}
