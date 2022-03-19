package alicloud

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAlicloudRAMAccountPasswordPolicy_basic(t *testing.T) {
	var v *ram.GetPasswordPolicyResponse

	resourceId := "alicloud_ram_account_password_policy.default"
	ra := resourceAttrInit(resourceId, ramAccountPasswordPolicyMap)
	serviceFunc := func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.RamNoSkipRegions)

		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRamAccountPasswordPolicyDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudRamAccountPasswordPolicyMinimumPasswordLength,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"minimum_password_length": "9",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAlicloudRamAccountPasswordPolicyRequireLowercaseCharacters,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"require_lowercase_characters": "false",
					}),
				),
			},
			{
				Config: testAccAlicloudRamAccountPasswordPolicyRequireUppercaseCharacters,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"require_uppercase_characters": "false",
					}),
				),
			},
			{
				Config: testAccAlicloudRamAccountPasswordPolicyRequireNumbers,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"require_numbers": "false",
					}),
				),
			},
			{
				Config: testAccAlicloudRamAccountPasswordPolicyRequireSymbols,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"require_symbols": "false",
					}),
				),
			},
			{
				Config: testAccAlicloudRamAccountPasswordPolicyHardExpiry,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"hard_expiry": "true",
					}),
				),
			},
			{
				Config: testAccAlicloudRamAccountPasswordPolicyMaxPasswordAge,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_password_age": "12",
					}),
				),
			},
			{
				Config: testAccAlicloudRamAccountPasswordPolicyPasswordReusePrevention,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_reuse_prevention": "5",
					}),
				),
			},
			{
				Config: testAccAlicloudRamAccountPasswordPolicyPasswordMaxLoginAttempts,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_login_attempts": "3",
					}),
				),
			},
			{
				Config: testAccAlicloudRamAccountPasswordPolicyPasswordAll,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"minimum_password_length":      "12",
						"require_lowercase_characters": "true",
						"require_uppercase_characters": "true",
						"require_numbers":              "true",
						"require_symbols":              "true",
						"hard_expiry":                  "false",
						"max_password_age":             "0",
						"password_reuse_prevention":    "0",
						"max_login_attempts":           "5",
					}),
				),
			},
		},
	})
}

const testAccAlicloudRamAccountPasswordPolicyMinimumPasswordLength = `
resource "alicloud_ram_account_password_policy" "default" {
	minimum_password_length = 9
}
`
const testAccAlicloudRamAccountPasswordPolicyRequireLowercaseCharacters = `
resource "alicloud_ram_account_password_policy" "default" {
	minimum_password_length = 9
	require_lowercase_characters = false
}
`
const testAccAlicloudRamAccountPasswordPolicyRequireUppercaseCharacters = `
resource "alicloud_ram_account_password_policy" "default" {
	minimum_password_length = 9
	require_lowercase_characters = false
	require_uppercase_characters = false
}
`
const testAccAlicloudRamAccountPasswordPolicyRequireNumbers = `
resource "alicloud_ram_account_password_policy" "default" {
	minimum_password_length = 9
	require_lowercase_characters = false
	require_uppercase_characters = false
	require_numbers = false
}
`
const testAccAlicloudRamAccountPasswordPolicyRequireSymbols = `
resource "alicloud_ram_account_password_policy" "default" {
	minimum_password_length = 9
	require_lowercase_characters = false
	require_uppercase_characters = false
	require_numbers = false
	require_symbols = false
}
`
const testAccAlicloudRamAccountPasswordPolicyHardExpiry = `
resource "alicloud_ram_account_password_policy" "default" {
	minimum_password_length = 9
	require_lowercase_characters = false
	require_uppercase_characters = false
	require_numbers = false
	require_symbols = false
	hard_expiry = true
}
`
const testAccAlicloudRamAccountPasswordPolicyMaxPasswordAge = `
resource "alicloud_ram_account_password_policy" "default" {
	minimum_password_length = 9
	require_lowercase_characters = false
	require_uppercase_characters = false
	require_numbers = false
	require_symbols = false
	hard_expiry = true
	max_password_age = 12
}
`
const testAccAlicloudRamAccountPasswordPolicyPasswordReusePrevention = `
resource "alicloud_ram_account_password_policy" "default" {
	minimum_password_length = 9
	require_lowercase_characters = false
	require_uppercase_characters = false
	require_numbers = false
	require_symbols = false
	hard_expiry = true
	max_password_age = 12
	password_reuse_prevention = 5
}
`
const testAccAlicloudRamAccountPasswordPolicyPasswordMaxLoginAttempts = `
resource "alicloud_ram_account_password_policy" "default" {
	minimum_password_length = 9
	require_lowercase_characters = false
	require_uppercase_characters = false
	require_numbers = false
	require_symbols = false
	hard_expiry = true
	max_password_age = 12
	password_reuse_prevention = 5
	max_login_attempts = 3
}
`
const testAccAlicloudRamAccountPasswordPolicyPasswordAll = `
resource "alicloud_ram_account_password_policy" "default" {
	minimum_password_length = 12
	require_lowercase_characters = true
	require_uppercase_characters = true
	require_numbers = true
	require_symbols = true
	hard_expiry = false
	max_password_age = 0
	password_reuse_prevention = 0
	max_login_attempts = 5
}
`

var ramAccountPasswordPolicyMap = map[string]string{
	"minimum_password_length":      "12",
	"require_lowercase_characters": "true",
	"require_uppercase_characters": "true",
	"require_numbers":              "true",
	"require_symbols":              "true",
	"hard_expiry":                  "false",
	"max_password_age":             "0",
	"password_reuse_prevention":    "0",
	"max_login_attempts":           "5",
}

func testAccCheckRamAccountPasswordPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	ramService := RamService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_account_password_policy" {
			continue
		}

		object, err := ramService.DescribeRamAccountPasswordPolicy(rs.Primary.ID)

		if err != nil {
			return WrapError(err)
		}

		var expectedPasswordPolicyDefault = ram.PasswordPolicyInGetPasswordPolicy{
			MinimumPasswordLength:      12,
			RequireLowercaseCharacters: true,
			RequireUppercaseCharacters: true,
			RequireNumbers:             true,
			RequireSymbols:             true,
			HardExpiry:                 false,
			MaxPasswordAge:             0,
			PasswordReusePrevention:    0,
			MaxLoginAttemps:            5,
		}

		if object.PasswordPolicy != expectedPasswordPolicyDefault {
			return WrapError(Error("Value set in Alicloud PasswordPolicy:\n %+v\nexpected PasswordPolicy:\n %+v\n",
				object.PasswordPolicy, expectedPasswordPolicyDefault))
		} else {
			return nil
		}
	}
	return nil
}
