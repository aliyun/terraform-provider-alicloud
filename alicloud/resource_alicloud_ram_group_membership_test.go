package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			{
				Config: testAccRamGroupMembershipConfig(acctest.RandIntRange(1000000, 99999999)),
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

func TestAccAlicloudRamGroupMembership_update(t *testing.T) {
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
			{
				Config: testAccRamGroupMembershipConfig(acctest.RandIntRange(1000000, 99999999)),
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
			{
				Config: testAccRamGroupMembershipUpdate(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamUserExists(
						"alicloud_ram_user.user", &u),
					testAccCheckRamGroupExists(
						"alicloud_ram_group.group", &g),
					testAccCheckRamGroupMembershipExistsOne(
						"alicloud_ram_group_membership.membership", &u, &g),
				),
			},
			{
				Config: testAccRamGroupMembershipConfig(acctest.RandIntRange(1000000, 99999999)),
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
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No membership ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateListUsersForGroupRequest()
		request.GroupName = rs.Primary.ID

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListUsersForGroup(request)
		})

		if err == nil {
			response, _ := raw.(*ram.ListUsersForGroupResponse)
			if len(response.Users.User) == 2 {
				return nil
			}
			return WrapError(fmt.Errorf("Membership %s not found.", rs.Primary.ID))
		}
		return WrapError(err)
	}
}

func testAccCheckRamGroupMembershipExistsOne(n string, user *ram.User, group *ram.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No membership ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateListUsersForGroupRequest()
		request.GroupName = rs.Primary.ID

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListUsersForGroup(request)
		})

		if err == nil {
			response, _ := raw.(*ram.ListUsersForGroupResponse)
			if len(response.Users.User) == 1 {
				return nil
			}
			return WrapError(fmt.Errorf("Membership %s not found.", rs.Primary.ID))
		}
		return WrapError(err)
	}
}

func testAccCheckRamGroupMembershipDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_group_membership" {
			continue
		}

		// Try to find the membership
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateListUsersForGroupRequest()
		request.GroupName = rs.Primary.ID

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListUsersForGroup(request)
		})

		if err != nil && !RamEntityNotExist(err) {
			return WrapError(err)
		}
		response, _ := raw.(*ram.ListUsersForGroupResponse)
		if len(response.Users.User) > 0 {
			for _, v := range response.Users.User {
				for _, u := range rs.Primary.Meta["user_names"].([]string) {
					if v.UserName == u {
						return WrapError(Error("Error membership still exist."))
					}
				}
			}
		}
	}
	return nil
}

func testAccRamGroupMembershipConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAccRamGroupMembershipConfig-%d"
	}
	resource "alicloud_ram_user" "user" {
	  name = "${var.name}"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}

	resource "alicloud_ram_user" "user1" {
	  name = "${var.name}1"
	  display_name = "displayname1"
	  mobile = "86-18888888888"
	  email = "hello.uuuu@aaa.com"
	  comments = "yoyoyo1"
	}

	resource "alicloud_ram_group" "group" {
	  name = "${var.name}"
	  comments = "group comments"
	  force=true
	}

	resource "alicloud_ram_group_membership" "membership" {
	  group_name = "${alicloud_ram_group.group.name}"
	  user_names = ["${alicloud_ram_user.user.name}", "${alicloud_ram_user.user1.name}"]
	}`, rand)
}

func testAccRamGroupMembershipUpdate(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAccRamGroupMembershipConfig-%d"
	}
	resource "alicloud_ram_user" "user" {
	  name = "${var.name}"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}

	resource "alicloud_ram_group" "group" {
	  name = "${var.name}"
	  comments = "group comments"
	  force=true
	}

	resource "alicloud_ram_group_membership" "membership" {
	  group_name = "${alicloud_ram_group.group.name}"
	  user_names = ["${alicloud_ram_user.user.name}"]
	}`, rand)
}
