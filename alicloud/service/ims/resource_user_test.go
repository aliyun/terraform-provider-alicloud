package ims_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudImsUser(t *testing.T) {
	// Bounded range: display_name is capped at 24 UTF-8 characters, so the
	// full-width RandInt does not fit once the tf- prefix and -upd suffix are on.
	r := tfacctest.RandIntRange(100000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudImsUserConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_ims_user.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_ims_user.default", "user_id"),
					resource.TestCheckResourceAttrSet("alicloud_ims_user.default", "user_principal_name"),
					resource.TestCheckResourceAttr("alicloud_ims_user.default", "user_name", fmt.Sprintf("tf-testacc%d", r)),
					resource.TestCheckResourceAttr("alicloud_ims_user.default", "display_name", fmt.Sprintf("tf-%d", r)),
					resource.TestCheckResourceAttr("alicloud_ims_user.default", "comments", fmt.Sprintf("tf-%d-comments", r)),
					resource.TestCheckResourceAttr("alicloud_ims_user.default", "provision_type", "Manual"),
					resource.TestCheckResourceAttrSet("alicloud_ims_user.default", "create_date"),
				),
			},
			{
				Config: testAccAliCloudImsUserConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_ims_user.default", "display_name", fmt.Sprintf("tf-%d-upd", r)),
					resource.TestCheckResourceAttr("alicloud_ims_user.default", "email", fmt.Sprintf("tf-%d@example.com", r)),
					resource.TestCheckResourceAttr("alicloud_ims_user.default", "comments", fmt.Sprintf("tf-%d-comments-2", r)),
				),
			},
			{
				Config: testAccAliCloudImsUserConfigRename(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_ims_user.default", "user_principal_name"),
					resource.TestCheckResourceAttr("alicloud_ims_user.default", "user_name", fmt.Sprintf("tf-testacc%d-renamed", r)),
					resource.TestCheckResourceAttr("alicloud_ims_user.default", "display_name", fmt.Sprintf("tf-%d-upd", r)),
					resource.TestCheckResourceAttr("alicloud_ims_user.default", "email", fmt.Sprintf("tf-%d@example.com", r)),
				),
			},
			{
				Config: testAccAliCloudImsUserConfigClear(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_ims_user.default", "user_name", fmt.Sprintf("tf-testacc%d-renamed", r)),
					resource.TestCheckResourceAttr("alicloud_ims_user.default", "display_name", fmt.Sprintf("tf-%d-upd", r)),
					resource.TestCheckNoResourceAttr("alicloud_ims_user.default", "email"),
					resource.TestCheckNoResourceAttr("alicloud_ims_user.default", "comments"),
				),
			},
			{
				ResourceName:      "alicloud_ims_user.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAliCloudImsUserConfig(r int) string {
	return fmt.Sprintf(`
data "alicloud_ims_default_domain" "default" {}

resource "alicloud_ims_user" "default" {
  user_principal_name = "tf-testacc%[1]d@${data.alicloud_ims_default_domain.default.default_domain}"
  display_name        = "tf-%[1]d"
  comments            = "tf-%[1]d-comments"
}
`, r)
}

func testAccAliCloudImsUserConfigUpdate(r int) string {
	return fmt.Sprintf(`
data "alicloud_ims_default_domain" "default" {}

resource "alicloud_ims_user" "default" {
  user_principal_name = "tf-testacc%[1]d@${data.alicloud_ims_default_domain.default.default_domain}"
  display_name        = "tf-%[1]d-upd"
  email               = "tf-%[1]d@example.com"
  comments            = "tf-%[1]d-comments-2"
}
`, r)
}

func testAccAliCloudImsUserConfigRename(r int) string {
	return fmt.Sprintf(`
data "alicloud_ims_default_domain" "default" {}

resource "alicloud_ims_user" "default" {
  user_principal_name = "tf-testacc%[1]d-renamed@${data.alicloud_ims_default_domain.default.default_domain}"
  display_name        = "tf-%[1]d-upd"
  email               = "tf-%[1]d@example.com"
  comments            = "tf-%[1]d-comments-2"
}
`, r)
}

func testAccAliCloudImsUserConfigClear(r int) string {
	return fmt.Sprintf(`
data "alicloud_ims_default_domain" "default" {}

resource "alicloud_ims_user" "default" {
  user_principal_name = "tf-testacc%[1]d-renamed@${data.alicloud_ims_default_domain.default.default_domain}"
  display_name        = "tf-%[1]d-upd"
}
`, r)
}
