package alicloud

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_kms_key", &resource.Sweeper{
		Name: "alicloud_kms_key",
		F:    testSweepKmsKey,
		Dependencies: []string{
			"alicloud_kms_alias",
		},
	})
}

func testSweepKmsKey(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	request := map[string]interface{}{
		"PageSize":   PageSizeXLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
		"Filters":    "[ {\"Key\":\"KeyState\", \"Values\":[\"Enabled\",\"Disabled\",\"PendingImport\"]} ]",
	}
	action := "ListKeys"

	var response map[string]interface{}
	for {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kms_key", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Keys.Key", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Keys.Key", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			if !sweepAll() {
				if _, ok := item["Description"]; !ok {
					continue
				}
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(item["Description"].(string)), strings.ToLower(prefix)) {
						skip = false
					}
				}
				if skip {
					log.Printf("[INFO] Skipping Kms Key: %s", item["Description"].(string))
					continue
				}
			}
			action = "SetDeletionProtection"
			request := map[string]interface{}{
				"ProtectedResourceArn":     item["KeyArn"],
				"EnableDeletionProtection": false,
			}
			_, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to cancel Kms Key DeletionProtection %s (%s): %s", item["Description"], item["KeyId"], err)
			}

			action = "ScheduleKeyDeletion"
			request = map[string]interface{}{
				"KeyId":               item["KeyId"],
				"PendingWindowInDays": 7,
			}
			_, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
			if err != nil {
				log.Printf("[ERROR] Failed to delete Kms Key %s (%s): %s", item["Description"], item["KeyId"], err)
			}
		}
		if len(result) < PageSizeXLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return nil
}

// Test Kms Key. >>> Resource test cases, automatically generated.
// Case 全生命周期 8855
func TestAccAliCloudKmsKey_basic8855(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AliCloudKmsKeyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsKey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKmsKeyBasicDependence8855)
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
					"pending_window_in_days": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pending_window_in_days": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "Enabled",
					"rotation_interval":  "605800s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Enabled",
						"rotation_interval":  "605800s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "Disabled",
					"rotation_interval":  REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Disabled",
						"rotation_interval":  REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection":             "Disabled",
					"deletion_protection_description": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection":             "Disabled",
						"deletion_protection_description": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Enabled",
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
				ImportStateVerifyIgnore: []string{"pending_window_in_days", "deletion_window_in_days"},
			},
		},
	})
}

func TestAccAliCloudKmsKey_basic8855_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AliCloudKmsKeyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsKey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKmsKeyBasicDependence8855)
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
					"key_usage":                       "ENCRYPT/DECRYPT",
					"origin":                          "Aliyun_KMS",
					"key_spec":                        "Aliyun_AES_256",
					"protection_level":                "SOFTWARE",
					"automatic_rotation":              "Enabled",
					"rotation_interval":               "605800s",
					"description":                     name,
					"status":                          "Enabled",
					"deletion_protection":             "Enabled",
					"deletion_protection_description": name,
					"pending_window_in_days":          "7",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_usage":                       "ENCRYPT/DECRYPT",
						"origin":                          "Aliyun_KMS",
						"key_spec":                        "Aliyun_AES_256",
						"protection_level":                "SOFTWARE",
						"automatic_rotation":              "Enabled",
						"rotation_interval":               "605800s",
						"description":                     name,
						"status":                          "Enabled",
						"deletion_protection":             "Enabled",
						"deletion_protection_description": name,
						"tags.%":                          "2",
						"tags.Created":                    "TF",
						"tags.For":                        "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection":             "Disabled",
					"deletion_protection_description": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection":             "Disabled",
						"deletion_protection_description": REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pending_window_in_days", "deletion_window_in_days"},
			},
		},
	})
}

// Case 全生命周期dkms_instance_id, policy 8856
func TestAccAliCloudKmsKey_basic8856(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AliCloudKmsKeyMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsKey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKmsKeyBasicDependence8856)
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
					"dkms_instance_id":       "${alicloud_kms_instance.default.id}",
					"pending_window_in_days": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dkms_instance_id":       CHECKSET,
						"pending_window_in_days": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "Enabled",
					"rotation_interval":  "605800s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Enabled",
						"rotation_interval":  "605800s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "Disabled",
					"rotation_interval":  REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Disabled",
						"rotation_interval":  REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection":             "Disabled",
					"deletion_protection_description": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection":             "Disabled",
						"deletion_protection_description": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": `{\"Statement\": [{\"Action\": [\"kms:*\"],\"Effect\": \"Allow\",\"Principal\": {\"RAM\": [\"acs:ram::` + "${data.alicloud_account.default.id}" + `:*\"]},\"Resource\": [\"*\"],\"Sid\": \"Key\"}],\"Version\": \"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": CHECKSET,
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
					"status": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Enabled",
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
				ImportStateVerifyIgnore: []string{"pending_window_in_days", "deletion_window_in_days"},
			},
		},
	})
}

func TestAccAliCloudKmsKey_basic8856_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AliCloudKmsKeyMap1)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsKey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKmsKeyBasicDependence8856)
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
					"key_usage":                       "ENCRYPT/DECRYPT",
					"origin":                          "Aliyun_KMS",
					"key_spec":                        "Aliyun_AES_256",
					"dkms_instance_id":                "${alicloud_kms_instance.default.id}",
					"protection_level":                "SOFTWARE",
					"automatic_rotation":              "Enabled",
					"rotation_interval":               "605800s",
					"policy":                          `{\"Statement\": [{\"Action\": [\"kms:*\"],\"Effect\": \"Allow\",\"Principal\": {\"RAM\": [\"acs:ram::` + "${data.alicloud_account.default.id}" + `:*\"]},\"Resource\": [\"*\"],\"Sid\": \"Key\"}],\"Version\": \"1\"}`,
					"description":                     name,
					"deletion_protection":             "Enabled",
					"deletion_protection_description": name,
					"status":                          "Enabled",
					"pending_window_in_days":          "7",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_usage":                       "ENCRYPT/DECRYPT",
						"origin":                          "Aliyun_KMS",
						"key_spec":                        "Aliyun_AES_256",
						"dkms_instance_id":                CHECKSET,
						"protection_level":                "SOFTWARE",
						"automatic_rotation":              "Enabled",
						"rotation_interval":               "605800s",
						"policy":                          CHECKSET,
						"description":                     name,
						"deletion_protection":             "Enabled",
						"deletion_protection_description": name,
						"status":                          "Enabled",
						"tags.%":                          "2",
						"tags.Created":                    "TF",
						"tags.For":                        "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection":             "Disabled",
					"deletion_protection_description": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection":             "Disabled",
						"deletion_protection_description": REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pending_window_in_days", "deletion_window_in_days"},
			},
		},
	})
}

// Case 全生命周期, 适配废弃字段key_state 8857
func TestAccAliCloudKmsKey_basic8857(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AliCloudKmsKeyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsKey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKmsKeyBasicDependence8855)
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
					"pending_window_in_days": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pending_window_in_days": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "Enabled",
					"rotation_interval":  "605800s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Enabled",
						"rotation_interval":  "605800s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "Disabled",
					"rotation_interval":  REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Disabled",
						"rotation_interval":  REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection":             "Disabled",
					"deletion_protection_description": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection":             "Disabled",
						"deletion_protection_description": REMOVEKEY,
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
					"key_state": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_state": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_state": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_state": "Enabled",
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
				ImportStateVerifyIgnore: []string{"pending_window_in_days", "deletion_window_in_days"},
			},
		},
	})
}

func TestAccAliCloudKmsKey_basic8857_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AliCloudKmsKeyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsKey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKmsKeyBasicDependence8855)
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
					"key_usage":                       "ENCRYPT/DECRYPT",
					"origin":                          "Aliyun_KMS",
					"key_spec":                        "Aliyun_AES_256",
					"protection_level":                "SOFTWARE",
					"automatic_rotation":              "Enabled",
					"rotation_interval":               "605800s",
					"description":                     name,
					"key_state":                       "Enabled",
					"deletion_protection":             "Enabled",
					"deletion_protection_description": name,
					"pending_window_in_days":          "7",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_usage":                       "ENCRYPT/DECRYPT",
						"origin":                          "Aliyun_KMS",
						"key_spec":                        "Aliyun_AES_256",
						"protection_level":                "SOFTWARE",
						"automatic_rotation":              "Enabled",
						"rotation_interval":               "605800s",
						"description":                     name,
						"key_state":                       "Enabled",
						"deletion_protection":             "Enabled",
						"deletion_protection_description": name,
						"tags.%":                          "2",
						"tags.Created":                    "TF",
						"tags.For":                        "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection":             "Disabled",
					"deletion_protection_description": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection":             "Disabled",
						"deletion_protection_description": REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pending_window_in_days", "deletion_window_in_days"},
			},
		},
	})
}

// Case 全生命周期, 适配废弃字段deletion_window_in_days, is_enabled 8858
func TestAccAliCloudKmsKey_basic8858(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AliCloudKmsKeyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsKey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKmsKeyBasicDependence8855)
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
					"deletion_window_in_days": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_window_in_days": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "Enabled",
					"rotation_interval":  "605800s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Enabled",
						"rotation_interval":  "605800s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "Disabled",
					"rotation_interval":  REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Disabled",
						"rotation_interval":  REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection_description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection":             "Disabled",
					"deletion_protection_description": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection":             "Disabled",
						"deletion_protection_description": REMOVEKEY,
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
					"is_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"is_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_enabled": "true",
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
				ImportStateVerifyIgnore: []string{"pending_window_in_days", "deletion_window_in_days"},
			},
		},
	})
}

func TestAccAliCloudKmsKey_basic8858_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, AliCloudKmsKeyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsKey%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudKmsKeyBasicDependence8855)
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
					"key_usage":                       "ENCRYPT/DECRYPT",
					"origin":                          "Aliyun_KMS",
					"key_spec":                        "Aliyun_AES_256",
					"protection_level":                "SOFTWARE",
					"automatic_rotation":              "Enabled",
					"rotation_interval":               "605800s",
					"description":                     name,
					"is_enabled":                      "true",
					"deletion_protection":             "Enabled",
					"deletion_protection_description": name,
					"deletion_window_in_days":         "7",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_usage":                       "ENCRYPT/DECRYPT",
						"origin":                          "Aliyun_KMS",
						"key_spec":                        "Aliyun_AES_256",
						"protection_level":                "SOFTWARE",
						"automatic_rotation":              "Enabled",
						"rotation_interval":               "605800s",
						"description":                     name,
						"is_enabled":                      "true",
						"deletion_protection":             "Enabled",
						"deletion_protection_description": name,
						"tags.%":                          "2",
						"tags.Created":                    "TF",
						"tags.For":                        "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection":             "Disabled",
					"deletion_protection_description": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection":             "Disabled",
						"deletion_protection_description": REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pending_window_in_days", "deletion_window_in_days"},
			},
		},
	})
}

var AliCloudKmsKeyMap0 = map[string]string{
	"key_usage":           CHECKSET,
	"origin":              CHECKSET,
	"key_spec":            CHECKSET,
	"automatic_rotation":  CHECKSET,
	"status":              CHECKSET,
	"arn":                 CHECKSET,
	"primary_key_version": CHECKSET,
	"last_rotation_date":  CHECKSET,
	"creator":             CHECKSET,
	"creation_date":       CHECKSET,
	"key_state":           CHECKSET,
	"is_enabled":          CHECKSET,
}

var AliCloudKmsKeyMap1 = map[string]string{
	"key_usage":           CHECKSET,
	"origin":              CHECKSET,
	"key_spec":            CHECKSET,
	"automatic_rotation":  CHECKSET,
	"policy":              CHECKSET,
	"status":              CHECKSET,
	"arn":                 CHECKSET,
	"primary_key_version": CHECKSET,
	"last_rotation_date":  CHECKSET,
	"creator":             CHECKSET,
	"creation_date":       CHECKSET,
	"key_state":           CHECKSET,
	"is_enabled":          CHECKSET,
}

func AliCloudKmsKeyBasicDependence8855(name string) string {
	return ""
}

func AliCloudKmsKeyBasicDependence8856(name string) string {
	return fmt.Sprintf(`
	variable "name" {
    	default = "%s"
	}

	data "alicloud_account" "default" {
	}

	data "alicloud_zones" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_kms_instance" "default" {
  		product_version = "3"
  		vpc_num         = "1"
  		key_num         = "1000"
  		secret_num      = "0"
  		spec            = "1000"
  		vpc_id          = data.alicloud_vpcs.default.ids.0
  		vswitch_ids = [
    		data.alicloud_vswitches.default.ids.0
  		]
  		zone_ids = [
    		data.alicloud_zones.default.zones.0.id,
    		data.alicloud_zones.default.zones.1.id
  		]
  		timeouts {
    		delete = "60m"
  		}
	}
`, name)
}

func TestUnitAliCloudKmsKey(t *testing.T) {
	p := Provider().(*schema.Provider).ResourcesMap
	d, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, nil)
	dCreate, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, nil)
	dCreate.MarkNewResource()
	for key, value := range map[string]interface{}{
		"description":            "description",
		"key_spec":               "Aliyun_AES_256",
		"protection_level":       "HSM",
		"pending_window_in_days": 7,
		"key_usage":              "key_usage",
		"origin":                 "origin",
		"rotation_interval":      "rotation_interval",
		"automatic_rotation":     "Disabled",
	} {
		err := dCreate.Set(key, value)
		assert.Nil(t, err)
		err = d.Set(key, value)
		assert.Nil(t, err)
	}
	region := os.Getenv("ALICLOUD_REGION")
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Skipf("Skipping the test case with err: %s", err)
		t.Skipped()
	}
	rawClient = rawClient.(*connectivity.AliyunClient)
	ReadMockResponse := map[string]interface{}{
		"KeyMetadata": map[string]interface{}{
			"Arn":                "arn",
			"AutomaticRotation":  "automatic_rotation",
			"Creator":            "creator",
			"CreationDate":       "creation_date",
			"DeleteDate":         "delete_date",
			"Description":        "description",
			"KeySpec":            "key_spec",
			"KeyUsage":           "key_usage",
			"LastRotationDate":   "last_rotation_date",
			"MaterialExpireTime": "material_expire_time",
			"NextRotationDate":   "next_rotation_date",
			"Origin":             "origin",
			"PrimaryKeyVersion":  "primary_key_version",
			"ProtectionLevel":    "protection_level",
			"RotationInterval":   "rotation_interval",
			"KeyState":           "status",
			"KeyId":              "MockKeyId",
		},
	}

	responseMock := map[string]func(errorCode string) (map[string]interface{}, error){
		"RetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"NotFoundError": func(errorCode string) (map[string]interface{}, error) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("alicloud_kms_key", "MockKeyId"))
		},
		"NoRetryError": func(errorCode string) (map[string]interface{}, error) {
			return nil, &tea.SDKError{
				Code:       String(errorCode),
				Data:       String(errorCode),
				Message:    String(errorCode),
				StatusCode: tea.Int(400),
			}
		},
		"CreateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"UpdateNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"DeleteNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"ReadNormal": func(errorCode string) (map[string]interface{}, error) {
			result := ReadMockResponse
			return result, nil
		},
		"UpdateDisableKey": func(errorCode string) (map[string]interface{}, error) {
			result := map[string]interface{}{
				"KeyState": "Disable",
				"KeyId":    "MockKeyId",
			}
			return result, nil
		},
		"UpdateEnableKey": func(errorCode string) (map[string]interface{}, error) {
			result := map[string]interface{}{
				"KeyState": "Enable",
				"KeyId":    "MockKeyId",
			}
			return result, nil
		},
	}
	// Create
	t.Run("CreateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewKmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudKmsKeyCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateAbnormal", func(t *testing.T) {
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudKmsKeyCreate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("CreateNormal", func(t *testing.T) {
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["CreateNormal"]("")
		})
		err := resourceAliCloudKmsKeyCreate(dCreate, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	// Set ID for Update and Delete Method
	d.SetId("MockKeyId")
	// Update
	t.Run("UpdateClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewKmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})

		err := resourceAliCloudKmsKeyUpdate(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateKeyDescriptionAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"description"} {
			switch p["alicloud_kms_key"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudKmsKeyUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateKeyDescriptionNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"description"} {
			switch p["alicloud_kms_key"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudKmsKeyUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("UpdateRotationPolicyAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"automatic_rotation", "rotation_interval"} {
			switch p["alicloud_kms_key"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudKmsKeyUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateRotationPolicyNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"automatic_rotation", "rotation_interval"} {
			switch p["alicloud_kms_key"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: d.Get(key).(string) + "_update"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: strconv.FormatBool(true)})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			case schema.TypeMap:
				diff.SetAttribute("tags.%", &terraform.ResourceAttrDiff{Old: "0", New: "2"})
				diff.SetAttribute("tags.For", &terraform.ResourceAttrDiff{Old: "", New: "Test"})
				diff.SetAttribute("tags.Created", &terraform.ResourceAttrDiff{Old: "", New: "TF"})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		err := resourceAliCloudKmsKeyUpdate(resourceData1, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	t.Run("UpdateDisableKeyAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"status"} {
			switch p["alicloud_kms_key"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: "Enabled"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: "false"})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		patcheDescribeKmsKey := gomonkey.ApplyMethod(reflect.TypeOf(&KmsServiceV2{}), "DescribeKmsKey", func(*KmsService, string) (map[string]interface{}, error) {
			return responseMock["UpdateDisableKey"]("")
		})
		err := resourceAliCloudKmsKeyUpdate(resourceData1, rawClient)
		patches.Reset()
		patcheDescribeKmsKey.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateDisableKeyNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"status"} {
			switch p["alicloud_kms_key"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: "Enabled"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: "false"})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		patcheDescribeKmsKey := gomonkey.ApplyMethod(reflect.TypeOf(&KmsServiceV2{}), "DescribeKmsKey", func(*KmsService, string) (map[string]interface{}, error) {
			return responseMock["UpdateDisableKey"]("")
		})
		err := resourceAliCloudKmsKeyUpdate(resourceData1, rawClient)
		patches.Reset()
		patcheDescribeKmsKey.Reset()
		assert.Nil(t, err)
	})

	t.Run("UpdateEnableKeyAbnormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"status"} {
			switch p["alicloud_kms_key"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: "Disabled"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: "true"})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		patcheDescribeKmsKey := gomonkey.ApplyMethod(reflect.TypeOf(&KmsServiceV2{}), "DescribeKmsKey", func(*KmsService, string) (map[string]interface{}, error) {
			return responseMock["UpdateEnableKey"]("")
		})
		err := resourceAliCloudKmsKeyUpdate(resourceData1, rawClient)
		patches.Reset()
		patcheDescribeKmsKey.Reset()
		assert.NotNil(t, err)
	})

	t.Run("UpdateEnableKeyNormal", func(t *testing.T) {
		diff := terraform.NewInstanceDiff()
		for _, key := range []string{"status"} {
			switch p["alicloud_kms_key"].Schema[key].Type {
			case schema.TypeString:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: d.Get(key).(string), New: "Disabled"})
			case schema.TypeBool:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.FormatBool(d.Get(key).(bool)), New: "true"})
			case schema.TypeInt:
				diff.SetAttribute(key, &terraform.ResourceAttrDiff{Old: strconv.Itoa(d.Get(key).(int)), New: strconv.Itoa(3)})
			}
		}
		resourceData1, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, diff)
		resourceData1.SetId(d.Id())
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["UpdateNormal"]("")
		})
		patcheDescribeKmsKey := gomonkey.ApplyMethod(reflect.TypeOf(&KmsServiceV2{}), "DescribeKmsKey", func(*KmsService, string) (map[string]interface{}, error) {
			return responseMock["UpdateEnableKey"]("")
		})
		err := resourceAliCloudKmsKeyUpdate(resourceData1, rawClient)
		patches.Reset()
		patcheDescribeKmsKey.Reset()
		assert.Nil(t, err)
	})

	// Delete
	t.Run("DeleteClientAbnormal", func(t *testing.T) {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&connectivity.AliyunClient{}), "NewKmsClient", func(_ *connectivity.AliyunClient) (*client.Client, error) {
			return nil, &tea.SDKError{
				Code:       String("loadEndpoint error"),
				Data:       String("loadEndpoint error"),
				Message:    String("loadEndpoint error"),
				StatusCode: tea.Int(400),
			}
		})
		err := resourceAliCloudKmsKeyDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockAbnormal", func(t *testing.T) {
		resourceData, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, nil)
		resourceData.SetId(d.Id())
		_ = resourceData.Set("pending_window_in_days", 7)
		retryFlag := true
		noRetryFlag := true
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudKmsKeyDelete(d, rawClient)
		patches.Reset()
		assert.NotNil(t, err)
	})
	t.Run("DeleteMockNormal", func(t *testing.T) {
		resourceData, _ := schema.InternalMap(p["alicloud_kms_key"].Schema).Data(nil, nil)
		resourceData.SetId(d.Id())
		_ = resourceData.Set("deletion_window_in_days", 7)
		retryFlag := false
		noRetryFlag := false
		patches := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			if retryFlag {
				retryFlag = false
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				noRetryFlag = false
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["DeleteNormal"]("")
		})
		err := resourceAliCloudKmsKeyDelete(resourceData, rawClient)
		patches.Reset()
		assert.Nil(t, err)
	})

	//Read
	t.Run("ReadDescribeKmsKeyNotFound", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			NotFoundFlag := true
			noRetryFlag := false
			if NotFoundFlag {
				return responseMock["NotFoundError"]("ResourceNotfound")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NoRetryError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudKmsKeyRead(d, rawClient)
		patcheDorequest.Reset()
		assert.Nil(t, err)
	})

	t.Run("ReadDescribeKmsKeyAbnormal", func(t *testing.T) {
		patcheDorequest := gomonkey.ApplyMethod(reflect.TypeOf(&client.Client{}), "DoRequest", func(_ *client.Client, _ *string, _ *string, _ *string, _ *string, _ *string, _ map[string]interface{}, _ map[string]interface{}, _ *util.RuntimeOptions) (map[string]interface{}, error) {
			retryFlag := false
			noRetryFlag := true
			if retryFlag {
				return responseMock["RetryError"]("Throttling")
			} else if noRetryFlag {
				return responseMock["NoRetryError"]("NonRetryableError")
			}
			return responseMock["ReadNormal"]("")
		})
		err := resourceAliCloudKmsKeyRead(d, rawClient)
		patcheDorequest.Reset()
		assert.NotNil(t, err)
	})
}

// Test Ecs SecurityGroup. <<< Resource test cases, automatically generated.
