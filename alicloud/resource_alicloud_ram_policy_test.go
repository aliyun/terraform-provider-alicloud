package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudRamPolicy_basic(t *testing.T) {
	var v ram.Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_policy.policy",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRamPolicyConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamPolicyExists(
						"alicloud_ram_policy.policy", &v),
					resource.TestCheckResourceAttr(
						"alicloud_ram_policy.policy",
						"name",
						"policyname"),
					resource.TestCheckResourceAttr(
						"alicloud_ram_policy.policy",
						"description",
						"this is a policy test"),
				),
			},
		},
	})

}

func testAccCheckRamPolicyExists(n string, policy *ram.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Policy ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		request := ram.PolicyRequest{
			PolicyName: rs.Primary.ID,
			PolicyType: ram.Custom,
		}

		response, err := conn.GetPolicy(request)
		log.Printf("[WARN] Policy id %#v", rs.Primary.ID)

		if err == nil {
			*policy = response.Policy
			return nil
		}
		return fmt.Errorf("Error finding policy %#v", rs.Primary.ID)
	}
}

func testAccCheckRamPolicyDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_policy" {
			continue
		}

		// Try to find the policy
		client := testAccProvider.Meta().(*AliyunClient)
		conn := client.ramconn

		request := ram.PolicyRequest{
			PolicyName: rs.Primary.ID,
			PolicyType: ram.Custom,
		}

		_, err := conn.GetPolicy(request)

		if err != nil && !RamEntityNotExist(err) {
			return err
		}
	}
	return nil
}

const testAccRamPolicyConfig = `
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
}`
