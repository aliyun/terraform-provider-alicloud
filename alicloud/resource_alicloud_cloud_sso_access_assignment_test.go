package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudSSOAccessAssignment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_access_assignment.default"
	checkoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudCloudSSOAccessAssignmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudssoService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSsoAccessAssignment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccconfig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudSSOAccessAssignmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEnterpriseAccountEnabled(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"directory_id":            "${local.directory_id}",
					"access_configuration_id": "${alicloud_cloud_sso_access_configuration.default.access_configuration_id}",
					"target_type":             "RD-Account",
					"target_id":               "${data.alicloud_resource_manager_resource_directories.default.directories.0.master_account_id}",
					"principal_type":          "User",
					"principal_id":            "${alicloud_cloud_sso_user.default.user_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id":            CHECKSET,
						"access_configuration_id": CHECKSET,
						"target_type":             "RD-Account",
						"target_id":               CHECKSET,
						"principal_type":          "User",
						"principal_id":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"deprovision_strategy"},
			},
		},
	})
}

var AlicloudCloudSSOAccessAssignmentMap0 = map[string]string{}

func AlicloudCloudSSOAccessAssignmentBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_cloud_sso_directories" "default" {}
data "alicloud_resource_manager_resource_directories" "default" {}
resource "alicloud_cloud_sso_directory" "default" {
  count             = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  directory_name    = var.name
}
locals{
  directory_id = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
}
resource "alicloud_cloud_sso_access_configuration" "default" {
  access_configuration_name = var.name
  directory_id = local.directory_id
}
resource "alicloud_cloud_sso_user" "default" {
  directory_id = local.directory_id
  user_name = var.name
}
`, name)
}
