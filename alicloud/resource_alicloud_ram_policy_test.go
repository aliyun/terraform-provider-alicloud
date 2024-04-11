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

func testAccCheckRamPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	ramService := RamService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_policy" {
			continue
		}

		_, err := ramService.DescribeRamPolicy(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		return WrapError(Error("Error SecurityGroup still exist"))
	}

	return nil
}

func TestAccAliCloudRamPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudRamPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sRamPolicyConfig-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_name":     name,
					"policy_document": `{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"oss:*\",\"Resource\":[\"acs:oss:*:*:myphotos\",\"acs:oss:*:*:myphotos/*\"]}],\"Version\":\"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name":     name,
						"policy_document": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_document": `{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"oss:*\",\"Resource\":[\"acs:oss:*:*:myphotos\"]}],\"Version\":\"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_document": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"force": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"force": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rotate_strategy", "force"},
			},
		},
	})
}

func TestAccAliCloudRamPolicy_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudRamPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sRamPolicyConfig-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_name":     name,
					"policy_document": `{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"oss:*\",\"Resource\":[\"acs:oss:*:*:myphotos\",\"acs:oss:*:*:myphotos/*\"]}],\"Version\":\"1\"}`,
					"description":     name,
					"rotate_strategy": "None",
					"force":           "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name":     name,
						"policy_document": CHECKSET,
						"description":     name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rotate_strategy", "force"},
			},
		},
	})
}

func TestAccAliCloudRamPolicy_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudRamPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sRamPolicyConfig-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":     name,
					"document": `{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"oss:*\",\"Resource\":[\"acs:oss:*:*:myphotos\",\"acs:oss:*:*:myphotos/*\"]}],\"Version\":\"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":     name,
						"document": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"document": `{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"oss:*\",\"Resource\":[\"acs:oss:*:*:myphotos\"]}],\"Version\":\"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"force": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"force": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rotate_strategy", "force"},
			},
		},
	})
}

func TestAccAliCloudRamPolicy_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudRamPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sRamPolicyConfig-%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":            name,
					"document":        `{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"oss:*\",\"Resource\":[\"acs:oss:*:*:myphotos\",\"acs:oss:*:*:myphotos/*\"]}],\"Version\":\"1\"}`,
					"description":     name,
					"rotate_strategy": "None",
					"force":           "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"document":    CHECKSET,
						"description": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rotate_strategy", "force"},
			},
		},
	})
}

func TestAccAliCloudRamPolicy_multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_policy.default.9"
	ra := resourceAttrInit(resourceId, AliCloudRamPolicyMap0)
	serviceFunc := func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sRamPolicyConfig-%d", defaultRegionToTest, rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRamPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamPolicyMultiConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var AliCloudRamPolicyMap0 = map[string]string{
	"type":             CHECKSET,
	"version_id":       CHECKSET,
	"default_version":  CHECKSET,
	"attachment_count": CHECKSET,
}

func AliCloudRamPolicyBasicDependence0(name string) string {
	return ""
}

func testAccRamPolicyMultiConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_ram_policy" "default" {
  		count           = 10
  		policy_name     = "${var.name}-${count.index}"
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
		description     = var.name
  		force           = true
}
`, name)
}
