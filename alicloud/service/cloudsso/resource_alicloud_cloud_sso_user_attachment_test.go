package cloudsso_test

import (
	"fmt"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	cloudsso "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudsso"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test CloudSSO UserAttachment. >>> Resource test cases, automatically generated.
// Case UserAttachment1_online 10386
func TestAccAliCloudCloudSSOUserAttachment_basic10386(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_sso_user_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudSSOUserAttachmentMap10386)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudsso.NewCloudSSOServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudSSOUserAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudsso%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudSSOUserAttachmentBasicDependence10386)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.CloudSsoSupportRegions)
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"directory_id": "${local.directory_id}",
					"group_id":     "${alicloud_cloud_sso_group.default.group_id}",
					"user_id":      "${alicloud_cloud_sso_user.default.user_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"directory_id": acctest.CHECKSET,
						"group_id":     acctest.CHECKSET,
						"user_id":      acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudCloudSSOUserAttachmentMap10386 = map[string]string{}

func AliCloudCloudSSOUserAttachmentBasicDependence10386(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}


	resource "alicloud_cloud_sso_directory" "default" {
  		directory_name = var.name
	}

	resource "alicloud_cloud_sso_group" "default" {
  		directory_id = local.directory_id
  		group_name   = var.name
  		description  = var.name
	}

	resource "alicloud_cloud_sso_user" "default" {
  		directory_id = local.directory_id
  		user_name    = var.name
	}

	locals {
  		directory_id = alicloud_cloud_sso_directory.default.id
	}
`, name)
}

// Test CloudSSO UserAttachment. <<< Resource test cases, automatically generated.
