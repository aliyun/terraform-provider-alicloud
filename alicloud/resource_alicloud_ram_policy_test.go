package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"
	"time"

	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_ram_policy", &resource.Sweeper{
		Name: "alicloud_ram_policy",
		F:    testSweepRamPolicies,
		Dependencies: []string{
			"alicloud_ram_user",
			"alicloud_ram_role",
			"alicloud_ram_group",
		},
	})
}

func testSweepRamPolicies(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"tftest",
	}

	request := ram.CreateListPoliciesRequest()
	sweeped := false
	for {
		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPolicies(request)
		})
		if err != nil {
			return WrapError(err)
		}
		resp, _ := raw.(*ram.ListPoliciesResponse)

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
			request := ram.CreateDeletePolicyRequest()
			request.PolicyName = name

			_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
				return ramClient.DeletePolicy(request)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ram Policy (%s): %s", name, err)
			}
		}
		if !resp.IsTruncated {
			break
		}
		request.Marker = resp.Marker
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
			{
				Config: testAccRamPolicyConfig(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamPolicyExists(
						"alicloud_ram_policy.policy", &v),
					resource.TestMatchResourceAttr(
						"alicloud_ram_policy.policy",
						"name",
						regexp.MustCompile("^tf-testAccRamPolicyConfig-*")),
					resource.TestCheckResourceAttr(
						"alicloud_ram_policy.policy",
						"description",
						"this is a policy test"),
				),
			},
		},
	})

}

func TestAccAlicloudRamPolicy_reDocument(t *testing.T) {
	var v ram.Policy
	randInt := acctest.RandIntRange(1000000, 99999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_ram_policy.policy",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRamPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamPolicyConfig_reDocument1(testPolicyTemplate1, randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamPolicyExists("alicloud_ram_policy.policy", &v),
					resource.TestMatchResourceAttr("alicloud_ram_policy.policy", "name", regexp.MustCompile("^tf-testAccRamPolicyConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_policy.policy", "description", "this is a policy test"),
					resource.TestCheckResourceAttrSet("alicloud_ram_policy.policy", "document"),
					resource.TestCheckResourceAttr("alicloud_ram_policy.policy", "force", "true"),
				),
			},
			{
				Config: testAccRamPolicyConfig_reDocument2(testPolicyTemplate2, randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamPolicyExists("alicloud_ram_policy.policy", &v),
					resource.TestMatchResourceAttr("alicloud_ram_policy.policy", "name", regexp.MustCompile("^tf-testAccRamPolicyConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_policy.policy", "description", "this is a policy test"),
					resource.TestCheckResourceAttrSet("alicloud_ram_policy.policy", "document"),
					resource.TestCheckResourceAttr("alicloud_ram_policy.policy", "force", "true"),
				),
			},
		},
	})

}

func TestAccAlicloudRamPolicy_reStatement(t *testing.T) {
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
			{
				Config: testAccRamPolicyConfig(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamPolicyExists("alicloud_ram_policy.policy", &v),
					resource.TestMatchResourceAttr("alicloud_ram_policy.policy", "name", regexp.MustCompile("^tf-testAccRamPolicyConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_policy.policy", "description", "this is a policy test"),
					resource.TestCheckResourceAttr("alicloud_ram_policy.policy", "statement.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ram_policy.policy", "version", "1"),
				),
			},
			{
				Config: testAccRamPolicyConfig_reStatement(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamPolicyExists("alicloud_ram_policy.policy", &v),
					resource.TestMatchResourceAttr("alicloud_ram_policy.policy", "name", regexp.MustCompile("^tf-testAccRamPolicyConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_policy.policy", "description", "this is a policy test"),
					resource.TestCheckResourceAttr("alicloud_ram_policy.policy", "statement.#", "2"),
					resource.TestCheckResourceAttr("alicloud_ram_policy.policy", "version", "1"),
				),
			},
		},
	})

}

func TestAccAlicloudRamPolicy_reVersion(t *testing.T) {
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
			{
				Config: testAccRamPolicyConfig_defaultVersion(acctest.RandIntRange(1000000, 99999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRamPolicyExists("alicloud_ram_policy.policy", &v),
					resource.TestMatchResourceAttr("alicloud_ram_policy.policy", "name", regexp.MustCompile("^tf-testAccRamPolicyConfig-*")),
					resource.TestCheckResourceAttr("alicloud_ram_policy.policy", "description", "this is a policy test"),
					resource.TestCheckResourceAttr("alicloud_ram_policy.policy", "statement.#", "1"),
					resource.TestCheckResourceAttr("alicloud_ram_policy.policy", "version", "1"),
					resource.TestCheckResourceAttr("alicloud_ram_policy.policy", "force", "true"),
				),
			},
		},
	})

}

func testAccCheckRamPolicyExists(n string, policy *ram.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No Policy ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateGetPolicyRequest()
		request.PolicyName = rs.Primary.ID
		request.PolicyType = "Custom"

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetPolicy(request)
		})
		log.Printf("[WARN] Policy id %#v", rs.Primary.ID)

		if err == nil {
			response, _ := raw.(*ram.GetPolicyResponse)
			*policy = response.Policy
			return nil
		}
		return WrapError(err)
	}
}

func testAccCheckRamPolicyDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_policy" {
			continue
		}

		// Try to find the policy
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateGetPolicyRequest()
		request.PolicyName = rs.Primary.ID
		request.PolicyType = "Custom"

		_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.GetPolicy(request)
		})

		if err != nil && !RamEntityNotExist(err) {
			return WrapError(err)
		}
	}
	return nil
}

func testAccRamPolicyConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_policy" "policy" {
	  name = "tf-testAccRamPolicyConfig-%d"
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
	}`, rand)
}

func testAccRamPolicyConfig_reStatement(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_policy" "policy" {
	  name = "tf-testAccRamPolicyConfig-%d"
	  statement = [
	    {
	      effect = "Deny"
	      action = [
		"oss:ListObjects",
		"oss:ListObjects"]
	      resource = [
		"acs:oss:*:*:mybucket",
		"acs:oss:*:*:mybucket/*"]
	    },
		{
		  effect = "Allow"
		  action = ["CreateInstance"]
		  resource = ["*"]
		}]
	  description = "this is a policy test"
	  force = true
	}`, rand)
}

func testAccRamPolicyConfig_defaultVersion(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_policy" "policy" {
	  name = "tf-testAccRamPolicyConfig-%d"
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
	  version = 1
	  force = true
	}`, rand)
}

func testAccRamPolicyConfig_reDocument1(policy string, rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_policy" "policy" {
	  name = "tf-testAccRamPolicyConfig-%d"
	  document = <<EOF
      %s
      EOF
	  description = "this is a policy test"
	  force = true
	}`, rand, policy)
}

func testAccRamPolicyConfig_reDocument2(policy string, rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_policy" "policy" {
	  name = "tf-testAccRamPolicyConfig-%d"
	  document = <<DEFINITION
      %s
      DEFINITION
	  description = "this is a policy test"
	  force = true
	}`, rand, policy)
}

var testPolicyTemplate1 = `
    {
      "Version": "1",
      "Statement": [
        {
          "Action": ["CreateInstance"],
          "Resource": "*",
          "Effect": "Allow"
        }
      ]
    }
`

var testPolicyTemplate2 = `
    {
      "Version": "1",
      "Statement": [
        {
          "Action": ["CreateInstance"],
          "Resource": "*",
          "Effect": "Deny"
        }
      ]
    }
`
