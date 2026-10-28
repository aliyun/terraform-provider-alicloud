package ram_test

import (
	"fmt"
	"testing"

	ramsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ram"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudRAMGroupMembership_basic(t *testing.T) {
	var v *ramsdk.ListUsersForGroupResponse
	resourceId := "alicloud_ram_group_membership.default"
	ra := acctest.ResourceAttrInit(resourceId, groupMenbershipMap)
	serviceFunc := func() interface{} {
		return ram.NewRamService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	rand := tfacctest.RandIntRange(1000000, 9999999)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckRamGroupMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamGroupMembershipCreateConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"group_name": fmt.Sprintf("tf-testAcc%sRamGroupMembershipConfig-%d", acctest.DefaultRegionToTest, rand)}),
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
	  force = true
	}

	resource "alicloud_ram_user" "default1" {
	  name = "${var.name}1"
	  display_name = "displayname1"
	  mobile = "86-18888888888"
	  email = "hello.uuuu@aaa.com"
	  comments = "yoyoyo1"
	  force = true
	}

	resource "alicloud_ram_group" "default" {
	  name = "${var.name}"
	  comments = "group comments"
	  force=true
	}

	resource "alicloud_ram_group_membership" "default" {
	  group_name = "${alicloud_ram_group.default.name}"
	  user_names = ["${alicloud_ram_user.default.name}", "${alicloud_ram_user.default1.name}"]
	}`, acctest.DefaultRegionToTest, rand)
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
	  force = true
	}

	resource "alicloud_ram_user" "default1" {
	  name = "${var.name}1"
	  display_name = "displayname1"
	  mobile = "86-18888888888"
	  email = "hello.uuuu@aaa.com"
	  comments = "yoyoyo1"
	  force = true
	}

	resource "alicloud_ram_group" "default" {
	  name = "${var.name}"
	  comments = "group comments"
	  force=true
	}

	resource "alicloud_ram_group_membership" "default" {
	  group_name = "${alicloud_ram_group.default.name}"
	  user_names = ["${alicloud_ram_user.default.name}"]
	}`, acctest.DefaultRegionToTest, rand)
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
	  force = true
	}
	resource "alicloud_ram_user" "default1" {
	  name = "${var.name}1"
	  display_name = "displayname1"
	  mobile = "86-18888888888"
	  email = "hello.uuuu@aaa.com"
	  comments = "yoyoyo1"
	  force = true
	}

	resource "alicloud_ram_group" "default" {
	  name = "${var.name}"
	  comments = "group comments"
	  force=true
	}

	resource "alicloud_ram_group_membership" "default" {
	  group_name = "${alicloud_ram_group.default.name}"
	  user_names = ["${alicloud_ram_user.default.name}", "${alicloud_ram_user.default1.name}"]
	}`, acctest.DefaultRegionToTest, rand)
}

func testAccCheckRamGroupMembershipDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_group_membership" {
			continue
		}

		// Try to find the membership
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)

		request := ramsdk.CreateListUsersForGroupRequest()
		request.GroupName = rs.Primary.ID
		users := []ramsdk.User{}

		for {
			raw, err := client.WithRamClient(func(ramClient *ramsdk.Client) (interface{}, error) {
				return ramClient.ListUsersForGroup(request)
			})

			if err != nil && !errs.IsExpectedErrors(err, []string{"EntityNotExist.Group"}) {
				return errs.WrapError(err)
			}

			response, _ := raw.(*ramsdk.ListUsersForGroupResponse)
			users = append(users, response.Users.User...)
			if !response.IsTruncated {
				break
			}
			request.Marker = response.Marker
		}
		if len(users) > 0 {
			for _, v := range users {
				for _, u := range rs.Primary.Meta["user_names"].([]string) {
					if v.UserName == u {
						return errs.WrapError(fmt.Errorf("Error membership still exist."))
					}
				}
			}
		}
	}
	return nil
}
