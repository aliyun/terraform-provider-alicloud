package dataworks_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dataworks"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const dsgUserGroupResourceId = "alicloud_data_works_dsg_user_group.default"

// TestAccAliCloudDataWorksDsgUserGroup_basic covers create, in-place update of
// owner/accounts and import. A RAM user is created so the owner and accounts
// reference real Alibaba Cloud users in the tenant; otherwise the API rejects
// the request with USERACCOUNT.OWNER.ERROR. user_group_type and project_name are
// write-only (DsgUserGroupQueryList never returns them), so they are ignored
// during import verification and preserved from configured state otherwise.
func TestAccAliCloudDataWorksDsgUserGroup_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-hangzhou"
	}
	name := fmt.Sprintf("tf-testacc%sdsgusergroup%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		IDRefreshName:            dsgUserGroupResourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckDataWorksDsgUserGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataWorksDsgUserGroupBasicConfig(name, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsgUserGroupResourceId, "name", name),
					resource.TestCheckResourceAttrPair(dsgUserGroupResourceId, "owner", "alicloud_ram_user.default", "name"),
					resource.TestCheckResourceAttr(dsgUserGroupResourceId, "user_group_type", "1"),
					resource.TestCheckResourceAttr(dsgUserGroupResourceId, "accounts.#", "1"),
					resource.TestCheckResourceAttrSet(dsgUserGroupResourceId, "id"),
					resource.TestCheckResourceAttrSet(dsgUserGroupResourceId, "gmt_create"),
				),
			},
			{
				// owner and accounts are in-place updatable via DsgUserGroupAddOrUpdate
				// with the existing Id. A second RAM user exercises the change.
				Config: testAccDataWorksDsgUserGroupBasicConfig(name, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dsgUserGroupResourceId, "owner", "alicloud_ram_user.other", "name"),
					resource.TestCheckResourceAttr(dsgUserGroupResourceId, "accounts.#", "1"),
				),
			},
			{
				ResourceName:            dsgUserGroupResourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"user_group_type", "project_name"},
			},
		},
	})
}

func testAccCheckDataWorksDsgUserGroupDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	service := dataworks.NewService(client)
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_data_works_dsg_user_group" {
			continue
		}
		_, err := service.DsgUserGroupDescribe(ctx, rs.Primary.ID, "")
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing dsg user group %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("dsg user group %s still exists", rs.Primary.ID)
	}
	return nil
}

// testAccDataWorksDsgUserGroupBasicConfig creates a RAM user whose name is used
// as owner and the single account. when first is true the primary user owns
// the group; when false a second user takes over to exercise the in-place
// update path.
func testAccDataWorksDsgUserGroupBasicConfig(name string, first bool) string {
	ownerRef := "alicloud_ram_user.default.name"
	if !first {
		ownerRef = "alicloud_ram_user.other.name"
	}
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_ram_user" "default" {
  name        = format("%%s-owner1", var.name)
  display_name = format("%%s-owner1", var.name)
}

resource "alicloud_ram_user" "other" {
  name        = format("%%s-owner2", var.name)
  display_name = format("%%s-owner2", var.name)
}

resource "alicloud_data_works_dsg_user_group" "default" {
  name            = var.name
  owner           = %s
  user_group_type = 1
  accounts        = [%s]
}
`, name, ownerRef, ownerRef)
}
