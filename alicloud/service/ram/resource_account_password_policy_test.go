package ram_test

import (
	"fmt"
	"testing"

	ramsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ram"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudRAMAccountPasswordPolicy_basic(t *testing.T) {
	var v *ramsdk.GetPasswordPolicyResponse

	resourceId := "alicloud_ram_account_password_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, ramAccountPasswordPolicyMap)
	serviceFunc := func() interface{} {
		return ram.NewRamService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.RamNoSkipRegions)

		},

		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckRamAccountPasswordPolicyDestroy,

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
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	ramService := ram.NewRamService(client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_account_password_policy" {
			continue
		}

		object, err := ramService.DescribeRamAccountPasswordPolicy(rs.Primary.ID)

		if err != nil {
			return errs.WrapError(err)
		}

		var expectedPasswordPolicyDefault = ramsdk.PasswordPolicy{
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
			return errs.WrapError(fmt.Errorf("Value set in Alicloud PasswordPolicy:\n %+v\nexpected PasswordPolicy:\n %+v\n",
				object.PasswordPolicy, expectedPasswordPolicyDefault))
		} else {
			return nil
		}
	}
	return nil
}
