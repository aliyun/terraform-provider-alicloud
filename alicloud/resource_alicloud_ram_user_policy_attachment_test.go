package alicloud

import (
	"fmt"
	"testing"

	"strings"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudRamUserPolicyAttachment_basic(t *testing.T) {
	var p ram.Policy
	var u ram.User

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_user_policy_attachment.attach",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamUserPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRamUserPolicyAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamPolicyExists(
						"alicloud_ram_policy.policy", &p),
					testAccCheckRamUserExists(
						"alicloud_ram_user.user", &u),
					testAccCheckRamUserPolicyAttachmentExists(
						"alicloud_ram_user_policy_attachment.attach", &p, &u),
				),
			},
		},
	})

}

func testAccCheckRamUserPolicyAttachmentExists(n string, policy *ram.Policy, user *ram.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Attachment ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn
		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		request := ram.UserQueryRequest{
			UserName: split[0],
		}

		response, err := conn.ListPoliciesForUser(request)
		if err == nil {
			if len(response.Policies.Policy) > 0 {
				for _, v := range response.Policies.Policy {
					if v.PolicyName == policy.PolicyName && v.PolicyType == policy.PolicyType {
						return nil
					}
				}
			}
			return fmt.Errorf("Error finding attach %s", rs.Primary.ID)
		}
		return fmt.Errorf("Error finding attach %s: %#v", rs.Primary.ID, err)
	}
}

func testAccCheckRamUserPolicyAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_user_policy_attachment" {
			continue
		}

		// Try to find the attachment
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		request := ram.UserQueryRequest{
			UserName: rs.Primary.Attributes["user_name"],
		}

		response, err := conn.ListPoliciesForUser(request)

		if err != nil && !RamEntityNotExist(err) {
			return err
		}

		if len(response.Policies.Policy) > 0 {
			for _, v := range response.Policies.Policy {
				if v.PolicyName == rs.Primary.Attributes["policy_name"] && v.PolicyType == rs.Primary.Attributes["policy_type"] {
					return fmt.Errorf("Error attachment still exist.")
				}
			}
		}
	}
	return nil
}

const testAccRamUserPolicyAttachmentConfig = `
variable "name" {
  default = "testAccRamUserPolicyAttachmentConfig"
}
resource "alicloud_ram_policy" "policy" {
  name = "${var.name}"
  statement = [
    {
      effect = "Deny"
      action = [
        "oss:ListObjects",
        "oss:ListObjects"]
      resource = [
        "acs:oss:*:*:mybucket",
        "acs:oss:*:*:mybucket/*"]
    }]
  description = "this is a policy test"
  force = true
}

resource "alicloud_ram_user" "user" {
  name = "${var.name}"
  display_name = "displayname"
  mobile = "86-18888888888"
  email = "hello.uuu@aaa.com"
  comments = "yoyoyo"
}

resource "alicloud_ram_user_policy_attachment" "attach" {
  policy_name = "${alicloud_ram_policy.policy.name}"
  user_name = "${alicloud_ram_user.user.name}"
  policy_type = "${alicloud_ram_policy.policy.type}"
}`
