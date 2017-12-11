package alicloud

import (
	"fmt"
	"testing"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
)

func TestAccAlicloudRamUser_basic(t *testing.T) {
	var v ram.User

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_user.user",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamUserDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRamUserConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists(
						"alicloud_ram_user.user", &v),
					resource.TestCheckResourceAttr(
						"alicloud_ram_user.user",
						"name",
						"username"),
					resource.TestCheckResourceAttr(
						"alicloud_ram_user.user",
						"display_name",
						"displayname"),
					resource.TestCheckResourceAttr(
						"alicloud_ram_user.user",
						"comments",
						"yoyoyo"),
				),
			},
		},
	})

}

func testAccCheckRamUserExists(n string, user *ram.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No User ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		request := ram.UserQueryRequest{
			UserName: rs.Primary.Attributes["user_name"],
		}

		response, err := conn.GetUser(request)
		log.Printf("[WARN] User id %#v", rs.Primary.ID)

		if err == nil {
			*user = response.User
			return nil
		}
		return fmt.Errorf("Error finding user %#v", rs.Primary.ID)
	}
}

func testAccCheckRamUserDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_user" {
			continue
		}

		// Try to find the user
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		request := ram.UserQueryRequest{
			UserName: rs.Primary.Attributes["user_name"],
		}

		_, err := conn.GetUser(request)

		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return err
		}
	}
	return nil
}

const testAccRamUserConfig = `
resource "alicloud_ram_user" "user" {
  name = "username"
  display_name = "displayname"
  mobile = "86-18888888888"
  email = "hello.uuu@aaa.com"
  comments = "yoyoyo"
}`
