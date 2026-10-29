package cloudsso_test

import (
	"fmt"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCloudSSOGroupsDataSource(t *testing.T) {

	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_cloud_sso_groups.default"
	name := fmt.Sprintf("tf-testacc-cloudssogroups%v", rand)

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceCloudSsoGroupsConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"directory_id": "${local.directory_id}",
			"name_regex":   "${alicloud_cloud_sso_group.default.group_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"directory_id": "${local.directory_id}",
			"name_regex":   "fake_tf-testacc*",
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"directory_id": "${local.directory_id}",
			"ids":          []string{"${alicloud_cloud_sso_group.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"directory_id": "${local.directory_id}",
			"ids":          []string{"${alicloud_cloud_sso_group.default.id}_fake"},
		}),
	}

	typeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_cloud_sso_group.default.id}"},
			"directory_id":   "${local.directory_id}",
			"provision_type": "Manual",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"directory_id":   "${local.directory_id}",
			"ids":            []string{"${alicloud_cloud_sso_group.default.id}_fake"},
			"provision_type": "Synchronized",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_cloud_sso_group.default.id}"},
			"name_regex":     "${alicloud_cloud_sso_group.default.group_name}",
			"directory_id":   "${local.directory_id}",
			"provision_type": "Manual",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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

	var CloudSsoGroupsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existCloudSsoGroupsMapFunc,
		FakeMapFunc:  fakeCloudSsoGroupsMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckEnterpriseAccountEnabled(t)
		acctest.PreCheckWithRegions(t, true, connectivity.CloudSsoSupportRegions)
	}
	CloudSsoGroupsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, typeConf, allConf)
}

func dataSourceCloudSsoGroupsConfigDependence(name string) string {
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
		resource "alicloud_cloud_sso_group" "default" {
		  directory_id = local.directory_id
		  group_name   = var.name
		  description  = var.name
		}
		`, name)
}
