package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudResourceManagerSharedTarget_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_shared_target.default"
	ra := resourceAttrInit(resourceId, AlicloudResourceManagerSharedTargetMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourcesharingService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerSharedTarget")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccResourceManagerSharedTarget%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudResourceManagerSharedTargetBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_share_id": "${alicloud_resource_manager_resource_share.default.id}",
					"target_id":         "${data.alicloud_resource_manager_accounts.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_share_id": CHECKSET,
						"target_id":         CHECKSET,
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

var AlicloudResourceManagerSharedTargetMap = map[string]string{
	"status": "Associated",
}

func AlicloudResourceManagerSharedTargetBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_resource_manager_accounts" "default" {}

resource "alicloud_resource_manager_resource_share" "default" {
	resource_share_name = var.name
}

`, name)
}
