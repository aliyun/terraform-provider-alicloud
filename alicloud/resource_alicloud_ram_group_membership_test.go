package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAlicloudRAMGroupMembership_basic(t *testing.T) {
	var v *ram.ListUsersForGroupResponse
	resourceId := "alicloud_ram_group_membership.default"
	ra := resourceAttrInit(resourceId, groupMenbershipMap)
	serviceFunc := func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRamGroupMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamGroupMembershipCreateConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"group_name": fmt.Sprintf("tf-testAcc%sRamGroupMembershipConfig-%d", defaultRegionToTest, rand)}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRamGroupMembershipUserNameConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"user_names.#": "1"}),
				),
			},
			{
				Config: testAccRamGroupMembershipAllConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(groupMenbershipMap),
				),
			},
		},
	})
}

var groupMenbershipMap = map[string]string{
	"user_names.#": "2",
}

func testAccRamGroupMembershipCreateConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%sRamGroupMembershipConfig-%d"
	}
	resource "alicloud_ram_user" "default" {
	  name = "${var.name}"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}

	resource "alicloud_ram_user" "default1" {
	  name = "${var.name}1"
	  display_name = "displayname1"
	  mobile = "86-18888888888"
	  email = "hello.uuuu@aaa.com"
	  comments = "yoyoyo1"
	}

	resource "alicloud_ram_group" "default" {
	  name = "${var.name}"
	  comments = "group comments"
	  force=true
	}

	resource "alicloud_ram_group_membership" "default" {
	  group_name = "${alicloud_ram_group.default.name}"
	  user_names = ["${alicloud_ram_user.default.name}", "${alicloud_ram_user.default1.name}"]
	}`, defaultRegionToTest, rand)
}

func testAccRamGroupMembershipUserNameConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%sRamGroupMembershipConfig-%d"
	}
	resource "alicloud_ram_user" "default" {
	  name = "${var.name}"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}

	resource "alicloud_ram_group" "default" {
	  name = "${var.name}"
	  comments = "group comments"
	  force=true
	}

	resource "alicloud_ram_group_membership" "default" {
	  group_name = "${alicloud_ram_group.default.name}"
	  user_names = ["${alicloud_ram_user.default.name}"]
	}`, defaultRegionToTest, rand)
}

func testAccRamGroupMembershipAllConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%sRamGroupMembershipConfig-%d"
	}
	resource "alicloud_ram_user" "default" {
	  name = "${var.name}"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}
	resource "alicloud_ram_user" "default1" {
	  name = "${var.name}1"
	  display_name = "displayname1"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo1"
	}

	resource "alicloud_ram_group" "default" {
	  name = "${var.name}"
	  comments = "group comments"
	  force=true
	}

	resource "alicloud_ram_group_membership" "default" {
	  group_name = "${alicloud_ram_group.default.name}"
	  user_names = ["${alicloud_ram_user.default.name}", "${alicloud_ram_user.default1.name}"]
	}`, defaultRegionToTest, rand)
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

		if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist.Group"}) {
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
