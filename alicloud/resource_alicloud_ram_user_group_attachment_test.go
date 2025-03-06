package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ram UserGroupAttachment. >>> Resource test cases, automatically generated.
// Case UserGroupAttachment测试_副本1737429710156 10095
func TestAccAliCloudRamUserGroupAttachment_basic10095(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_user_group_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudRamUserGroupAttachmentMap10095)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamUserGroupAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRamUserGroupAttachmentBasicDependence10095)
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
					"group_name": "${alicloud_ram_group.defaultieyhdn.id}",
					"user_name":  "${alicloud_ram_user.defaultJSblfg.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": CHECKSET,
						"user_name":  CHECKSET,
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

var AlicloudRamUserGroupAttachmentMap10095 = map[string]string{}

func AlicloudRamUserGroupAttachmentBasicDependence10095(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ram_user" "defaultJSblfg" {
  display_name = var.name
  name = var.name
}

resource "alicloud_ram_group" "defaultieyhdn" {
  name = var.name
}


`, name)
}

// Test Ram UserGroupAttachment. <<< Resource test cases, automatically generated.
