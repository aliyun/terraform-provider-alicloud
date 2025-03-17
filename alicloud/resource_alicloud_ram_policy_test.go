package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
		response, err = client.RpcPost("Ram", "2015-05-01", action, nil, request, true)
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
			_, err = client.RpcPost("Ram", "2015-05-01", action, nil, request, true)
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

// Test Ram Policy. >>> Resource test cases, automatically generated.
// Case  Policy测试增加tag测试 10003
func TestAccAliCloudRamPolicy_basic10003(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudRamPolicyMap10003)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamPolicyBasicDependence10003)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_name":     name,
					"policy_document": "{\\n    \\\"Version\\\": \\\"1\\\",\\n    \\\"Statement\\\": [\\n        {\\n            \\\"Effect\\\": \\\"Deny\\\",\\n            \\\"Action\\\": \\\"*\\\",\\n            \\\"Resource\\\": \\\"*\\\"\\n        }\\n    ]\\n}",
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
					"policy_document": "{     \\\"Version\\\": \\\"1\\\",     \\\"Statement\\\": [         {             \\\"Effect\\\": \\\"Allow\\\",             \\\"Action\\\": \\\"*\\\",             \\\"Resource\\\": \\\"*\\\"         }     ] }",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_document": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

func TestAccAliCloudRamPolicy_basic10003_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudRamPolicyMap10003)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamPolicyBasicDependence10003)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":     name,
					"policy_name":     name,
					"policy_document": "{\\n    \\\"Version\\\": \\\"1\\\",\\n    \\\"Statement\\\": [\\n        {\\n            \\\"Effect\\\": \\\"Deny\\\",\\n            \\\"Action\\\": \\\"*\\\",\\n            \\\"Resource\\\": \\\"*\\\"\\n        }\\n    ]\\n}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"rotate_strategy": "DeleteOldestNonDefaultVersionWhenLimitExceeded",
					"force":           "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":     name,
						"policy_name":     name,
						"policy_document": CHECKSET,
						"tags.%":          "2",
						"tags.Created":    "TF",
						"tags.For":        "Test",
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

// Case  Policy适配废弃字段name, document
func TestAccAliCloudRamPolicy_basic10005(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudRamPolicyMap10003)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamPolicyBasicDependence10003)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":     name,
					"document": "{\\n    \\\"Version\\\": \\\"1\\\",\\n    \\\"Statement\\\": [\\n        {\\n            \\\"Effect\\\": \\\"Deny\\\",\\n            \\\"Action\\\": \\\"*\\\",\\n            \\\"Resource\\\": \\\"*\\\"\\n        }\\n    ]\\n}",
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
					"document": "{     \\\"Version\\\": \\\"1\\\",     \\\"Statement\\\": [         {             \\\"Effect\\\": \\\"Allow\\\",             \\\"Action\\\": \\\"*\\\",             \\\"Resource\\\": \\\"*\\\"         }     ] }",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

func TestAccAliCloudRamPolicy_basic10005_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudRamPolicyMap10003)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamPolicyBasicDependence10003)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
					"name":        name,
					"document":    "{\\n    \\\"Version\\\": \\\"1\\\",\\n    \\\"Statement\\\": [\\n        {\\n            \\\"Effect\\\": \\\"Deny\\\",\\n            \\\"Action\\\": \\\"*\\\",\\n            \\\"Resource\\\": \\\"*\\\"\\n        }\\n    ]\\n}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"rotate_strategy": "DeleteOldestNonDefaultVersionWhenLimitExceeded",
					"force":           "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  name,
						"name":         name,
						"document":     CHECKSET,
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
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

// Case  Policy适配废弃字段name, version, statement
func TestAccAliCloudRamPolicy_basic10006(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudRamPolicyMap10003)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamPolicyBasicDependence10003)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":    name,
					"version": "1",
					"statement": []map[string]interface{}{
						{
							"effect":   "Deny",
							"action":   []string{"*"},
							"resource": []string{"*"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":        name,
						"version":     "1",
						"statement.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"statement": []map[string]interface{}{
						{
							"effect":   "Allow",
							"action":   []string{"kms:DescribeKey"},
							"resource": []string{"acs:kms:*:*:*"},
						},
					}}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"statement.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

func TestAccAliCloudRamPolicy_basic10006_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudRamPolicyMap10003)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamPolicyBasicDependence10003)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
					"name":        name,
					"version":     "1",
					"statement": []map[string]interface{}{
						{
							"effect":   "Deny",
							"action":   []string{"*"},
							"resource": []string{"*"},
						},
					},
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"rotate_strategy": "DeleteOldestNonDefaultVersionWhenLimitExceeded",
					"force":           "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  name,
						"name":         name,
						"version":      "1",
						"statement.#":  "1",
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
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
	resourceId := "alicloud_ram_policy.default.5"
	ra := resourceAttrInit(resourceId, AliCloudRamPolicyMap10003)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamPolicyBasicDependence10003)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":           "6",
					"description":     name + "-${count.index}",
					"policy_name":     name + "-${count.index}",
					"policy_document": "{\\n    \\\"Version\\\": \\\"1\\\",\\n    \\\"Statement\\\": [\\n        {\\n            \\\"Effect\\\": \\\"Deny\\\",\\n            \\\"Action\\\": \\\"*\\\",\\n            \\\"Resource\\\": \\\"*\\\"\\n        }\\n    ]\\n}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
					"rotate_strategy": "DeleteOldestNonDefaultVersionWhenLimitExceeded",
					"force":           "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":     name + fmt.Sprint(-5),
						"policy_name":     name + fmt.Sprint(-5),
						"policy_document": CHECKSET,
						"tags.%":          "2",
						"tags.Created":    "TF",
						"tags.For":        "Test",
					}),
				),
			},
		},
	})
}

var AliCloudRamPolicyMap10003 = map[string]string{
	"create_time":      CHECKSET,
	"type":             CHECKSET,
	"attachment_count": CHECKSET,
	"version_id":       CHECKSET,
	"default_version":  CHECKSET,
}

func AliCloudRamPolicyBasicDependence10003(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Ram Policy. <<< Resource test cases, automatically generated.
