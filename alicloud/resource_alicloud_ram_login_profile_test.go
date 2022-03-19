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

func TestAccAlicloudRAMLoginProfile_basic(t *testing.T) {
	var v *ram.GetLoginProfileResponse
	resourceId := "alicloud_ram_login_profile.default"
	ra := resourceAttrInit(resourceId, ramLoginProfilMap)
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
		CheckDestroy:  testAccCheckRamLoginProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamLoginProfileCreateConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"user_name": fmt.Sprintf("tf-testAcc%sRamLoginProfileConfig-%d", defaultRegionToTest, rand)}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			{
				Config: testAccRamLoginProfileUserNameConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"user_name": fmt.Sprintf("tf-testAcc%sRamLoginProfileConfig-%d-N", defaultRegionToTest, rand)}),
				),
			},
			{
				Config: testAccRamLoginProfilePasswordConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"password": "Yourpassword_1235"}),
				),
			},
			{
				Config: testAccRamLoginProfilePasswordResetRequiredConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"password_reset_required": "true"}),
				),
			},
			{
				Config: testAccRamLoginProfileMfaBindRequiredConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"mfa_bind_required": "true"}),
				),
			},
			{
				Config: testAccRamLoginProfileAllConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_name":               fmt.Sprintf("tf-testAcc%sRamLoginProfileConfig-%d", defaultRegionToTest, rand),
						"password":                "Yourpassword_1234",
						"password_reset_required": "false",
						"mfa_bind_required":       "false",
					}),
				),
			},
		},
	})
}

var ramLoginProfilMap = map[string]string{
	"user_name":               CHECKSET,
	"password":                CHECKSET,
	"password_reset_required": CHECKSET,
	"mfa_bind_required":       CHECKSET,
}

func testAccRamLoginProfileCreateConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
	  name = "tf-testAcc%sRamLoginProfileConfig-%d"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}

	resource "alicloud_ram_login_profile" "default" {
	  user_name = "${alicloud_ram_user.default.name}"
	  password = "Yourpassword_1234"
	}`, defaultRegionToTest, rand)
}
func testAccRamLoginProfileUserNameConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
	  name = "tf-testAcc%sRamLoginProfileConfig-%d-N"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}

	resource "alicloud_ram_login_profile" "default" {
	  user_name = "${alicloud_ram_user.default.name}"
	  password = "Yourpassword_1234"
	}`, defaultRegionToTest, rand)
}
func testAccRamLoginProfilePasswordConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
	  name = "tf-testAcc%sRamLoginProfileConfig-%d-N"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}

	resource "alicloud_ram_login_profile" "default" {
	  user_name = "${alicloud_ram_user.default.name}"
	  password = "Yourpassword_1235"
	}`, defaultRegionToTest, rand)
}
func testAccRamLoginProfilePasswordResetRequiredConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
	  name = "tf-testAcc%sRamLoginProfileConfig-%d-N"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}

	resource "alicloud_ram_login_profile" "default" {
	  user_name = "${alicloud_ram_user.default.name}"
	  password = "Yourpassword_1235"
	  password_reset_required = "true"
	}`, defaultRegionToTest, rand)
}
func testAccRamLoginProfileMfaBindRequiredConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
	  name = "tf-testAcc%sRamLoginProfileConfig-%d-N"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}

	resource "alicloud_ram_login_profile" "default" {
	  user_name = "${alicloud_ram_user.default.name}"
	  password = "Yourpassword_1235"
	  password_reset_required = "true"
	  mfa_bind_required = "true"
	}`, defaultRegionToTest, rand)
}
func testAccRamLoginProfileAllConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "default" {
	  name = "tf-testAcc%sRamLoginProfileConfig-%d"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}

	resource "alicloud_ram_login_profile" "default" {
	  user_name = "${alicloud_ram_user.default.name}"
	  password = "Yourpassword_1234"
	  password_reset_required = "false"
	  mfa_bind_required = "false"
	}`, defaultRegionToTest, rand)
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

		if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist.User"}) {
			return WrapError(err)
		}
	}
	return nil
}
