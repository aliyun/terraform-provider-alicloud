package cloudsso_test

import (
	"fmt"
	"strings"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudCloudSsoAccessAssignmentsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudSsoAccessAssignmentDataSourceName(rand, map[string]string{
			"directory_id": `"${local.directory_id}"`,
			"ids":          `["${alicloud_cloud_sso_access_assignment.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudCloudSsoAccessAssignmentDataSourceName(rand, map[string]string{
			"directory_id": `"${local.directory_id}"`,
			"ids":          `["${alicloud_cloud_sso_access_assignment.default.id}_fake"]`,
		}),
	}
	principalType := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudSsoAccessAssignmentDataSourceName(rand, map[string]string{
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_access_assignment.default.id}"]`,
			"principal_type": `"User"`,
		}),
		FakeConfig: testAccCheckAlicloudCloudSsoAccessAssignmentDataSourceName(rand, map[string]string{
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_access_assignment.default.id}_fake"]`,
			"principal_type": `"Group"`,
		}),
	}
	accessConfigurationId := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudSsoAccessAssignmentDataSourceName(rand, map[string]string{
			"directory_id":            `"${local.directory_id}"`,
			"ids":                     `["${alicloud_cloud_sso_access_assignment.default.id}"]`,
			"access_configuration_id": `"${alicloud_cloud_sso_access_assignment.default.access_configuration_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudCloudSsoAccessAssignmentDataSourceName(rand, map[string]string{
			"directory_id":            `"${local.directory_id}"`,
			"ids":                     `["${alicloud_cloud_sso_access_assignment.default.id}_fake"]`,
			"access_configuration_id": `"${alicloud_cloud_sso_access_assignment.default.access_configuration_id}_fake"`,
		}),
	}
	targetId := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudSsoAccessAssignmentDataSourceName(rand, map[string]string{
			"directory_id": `"${local.directory_id}"`,
			"ids":          `["${alicloud_cloud_sso_access_assignment.default.id}"]`,
			"target_id":    `"${alicloud_cloud_sso_access_assignment.default.target_id}"`,
		}),
		FakeConfig: testAccCheckAlicloudCloudSsoAccessAssignmentDataSourceName(rand, map[string]string{
			"directory_id": `"${local.directory_id}"`,
			"ids":          `["${alicloud_cloud_sso_access_assignment.default.id}_fake"]`,
			"target_id":    `"${alicloud_cloud_sso_access_assignment.default.target_id}_fake"`,
		}),
	}
	targetType := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudSsoAccessAssignmentDataSourceName(rand, map[string]string{
			"directory_id": `"${local.directory_id}"`,
			"ids":          `["${alicloud_cloud_sso_access_assignment.default.id}"]`,
			"target_type":  `"RD-Account"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudSsoAccessAssignmentDataSourceName(rand, map[string]string{
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_access_assignment.default.id}"]`,
			"principal_type": `"User"`,
			"target_id":      `"${alicloud_cloud_sso_access_assignment.default.target_id}"`,
			"target_type":    `"RD-Account"`,
		}),
		FakeConfig: testAccCheckAlicloudCloudSsoAccessAssignmentDataSourceName(rand, map[string]string{
			"directory_id":   `"${local.directory_id}"`,
			"ids":            `["${alicloud_cloud_sso_access_assignment.default.id}_fake"]`,
			"principal_type": `"Group"`,
			"target_id":      `"${alicloud_cloud_sso_access_assignment.default.target_id}_fake"`,
		}),
	}
	var existAlicloudCloudSsoAccessAssignmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "1",
			"assignments.#": "1",
			"assignments.0.access_configuration_name": fmt.Sprintf("tf-testaccconfiguration%d", rand),
			"assignments.0.directory_id":              acctest.CHECKSET,
			"assignments.0.access_configuration_id":   acctest.CHECKSET,
			"assignments.0.principal_id":              acctest.CHECKSET,
			"assignments.0.principal_name":            acctest.CHECKSET,
			"assignments.0.principal_type":            "User",
			"assignments.0.target_name":               acctest.CHECKSET,
			"assignments.0.target_type":               "RD-Account",
			"assignments.0.target_id":                 acctest.CHECKSET,
		}
	}
	var fakeAlicloudCloudSsoAccessAssignmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "0",
			"assignments.#": "0",
		}
	}
	var AlicloudCloudSsoAccessAssignmentCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cloud_sso_access_assignments.default",
		ExistMapFunc: existAlicloudCloudSsoAccessAssignmentsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudCloudSsoAccessAssignmentsDataSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.CheckoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	}
	AlicloudCloudSsoAccessAssignmentCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, principalType, accessConfigurationId, targetId, targetType, allConf)
}
func testAccCheckAlicloudCloudSsoAccessAssignmentDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testaccconfiguration%d"
}


data "alicloud_resource_manager_resource_directories" "default" {}

resource "alicloud_cloud_sso_directory" "default" {
  directory_name = var.name
}

resource "alicloud_cloud_sso_user" "default" {
  directory_id = alicloud_cloud_sso_directory.default.id
  user_name    = var.name
}

resource "alicloud_cloud_sso_access_configuration" "default" {
  access_configuration_name = var.name
  directory_id = alicloud_cloud_sso_directory.default.id
}

locals{
  directory_id = alicloud_cloud_sso_directory.default.id
}

resource "alicloud_cloud_sso_access_configuration_provisioning" "default" {
  directory_id = alicloud_cloud_sso_directory.default.id
  access_configuration_id = alicloud_cloud_sso_access_configuration.default.access_configuration_id
  target_type             = "RD-Account"
  target_id               = data.alicloud_resource_manager_resource_directories.default.directories.0.master_account_id
}

resource "alicloud_cloud_sso_access_assignment" "default" {
  directory_id = alicloud_cloud_sso_directory.default.id
  access_configuration_id = alicloud_cloud_sso_access_configuration.default.access_configuration_id
  target_type             = "RD-Account"
  target_id               = data.alicloud_resource_manager_resource_directories.default.directories.0.master_account_id
  principal_type          = "User"
  principal_id            = alicloud_cloud_sso_user.default.user_id
  deprovision_strategy    = "DeprovisionForLastAccessAssignmentOnAccount"
  depends_on              = [alicloud_cloud_sso_access_configuration_provisioning.default]
}

data "alicloud_cloud_sso_access_assignments" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
