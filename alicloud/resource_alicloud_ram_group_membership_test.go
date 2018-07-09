package alicloud

import (
	"fmt"
	"testing"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudRamGroupMembership_basic(t *testing.T) {
	var u, u1 ram.User
	var g ram.Group

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_group_membership.membership",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamGroupMembershipDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRamGroupMembershipConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists(
						"alicloud_ram_user.user", &u),
					testAccCheckRamUserExists(
						"alicloud_ram_user.user1", &u1),
					testAccCheckRamGroupExists(
						"alicloud_ram_group.group", &g),
					testAccCheckRamGroupMembershipExists(
						"alicloud_ram_group_membership.membership", &u, &u1, &g),
				),
			},
		},
	})

}

func testAccCheckRamGroupMembershipExists(n string, user *ram.User, user1 *ram.User, group *ram.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No membership ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		request := ram.GroupQueryRequest{
			GroupName: rs.Primary.ID,
		}

		response, err := conn.ListUsersForGroup(request)

		if err == nil {
			if len(response.Users.User) == 2 {
				return nil
			}
			return fmt.Errorf("Membership %s not found.", rs.Primary.ID)
		}
		return fmt.Errorf("Error finding membership %s: %#v", rs.Primary.ID, err)
	}
}

func testAccCheckRamGroupMembershipDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_group_membership" {
			continue
		}

		// Try to find the membership
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		request := ram.GroupQueryRequest{
			GroupName: rs.Primary.ID,
		}

		response, err := conn.ListUsersForGroup(request)

		if err != nil && !RamEntityNotExist(err) {
			return err
		}

		if len(response.Users.User) > 0 {
			for _, v := range response.Users.User {
				for _, u := range rs.Primary.Meta["user_names"].([]string) {
					if v.UserName == u {
						return fmt.Errorf("Error membership still exist.")
					}
				}
			}
		}
	}
	return nil
}

const testAccRamGroupMembershipConfig = `
resource "alicloud_ram_user" "user" {
  name = "username"
  display_name = "displayname"
  mobile = "86-18888888888"
  email = "hello.uuu@aaa.com"
  comments = "yoyoyo"
}

resource "alicloud_ram_user" "user1" {
  name = "username1"
  display_name = "displayname1"
  mobile = "86-18888888888"
  email = "hello.uuuu@aaa.com"
  comments = "yoyoyo1"
}

resource "alicloud_ram_group" "group" {
  name = "groupname"
  comments = "group comments"
  force=true
}

resource "alicloud_ram_group_membership" "membership" {
  group_name = "${alicloud_ram_group.group.name}"
  user_names = ["${alicloud_ram_user.user.name}", "${alicloud_ram_user.user1.name}"]
}`
