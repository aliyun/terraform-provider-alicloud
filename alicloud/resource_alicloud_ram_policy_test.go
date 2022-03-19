package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
	conn, err := client.NewRamClient()
	if err != nil {
		return WrapError(err)
	}
	action := "ListPolicies"
	request := map[string]interface{}{
		"PolicyType": "Custom",
		"MaxItems":   PageSizeLarge,
	}

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}
	var response map[string]interface{}
	sweeped := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_policies", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Policies.Policy", response)

		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Policies.Policy", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := item["PolicyName"]
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name.(string)), strings.ToLower(prefix)) {
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

			action = "DeletePolicy"
			request := map[string]interface{}{
				"PolicyName": name,
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Ram Policy (%s): %s", name, err)
			}
		}
		if !response["IsTruncated"].(bool) {
			break
		}
		request["Marker"] = response["Marker"]
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudRAMPolicy_basic(t *testing.T) {
	var v map[string]interface{}
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
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAcc%sRamPolicyConfig-%d", defaultRegionToTest, rand),
						"policy_name": fmt.Sprintf("tf-testAcc%sRamPolicyConfig-%d", defaultRegionToTest, rand),
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force", "rotate_strategy"},
			},
			{
				Config: testAccRamPolicyNameConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        fmt.Sprintf("tf-testAcc%sRamPolicyConfig-%d-N", defaultRegionToTest, rand),
						"policy_name": fmt.Sprintf("tf-testAcc%sRamPolicyConfig-%d-N", defaultRegionToTest, rand),
					}),
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
						"policy_name": fmt.Sprintf("tf-testAcc%sRamPolicyConfig-%d", defaultRegionToTest, rand),
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

func TestAccAlicloudRAMPolicy_multi(t *testing.T) {
	var v map[string]interface{}
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
	"policy_name":      CHECKSET,
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
	  policy_name = "tf-testAcc%sRamPolicyConfig-%d"
	  policy_document = <<EOF
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
	  policy_name = "tf-testAcc%sRamPolicyConfig-%d-N"
	  policy_document = <<EOF
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
	  policy_name = "tf-testAcc%sRamPolicyConfig-%d-N"
	  policy_document = <<EOF
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
	  policy_name = "tf-testAcc%sRamPolicyConfig-%d-N"
	  policy_document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"oss:ListObjects",
				"oss:ListObjects"
			  ],
			  "Effect": "Allow",
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
	  policy_name = "tf-testAcc%sRamPolicyConfig-%d-${count.index}"
	  policy_document = <<EOF
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

		// Try to find the policy
		conn, err := client.NewRamClient()
		if err != nil {
			return WrapError(err)
		}
		action := "GetPolicy"
		request := map[string]interface{}{
			"PolicyName": rs.Primary.ID,
			"PolicyType": "Custom",
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist.Policy"}) {
			return WrapError(err)
		}
	}
	return nil
}
