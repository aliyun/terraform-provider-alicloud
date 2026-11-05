package kms_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/kms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	request := map[string]interface{}{
		"PageSize":   helper.PageSizeXLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
		"Filters":    "[ {\"Key\":\"KeyState\", \"Values\":[\"Enabled\",\"Disabled\",\"PendingImport\"]} ]",
	}
	action := "ListKeys"

	var response map[string]interface{}
	for {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_kms_key", action, errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(action, response, request)

		resp, err := jsonpath.Get("$.Keys.Key", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Keys.Key", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			if !sweep.SweepAll() {
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
		if len(result) < helper.PageSizeXLarge {
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudKmsKeyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsKey")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsKey%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudKmsKeyBasicDependence8855)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
					"rotation_interval":  acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Disabled",
						"rotation_interval":  acctest.REMOVEKEY,
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
					"deletion_protection_description": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection":             "Disabled",
						"deletion_protection_description": acctest.REMOVEKEY,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudKmsKeyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsKey")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsKey%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudKmsKeyBasicDependence8855)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
					"deletion_protection_description": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection":             "Disabled",
						"deletion_protection_description": acctest.REMOVEKEY,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudKmsKeyMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsKey")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsKey%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudKmsKeyBasicDependence8856)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"}) // fixed to cn-hangzhou region
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dkms_instance_id":       "${alicloud_kms_instance.default.id}",
					"pending_window_in_days": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dkms_instance_id":       acctest.CHECKSET,
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
					"rotation_interval":  acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Disabled",
						"rotation_interval":  acctest.REMOVEKEY,
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
					"deletion_protection_description": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection":             "Disabled",
						"deletion_protection_description": acctest.REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy": `{\"Statement\": [{\"Action\": [\"kms:*\"],\"Effect\": \"Allow\",\"Principal\": {\"RAM\": [\"acs:ram::` + "${data.alicloud_account.default.id}" + `:*\"]},\"Resource\": [\"*\"],\"Sid\": \"Key\"}],\"Version\": \"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy": acctest.CHECKSET,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudKmsKeyMap1)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsKey")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsKey%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudKmsKeyBasicDependence8856)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"}) // fixed to cn-hangzhou region
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"dkms_instance_id":                acctest.CHECKSET,
						"protection_level":                "SOFTWARE",
						"automatic_rotation":              "Enabled",
						"rotation_interval":               "605800s",
						"policy":                          acctest.CHECKSET,
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
					"deletion_protection_description": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection":             "Disabled",
						"deletion_protection_description": acctest.REMOVEKEY,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudKmsKeyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsKey")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsKey%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudKmsKeyBasicDependence8855)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
					"rotation_interval":  acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Disabled",
						"rotation_interval":  acctest.REMOVEKEY,
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
					"deletion_protection_description": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection":             "Disabled",
						"deletion_protection_description": acctest.REMOVEKEY,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudKmsKeyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsKey")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsKey%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudKmsKeyBasicDependence8855)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
					"deletion_protection_description": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection":             "Disabled",
						"deletion_protection_description": acctest.REMOVEKEY,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudKmsKeyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsKey")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsKey%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudKmsKeyBasicDependence8855)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
					"rotation_interval":  acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Disabled",
						"rotation_interval":  acctest.REMOVEKEY,
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
					"deletion_protection_description": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection":             "Disabled",
						"deletion_protection_description": acctest.REMOVEKEY,
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
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudKmsKeyMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return kms.NewKmsServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKmsKey")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sKmsKey%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudKmsKeyBasicDependence8855)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
					"deletion_protection_description": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection":             "Disabled",
						"deletion_protection_description": acctest.REMOVEKEY,
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
	"key_usage":           acctest.CHECKSET,
	"origin":              acctest.CHECKSET,
	"key_spec":            acctest.CHECKSET,
	"automatic_rotation":  acctest.CHECKSET,
	"status":              acctest.CHECKSET,
	"arn":                 acctest.CHECKSET,
	"primary_key_version": acctest.CHECKSET,
	"last_rotation_date":  acctest.CHECKSET,
	"creator":             acctest.CHECKSET,
	"creation_date":       acctest.CHECKSET,
	"key_state":           acctest.CHECKSET,
	"is_enabled":          acctest.CHECKSET,
}

var AliCloudKmsKeyMap1 = map[string]string{
	"key_usage":           acctest.CHECKSET,
	"origin":              acctest.CHECKSET,
	"key_spec":            acctest.CHECKSET,
	"automatic_rotation":  acctest.CHECKSET,
	"policy":              acctest.CHECKSET,
	"status":              acctest.CHECKSET,
	"arn":                 acctest.CHECKSET,
	"primary_key_version": acctest.CHECKSET,
	"last_rotation_date":  acctest.CHECKSET,
	"creator":             acctest.CHECKSET,
	"creation_date":       acctest.CHECKSET,
	"key_state":           acctest.CHECKSET,
	"is_enabled":          acctest.CHECKSET,
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

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = "cn-hangzhou-k"
	}

	resource "alicloud_kms_instance" "default" {
  		product_version = "3"
  		vpc_num         = "1"
  		key_num         = "1000"
  		secret_num      = "0"
  		spec            = "1000"
  		vpc_id          = data.alicloud_vpcs.default.ids.0
		payment_type    = "PayAsYouGo"
  		vswitch_ids = [
    		data.alicloud_vswitches.default.ids.0
  		]
  		zone_ids = [
    		"cn-hangzhou-k",
    		"cn-hangzhou-j"
  		]
  		timeouts {
    		delete = "60m"
  		}
	}
`, name)
}

// lintignore: R001
// Test Ecs SecurityGroup. <<< Resource test cases, automatically generated.
