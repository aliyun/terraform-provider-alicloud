package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"
	"time"

	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_ram_policy", &resource.Sweeper{
		Name: "alicloud_ram_policy",
		F:    testSweepRamPolicies,
	})
}

func testSweepRamPolicies(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	conn := client.(*AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"tftest",
	}

	args := ram.PolicyQueryRequest{}
	resp, err := conn.ramconn.ListPolicies(args)
	if err != nil {
		return fmt.Errorf("Error retrieving Ram policies: %s", err)
	}

	sweeped := false

	for _, v := range resp.Policies.Policy {
		name := v.PolicyName
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Ram policy: %s", name)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Deleting Ram Policy: %s", name)
		req := ram.PolicyRequest{
			PolicyName: name,
		}
		if _, err := conn.ramconn.DeletePolicy(req); err != nil {
			log.Printf("[ERROR] Failed to delete Ram Policy (%s): %s", name, err)
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

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
						"tf-testAccRamPolicyConfig"),
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
  name = "tf-testAccRamPolicyConfig"
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
