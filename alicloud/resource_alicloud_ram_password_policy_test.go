package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ram PasswordPolicy. >>> Resource test cases, automatically generated.
// Case  PasswordPolicy测试 9035
func TestAccAliCloudRamPasswordPolicy_basic9035(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_password_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudRamPasswordPolicyMap9035)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamPasswordPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRamPasswordPolicyBasicDependence9035)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"minimum_password_length":              "8",
					"require_lowercase_characters":         "false",
					"require_numbers":                      "false",
					"max_password_age":                     "0",
					"password_reuse_prevention":            "1",
					"max_login_attemps":                    "1",
					"hard_expiry":                          "false",
					"require_uppercase_characters":         "false",
					"require_symbols":                      "false",
					"password_not_contain_user_name":       "false",
					"minimum_password_different_character": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"minimum_password_length":              "8",
						"require_lowercase_characters":         "false",
						"require_numbers":                      "false",
						"max_password_age":                     "0",
						"password_reuse_prevention":            "1",
						"max_login_attemps":                    "1",
						"hard_expiry":                          "false",
						"require_uppercase_characters":         "false",
						"require_symbols":                      "false",
						"password_not_contain_user_name":       "false",
						"minimum_password_different_character": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"minimum_password_length":              "10",
					"require_lowercase_characters":         "true",
					"require_numbers":                      "true",
					"max_password_age":                     "99",
					"password_reuse_prevention":            "24",
					"max_login_attemps":                    "32",
					"hard_expiry":                          "true",
					"require_uppercase_characters":         "true",
					"require_symbols":                      "true",
					"password_not_contain_user_name":       "true",
					"minimum_password_different_character": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"minimum_password_length":              "10",
						"require_lowercase_characters":         "true",
						"require_numbers":                      "true",
						"max_password_age":                     "99",
						"password_reuse_prevention":            "24",
						"max_login_attemps":                    "32",
						"hard_expiry":                          "true",
						"require_uppercase_characters":         "true",
						"require_symbols":                      "true",
						"password_not_contain_user_name":       "true",
						"minimum_password_different_character": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"minimum_password_length":              "8",
					"require_lowercase_characters":         "false",
					"require_numbers":                      "false",
					"max_password_age":                     "0",
					"password_reuse_prevention":            "0",
					"max_login_attemps":                    "0",
					"hard_expiry":                          "false",
					"require_uppercase_characters":         "false",
					"require_symbols":                      "false",
					"password_not_contain_user_name":       "false",
					"minimum_password_different_character": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"minimum_password_length":              "8",
						"require_lowercase_characters":         "false",
						"require_numbers":                      "false",
						"max_password_age":                     "0",
						"password_reuse_prevention":            "0",
						"max_login_attemps":                    "0",
						"hard_expiry":                          "false",
						"require_uppercase_characters":         "false",
						"require_symbols":                      "false",
						"password_not_contain_user_name":       "false",
						"minimum_password_different_character": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password_reuse_prevention":            "1",
					"max_login_attemps":                    "1",
					"minimum_password_different_character": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password_reuse_prevention":            "1",
						"max_login_attemps":                    "1",
						"minimum_password_different_character": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"minimum_password_length":              "10",
					"require_lowercase_characters":         "true",
					"require_numbers":                      "true",
					"max_password_age":                     "99",
					"password_reuse_prevention":            "24",
					"max_login_attemps":                    "32",
					"hard_expiry":                          "true",
					"require_uppercase_characters":         "true",
					"require_symbols":                      "true",
					"password_not_contain_user_name":       "true",
					"minimum_password_different_character": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"minimum_password_length":              "10",
						"require_lowercase_characters":         "true",
						"require_numbers":                      "true",
						"max_password_age":                     "99",
						"password_reuse_prevention":            "24",
						"max_login_attemps":                    "32",
						"hard_expiry":                          "true",
						"require_uppercase_characters":         "true",
						"require_symbols":                      "true",
						"password_not_contain_user_name":       "true",
						"minimum_password_different_character": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"minimum_password_length":              "8",
					"require_lowercase_characters":         "false",
					"require_numbers":                      "false",
					"max_password_age":                     "0",
					"password_reuse_prevention":            "1",
					"max_login_attemps":                    "1",
					"hard_expiry":                          "false",
					"require_uppercase_characters":         "false",
					"require_symbols":                      "false",
					"password_not_contain_user_name":       "false",
					"minimum_password_different_character": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"minimum_password_length":              "8",
						"require_lowercase_characters":         "false",
						"require_numbers":                      "false",
						"max_password_age":                     "0",
						"password_reuse_prevention":            "1",
						"max_login_attemps":                    "1",
						"hard_expiry":                          "false",
						"require_uppercase_characters":         "false",
						"require_symbols":                      "false",
						"password_not_contain_user_name":       "false",
						"minimum_password_different_character": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudRamPasswordPolicyMap9035 = map[string]string{}

func AlicloudRamPasswordPolicyBasicDependence9035(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Ram PasswordPolicy. <<< Resource test cases, automatically generated.
