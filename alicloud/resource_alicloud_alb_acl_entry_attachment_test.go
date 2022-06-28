package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudAlbAclEntryAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_alb_acl_entry_attachment.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAlbAclEntryAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlbAclEntryAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAlbAclEntryAttachmentBasicDependence0)
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
					"acl_id":      "${alicloud_alb_acl.default.id}",
					"entry":       "10.10.10.0/24",
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl_id":      CHECKSET,
						"entry":       "10.10.10.0/24",
						"description": name,
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

func AlicloudAlbAclEntryAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_alb_acl" "default" {
  acl_name          = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}
`, name)
}
