package alicloud

import (
	"fmt"
	"testing"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudRamGroupPolicyAttachment_basic(t *testing.T) {
	var p ram.Policy
	var g ram.Group

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_group_policy_attachment.attach",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamGroupPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRamGroupPolicyAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamPolicyExists(
						"alicloud_ram_policy.policy", &p),
					testAccCheckRamGroupExists(
						"alicloud_ram_group.group", &g),
					testAccCheckRamGroupPolicyAttachmentExists(
						"alicloud_ram_group_policy_attachment.attach", &p, &g),
				),
			},
		},
	})

}

func testAccCheckRamGroupPolicyAttachmentExists(n string, policy *ram.Policy, group *ram.Group) resource.TestCheckFunc {
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

		request := ram.GroupQueryRequest{
			GroupName: group.GroupName,
		}

		response, err := conn.ListPoliciesForGroup(request)
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

func testAccCheckRamGroupPolicyAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_group_policy_attachment" {
			continue
		}

		// Try to find the attachment
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		request := ram.GroupQueryRequest{
			GroupName: rs.Primary.Attributes["group_name"],
		}

		response, err := conn.ListPoliciesForGroup(request)

		if err != nil && !RamEntityNotExist(err) {
			return err
		}

		if len(response.Policies.Policy) > 0 {
			for _, v := range response.Policies.Policy {
				if v.PolicyName == rs.Primary.Attributes["name"] && v.PolicyType == rs.Primary.Attributes["policy_type"] {
					return fmt.Errorf("Error attachment still exist.")
				}
			}
		}
	}
	return nil
}

const testAccRamGroupPolicyAttachmentConfig = `
resource "alicloud_ram_policy" "policy" {
  name = "policyname"
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

resource "alicloud_ram_group" "group" {
  name = "groupname"
  comments = "group comments"
  force=true
}

resource "alicloud_ram_group_policy_attachment" "attach" {
  policy_name = "${alicloud_ram_policy.policy.name}"
  group_name = "${alicloud_ram_group.group.name}"
  policy_type = "${alicloud_ram_policy.policy.type}"
}`
