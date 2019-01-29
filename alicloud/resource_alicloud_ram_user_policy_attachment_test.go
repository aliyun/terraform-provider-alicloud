package alicloud

import (
	"fmt"
	"testing"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			{
				Config: testAccRamUserPolicyAttachmentConfig(acctest.RandIntRange(1000000, 99999999)),
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
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No attachment ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		split := strings.Split(rs.Primary.ID, COLON_SEPARATED)
		request := ram.CreateListPoliciesForUserRequest()
		request.UserName = split[0]

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForUser(request)
		})
		if err == nil {
			response, _ := raw.(*ram.ListPoliciesForUserResponse)
			if len(response.Policies.Policy) > 0 {
				for _, v := range response.Policies.Policy {
					if v.PolicyName == policy.PolicyName && v.PolicyType == policy.PolicyType {
						return nil
					}
				}
			}
			return WrapError(fmt.Errorf("Error finding attach %s", rs.Primary.ID))
		}
		return WrapError(err)
	}
}

func testAccCheckRamUserPolicyAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_user_policy_attachment" {
			continue
		}

		// Try to find the attachment
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateListPoliciesForUserRequest()
		request.UserName = rs.Primary.Attributes["user_name"]

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForUser(request)
		})

		if err != nil && !RamEntityNotExist(err) {
			return WrapError(err)
		}
		response, _ := raw.(*ram.ListPoliciesForUserResponse)
		if len(response.Policies.Policy) > 0 {
			for _, v := range response.Policies.Policy {
				if v.PolicyName == rs.Primary.Attributes["policy_name"] && v.PolicyType == rs.Primary.Attributes["policy_type"] {
					return WrapError(Error("Error attachment still exist."))
				}
			}
		}
	}
	return nil
}

func testAccRamUserPolicyAttachmentConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAccRamUserPolicyAttachmentConfig-%d"
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
	}`, rand)
}
