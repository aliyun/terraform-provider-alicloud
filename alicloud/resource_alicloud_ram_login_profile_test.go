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
			{
				Config: testAccRamLoginProfileConfig(acctest.RandIntRange(1000000, 99999999)),
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
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No LoginProfile ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateGetLoginProfileRequest()
		request.UserName = rs.Primary.Attributes["user_name"]

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetLoginProfile(request)
		})

		if err != nil {
			return WrapError(err)
		}
		response, _ := raw.(*ram.GetLoginProfileResponse)
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
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateGetLoginProfileRequest()
		request.UserName = rs.Primary.Attributes["user_name"]

		_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetLoginProfile(request)
		})

		if err != nil && !RamEntityNotExist(err) {
			return WrapError(err)
		}
	}
	return nil
}

func testAccRamLoginProfileConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	  name = "tf-testAccRamLoginProfileConfig-%d"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}

	resource "alicloud_ram_login_profile" "profile" {
	  user_name = "${alicloud_ram_user.user.name}"
	  password = "World.123456"
	}`, rand)
}
