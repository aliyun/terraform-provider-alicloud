package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
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
		fmt.Sprintf("tf-testAcc%s", region),
		fmt.Sprintf("tf_testAcc%s", region),
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
	var v *ram.GetPolicyResponse
	resourceId := "alicloud_ram_policy.default"
	ra := resourceAttrInit(resourceId, ramPolicyMap)
	serviceFunc := func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRamPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamPolicyCreateConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": fmt.Sprintf("tf-testAcc%sRamPolicyConfig-%d", defaultRegionToTest, rand)}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
			{
				Config: testAccRamPolicyNameConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"name": fmt.Sprintf("tf-testAcc%sRamPolicyConfig-%d-N", defaultRegionToTest, rand)}),
				),
			},
			{
				Config: testAccRamPolicyDescriptionConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"description": "this is a policy description test"}),
				),
			},
			{
				Config: testAccRamPolicyStatementConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccRamPolicyCreateConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAcc%sRamPolicyConfig-%d", defaultRegionToTest, rand),
						"type":        "Custom",
						"description": "this is a policy test",
						"version":     "1",
						"force":       "true",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudRamPolicy_multi(t *testing.T) {
	var v *ram.GetPolicyResponse
	resourceId := "alicloud_ram_policy.default.9"
	ra := resourceAttrInit(resourceId, ramPolicyMap)
	serviceFunc := func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRamPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamPolicyMultiConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var ramPolicyMap = map[string]string{
	"name":             CHECKSET,
	"type":             "Custom",
	"description":      "this is a policy test",
	"version":          "1",
	"attachment_count": CHECKSET,
	"force":            "true",
	"statement.#":      "1",
}

func testAccRamPolicyCreateConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_policy" "default" {
	  name = "tf-testAcc%sRamPolicyConfig-%d"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"oss:ListObjects",
				"oss:ListObjects"
			  ],
			  "Effect": "Deny",
			  "Resource": [
				"acs:oss:*:*:mybucket",
				"acs:oss:*:*:mybucket/*"
			  ]
			}
		  ],
			"Version": "1"
		}
	  EOF
	  description = "this is a policy test"
	  force = true
	}`, defaultRegionToTest, rand)
}

func testAccRamPolicyNameConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_policy" "default" {
	  name = "tf-testAcc%sRamPolicyConfig-%d-N"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"oss:ListObjects",
				"oss:ListObjects"
			  ],
			  "Effect": "Deny",
			  "Resource": [
				"acs:oss:*:*:mybucket",
				"acs:oss:*:*:mybucket/*"
			  ]
			}
		  ],
			"Version": "1"
		}
	  EOF
	  description = "this is a policy test"
	  force = true
	}`, defaultRegionToTest, rand)
}

func testAccRamPolicyDescriptionConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_policy" "default" {
	  name = "tf-testAcc%sRamPolicyConfig-%d-N"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"oss:ListObjects",
				"oss:ListObjects"
			  ],
			  "Effect": "Deny",
			  "Resource": [
				"acs:oss:*:*:mybucket",
				"acs:oss:*:*:mybucket/*"
			  ]
			}
		  ],
			"Version": "1"
		}
	  EOF
	  description = "this is a policy description test"
	  force = true
	}`, defaultRegionToTest, rand)
}
func testAccRamPolicyStatementConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_policy" "default" {
	  name = "tf-testAcc%sRamPolicyConfig-%d-N"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"oss:ListObjects",
				"oss:ListObjects"
			  ],
			  "Effect": "Deny",
			  "Resource": [
				"acs:oss:*:*:mybucket",
				"acs:oss:*:*:mybucket/*"
			  ]
			}
		  ],
			"Version": "1"
		}
	  EOF
	  description = "this is a policy description test"
	  force = true
	}`, defaultRegionToTest, rand)
}

func testAccRamPolicyMultiConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_policy" "default" {
	  name = "tf-testAcc%sRamPolicyConfig-%d-${count.index}"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"oss:ListObjects",
				"oss:ListObjects"
			  ],
			  "Effect": "Deny",
			  "Resource": [
				"acs:oss:*:*:mybucket",
				"acs:oss:*:*:mybucket/*"
			  ]
			}
		  ],
			"Version": "1"
		}
	  EOF
	  description = "this is a policy test"
	  force = true
	  count = 10
	}`, defaultRegionToTest, rand)
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
