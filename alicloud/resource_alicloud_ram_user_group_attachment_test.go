package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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

// Case UserGroupAttachment idempotency: adding a user that is already a member of the group
// makes AddUserToGroup return 409 EntityAlreadyExists.User.Group, which Create must accept
// instead of failing.
func TestAccAliCloudRamUserGroupAttachment_alreadyExists(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamUserGroupAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamUserGroupAttachmentAlreadyExistsConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_ram_user_group_attachment.default", "group_name"),
					resource.TestCheckResourceAttrSet("alicloud_ram_user_group_attachment.default", "user_name"),
					// The duplicate attachment only reaches state if Create accepted the 409 and
					// the state sync confirmed the membership.
					resource.TestCheckResourceAttrSet("alicloud_ram_user_group_attachment.duplicate", "group_name"),
					resource.TestCheckResourceAttrSet("alicloud_ram_user_group_attachment.duplicate", "user_name"),
				),
			},
		},
	})
}

func testAccCheckRamUserGroupAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	ramServiceV2 := RamServiceV2{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_user_group_attachment" {
			continue
		}
		_, err := ramServiceV2.DescribeRamUserGroupAttachment(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		return fmt.Errorf("RAM UserGroupAttachment %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccRamUserGroupAttachmentAlreadyExistsConfig(name string) string {
	return fmt.Sprintf(`
resource "alicloud_ram_user" "default" {
  name = "%s"
}

resource "alicloud_ram_group" "default" {
  name = "%s"
}

resource "alicloud_ram_user_group_attachment" "default" {
  group_name = alicloud_ram_group.default.name
  user_name  = alicloud_ram_user.default.name
}

# Same user and group as above: AddUserToGroup returns 409 EntityAlreadyExists.User.Group,
# which Create must accept for this to apply cleanly.
resource "alicloud_ram_user_group_attachment" "duplicate" {
  group_name = alicloud_ram_group.default.name
  user_name  = alicloud_ram_user.default.name
}
`, name, name)
}

var AlicloudRamUserGroupAttachmentMap10095 = map[string]string{}

func AlicloudRamUserGroupAttachmentBasicDependence10095(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_ram_user" "defaultJSblfg" {
  display_name = var.name
  name         = var.name
}

resource "alicloud_ram_group" "defaultieyhdn" {
  name = var.name
}`, name)
}

// Test Ram UserGroupAttachment. <<< Resource test cases, automatically generated.
