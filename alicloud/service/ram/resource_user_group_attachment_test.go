package ram_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ram"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Test Ram UserGroupAttachment. >>> Resource test cases, automatically generated.
// Case UserGroupAttachment测试_副本1737429710156 10095
func TestAccAliCloudRamUserGroupAttachment_basic10095(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_user_group_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudRamUserGroupAttachmentMap10095)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ram.NewRamServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRamUserGroupAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudRamUserGroupAttachmentBasicDependence10095)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name": "${alicloud_ram_group.defaultieyhdn.id}",
					"user_name":  "${alicloud_ram_user.defaultJSblfg.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": acctest.CHECKSET,
						"user_name":  acctest.CHECKSET,
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
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
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
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	ramServiceV2 := ram.NewRamServiceV2(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_user_group_attachment" {
			continue
		}
		_, err := ramServiceV2.DescribeRamUserGroupAttachment(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return errs.WrapError(err)
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
