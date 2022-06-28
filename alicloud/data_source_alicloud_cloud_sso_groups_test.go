package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudSSOGroupsDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_cloud_sso_groups.default"
	name := fmt.Sprintf("tf-testacc-cloudssogroups%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceCloudSsoGroupsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"directory_id": "${local.directory_id}",
			"name_regex":   "${alicloud_cloud_sso_group.default.group_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"directory_id": "${local.directory_id}",
			"name_regex":   "fake_tf-testacc*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"directory_id": "${local.directory_id}",
			"ids":          []string{"${alicloud_cloud_sso_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"directory_id": "${local.directory_id}",
			"ids":          []string{"${alicloud_cloud_sso_group.default.id}_fake"},
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_cloud_sso_group.default.id}"},
			"directory_id":   "${local.directory_id}",
			"provision_type": "Manual",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"directory_id":   "${local.directory_id}",
			"ids":            []string{"${alicloud_cloud_sso_group.default.id}_fake"},
			"provision_type": "Synchronized",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_cloud_sso_group.default.id}"},
			"name_regex":     "${alicloud_cloud_sso_group.default.group_name}",
			"directory_id":   "${local.directory_id}",
			"provision_type": "Manual",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"directory_id":   "${local.directory_id}",
			"ids":            []string{"${alicloud_cloud_sso_group.default.id}_fake"},
			"provision_type": "Synchronized",
			"name_regex":     "fake_tf-testacc*",
		}),
	}
	var existCloudSsoGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":               "1",
			"names.#":             "1",
			"groups.#":            "1",
			"groups.0.group_name": name,
		}
	}

	var fakeCloudSsoGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#": "0",
			"names.#":  "0",
			"ids.#":    "0",
		}
	}

	var CloudSsoGroupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCloudSsoGroupsMapFunc,
		fakeMapFunc:  fakeCloudSsoGroupsMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckEnterpriseAccountEnabled(t)
		testAccPreCheckWithRegions(t, true, connectivity.CloudSsoSupportRegions)
	}
	CloudSsoGroupsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, typeConf, allConf)
}

func dataSourceCloudSsoGroupsConfigDependence(name string) string {
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
		resource "alicloud_cloud_sso_group" "default" {
		  directory_id = local.directory_id
		  group_name   = var.name
		  description  = var.name
		}
		`, name)
}
