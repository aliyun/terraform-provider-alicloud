package alicloud

import (
	"fmt"
	"testing"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudRamLoginProfile_basic(t *testing.T) {
	var v ram.LoginProfile
	var u ram.User

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_login_profile.profile",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamLoginProfileDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRamLoginProfileConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists(
						"alicloud_ram_user.user", &u),
					testAccCheckRamLoginProfileExists(
						"alicloud_ram_login_profile.profile", &v),
				),
			},
		},
	})

}

func testAccCheckRamLoginProfileExists(n string, profile *ram.LoginProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LoginProfile ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		request := ram.UserQueryRequest{
			UserName: rs.Primary.Attributes["user_name"],
		}

		response, err := conn.GetLoginProfile(request)

		if err != nil {
			return fmt.Errorf("Error finding login profile %#v", rs.Primary.ID)
		}
		*profile = response.LoginProfile
		return nil
	}
}

func testAccCheckRamLoginProfileDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_login_profile" {
			continue
		}

		// Try to find the login profile
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		request := ram.UserQueryRequest{
			UserName: rs.Primary.Attributes["user_name"],
		}

		_, err := conn.GetLoginProfile(request)

		if err != nil && !RamEntityNotExist(err) {
			return err
		}
	}
	return nil
}

const testAccRamLoginProfileConfig = `
resource "alicloud_ram_user" "user" {
  name = "username"
  display_name = "displayname"
  mobile = "86-18888888888"
  email = "hello.uuu@aaa.com"
  comments = "yoyoyo"
}

resource "alicloud_ram_login_profile" "profile" {
  user_name = "${alicloud_ram_user.user.name}"
  password = "World.123456"
}`
