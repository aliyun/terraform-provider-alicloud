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
			{
				Config: testAccRamGroupPolicyAttachmentConfig(acctest.RandIntRange(1000000, 99999999)),
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

func TestAccAlicloudRamGroupPolicyAttachment_update(t *testing.T) {
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
			{
				Config: testAccRamGroupPolicyAttachmentConfig(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamPolicyExists(
						"alicloud_ram_policy.policy", &p),
					testAccCheckRamGroupExists(
						"alicloud_ram_group.group", &g),
					testAccCheckRamGroupPolicyAttachmentExists(
						"alicloud_ram_group_policy_attachment.attach", &p, &g),
				),
			},
			{
				Config: testAccRamGroupPolicyAttachmentUpdate(acctest.RandIntRange(1000000, 99999999)),
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
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No Attachment ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateListPoliciesForGroupRequest()
		request.GroupName = group.GroupName

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForGroup(request)
		})
		if err == nil {
			response, _ := raw.(*ram.ListPoliciesForGroupResponse)
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

func testAccCheckRamGroupPolicyAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_group_policy_attachment" {
			continue
		}

		// Try to find the attachment
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateListPoliciesForGroupRequest()
		request.GroupName = rs.Primary.Attributes["group_name"]

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForGroup(request)
		})

		if err != nil && !RamEntityNotExist(err) {
			return WrapError(err)
		}
		response, _ := raw.(*ram.ListPoliciesForGroupResponse)
		if len(response.Policies.Policy) > 0 {
			for _, v := range response.Policies.Policy {
				if v.PolicyName == rs.Primary.Attributes["name"] && v.PolicyType == rs.Primary.Attributes["policy_type"] {
					return WrapError(Error("Error attachment still exist."))
				}
			}
		}
	}
	return nil
}

func testAccRamGroupPolicyAttachmentConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAccRamGroupPolicyAttachmentConfig-%d"
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

	resource "alicloud_ram_group" "group" {
	  name = "${var.name}"
	  comments = "group comments"
	  force=true
	}

	resource "alicloud_ram_group_policy_attachment" "attach" {
	  policy_name = "${alicloud_ram_policy.policy.name}"
	  group_name = "${alicloud_ram_group.group.name}"
	  policy_type = "${alicloud_ram_policy.policy.type}"
	}`, rand)
}

func testAccRamGroupPolicyAttachmentUpdate(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAccRamGroupPolicyAttachmentConfig-%d"
	}
	resource "alicloud_ram_policy" "policy" {
	  name = "${var.name}1"
	  statement = [
	    {
	      effect = "Allow"
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
	  name = "${var.name}"
	  comments = "group comments"
	  force=true
	}

	resource "alicloud_ram_group_policy_attachment" "attach" {
	  policy_name = "${alicloud_ram_policy.policy.name}"
	  group_name = "${alicloud_ram_group.group.name}"
	  policy_type = "${alicloud_ram_policy.policy.type}"
	}`, rand)
}
