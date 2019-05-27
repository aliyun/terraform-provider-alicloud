package alicloud

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudRamAccountPasswordPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckAlicloudRamAccountPasswordPolicyDestroy(),
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudRamAccountPasswordPolicyDefault,
				Check:  resource.ComposeTestCheckFunc(testAccCheckAlicloudRamAccountPasswordPolicyIsSet(expectedPasswordPolicyDefault)),
			},
			{
				Config: testAccAlicloudRamAccountPasswordPolicyModified,
				Check:  resource.ComposeTestCheckFunc(testAccCheckAlicloudRamAccountPasswordPolicyIsSet(expectedPasswordPolicyModified)),
			},
		},
	},
	)

}

func testAccCheckAlicloudRamAccountPasswordPolicyDestroy() resource.TestCheckFunc {
	// default configuration is also expected after Terraform destroy action
	return testAccCheckAlicloudRamAccountPasswordPolicyIsSet(expectedPasswordPolicyDefault)
}

func testAccCheckAlicloudRamAccountPasswordPolicyIsSet(expectedPasswordPolicy ram.PasswordPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateGetPasswordPolicyRequest()
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetPasswordPolicy(request)
		})
		if err != nil {
			WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_account_password_policy", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*ram.GetPasswordPolicyResponse)
		currentPasswordPolicy := response.PasswordPolicy

		if currentPasswordPolicy != expectedPasswordPolicy {
			return WrapError(Error("Value set in Alicloud PasswordPolicy:\n %+v\nexpected PasswordPolicy:\n %+v\n",
				currentPasswordPolicy, expectedPasswordPolicy))
		} else {
			return nil
		}
	}
}

var expectedPasswordPolicyDefault = ram.PasswordPolicy{
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

var expectedPasswordPolicyModified = ram.PasswordPolicy{
	MinimumPasswordLength:      9,
	RequireLowercaseCharacters: false,
	RequireUppercaseCharacters: false,
	RequireNumbers:             false,
	RequireSymbols:             false,
	HardExpiry:                 true,
	MaxPasswordAge:             12,
	PasswordReusePrevention:    5,
	MaxLoginAttemps:            3,
}

const testAccAlicloudRamAccountPasswordPolicyDefault = `
resource "alicloud_ram_account_password_policy" "default" {

}
`
const testAccAlicloudRamAccountPasswordPolicyModified = `
resource "alicloud_ram_account_password_policy" "modified" {
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
