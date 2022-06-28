package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudSSOUserAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_user_attachment.default"
	checkoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudCloudSSOUserAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudssoService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudSsoUserAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudssouserattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudSSOUserAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEnterpriseAccountEnabled(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"directory_id": "${local.directory_id}",
					"user_id":      "${alicloud_cloud_sso_user.default.user_id}",
					"group_id":     "${alicloud_cloud_sso_group.default.group_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id": CHECKSET,
						"user_id":      CHECKSET,
						"group_id":     CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudCloudSSOUserAttachmentMap0 = map[string]string{}

func AlicloudCloudSSOUserAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_cloud_sso_directories" "default" {}
resource "alicloud_cloud_sso_directory" "default" {
  count             = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? 0 : 1
  directory_name    = var.name
}
locals{
  directory_id = length(data.alicloud_cloud_sso_directories.default.ids) > 0 ? data.alicloud_cloud_sso_directories.default.ids[0] : concat(alicloud_cloud_sso_directory.default.*.id, [""])[0]
}
resource "alicloud_cloud_sso_user" "default" {
  directory_id = local.directory_id
  user_name = var.name
}
resource "alicloud_cloud_sso_group" "default" {
  directory_id = local.directory_id
  group_name   = var.name
  description  = var.name
}
`, name)
}
